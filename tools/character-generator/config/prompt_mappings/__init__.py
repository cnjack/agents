# 角色属性定义 - 基于CHARACTER_GENERATOR_DESIGN.md
# 用于SD提示词生成

# ==================== 发型定义 ====================

HAIR_STYLES = {
    "hair_01": {
        "id": "short_spiky",
        "labels": {"en": "Short Spiky", "zh": "短刺头"},
        "prompt": "short spiky hair, textured short hair, anime style"
    },
    "hair_02": {
        "id": "short_neat",
        "labels": {"en": "Short Neat", "zh": "短发整齐"},
        "prompt": "short neat hair, tidy short hair, anime style"
    },
    "hair_03": {
        "id": "medium_wavy",
        "labels": {"en": "Medium Wavy", "zh": "中长波浪"},
        "prompt": "medium length wavy hair, shoulder length waves, anime style"
    },
    "hair_04": {
        "id": "medium_straight",
        "labels": {"en": "Medium Straight", "zh": "中长直发"},
        "prompt": "medium length straight hair, shoulder length straight, anime style"
    },
    "hair_05": {
        "id": "long_straight",
        "labels": {"en": "Long Straight", "zh": "长直发"},
        "prompt": "long straight hair, flowing straight hair, anime style"
    },
    "hair_06": {
        "id": "long_wavy",
        "labels": {"en": "Long Wavy", "zh": "长波浪"},
        "prompt": "long wavy hair, flowing waves, anime style"
    },
    "hair_07": {
        "id": "ponytail",
        "labels": {"en": "Ponytail", "zh": "马尾辫"},
        "prompt": "ponytail, hair tied back, anime style"
    },
    "hair_08": {
        "id": "braid",
        "labels": {"en": "Braid", "zh": "麻花辫"},
        "prompt": "braided hair, plait, anime style"
    },
    "hair_00": {
        "id": "bald",
        "labels": {"en": "Bald", "zh": "光头"},
        "prompt": "bald, no hair"
    },
    "hair_09": {
        "id": "curly_short",
        "labels": {"en": "Curly Short", "zh": "短卷发"},
        "prompt": "short curly hair, tight curls, anime style"
    },
    "hair_10": {
        "id": "curly_long",
        "labels": {"en": "Curly Long", "zh": "长卷发"},
        "prompt": "long curly hair, flowing curls, anime style"
    },
    "hair_11": {
        "id": "mohawk",
        "labels": {"en": "Mohawk", "zh": "莫霍克"},
        "prompt": "mohawk hairstyle, spiked crest, anime style"
    },
    "hair_12": {
        "id": "buns",
        "labels": {"en": "Double Buns", "zh": "双丸子头"},
        "prompt": "double bun hairstyle, odango hair, anime style"
    },
    "hair_13": {
        "id": "undercut",
        "labels": {"en": "Undercut", "zh": "铲青"},
        "prompt": "undercut hairstyle, shaved sides, anime style"
    },
    "hair_14": {
        "id": "afro",
        "labels": {"en": "Afro", "zh": "爆炸头"},
        "prompt": "afro hairstyle, big curly hair, anime style"
    },
    "hair_15": {
        "id": "dreadlocks",
        "labels": {"en": "Dreadlocks", "zh": "脏辫"},
        "prompt": "dreadlocks, dreads, anime style"
    }
}

HAIR_COLORS = {
    "black":       {"hex": "1a1a1a", "prompt": "black hair"},
    "brown_dark":  {"hex": "3d2314", "prompt": "dark brown hair"},
    "brown_light": {"hex": "8b4513", "prompt": "light brown hair"},
    "blonde":      {"hex": "f4d03f", "prompt": "blonde hair, golden hair"},
    "platinum":    {"hex": "e8e4c9", "prompt": "platinum blonde hair, silver hair"},
    "red":         {"hex": "8b2500", "prompt": "red hair, dark red hair"},
    "ginger":      {"hex": "d35400", "prompt": "ginger hair, orange hair"},
    "auburn":      {"hex": "922b21", "prompt": "auburn hair, reddish brown hair"},
    "gray":        {"hex": "808080", "prompt": "gray hair, grey hair"},
    "white":       {"hex": "f5f5f5", "prompt": "white hair"},
    "blue":        {"hex": "3498db", "prompt": "blue hair, anime blue hair"},
    "pink":        {"hex": "ff69b4", "prompt": "pink hair, anime pink hair"},
    "purple":      {"hex": "9b59b6", "prompt": "purple hair, anime purple hair"},
    "green":       {"hex": "27ae60", "prompt": "green hair, anime green hair"},
    "teal":        {"hex": "1abc9c", "prompt": "teal hair, cyan hair"}
}

# ==================== 眼睛定义 ====================

EYE_SHAPES = {
    "eyes_01": {
        "id": "round",
        "labels": {"en": "Round", "zh": "圆眼"},
        "prompt": "round eyes, large round eyes, anime style"
    },
    "eyes_02": {
        "id": "almond",
        "labels": {"en": "Almond", "zh": "杏眼"},
        "prompt": "almond shaped eyes, anime style"
    },
    "eyes_03": {
        "id": "narrow",
        "labels": {"en": "Narrow", "zh": "细长眼"},
        "prompt": "narrow eyes, thin eyes, anime style"
    },
    "eyes_04": {
        "id": "wide",
        "labels": {"en": "Wide", "zh": "大眼"},
        "prompt": "wide eyes, large expressive eyes, anime style"
    },
    "eyes_05": {
        "id": "droopy",
        "labels": {"en": "Droopy", "zh": "下垂眼"},
        "prompt": "droopy eyes, downward slanting eyes, anime style"
    },
    "eyes_06": {
        "id": "cat_eye",
        "labels": {"en": "Cat Eye", "zh": "猫眼"},
        "prompt": "cat eyes, upward slanting eyes, anime style"
    },
    "eyes_07": {
        "id": "monolid",
        "labels": {"en": "Monolid", "zh": "单眼皮"},
        "prompt": "monolid eyes, single eyelid, anime style"
    },
    "eyes_08": {
        "id": "hooded",
        "labels": {"en": "Hooded", "zh": "内双"},
        "prompt": "hooded eyes, partially hidden eyelid, anime style"
    }
}

EYE_COLORS = {
    "brown_dark":  {"hex": "3d2314", "prompt": "dark brown eyes"},
    "brown_light": {"hex": "8b4513", "prompt": "light brown eyes"},
    "hazel":       {"hex": "a67c52", "prompt": "hazel eyes"},
    "amber":       {"hex": "ffbf00", "prompt": "amber eyes, golden eyes"},
    "blue":        {"hex": "3498db", "prompt": "blue eyes"},
    "blue_light":  {"hex": "87ceeb", "prompt": "light blue eyes, sky blue eyes"},
    "green":       {"hex": "27ae60", "prompt": "green eyes"},
    "gray":        {"hex": "808080", "prompt": "gray eyes, grey eyes"},
    "violet":      {"hex": "9b59b6", "prompt": "violet eyes, purple eyes"},
    "red":         {"hex": "c0392b", "prompt": "red eyes, crimson eyes"}
}

# ==================== 肤色定义 ====================

SKIN_TONES = {
    "pale":           {"hex": "ffe4c4", "prompt": "pale skin, fair skin"},
    "fair":           {"hex": "f5deb3", "prompt": "fair skin, light skin"},
    "light":          {"hex": "deb887", "prompt": "light skin"},
    "medium":         {"hex": "d2a679", "prompt": "medium skin tone"},
    "tan":            {"hex": "c68642", "prompt": "tan skin, sun-kissed"},
    "olive":          {"hex": "b5a642", "prompt": "olive skin tone"},
    "brown":          {"hex": "8d5524", "prompt": "brown skin"},
    "dark_brown":     {"hex": "5c4033", "prompt": "dark brown skin"},
    "black":          {"hex": "3d2314", "prompt": "black skin, dark skin"},
    "fantasy_blue":   {"hex": "87ceeb", "prompt": "blue skin, fantasy skin"},
    "fantasy_green":  {"hex": "90ee90", "prompt": "green skin, fantasy skin"},
    "fantasy_purple": {"hex": "dda0dd", "prompt": "purple skin, fantasy skin"}
}

# ==================== 性别定义 ====================

GENDERS = {
    "male": {
        "id": "male",
        "labels": {"en": "Male", "zh": "男"},
        "prompt": "male, man, masculine",
        "sprite_suffix": "_m"
    },
    "female": {
        "id": "female",
        "labels": {"en": "Female", "zh": "女"},
        "prompt": "female, woman, feminine",
        "sprite_suffix": "_f"
    },
    "non_binary": {
        "id": "non_binary",
        "labels": {"en": "Non-binary", "zh": "非二元"},
        "prompt": "androgynous, neutral",
        "sprite_suffix": "_nb"
    }
}

# ==================== 体型定义 ====================

BODY_TYPES = {
    "slim": {
        "height": 0.85,
        "build": 0.8,
        "prompt": "slim body, slender build"
    },
    "average": {
        "height": 1.0,
        "build": 1.0,
        "prompt": "average body, normal build"
    },
    "broad": {
        "height": 1.15,
        "build": 1.2,
        "prompt": "broad body, muscular build"
    }
}

# ==================== 服装定义 ====================

TOP_STYLES = {
    "top_01": {"id": "t_shirt", "labels": {"en": "T-Shirt", "zh": "T恤"}, "prompt": "t-shirt, casual shirt"},
    "top_02": {"id": "shirt", "labels": {"en": "Shirt", "zh": "衬衫"}, "prompt": "button-up shirt, collared shirt"},
    "top_03": {"id": "sweater", "labels": {"en": "Sweater", "zh": "毛衣"}, "prompt": "sweater, knit sweater"},
    "top_04": {"id": "hoodie", "labels": {"en": "Hoodie", "zh": "连帽衫"}, "prompt": "hoodie, hooded sweatshirt"},
    "top_05": {"id": "overalls", "labels": {"en": "Overalls", "zh": "工装背带裤"}, "prompt": "overalls, dungarees"},
    "top_06": {"id": "vest", "labels": {"en": "Vest", "zh": "马甲"}, "prompt": "vest, waistcoat"},
    "top_07": {"id": "flannel", "labels": {"en": "Flannel", "zh": "法兰绒衬衫"}, "prompt": "flannel shirt, plaid shirt"},
    "top_08": {"id": "dress", "labels": {"en": "Dress", "zh": "连衣裙"}, "prompt": "dress, one-piece dress"}
}

BOTTOM_STYLES = {
    "bottom_01": {"id": "jeans", "labels": {"en": "Jeans", "zh": "牛仔裤"}, "prompt": "jeans, denim pants"},
    "bottom_02": {"id": "shorts", "labels": {"en": "Shorts", "zh": "短裤"}, "prompt": "shorts, short pants"},
    "bottom_03": {"id": "cargo_pants", "labels": {"en": "Cargo Pants", "zh": "工装裤"}, "prompt": "cargo pants, utility pants"},
    "bottom_04": {"id": "skirt", "labels": {"en": "Skirt", "zh": "裙子"}, "prompt": "skirt"},
    "bottom_05": {"id": "overalls", "labels": {"en": "Overalls", "zh": "背带裤"}, "prompt": "overalls pants"},
    "bottom_06": {"id": "sweatpants", "labels": {"en": "Sweatpants", "zh": "运动裤"}, "prompt": "sweatpants, jogging pants"}
}

SHOES_STYLES = {
    "shoes_01": {"id": "sneakers", "labels": {"en": "Sneakers", "zh": "运动鞋"}, "prompt": "sneakers, athletic shoes"},
    "shoes_02": {"id": "boots", "labels": {"en": "Boots", "zh": "靴子"}, "prompt": "boots, ankle boots"},
    "shoes_03": {"id": "sandals", "labels": {"en": "Sandals", "zh": "凉鞋"}, "prompt": "sandals, open-toe shoes"},
    "shoes_04": {"id": "work_boots", "labels": {"en": "Work Boots", "zh": "工装靴"}, "prompt": "work boots, heavy boots"}
}

# 服装颜色
CLOTHING_COLORS = {
    "white":      {"hex": "ffffff", "prompt": "white"},
    "black":      {"hex": "1a1a1a", "prompt": "black"},
    "red":        {"hex": "c0392b", "prompt": "red"},
    "blue":       {"hex": "3498db", "prompt": "blue"},
    "green":      {"hex": "27ae60", "prompt": "green"},
    "yellow":     {"hex": "f1c40f", "prompt": "yellow"},
    "orange":     {"hex": "e67e22", "prompt": "orange"},
    "purple":     {"hex": "9b59b6", "prompt": "purple"},
    "pink":       {"hex": "ff69b4", "prompt": "pink"},
    "brown":      {"hex": "8b4513", "prompt": "brown"},
    "gray":       {"hex": "808080", "prompt": "gray"},
    "navy":       {"hex": "1a3a5c", "prompt": "navy blue"}
}

BOTTOM_COLORS = {
    "blue_denim":  {"hex": "4a6fa5", "prompt": "blue denim"},
    "black_denim": {"hex": "2c3e50", "prompt": "black denim"},
    "khaki":       {"hex": "c3b091", "prompt": "khaki"},
    "brown":       {"hex": "8b4513", "prompt": "brown"},
    "white":       {"hex": "f5f5f5", "prompt": "white"},
    "gray":        {"hex": "808080", "prompt": "gray"}
}

# ==================== 配饰定义 ====================

HAT_STYLES = {
    "hat_01": {"id": "baseball_cap", "labels": {"en": "Baseball Cap", "zh": "棒球帽"}, "prompt": "baseball cap"},
    "hat_02": {"id": "straw_hat", "labels": {"en": "Straw Hat", "zh": "草帽"}, "prompt": "straw hat, sun hat"},
    "hat_03": {"id": "beanie", "labels": {"en": "Beanie", "zh": "针织帽"}, "prompt": "beanie, knit cap"},
    "hat_04": {"id": "cowboy_hat", "labels": {"en": "Cowboy Hat", "zh": "牛仔帽"}, "prompt": "cowboy hat, western hat"},
    "hat_05": {"id": "sun_hat", "labels": {"en": "Sun Hat", "zh": "遮阳帽"}, "prompt": "sun hat, wide brim hat"},
    "hat_06": {"id": "bowler", "labels": {"en": "Bowler Hat", "zh": "圆顶礼帽"}, "prompt": "bowler hat, derby hat"},
    "hat_07": {"id": "witch_hat", "labels": {"en": "Witch Hat", "zh": "巫师帽"}, "prompt": "witch hat, pointed hat"}
}

GLASSES_STYLES = {
    "glasses_01": {"id": "round_glasses", "labels": {"en": "Round Glasses", "zh": "圆框眼镜"}, "prompt": "round glasses"},
    "glasses_02": {"id": "square_glasses", "labels": {"en": "Square Glasses", "zh": "方框眼镜"}, "prompt": "square glasses"},
    "glasses_03": {"id": "sunglasses", "labels": {"en": "Sunglasses", "zh": "墨镜"}, "prompt": "sunglasses, shades"},
    "glasses_04": {"id": "reading_glasses", "labels": {"en": "Reading Glasses", "zh": "老花镜"}, "prompt": "reading glasses"},
    "glasses_05": {"id": "monocle", "labels": {"en": "Monocle", "zh": "单片眼镜"}, "prompt": "monocle"}
}

JEWELRY_STYLES = {
    "jewelry_01": {"id": "earrings_stud", "labels": {"en": "Stud Earrings", "zh": "耳钉"}, "prompt": "stud earrings"},
    "jewelry_02": {"id": "earrings_hoop", "labels": {"en": "Hoop Earrings", "zh": "耳环"}, "prompt": "hoop earrings"},
    "jewelry_03": {"id": "necklace", "labels": {"en": "Necklace", "zh": "项链"}, "prompt": "necklace"},
    "jewelry_04": {"id": "scarf", "labels": {"en": "Scarf", "zh": "围巾"}, "prompt": "scarf"},
    "jewelry_05": {"id": "watch", "labels": {"en": "Watch", "zh": "手表"}, "prompt": "wrist watch"}
}

# ==================== 组合数量统计 ====================

COMBINATION_COUNTS = {
    "hair": len(HAIR_STYLES) * len(HAIR_COLORS),      # 16 * 15 = 240
    "eyes": len(EYE_SHAPES) * len(EYE_COLORS),        # 8 * 10 = 80
    "body_base": len(GENDERS) * len(SKIN_TONES) * len(BODY_TYPES),  # 3 * 12 * 3 = 108
    "top": len(TOP_STYLES) * len(CLOTHING_COLORS),    # 8 * 12 = 96
    "bottom": len(BOTTOM_STYLES) * len(BOTTOM_COLORS), # 6 * 6 = 36
    "shoes": len(SHOES_STYLES),                        # 4
    "hat": len(HAT_STYLES),                            # 7
    "glasses": len(GLASSES_STYLES),                    # 5
    "jewelry": len(JEWELRY_STYLES)                     # 5
}

def get_total_assets():
    """获取总资源数量"""
    return sum(COMBINATION_COUNTS.values())

if __name__ == "__main__":
    print("角色属性组合统计:")
    for category, count in COMBINATION_COUNTS.items():
        print(f"  {category}: {count}")
    print(f"总计: {get_total_assets()} 个独立素材")
