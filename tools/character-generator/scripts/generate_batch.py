#!/usr/bin/env python3
"""
Batch Generator - 基于属性定义的分批生成脚本

Usage:
    python generate_batch.py --batch 1  # 生成基础身体层
    python generate_batch.py --batch 2  # 生成发型层
    python generate_batch.py --all      # 生成所有批次
"""

import argparse
import asyncio
import json
import hashlib
from pathlib import Path
from datetime import datetime
from dataclasses import dataclass, asdict
from typing import List, Dict, Optional, Tuple
import sys

# 添加配置路径
sys.path.insert(0, str(Path(__file__).parent.parent))
from config.prompt_mappings import (
    HAIR_STYLES, HAIR_COLORS,
    EYE_SHAPES, EYE_COLORS,
    SKIN_TONES, GENDERS, BODY_TYPES,
    TOP_STYLES, BOTTOM_STYLES, SHOES_STYLES,
    CLOTHING_COLORS, BOTTOM_COLORS,
    HAT_STYLES, GLASSES_STYLES, JEWELRY_STYLES,
    COMBINATION_COUNTS
)


# ==================== 数据结构 ====================

@dataclass
class GenerationTask:
    """单个生成任务"""
    task_id: str
    category: str
    sub_category: str
    style_id: str
    style_name: Dict[str, str]
    color_id: Optional[str]
    color_hex: Optional[str]
    color_name: Optional[Dict[str, str]]
    output_filename: str
    output_path: str
    prompt: str
    negative_prompt: str
    metadata: Dict


# ==================== 提示词模板 ====================

PROMPT_TEMPLATES = {
    "body": {
        "positive": (
            "anime character sprite, {gender}, {skin_tone}, {body_type}, "
            "standing pose, front view, pixel art style, 64x64 pixels, "
            "simple shading, transparent background, no clothes, bald, simple line art"
        ),
        "negative": (
            "clothes, hair, eyes, detailed shading, complex background, "
            "watermark, signature, outline, shadow, accessories, 3d, realistic"
        )
    },
    "hair": {
        "positive": (
            "anime hair sprite, {style}, {color}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated hair only, simple shading, no face, no body"
        ),
        "negative": (
            "face, body, eyes, complex shading, background, "
            "watermark, signature, outline, shadow, 3d, realistic"
        )
    },
    "eyes": {
        "positive": (
            "anime eyes sprite, {style}, {color}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated eyes only, simple shading, no face, no hair"
        ),
        "negative": (
            "face, hair, body, complex shading, background, "
            "watermark, signature, outline, 3d, realistic"
        )
    },
    "top": {
        "positive": (
            "anime clothing sprite, {style}, {color}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated clothing only, simple shading, no body, no face"
        ),
        "negative": (
            "body, face, hair, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    },
    "bottom": {
        "positive": (
            "anime pants sprite, {style}, {color}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated clothing only, simple shading, no body, no face"
        ),
        "negative": (
            "body, face, hair, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    },
    "shoes": {
        "positive": (
            "anime shoes sprite, {style}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated shoes only, simple shading"
        ),
        "negative": (
            "body, face, legs, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    },
    "hat": {
        "positive": (
            "anime hat sprite, {style}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated accessory only, simple shading, no head"
        ),
        "negative": (
            "head, face, hair, body, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    },
    "glasses": {
        "positive": (
            "anime glasses sprite, {style}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated accessory only, simple shading"
        ),
        "negative": (
            "face, eyes, head, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    },
    "jewelry": {
        "positive": (
            "anime jewelry sprite, {style}, "
            "pixel art style, 64x64 pixels, transparent background, "
            "isolated accessory only, simple shading"
        ),
        "negative": (
            "body, face, complex shading, background, "
            "watermark, signature, 3d, realistic"
        )
    }
}


# ==================== 批次生成器 ====================

class BatchPlanGenerator:
    """生成批次计划"""

    def __init__(self, output_base: str = "output/assets"):
        self.output_base = Path(output_base)

    def generate_all_tasks(self) -> Dict[str, List[GenerationTask]]:
        """生成所有批次的任务"""
        return {
            "batch_1_body": self._generate_body_tasks(),
            "batch_2_hair": self._generate_hair_tasks(),
            "batch_3_eyes": self._generate_eyes_tasks(),
            "batch_4_top": self._generate_top_tasks(),
            "batch_5_bottom": self._generate_bottom_tasks(),
            "batch_6_shoes": self._generate_shoes_tasks(),
            "batch_7_hat": self._generate_hat_tasks(),
            "batch_8_glasses": self._generate_glasses_tasks(),
            "batch_9_jewelry": self._generate_jewelry_tasks()
        }

    def generate_batch(self, batch_num: int) -> List[GenerationTask]:
        """生成指定批次的任务"""
        batch_map = {
            1: self._generate_body_tasks,
            2: self._generate_hair_tasks,
            3: self._generate_eyes_tasks,
            4: self._generate_top_tasks,
            5: self._generate_bottom_tasks,
            6: self._generate_shoes_tasks,
            7: self._generate_hat_tasks,
            8: self._generate_glasses_tasks,
            9: self._generate_jewelry_tasks
        }
        if batch_num in batch_map:
            return batch_map[batch_num]()
        return []

    def _generate_body_tasks(self) -> List[GenerationTask]:
        """生成基础身体任务"""
        tasks = []
        template = PROMPT_TEMPLATES["body"]

        for gender_id, gender_data in GENDERS.items():
            for skin_id, skin_data in SKIN_TONES.items():
                for body_id, body_data in BODY_TYPES.items():
                    # 生成文件名
                    filename = f"body_{gender_data['sprite_suffix']}_{skin_id}_{skin_data['hex']}.png"
                    output_path = self.output_base / "base" / filename

                    # 构建提示词
                    prompt = template["positive"].format(
                        gender=gender_data["prompt"],
                        skin_tone=skin_data["prompt"],
                        body_type=body_data["prompt"]
                    )

                    task = GenerationTask(
                        task_id=f"body_{gender_id}_{skin_id}_{body_id}",
                        category="body",
                        sub_category="base",
                        style_id=gender_id,
                        style_name=gender_data["labels"],
                        color_id=skin_id,
                        color_hex=f"#{skin_data['hex']}",
                        color_name={"en": skin_id.replace("_", " ").title()},
                        output_filename=filename,
                        output_path=str(output_path),
                        prompt=prompt,
                        negative_prompt=template["negative"],
                        metadata={
                            "gender": gender_id,
                            "skin_tone": skin_id,
                            "body_type": body_id,
                            "height": body_data["height"],
                            "build": body_data["build"]
                        }
                    )
                    tasks.append(task)

        return tasks

    def _generate_hair_tasks(self) -> List[GenerationTask]:
        """生成发型任务"""
        tasks = []
        template = PROMPT_TEMPLATES["hair"]

        for style_id, style_data in HAIR_STYLES.items():
            # 跳过光头
            if style_id == "hair_00":
                continue

            for color_id, color_data in HAIR_COLORS.items():
                # 生成文件名 (hair_01_blue_3498db.png)
                filename = f"{style_id}_{color_id}_{color_data['hex']}.png"
                output_path = self.output_base / "hair" / filename

                # 构建提示词
                prompt = template["positive"].format(
                    style=style_data["prompt"],
                    color=color_data["prompt"]
                )

                task = GenerationTask(
                    task_id=f"hair_{style_id}_{color_id}",
                    category="hair",
                    sub_category="hair",
                    style_id=style_id,
                    style_name=style_data["labels"],
                    color_id=color_id,
                    color_hex=f"#{color_data['hex']}",
                    color_name={"en": color_id.replace("_", " ").title()},
                    output_filename=filename,
                    output_path=str(output_path),
                    prompt=prompt,
                    negative_prompt=template["negative"],
                    metadata={
                        "hair_style": style_data["id"],
                        "hair_color": color_id
                    }
                )
                tasks.append(task)

        return tasks

    def _generate_eyes_tasks(self) -> List[GenerationTask]:
        """生成眼睛任务"""
        tasks = []
        template = PROMPT_TEMPLATES["eyes"]

        for shape_id, shape_data in EYE_SHAPES.items():
            for color_id, color_data in EYE_COLORS.items():
                filename = f"{shape_id}_{color_id}_{color_data['hex']}.png"
                output_path = self.output_base / "eyes" / filename

                prompt = template["positive"].format(
                    style=shape_data["prompt"],
                    color=color_data["prompt"]
                )

                task = GenerationTask(
                    task_id=f"eyes_{shape_id}_{color_id}",
                    category="eyes",
                    sub_category="eyes",
                    style_id=shape_id,
                    style_name=shape_data["labels"],
                    color_id=color_id,
                    color_hex=f"#{color_data['hex']}",
                    color_name={"en": color_id.replace("_", " ").title()},
                    output_filename=filename,
                    output_path=str(output_path),
                    prompt=prompt,
                    negative_prompt=template["negative"],
                    metadata={
                        "eye_shape": shape_data["id"],
                        "eye_color": color_id
                    }
                )
                tasks.append(task)

        return tasks

    def _generate_top_tasks(self) -> List[GenerationTask]:
        """生成上衣任务"""
        tasks = []
        template = PROMPT_TEMPLATES["top"]

        for style_id, style_data in TOP_STYLES.items():
            for color_id, color_data in CLOTHING_COLORS.items():
                filename = f"{style_id}_{color_id}_{color_data['hex']}.png"
                output_path = self.output_base / "clothing" / "top" / filename

                prompt = template["positive"].format(
                    style=style_data["prompt"],
                    color=color_data["prompt"]
                )

                task = GenerationTask(
                    task_id=f"top_{style_id}_{color_id}",
                    category="clothing",
                    sub_category="top",
                    style_id=style_id,
                    style_name=style_data["labels"],
                    color_id=color_id,
                    color_hex=f"#{color_data['hex']}",
                    color_name={"en": color_id.replace("_", " ").title()},
                    output_filename=filename,
                    output_path=str(output_path),
                    prompt=prompt,
                    negative_prompt=template["negative"],
                    metadata={
                        "top_style": style_data["id"],
                        "top_color": color_id
                    }
                )
                tasks.append(task)

        return tasks

    def _generate_bottom_tasks(self) -> List[GenerationTask]:
        """生成下装任务"""
        tasks = []
        template = PROMPT_TEMPLATES["bottom"]

        for style_id, style_data in BOTTOM_STYLES.items():
            for color_id, color_data in BOTTOM_COLORS.items():
                filename = f"{style_id}_{color_id}_{color_data['hex']}.png"
                output_path = self.output_base / "clothing" / "bottom" / filename

                prompt = template["positive"].format(
                    style=style_data["prompt"],
                    color=color_data["prompt"]
                )

                task = GenerationTask(
                    task_id=f"bottom_{style_id}_{color_id}",
                    category="clothing",
                    sub_category="bottom",
                    style_id=style_id,
                    style_name=style_data["labels"],
                    color_id=color_id,
                    color_hex=f"#{color_data['hex']}",
                    color_name={"en": color_id.replace("_", " ").title()},
                    output_filename=filename,
                    output_path=str(output_path),
                    prompt=prompt,
                    negative_prompt=template["negative"],
                    metadata={
                        "bottom_style": style_data["id"],
                        "bottom_color": color_id
                    }
                )
                tasks.append(task)

        return tasks

    def _generate_shoes_tasks(self) -> List[GenerationTask]:
        """生成鞋子任务"""
        tasks = []
        template = PROMPT_TEMPLATES["shoes"]

        for style_id, style_data in SHOES_STYLES.items():
            filename = f"{style_id}.png"
            output_path = self.output_base / "clothing" / "shoes" / filename

            prompt = template["positive"].format(style=style_data["prompt"])

            task = GenerationTask(
                task_id=f"shoes_{style_id}",
                category="clothing",
                sub_category="shoes",
                style_id=style_id,
                style_name=style_data["labels"],
                color_id=None,
                color_hex=None,
                color_name=None,
                output_filename=filename,
                output_path=str(output_path),
                prompt=prompt,
                negative_prompt=template["negative"],
                metadata={"shoes_style": style_data["id"]}
            )
            tasks.append(task)

        return tasks

    def _generate_hat_tasks(self) -> List[GenerationTask]:
        """生成帽子任务"""
        tasks = []
        template = PROMPT_TEMPLATES["hat"]

        for style_id, style_data in HAT_STYLES.items():
            filename = f"{style_id}.png"
            output_path = self.output_base / "accessories" / "hat" / filename

            prompt = template["positive"].format(style=style_data["prompt"])

            task = GenerationTask(
                task_id=f"hat_{style_id}",
                category="accessories",
                sub_category="hat",
                style_id=style_id,
                style_name=style_data["labels"],
                color_id=None,
                color_hex=None,
                color_name=None,
                output_filename=filename,
                output_path=str(output_path),
                prompt=prompt,
                negative_prompt=template["negative"],
                metadata={"hat_style": style_data["id"]}
            )
            tasks.append(task)

        return tasks

    def _generate_glasses_tasks(self) -> List[GenerationTask]:
        """生成眼镜任务"""
        tasks = []
        template = PROMPT_TEMPLATES["glasses"]

        for style_id, style_data in GLASSES_STYLES.items():
            filename = f"{style_id}.png"
            output_path = self.output_base / "accessories" / "glasses" / filename

            prompt = template["positive"].format(style=style_data["prompt"])

            task = GenerationTask(
                task_id=f"glasses_{style_id}",
                category="accessories",
                sub_category="glasses",
                style_id=style_id,
                style_name=style_data["labels"],
                color_id=None,
                color_hex=None,
                color_name=None,
                output_filename=filename,
                output_path=str(output_path),
                prompt=prompt,
                negative_prompt=template["negative"],
                metadata={"glasses_style": style_data["id"]}
            )
            tasks.append(task)

        return tasks

    def _generate_jewelry_tasks(self) -> List[GenerationTask]:
        """生成首饰任务"""
        tasks = []
        template = PROMPT_TEMPLATES["jewelry"]

        for style_id, style_data in JEWELRY_STYLES.items():
            filename = f"{style_id}.png"
            output_path = self.output_base / "accessories" / "jewelry" / filename

            prompt = template["positive"].format(style=style_data["prompt"])

            task = GenerationTask(
                task_id=f"jewelry_{style_id}",
                category="accessories",
                sub_category="jewelry",
                style_id=style_id,
                style_name=style_data["labels"],
                color_id=None,
                color_hex=None,
                color_name=None,
                output_filename=filename,
                output_path=str(output_path),
                prompt=prompt,
                negative_prompt=template["negative"],
                metadata={"jewelry_style": style_data["id"]}
            )
            tasks.append(task)

        return tasks


# ==================== 任务导出 ====================

def export_task_list(tasks: List[GenerationTask], output_file: str):
    """导出任务列表为JSON"""
    data = {
        "exported_at": datetime.now().isoformat(),
        "total_tasks": len(tasks),
        "tasks": [asdict(t) for t in tasks]
    }

    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

    print(f"Exported {len(tasks)} tasks to {output_file}")


def export_manifest(all_tasks: Dict[str, List[GenerationTask]], output_file: str):
    """导出完整清单"""
    manifest = {
        "manifest_version": "1.0.0",
        "generated_at": datetime.now().isoformat(),
        "categories": {},
        "total_files": 0
    }

    for batch_name, tasks in all_tasks.items():
        category = tasks[0].category if tasks else "unknown"
        manifest["categories"][batch_name] = {
            "category": category,
            "total": len(tasks),
            "files": [t.output_filename for t in tasks[:10]]  # 前10个示例
        }
        manifest["total_files"] += len(tasks)

    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(manifest, f, indent=2, ensure_ascii=False)

    print(f"Exported manifest to {output_file}")
    print(f"Total files to generate: {manifest['total_files']}")


# ==================== 主入口 ====================

def main():
    parser = argparse.ArgumentParser(description="Generate character asset batch plan")
    parser.add_argument("--batch", "-b", type=int, help="Batch number (1-9)")
    parser.add_argument("--all", "-a", action="store_true", help="Generate all batches")
    parser.add_argument("--output", "-o", type=str, default="output", help="Output directory")
    parser.add_argument("--export-tasks", action="store_true", help="Export task list to JSON")
    parser.add_argument("--stats", action="store_true", help="Show statistics only")

    args = parser.parse_args()

    planner = BatchPlanGenerator(output_base=f"{args.output}/assets")

    # 统计模式
    if args.stats:
        print("\n" + "="*60)
        print("角色资源组合统计")
        print("="*60)
        for category, count in COMBINATION_COUNTS.items():
            print(f"  {category}: {count} 个素材")
        print("-"*60)
        print(f"  总计: {sum(COMBINATION_COUNTS.values())} 个独立素材")
        print("="*60)
        return

    # 生成所有批次
    if args.all:
        all_tasks = planner.generate_all_tasks()

        print("\n" + "="*60)
        print("批量生成计划")
        print("="*60)

        for batch_name, tasks in all_tasks.items():
            print(f"  {batch_name}: {len(tasks)} 个任务")

        total = sum(len(t) for t in all_tasks.values())
        print("-"*60)
        print(f"  总计: {total} 个生成任务")
        print("="*60)

        if args.export_tasks:
            export_manifest(all_tasks, f"{args.output}/manifest.json")

            # 导出每个批次的任务列表
            for batch_name, tasks in all_tasks.items():
                export_task_list(tasks, f"{args.output}/{batch_name}_tasks.json")

        return

    # 生成指定批次
    if args.batch:
        tasks = planner.generate_batch(args.batch)

        print(f"\nBatch {args.batch}: {len(tasks)} 个任务")
        print("-"*40)

        # 显示前5个任务示例
        for task in tasks[:5]:
            print(f"  [{task.task_id}]")
            print(f"    文件: {task.output_filename}")
            print(f"    提示词: {task.prompt[:80]}...")

        if len(tasks) > 5:
            print(f"  ... 还有 {len(tasks)-5} 个任务")

        if args.export_tasks:
            export_task_list(tasks, f"{args.output}/batch_{args.batch}_tasks.json")

        return

    # 默认显示帮助
    parser.print_help()


if __name__ == "__main__":
    main()
