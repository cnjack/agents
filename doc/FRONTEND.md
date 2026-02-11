# 前端开发文档

## 快速开始

### 安装依赖
```bash
cd frontend
npm install
```

### 开发模式
```bash
npm run dev
# 访问 http://localhost:5173
```

### 构建生产版本
```bash
npm run build
```

---

## 项目结构

```
frontend/
├── src/
│   ├── components/           # Vue组件
│   ├── game/                 # 游戏引擎
│   ├── stores/               # 状态管理
│   ├── api/                  # API通信
│   ├── types/                # TypeScript类型
│   └── styles/               # 样式文件
├── public/                   # 静态资源
├── index.html               # HTML入口
├── vite.config.ts           # Vite配置
├── tailwind.config.js       # Tailwind配置
└── tsconfig.json            # TypeScript配置
```

---

## 核心模块

### 1. 渲染器 (Renderer)

位置: `src/game/engine/Renderer.ts`

渲染器负责所有游戏画面的绘制，使用Canvas 2D API。

#### 初始化
```typescript
import { Renderer } from '@/game/engine/Renderer'

const canvas = document.getElementById('game-canvas') as HTMLCanvasElement
const renderer = new Renderer(canvas)
```

#### 渲染循环
```typescript
function gameLoop() {
  renderer.render(gameState)
  requestAnimationFrame(gameLoop)
}
```

#### 绘制元素

**地图元素**:
- `drawGrassTile()` - 草地格子
- `drawPathTile()` - 道路格子
- `drawWaterTile()` - 水面 (带动画)
- `drawBuilding()` - 建筑
- `drawTree()` - 树木
- `drawFence()` - 栅栏

**角色元素**:
- `drawPlayer()` - 玩家
- `drawNPC()` - NPC

**游戏元素**:
- `drawCrop()` - 作物 (4个生长阶段)
- `drawMinimap()` - 小地图

#### 相机系统
```typescript
// 相机跟随玩家
this.targetCameraX = state.player.position.x * TILE_SIZE - canvas.width / 2
this.targetCameraY = state.player.position.y * TILE_SIZE - canvas.height / 2

// 平滑移动
this.cameraX += (this.targetCameraX - this.cameraX) * 0.15
```

#### 颜色主题
```typescript
const COLORS = {
  // 地面
  grass: '#5a8f3a',
  dirt: '#a08060',
  path: '#c4a87a',
  water: '#4a90c2',

  // 建筑
  buildingWall: '#b87333',
  buildingRoof: '#8b4513',

  // 角色
  player: '#ff6b6b',
  npc: '#9370db',
}
```

---

### 2. 游戏循环 (GameLoop)

位置: `src/game/engine/GameLoop.ts`

管理60fps渲染循环。

```typescript
import { GameLoop } from '@/game/engine/GameLoop'

const gameLoop = new GameLoop(canvas, () => store.gameState)
gameLoop.start()

// 停止
gameLoop.stop()

// 调整大小
gameLoop.resize(width, height)

// 设置每帧回调
gameLoop.setOnTick((deltaTime) => {
  // 自定义逻辑
})
```

---

### 3. 输入处理 (InputHandler)

位置: `src/game/engine/InputHandler.ts`

处理键盘输入。

#### 键位绑定
| 按键 | 功能 |
|------|------|
| W / ↑ | 向上移动 |
| S / ↓ | 向下移动 |
| A / ← | 向左移动 |
| D / → | 向右移动 |
| Space | 使用当前工具 |
| E | 与NPC交互 |
| Q | 收获作物 |
| 1-5 | 选择工具 |

#### 使用方式
```typescript
import { InputHandler } from '@/game/engine/InputHandler'

const inputHandler = new InputHandler(store)
inputHandler.attach()

// 清理
inputHandler.detach()
```

---

### 4. 状态管理 (gameStore)

位置: `src/stores/gameStore.ts`

使用Pinia进行全局状态管理。

#### 状态
```typescript
const store = useGameStore()

store.gameState     // 完整游戏状态
store.isConnected   // WebSocket连接状态
store.isLoading     // 加载状态
store.error         // 错误信息
store.messages      // 消息列表
```

#### 计算属性
```typescript
store.player        // 玩家状态
store.time          // 时间状态
store.farm          // 农场状态
store.npcs          // NPC列表
store.quests        // 任务列表

store.activeQuests      // 进行中的任务
store.availableQuests   // 可接受的任务
store.completedQuests   // 已完成的任务

store.timeString    // "8:30 AM"
store.dateString    // "Spring 1, Year 1"
store.energyPercent // 0-100
```

#### 动作
```typescript
// 获取状态
await store.fetchState()

// 移动
await store.move('down')

// 使用工具
await store.useTool('hoe')

// 种植
await store.plant('potato')

// 收获
await store.harvest()

// 购买/出售
await store.buy('seed_potato', 5)
await store.sell('potato', 3)

// NPC交互
await store.talk('lewis')
await store.giveGift('haley', 'flower')

// 任务
await store.acceptQuest('q001')

// 时间
await store.wait(10)
await store.sleep()

// 重置
await store.resetGame()
```

---

### 5. WebSocket通信

位置: `src/api/websocket.ts`

#### 连接
```typescript
import { connectWebSocket, disconnectWebSocket } from '@/api/websocket'

// 连接
connectWebSocket()

// 断开
disconnectWebSocket()
```

#### 发送动作
```typescript
import { sendAction } from '@/api/websocket'

sendAction({
  type: 'move',
  params: { direction: 'down' }
})
```

---

## Vue组件

### GameCanvas.vue
游戏画布容器组件。

```vue
<template>
  <canvas ref="canvasRef" class="game-canvas" />
</template>

<script setup lang="ts">
// 自动初始化渲染器和游戏循环
</script>
```

### StatusPanel.vue
状态面板，显示体力和金币。

```vue
<StatusPanel />
```

### Inventory.vue
背包界面，4x4格子。

```vue
<Inventory />
```

### Toolbar.vue
工具栏，显示5个工具槽。

```vue
<Toolbar />
```

### QuestPanel.vue
任务面板，显示进行中和可接受的任务。

```vue
<QuestPanel />
```

### TimeDisplay.vue
时间显示，包含时钟、日期和睡眠按钮。

```vue
<TimeDisplay />
```

### NPCDialog.vue
NPC对话框。

```vue
<NPCDialog :npc-id="currentNpc" :visible="showDialog" @close="showDialog = false" />
```

---

## 样式系统

### TailwindCSS配置
```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        'pixel-green': '#4a7c23',
        'pixel-brown': '#6b4423',
        'pixel-blue': '#4a90c2',
        'pixel-gold': '#d4a017',
      },
      fontFamily: {
        pixel: ['"Press Start 2P"', 'monospace'],
      },
    },
  },
}
```

### 主要样式类

```css
/* 游戏容器 */
.game-container { display: flex; }

/* 游戏视口 */
.game-viewport { flex: 1; }

/* UI侧边栏 */
.ui-sidebar { width: 320px; }

/* UI面板 */
.ui-panel { background: rgba(74, 55, 40, 0.9); }

/* 像素按钮 */
.pixel-btn { /* 像素风格按钮 */ }

/* 背包格子 */
.inventory-slot { width: 48px; height: 48px; }

/* 工具栏 */
.toolbar { /* 底部工具栏 */ }

/* 对话框 */
.dialog-box { /* NPC对话框 */ }
```

---

## 类型定义

位置: `src/types/game.ts`

### 主要类型

```typescript
interface Position { x: number; y: number }
type Direction = 'up' | 'down' | 'left' | 'right'
type Season = 'spring' | 'summer' | 'fall' | 'winter'
type TileType = 'grass' | 'dirt' | 'tilled' | 'watered' | 'water' | 'building'
type CropStage = 'seed' | 'sprout' | 'growing' | 'mature'
type CropType = 'potato' | 'tomato' | 'corn' | 'pumpkin' | 'strawberry' | 'cauliflower'

interface GameState { /* 见数据模型文档 */ }
interface PlayerState { /* ... */ }
interface TimeState { /* ... */ }
interface FarmState { /* ... */ }
interface NPCState { /* ... */ }
interface QuestData { /* ... */ }
interface Action { /* ... */ }
interface ActionResult { /* ... */ }
```

---

## 游戏地图

### 地图特征 (MAP_FEATURES)

```typescript
const MAP_FEATURES: MapFeature[] = [
  // 建筑
  { type: 'building', x: 18, y: 4, width: 12, height: 8, variant: 'community_center' },
  { type: 'building', x: 4, y: 16, width: 8, height: 6, variant: 'shop' },
  { type: 'building', x: 6, y: 4, width: 8, height: 6, variant: 'carpenter' },
  { type: 'building', x: 2, y: 36, width: 8, height: 6, variant: 'fish_shop' },
  { type: 'building', x: 36, y: 20, width: 8, height: 6, variant: 'saloon' },

  // 水域
  { type: 'water', x: 38, y: 2, width: 8, height: 6 },

  // 树木
  { type: 'tree', x: 1, y: 1 },
  { type: 'tree', x: 44, y: 1 },
  // ...

  // 道路
  { type: 'path_h', x: 14, y: 24, width: 20 },
  { type: 'path_v', x: 24, y: 12, height: 12 },

  // 栅栏
  { type: 'fence', x: 12, y: 28 },
]
```

### 地图尺寸
- 总大小: 48 x 48 格子
- 格子像素: 32 x 32
- 农场区域: 32 x 32 (居中)

---

## 开发建议

### 添加新作物
1. 在 `Crop.ts` 的 `CROP_DATA` 中添加配置
2. 在渲染器中添加对应的绘制逻辑 (如需要)

### 添加新NPC
1. 在后端 `engine.go` 的 `createNPCs()` 中添加
2. 在前端渲染器的 `npcColors` 中添加颜色

### 添加新建筑
1. 在后端 `world.go` 的 `setupBuildings()` 中添加
2. 在前端 `MAP_FEATURES` 中添加对应的渲染配置

### 自定义渲染
修改 `Renderer.ts` 中的绘制函数:
- `drawGrassTile()` - 草地样式
- `drawBuilding()` - 建筑外观
- `drawPlayer()` - 玩家角色
- 等等

---

## 性能优化

1. **视口裁剪**: 只渲染可见区域的格子
2. **对象池**: 复用Canvas绘制状态
3. **请求动画帧**: 使用 `requestAnimationFrame`
4. **状态缓存**: 避免每帧重新计算不变的数据

---

*文档版本: 1.0*
