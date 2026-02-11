# AI驱动游戏玩法设计规范

## 一、设计目标

将AI深度集成到游戏核心,让AI成为:
- **NPC的"大脑"** - 动态生成对话,智能决定行为
- **玩家的"翻译官"** - 理解自然语言,转换为游戏动作
- **世界的"叙事者"** - 生成动态事件和故事

## 二、核心系统设计

### 2.1 AI对话系统 (ConversationEngine)

#### 设计原理
```
玩家交互触发 → 构建对话上下文 → AI生成回复 → 返回游戏
```

#### 新增数据结构

```go
// NPC人格档案
type NPCPersonality struct {
    ID              string            `json:"id"`
    Name            string            `json:"name"`
    Role            string            `json:"role"`        // 市长、店主、木匠等
    Traits          []string          `json:"traits"`      // 性格特质: 开朗、害羞、务实、傲慢
    Background      string            `json:"background"`  // 背景故事
    Values          []string          `json:"values"`      // 重视的事物
    Dislikes        []string          `json:"dislikes"`    // 讨厌的事物
    SpeechStyle     string            `json:"speech_style"`// 说话风格描述
    Hobby           string            `json:"hobby"`       // 爱好
    Secrets         []string          `json:"secrets"`     // 秘密(高友谊度解锁)
}

// 对话上下文
type DialogueContext struct {
    NPC             NPCPersonality    `json:"npc"`
    Time            TimeState         `json:"time"`          // 时间、季节
    Weather         Weather           `json:"weather"`       // 天气
    Friendship      int               `json:"friendship"`    // 友谊等级
    Hearts          int               `json:"hearts"`        // 心数
    PlayerName      string            `json:"player_name"`
    RecentInteractions []Interaction  `json:"recent_interactions"` // 最近互动
    RecentEvents    []string          `json:"recent_events"` // 最近游戏事件
    CurrentMood     string            `json:"current_mood"`  // 当前心情
    Location        string            `json:"location"`      // 当前位置
    PlayerInput     string            `json:"player_input"`  // 玩家输入的自然语言
}

// 互动记录
type Interaction struct {
    Type        string    `json:"type"`        // talk, gift, quest
    Timestamp   int64     `json:"timestamp"`
    Summary     string    `json:"summary"`     // 互动摘要
    Result      string    `json:"result"`      // 结果
}

// AI对话响应
type AIDialogueResponse struct {
    NPCID           string   `json:"npc_id"`
    Dialogue        string   `json:"dialogue"`        // 生成的对话
    MoodChange      string   `json:"mood_change"`     // 心情变化
    FriendshipDelta int      `json:"friendship_delta"`// 友谊变化
    SuggestedActions []string `json:"suggested_actions"` // 建议的后续动作
    TriggerEvent    string   `json:"trigger_event"`   // 可能触发的事件
}
```

#### NPC人格档案示例

```json
{
  "id": "lewis",
  "name": "Lewis",
  "role": "mayor",
  "traits": ["friendly", "responsible", "nostalgic", "secretive"],
  "background": "Stardew Valley的市长,服务小镇多年。年轻时也曾是农民。",
  "values": ["community", "tradition", "hard_work"],
  "dislikes": ["laziness", "disrespect"],
  "speech_style": "热情但略带官腔,喜欢用'我们小镇'这类词",
  "hobby": "园艺,喜欢紫色短裤(秘密)",
  "secrets": ["有一条紫色短裤", "对Jodi有好感"]
}
```

```json
{
  "id": "haley",
  "name": "Haley",
  "role": "villager",
  "traits": ["superficial", "photography_lover", "gradually_warming"],
  "background": "来自城市的女孩,和姐姐Emily住在一起。一开始看不起农村生活。",
  "values": ["beauty", "fashion", "photography"],
  "dislikes": ["dirt", "bugs", "farming"],
  "speech_style": "傲慢但逐渐软化,高友谊度会变得温柔",
  "hobby": "摄影",
  "secrets": ["其实很寂寞", "喜欢椰子"]
}
```

---

### 2.2 AI行为决策系统 (BehaviorEngine)

#### 设计原理
```
游戏tick → AI评估NPC状态 → 决定下一步行为 → 执行移动/动作
```

#### 新增数据结构

```go
// NPC行为状态
type NPCBehaviorState struct {
    NPCID           string        `json:"npc_id"`
    CurrentGoal     string        `json:"current_goal"`     // 当前目标
    CurrentMood     string        `json:"current_mood"`     // 当前心情
    Energy          int           `json:"energy"`           // 精力值
    SocialNeed      int           `json:"social_need"`      // 社交需求
    LastDecision    string        `json:"last_decision"`    // 上次决策
    DecisionReason  string        `json:"decision_reason"`  // 决策原因
}

// AI行为决策请求
type AIBehaviorRequest struct {
    NPCID           string              `json:"npc_id"`
    Personality     NPCPersonality      `json:"personality"`
    BehaviorState   NPCBehaviorState    `json:"behavior_state"`
    GameState       GameState           `json:"game_state"`
    PlayerNearby    bool                `json:"player_nearby"`
    Weather         Weather             `json:"weather"`
    Time            TimeState           `json:"time"`
    PossibleActions []string            `json:"possible_actions"`
}

// AI行为决策响应
type AIBehaviorResponse struct {
    NPCID           string   `json:"npc_id"`
    Action          string   `json:"action"`           // move, stay, interact, talk
    TargetLocation  string   `json:"target_location"`  // 目标位置
    TargetPosition  Position `json:"target_position"`  // 目标坐标
    Reason          string   `json:"reason"`           // 决策原因
    MoodAfterAction string   `json:"mood_after_action"`// 行动后心情
    CanInterrupt    bool     `json:"can_interrupt"`    // 是否可被玩家打断
}
```

#### 决策因素权重

| 因素 | 权重 | 影响描述 |
|------|------|----------|
| 天气 | 高 | 下雨天倾向室内活动 |
| 友谊度 | 高 | 高友谊度NPC更愿意主动接近玩家 |
| 时间 | 中 | 影响NPC的日常作息 |
| 社交需求 | 中 | 长时间不互动会想去公共场所 |
| 心情 | 中 | 心情差可能不愿互动 |
| 个人目标 | 低 | NPC的职业/爱好目标 |

---

### 2.3 自然语言理解系统 (NLPEngine)

#### 设计原理
```
玩家输入自然语言 → AI理解意图 → 转换为游戏动作 → 执行
```

#### 新增数据结构

```go
// 自然语言理解请求
type NLPRequest struct {
    PlayerInput     string      `json:"player_input"`
    GameState       GameState   `json:"game_state"`
    NearbyNPCs      []NPCState  `json:"nearby_npcs"`
    PlayerInventory []InventoryItem `json:"player_inventory"`
    Context         string      `json:"context"`         // 当前上下文
}

// 自然语言理解响应
type NLPResponse struct {
    Understood      bool        `json:"understood"`
    Intent          string      `json:"intent"`          // talk, gift, buy, sell, move, etc.
    Action          Action      `json:"action"`          // 转换后的游戏动作
    TargetNPC       string      `json:"target_npc"`      // 目标NPC(如果有)
    TargetItem      string      `json:"target_item"`     // 目标物品(如果有)
    ResponseText    string      `json:"response_text"`   // AI生成的回应文本
    NeedsClarification bool     `json:"needs_clarification"` // 是否需要澄清
    ClarificationQuestion string `json:"clarification_question"` // 澄清问题
}
```

#### 支持的自然语言类型

| 玩家输入示例 | 识别意图 | 转换动作 |
|-------------|----------|----------|
| "你好啊刘易斯" | 问候/对话 | talk(npc: lewis) |
| "把这个花送给Haley" | 送礼 | give_gift(npc: haley, item: flower) |
| "我想买点种子" | 购买 | interact(shop) |
| "今天天气真好" | 闲聊 | talk(context: weather) |
| "你在干什么?" | 询问 | talk(context: current_activity) |
| "这鱼怎么卖?" | 交易 | interact(shop, item: fish) |

---

### 2.4 NPC心情系统 (MoodSystem)

#### 新增数据结构

```go
// NPC心情状态
type NPCMood struct {
    NPCID           string    `json:"npc_id"`
    CurrentMood     string    `json:"current_mood"`   // happy, neutral, sad, angry, excited
    MoodValue       int       `json:"mood_value"`     // -100 to 100
    MoodFactors     []MoodFactor `json:"mood_factors"`
    LastMoodChange  int64     `json:"last_mood_change"`
}

// 心情影响因子
type MoodFactor struct {
    Factor          string    `json:"factor"`         // weather, player_interaction, event
    Impact          int       `json:"impact"`         // -100 to 100
    Duration        int       `json:"duration"`       // 持续时间(游戏小时)
    Description     string    `json:"description"`
}

// 心情对对话的影响
var MoodDialogueModifiers = map[string]string{
    "happy":    "语气愉快,愿意分享更多信息",
    "neutral":  "正常语气",
    "sad":      "语气低落,回复简短",
    "angry":    "语气生硬,可能拒绝互动",
    "excited":  "语气兴奋,可能主动分享秘密",
}
```

#### 心情影响因素

| 因素 | 影响 | 持续时间 |
|------|------|----------|
| 收到喜欢的礼物 | +30 ~ +50 | 24小时 |
| 收到讨厌的礼物 | -30 ~ -50 | 24小时 |
| 下雨天(某些NPC) | -10 ~ -20 | 整个雨天 |
| 节日 | +20 ~ +40 | 节日当天 |
| 玩家长时间不互动 | -5 每天 | 累积 |
| 生日收到礼物 | +50 ~ +80 | 3天 |
| 完成NPC委托的任务 | +20 ~ +40 | 48小时 |

---

## 三、API设计

### 3.1 新增API端点

```
POST /api/v1/ai/dialogue
请求: DialogueContext
响应: AIDialogueResponse

POST /api/v1/ai/behavior
请求: AIBehaviorRequest
响应: AIBehaviorResponse

POST /api/v1/ai/interpret
请求: NLPRequest
响应: NLPResponse

POST /api/v1/ai/narrate
请求: EventContext
响应: NarrativeText

GET /api/v1/npc/{id}/mood
响应: NPCMood
```

### 3.2 WebSocket消息类型

```json
// AI生成对话
{
  "type": "ai_dialogue",
  "data": {
    "npc_id": "lewis",
    "dialogue": "今天的天气让我想起了老农场...",
    "mood": "nostalgic"
  }
}

// NPC行为变化
{
  "type": "npc_action",
  "data": {
    "npc_id": "haley",
    "action": "move_to",
    "target": "beach",
    "reason": "想去拍照"
  }
}

// 自然语言理解结果
{
  "type": "nlp_result",
  "data": {
    "understood": true,
    "intent": "talk",
    "response": "Lewis转过身来..."
  }
}
```

---

## 四、实现优先级

### Phase 1: 基础AI对话 (优先级: 高)
1. 创建NPC人格档案数据结构
2. 实现DialogueContext构建器
3. 集成LLM API
4. 实现基础对话生成

### Phase 2: 自然语言交互 (优先级: 高)
1. 实现NLP请求处理
2. 创建意图识别系统
3. 连接游戏动作执行

### Phase 3: AI行为决策 (优先级: 中)
1. 实现NPCBehaviorState管理
2. 创建AI行为决策引擎
3. 集成到游戏主循环

### Phase 4: 心情系统 (优先级: 中)
1. 实现MoodSystem核心逻辑
2. 连接心情影响因子
3. 集成到对话和行为系统

---

## 五、与Engineer协作要点

### 需要Engineer实现:
1. **LLM API集成层** - 调用AI模型的封装
2. **Context构建器** - 将游戏状态转为AI prompt
3. **响应解析器** - 解析AI返回的结构化数据
4. **缓存机制** - 减少重复AI调用
5. **异步处理** - 不阻塞游戏主循环

### 需要Game Designer提供:
1. NPC人格档案内容(每个NPC的详细设定)
2. 心情影响因素的具体数值
3. 对话风格指南
4. 行为决策规则

---

## 六、示例场景流程

### 场景: 玩家与Lewis自然对话

```
1. 玩家输入: "刘易斯市长,今天有什么需要帮忙的吗?"

2. NLP解析:
   - intent: talk
   - target_npc: lewis
   - context: asking_for_help

3. 构建DialogueContext:
   - NPC: Lewis (市长,热情,怀旧)
   - Time: Spring 5, 10:30 AM
   - Weather: sunny
   - Friendship: 2 hearts
   - RecentEvents: ["completed_quest_q001"]
   - Mood: happy

4. AI生成对话:
   "哦,我们年轻的农民!你刚完成了第一次收获,
    做得很好!今天社区中心有些事情需要处理,
    你有兴趣帮忙吗?"

5. 返回结果:
   - dialogue: (上述内容)
   - mood_change: "excited"
   - suggested_actions: ["accept_quest", "ask_about_community_center"]
```

---

*文档版本: 1.0*
*创建者: Game Designer*
*日期: 2026-02-11*
