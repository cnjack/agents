# 后端开发文档

## 快速开始

### 安装依赖
```bash
cd backend
go mod download
```

### 开发模式
```bash
go run ./cmd/server
# 服务运行在 http://localhost:8080
```

### 构建
```bash
go build -o server ./cmd/server
./server
```

---

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 入口点
│
├── internal/
│   ├── game/                    # 游戏引擎
│   │   ├── engine.go            # 核心引擎
│   │   ├── world.go             # 世界管理
│   │   ├── player.go            # 玩家逻辑
│   │   ├── npc.go               # NPC逻辑
│   │   ├── farm.go              # 农场系统
│   │   ├── quest.go             # 任务系统
│   │   ├── time.go              # 时间系统
│   │   └── friendship.go        # 友好度系统
│   │
│   ├── api/                     # API层
│   │   ├── handler.go           # 处理器
│   │   ├── routes.go            # 路由
│   │   └── middleware.go        # 中间件
│   │
│   ├── models/                  # 数据模型
│   │   ├── state.go             # 状态
│   │   ├── action.go            # 动作
│   │   └── entity.go            # 实体
│   │
│   └── websocket/               # WebSocket
│       └── hub.go               # 连接管理
│
├── go.mod
└── go.sum
```

---

## 核心模块

### 1. 游戏引擎 (Engine)

位置: `internal/game/engine.go`

#### 结构
```go
type Engine struct {
    mu            sync.RWMutex
    state         *models.GameState
    world         *World
    timeSystem    *TimeSystem
    farmSystem    *FarmSystem
    questSystem   *QuestSystem
    friendshipSystem *FriendshipSystem
    ticker        *time.Ticker
    running       bool
    eventChan     chan models.GameEvent
    actionChan    chan *models.Action
}
```

#### 生命周期
```go
// 创建引擎
engine := game.NewEngine()

// 启动游戏循环
engine.Start()

// 停止
engine.Stop()
```

#### 游戏循环
```go
func (e *Engine) tick() {
    // 1. 更新时间
    e.timeSystem.AdvanceTime(&e.state.Time)

    // 2. 更新NPC
    e.updateNPCs()

    // 3. 更新农场
    e.farmSystem.Update(e.state.Time, e.state.Farm)

    // 4. 检查天结束
    if e.state.Time.Hour >= 2 && e.state.Time.Hour < 6 {
        e.state.DayEnded = true
    }
}
```

---

### 2. 世界管理 (World)

位置: `internal/game/world.go`

#### 结构
```go
type World struct {
    config    *models.GameConfig
    tiles     [][]models.Tile
    buildings []models.Building
    areas     []models.MapArea
    trees     []models.Position
}
```

#### 初始化
```go
func (w *World) Initialize(config *models.GameConfig) {
    // 1. 初始化格子
    // 2. 设置农场区域
    // 3. 放置建筑
    // 4. 添加水域
    // 5. 放置树木
    // 6. 定义区域
}
```

#### 行走检测
```go
func (w *World) IsWalkable(x, y int) bool {
    tile := w.GetTile(x, y)
    if tile == nil {
        return false
    }
    return tile.Type != models.TileTypeWater &&
        tile.Type != models.TileTypeBuilding &&
        !tile.Occupied
}
```

---

### 3. 玩家管理 (PlayerManager)

位置: `internal/game/player.go`

#### 移动
```go
func (pm *PlayerManager) Move(player *models.PlayerState, direction models.Direction) error {
    player.Direction = direction

    newPos := player.Position
    switch direction {
    case models.DirectionUp:    newPos.Y--
    case models.DirectionDown:  newPos.Y++
    case models.DirectionLeft:  newPos.X--
    case models.DirectionRight: newPos.X++
    }

    if !pm.world.IsWalkable(newPos.X, newPos.Y) {
        return errors.New("cannot move there")
    }

    player.Position = newPos
    return nil
}
```

#### 使用工具
```go
func (pm *PlayerManager) UseTool(player *models.PlayerState, tool string, farm *models.FarmState) (*models.FarmActionResult, error) {
    if player.Energy < 2 {
        return nil, errors.New("not enough energy")
    }

    targetPos := pm.GetFacingPosition(player)

    switch tool {
    case "hoe":
        // 翻地
    case "watering_can":
        // 浇水
    case "scythe":
        // 收获
    }

    player.Energy -= energyCost
    return result, nil
}
```

---

### 4. 农场系统 (FarmSystem)

位置: `internal/game/farm.go`

#### 作物配置
```go
var cropInfo = map[models.CropType]models.CropInfo{
    models.CropPotato: {
        Type:       models.CropPotato,
        Name:       "Potato",
        GrowthDays: 6,
        Seasons:    []models.Season{models.SeasonSpring},
        BasePrice:  80,
        SeedPrice:  50,
        Stages:     []int{1, 2, 2, 1}, // 每阶段天数
    },
    // ...
}
```

#### 每日更新
```go
func (fs *FarmSystem) NewDayReset(farm *models.FarmState) {
    for y := 0; y < farm.Height; y++ {
        for x := 0; x < farm.Width; x++ {
            tile := &farm.Tiles[y][x]

            // 重置浇水状态
            if tile.Type == models.TileTypeWatered {
                tile.Type = models.TileTypeTilled
            }

            // 作物生长
            if tile.Crop != nil && tile.Crop.WateredToday {
                tile.Crop.DaysGrown++
                fs.updateCropStage(tile.Crop)
            }
            tile.Crop.WateredToday = false
        }
    }
}
```

---

### 5. 任务系统 (QuestSystem)

位置: `internal/game/quest.go`

#### 任务定义
```go
var quests = map[string]*models.QuestData{
    "q001": {
        ID:          "q001",
        Name:        "First Harvest",
        Description: "Harvest your first crop.",
        Objectives: []models.Objective{{
            ID:          "q001_obj1",
            Description: "Harvest any crop",
            Type:        "harvest",
            Target:      "any",
            Count:       1,
        }},
        Rewards: models.Reward{
            Gold: 100,
            Items: []models.InventoryItem{...},
        },
        Status: models.QuestStatusAvailable,
        Giver:  "lewis",
    },
}
```

#### 目标更新
```go
func (qs *QuestSystem) UpdateObjective(questID, objectiveType, target string, amount int) bool {
    quest := qs.quests[questID]

    for i := range quest.Objectives {
        obj := &quest.Objectives[i]
        if obj.Type == objectiveType && obj.Target == target {
            obj.Progress += amount
            if obj.Progress >= obj.Count {
                obj.Completed = true
            }
        }
    }

    qs.checkQuestCompletion(questID)
    return true
}
```

---

### 6. 时间系统 (TimeSystem)

位置: `internal/game/time.go`

#### 时间推进
```go
func (ts *TimeSystem) AdvanceTime(time *models.TimeState) {
    time.Minute++

    if time.Minute >= 60 {
        time.Minute = 0
        time.Hour++

        if time.Hour >= 26 { // 2 AM
            ts.AdvanceDay(time)
        }
    }
}

func (ts *TimeSystem) AdvanceDay(time *models.TimeState) {
    time.Day++
    time.Hour = 6
    time.Minute = 0

    if time.Day > 28 {
        time.Day = 1
        ts.AdvanceSeason(time)
    }
}

func (ts *TimeSystem) AdvanceSeason(time *models.TimeState) {
    switch time.Season {
    case models.SeasonSpring: time.Season = models.SeasonSummer
    case models.SeasonSummer: time.Season = models.SeasonFall
    case models.SeasonFall:   time.Season = models.SeasonWinter
    case models.SeasonWinter:
        time.Season = models.SeasonSpring
        time.Year++
    }
}
```

---

## API处理器

位置: `internal/api/handler.go`

### 处理器结构
```go
type Handler struct {
    engine        *game.Engine
    playerManager *game.PlayerManager
    npcManager    *game.NPCManager
    questSystem   *game.QuestSystem
    farmSystem    *game.FarmSystem
}
```

### 主要处理函数

```go
// 获取游戏状态
func (h *Handler) GetState(c *gin.Context) {
    state := h.engine.GetState()
    c.JSON(http.StatusOK, state)
}

// 执行动作
func (h *Handler) ExecuteAction(c *gin.Context) {
    var action models.Action
    if err := c.ShouldBindJSON(&action); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result := h.processAction(&action)
    c.JSON(http.StatusOK, result)
}

// 重置游戏
func (h *Handler) ResetGame(c *gin.Context) {
    state := h.engine.Reset()
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "state":   state,
    })
}
```

### 动作处理
```go
func (h *Handler) processAction(action *models.Action) *models.ActionResult {
    state := h.engine.GetState()

    switch action.Type {
    case models.ActionMove:
        return h.handleMove(state, action)
    case models.ActionUseTool:
        return h.handleUseTool(state, action)
    case models.ActionPlant:
        return h.handlePlant(state, action)
    case models.ActionHarvest:
        return h.handleHarvest(state, action)
    case models.ActionBuy:
        return h.handleBuy(state, action)
    case models.ActionSell:
        return h.handleSell(state, action)
    case models.ActionTalk:
        return h.handleTalk(state, action)
    case models.ActionGiveGift:
        return h.handleGiveGift(state, action)
    case models.ActionAcceptQuest:
        return h.handleAcceptQuest(state, action)
    case models.ActionWait:
        return h.handleWait(state, action)
    case models.ActionSleep:
        return h.handleSleep(state, action)
    }
}
```

---

## WebSocket管理

位置: `internal/websocket/hub.go`

### Hub结构
```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}
```

### 连接处理
```go
func HandleWebSocket(hub *Hub, engine *game.Engine, c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }

    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
    }

    hub.register <- client

    // 发送初始状态
    state := engine.GetState()
    initialState, _ := json.Marshal(map[string]interface{}{
        "type": "state",
        "data": state,
    })
    client.send <- initialState

    go client.writePump()
    go client.readPump(engine)
}
```

---

## 数据模型

位置: `internal/models/`

### 状态模型 (state.go)
```go
type GameState struct {
    Player    PlayerState
    Time      TimeState
    Farm      FarmState
    NPCs      []NPCState
    Quests    []QuestData
    MapWidth  int
    MapHeight int
    GameOver  bool
    DayEnded  bool
}
```

### 动作模型 (action.go)
```go
type Action struct {
    Type   ActionType
    Params ActionParams
}

type ActionResult struct {
    Success  bool
    Message  string
    NewState *GameState
    Events   []GameEvent
}
```

---

## 配置

### 游戏配置
```go
func DefaultConfig() *models.GameConfig {
    return &models.GameConfig{
        MapWidth:        48,
        MapHeight:       48,
        FarmWidth:       32,
        FarmHeight:      32,
        TickRate:        100,  // ms
        TimeScale:       10,   // game minutes per real second
        StartingGold:    500,
        StartingEnergy:  100,
        EnergyPerAction: 2,
    }
}
```

---

## 扩展开发

### 添加新动作类型

1. 在 `models/action.go` 添加类型:
```go
const ActionNewAction ActionType = "new_action"
```

2. 在 `api/handler.go` 添加处理:
```go
case models.ActionNewAction:
    return h.handleNewAction(state, action)
```

3. 实现处理函数:
```go
func (h *Handler) handleNewAction(state *models.GameState, action *models.Action) *models.ActionResult {
    // 实现逻辑
}
```

### 添加新NPC

在 `game/engine.go` 的 `createNPCs()` 中添加:
```go
{
    ID:           "new_npc",
    Name:         "New NPC",
    Position:     models.Position{X: 20, Y: 20},
    Direction:    models.DirectionDown,
    Friendship:   0,
    MaxFriendship: 250,
    Location:     "town",
    Dialogue:     []string{"Hello!", "Nice to meet you!"},
    Schedule:     []models.ScheduleEntry{...},
}
```

### 添加新任务

在 `game/quest.go` 的 `initQuests()` 中添加:
```go
"q007": {
    ID:          "q007",
    Name:        "New Quest",
    Description: "Quest description.",
    Objectives:  []models.Objective{...},
    Rewards:     models.Reward{...},
    Status:      models.QuestStatusAvailable,
    Giver:       "npc_id",
},
```

---

*文档版本: 1.0*
