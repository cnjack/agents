# 角色图像预生成方案设计

## 基于CHARACTER_GENERATOR_DESIGN.md的预生成技术方案

---

## 1. 组合数量计算

### 1.1 属性组合统计

| 属性类别 | 选项数 | 颜色数 | 组合数 |
|---------|-------|-------|--------|
| **发型** | 16 styles | 15 colors | **240** |
| **眼睛** | 8 shapes | 10 colors | **80** |
| **肤色** | 12 tones | - | **12** |
| **上衣** | 8 styles | 12 colors | **96** |
| **下装** | 6 styles | 6 colors | **36** |
| **鞋子** | 4 styles | - | **4** |
| **帽子** | 8 styles (含none) | - | **8** |
| **眼镜** | 6 styles (含none) | - | **6** |
| **性别** | 3 types | - | **3** |

### 1.2 分层生成策略

由于全组合数量巨大（240×80×12×96×36×...），采用**分层独立生成**策略：

```
Layer 1: 基础层 (性别 + 体型 + 肤色)
         3 × 9 × 12 = 324 种基础图

Layer 2: 发型层 (发型样式 + 发色)
         16 × 15 = 240 种发型图

Layer 3: 眼睛层 (眼型 + 瞳色)
         8 × 10 = 80 种眼睛图

Layer 4: 服装层
         上衣: 8 × 12 = 96 种
         下装: 6 × 6 = 36 种
         鞋子: 4 种

Layer 5: 配饰层
         帽子: 7 种 (不含none)
         眼镜: 5 种 (不含none)
         首饰: 组合选择
```

### 1.3 总计预生成数量

```
基础层:    324 张
发型层:    240 张
眼睛层:     80 张
上衣层:     96 张
下装层:     36 张
鞋子层:      4 张
帽子层:      7 张
眼镜层:      5 张
首饰层:     ~30 张 (组合)
─────────────────────
总计:     ~822 张独立素材
```

**游戏内动态合成**: 通过图层叠加实现 ~数百万种角色组合

---

## 2. 属性到提示词映射

### 2.1 发型映射

```python
# config/prompt_mappings/hair.py

HAIR_STYLE_PROMPTS = {
    "short_spiky": {
        "en": "short spiky hair, textured short hair",
        "zh": "短刺头",
        "tags": ["short hair", "spiky", "textured"]
    },
    "short_neat": {
        "en": "short neat hair, tidy short hair",
        "zh": "短发整齐",
        "tags": ["short hair", "neat", "tidy"]
    },
    "medium_wavy": {
        "en": "medium length wavy hair, shoulder length waves",
        "zh": "中长波浪",
        "tags": ["medium hair", "wavy", "shoulder length"]
    },
    "medium_straight": {
        "en": "medium length straight hair, shoulder length straight",
        "zh": "中长直发",
        "tags": ["medium hair", "straight", "shoulder length"]
    },
    "long_straight": {
        "en": "long straight hair, flowing straight hair",
        "zh": "长直发",
        "tags": ["long hair", "straight", "flowing"]
    },
    "long_wavy": {
        "en": "long wavy hair, flowing waves",
        "zh": "长波浪",
        "tags": ["long hair", "wavy", "flowing"]
    },
    "ponytail": {
        "en": "ponytail, hair tied back",
        "zh": "马尾辫",
        "tags": ["ponytail", "tied hair"]
    },
    "braid": {
        "en": "braided hair, plait",
        "zh": "麻花辫",
        "tags": ["braid", "plait", "long hair"]
    },
    "bald": {
        "en": "bald, no hair",
        "zh": "光头",
        "tags": ["bald", "no hair"]
    },
    "curly_short": {
        "en": "short curly hair, tight curls",
        "zh": "短卷发",
        "tags": ["short hair", "curly", "tight curls"]
    },
    "curly_long": {
        "en": "long curly hair, flowing curls",
        "zh": "长卷发",
        "tags": ["long hair", "curly", "flowing"]
    },
    "mohawk": {
        "en": "mohawk hairstyle, spiked crest",
        "zh": "莫霍克",
        "tags": ["mohawk", "spiked", "edgy"]
    },
    "buns": {
        "en": "double bun hairstyle, odango hair",
        "zh": "双丸子头",
        "tags": ["double bun", "odango", "tied hair"]
    },
    "undercut": {
        "en": "undercut hairstyle, shaved sides",
        "zh": "铲青",
        "tags": ["undercut", "shaved sides", "modern"]
    },
    "afro": {
        "en": "afro hairstyle, big curly hair",
        "zh": "爆炸头",
        "tags": ["afro", "big hair", "curly"]
    },
    "dreadlocks": {
        "en": "dreadlocks, dreads",
        "zh": "脏辫",
        "tags": ["dreadlocks", "dreads", "long hair"]
    }
}

HAIR_COLOR_PROMPTS = {
    "black":       {"hex": "#1a1a1a", "prompt": "black hair"},
    "brown_dark":  {"hex": "#3d2314", "prompt": "dark brown hair"},
    "brown_light": {"hex": "#8b4513", "prompt": "light brown hair"},
    "blonde":      {"hex": "#f4d03f", "prompt": "blonde hair, golden hair"},
    "platinum":    {"hex": "#e8e4c9", "prompt": "platinum blonde hair, silver hair"},
    "red":         {"hex": "#8b2500", "prompt": "red hair, dark red hair"},
    "ginger":      {"hex": "#d35400", "prompt": "ginger hair, orange hair"},
    "auburn":      {"hex": "#922b21", "prompt": "auburn hair, reddish brown hair"},
    "gray":        {"hex": "#808080", "prompt": "gray hair, grey hair"},
    "white":       {"hex": "#f5f5f5", "prompt": "white hair"},
    "blue":        {"hex": "#3498db", "prompt": "blue hair, anime blue hair"},
    "pink":        {"hex": "#ff69b4", "prompt": "pink hair, anime pink hair"},
    "purple":      {"hex": "#9b59b6", "prompt": "purple hair, anime purple hair"},
    "green":       {"hex": "#27ae60", "prompt": "green hair, anime green hair"},
    "teal":        {"hex": "#1abc9c", "prompt": "teal hair, cyan hair"}
}
```

### 2.2 眼睛映射

```python
# config/prompt_mappings/eyes.py

EYE_SHAPE_PROMPTS = {
    "round": {
        "en": "round eyes, large round eyes",
        "zh": "圆眼",
        "tags": ["round eyes", "large eyes", "cute"]
    },
    "almond": {
        "en": "almond shaped eyes",
        "zh": "杏眼",
        "tags": ["almond eyes", "elegant"]
    },
    "narrow": {
        "en": "narrow eyes, thin eyes",
        "zh": "细长眼",
        "tags": ["narrow eyes", "thin eyes", "sharp"]
    },
    "wide": {
        "en": "wide eyes, large expressive eyes",
        "zh": "大眼",
        "tags": ["wide eyes", "large eyes", "expressive"]
    },
    "droopy": {
        "en": "droopy eyes, downward slanting eyes",
        "zh": "下垂眼",
        "tags": ["droopy eyes", "gentle expression"]
    },
    "cat_eye": {
        "en": "cat eyes, upward slanting eyes",
        "zh": "猫眼",
        "tags": ["cat eyes", "upward slant", "sharp"]
    },
    "monolid": {
        "en": "monolid eyes, single eyelid",
        "zh": "单眼皮",
        "tags": ["monolid", "single eyelid"]
    },
    "hooded": {
        "en": "hooded eyes, partially hidden eyelid",
        "zh": "内双",
        "tags": ["hooded eyes", "subtle eyelid"]
    }
}

EYE_COLOR_PROMPTS = {
    "brown_dark":  {"hex": "#3d2314", "prompt": "dark brown eyes"},
    "brown_light": {"hex": "#8b4513", "prompt": "light brown eyes"},
    "hazel":       {"hex": "#a67c52", "prompt": "hazel eyes"},
    "amber":       {"hex": "#ffbf00", "prompt": "amber eyes, golden eyes"},
    "blue":        {"hex": "#3498db", "prompt": "blue eyes"},
    "blue_light":  {"hex": "#87ceeb", "prompt": "light blue eyes, sky blue eyes"},
    "green":       {"hex": "#27ae60", "prompt": "green eyes"},
    "gray":        {"hex": "#808080", "prompt": "gray eyes, grey eyes"},
    "violet":      {"hex": "#9b59b6", "prompt": "violet eyes, purple eyes"},
    "red":         {"hex": "#c0392b", "prompt": "red eyes, crimson eyes"}
}
```

### 2.3 服装映射

```python
# config/prompt_mappings/clothing.py

TOP_STYLE_PROMPTS = {
    "t_shirt":    {"en": "t-shirt, casual shirt", "tags": ["t-shirt", "casual"]},
    "shirt":      {"en": "button-up shirt, collared shirt", "tags": ["shirt", "collared"]},
    "sweater":    {"en": "sweater, knit sweater", "tags": ["sweater", "warm", "knit"]},
    "hoodie":     {"en": "hoodie, hooded sweatshirt", "tags": ["hoodie", "casual", "hooded"]},
    "overalls":   {"en": "overalls, dungarees", "tags": ["overalls", "workwear"]},
    "vest":       {"en": "vest, waistcoat", "tags": ["vest", "layered"]},
    "flannel":    {"en": "flannel shirt, plaid shirt", "tags": ["flannel", "plaid", "checkered"]},
    "dress":      {"en": "dress, one-piece dress", "tags": ["dress", "feminine"]}
}

BOTTOM_STYLE_PROMPTS = {
    "jeans":       {"en": "jeans, denim pants", "tags": ["jeans", "denim"]},
    "shorts":      {"en": "shorts, short pants", "tags": ["shorts", "casual"]},
    "cargo_pants": {"en": "cargo pants, utility pants", "tags": ["cargo pants", "utility"]},
    "skirt":       {"en": "skirt", "tags": ["skirt", "feminine"]},
    "overalls":    {"en": "overalls pants", "tags": ["overalls", "workwear"]},
    "sweatpants":  {"en": "sweatpants, jogging pants", "tags": ["sweatpants", "casual", "sporty"]}
}

SHOES_STYLE_PROMPTS = {
    "sneakers":    {"en": "sneakers, athletic shoes", "tags": ["sneakers", "casual"]},
    "boots":       {"en": "boots, ankle boots", "tags": ["boots", "sturdy"]},
    "sandals":     {"en": "sandals, open-toe shoes", "tags": ["sandals", "summer"]},
    "work_boots":  {"en": "work boots, heavy boots", "tags": ["work boots", "rugged"]}
}

# 颜色映射
CLOTHING_COLORS = {
    "white":      "#ffffff",
    "black":      "#1a1a1a",
    "red":        "#c0392b",
    "blue":       "#3498db",
    "green":      "#27ae60",
    "yellow":     "#f1c40f",
    "orange":     "#e67e22",
    "purple":     "#9b59b6",
    "pink":       "#ff69b4",
    "brown":      "#8b4513",
    "gray":       "#808080",
    "navy":       "#1a3a5c",
    "blue_denim": "#4a6fa5",
    "black_denim":"#2c3e50",
    "khaki":      "#c3b091"
}
```

### 2.4 其他属性映射

```python
# config/prompt_mappings/body.py

GENDER_PROMPTS = {
    "male":       {"prompt": "male, man, masculine", "sprite_suffix": "_m"},
    "female":     {"prompt": "female, woman, feminine", "sprite_suffix": "_f"},
    "non_binary": {"prompt": "androgynous, neutral", "sprite_suffix": "_nb"}
}

SKIN_TONE_PROMPTS = {
    "pale":           {"hex": "#ffe4c4", "prompt": "pale skin, fair skin"},
    "fair":           {"hex": "#f5deb3", "prompt": "fair skin, light skin"},
    "light":          {"hex": "#deb887", "prompt": "light skin"},
    "medium":         {"hex": "#d2a679", "prompt": "medium skin tone"},
    "tan":            {"hex": "#c68642", "prompt": "tan skin, sun-kissed"},
    "olive":          {"hex": "#b5a642", "prompt": "olive skin tone"},
    "brown":          {"hex": "#8d5524", "prompt": "brown skin"},
    "dark_brown":     {"hex": "#5c4033", "prompt": "dark brown skin"},
    "black":          {"hex": "#3d2314", "prompt": "black skin, dark skin"},
    "fantasy_blue":   {"hex": "#87ceeb", "prompt": "blue skin, fantasy skin"},
    "fantasy_green":  {"hex": "#90ee90", "prompt": "green skin, fantasy skin"},
    "fantasy_purple": {"hex": "#dda0dd", "prompt": "purple skin, fantasy skin"}
}

BODY_TYPE_PROMPTS = {
    "slim":    {"height": 0.85, "build": 0.8, "prompt": "slim body, slender build"},
    "average": {"height": 1.0, "build": 1.0, "prompt": "average body, normal build"},
    "broad":   {"height": 1.15, "build": 1.2, "prompt": "broad body, muscular build"}
}

# config/prompt_mappings/accessories.py

HAT_PROMPTS = {
    "none":          None,
    "baseball_cap":  {"prompt": "baseball cap", "tags": ["cap", "casual"]},
    "straw_hat":     {"prompt": "straw hat, sun hat", "tags": ["straw hat", "farmer"]},
    "beanie":        {"prompt": "beanie, knit cap", "tags": ["beanie", "winter"]},
    "cowboy_hat":    {"prompt": "cowboy hat, western hat", "tags": ["cowboy hat", "western"]},
    "sun_hat":       {"prompt": "sun hat, wide brim hat", "tags": ["sun hat", "summer"]},
    "bowler":        {"prompt": "bowler hat, derby hat", "tags": ["bowler hat", "formal"]},
    "witch_hat":     {"prompt": "witch hat, pointed hat", "tags": ["witch hat", "fantasy"]}
}

GLASSES_PROMPTS = {
    "none":            None,
    "round_glasses":   {"prompt": "round glasses", "tags": ["glasses", "round frames"]},
    "square_glasses":  {"prompt": "square glasses", "tags": ["glasses", "square frames"]},
    "sunglasses":      {"prompt": "sunglasses, shades", "tags": ["sunglasses", "cool"]},
    "reading_glasses": {"prompt": "reading glasses", "tags": ["glasses", "intellectual"]},
    "monocle":         {"prompt": "monocle", "tags": ["monocle", "fancy"]}
}
```

---

## 3. 文件命名规范

### 3.1 命名模式

```
{category}_{style_id}_{color_id}_{hex_hash}.{ext}
```

### 3.2 各类别命名示例

```
# 发型
hair_01_blue_3498db.png          # 短刺头 + 蓝色
hair_06_blonde_f4d03f.png        # 长波浪 + 金色
hair_00_none_none.png            # 光头

# 眼睛
eyes_01_blue_3498db.png          # 圆眼 + 蓝色
eyes_06_violet_9b59b6.png        # 猫眼 + 紫色

# 基础身体
body_m_light_deb887.png          # 男性 + 浅肤色
body_f_fair_f5deb3.png           # 女性 + 白皙肤色

# 服装
top_01_red_c0392b.png            # T恤 + 红色
bottom_01_blue_denim_4a6fa5.png  # 牛仔裤 + 蓝色牛仔
shoes_01_none_none.png           # 运动鞋

# 配饰
hat_02_none_none.png             # 草帽
glasses_03_none_none.png         # 墨镜
```

### 3.3 完整文件路径结构

```
assets/sprites/character/
├── base/
│   ├── body_m_light_deb887.png
│   ├── body_m_medium_d2a679.png
│   ├── body_f_fair_f5deb3.png
│   └── ...
├── hair/
│   ├── hair_01_black_1a1a1a.png
│   ├── hair_01_brown_dark_3d2314.png
│   ├── hair_01_blonde_f4d03f.png
│   ├── hair_02_black_1a1a1a.png
│   └── ...  (240 files)
├── eyes/
│   ├── eyes_01_brown_dark_3d2314.png
│   ├── eyes_01_blue_3498db.png
│   └── ...  (80 files)
├── clothing/
│   ├── top/
│   │   ├── top_01_white_ffffff.png
│   │   ├── top_01_red_c0392b.png
│   │   └── ...  (96 files)
│   ├── bottom/
│   │   ├── bottom_01_blue_denim_4a6fa5.png
│   │   └── ...  (36 files)
│   └── shoes/
│       ├── shoes_01_none_none.png
│       └── ...  (4 files)
└── accessories/
    ├── hat/
    │   ├── hat_01_none_none.png
    │   └── ...  (7 files)
    ├── glasses/
    │   └── ...  (5 files)
    └── jewelry/
        └── ...  (~30 files)
```

---

## 4. 元数据格式

### 4.1 单个素材元数据

```json
{
  "id": "hair_01_blue_3498db",
  "category": "hair",
  "file": {
    "path": "assets/sprites/character/hair/hair_01_blue_3498db.png",
    "format": "PNG",
    "size": {"width": 64, "height": 64},
    "file_size_bytes": 2048
  },
  "attributes": {
    "style_id": "hair_01",
    "style_name": {
      "en": "Short Spiky",
      "zh": "短刺头"
    },
    "color_id": "blue",
    "color_hex": "#3498db",
    "color_name": {
      "en": "Blue",
      "zh": "蓝色"
    }
  },
  "generation": {
    "model": "anything_v5",
    "prompt": "character sprite, short spiky hair, blue hair, anime style, pixel art, transparent background, simple shading",
    "negative_prompt": "complex, detailed background, shading, shadows, outline, watermark",
    "seed": 1234567890,
    "generated_at": "2026-02-12T10:00:00Z"
  },
  "tags": ["short hair", "spiky", "textured", "blue", "anime"],
  "compatible_with": ["body_m", "body_f", "body_nb"]
}
```

### 4.2 批量生成清单

```json
{
  "manifest_version": "1.0.0",
  "generated_at": "2026-02-12T12:00:00Z",
  "categories": {
    "hair": {
      "total": 240,
      "styles": 16,
      "colors": 15,
      "files": [
        "hair_01_black_1a1a1a.png",
        "hair_01_brown_dark_3d2314.png",
        "..."
      ]
    },
    "eyes": {
      "total": 80,
      "styles": 8,
      "colors": 10,
      "files": ["..."]
    },
    "body": {
      "total": 324,
      "genders": 3,
      "skin_tones": 12,
      "files": ["..."]
    },
    "clothing": {
      "top": {"total": 96},
      "bottom": {"total": 36},
      "shoes": {"total": 4}
    },
    "accessories": {
      "hat": {"total": 7},
      "glasses": {"total": 5},
      "jewelry": {"total": 30}
    }
  },
  "total_files": 822,
  "total_size_mb": 15.6
}
```

---

## 5. 分批生成计划

### 5.1 批次划分

| 批次 | 内容 | 数量 | 预估时间 | 优先级 |
|-----|------|-----|---------|-------|
| **Batch 1** | 基础身体 (性别×肤色) | 324 | 2h | P0 |
| **Batch 2** | 发型 (样式×颜色) | 240 | 1.5h | P0 |
| **Batch 3** | 眼睛 (眼型×颜色) | 80 | 30min | P1 |
| **Batch 4** | 上衣 (样式×颜色) | 96 | 45min | P1 |
| **Batch 5** | 下装+鞋子 | 40 | 20min | P1 |
| **Batch 6** | 配饰 | 42 | 20min | P2 |

### 5.2 生成脚本

```python
# scripts/batch_generator.py

from dataclasses import dataclass
from typing import List, Dict, Iterator
from itertools import product

@dataclass
class GenerationTask:
    """单个生成任务"""
    category: str
    style_id: str
    color_id: str
    color_hex: str
    output_filename: str
    prompt: str
    negative_prompt: str

class BatchPlanGenerator:
    """生成批次计划"""

    def __init__(self, config: dict):
        self.config = config
        self.mappings = self._load_mappings()

    def _load_mappings(self) -> dict:
        """加载提示词映射"""
        return {
            "hair": HAIR_STYLE_PROMPTS,
            "hair_color": HAIR_COLOR_PROMPTS,
            "eyes": EYE_SHAPE_PROMPTS,
            "eye_color": EYE_COLOR_PROMPTS,
            # ... 其他映射
        }

    def generate_batch_plan(self, batch_id: int) -> List[GenerationTask]:
        """生成指定批次的任务列表"""
        tasks = []

        if batch_id == 1:
            # 基础身体层
            for gender in ["male", "female", "non_binary"]:
                for skin_id, skin_data in SKIN_TONE_PROMPTS.items():
                    for body_type in ["slim", "average", "broad"]:
                        task = self._create_body_task(
                            gender, skin_id, skin_data, body_type
                        )
                        tasks.append(task)

        elif batch_id == 2:
            # 发型层
            for style_id, style_data in HAIR_STYLE_PROMPTS.items():
                for color_id, color_data in HAIR_COLOR_PROMPTS.items():
                    task = self._create_hair_task(style_id, style_data, color_id, color_data)
                    tasks.append(task)

        elif batch_id == 3:
            # 眼睛层
            for shape_id, shape_data in EYE_SHAPE_PROMPTS.items():
                for color_id, color_data in EYE_COLOR_PROMPTS.items():
                    task = self._create_eyes_task(shape_id, shape_data, color_id, color_data)
                    tasks.append(task)

        # ... 其他批次

        return tasks

    def _create_hair_task(self, style_id, style_data, color_id, color_data) -> GenerationTask:
        """创建发型生成任务"""
        filename = f"hair_{style_id.split('_')[1]}_{color_id}_{color_data['hex'][1:]}.png"

        prompt = (
            f"character sprite, {style_data['en']}, {color_data['prompt']}, "
            f"anime style, pixel art, transparent background, simple shading, "
            f"isolated on transparent background, no body"
        )

        negative = (
            "complex background, body, face, detailed shading, "
            "watermark, signature, outline, shadow"
        )

        return GenerationTask(
            category="hair",
            style_id=style_id,
            color_id=color_id,
            color_hex=color_data["hex"],
            output_filename=filename,
            prompt=prompt,
            negative_prompt=negative
        )

    def get_total_tasks(self) -> Dict[str, int]:
        """获取各批次任务总数"""
        return {
            "batch_1_body": 3 * 12 * 9,      # 324
            "batch_2_hair": 16 * 15,          # 240
            "batch_3_eyes": 8 * 10,           # 80
            "batch_4_top": 8 * 12,            # 96
            "batch_5_bottom": 6 * 6 + 4,      # 40
            "batch_6_accessories": 7 + 5 + 30 # 42
        }


# 使用示例
def main():
    planner = BatchPlanGenerator(config={})

    # 生成批次1的任务
    batch1_tasks = planner.generate_batch_plan(1)
    print(f"Batch 1: {len(batch1_tasks)} tasks")

    # 执行生成
    generator = CharacterGenerator(config)
    for task in batch1_tasks:
        result = await generator.generate_asset(task)
        save_result(result)
```

---

## 6. 提示词模板

### 6.1 分层提示词模板

```yaml
# config/prompt_templates.yaml

templates:
  # 基础身体层模板
  body_base:
    positive: |
      anime character sprite, {gender} body, {skin_tone} skin, {body_type} build,
      standing pose, front view, pixel art style, 64x64,
      simple shading, transparent background, no clothes, bald
    negative: |
      clothes, hair, eyes, detailed shading, complex background,
      watermark, signature, outline

  # 发型层模板
  hair_layer:
    positive: |
      anime hair sprite, {hair_style}, {hair_color} hair,
      pixel art style, transparent background, isolated hair only,
      simple shading, no face, no body
    negative: |
      face, body, eyes, complex shading, background,
      watermark, signature

  # 眼睛层模板
  eyes_layer:
    positive: |
      anime eyes sprite, {eye_shape}, {eye_color} eyes,
      pixel art style, transparent background, isolated eyes only,
      simple shading, no face, no hair
    negative: |
      face, hair, body, complex shading, background,
      watermark, signature

  # 服装层模板
  clothing_top:
    positive: |
      anime clothing sprite, {top_style}, {top_color} color,
      pixel art style, transparent background, isolated clothing only,
      simple shading, no body, no face
    negative: |
      body, face, hair, complex shading, background

  clothing_bottom:
    positive: |
      anime pants sprite, {bottom_style}, {bottom_color} color,
      pixel art style, transparent background, isolated clothing only,
      simple shading, no body, no face
    negative: |
      body, face, hair, complex shading, background

  # 配饰层模板
  accessory_hat:
    positive: |
      anime hat sprite, {hat_style},
      pixel art style, transparent background, isolated accessory only,
      simple shading, no head
    negative: |
      head, face, hair, body, complex shading, background
```

---

## 7. 执行命令

```bash
# 初始化项目
python scripts/setup.py

# 生成批次1（基础身体）
python scripts/generate_batch.py --batch 1 --output assets/sprites/character/base/

# 生成批次2（发型）
python scripts/generate_batch.py --batch 2 --output assets/sprites/character/hair/

# 生成所有批次
python scripts/generate_all.py --parallel 4

# 验证生成结果
python scripts/validate_assets.py --manifest assets/manifest.json

# 导出游戏资源
python scripts/export_game_assets.py --format webp --quality 90
```

---

*文档版本: 1.0*
*创建者: Game Engine Developer*
*日期: 2026-02-12*
