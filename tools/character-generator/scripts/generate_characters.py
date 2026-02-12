#!/usr/bin/env python3
"""
Character Generator - Main script for generating anime OC characters
Usage:
    python generate_characters.py --character lewis --template full_body
    python generate_characters.py --all --templates full_body,portrait
"""

import argparse
import asyncio
import json
import hashlib
from pathlib import Path
from datetime import datetime
from dataclasses import dataclass, field, asdict
from typing import List, Dict, Optional, Tuple
from enum import Enum
import yaml


# ============ Enums ============

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


class TemplateType(Enum):
    FULL_BODY = "full_body"
    PORTRAIT = "portrait"
    EXPRESSIONS = "expressions"
    POSES = "poses"
    TURNAROUND = "turnaround"


# ============ Data Classes ============

@dataclass
class CharacterAttributes:
    """Character attributes for image generation"""
    id: str
    name: str
    role: str
    age: str
    gender: str
    hair_color: str
    hair_style: str
    eye_color: str
    body_type: str = "average"
    clothing_style: str = ""
    accessories: List[str] = field(default_factory=list)
    personality_traits: List[str] = field(default_factory=list)
    skin_tone: str = "fair"
    distinctive_features: List[str] = field(default_factory=list)


@dataclass
class GenerationConfig:
    """Generation configuration"""
    model: str = "anything_v5"
    lora: Optional[str] = None
    lora_weight: float = 0.7
    width: int = 1024
    height: int = 1024
    steps: int = 28
    cfg_scale: float = 7.0
    sampler: str = "DPM++ 2M Karras"
    batch_size: int = 4
    output_dir: str = "output/characters"
    save_metadata: bool = True
    seed: int = -1


# ============ Prompt Builder ============

class PromptBuilder:
    """Build prompts for character generation"""

    QUALITY_PROMPTS = [
        "masterpiece", "best quality", "highly detailed",
        "ultra detailed", "8k resolution", "sharp focus"
    ]

    ANIME_STYLE_PROMPTS = ["anime style", "cel shaded", "anime artwork"]

    NEGATIVE_PROMPTS = [
        "low quality", "worst quality", "bad anatomy", "bad hands",
        "missing fingers", "extra digits", "cropped", "watermark",
        "signature", "blurry", "mutated", "deformed"
    ]

    CHARACTER_SHEET_TEMPLATES = {
        "full_body": "character sheet, full body, multiple angles, front view, side view, back view, white background",
        "portrait": "portrait, head shot, facing viewer, simple background",
        "expressions": "expression sheet, multiple expressions, happy, sad, angry, surprised, neutral",
        "poses": "pose sheet, multiple poses, standing, sitting, walking, running",
        "turnaround": "character turnaround, front view, side view, back view, three-quarter view"
    }

    ROLE_CLOTHING_MAP = {
        "mayor": ["formal suit", "mayoral sash", "dress shoes"],
        "shopkeeper": ["casual shirt", "apron", "comfortable shoes"],
        "carpenter": ["work clothes", "tool belt", "work boots"],
        "fisherman": ["fisherman outfit", "rain coat", "rubber boots"],
        "villager": ["casual clothes", "seasonal outfit"]
    }

    AGE_APPEARANCE_MAP = {
        "young": ["young adult", "youthful appearance", "smooth skin"],
        "middle": ["middle-aged", "mature appearance"],
        "senior": ["elderly", "aged appearance", "wrinkles"]
    }

    @classmethod
    def build_prompt(
        cls,
        character: CharacterAttributes,
        template: str = "full_body",
        add_quality: bool = True,
        add_style: bool = True
    ) -> Tuple[str, str]:
        """Build positive and negative prompts"""

        parts = []

        # 1. Template type
        if template in cls.CHARACTER_SHEET_TEMPLATES:
            parts.append(cls.CHARACTER_SHEET_TEMPLATES[template])

        # 2. Base attributes
        age_desc = cls.AGE_APPEARANCE_MAP.get(character.age, [""])[0]
        parts.append(f"{age_desc} {character.gender}")

        # 3. Appearance
        appearance_parts = [
            f"{character.hair_color} {character.hair_style} hair",
            f"{character.eye_color} eyes",
            f"{character.skin_tone} skin"
        ]
        if character.body_type != "average":
            appearance_parts.append(f"{character.body_type} body")
        parts.append(", ".join(appearance_parts))

        # 4. Clothing
        if character.clothing_style:
            parts.append(character.clothing_style)
        else:
            role_clothes = cls.ROLE_CLOTHING_MAP.get(character.role, [])
            if role_clothes:
                parts.append(", ".join(role_clothes[:2]))

        # 5. Accessories
        if character.accessories:
            parts.append(", ".join(character.accessories[:3]))

        # 6. Quality and style
        if add_quality:
            parts.extend(cls.QUALITY_PROMPTS[:4])
        if add_style:
            parts.extend(cls.ANIME_STYLE_PROMPTS[:2])

        positive = ", ".join(parts)
        negative = ", ".join(cls.NEGATIVE_PROMPTS)

        return positive, negative


# ============ Character Profiles ============

class CharacterProfileLoader:
    """Load character profiles from game data"""

    DEFAULT_CHARACTERS = [
        CharacterAttributes(
            id="lewis",
            name="Lewis",
            role="mayor",
            age="senior",
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
            role="shopkeeper",
            age="middle",
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
            role="carpenter",
            age="middle",
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
            role="villager",
            age="young",
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
            role="fisherman",
            age="senior",
            gender="male",
            hair_color="white",
            hair_style="balding with beard",
            eye_color="blue",
            body_type="weathered",
            clothing_style="fisherman outfit, rain coat",
            accessories=["fishing hat", "pipe"],
            personality_traits=["friendly", "wise", "patient"]
        ),
    ]

    @classmethod
    def load(cls, profile_path: Optional[str] = None) -> List[CharacterAttributes]:
        """Load character profiles from file or use defaults"""
        if profile_path and Path(profile_path).exists():
            with open(profile_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
                return [CharacterAttributes(**char) for char in data.get("characters", [])]
        return cls.DEFAULT_CHARACTERS

    @classmethod
    def get_by_id(cls, char_id: str) -> Optional[CharacterAttributes]:
        """Get a specific character by ID"""
        for char in cls.DEFAULT_CHARACTERS:
            if char.id == char_id:
                return char
        return None


# ============ Image Generator ============

class CharacterGenerator:
    """Generate character images using SD or API"""

    def __init__(self, config: GenerationConfig):
        self.config = config
        self.pipeline = None
        self.output_dir = Path(config.output_dir)

    def _init_pipeline(self):
        """Initialize the diffusion pipeline"""
        try:
            import torch
            from diffusers import StableDiffusionXLPipeline

            print(f"Loading model: {self.config.model}")
            self.pipeline = StableDiffusionXLPipeline.from_pretrained(
                self.config.model,
                torch_dtype=torch.float16,
                use_safetensors=True,
                variant="fp16"
            )
            self.pipeline.to("cuda")
            self.pipeline.enable_attention_slicing()
            print("Model loaded successfully")
        except ImportError:
            print("Warning: diffusers not installed. Using mock generation.")
            self.pipeline = None

    async def generate(
        self,
        character: CharacterAttributes,
        template: str = "full_body",
        variations: int = 1
    ) -> List[Dict]:
        """Generate character images"""
        if self.pipeline is None:
            self._init_pipeline()

        results = []
        positive, negative = PromptBuilder.build_prompt(character, template)

        print(f"\nGenerating {character.name} - {template}")
        print(f"Prompt: {positive[:100]}...")

        # Create output directory
        char_dir = self.output_dir / character.id
        char_dir.mkdir(parents=True, exist_ok=True)

        for v in range(variations):
            for i in range(self.config.batch_size):
                # Generate image
                if self.pipeline:
                    import torch
                    images = self.pipeline(
                        prompt=positive,
                        negative_prompt=negative,
                        width=self.config.width,
                        height=self.config.height,
                        num_inference_steps=self.config.steps,
                        guidance_scale=self.config.cfg_scale,
                        num_images_per_prompt=1
                    ).images
                    image = images[0]
                else:
                    # Mock generation for testing
                    from PIL import Image
                    import random
                    image = Image.new('RGB', (self.config.width, self.config.height),
                                     (random.randint(100, 200), random.randint(100, 200), random.randint(100, 200)))

                # Save result
                index = v * self.config.batch_size + i
                result = await self._save_result(image, character, template, index, positive, negative)
                results.append(result)

        return results

    async def _save_result(
        self,
        image,
        character: CharacterAttributes,
        template: str,
        index: int,
        prompt: str,
        negative_prompt: str
    ) -> Dict:
        """Save generated image and metadata"""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"{character.id}_{template}_{index:03d}_{timestamp}"
        char_dir = self.output_dir / character.id

        # Save image
        image_path = char_dir / f"{filename}.png"
        image.save(image_path, "PNG")

        # Calculate hash
        img_hash = hashlib.md5(image.tobytes()).hexdigest()

        # Build metadata
        metadata = {
            "schema_version": "1.0",
            "character": {
                "id": character.id,
                "name": character.name,
                "role": character.role,
                "age": character.age
            },
            "image": {
                "file_name": f"{filename}.png",
                "file_path": str(image_path),
                "dimensions": {"width": image.width, "height": image.height},
                "format": "PNG",
                "hash_md5": img_hash
            },
            "generation": {
                "model": self.config.model,
                "prompt": prompt,
                "negative_prompt": negative_prompt,
                "steps": self.config.steps,
                "cfg_scale": self.config.cfg_scale,
                "sampler": self.config.sampler,
                "generated_at": datetime.now().isoformat()
            },
            "attributes": asdict(character)
        }

        # Save metadata
        if self.config.save_metadata:
            metadata_path = char_dir / f"{filename}.json"
            with open(metadata_path, 'w', encoding='utf-8') as f:
                json.dump(metadata, f, indent=2, ensure_ascii=False)

        return {"image_path": str(image_path), "metadata": metadata}


# ============ Batch Generator ============

class BatchGenerator:
    """Batch generation manager"""

    def __init__(self, config: GenerationConfig):
        self.config = config
        self.generator = CharacterGenerator(config)

    async def generate_all(
        self,
        characters: List[CharacterAttributes],
        templates: List[str],
        variations_per_template: int = 2
    ) -> Dict[str, Dict[str, List]]:
        """Generate all characters with specified templates"""
        results = {}

        for char in characters:
            print(f"\n{'='*50}")
            print(f"Processing: {char.name}")
            print(f"{'='*50}")

            results[char.id] = {}

            for template in templates:
                images = await self.generator.generate(
                    char,
                    template=template,
                    variations=variations_per_template
                )
                results[char.id][template] = images

        return results


# ============ CLI ============

def main():
    parser = argparse.ArgumentParser(description="Generate anime character images")
    parser.add_argument("--character", "-c", type=str, help="Character ID to generate")
    parser.add_argument("--all", "-a", action="store_true", help="Generate all characters")
    parser.add_argument("--template", "-t", type=str, default="full_body",
                       help="Template type: full_body, portrait, expressions, poses, turnaround")
    parser.add_argument("--templates", type=str, help="Comma-separated list of templates")
    parser.add_argument("--variations", "-v", type=int, default=2, help="Number of variations")
    parser.add_argument("--output", "-o", type=str, default="output/characters", help="Output directory")
    parser.add_argument("--model", "-m", type=str, default="anything_v5", help="Model name")
    parser.add_argument("--steps", "-s", type=int, default=28, help="Generation steps")
    parser.add_argument("--cfg", type=float, default=7.0, help="CFG scale")
    parser.add_argument("--width", type=int, default=1024, help="Image width")
    parser.add_argument("--height", type=int, default=1024, help="Image height")
    parser.add_argument("--batch-size", "-b", type=int, default=4, help="Batch size")

    args = parser.parse_args()

    # Build config
    config = GenerationConfig(
        model=args.model,
        steps=args.steps,
        cfg_scale=args.cfg,
        width=args.width,
        height=args.height,
        batch_size=args.batch_size,
        output_dir=args.output
    )

    # Determine templates
    if args.templates:
        templates = [t.strip() for t in args.templates.split(",")]
    else:
        templates = [args.template]

    # Load characters
    if args.all:
        characters = CharacterProfileLoader.load()
    elif args.character:
        char = CharacterProfileLoader.get_by_id(args.character)
        if char:
            characters = [char]
        else:
            print(f"Error: Character '{args.character}' not found")
            return
    else:
        print("Please specify --character ID or --all")
        return

    # Run generation
    batch_gen = BatchGenerator(config)
    results = asyncio.run(batch_gen.generate_all(
        characters,
        templates=templates,
        variations_per_template=args.variations
    ))

    # Summary
    total = sum(len(imgs) for char_results in results.values()
                for imgs in char_results.values())
    print(f"\n{'='*50}")
    print(f"Generation complete!")
    print(f"Total images generated: {total}")
    print(f"Output directory: {args.output}")


if __name__ == "__main__":
    main()
