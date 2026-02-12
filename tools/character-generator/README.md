# 二次元角色预生成系统

## 概述

本系统用于批量预生成游戏NPC角色图像，基于Stable Diffusion + LoRA技术实现角色一致性，支持游戏内动态使用。

---

## 1. SD/LoRA 配置方案

### 1.1 硬件配置

#### 方案A: GPU服务器（推荐生产环境）
```yaml
# 推荐配置
GPU: NVIDIA RTX 4090 / A100
VRAM: 24GB+
CPU: 16核+
RAM: 64GB+
Storage: 500GB+ SSD
```

#### 方案B: 云端API（适合开发/测试）
```yaml
# 可选API服务
- Stability AI API: https://platform.stability.ai/
- Replicate: https://replicate.com/
- Hugging Face Inference API: https://huggingface.co/inference-api
- RunPod Serverless: https://www.runpod.io/
```

### 1.2 软件环境

```bash
# 基础环境
Python 3.10+
CUDA 11.8+
cuDNN 8.6+

# 核心依赖
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu118
pip install diffusers transformers accelerate safetensors
pip install xformers  # 可选，加速推理
pip install omegaconf einops

# WebUI（可选）
# AUTOMATIC1111 Stable Diffusion WebUI
git clone https://github.com/AUTOMATIC1111/stable-diffusion-webui.git
```

### 1.3 模型配置

```python
# config/models.py

MODEL_CONFIG = {
    # 基础动漫模型
    "base_models": {
        "anything_v5": {
            "path": "models/AnythingV5_v5.safetensors",
            "type": "sd15",
            "description": "通用动漫风格"
        },
        "counterfeit_xl": {
            "path": "models/CounterfeitXL_v10.safetensors",
            "type": "sdxl",
            "description": "高质量动漫SDXL模型"
        },
        "illustrious": {
            "path": "models/illustrious_v2.safetensors",
            "type": "sdxl",
            "description": "高质量动漫插图模型"
        }
    },

    # LoRA模型（角色一致性）
    "lora_models": {
        "character_maker": {
            "path": "loras/character_sheet_maker.safetensors",
            "weight": 0.7,
            "trigger_word": "character sheet"
        },
        "anime_detailer": {
            "path": "loras/anime_detailer.safetensors",
            "weight": 0.5
        }
    },

    # VAE（可选）
    "vae": {
        "path": "vae/anime.vae.safetensors"
    },

    # ControlNet
    "controlnet": {
        "openpose": {
            "path": "controlnet/control_v11p_sd15_openpose.safetensors"
        },
        "canny": {
            "path": "controlnet/control_v11p_sd15_canny.safetensors"
        }
    }
}
```

---

## 2. 批量生成流程设计

### 2.1 整体流程

```
┌─────────────────────────────────────────────────────────────────┐
│                     角色预生成流程                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. 角色定义                                                    │
│     ├── 读取NPC_PROFILES.md                                     │
│     └── 生成角色属性JSON                                        │
│              │                                                  │
│              ▼                                                  │
│  2. 提示词生成                                                  │
│     ├── 基于属性生成基础提示词                                   │
│     ├── 添加风格修饰词                                          │
│     └── 添加角色设定图关键词                                    │
│              │                                                  │
│              ▼                                                  │
│  3. 角色设定图生成                                              │
│     ├── 生成Character Sheet (多角度)                            │
│     ├── ControlNet控制姿势                                      │
│     └── 质量筛选                                                │
│              │                                                  │
│              ▼                                                  │
│  4. LoRA训练（可选）                                            │
│     ├── 构建训练数据集                                          │
│     ├── 训练角色LoRA                                            │
│     └── 验证生成效果                                            │
│              │                                                  │
│              ▼                                                  │
│  5. 批量变体生成                                                │
│     ├── 不同表情                                                 │
│     ├── 不同姿势                                                │
│     └── 不同服装/场景                                           │
│              │                                                  │
│              ▼                                                  │
│  6. 后处理与导出                                                │
│     ├── 图像裁剪/缩放                                           │
│     ├── 背景处理                                                │
│     ├── 元数据写入                                              │
│     └── 打包输出                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 生成脚本框架

```python
# scripts/generate_characters.py

import asyncio
import json
from pathlib import Path
from typing import List, Dict, Optional
from dataclasses import dataclass
from enum import Enum

# ============ 数据结构定义 ============

class AgeGroup(Enum):
    YOUNG = "young"
    MIDDLE = "middle"
    SENIOR = "senior"

class CharacterRole(Enum):
    MAYOR = "mayor"
    SHOPKEEPER = "shopkeeper"
    CARPENTER = "carpenter"
    FISHERMAN = "fisherman"
    VILLAGER = "villager"

@dataclass
class CharacterAttributes:
    """角色属性定义 - 基于ui-designer定义的属性"""
    id: str
    name: str
    role: CharacterRole
    age: AgeGroup
    gender: str
    hair_color: str
    hair_style: str
    eye_color: str
    body_type: str
    clothing_style: str
    accessories: List[str]
    personality_traits: List[str]

    # 扩展属性
    skin_tone: str = "fair"
    height: str = "average"
    distinctive_features: List[str] = None

    def to_prompt_dict(self) -> Dict:
        """转换为提示词字典"""
        return {
            "base": self._generate_base_prompt(),
            "appearance": self._generate_appearance_prompt(),
            "clothing": self._generate_clothing_prompt(),
            "personality": self._generate_personality_prompt()
        }

    def _generate_base_prompt(self) -> str:
        age_map = {
            AgeGroup.YOUNG: "young adult",
            AgeGroup.MIDDLE: "middle-aged",
            AgeGroup.SENIOR: "elderly"
        }
        return f"{age_map[self.age]} {self.gender}"

    def _generate_appearance_prompt(self) -> str:
        parts = [
            f"{self.hair_color} {self.hair_style} hair",
            f"{self.eye_color} eyes",
            f"{self.skin_tone} skin",
            f"{self.body_type} body"
        ]
        if self.distinctive_features:
            parts.extend(self.distinctive_features)
        return ", ".join(parts)

    def _generate_clothing_prompt(self) -> str:
        parts = [self.clothing_style]
        parts.extend(self.accessories)
        return ", ".join(parts)

    def _generate_personality_prompt(self) -> str:
        return ", ".join(self.personality_traits[:3])


@dataclass
class GenerationConfig:
    """生成配置"""
    model: str = "anything_v5"
    lora: Optional[str] = None
    lora_weight: float = 0.7

    # 图像尺寸
    width: int = 1024
    height: int = 1024

    # 生成参数
    steps: int = 28
    cfg_scale: float = 7.0
    sampler: str = "DPM++ 2M Karras"

    # 批量生成
    batch_size: int = 4
    num_batches: int = 2

    # 输出
    output_dir: str = "output/characters"
    save_metadata: bool = True


# ============ 提示词构建器 ============

class PromptBuilder:
    """提示词构建器"""

    # 通用质量词
    QUALITY_PROMPTS = [
        "masterpiece",
        "best quality",
        "highly detailed",
        "ultra detailed",
        "8k resolution",
        "sharp focus"
    ]

    # 动漫风格词
    ANIME_STYLE_PROMPTS = [
        "anime style",
        "cel shaded",
        "anime artwork",
        "studio ghibli inspired"
    ]

    # 负面提示词
    NEGATIVE_PROMPTS = [
        "low quality",
        "worst quality",
        "bad anatomy",
        "bad hands",
        "missing fingers",
        "extra digits",
        "cropped",
        "watermark",
        "signature",
        "blurry",
        "mutated",
        "deformed"
    ]

    # 角色设定图模板
    CHARACTER_SHEET_TEMPLATES = {
        "full_body": "character sheet, full body, multiple angles, front view, side view, back view, white background",
        "portrait": "portrait, head shot, facing viewer, simple background",
        "expressions": "expression sheet, multiple expressions, happy, sad, angry, surprised, neutral",
        "poses": "pose sheet, multiple poses, standing, sitting, walking, running"
    }

    @classmethod
    def build_character_prompt(
        cls,
        character: CharacterAttributes,
        template: str = "full_body",
        add_quality: bool = True,
        add_style: bool = True
    ) -> tuple[str, str]:
        """构建角色生成提示词"""

        prompt_parts = []

        # 1. 角色设定图类型
        prompt_parts.append(cls.CHARACTER_SHEET_TEMPLATES[template])

        # 2. 角色基础属性
        attr_prompts = character.to_prompt_dict()
        prompt_parts.append(attr_prompts["base"])
        prompt_parts.append(attr_prompts["appearance"])

        # 3. 服装
        prompt_parts.append(attr_prompts["clothing"])

        # 4. 质量词
        if add_quality:
            prompt_parts.extend(cls.QUALITY_PROMPTS[:4])

        # 5. 风格词
        if add_style:
            prompt_parts.extend(cls.ANIME_STYLE_PROMPTS[:2])

        # 构建最终提示词
        positive_prompt = ", ".join(prompt_parts)
        negative_prompt = ", ".join(cls.NEGATIVE_PROMPTS)

        return positive_prompt, negative_prompt


# ============ 图像生成器 ============

class CharacterGenerator:
    """角色图像生成器"""

    def __init__(self, config: GenerationConfig):
        self.config = config
        self.pipeline = None

    def load_model(self):
        """加载模型"""
        # 使用diffusers加载模型
        from diffusers import StableDiffusionXLPipeline
        import torch

        self.pipeline = StableDiffusionXLPipeline.from_pretrained(
            self.config.model,
            torch_dtype=torch.float16,
            use_safetensors=True
        )
        self.pipeline.to("cuda")

        # 启用内存优化
        self.pipeline.enable_attention_slicing()
        self.pipeline.enable_vae_slicing()

    async def generate_character(
        self,
        character: CharacterAttributes,
        template: str = "full_body",
        variations: int = 1
    ) -> List[Dict]:
        """生成角色图像"""

        results = []
        prompt, negative_prompt = PromptBuilder.build_character_prompt(
            character, template
        )

        for i in range(variations):
            # 生成图像
            images = self.pipeline(
                prompt=prompt,
                negative_prompt=negative_prompt,
                width=self.config.width,
                height=self.config.height,
                num_inference_steps=self.config.steps,
                guidance_scale=self.config.cfg_scale,
                num_images_per_prompt=self.config.batch_size
            ).images

            # 保存结果
            for j, img in enumerate(images):
                result = await self._save_image(
                    img, character, template, i * self.config.batch_size + j
                )
                results.append(result)

        return results

    async def _save_image(
        self,
        image,
        character: CharacterAttributes,
        template: str,
        index: int
    ) -> Dict:
        """保存图像和元数据"""
        import hashlib
        from datetime import datetime

        # 生成文件名
        filename = f"{character.id}_{template}_{index:03d}"
        output_path = Path(self.config.output_dir) / character.id
        output_path.mkdir(parents=True, exist_ok=True)

        # 保存图像
        image_path = output_path / f"{filename}.png"
        image.save(image_path, quality=95)

        # 生成图像哈希
        img_hash = hashlib.md5(image.tobytes()).hexdigest()

        # 构建元数据
        metadata = {
            "character_id": character.id,
            "character_name": character.name,
            "template": template,
            "index": index,
            "hash": img_hash,
            "generated_at": datetime.now().isoformat(),
            "config": {
                "model": self.config.model,
                "steps": self.config.steps,
                "cfg_scale": self.config.cfg_scale,
                "size": f"{self.config.width}x{self.config.height}"
            },
            "attributes": {
                "role": character.role.value,
                "age": character.age.value,
                "gender": character.gender,
                "hair_color": character.hair_color,
                "eye_color": character.eye_color
            }
        }

        # 保存元数据
        if self.config.save_metadata:
            metadata_path = output_path / f"{filename}.json"
            with open(metadata_path, 'w', encoding='utf-8') as f:
                json.dump(metadata, f, indent=2, ensure_ascii=False)

        return {
            "image_path": str(image_path),
            "metadata_path": str(metadata_path) if self.config.save_metadata else None,
            "metadata": metadata
        }


# ============ 批量生成器 ============

class BatchGenerator:
    """批量生成管理器"""

    def __init__(self, config: GenerationConfig):
        self.config = config
        self.generator = CharacterGenerator(config)
        self.characters: List[CharacterAttributes] = []

    def load_characters_from_profiles(self, profiles_path: str):
        """从NPC档案加载角色"""
        # 这里解析NPC_PROFILES.md或从数据库加载
        # 示例数据
        self.characters = [
            CharacterAttributes(
                id="lewis",
                name="Lewis",
                role=CharacterRole.MAYOR,
                age=AgeGroup.SENIOR,
                gender="male",
                hair_color="gray",
                hair_style="short receding",
                eye_color="brown",
                body_type="average",
                clothing_style="formal suit, mayoral sash",
                accessories=["glasses", "pocket watch"],
                personality_traits=["friendly", "responsible", "diplomatic"]
            ),
            CharacterAttributes(
                id="pierre",
                name="Pierre",
                role=CharacterRole.SHOPKEEPER,
                age=AgeGroup.MIDDLE,
                gender="male",
                hair_color="brown",
                hair_style="short neat",
                eye_color="hazel",
                body_type="slightly plump",
                clothing_style="casual shirt, apron",
                accessories=["spectacles"],
                personality_traits=["business_minded", "friendly", "ambitious"]
            ),
            CharacterAttributes(
                id="robin",
                name="Robin",
                role=CharacterRole.CARPENTER,
                age=AgeGroup.MIDDLE,
                gender="female",
                hair_color="orange",
                hair_style="short wavy",
                eye_color="green",
                body_type="athletic",
                clothing_style="work clothes, tool belt",
                accessories=["safety goggles"],
                personality_traits=["hardworking", "creative", "confident"]
            ),
            CharacterAttributes(
                id="haley",
                name="Haley",
                role=CharacterRole.VILLAGER,
                age=AgeGroup.YOUNG,
                gender="female",
                hair_color="blonde",
                hair_style="long wavy",
                eye_color="blue",
                body_type="slender",
                clothing_style="casual summer dress, flower accessories",
                accessories=["hair ribbon", "camera"],
                personality_traits=["beautiful", "photography_lover", "gradually_warming"]
            ),
            CharacterAttributes(
                id="willy",
                name="Willy",
                role=CharacterRole.FISHERMAN,
                age=AgeGroup.SENIOR,
                gender="male",
                hair_color="white",
                hair_style="balding with beard",
                eye_color="blue",
                body_type="weathered",
                clothing_style="fisherman outfit, rain coat",
                accessories=["fishing hat", "pipe"],
                personality_traits=["friendly", "wise", "patient"]
            )
        ]

    async def generate_all(
        self,
        templates: List[str] = None,
        variations_per_template: int = 2
    ):
        """批量生成所有角色"""
        if templates is None:
            templates = ["full_body", "portrait", "expressions"]

        self.generator.load_model()

        results = {}
        for character in self.characters:
            print(f"Generating character: {character.name}")
            results[character.id] = {}

            for template in templates:
                print(f"  Template: {template}")
                images = await self.generator.generate_character(
                    character,
                    template=template,
                    variations=variations_per_template
                )
                results[character.id][template] = images

        return results


# ============ 主入口 ============

async def main():
    config = GenerationConfig(
        model="models/AnythingV5_v5.safetensors",
        output_dir="output/characters",
        batch_size=4,
        steps=28
    )

    batch_gen = BatchGenerator(config)
    batch_gen.load_characters_from_profiles("doc/NPC_PROFILES.md")

    results = await batch_gen.generate_all(
        templates=["full_body", "portrait"],
        variations_per_template=2
    )

    print(f"Generated {sum(len(v) for r in results.values() for v in r.values())} images")


if __name__ == "__main__":
    asyncio.run(main())
```

---

## 3. 角色数据集构建策略

### 3.1 数据集结构

```
datasets/
├── characters/
│   ├── lewis/
│   │   ├── images/
│   │   │   ├── 001_full_body.png
│   │   │   ├── 002_full_body.png
│   │   │   ├── 003_portrait.png
│   │   │   └── ...
│   │   ├── metadata/
│   │   │   ├── 001_full_body.json
│   │   │   └── ...
│   │   └── character.json       # 角色总配置
│   ├── pierre/
│   │   └── ...
│   └── ...
├── prompts/
│   ├── base_prompts.yaml        # 基础提示词模板
│   ├── style_prompts.yaml       # 风格修饰词
│   └── negative_prompts.yaml    # 负面提示词
└── config/
    ├── generation_config.yaml   # 生成配置
    └── training_config.yaml     # LoRA训练配置
```

### 3.2 角色属性映射

```python
# scripts/attribute_mapper.py

from typing import Dict, List

class AttributeMapper:
    """将NPC档案属性映射为图像生成属性"""

    # 角色类型到服装的映射
    ROLE_CLOTHING_MAP = {
        "mayor": {
            "clothing": ["formal suit", "mayoral sash", "dress shoes"],
            "accessories": ["glasses", "pocket watch", "tie"],
            "setting": "office, town hall"
        },
        "shopkeeper": {
            "clothing": ["casual shirt", "apron", "comfortable shoes"],
            "accessories": ["spectacles", "name tag"],
            "setting": "shop interior, store counter"
        },
        "carpenter": {
            "clothing": ["work clothes", "tool belt", "work boots"],
            "accessories": ["safety goggles", "tool belt"],
            "setting": "workshop, construction site"
        },
        "fisherman": {
            "clothing": ["fisherman outfit", "rain coat", "rubber boots"],
            "accessories": ["fishing hat", "pipe", "fishing rod"],
            "setting": "dock, beach, boat"
        },
        "villager": {
            "clothing": ["casual clothes", "seasonal outfit"],
            "accessories": [],
            "setting": "town, home, outdoors"
        }
    }

    # 年龄到外观的映射
    AGE_APPEARANCE_MAP = {
        "young": {
            "skin": "smooth youthful skin",
            "features": "energetic expression",
            "body": "slim athletic build"
        },
        "middle": {
            "skin": "mature skin",
            "features": "confident expression",
            "body": "average build"
        },
        "senior": {
            "skin": "weathered skin with wrinkles",
            "features": "wise expression",
            "body": "slightly frail build"
        }
    }

    # 性格特征到表情的映射
    PERSONALITY_EXPRESSION_MAP = {
        "friendly": "warm smile, kind eyes",
        "responsible": "serious but approachable",
        "nostalgic": "wistful expression",
        "diplomatic": "calm composed expression",
        "business_minded": "confident professional look",
        "hardworking": "determined focused expression",
        "creative": "thoughtful inspired look",
        "beautiful": "elegant graceful expression",
        "wise": "knowing gentle smile",
        "patient": "calm serene expression"
    }

    @classmethod
    def map_npc_to_attributes(cls, npc_profile: Dict) -> Dict:
        """将NPC档案转换为图像生成属性"""

        result = {
            "id": npc_profile["id"],
            "name": npc_profile["name"],
            "role": npc_profile["role"],
            "age": npc_profile["age"],
            "traits": npc_profile["traits"]
        }

        # 映射角色类型相关属性
        role_config = cls.ROLE_CLOTHING_MAP.get(npc_profile["role"], {})
        result["clothing"] = role_config.get("clothing", [])
        result["accessories"] = role_config.get("accessories", [])
        result["setting"] = role_config.get("setting", "town")

        # 映射年龄相关属性
        age_config = cls.AGE_APPEARANCE_MAP.get(npc_profile["age"], {})
        result["skin"] = age_config.get("skin", "normal skin")
        result["body"] = age_config.get("body", "average build")

        # 映射性格特征
        expressions = []
        for trait in npc_profile.get("traits", []):
            if trait in cls.PERSONALITY_EXPRESSION_MAP:
                expressions.append(cls.PERSONALITY_EXPRESSION_MAP[trait])
        result["expressions"] = expressions[:3]  # 取前3个

        # 从对话风格推断更多信息
        speech_style = npc_profile.get("speech_style", "")
        result["speech_hints"] = cls._extract_speech_hints(speech_style)

        return result

    @classmethod
    def _extract_speech_hints(cls, speech_style: str) -> List[str]:
        """从对话风格提取视觉提示"""
        hints = []
        if "热情" in speech_style:
            hints.append("warm approachable demeanor")
        if "直接" in speech_style:
            hints.append("straightforward confident posture")
        if "傲慢" in speech_style:
            hints.append("proud haughty posture")
        if "温和" in speech_style:
            hints.append("gentle soft demeanor")
        return hints
```

### 3.3 LoRA训练数据集准备

```python
# scripts/prepare_lora_dataset.py

from pathlib import Path
import json
import shutil

class LoRADatasetPreparer:
    """LoRA训练数据集准备器"""

    def __init__(self, output_dir: str = "datasets/lora_training"):
        self.output_dir = Path(output_dir)

    def prepare_character_dataset(
        self,
        character_id: str,
        source_dir: str,
        num_images: int = 30,
        image_size: int = 512
    ):
        """准备单个角色的LoRA训练数据集"""

        char_output = self.output_dir / character_id
        char_output.mkdir(parents=True, exist_ok=True)

        # 1. 收集并筛选图像
        source_path = Path(source_dir) / character_id / "images"
        images = list(source_path.glob("*.png"))[:num_images]

        # 2. 复制并重命名图像
        for i, img_path in enumerate(images):
            dest = char_output / f"{i:04d}.jpg"
            self._process_image(img_path, dest, image_size)

        # 3. 生成训练标签文件
        self._generate_training_tags(character_id, char_output, len(images))

        # 4. 生成配置文件
        self._generate_dataset_config(character_id, char_output, len(images))

    def _process_image(self, source, dest, size):
        """处理图像：调整大小、格式转换"""
        from PIL import Image

        img = Image.open(source)
        img = img.convert("RGB")
        img = img.resize((size, size), Image.LANCZOS)
        img.save(dest, "JPEG", quality=95)

    def _generate_training_tags(self, character_id, output_dir, num_images):
        """生成训练标签文件"""
        # 读取角色元数据
        metadata_path = Path(f"datasets/characters/{character_id}/character.json")
        with open(metadata_path, 'r') as f:
            char_data = json.load(f)

        # 构建触发词
        trigger_words = [character_id, char_data.get("name", "")]
        trigger_words.extend(char_data.get("traits", [])[:2])

        # 为每张图像生成标签文件
        for i in range(num_images):
            tag_file = output_dir / f"{i:04d}.txt"
            tags = ", ".join(trigger_words)
            tag_file.write_text(tags)

    def _generate_dataset_config(self, character_id, output_dir, num_images):
        """生成数据集配置"""
        config = {
            "character_id": character_id,
            "num_images": num_images,
            "image_size": 512,
            "trigger_words": [character_id],
            "training_params": {
                "learning_rate": 1e-4,
                "train_batch_size": 4,
                "num_train_epochs": 10,
                "save_every_n_epochs": 2
            }
        }

        config_path = output_dir / "dataset_config.json"
        with open(config_path, 'w') as f:
            json.dump(config, f, indent=2)
```

---

## 4. 输出格式规范

### 4.1 图片格式规范

```yaml
# 图片输出规范
image_formats:
  # 主图像（高质量）
  main:
    format: PNG
    color_mode: RGBA
    resolution:
      - 1024x1024  # 角色设定图
      - 512x512    # 肖像
      - 256x256    # 图标
    quality: lossless

  # 游戏资源
  game_ready:
    format: PNG
    color_mode: RGBA
    resolution:
      - 128x128    # 小头像
      - 64x64      # 地图图标
      - 32x32      # 迷你图标
    background: transparent

  # Web优化
  web_optimized:
    format: WebP
    quality: 85
    resolution: 512x512

# 命名规范
naming_convention:
  pattern: "{character_id}_{type}_{variant}_{index}.{ext}"
  examples:
    - "lewis_full_body_v1_001.png"
    - "haley_portrait_v2_003.webp"
    - "pierre_expression_happy_001.png"
```

### 4.2 元数据格式

```json
{
  "schema_version": "1.0",
  "character": {
    "id": "lewis",
    "name": "Lewis",
    "display_name": "路易斯",
    "role": "mayor",
    "age_group": "senior"
  },
  "image": {
    "file_name": "lewis_full_body_v1_001.png",
    "file_path": "characters/lewis/images/lewis_full_body_v1_001.png",
    "file_size": 2048576,
    "dimensions": {
      "width": 1024,
      "height": 1024
    },
    "format": "PNG",
    "hash": {
      "md5": "a1b2c3d4e5f6...",
      "sha256": "abc123..."
    }
  },
  "generation": {
    "model": "anything_v5",
    "lora": null,
    "prompt": "character sheet, full body...",
    "negative_prompt": "low quality, worst quality...",
    "seed": 1234567890,
    "steps": 28,
    "cfg_scale": 7.0,
    "sampler": "DPM++ 2M Karras",
    "generated_at": "2026-02-12T10:30:00Z"
  },
  "attributes": {
    "gender": "male",
    "hair_color": "gray",
    "hair_style": "short receding",
    "eye_color": "brown",
    "skin_tone": "fair",
    "body_type": "average",
    "clothing": ["formal suit", "mayoral sash"],
    "accessories": ["glasses", "pocket watch"],
    "expression": "friendly"
  },
  "variants": {
    "type": "full_body",
    "variant": "v1",
    "index": 1,
    "total_variants": 4
  },
  "game_usage": {
    "sprite_type": "character_portrait",
    "animations": ["idle", "talk", "walk"],
    "seasons": ["spring", "summer", "fall", "winter"]
  }
}
```

### 4.3 游戏资源清单

```json
{
  "manifest_version": "1.0",
  "characters": [
    {
      "id": "lewis",
      "name": "Lewis",
      "assets": {
        "portrait": {
          "default": "lewis_portrait_v1_001.png",
          "happy": "lewis_expression_happy_001.png",
          "sad": "lewis_expression_sad_001.png",
          "angry": "lewis_expression_angry_001.png"
        },
        "sprite": {
          "idle": "lewis_sprite_idle.png",
          "walk": "lewis_sprite_walk.png"
        },
        "icon": {
          "map": "lewis_icon_64.png",
          "menu": "lewis_icon_128.png"
        }
      },
      "metadata": "lewis/character.json"
    }
  ],
  "total_characters": 5,
  "total_assets": 45,
  "generated_at": "2026-02-12T12:00:00Z"
}
```

---

## 5. 脚本目录结构设计

```
tools/character-generator/
├── README.md                          # 本文档
├── requirements.txt                   # Python依赖
├── setup.py                          # 安装脚本
│
├── config/                           # 配置文件
│   ├── models.yaml                   # 模型配置
│   ├── prompts/
│   │   ├── base_prompts.yaml        # 基础提示词
│   │   ├── style_prompts.yaml       # 风格提示词
│   │   └── negative_prompts.yaml    # 负面提示词
│   ├── generation.yaml              # 生成配置
│   └── training.yaml                # LoRA训练配置
│
├── scripts/                          # 脚本
│   ├── generate_characters.py       # 主生成脚本
│   ├── batch_generate.py            # 批量生成
│   ├── prepare_lora_dataset.py      # LoRA数据集准备
│   ├── train_lora.py               # LoRA训练脚本
│   ├── process_output.py           # 输出处理
│   └── export_game_assets.py       # 导出游戏资源
│
├── src/                             # 源代码
│   ├── __init__.py
│   ├── generator/
│   │   ├── __init__.py
│   │   ├── sd_generator.py         # SD生成器
│   │   ├── api_generator.py        # API生成器
│   │   └── controlnet.py           # ControlNet控制
│   ├── prompt/
│   │   ├── __init__.py
│   │   ├── builder.py              # 提示词构建
│   │   └── templates.py            # 提示词模板
│   ├── dataset/
│   │   ├── __init__.py
│   │   ├── prepare.py              # 数据集准备
│   │   └── augment.py              # 数据增强
│   ├── training/
│   │   ├── __init__.py
│   │   └── lora_trainer.py         # LoRA训练
│   └── utils/
│       ├── __init__.py
│       ├── image_utils.py          # 图像处理
│       ├── metadata.py             # 元数据处理
│       └── config.py               # 配置加载
│
├── models/                          # 模型文件（git忽略）
│   ├── base/                       # 基础模型
│   │   └── .gitkeep
│   ├── lora/                       # LoRA模型
│   │   └── .gitkeep
│   ├── vae/                        # VAE模型
│   │   └── .gitkeep
│   └── controlnet/                 # ControlNet模型
│       └── .gitkeep
│
├── datasets/                        # 数据集
│   ├── characters/                 # 角色图像数据
│   ├── prompts/                    # 提示词数据
│   └── lora_training/             # LoRA训练数据
│
├── output/                         # 输出目录
│   ├── characters/                # 生成的角色图像
│   ├── lora/                      # 训练的LoRA模型
│   └── game_assets/              # 游戏资源
│
├── tests/                          # 测试
│   ├── test_generator.py
│   ├── test_prompt_builder.py
│   └── test_dataset.py
│
└── examples/                       # 示例
    ├── example_config.yaml
    └── example_generation.py
```

---

## 6. 使用示例

### 6.1 快速开始

```bash
# 安装依赖
pip install -r requirements.txt

# 生成单个角色
python scripts/generate_characters.py --character lewis --template full_body

# 批量生成所有角色
python scripts/batch_generate.py --all --templates full_body,portrait,expressions

# 准备LoRA训练数据集
python scripts/prepare_lora_dataset.py --character lewis --num-images 30

# 训练角色LoRA
python scripts/train_lora.py --character lewis --epochs 10

# 导出游戏资源
python scripts/export_game_assets.py --output-dir ../frontend/public/assets/characters/
```

### 6.2 API调用示例

```python
from src.generator import CharacterGenerator
from src.prompt import PromptBuilder

# 初始化生成器
generator = CharacterGenerator(model="anything_v5")

# 定义角色属性
character = {
    "id": "lewis",
    "name": "Lewis",
    "role": "mayor",
    "age": "senior",
    "gender": "male",
    "hair_color": "gray",
    "eye_color": "brown"
}

# 生成提示词
prompt, negative = PromptBuilder.build_character_prompt(
    character,
    template="full_body"
)

# 生成图像
images = generator.generate(
    prompt=prompt,
    negative_prompt=negative,
    num_images=4
)

# 保存结果
for i, img in enumerate(images):
    img.save(f"output/lewis_{i}.png")
```

---

## 7. 后续扩展

### 7.1 待实现功能

1. **角色一致性增强**
   - IP-Adapter集成
   - 角色嵌入训练
   - 多角色一致性控制

2. **动态生成接口**
   - 实时生成API
   - 与游戏引擎集成
   - 运行时角色变体

3. **质量保证**
   - 自动质量评分
   - 异常检测
   - 人工审核工作流

4. **性能优化**
   - 批处理优化
   - 模型量化
   - 缓存策略

---

*文档版本: 1.0*
*创建者: Game Engine Developer*
*日期: 2026-02-12*
