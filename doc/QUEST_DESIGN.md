# 任务系统设计文档

## 一、现有任务系统分析

### 1.1 任务类型

| 类型 | 标识 | 描述 | 示例 |
|------|------|------|------|
| 收获 | `harvest` | 收获指定数量的作物 | 收获5个土豆 |
| 对话 | `talk` | 与指定NPC对话 | 与Pierre交谈 |
| 收集 | `collect` | 收集指定数量的资源 | 收集10木材 |
| 送礼 | `give_gift` | 向NPC赠送礼物 | 给Haley送礼 |

### 1.2 现有任务列表

| ID | 名称 | 发布者 | 目标 | 截止(天) | 奖励 |
|----|------|--------|------|----------|------|
| q001 | First Harvest | Lewis | 收获任意作物 x1 | 无 | 100g + 土豆种子x10 |
| q002 | Potato Farmer | Lewis | 收获土豆 x5 | 14 | 250g + 番茄种子x15 |
| q003 | Making Friends | Lewis | 与Pierre对话 | 7 | 50g |
| q004 | Robin's Request | Robin | 收集木材 x10 | 10 | 150g + 箱子x1 |
| q005 | A Gift for Haley | Lewis | 给Haley送礼 | 无 | 100g + 郁金香x3 |
| q006 | Corn for the Community | Lewis | 收获玉米 x10 | 28 | 500g + 南瓜种子x20 |

### 1.3 存在的问题

1. **任务数量不足** - 仅6个任务,游戏后期缺乏内容
2. **缺少日常任务** - 无每日/每周可重复任务,降低玩家留存
3. **任务链缺失** - 任务之间缺乏关联,无成长感
4. **奖励类型单一** - 仅金币和种子,缺少解锁类奖励
5. **无随机任务** - 所有任务静态预设,缺乏惊喜
6. **好感度关联弱** - 任务与NPC关系系统联动不足

---

## 二、NPC日常任务系统

### 2.1 每日委托系统

每日早上6:00刷新,玩家最多接取3个每日委托。

#### Lewis (村长) - 社区服务类

| 任务ID | 名称 | 描述 | 目标 | 奖励 | 刷新权重 |
|--------|------|------|------|------|----------|
| daily_lewis_01 | 镇务帮忙 | 协助处理日常镇务 | 与Lewis对话 | 50g + 好感度+10 | 30% |
| daily_lewis_02 | 社区巡逻 | 巡视小镇各区域 | 访问3个城镇建筑 | 30g | 20% |
| daily_lewis_03 | 居民走访 | 了解居民近况 | 与2名不同NPC对话 | 40g | 25% |
| daily_lewis_04 | 信件投递 | 帮村长送信 | 将信件交给指定NPC | 60g + 好感度+5 | 25% |

#### Pierre (商店老板) - 商业类

| 任务ID | 名称 | 描述 | 目标 | 奖励 | 刷新权重 |
|--------|------|------|------|------|----------|
| daily_pierre_01 | 新鲜货源 | 需要新鲜作物补货 | 出售任意作物x3 | 80g + 商店折扣券(9折) | 40% |
| daily_pierre_02 | 稀有订单 | 客户需要稀有物品 | 出售指定稀有物品x1 | 200g | 15% |
| daily_pierre_03 | 库存盘点 | 整理店铺库存 | 收集指定种子x5 | 50g + 种子x5 | 25% |
| daily_pierre_04 | 市场调研 | 了解竞争对手情况 | 访问Joja超市 | 40g | 20% |

#### Robin (木匠) - 建设类

| 任务ID | 名称 | 描述 | 目标 | 奖励 | 刷新权重 |
|--------|------|------|------|------|----------|
| daily_robin_01 | 材料收集 | 工坊需要原材料 | 收集木材x10 或 石材x10 | 100g + 建筑材料x5 | 35% |
| daily_robin_02 | 工具维护 | 帮忙维护工具 | 与Robin对话(携带工具) | 30g + 工具效率+5%(当天) | 25% |
| daily_robin_03 | 建筑咨询 | 学习建筑知识 | 与Robin对话 | 20g + 建筑经验+10 | 20% |
| daily_robin_04 | 农场测量 | 测量农场面积 | 在农场区域停留并返回 | 50g | 20% |

#### Willy (渔夫) - 钓鱼类

| 任务ID | 名称 | 描述 | 目标 | 奖励 | 刷新权重 |
|--------|------|------|------|------|----------|
| daily_willy_01 | 今日渔获 | 需要新鲜鱼类 | 钓到任意鱼x2 | 60g + 鱼饵x3 | 40% |
| daily_willy_02 | 稀有鱼种 | 寻找特定稀有鱼 | 钓到指定稀有鱼x1 | 300g + 高级鱼饵x1 | 10% |
| daily_willy_03 | 海滩清理 | 清理海滩垃圾 | 在海滩收集物品x5 | 40g | 25% |
| daily_willy_04 | 船只维护 | 帮忙维护渔船 | 与Willy对话 | 30g + 钓鱼技能+5%(当天) | 25% |

#### Haley (村民) - 社交类

| 任务ID | 名称 | 描述 | 目标 | 奖励 | 刷新权重 |
|--------|------|------|------|------|----------|
| daily_haley_01 | 寻找美景 | 发现适合拍照的地点 | 访问海滩或森林 | 40g | 30% |
| daily_haley_02 | 花朵收集 | 需要花朵装饰 | 采集任意花x3 | 50g + 好感度+10 | 25% |
| daily_haley_03 | 时尚建议 | 讨论时尚话题 | 与Haley对话 | 20g | 25% |
| daily_haley_04 | 照片冲洗 | 帮忙冲洗照片 | 收集指定物品(作为交换) | 60g | 20% |

### 2.2 每周特别任务

每周一刷新,周日截止,难度较高但奖励丰厚。

| 任务ID | 名称 | 发布者 | 描述 | 目标 | 奖励 |
|--------|------|--------|------|------|------|
| weekly_001 | 社区贡献 | Lewis | 为社区做出贡献 | 完成任意3个每日委托 | 300g + 社区贡献点x10 |
| weekly_002 | 农场丰收 | Pierre | 大量收购作物 | 出售作物总计价值500g | 500g + 商店VIP卡(永久9折) |
| weekly_003 | 建设周 | Robin | 推进农场建设 | 开垦新土地x20 或 建造设施x1 | 400g + 稀有建筑材料x10 |
| weekly_004 | 钓鱼大赛 | Willy | 参加钓鱼比赛 | 累计钓鱼x10 | 600g + 金色鱼竿(耐久+50%) |
| weekly_005 | 社交达人 | Haley | 提升与村民关系 | 累计好感度+100 | 200g + 全员好感度+5 |

---

## 三、主线任务链设计

### 3.1 任务链总览

```
┌─────────────────────────────────────────────────────────────────┐
│                        主线任务链结构                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  【新手引导线】                                                  │
│  q001 → q002 → q003 → q007 → q008                               │
│                                                                 │
│  【农场发展线】                                                  │
│  q010 → q011 → q012 → q013 → q014                               │
│                                                                 │
│  【社区复兴线】                                                  │
│  q020 → q021 → q022 → q023 → q024 → q025                        │
│                                                                 │
│  【社交关系线】                                                  │
│  q030 → q031 → q032 → q033                                      │
│           ↓                                                     │
│  【婚姻系统线】                                                  │
│  q040 → q041 → q042 → q043                                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 新手引导线

引导玩家熟悉游戏基础操作,逐步解锁游戏功能。

| ID | 名称 | 描述 | 前置条件 | 目标 | 截止 | 奖励 |
|----|------|------|----------|------|------|------|
| q001 | First Harvest | 开始你的农耕之旅 | 无 | 收获任意作物x1 | 无 | 100g + 土豆种子x10 |
| q002 | Potato Farmer | 扩大种植规模 | q001完成 | 收获土豆x5 | 14天 | 250g + 番茄种子x15 |
| q003 | Making Friends | 认识镇上的居民 | q002完成 | 与Pierre对话 | 7天 | 50g + 解锁商店折扣 |
| q007 | Tool Basics | 学习使用工具 | q003完成 | 使用锄头x5, 使用水壶x5 | 7天 | 50g + 铜工具解锁 |
| q008 | First Sale | 学会出售物品 | q007完成 | 向Pierre出售物品x3 | 无 | 100g + 解锁背包扩展 |

### 3.3 农场发展线

引导玩家发展农场,解锁更多功能和建筑。

| ID | 名称 | 描述 | 前置条件 | 目标 | 截止 | 奖励 |
|----|------|------|----------|------|------|------|
| q010 | Farm Expansion | 扩大农场规模 | q002完成 | 开垦农田x50 | 21天 | 500g + 洒水器配方 |
| q011 | Tool Upgrade | 升级你的工具 | q010完成 | 收集铜矿x20 + 支付500g | 无 | 铜锄头(效率+1, 体力消耗-1) |
| q012 | Irrigation System | 建立灌溉系统 | q011完成 | 制作洒水器x3 + 放置在农田 | 14天 | 800g + 质量洒水器配方 |
| q013 | Barn Dreams | 建造畜棚 | q012完成 | 收集木材x200 + 石材x100 + 支付1000g | 无 | 解锁畜牧系统 + 鸡舍建造权限 |
| q014 | Greenhouse Dreams | 拥有温室 | q013完成 + 社区中心进度50% | 完成社区中心4个房间 | 无 | 温室解锁 + 四季种植权限 |

### 3.4 社区复兴线

重建社区中心,解锁游戏核心内容和隐藏区域。

| ID | 名称 | 描述 | 前置条件 | 目标 | 截止 | 奖励 |
|----|------|------|----------|------|------|------|
| q020 | Community Center | 发现社区中心 | q003完成 | 进入社区中心 + 与Lewis对话 | 无 | 解锁社区中心 + 捐赠系统 |
| q021 | Boiler Room | 修复锅炉房 | q020完成 | 捐赠: 铜矿x10 + 铁矿x5 + 石材x50 | 无 | 解锁矿车系统(快速旅行) |
| q022 | Crafts Room | 修复工艺室 | q020完成 | 捐赠: 木材x100 + 硬木x10 + 树脂x5 | 无 | 解锁木匠高级配方 |
| q023 | Pantry | 修复茶水间 | q020完成 | 捐赠: 春季作物x10 + 夏季作物x10 | 无 | 温室修复进度+25% |
| q024 | Fish Tank | 修复鱼缸 | q020完成 | 捐赠: 河鱼x3 + 海鱼x3 + 湖鱼x3 | 无 | 解锁高级鱼竿 + 稀有鱼种 |
| q025 | Vault | 修复金库 | q020完成 | 捐赠: 金币累计10000g | 无 | 解锁商店VIP + 特殊物品折扣 |

**社区中心完成奖励**:
- 温室完全修复
- 解锁废弃矿井
- 解锁秘密森林
- 全员好感度+2心

### 3.5 社交关系线

深化与NPC的关系,解锁个人故事线。

| ID | 名称 | 描述 | 前置条件 | 目标 | 截止 | 奖励 |
|----|------|------|----------|------|------|------|
| q030 | Getting to Know You | 了解小镇居民 | q003完成 | 与所有NPC对话一次 | 28天 | 200g + 解锁NPC日程表 |
| q031 | A Friend in Need | 帮助有困难的居民 | q030完成 + 任意NPC好感度2心 | 完成NPC求助事件 | 14天 | 150g + 好感度+1心(求助NPC) |
| q032 | The Heart of the Valley | 感受小镇的温暖 | q031完成 + 3名NPC好感度4心 | 累计好感度达到500点 | 28天 | 500g + 解锁友谊手镯(好感度+10%) |
| q033 | Valley Legend | 成为小镇传奇 | q032完成 + 5名NPC好感度6心 | 累计好感度达到1500点 | 无 | 1000g + 小镇荣誉市民称号 + 解锁特殊对话 |

### 3.6 婚姻系统线

与心仪NPC发展恋爱关系直至结婚(以Haley为例)。

| ID | 名称 | 描述 | 前置条件 | 目标 | 截止 | 奖励 |
|----|------|------|----------|------|------|------|
| q040_haley | Haley's Heart | 赢得Haley的好感 | q030完成 + Haley好感度2心 | 与Haley好感度达到4心 | 无 | 解锁Haley个人事件 |
| q041_haley | A Stormy Night | 暴风雨之夜的邂逅 | q040完成 + Haley好感度4心 | 在雨天与Haley对话 | 无 | 好感度+1心 + 解锁特殊对话 |
| q042_haley | Photography Contest | 帮助Haley参加摄影比赛 | q041完成 + Haley好感度6心 | 收集: 花朵x5 + 风景点拍照x3 | 14天 | 好感度+1心 + 解锁订婚选项 |
| q043_haley | A Bouquet of Love | 表达你的心意 | q042完成 + Haley好感度8心 | 购买并赠送花束给Haley | 无 | 解锁约会系统 |
| q044_haley | The Mermaid Pendant | 求婚 | q043完成 + Haley好感度10心 | 购买美人鱼吊坠 + 赠送 | 无 | 解锁婚姻系统 + 婚后事件 |

**其他NPC婚姻线**: 每个可婚NPC都有独立的婚姻任务链,结构类似但事件不同。

---

## 四、随机事件任务

### 4.1 季节性节日任务

| 季节 | 节日 | 触发日期 | 任务ID | 名称 | 目标 | 奖励 |
|------|------|----------|--------|------|------|------|
| 春季 | 复活节 | 24日 | event_easter_001 | Egg Hunt | 收集彩蛋x9 | 500g + 草莓种子x20 |
| 夏季 | 夏日宴会 | 11日 | event_summer_001 | Luau Soup | 准备宴会食材(任意作物x5) | 全员好感度+1心 |
| 夏季 | 月光水母舞 | 28日 | event_summer_002 | Moonlight Jellies | 观看水母群(参与活动) | 300g + 珍珠x1 |
| 秋季 | 万灵节 | 27日 | event_fall_001 | Spirit's Eve | 完成迷宫挑战 | 500g + 南瓜灯配方 |
| 秋季 | 星露谷展览会 | 16日 | event_fall_002 | Grange Display | 展示9格高品质物品 | 根据评分奖励金币和星尘 |
| 冬季 | 冰雪节 | 8日 | event_winter_001 | Ice Fishing | 钓到冰鱼x3 | 300g + 冬季鱼饵x10 |
| 冬季 | 冬日星盛宴 | 25日 | event_winter_002 | Secret Santa | 准备神秘礼物送给指定NPC | 全员好感度+1心 |

### 4.2 随机突发事件

#### 天气相关事件

| 事件ID | 名称 | 触发条件 | 描述 | 目标 | 奖励 | 持续时间 |
|--------|------|----------|------|------|------|----------|
| random_storm_001 | Storm Cleanup | 暴风雨后次日 | 清理暴风雨造成的损坏 | 收集倒下的树木x5 | 150g + 硬木x10 | 1天 |
| random_rain_001 | Rainy Day Request | 雨天 | NPC需要雨天的帮助 | 与指定NPC对话 | 好感度+20 | 当天 |

#### NPC生日事件

| NPC | 生日 | 任务ID | 名称 | 目标 | 奖励 |
|-----|------|--------|------|------|------|
| Haley | 春季14日 | birthday_haley | Haley's Birthday | 赠送生日礼物 | 好感度x8倍奖励 |
| Lewis | 春季7日 | birthday_lewis | Lewis's Birthday | 赠送生日礼物 | 好感度x8倍奖励 |
| Pierre | 春季26日 | birthday_pierre | Pierre's Birthday | 赠送生日礼物 | 好感度x8倍奖励 |
| Robin | 夏季21日 | birthday_robin | Robin's Birthday | 赠送生日礼物 | 好感度x8倍奖励 |
| Willy | 夏季10日 | birthday_willy | Willy's Birthday | 赠送生日礼物 | 好感度x8倍奖励 |

#### 特殊NPC事件

| 事件ID | 名称 | 触发条件 | 描述 | 目标 | 奖励 |
|--------|------|----------|------|------|------|
| random_merchant_001 | Traveling Merchant | 随机(周五/周日) | 神秘商人带来稀有物品 | 与商人交易 | 稀有物品 |
| random_mine_001 | Mine Discovery | 进入矿洞 | 发现矿洞深处有奇怪的声音 | 探索矿洞深层 | 解锁新矿层 + 矿石奖励 |
| random_cat_001 | Lost Cat | 随机 | 帮助NPC找回走失的猫 | 找到并归还猫 | 好感度+30 + 猫咪好感 |

### 4.3 好感度触发事件

| NPC | 好感度 | 事件ID | 名称 | 描述 | 奖励 |
|-----|--------|--------|------|------|------|
| Haley | 2心 | heart_haley_2 | A Photo by the Sea | 在海边遇到Haley拍照 | 解锁Haley日常对话变化 |
| Haley | 4心 | heart_haley_4 | Saving the Bracelet | 帮Haley找回手镯 | 好感度+20 + 手镯(饰品) |
| Haley | 6心 | heart_haley_6 | Photography Lesson | 学习摄影技巧 | 解锁相机使用 |
| Haley | 8心 | heart_haley_8 | The Dark Room | 参观Haley的暗房 | 解锁Haley特殊礼物 |
| Lewis | 2心 | heart_lewis_2 | Mayor's Office | 参观市长办公室 | 解锁社区中心任务 |
| Lewis | 4心 | heart_lewis_4 | The Purple Shorts | 发现市长的秘密 | 750g + 特殊对话 |
| Pierre | 2心 | heart_pierre_2 | Shop Storage | 帮忙整理仓库 | 解锁商店折扣 |
| Pierre | 6心 | heart_pierre_6 | Abigail's Secret | 发现Abigail的秘密 | 解锁Abigail任务线 |
| Robin | 2心 | heart_robin_2 | Workshop Tour | 参观木匠工坊 | 解锁建筑蓝图 |
| Willy | 2心 | heart_willy_2 | Fishing Lessons | 学习钓鱼基础 | 解锁钓鱼教程 |

---

## 五、好感度关联任务

### 5.1 好感度解锁任务列表

#### Lewis (村长)

| 心数 | 任务ID | 名称 | 描述 | 目标 | 奖励 |
|------|--------|------|------|------|------|
| 2心 | lewis_q_01 | Mayor's Archives | 整理市长的档案室 | 收集并分类文件 | 100g + 小镇历史书 |
| 4心 | lewis_q_02 | The Purple Shorts | 帮市长找回丢失的物品 | 在特定地点找到紫色短裤 | 750g + 市长的信任 |
| 6心 | lewis_q_03 | Community Garden | 规划社区花园 | 设计并种植指定区域 | 500g + 花园装饰 |
| 8心 | lewis_q_04 | Valley Legacy | 完成小镇历史记录 | 收集所有NPC的故事 | 1000g + 荣誉市民证书 |

#### Pierre (商店老板)

| 心数 | 任务ID | 名称 | 描述 | 目标 | 奖励 |
|------|--------|------|------|------|------|
| 2心 | pierre_q_01 | Stock Rotation | 帮忙清点库存 | 记录所有物品 | 80g + 商店折扣券 |
| 4心 | pierre_q_02 | Competition Research | 调查竞争对手 | 访问Joja超市并记录价格 | 150g + 商业头脑技能 |
| 6心 | pierre_q_03 | Abigail's Favorite | 找到Abigail喜欢的食物 | 准备并赠送特定食物 | 200g + Abigail好感度+1心 |
| 8心 | pierre_q_04 | Family Recipe | 学习Pierre家传食谱 | 收集食材并制作 | 500g + 特殊食谱解锁 |

#### Robin (木匠)

| 心数 | 任务ID | 名称 | 描述 | 目标 | 奖励 |
|------|--------|------|------|------|------|
| 2心 | robin_q_01 | Tool Maintenance | 维护工坊工具 | 收集材料修复工具 | 100g + 工具效率+5% |
| 4心 | robin_q_02 | Demetrius's Research | 协助科学研究 | 收集特定样本 | 200g + 科学知识 |
| 6心 | robin_q_03 | Maru's Invention | 帮Maru完成发明 | 收集电子元件x5 | 300g + 解锁Maru任务线 |
| 8心 | robin_q_04 | Dream House Design | 设计理想之家 | 与Robin讨论并确认设计 | 1000g + 房屋升级折扣 |

#### Willy (渔夫)

| 心数 | 任务ID | 名称 | 描述 | 目标 | 奖励 |
|------|--------|------|------|------|------|
| 2心 | willy_q_01 | First Catch | 学习钓鱼基础 | 钓到任意鱼x5 | 50g + 初级鱼竿 |
| 4心 | willy_q_02 | Secret Spot | 发现秘密钓点 | 访问特定地点 | 100g + 秘密钓点标记 |
| 6心 | willy_q_03 | Legendary Fish | 寻找传说之鱼 | 钓到传说鱼种x1 | 1000g + 金色鱼竿 |
| 8心 | willy_q_04 | Ocean's Guardian | 保护海洋生态 | 清理海滩垃圾x20 | 500g + 海洋守护者称号 |

#### Haley (村民)

| 心数 | 任务ID | 名称 | 描述 | 目标 | 奖励 |
|------|--------|------|------|------|------|
| 2心 | haley_q_01 | Beach Photography | 海滩摄影 | 在海滩找到拍摄点 | 50g + 好感度+20 |
| 4心 | haley_q_02 | Stormy Rescue | 暴风雨救援 | 在暴风雨中帮助Haley | 100g + Haley的感谢信 |
| 6心 | haley_q_03 | Photo Contest | 参加摄影比赛 | 准备并提交照片 | 300g + 相机解锁 |
| 8心 | haley_q_04 | True Feelings | 表达真心 | 与Haley进行真诚对话 | 500g + 解锁恋爱选项 |

### 5.2 好感度奖励汇总

| 心数 | 通用奖励 | 特殊解锁 |
|------|----------|----------|
| 1心 | 解锁基础对话 | 可以进入NPC家中部分区域 |
| 2心 | 解锁日常问候 | 触发个人事件1 |
| 3心 | 好感度加成+5% | 可以查看NPC部分喜好 |
| 4心 | 触发个人事件2 | 解锁好感度任务 |
| 5心 | 好感度加成+10% | 可以进入NPC卧室 |
| 6心 | 触发个人事件3 | 解锁特殊互动 |
| 7心 | 好感度加成+15% | 解锁高级对话选项 |
| 8心 | 触发个人事件4 | 解锁恋爱/挚友选项 |
| 9心 | 好感度加成+20% | 解锁特殊礼物选项 |
| 10心 | 完成个人故事线 | 解锁婚姻/最佳朋友结局 |

---

## 六、奖励系统设计

### 6.1 奖励类型定义

| 类型 | 标识 | 描述 | 示例 |
|------|------|------|------|
| 金币 | `gold` | 直接获得金币 | 100g |
| 物品 | `item` | 获得指定物品 | 土豆种子x10 |
| 好感度 | `friendship` | 提升NPC好感度 | Haley好感度+1心 |
| 配方 | `recipe` | 解锁制作配方 | 洒水器配方 |
| 建筑 | `building` | 解锁建筑权限 | 鸡舍建造 |
| 区域 | `area` | 解锁地图区域 | 秘密森林 |
| 技能 | `skill` | 获得技能加成 | 钓鱼效率+5% |
| 称号 | `title` | 获得特殊称号 | 荣誉市民 |
| 折扣 | `discount` | 获得商店折扣 | 商店永久9折 |

### 6.2 任务星级系统

| 星级 | 奖励倍率 | 额外奖励 | 出现概率 |
|------|----------|----------|----------|
| ★ | 1.0x | 无 | 60% |
| ★★ | 1.5x | 随机物品x1 | 30% |
| ★★★ | 2.0x | 稀有物品x1 + 称号 | 10% |

### 6.3 任务难度与奖励对照表

| 难度 | 目标数量 | 截止时间 | 基础金币 | 经验值 |
|------|----------|----------|----------|--------|
| 简单 | 1-3 | 无限制 | 50-100 | 10 |
| 普通 | 3-5 | 7-14天 | 100-250 | 25 |
| 困难 | 5-10 | 14-21天 | 250-500 | 50 |
| 史诗 | 10-20 | 21-28天 | 500-1000 | 100 |
| 传说 | 20+ | 28天+ | 1000+ | 200 |

### 6.4 特殊奖励详解

#### 配方奖励
```json
{
  "type": "recipe",
  "id": "recipe_sprinkler",
  "name": "洒水器配方",
  "description": "学会制作洒水器,每天自动浇灌周围4格作物",
  "materials": ["铜矿x1", "铁矿石x1"],
  "unlocked_by": "q010"
}
```

#### 建筑解锁
```json
{
  "type": "building_unlock",
  "id": "building_coop",
  "name": "鸡舍建造权限",
  "description": "可以在Robin处建造鸡舍",
  "cost": "木材x200 + 石材x100 + 1000g",
  "unlocked_by": "q013"
}
```

#### 技能奖励
```json
{
  "type": "skill",
  "id": "skill_fishing_boost",
  "name": "钓鱼精通",
  "description": "钓鱼效率永久+5%",
  "effect": "fishing_efficiency +0.05",
  "unlocked_by": "willy_q_01"
}
```

---

## 七、数据结构定义

### 7.1 Go Struct 定义

```go
package models

import "time"

// ============================================================
// 基础任务结构
// ============================================================

// QuestStatus 任务状态
type QuestStatus string

const (
    QuestStatusAvailable QuestStatus = "available" // 可接取
    QuestStatusActive    QuestStatus = "active"    // 进行中
    QuestStatusCompleted QuestStatus = "completed" // 已完成
    QuestStatusFailed    QuestStatus = "failed"    // 已失败
    QuestStatusLocked    QuestStatus = "locked"    // 未解锁
)

// ObjectiveType 目标类型
type ObjectiveType string

const (
    ObjectiveHarvest   ObjectiveType = "harvest"    // 收获作物
    ObjectiveTalk      ObjectiveType = "talk"       // 对话
    ObjectiveCollect   ObjectiveType = "collect"    // 收集资源
    ObjectiveGiveGift  ObjectiveType = "give_gift"  // 送礼
    ObjectiveFish      ObjectiveType = "fish"       // 钓鱼
    ObjectiveMine      ObjectiveType = "mine"       // 挖矿
    ObjectiveCraft     ObjectiveType = "craft"      // 制作
    ObjectiveBuild     ObjectiveType = "build"      // 建造
    ObjectiveExplore   ObjectiveType = "explore"    // 探索
    ObjectiveVisit     ObjectiveType = "visit"      // 访问地点
)

// QuestData 任务数据
type QuestData struct {
    ID              string       `json:"id"`
    Name            string       `json:"name"`
    Description     string       `json:"description"`
    Type            string       `json:"type"`              // main, daily, weekly, event, friendship
    Objectives      []Objective  `json:"objectives"`
    Rewards         Reward       `json:"rewards"`
    Status          QuestStatus  `json:"status"`
    Giver           string       `json:"giver"`             // NPC ID
    Deadline        int          `json:"deadline"`          // 截止天数, 0表示无限制
    Prerequisites   []string     `json:"prerequisites"`     // 前置任务ID列表
    StarLevel       int          `json:"star_level"`        // 1-3星级
    Difficulty      string       `json:"difficulty"`        // easy, normal, hard, epic, legendary
    ChainID         string       `json:"chain_id"`          // 所属任务链ID
    ChainOrder      int          `json:"chain_order"`       // 任务链中的顺序
    SeasonLimit     []Season     `json:"season_limit"`      // 季节限定
    FriendshipReq   int          `json:"friendship_req"`    // 好感度需求(心数)
    Repeatable      bool         `json:"repeatable"`        // 是否可重复
    CooldownDays    int          `json:"cooldown_days"`     // 冷却天数
    LastCompletedAt *time.Time   `json:"last_completed_at"`
}

// Objective 任务目标
type Objective struct {
    ID          string        `json:"id"`
    Description string        `json:"description"`
    Type        ObjectiveType `json:"type"`
    Target      string        `json:"target"`       // 目标物品/NPC/地点ID
    Count       int           `json:"count"`        // 需要数量
    Progress    int           `json:"progress"`     // 当前进度
    Completed   bool          `json:"completed"`
    Optional    bool          `json:"optional"`     // 是否可选目标
}

// Reward 任务奖励
type Reward struct {
    Gold           int              `json:"gold"`
    Items          []InventoryItem  `json:"items"`
    Friendship     []FriendshipReward `json:"friendship"`
    Recipes        []string         `json:"recipes"`       // 解锁配方ID列表
    Buildings      []string         `json:"buildings"`     // 解锁建筑ID列表
    Areas          []string         `json:"areas"`         // 解锁区域ID列表
    Skills         []SkillReward    `json:"skills"`
    Titles         []string         `json:"titles"`        // 获得称号ID列表
    Discounts      []DiscountReward `json:"discounts"`
    Experience     int              `json:"experience"`
}

// FriendshipReward 好感度奖励
type FriendshipReward struct {
    NPCID    string `json:"npc_id"`
    Points   int    `json:"points"`    // 直接加点数
    Hearts   int    `json:"hearts"`    // 加心数
    AllNPC   bool   `json:"all_npc"`   // 是否作用于所有NPC
}

// SkillReward 技能奖励
type SkillReward struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Effect      string  `json:"effect"`
    Value       float64 `json:"value"`
    Duration    int     `json:"duration"`    // 持续天数, 0表示永久
    Description string  `json:"description"`
}

// DiscountReward 折扣奖励
type DiscountReward struct {
    ShopID    string  `json:"shop_id"`
    Discount  float64 `json:"discount"`   // 0.9 = 9折
    Permanent bool    `json:"permanent"`
    Duration  int     `json:"duration"`   // 持续天数
}

// ============================================================
// 日常任务系统
// ============================================================

// DailyQuestTemplate 每日任务模板
type DailyQuestTemplate struct {
    ID             string        `json:"id"`
    Name           string        `json:"name"`
    Description    string        `json:"description"`
    GiverNPC       string        `json:"giver_npc"`
    Objectives     []ObjectiveTemplate `json:"objectives"`
    BaseReward     Reward        `json:"base_reward"`
    StarLevel      int           `json:"star_level"`
    Probability    float64       `json:"probability"`    // 刷新概率 0-1
    SeasonLimit    []Season      `json:"season_limit"`
    WeatherLimit   []Weather     `json:"weather_limit"`  // 天气限制
    FriendshipReq  int           `json:"friendship_req"`
    MaxPerDay      int           `json:"max_per_day"`    // 每天最多出现次数
}

// ObjectiveTemplate 目标模板(用于生成随机目标)
type ObjectiveTemplate struct {
    Type        ObjectiveType `json:"type"`
    TargetPool  []string      `json:"target_pool"`  // 目标池,随机选择
    CountMin    int           `json:"count_min"`
    CountMax    int           `json:"count_max"`
    Description string        `json:"description_template"`
}

// WeeklyQuestTemplate 每周任务模板
type WeeklyQuestTemplate struct {
    ID            string       `json:"id"`
    Name          string       `json:"name"`
    Description   string       `json:"description"`
    GiverNPC      string       `json:"giver_npc"`
    Objectives    []Objective  `json:"objectives"`
    Rewards       Reward       `json:"rewards"`
    Difficulty    string       `json:"difficulty"`
    RefreshDay    int          `json:"refresh_day"`     // 每周几刷新 (1=周一)
    DeadlineDays  int          `json:"deadline_days"`   // 截止天数
}

// ============================================================
// 任务链系统
// ============================================================

// QuestChain 任务链定义
type QuestChain struct {
    ID              string   `json:"id"`
    Name            string   `json:"name"`
    Description     string   `json:"description"`
    Type            string   `json:"type"`            // tutorial, farm, community, social, marriage
    QuestIDs        []string `json:"quest_ids"`       // 任务ID列表(按顺序)
    CurrentProgress int      `json:"current_progress"` // 当前进度索引
    Completed       bool     `json:"completed"`
    FinalReward     Reward   `json:"final_reward"`    // 完成全部奖励
    Unlocks         []string `json:"unlocks"`         // 解锁的其他任务链
}

// ============================================================
// 随机事件系统
// ============================================================

// RandomEventType 随机事件类型
type RandomEventType string

const (
    EventTypeWeather    RandomEventType = "weather"    // 天气触发
    EventTypeDate       RandomEventType = "date"       // 日期触发
    EventTypeFriendship RandomEventType = "friendship" // 好感度触发
    EventTypeLocation   RandomEventType = "location"   // 地点触发
    EventTypeRandom     RandomEventType = "random"     // 随机触发
)

// RandomEvent 随机事件
type RandomEvent struct {
    ID            string          `json:"id"`
    Name          string          `json:"name"`
    Description   string          `json:"description"`
    Type          RandomEventType `json:"type"`
    TriggerType   string          `json:"trigger_type"`
    TriggerValue  string          `json:"trigger_value"`
    Quest         QuestData       `json:"quest"`
    Duration      int             `json:"duration"`       // 持续天数
    Probability   float64         `json:"probability"`    // 触发概率
    CooldownDays  int             `json:"cooldown_days"`
    LastTriggered *time.Time      `json:"last_triggered"`
}

// SeasonalEvent 季节性事件
type SeasonalEvent struct {
    ID           string     `json:"id"`
    Name         string     `json:"name"`
    Description  string     `json:"description"`
    Season       Season     `json:"season"`
    Day          int        `json:"day"`            // 季节中的日期
    Quest        QuestData  `json:"quest"`
    Location     string     `json:"location"`       // 触发地点
    StartTime    string     `json:"start_time"`     // 开始时间 "09:00"
    EndTime      string     `json:"end_time"`       // 结束时间 "18:00"
    RepeatYearly bool       `json:"repeat_yearly"`  // 是否每年重复
}

// ============================================================
// 好感度关联系统
// ============================================================

// FriendshipQuest 好感度关联任务
type FriendshipQuest struct {
    ID             string     `json:"id"`
    NPCID          string     `json:"npc_id"`
    RequiredHearts int        `json:"required_hearts"`
    Quest          QuestData  `json:"quest"`
    EventTrigger   string     `json:"event_trigger"`  // 触发事件ID
    UnlocksNext    []string   `json:"unlocks_next"`   // 解锁的下一级任务
}

// NPCBirthday NPC生日
type NPCBirthday struct {
    NPCID       string `json:"npc_id"`
    Season      Season `json:"season"`
    Day         int    `json:"day"`
    QuestID     string `json:"quest_id"`      // 生日任务ID
    GiftBonus   int    `json:"gift_bonus"`    // 礼物好感度倍率
}

// ============================================================
// 任务系统管理器
// ============================================================

// QuestSystemState 任务系统状态
type QuestSystemState struct {
    ActiveQuests       []string          `json:"active_quests"`
    CompletedQuests    []string          `json:"completed_quests"`
    AvailableQuests    []string          `json:"available_quests"`
    DailyQuests        []string          `json:"daily_quests"`        // 当天日常任务
    WeeklyQuest        string            `json:"weekly_quest"`        // 当周周常任务
    ActiveEvents       []string          `json:"active_events"`       // 当前活跃事件
    ChainProgress      map[string]int    `json:"chain_progress"`      // 任务链进度
    LastDailyRefresh   time.Time         `json:"last_daily_refresh"`
    LastWeeklyRefresh  time.Time         `json:"last_weekly_refresh"`
    CompletedCount     int               `json:"completed_count"`     // 累计完成任务数
}

// QuestProgressUpdate 任务进度更新
type QuestProgressUpdate struct {
    QuestID       string `json:"quest_id"`
    ObjectiveID   string `json:"objective_id"`
    ProgressDelta int    `json:"progress_delta"`
    Completed     bool   `json:"completed"`
    QuestComplete bool   `json:"quest_complete"`
}

// QuestActionResult 任务行动结果
type QuestActionResult struct {
    Success      bool          `json:"success"`
    QuestID      string        `json:"quest_id"`
    QuestName    string        `json:"quest_name"`
    Message      string        `json:"message"`
    Rewards      *Reward       `json:"rewards,omitempty"`
    NewQuests    []string      `json:"new_quests,omitempty"`    // 解锁的新任务
    Unlocks      []string      `json:"unlocks,omitempty"`       // 其他解锁内容
}
```

### 7.2 API 接口定义

```go
// 任务相关API
GET  /api/v1/quests                    // 获取所有可接取任务
GET  /api/v1/quests/active             // 获取进行中的任务
GET  /api/v1/quests/completed          // 获取已完成任务
GET  /api/v1/quests/{id}               // 获取任务详情
POST /api/v1/quests/{id}/accept        // 接受任务
POST /api/v1/quests/{id}/complete      // 完成任务
POST /api/v1/quests/{id}/abandon       // 放弃任务

// 日常任务API
GET  /api/v1/quests/daily              // 获取今日日常任务
POST /api/v1/quests/daily/refresh      // 刷新日常任务(消耗道具)

// 任务链API
GET  /api/v1/quests/chains             // 获取所有任务链
GET  /api/v1/quests/chains/{id}        // 获取任务链详情

// 随机事件API
GET  /api/v1/events                    // 获取当前活跃事件
POST /api/v1/events/{id}/trigger       // 触发事件

// 好感度任务API
GET  /api/v1/npc/{id}/quests           // 获取NPC关联任务
```

---

## 八、实现优先级

### Phase 1: 基础扩展 (高优先级)
1. 扩展现有任务类型
2. 实现任务前置条件检查
3. 实现任务链系统
4. 添加更多主线任务

### Phase 2: 日常任务系统 (高优先级)
1. 实现每日任务模板系统
2. 实现每日刷新逻辑
3. 实现每周特别任务
4. 实现任务星级系统

### Phase 3: 随机事件系统 (中优先级)
1. 实现季节性节日任务
2. 实现随机突发事件
3. 实现NPC生日系统
4. 实现天气关联事件

### Phase 4: 好感度关联 (中优先级)
1. 实现好感度解锁任务
2. 实现NPC个人事件
3. 实现婚姻任务链
4. 实现好感度奖励系统

### Phase 5: 奖励扩展 (低优先级)
1. 实现配方解锁系统
2. 实现建筑解锁系统
3. 实现区域解锁系统
4. 实现称号系统

---

*文档版本: 2.0*
*创建者: Game Designer*
*日期: 2026-02-12*
*更新记录: 完整任务系统设计,包含日常任务、任务链、随机事件、好感度关联*
