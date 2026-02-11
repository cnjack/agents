# 星露谷物语风格网页游戏 - 系统架构文档

## 目录
- [项目概述](#项目概述)
- [技术栈](#技术栈)
- [系统架构](#系统架构)
- [前端架构](#前端架构)
- [后端架构](#后端架构)
- [游戏系统](#游戏系统)
- [API设计](#api设计)
- [数据模型](#数据模型)
- [部署架构](#部署架构)

---

## 项目概述

### 目标
创建一个基于Vue.js的像素风格农场模拟游戏，作为AI Agent智能体的开发和测试环境。Agent通过RESTful API控制游戏角色执行各种操作。

### 核心特性
- 像素风格2D渲染 (Canvas 2D)
- 实时状态同步 (WebSocket)
- AI Agent API接口
- 完整的游戏系统 (农场、NPC、任务、时间)

---

## 技术栈

### 前端
| 技术 | 版本 | 用途 |
|------|------|------|
| Vue.js | 3.4+ | 前端框架 |
| TypeScript | 5.3+ | 类型安全 |
| Pinia | 2.1+ | 状态管理 |
| Vite | 5.0+ | 构建工具 |
| Canvas 2D | - | 游戏渲染 |
| TailwindCSS | 3.4+ | UI样式 |
| WebSocket | - | 实时通信 |

### 后端
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端语言 |
| Gin | 1.9+ | Web框架 |
| Gorilla WebSocket | 1.5+ | WebSocket支持 |

### 部署
| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 服务编排 |
| Nginx | 反向代理 |

---

## 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        用户界面层                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Vue 组件   │  │  Canvas渲染 │  │  WebSocket  │             │
│  │  (UI界面)   │  │  (游戏画面) │  │  (实时同步) │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        状态管理层                                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Pinia Store                           │   │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐              │   │
│  │  │ gameStore │ │playerStore│ │ uiStore   │              │   │
│  │  └───────────┘ └───────────┘ └───────────┘              │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API通信层                                 │
│  ┌─────────────┐  ┌─────────────┐                              │
│  │  REST API   │  │  WebSocket  │                              │
│  │  (动作执行) │  │  (状态推送) │                              │
│  └─────────────┘  └─────────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        后端服务层                                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Gin Web Server                        │   │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐              │   │
│  │  │ API路由   │ │ WebSocket │ │ 中间件    │              │   │
│  │  └───────────┘ └───────────┘ └───────────┘              │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    游戏引擎 (Game Engine)                 │   │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐              │   │
│  │  │ World     │ │ Player    │ │ NPC       │              │   │
│  │  │ (世界)    │ │ (玩家)    │ │ (NPC)     │              │   │
│  │  └───────────┘ └───────────┘ └───────────┘              │   │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐              │   │
│  │  │ Farm      │ │ Quest     │ │ Time      │              │   │
│  │  │ (农场)    │ │ (任务)    │ │ (时间)    │              │   │
│  │  └───────────┘ └───────────┘ └───────────┘              │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 前端架构

### 目录结构
```
frontend/
├── src/
│   ├── components/          # Vue组件
│   │   ├── GameCanvas.vue   # 游戏画布容器
│   │   ├── Inventory.vue    # 背包界面
│   │   ├── Toolbar.vue      # 工具栏
│   │   ├── QuestPanel.vue   # 任务面板
│   │   ├── StatusPanel.vue  # 状态面板
│   │   ├── NPCDialog.vue    # NPC对话
│   │   └── TimeDisplay.vue  # 时间显示
│   │
│   ├── game/                # 游戏引擎
│   │   ├── engine/          # 核心引擎
│   │   │   ├── Renderer.ts      # 渲染器
│   │   │   ├── GameLoop.ts      # 游戏循环
│   │   │   └── InputHandler.ts  # 输入处理
│   │   │
│   │   ├── entities/        # 游戏实体
│   │   │   ├── Player.ts        # 玩家实体
│   │   │   ├── NPC.ts           # NPC实体
│   │   │   └── Crop.ts          # 作物实体
│   │   │
│   │   └── systems/         # 游戏系统 (预留)
│   │
│   ├── stores/              # Pinia状态管理
│   │   └── gameStore.ts     # 游戏状态存储
│   │
│   ├── api/                 # API通信
│   │   └── websocket.ts     # WebSocket客户端
│   │
│   ├── types/               # TypeScript类型
│   │   └── game.ts          # 游戏类型定义
│   │
│   ├── styles/              # 样式文件
│   │   └── main.css         # 主样式
│   │
│   ├── App.vue              # 根组件
│   └── main.ts              # 入口文件
│
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

### 核心模块

#### 1. 渲染器 (Renderer)
```typescript
class Renderer {
  // 职责:
  // - Canvas 2D 渲染
  // - 相机跟随
  // - 地图绘制 (草地、建筑、树木、水)
  // - 角色绘制 (玩家、NPC)
  // - 作物绘制
  // - UI覆盖层 (时间、小地图)

  render(state: GameState): void
  resize(width, height): void

  // 私有方法
  - drawGrassTile()
  - drawBuilding()
  - drawTree()
  - drawWater()
  - drawPlayer()
  - drawNPC()
  - drawCrop()
  - drawMinimap()
}
```

#### 2. 游戏循环 (GameLoop)
```typescript
class GameLoop {
  // 职责:
  // - 60fps 渲染循环
  // - 帧时间管理
  // - 状态获取与渲染调度

  start(): void
  stop(): void
  resize(width, height): void
  setOnTick(callback): void
}
```

#### 3. 输入处理 (InputHandler)
```typescript
class InputHandler {
  // 职责:
  // - 键盘事件监听
  // - 移动处理 (WASD/方向键)
  // - 工具选择 (1-5)
  // - 动作处理 (空格=使用工具, E=交互, Q=收获)

  attach(): void
  detach(): void

  // 私有方法
  - handleKeyDown()
  - handleKeyUp()
  - handleAction()
  - handleInteract()
}
```

#### 4. 状态管理 (gameStore)
```typescript
// Pinia Store
const useGameStore = defineStore('game', () => {
  // 状态
  gameState: GameState | null
  isConnected: boolean
  isLoading: boolean
  error: string | null
  messages: string[]

  // 计算属性
  player, time, farm, npcs, quests
  activeQuests, availableQuests
  timeString, dateString, energyPercent

  // 动作
  fetchState(): Promise<void>
  executeAction(action): Promise<boolean>
  move(direction): Promise<boolean>
  useTool(tool): Promise<boolean>
  plant(seed): Promise<boolean>
  harvest(): Promise<boolean>
  buy/sell/talk/giveGift/acceptQuest/wait/sleep/resetGame
})
```

---

## 后端架构

### 目录结构
```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口
│
├── internal/
│   ├── game/                    # 游戏引擎
│   │   ├── engine.go            # 游戏引擎核心
│   │   ├── world.go             # 世界管理
│   │   ├── player.go            # 玩家逻辑
│   │   ├── npc.go               # NPC逻辑
│   │   ├── farm.go              # 农场系统
│   │   ├── quest.go             # 任务系统
│   │   ├── time.go              # 时间系统
│   │   └── friendship.go        # 友好度系统
│   │
│   ├── api/                     # API层
│   │   ├── handler.go           # 请求处理器
│   │   ├── routes.go            # 路由配置
│   │   └── middleware.go        # 中间件
│   │
│   ├── models/                  # 数据模型
│   │   ├── state.go             # 状态模型
│   │   ├── action.go            # 动作模型
│   │   └── entity.go            # 实体模型
│   │
│   └── websocket/               # WebSocket
│       └── hub.go               # 连接管理
│
├── go.mod
└── go.sum
```

### 核心模块

#### 1. 游戏引擎 (Engine)
```go
type Engine struct {
    state         *GameState
    world         *World
    timeSystem    *TimeSystem
    farmSystem    *FarmSystem
    questSystem   *QuestSystem
    friendshipSystem *FriendshipSystem
    ticker        *time.Ticker
    running       bool
    eventChan     chan GameEvent
    actionChan    chan *Action
}

// 核心方法
func NewEngine() *Engine
func (e *Engine) Start()
func (e *Engine) Stop()
func (e *Engine) tick()
func (e *Engine) GetState() *GameState
func (e *Engine) Reset() *GameState
func (e *Engine) GetObservation() *Observation
```

#### 2. 世界管理 (World)
```go
type World struct {
    config    *GameConfig
    tiles     [][]Tile
    buildings []Building
    areas     []MapArea
    trees     []Position
}

// 核心方法
func (w *World) Initialize(config *GameConfig)
func (w *World) GetTile(x, y int) *Tile
func (w *World) IsWalkable(x, y int) bool
func (w *World) GetBuildingAt(x, y int) *Building
```

#### 3. 玩家管理 (PlayerManager)
```go
type PlayerManager struct {
    world *World
}

// 核心方法
func (pm *PlayerManager) Move(player *PlayerState, direction Direction) error
func (pm *PlayerManager) UseTool(player *PlayerState, tool string, farm *FarmState) (*FarmActionResult, error)
func (pm *PlayerManager) Plant(player *PlayerState, seedType string, farm *FarmState) (*FarmActionResult, error)
func (pm *PlayerManager) Harvest(player *PlayerState, farm *FarmState, pos Position) (*FarmActionResult, error)
```

#### 4. 农场系统 (FarmSystem)
```go
type FarmSystem struct {
    config   *GameConfig
    cropInfo map[CropType]CropInfo
}

// 核心方法
func (fs *FarmSystem) CreateFarm(width, height int) FarmState
func (fs *FarmSystem) NewDayReset(farm *FarmState)
func (fs *FarmSystem) Update(time TimeState, farm FarmState)
func (fs *FarmSystem) CanPlant(cropType CropType, season Season) bool
```

---

## 游戏系统

### 1. 地图系统
```
地图尺寸: 48 x 48 格子
农场区域: 32 x 32 格子 (居中)
格子大小: 32 x 32 像素

地图元素:
- 草地 (Grass) - 基础地形
- 泥土 (Dirt) - 农场区域
- 耕地 (Tilled) - 锄过的地
- 浇水 (Watered) - 浇过水的地
- 水 (Water) - 池塘/河流
- 建筑 (Building) - 房屋/商店

建筑布局:
- Community Center: (18, 4) 12x8
- Pierre's Shop: (4, 16) 8x6
- Robin's Carpenter: (6, 4) 8x6
- Willy's Fish Shop: (2, 36) 8x6
- Saloon: (36, 20) 8x6
- Houses: 多个住宅
```

### 2. 角色系统

#### 玩家属性
```
位置 (Position): x, y 坐标
朝向 (Direction): up/down/left/right
体力 (Energy): 0-100, 每次动作消耗
金币 (Gold): 货币
背包 (Inventory): 物品列表
工具 (Tools): 锄头、水壶、斧头、镐子、镰刀
```

#### NPC列表
| ID | 名字 | 角色 | 初始位置 |
|----|------|------|----------|
| lewis | Lewis | 村长 | (40, 10) |
| pierre | Pierre | 商店老板 | (10, 30) |
| robin | Robin | 木匠 | (15, 35) |
| haley | Haley | 邻居 | (25, 15) |
| willy | Willy | 渔夫 | (5, 45) |

### 3. 农耕系统

#### 作物类型
| 作物 | 季节 | 生长天数 | 基础售价 | 种子价格 |
|------|------|----------|----------|----------|
| 土豆 | 春季 | 6天 | 80g | 50g |
| 番茄 | 夏季 | 11天 | 60g | 50g |
| 玉米 | 夏/秋 | 14天 | 50g | 75g |
| 南瓜 | 秋季 | 13天 | 320g | 100g |
| 草莓 | 春季 | 8天 | 120g | 100g |
| 花椰菜 | 春季 | 12天 | 175g | 80g |

#### 生长阶段
```
种子 (Seed) -> 发芽 (Sprout) -> 成长 (Growing) -> 成熟 (Mature)
```

### 4. 时间系统

```
时间比例: 1秒现实时间 = 10分钟游戏时间
一天时长: 6:00 AM - 2:00 AM (20小时)
季节天数: 28天
四季循环: 春 -> 夏 -> 秋 -> 冬

昼夜循环:
- 早晨 (6:00-12:00): 天空亮蓝色
- 下午 (12:00-18:00): 天空蓝色
- 傍晚 (18:00-21:00): 天空橙红色
- 夜晚 (21:00-6:00): 天空深蓝色
```

### 5. 任务系统

#### 任务列表
| ID | 名称 | 目标 | 奖励 |
|----|------|------|------|
| q001 | First Harvest | 收获1个作物 | 100g + 土豆种子x10 |
| q002 | Potato Farmer | 收获5个土豆 | 250g + 番茄种子x15 |
| q003 | Making Friends | 与Pierre对话 | 50g |
| q004 | Robin's Request | 收集10木头 | 150g + 箱子x1 |
| q005 | A Gift for Haley | 送Haley礼物 | 100g + 花x3 |
| q006 | Corn for Community | 收获10个玉米 | 500g + 南瓜种子x20 |

### 6. 友好度系统

```
最大友好度: 2500 (10颗心)
每颗心: 250点

礼物好感度:
- Love: +80点
- Like: +45点
- Neutral: +20点
- Dislike: -20点
- Hate: -40点

每日衰减: -2点 (无互动时)
```

---

## API设计

### RESTful API

#### 游戏状态
```
GET /api/v1/game/state
响应: 完整游戏状态 (GameState)

GET /api/v1/game/observation
响应: Agent可观察状态 (Observation)

GET /api/v1/game/map
响应: 地图数据

POST /api/v1/game/reset
响应: 重置后的游戏状态
```

#### 动作执行
```
POST /api/v1/game/action
请求体: Action
响应: ActionResult
```

### WebSocket
```
连接: GET /api/v1/ws
协议: 文本帧 (JSON)

消息类型:
- state: 完整状态推送
- action_result: 动作执行结果
```

### 动作类型

```json
// 移动
{
  "type": "move",
  "params": { "direction": "up|down|left|right" }
}

// 使用工具
{
  "type": "use_tool",
  "params": { "tool": "hoe|watering_can|axe|pickaxe|scythe" }
}

// 种植
{
  "type": "plant",
  "params": { "seed": "potato|tomato|corn|pumpkin|strawberry|cauliflower" }
}

// 收获
{
  "type": "harvest",
  "params": {}
}

// 购买
{
  "type": "buy",
  "params": { "item": "seed_potato", "quantity": 5 }
}

// 出售
{
  "type": "sell",
  "params": { "item": "potato", "quantity": 3 }
}

// 对话
{
  "type": "talk",
  "params": { "npc": "lewis|pierre|robin|haley|willy" }
}

// 送礼
{
  "type": "give_gift",
  "params": { "npc": "haley", "gift_item": "flower" }
}

// 接受任务
{
  "type": "accept_quest",
  "params": { "quest_id": "q001" }
}

// 等待
{
  "type": "wait",
  "params": { "ticks": 10 }
}

// 睡觉
{
  "type": "sleep",
  "params": {}
}
```

---

## 数据模型

### GameState
```typescript
interface GameState {
  player: PlayerState      // 玩家状态
  time: TimeState          // 时间状态
  farm: FarmState          // 农场状态
  npcs: NPCState[]         // NPC列表
  quests: QuestData[]      // 任务列表
  map_width: number        // 地图宽度
  map_height: number       // 地图高度
  game_over: boolean       // 游戏结束
  day_ended: boolean       // 一天结束
}
```

### PlayerState
```typescript
interface PlayerState {
  position: Position       // 位置
  direction: Direction     // 朝向
  energy: number           // 当前体力
  max_energy: number       // 最大体力
  gold: number             // 金币
  inventory: InventoryItem[] // 背包
  tools: string[]          // 工具列表
  active_tool: string      // 当前工具
}
```

### TimeState
```typescript
interface TimeState {
  day: number              // 天数 (1-28)
  hour: number             // 小时 (6-26)
  minute: number           // 分钟 (0-59)
  season: Season           // 季节
  year: number             // 年份
  paused: boolean          // 是否暂停
}
```

### FarmState
```typescript
interface FarmState {
  width: number            // 农场宽度
  height: number           // 农场高度
  tiles: Tile[][]          // 格子数据
}

interface Tile {
  x: number
  y: number
  type: TileType           // 格子类型
  crop?: CropData          // 作物数据
  occupied: boolean        // 是否被占用
}
```

---

## 部署架构

### Docker Compose
```yaml
services:
  frontend:
    build: ./frontend
    ports: ["3000:80"]
    depends_on: [backend]

  backend:
    build: ./backend
    ports: ["8080:8080"]
    environment:
      - PORT=8080

  nginx:
    image: nginx:alpine
    ports: ["80:80"]
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on: [frontend, backend]
```

### Nginx配置
```
- 前端静态文件: /
- API代理: /api/ -> backend:8080
- WebSocket代理: /api/v1/ws -> backend:8080
```

---

## 当前问题与改进方向

### 已知问题
1. **地图渲染简单** - 需要更精美的像素贴图
2. **移动问题** - 输入响应需要优化
3. **无碰撞检测** - 玩家可以穿过建筑

### 改进计划

#### 地图改进
- [ ] 使用精灵图(Sprite Sheet)替代程序生成
- [ ] 添加地图编辑器支持
- [ ] 实现多层渲染 (地面层、物体层、角色层)
- [ ] 添加天气效果

#### 游戏机制
- [ ] 完善碰撞检测
- [ ] 添加动物系统
- [ ] 实现钓鱼小游戏
- [ ] 添加挖矿系统

#### AI Agent
- [ ] 添加更多观察数据
- [ ] 实现奖励系统
- [ ] 添加环境重置API
- [ ] 支持多Agent并发

---

*文档版本: 1.0*
*最后更新: 2024*
