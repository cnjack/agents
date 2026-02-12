# Character Generator UI/UX Design Document

## Overview

This document defines the character generator system for the Stardew Valley-style farming game. It includes attribute categories, UI/UX design specifications, and JSON Schema for character definitions with bilingual support (Chinese/English).

---

## 1. Character Attribute Categories

### 1.1 Category Overview

| Category | Description | UI Component Type |
|----------|-------------|-------------------|
| Basic Info | Name, Gender, Birthday | Text Input, Select |
| Hairstyle | Hair style and color | Grid Selection + Color Picker |
| Eyes | Eye shape and color | Grid Selection + Color Picker |
| Skin | Skin tone | Color Palette |
| Clothing | Outfit type and color | Grid Selection + Color Picker |
| Accessories | Hats, glasses, jewelry | Grid Selection |
| Body Type | Height, build | Slider + Presets |
| Personality | Character traits | Tag Selection |

---

## 2. Detailed Attribute Definitions

### 2.1 Basic Information (基础信息)

```json
{
  "category_id": "basic_info",
  "labels": {
    "en": "Basic Information",
    "zh": "基础信息"
  },
  "attributes": [
    {
      "id": "name",
      "labels": {
        "en": "Name",
        "zh": "姓名"
      },
      "type": "text",
      "required": true,
      "validation": {
        "min_length": 1,
        "max_length": 20,
        "pattern": "^[a-zA-Z\\u4e00-\\u9fa5]+$"
      },
      "placeholder": {
        "en": "Enter your name",
        "zh": "输入你的名字"
      }
    },
    {
      "id": "gender",
      "labels": {
        "en": "Gender",
        "zh": "性别"
      },
      "type": "select",
      "required": true,
      "options": [
        {
          "id": "male",
          "labels": { "en": "Male", "zh": "男" },
          "sprite_suffix": "_m"
        },
        {
          "id": "female",
          "labels": { "en": "Female", "zh": "女" },
          "sprite_suffix": "_f"
        },
        {
          "id": "non_binary",
          "labels": { "en": "Non-binary", "zh": "非二元" },
          "sprite_suffix": "_nb"
        }
      ],
      "default": "male"
    },
    {
      "id": "birthday",
      "labels": {
        "en": "Birthday",
        "zh": "生日"
      },
      "type": "date_select",
      "required": true,
      "fields": [
        {
          "id": "season",
          "labels": { "en": "Season", "zh": "季节" },
          "options": [
            { "id": "spring", "labels": { "en": "Spring", "zh": "春季" } },
            { "id": "summer", "labels": { "en": "Summer", "zh": "夏季" } },
            { "id": "fall", "labels": { "en": "Fall", "zh": "秋季" } },
            { "id": "winter", "labels": { "en": "Winter", "zh": "冬季" } }
          ]
        },
        {
          "id": "day",
          "labels": { "en": "Day", "zh": "日期" },
          "range": { "min": 1, "max": 28 }
        }
      ]
    },
    {
      "id": "favorite_thing",
      "labels": {
        "en": "Favorite Thing",
        "zh": "最喜欢的事物"
      },
      "type": "text",
      "required": false,
      "max_length": 30,
      "description": {
        "en": "This will appear in certain menus",
        "zh": "这将在某些菜单中显示"
      }
    }
  ]
}
```

### 2.2 Hairstyle (发型)

```json
{
  "category_id": "hairstyle",
  "labels": {
    "en": "Hairstyle",
    "zh": "发型"
  },
  "attributes": [
    {
      "id": "hair_style",
      "labels": {
        "en": "Hair Style",
        "zh": "发型样式"
      },
      "type": "grid_select",
      "options": [
        { "id": "short_spiky", "labels": { "en": "Short Spiky", "zh": "短刺头" }, "sprite_id": "hair_01" },
        { "id": "short_neat", "labels": { "en": "Short Neat", "zh": "短发整齐" }, "sprite_id": "hair_02" },
        { "id": "medium_wavy", "labels": { "en": "Medium Wavy", "zh": "中长波浪" }, "sprite_id": "hair_03" },
        { "id": "medium_straight", "labels": { "en": "Medium Straight", "zh": "中长直发" }, "sprite_id": "hair_04" },
        { "id": "long_straight", "labels": { "en": "Long Straight", "zh": "长直发" }, "sprite_id": "hair_05" },
        { "id": "long_wavy", "labels": { "en": "Long Wavy", "zh": "长波浪" }, "sprite_id": "hair_06" },
        { "id": "ponytail", "labels": { "en": "Ponytail", "zh": "马尾辫" }, "sprite_id": "hair_07" },
        { "id": "braid", "labels": { "en": "Braid", "zh": "麻花辫" }, "sprite_id": "hair_08" },
        { "id": "bald", "labels": { "en": "Bald", "zh": "光头" }, "sprite_id": "hair_00" },
        { "id": "curly_short", "labels": { "en": "Curly Short", "zh": "短卷发" }, "sprite_id": "hair_09" },
        { "id": "curly_long", "labels": { "en": "Curly Long", "zh": "长卷发" }, "sprite_id": "hair_10" },
        { "id": "mohawk", "labels": { "en": "Mohawk", "zh": "莫霍克" }, "sprite_id": "hair_11" },
        { "id": "buns", "labels": { "en": "Double Buns", "zh": "双丸子头" }, "sprite_id": "hair_12" },
        { "id": "undercut", "labels": { "en": "Undercut", "zh": "铲青" }, "sprite_id": "hair_13" },
        { "id": "afro", "labels": { "en": "Afro", "zh": "爆炸头" }, "sprite_id": "hair_14" },
        { "id": "dreadlocks", "labels": { "en": "Dreadlocks", "zh": "脏辫" }, "sprite_id": "hair_15" }
      ],
      "default": "short_neat",
      "grid_layout": { "columns": 4, "rows": 4 }
    },
    {
      "id": "hair_color",
      "labels": {
        "en": "Hair Color",
        "zh": "发色"
      },
      "type": "color_picker",
      "palette": [
        { "id": "black", "labels": { "en": "Black", "zh": "黑色" }, "hex": "#1a1a1a" },
        { "id": "brown_dark", "labels": { "en": "Dark Brown", "zh": "深棕色" }, "hex": "#3d2314" },
        { "id": "brown_light", "labels": { "en": "Light Brown", "zh": "浅棕色" }, "hex": "#8b4513" },
        { "id": "blonde", "labels": { "en": "Blonde", "zh": "金色" }, "hex": "#f4d03f" },
        { "id": "platinum", "labels": { "en": "Platinum", "zh": "白金色" }, "hex": "#e8e4c9" },
        { "id": "red", "labels": { "en": "Red", "zh": "红色" }, "hex": "#8b2500" },
        { "id": "ginger", "labels": { "en": "Ginger", "zh": "姜红色" }, "hex": "#d35400" },
        { "id": "auburn", "labels": { "en": "Auburn", "zh": "赤褐色" }, "hex": "#922b21" },
        { "id": "gray", "labels": { "en": "Gray", "zh": "灰色" }, "hex": "#808080" },
        { "id": "white", "labels": { "en": "White", "zh": "白色" }, "hex": "#f5f5f5" },
        { "id": "blue", "labels": { "en": "Blue", "zh": "蓝色" }, "hex": "#3498db" },
        { "id": "pink", "labels": { "en": "Pink", "zh": "粉色" }, "hex": "#ff69b4" },
        { "id": "purple", "labels": { "en": "Purple", "zh": "紫色" }, "hex": "#9b59b6" },
        { "id": "green", "labels": { "en": "Green", "zh": "绿色" }, "hex": "#27ae60" },
        { "id": "teal", "labels": { "en": "Teal", "zh": "青色" }, "hex": "#1abc9c" }
      ],
      "allow_custom": true,
      "default": "brown_dark"
    }
  ]
}
```

### 2.3 Eyes (眼睛)

```json
{
  "category_id": "eyes",
  "labels": {
    "en": "Eyes",
    "zh": "眼睛"
  },
  "attributes": [
    {
      "id": "eye_shape",
      "labels": {
        "en": "Eye Shape",
        "zh": "眼型"
      },
      "type": "grid_select",
      "options": [
        { "id": "round", "labels": { "en": "Round", "zh": "圆眼" }, "sprite_id": "eyes_01" },
        { "id": "almond", "labels": { "en": "Almond", "zh": "杏眼" }, "sprite_id": "eyes_02" },
        { "id": "narrow", "labels": { "en": "Narrow", "zh": "细长眼" }, "sprite_id": "eyes_03" },
        { "id": "wide", "labels": { "en": "Wide", "zh": "大眼" }, "sprite_id": "eyes_04" },
        { "id": "droopy", "labels": { "en": "Droopy", "zh": "下垂眼" }, "sprite_id": "eyes_05" },
        { "id": "cat_eye", "labels": { "en": "Cat Eye", "zh": "猫眼" }, "sprite_id": "eyes_06" },
        { "id": "monolid", "labels": { "en": "Monolid", "zh": "单眼皮" }, "sprite_id": "eyes_07" },
        { "id": "hooded", "labels": { "en": "Hooded", "zh": "内双" }, "sprite_id": "eyes_08" }
      ],
      "default": "round",
      "grid_layout": { "columns": 4, "rows": 2 }
    },
    {
      "id": "eye_color",
      "labels": {
        "en": "Eye Color",
        "zh": "瞳色"
      },
      "type": "color_picker",
      "palette": [
        { "id": "brown_dark", "labels": { "en": "Dark Brown", "zh": "深棕色" }, "hex": "#3d2314" },
        { "id": "brown_light", "labels": { "en": "Light Brown", "zh": "浅棕色" }, "hex": "#8b4513" },
        { "id": "hazel", "labels": { "en": "Hazel", "zh": "榛色" }, "hex": "#a67c52" },
        { "id": "amber", "labels": { "en": "Amber", "zh": "琥珀色" }, "hex": "#ffbf00" },
        { "id": "blue", "labels": { "en": "Blue", "zh": "蓝色" }, "hex": "#3498db" },
        { "id": "blue_light", "labels": { "en": "Light Blue", "zh": "浅蓝色" }, "hex": "#87ceeb" },
        { "id": "green", "labels": { "en": "Green", "zh": "绿色" }, "hex": "#27ae60" },
        { "id": "gray", "labels": { "en": "Gray", "zh": "灰色" }, "hex": "#808080" },
        { "id": "violet", "labels": { "en": "Violet", "zh": "紫色" }, "hex": "#9b59b6" },
        { "id": "red", "labels": { "en": "Red", "zh": "红色" }, "hex": "#c0392b" }
      ],
      "allow_custom": false,
      "default": "brown_dark"
    }
  ]
}
```

### 2.4 Skin Tone (肤色)

```json
{
  "category_id": "skin",
  "labels": {
    "en": "Skin Tone",
    "zh": "肤色"
  },
  "attributes": [
    {
      "id": "skin_tone",
      "labels": {
        "en": "Skin Tone",
        "zh": "肤色"
      },
      "type": "color_palette",
      "palette": [
        { "id": "pale", "labels": { "en": "Pale", "zh": "苍白" }, "hex": "#ffe4c4" },
        { "id": "fair", "labels": { "en": "Fair", "zh": "白皙" }, "hex": "#f5deb3" },
        { "id": "light", "labels": { "en": "Light", "zh": "浅色" }, "hex": "#deb887" },
        { "id": "medium", "labels": { "en": "Medium", "zh": "中等" }, "hex": "#d2a679" },
        { "id": "tan", "labels": { "en": "Tan", "zh": "古铜" }, "hex": "#c68642" },
        { "id": "olive", "labels": { "en": "Olive", "zh": "橄榄色" }, "hex": "#b5a642" },
        { "id": "brown", "labels": { "en": "Brown", "zh": "棕色" }, "hex": "#8d5524" },
        { "id": "dark_brown", "labels": { "en": "Dark Brown", "zh": "深棕色" }, "hex": "#5c4033" },
        { "id": "black", "labels": { "en": "Black", "zh": "黑色" }, "hex": "#3d2314" },
        { "id": "fantasy_blue", "labels": { "en": "Fantasy Blue", "zh": "奇幻蓝" }, "hex": "#87ceeb" },
        { "id": "fantasy_green", "labels": { "en": "Fantasy Green", "zh": "奇幻绿" }, "hex": "#90ee90" },
        { "id": "fantasy_purple", "labels": { "en": "Fantasy Purple", "zh": "奇幻紫" }, "hex": "#dda0dd" }
      ],
      "default": "fair",
      "layout": "horizontal_bar"
    }
  ]
}
```

### 2.5 Clothing (服装)

```json
{
  "category_id": "clothing",
  "labels": {
    "en": "Clothing",
    "zh": "服装"
  },
  "attributes": [
    {
      "id": "top",
      "labels": {
        "en": "Top",
        "zh": "上衣"
      },
      "type": "grid_select",
      "options": [
        { "id": "t_shirt", "labels": { "en": "T-Shirt", "zh": "T恤" }, "sprite_id": "top_01" },
        { "id": "shirt", "labels": { "en": "Shirt", "zh": "衬衫" }, "sprite_id": "top_02" },
        { "id": "sweater", "labels": { "en": "Sweater", "zh": "毛衣" }, "sprite_id": "top_03" },
        { "id": "hoodie", "labels": { "en": "Hoodie", "zh": "连帽衫" }, "sprite_id": "top_04" },
        { "id": "overalls", "labels": { "en": "Overalls", "zh": "工装背带裤" }, "sprite_id": "top_05" },
        { "id": "vest", "labels": { "en": "Vest", "zh": "马甲" }, "sprite_id": "top_06" },
        { "id": "flannel", "labels": { "en": "Flannel", "zh": "法兰绒衬衫" }, "sprite_id": "top_07" },
        { "id": "dress", "labels": { "en": "Dress", "zh": "连衣裙" }, "sprite_id": "top_08" }
      ],
      "default": "t_shirt",
      "grid_layout": { "columns": 4, "rows": 2 }
    },
    {
      "id": "top_color",
      "labels": {
        "en": "Top Color",
        "zh": "上衣颜色"
      },
      "type": "color_picker",
      "palette": [
        { "id": "white", "labels": { "en": "White", "zh": "白色" }, "hex": "#ffffff" },
        { "id": "black", "labels": { "en": "Black", "zh": "黑色" }, "hex": "#1a1a1a" },
        { "id": "red", "labels": { "en": "Red", "zh": "红色" }, "hex": "#c0392b" },
        { "id": "blue", "labels": { "en": "Blue", "zh": "蓝色" }, "hex": "#3498db" },
        { "id": "green", "labels": { "en": "Green", "zh": "绿色" }, "hex": "#27ae60" },
        { "id": "yellow", "labels": { "en": "Yellow", "zh": "黄色" }, "hex": "#f1c40f" },
        { "id": "orange", "labels": { "en": "Orange", "zh": "橙色" }, "hex": "#e67e22" },
        { "id": "purple", "labels": { "en": "Purple", "zh": "紫色" }, "hex": "#9b59b6" },
        { "id": "pink", "labels": { "en": "Pink", "zh": "粉色" }, "hex": "#ff69b4" },
        { "id": "brown", "labels": { "en": "Brown", "zh": "棕色" }, "hex": "#8b4513" },
        { "id": "gray", "labels": { "en": "Gray", "zh": "灰色" }, "hex": "#808080" },
        { "id": "navy", "labels": { "en": "Navy", "zh": "深蓝色" }, "hex": "#1a3a5c" }
      ],
      "allow_custom": true,
      "default": "blue"
    },
    {
      "id": "bottom",
      "labels": {
        "en": "Bottom",
        "zh": "下装"
      },
      "type": "grid_select",
      "options": [
        { "id": "jeans", "labels": { "en": "Jeans", "zh": "牛仔裤" }, "sprite_id": "bottom_01" },
        { "id": "shorts", "labels": { "en": "Shorts", "zh": "短裤" }, "sprite_id": "bottom_02" },
        { "id": "cargo_pants", "labels": { "en": "Cargo Pants", "zh": "工装裤" }, "sprite_id": "bottom_03" },
        { "id": "skirt", "labels": { "en": "Skirt", "zh": "裙子" }, "sprite_id": "bottom_04" },
        { "id": "overalls", "labels": { "en": "Overalls", "zh": "背带裤" }, "sprite_id": "bottom_05" },
        { "id": "sweatpants", "labels": { "en": "Sweatpants", "zh": "运动裤" }, "sprite_id": "bottom_06" }
      ],
      "default": "jeans",
      "grid_layout": { "columns": 3, "rows": 2 }
    },
    {
      "id": "bottom_color",
      "labels": {
        "en": "Bottom Color",
        "zh": "下装颜色"
      },
      "type": "color_picker",
      "palette": [
        { "id": "blue_denim", "labels": { "en": "Blue Denim", "zh": "蓝色牛仔" }, "hex": "#4a6fa5" },
        { "id": "black_denim", "labels": { "en": "Black Denim", "zh": "黑色牛仔" }, "hex": "#2c3e50" },
        { "id": "khaki", "labels": { "en": "Khaki", "zh": "卡其色" }, "hex": "#c3b091" },
        { "id": "brown", "labels": { "en": "Brown", "zh": "棕色" }, "hex": "#8b4513" },
        { "id": "white", "labels": { "en": "White", "zh": "白色" }, "hex": "#f5f5f5" },
        { "id": "gray", "labels": { "en": "Gray", "zh": "灰色" }, "hex": "#808080" }
      ],
      "allow_custom": true,
      "default": "blue_denim"
    },
    {
      "id": "shoes",
      "labels": {
        "en": "Shoes",
        "zh": "鞋子"
      },
      "type": "grid_select",
      "options": [
        { "id": "sneakers", "labels": { "en": "Sneakers", "zh": "运动鞋" }, "sprite_id": "shoes_01" },
        { "id": "boots", "labels": { "en": "Boots", "zh": "靴子" }, "sprite_id": "shoes_02" },
        { "id": "sandals", "labels": { "en": "Sandals", "zh": "凉鞋" }, "sprite_id": "shoes_03" },
        { "id": "work_boots", "labels": { "en": "Work Boots", "zh": "工装靴" }, "sprite_id": "shoes_04" }
      ],
      "default": "sneakers",
      "grid_layout": { "columns": 4, "rows": 1 }
    }
  ]
}
```

### 2.6 Accessories (配饰)

```json
{
  "category_id": "accessories",
  "labels": {
    "en": "Accessories",
    "zh": "配饰"
  },
  "attributes": [
    {
      "id": "hat",
      "labels": {
        "en": "Hat",
        "zh": "帽子"
      },
      "type": "grid_select",
      "optional": true,
      "options": [
        { "id": "none", "labels": { "en": "None", "zh": "无" }, "sprite_id": null },
        { "id": "baseball_cap", "labels": { "en": "Baseball Cap", "zh": "棒球帽" }, "sprite_id": "hat_01" },
        { "id": "straw_hat", "labels": { "en": "Straw Hat", "zh": "草帽" }, "sprite_id": "hat_02" },
        { "id": "beanie", "labels": { "en": "Beanie", "zh": "针织帽" }, "sprite_id": "hat_03" },
        { "id": "cowboy_hat", "labels": { "en": "Cowboy Hat", "zh": "牛仔帽" }, "sprite_id": "hat_04" },
        { "id": "sun_hat", "labels": { "en": "Sun Hat", "zh": "遮阳帽" }, "sprite_id": "hat_05" },
        { "id": "bowler", "labels": { "en": "Bowler Hat", "zh": "圆顶礼帽" }, "sprite_id": "hat_06" },
        { "id": "witch_hat", "labels": { "en": "Witch Hat", "zh": "巫师帽" }, "sprite_id": "hat_07" }
      ],
      "default": "none",
      "grid_layout": { "columns": 4, "rows": 2 }
    },
    {
      "id": "glasses",
      "labels": {
        "en": "Glasses",
        "zh": "眼镜"
      },
      "type": "grid_select",
      "optional": true,
      "options": [
        { "id": "none", "labels": { "en": "None", "zh": "无" }, "sprite_id": null },
        { "id": "round_glasses", "labels": { "en": "Round Glasses", "zh": "圆框眼镜" }, "sprite_id": "glasses_01" },
        { "id": "square_glasses", "labels": { "en": "Square Glasses", "zh": "方框眼镜" }, "sprite_id": "glasses_02" },
        { "id": "sunglasses", "labels": { "en": "Sunglasses", "zh": "墨镜" }, "sprite_id": "glasses_03" },
        { "id": "reading_glasses", "labels": { "en": "Reading Glasses", "zh": "老花镜" }, "sprite_id": "glasses_04" },
        { "id": "monocle", "labels": { "en": "Monocle", "zh": "单片眼镜" }, "sprite_id": "glasses_05" }
      ],
      "default": "none",
      "grid_layout": { "columns": 3, "rows": 2 }
    },
    {
      "id": "jewelry",
      "labels": {
        "en": "Jewelry",
        "zh": "首饰"
      },
      "type": "multi_select",
      "optional": true,
      "max_selections": 3,
      "options": [
        { "id": "none", "labels": { "en": "None", "zh": "无" }, "sprite_id": null },
        { "id": "earrings_stud", "labels": { "en": "Stud Earrings", "zh": "耳钉" }, "sprite_id": "jewelry_01" },
        { "id": "earrings_hoop", "labels": { "en": "Hoop Earrings", "zh": "耳环" }, "sprite_id": "jewelry_02" },
        { "id": "necklace", "labels": { "en": "Necklace", "zh": "项链" }, "sprite_id": "jewelry_03" },
        { "id": "scarf", "labels": { "en": "Scarf", "zh": "围巾" }, "sprite_id": "jewelry_04" },
        { "id": "watch", "labels": { "en": "Watch", "zh": "手表" }, "sprite_id": "jewelry_05" }
      ],
      "default": ["none"]
    }
  ]
}
```

### 2.7 Body Type (体型)

```json
{
  "category_id": "body_type",
  "labels": {
    "en": "Body Type",
    "zh": "体型"
  },
  "attributes": [
    {
      "id": "height",
      "labels": {
        "en": "Height",
        "zh": "身高"
      },
      "type": "slider",
      "range": { "min": 0.8, "max": 1.2, "step": 0.05 },
      "default": 1.0,
      "display_labels": {
        "en": { "min": "Short", "max": "Tall" },
        "zh": { "min": "矮", "max": "高" }
      },
      "presets": [
        { "id": "short", "labels": { "en": "Short", "zh": "矮" }, "value": 0.85 },
        { "id": "average", "labels": { "en": "Average", "zh": "中等" }, "value": 1.0 },
        { "id": "tall", "labels": { "en": "Tall", "zh": "高" }, "value": 1.15 }
      ]
    },
    {
      "id": "build",
      "labels": {
        "en": "Build",
        "zh": "体型"
      },
      "type": "slider",
      "range": { "min": 0.7, "max": 1.3, "step": 0.1 },
      "default": 1.0,
      "display_labels": {
        "en": { "min": "Slim", "max": "Broad" },
        "zh": { "min": "纤细", "max": "壮实" }
      },
      "presets": [
        { "id": "slim", "labels": { "en": "Slim", "zh": "纤细" }, "value": 0.8 },
        { "id": "average", "labels": { "en": "Average", "zh": "中等" }, "value": 1.0 },
        { "id": "broad", "labels": { "en": "Broad", "zh": "壮实" }, "value": 1.2 }
      ]
    }
  ]
}
```

### 2.8 Personality (性格)

```json
{
  "category_id": "personality",
  "labels": {
    "en": "Personality",
    "zh": "性格"
  },
  "attributes": [
    {
      "id": "primary_trait",
      "labels": {
        "en": "Primary Trait",
        "zh": "主要性格"
      },
      "type": "tag_select",
      "required": true,
      "options": [
        { "id": "friendly", "labels": { "en": "Friendly", "zh": "友好" }, "icon": "smile" },
        { "id": "shy", "labels": { "en": "Shy", "zh": "害羞" }, "icon": "blush" },
        { "id": "energetic", "labels": { "en": "Energetic", "zh": "活力充沛" }, "icon": "zap" },
        { "id": "calm", "labels": { "en": "Calm", "zh": "冷静" }, "icon": "leaf" },
        { "id": "curious", "labels": { "en": "Curious", "zh": "好奇" }, "icon": "search" },
        { "id": "serious", "labels": { "en": "Serious", "zh": "认真" }, "icon": "briefcase" },
        { "id": "cheerful", "labels": { "en": "Cheerful", "zh": "开朗" }, "icon": "sun" },
        { "id": "mysterious", "labels": { "en": "Mysterious", "zh": "神秘" }, "icon": "moon" }
      ],
      "default": "friendly"
    },
    {
      "id": "secondary_traits",
      "labels": {
        "en": "Secondary Traits",
        "zh": "次要性格"
      },
      "type": "multi_tag_select",
      "required": false,
      "min_selections": 0,
      "max_selections": 3,
      "options": [
        { "id": "hardworking", "labels": { "en": "Hardworking", "zh": "勤劳" }, "icon": "hammer" },
        { "id": "creative", "labels": { "en": "Creative", "zh": "创意" }, "icon": "palette" },
        { "id": "adventurous", "labels": { "en": "Adventurous", "zh": "冒险" }, "icon": "compass" },
        { "id": "romantic", "labels": { "en": "Romantic", "zh": "浪漫" }, "icon": "heart" },
        { "id": "practical", "labels": { "en": "Practical", "zh": "务实" }, "icon": "wrench" },
        { "id": "optimistic", "labels": { "en": "Optimistic", "zh": "乐观" }, "icon": "star" },
        { "id": "honest", "labels": { "en": "Honest", "zh": "诚实" }, "icon": "check" },
        { "id": "generous", "labels": { "en": "Generous", "zh": "慷慨" }, "icon": "gift" }
      ]
    },
    {
      "id": "hobby",
      "labels": {
        "en": "Hobby",
        "zh": "爱好"
      },
      "type": "tag_select",
      "required": false,
      "options": [
        { "id": "farming", "labels": { "en": "Farming", "zh": "种田" }, "icon": "wheat" },
        { "id": "fishing", "labels": { "en": "Fishing", "zh": "钓鱼" }, "icon": "fish" },
        { "id": "mining", "labels": { "en": "Mining", "zh": "采矿" }, "icon": "gem" },
        { "id": "foraging", "labels": { "en": "Foraging", "zh": "采集" }, "icon": "tree" },
        { "id": "cooking", "labels": { "en": "Cooking", "zh": "烹饪" }, "icon": "chef-hat" },
        { "id": "crafting", "labels": { "en": "Crafting", "zh": "手工" }, "icon": "scissors" },
        { "id": "socializing", "labels": { "en": "Socializing", "zh": "社交" }, "icon": "users" },
        { "id": "exploring", "labels": { "en": "Exploring", "zh": "探索" }, "icon": "map" }
      ]
    }
  ]
}
```

---

## 3. JSON Schema for Character Definition

### 3.1 Complete Character Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://stardew-agents.local/schemas/character.json",
  "title": "Character Definition",
  "description": "Schema for character creation and storage",
  "type": "object",
  "required": ["id", "basic_info", "appearance", "personality"],
  "properties": {
    "id": {
      "type": "string",
      "description": "Unique character identifier",
      "pattern": "^[a-z0-9_]+$"
    },
    "version": {
      "type": "string",
      "description": "Schema version",
      "default": "1.0.0"
    },
    "created_at": {
      "type": "string",
      "format": "date-time"
    },
    "updated_at": {
      "type": "string",
      "format": "date-time"
    },
    "basic_info": {
      "type": "object",
      "required": ["name", "gender", "birthday"],
      "properties": {
        "name": {
          "type": "string",
          "minLength": 1,
          "maxLength": 20,
          "description": "Character name"
        },
        "gender": {
          "type": "string",
          "enum": ["male", "female", "non_binary"]
        },
        "birthday": {
          "type": "object",
          "required": ["season", "day"],
          "properties": {
            "season": {
              "type": "string",
              "enum": ["spring", "summer", "fall", "winter"]
            },
            "day": {
              "type": "integer",
              "minimum": 1,
              "maximum": 28
            }
          }
        },
        "favorite_thing": {
          "type": "string",
          "maxLength": 30
        }
      }
    },
    "appearance": {
      "type": "object",
      "required": ["hairstyle", "eyes", "skin", "clothing", "body_type"],
      "properties": {
        "hairstyle": {
          "type": "object",
          "required": ["style", "color"],
          "properties": {
            "style": {
              "type": "string",
              "enum": ["short_spiky", "short_neat", "medium_wavy", "medium_straight",
                       "long_straight", "long_wavy", "ponytail", "braid", "bald",
                       "curly_short", "curly_long", "mohawk", "buns", "undercut",
                       "afro", "dreadlocks"]
            },
            "color": {
              "type": "string",
              "pattern": "^#[0-9a-fA-F]{6}$"
            }
          }
        },
        "eyes": {
          "type": "object",
          "required": ["shape", "color"],
          "properties": {
            "shape": {
              "type": "string",
              "enum": ["round", "almond", "narrow", "wide", "droopy", "cat_eye", "monolid", "hooded"]
            },
            "color": {
              "type": "string",
              "pattern": "^#[0-9a-fA-F]{6}$"
            }
          }
        },
        "skin": {
          "type": "object",
          "required": ["tone"],
          "properties": {
            "tone": {
              "type": "string",
              "pattern": "^#[0-9a-fA-F]{6}$"
            }
          }
        },
        "clothing": {
          "type": "object",
          "required": ["top", "bottom", "shoes"],
          "properties": {
            "top": {
              "type": "object",
              "required": ["style", "color"],
              "properties": {
                "style": {
                  "type": "string",
                  "enum": ["t_shirt", "shirt", "sweater", "hoodie", "overalls", "vest", "flannel", "dress"]
                },
                "color": {
                  "type": "string",
                  "pattern": "^#[0-9a-fA-F]{6}$"
                }
              }
            },
            "bottom": {
              "type": "object",
              "required": ["style", "color"],
              "properties": {
                "style": {
                  "type": "string",
                  "enum": ["jeans", "shorts", "cargo_pants", "skirt", "overalls", "sweatpants"]
                },
                "color": {
                  "type": "string",
                  "pattern": "^#[0-9a-fA-F]{6}$"
                }
              }
            },
            "shoes": {
              "type": "string",
              "enum": ["sneakers", "boots", "sandals", "work_boots"]
            }
          }
        },
        "accessories": {
          "type": "object",
          "properties": {
            "hat": {
              "type": ["string", "null"],
              "enum": [null, "baseball_cap", "straw_hat", "beanie", "cowboy_hat", "sun_hat", "bowler", "witch_hat"]
            },
            "glasses": {
              "type": ["string", "null"],
              "enum": [null, "round_glasses", "square_glasses", "sunglasses", "reading_glasses", "monocle"]
            },
            "jewelry": {
              "type": "array",
              "items": {
                "type": "string",
                "enum": ["earrings_stud", "earrings_hoop", "necklace", "scarf", "watch"]
              },
              "maxItems": 3
            }
          }
        },
        "body_type": {
          "type": "object",
          "required": ["height", "build"],
          "properties": {
            "height": {
              "type": "number",
              "minimum": 0.8,
              "maximum": 1.2
            },
            "build": {
              "type": "number",
              "minimum": 0.7,
              "maximum": 1.3
            }
          }
        }
      }
    },
    "personality": {
      "type": "object",
      "required": ["primary_trait"],
      "properties": {
        "primary_trait": {
          "type": "string",
          "enum": ["friendly", "shy", "energetic", "calm", "curious", "serious", "cheerful", "mysterious"]
        },
        "secondary_traits": {
          "type": "array",
          "items": {
            "type": "string",
            "enum": ["hardworking", "creative", "adventurous", "romantic", "practical", "optimistic", "honest", "generous"]
          },
          "maxItems": 3
        },
        "hobby": {
          "type": "string",
          "enum": ["farming", "fishing", "mining", "foraging", "cooking", "crafting", "socializing", "exploring"]
        }
      }
    },
    "game_state": {
      "type": "object",
      "description": "Runtime game state (not part of character creation)",
      "properties": {
        "position": {
          "type": "object",
          "properties": {
            "x": { "type": "number" },
            "y": { "type": "number" }
          }
        },
        "gold": {
          "type": "integer",
          "minimum": 0,
          "default": 500
        },
        "energy": {
          "type": "integer",
          "minimum": 0,
          "maximum": 100,
          "default": 100
        },
        "friendship": {
          "type": "object",
          "additionalProperties": {
            "type": "integer",
            "minimum": 0,
            "maximum": 2500
          }
        }
      }
    }
  }
}
```

### 3.2 Example Character Instance

```json
{
  "id": "player_001",
  "version": "1.0.0",
  "created_at": "2026-02-12T10:00:00Z",
  "updated_at": "2026-02-12T10:00:00Z",
  "basic_info": {
    "name": "Alex",
    "gender": "male",
    "birthday": {
      "season": "spring",
      "day": 14
    },
    "favorite_thing": "Crystal Gems"
  },
  "appearance": {
    "hairstyle": {
      "style": "short_spiky",
      "color": "#3d2314"
    },
    "eyes": {
      "shape": "round",
      "color": "#3498db"
    },
    "skin": {
      "tone": "#deb887"
    },
    "clothing": {
      "top": {
        "style": "flannel",
        "color": "#c0392b"
      },
      "bottom": {
        "style": "jeans",
        "color": "#4a6fa5"
      },
      "shoes": "work_boots"
    },
    "accessories": {
      "hat": "baseball_cap",
      "glasses": null,
      "jewelry": ["watch"]
    },
    "body_type": {
      "height": 1.0,
      "build": 1.1
    }
  },
  "personality": {
    "primary_trait": "friendly",
    "secondary_traits": ["hardworking", "adventurous"],
    "hobby": "farming"
  },
  "game_state": {
    "position": { "x": 16, "y": 16 },
    "gold": 500,
    "energy": 100,
    "friendship": {}
  }
}
```

---

## 4. UI/UX Design Specifications

### 4.1 Screen Layout

```
+------------------------------------------------------------------+
|  [LOGO]  Character Creator                    [Language: EN/ZH]  |
+------------------------------------------------------------------+
|                                                                   |
|  +------------------------+    +-----------------------------+   |
|  |                        |    |  [Basic Info] [Appearance]  |   |
|  |    CHARACTER           |    |  [Accessories] [Personality]|   |
|  |    PREVIEW             |    +-----------------------------+   |
|  |                        |    |                             |   |
|  |    +-------------+     |    |  Current Category Content:  |   |
|  |    |    [SPRITE] |     |    |                             |   |
|  |    |             |     |    |  +----+ +----+ +----+ +----+|   |
|  |    |   Animated  |     |    |  | 01 | | 02 | | 03 | | 04 ||   |
|  |    |   Preview   |     |    |  +----+ +----+ +----+ +----+|   |
|  |    |             |     |    |  +----+ +----+ +----+ +----+|   |
|  |    +-------------+     |    |  | 05 | | 06 | | 07 | | 08 ||   |
|  |                        |    |  +----+ +----+ +----+ +----+|   |
|  |    [Random] [Reset]    |    |                             |   |
|  +------------------------+    |  Color Palette:             |   |
|                                |  [][ ][ ][ ][ ][ ][ ][ ]    |   |
|                                |                             |   |
|                                +-----------------------------+   |
|                                                                   |
+------------------------------------------------------------------+
|                    [ Back ]          [ Create Character ]         |
+------------------------------------------------------------------+
```

### 4.2 UI Components

#### 4.2.1 Category Navigation
- **Type**: Horizontal tab bar
- **Style**: Pixel-art themed icons with labels
- **States**: Default, Hover, Active, Disabled
- **Animation**: Slide transition between categories

#### 4.2.2 Preview Panel
- **Size**: 280x320 pixels
- **Features**:
  - Animated idle sprite
  - Direction rotation buttons (4 directions)
  - Zoom controls (1x, 2x)
  - Background toggle (solid color / transparent)

#### 4.2.3 Grid Selection
- **Cell Size**: 64x64 pixels
- **Grid Gap**: 8 pixels
- **Selection Style**: Golden border with subtle glow
- **Hover**: Scale 1.05 + shadow

#### 4.2.4 Color Picker
- **Palette Style**: Horizontal bar or grid
- **Cell Size**: 32x32 pixels
- **Custom Picker**: Hex input + color wheel modal
- **Preview**: Real-time update on character

#### 4.2.5 Slider Controls
- **Track Height**: 8 pixels
- **Thumb Size**: 20x20 pixels
- **Preset Buttons**: Quick-select common values
- **Value Display**: Current value shown above slider

### 4.3 Interaction Flow

```
Start
  |
  v
[Select Category Tab]
  |
  v
[View Current Selection]
  |
  +--[Grid Select]--->[Click Option]--->[Update Preview]
  |                                            |
  +--[Color Picker]->[Select/Custom]-->[Update Preview]
  |                                            |
  +--[Slider]------->[Drag/Preset]----->[Update Preview]
  |                                            |
  +--[Text Input]--->[Type]----------->[Validate]
                                               |
                                               v
                                    [All Required Filled?]
                                               |
                                    +----Yes---+---No----+
                                    |                    |
                                    v                    v
                           [Enable Create]      [Disable Create]
                                    |
                                    v
                           [Create Character]
                                    |
                                    v
                           [Save to Profile]
                                    |
                                    v
                                 [Done]
```

### 4.4 Responsive Design

| Screen Size | Layout | Adjustments |
|-------------|--------|-------------|
| Desktop (>1200px) | Side-by-side | Full preview + full options |
| Tablet (768-1200px) | Stacked | Preview above, scrollable options |
| Mobile (<768px) | Full-screen | Tab navigation, modal preview |

### 4.5 Accessibility

- **Keyboard Navigation**: Tab through all elements, Enter to select
- **Screen Reader**: ARIA labels for all interactive elements
- **Color Contrast**: WCAG 2.1 AA compliance (4.5:1 ratio)
- **Focus Indicators**: Visible focus rings on all controls
- **Reduced Motion**: Respect `prefers-reduced-motion`

---

## 5. Localization Support

### 5.1 Language File Structure

```json
{
  "locale": "en",
  "character_creator": {
    "title": "Character Creator",
    "tabs": {
      "basic_info": "Basic Info",
      "appearance": "Appearance",
      "accessories": "Accessories",
      "personality": "Personality"
    },
    "actions": {
      "random": "Random",
      "reset": "Reset",
      "back": "Back",
      "create": "Create Character"
    },
    "validation": {
      "name_required": "Name is required",
      "name_too_long": "Name must be 20 characters or less",
      "name_invalid": "Name can only contain letters"
    }
  }
}
```

```json
{
  "locale": "zh",
  "character_creator": {
    "title": "角色创建器",
    "tabs": {
      "basic_info": "基础信息",
      "appearance": "外观",
      "accessories": "配饰",
      "personality": "性格"
    },
    "actions": {
      "random": "随机",
      "reset": "重置",
      "back": "返回",
      "create": "创建角色"
    },
    "validation": {
      "name_required": "请输入姓名",
      "name_too_long": "姓名不能超过20个字符",
      "name_invalid": "姓名只能包含字母"
    }
  }
}
```

### 5.2 Font Considerations

- **English**: Pixel font (e.g., Press Start 2P, m5x7)
- **Chinese**: Support for CJK characters in pixel style
- **Fallback**: System sans-serif for accessibility

---

## 6. Sprite Asset Requirements

### 6.1 Sprite Sheet Structure

```
sprites/
  character/
    base/
      body_m.png       # Male body base
      body_f.png       # Female body base
      body_nb.png      # Non-binary body base
    hair/
      hair_01.png      # Style variants
      hair_02.png
      ...
    eyes/
      eyes_01.png
      ...
    clothing/
      top_01.png
      bottom_01.png
      shoes_01.png
      ...
    accessories/
      hat_01.png
      glasses_01.png
      jewelry_01.png
      ...
```

### 6.2 Sprite Specifications

- **Frame Size**: 16x16 pixels (base), 32x32 (detailed)
- **Animation Frames**: 4 frames per direction (idle animation)
- **Directions**: 4 (up, down, left, right)
- **Format**: PNG with transparency
- **Color Depth**: 32-bit RGBA

---

*Document Version: 1.0*
*Created by: UI Designer*
*Date: 2026-02-12*
