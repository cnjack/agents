# API 接口文档

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **内容类型**: `application/json`
- **WebSocket**: `ws://localhost:8080/api/v1/ws`

---

## 游戏状态接口

### 获取完整游戏状态

```
GET /game/state
```

**响应示例**:
```json
{
  "player": {
    "position": { "x": 24, "y": 24 },
    "direction": "down",
    "energy": 100,
    "max_energy": 100,
    "gold": 500,
    "inventory": [
      { "id": "seed_potato", "name": "Potato Seeds", "type": "seed", "quantity": 15, "price": 50 }
    ],
    "tools": ["hoe", "watering_can", "axe", "pickaxe", "scythe"],
    "active_tool": "hoe"
  },
  "time": {
    "day": 1,
    "hour": 6,
    "minute": 0,
    "season": "spring",
    "year": 1,
    "paused": false
  },
  "farm": {
    "width": 32,
    "height": 32,
    "tiles": [[...]]
  },
  "npcs": [...],
  "quests": [...],
  "map_width": 48,
  "map_height": 48,
  "game_over": false,
  "day_ended": false
}
```

### 获取Agent观察状态

```
GET /game/observation
```

**响应示例**:
```json
{
  "player": { ... },
  "time": { ... },
  "nearby_tiles": [...],
  "nearby_npcs": [...],
  "active_quests": [...],
  "messages": ["Welcome to Stardew Valley!"]
}
```

### 获取地图数据

```
GET /game/map
```

**响应示例**:
```json
{
  "width": 48,
  "height": 48,
  "farm": { ... },
  "buildings": [...],
  "areas": [...]
}
```

### 重置游戏

```
POST /game/reset
```

**响应示例**:
```json
{
  "success": true,
  "message": "Game reset",
  "state": { ... }
}
```

---

## 动作执行接口

### 执行动作

```
POST /game/action
```

**请求体**:
```json
{
  "type": "<action_type>",
  "params": { ... }
}
```

**响应**:
```json
{
  "success": true,
  "message": "Action executed",
  "new_state": { ... },
  "events": [...]
}
```

---

## 动作类型详解

### 1. 移动 (move)

**请求**:
```json
{
  "type": "move",
  "params": {
    "direction": "up"
  }
}
```

**参数**:
- `direction`: 移动方向
  - `"up"` - 向上
  - `"down"` - 向下
  - `"left"` - 向左
  - `"right"` - 向右

**响应**:
```json
{
  "success": true,
  "message": "Moved up",
  "new_state": { ... }
}
```

**错误情况**:
- 目标位置不可行走 (水、建筑边界)
- `"success": false, "message": "cannot move there"`

---

### 2. 使用工具 (use_tool)

**请求**:
```json
{
  "type": "use_tool",
  "params": {
    "tool": "hoe"
  }
}
```

**工具列表**:
| 工具 | 用途 | 体力消耗 |
|------|------|----------|
| `hoe` | 翻地 | 2 |
| `watering_can` | 浇水 | 1 |
| `axe` | 砍树 | 2 |
| `pickaxe` | 挖矿 | 2 |
| `scythe` | 收割 | 0 |

**响应示例** (锄地):
```json
{
  "success": true,
  "message": "Used hoe",
  "events": [
    {
      "type": "tool_used",
      "data": {
        "success": true,
        "tile_pos": { "x": 24, "y": 23 },
        "old_type": "dirt",
        "new_type": "tilled"
      }
    }
  ]
}
```

---

### 3. 种植 (plant)

**请求**:
```json
{
  "type": "plant",
  "params": {
    "seed": "potato"
  }
}
```

**种子类型**:
- `potato` - 土豆
- `tomato` - 番茄
- `corn` - 玉米
- `pumpkin` - 南瓜
- `strawberry` - 草莓
- `cauliflower` - 花椰菜

**前置条件**:
- 玩家面向的格子必须是耕地 (tilled) 或已浇水 (watered)
- 玩家背包中必须有对应种子

**响应**:
```json
{
  "success": true,
  "message": "Planted potato",
  "new_state": { ... }
}
```

---

### 4. 收获 (harvest)

**请求**:
```json
{
  "type": "harvest",
  "params": {}
}
```

**前置条件**:
- 玩家面向的格子有成熟作物

**响应**:
```json
{
  "success": true,
  "message": "Harvested Potato",
  "events": [
    {
      "type": "item_harvested",
      "data": {
        "id": "potato",
        "name": "Potato",
        "type": "crop",
        "quantity": 1,
        "price": 80
      }
    }
  ]
}
```

---

### 5. 购买 (buy)

**请求**:
```json
{
  "type": "buy",
  "params": {
    "item": "seed_potato",
    "quantity": 5
  }
}
```

**前置条件**:
- 玩家必须靠近商店NPC (Pierre, Robin, Willy)
- 玩家必须有足够金币

**可购买物品**:
| 物品ID | 名称 | 价格 |
|--------|------|------|
| `seed_potato` | 土豆种子 | 50g |
| `seed_tomato` | 番茄种子 | 50g |
| `seed_corn` | 玉米种子 | 75g |
| `seed_pumpkin` | 南瓜种子 | 100g |
| `seed_strawberry` | 草莓种子 | 100g |
| `seed_cauliflower` | 花椰菜种子 | 80g |

---

### 6. 出售 (sell)

**请求**:
```json
{
  "type": "sell",
  "params": {
    "item": "potato",
    "quantity": 3
  }
}
```

**前置条件**:
- 玩家必须靠近商店NPC
- 玩家背包中必须有足够数量的物品

---

### 7. 对话 (talk)

**请求**:
```json
{
  "type": "talk",
  "params": {
    "npc": "lewis"
  }
}
```

**NPC列表**:
- `lewis` - 村长 Lewis
- `pierre` - 商店老板 Pierre
- `robin` - 木匠 Robin
- `haley` - 邻居 Haley
- `willy` - 渔夫 Willy

**响应**:
```json
{
  "success": true,
  "message": "Welcome to Stardew Valley! I'm Lewis, the mayor.",
  "events": [
    {
      "type": "npc_dialog",
      "data": {
        "npc_id": "lewis",
        "npc_name": "Lewis",
        "dialogue": "Welcome to Stardew Valley!"
      }
    }
  ]
}
```

---

### 8. 送礼 (give_gift)

**请求**:
```json
{
  "type": "give_gift",
  "params": {
    "npc": "haley",
    "gift_item": "flower"
  }
}
```

**响应**:
```json
{
  "success": true,
  "message": "Oh, for me? How kind of you!",
  "events": [...]
}
```

---

### 9. 接受任务 (accept_quest)

**请求**:
```json
{
  "type": "accept_quest",
  "params": {
    "quest_id": "q001"
  }
}
```

**响应**:
```json
{
  "success": true,
  "message": "Accepted quest: First Harvest"
}
```

---

### 10. 等待 (wait)

**请求**:
```json
{
  "type": "wait",
  "params": {
    "ticks": 10
  }
}
```

**说明**: 每tick = 1游戏分钟

---

### 11. 睡觉 (sleep)

**请求**:
```json
{
  "type": "sleep",
  "params": {}
}
```

**效果**:
- 结束当前一天
- 恢复满体力
- 进入下一天
- 作物生长 (已浇水的)

---

## WebSocket 协议

### 连接
```
ws://localhost:8080/api/v1/ws
```

### 消息格式

#### 服务端 -> 客户端

**状态推送**:
```json
{
  "type": "state",
  "data": { ... GameState }
}
```

**动作结果**:
```json
{
  "type": "action_result",
  "action": "move",
  "state": { ... GameState }
}
```

#### 客户端 -> 服务端

**执行动作**:
```json
{
  "type": "move",
  "params": { "direction": "up" }
}
```

---

## 错误处理

### 错误响应格式
```json
{
  "success": false,
  "message": "Error description"
}
```

### 常见错误

| 错误信息 | 原因 |
|----------|------|
| `cannot move there` | 目标位置不可行走 |
| `not enough energy` | 体力不足 |
| `not enough gold` | 金币不足 |
| `no seeds of this type` | 没有该类型种子 |
| `can only plant on tilled soil` | 只能在耕地上种植 |
| `something is already planted here` | 已经有作物 |
| `crop not ready for harvest` | 作物未成熟 |
| `no shopkeeper nearby` | 附近没有商店NPC |
| `NPC not found` | NPC不存在 |
| `quest not found` | 任务不存在 |
| `quest not available` | 任务不可用 |

---

## 使用示例

### cURL 示例

```bash
# 获取游戏状态
curl http://localhost:8080/api/v1/game/state

# 移动
curl -X POST http://localhost:8080/api/v1/game/action \
  -H "Content-Type: application/json" \
  -d '{"type":"move","params":{"direction":"down"}}'

# 使用工具
curl -X POST http://localhost:8080/api/v1/game/action \
  -H "Content-Type: application/json" \
  -d '{"type":"use_tool","params":{"tool":"hoe"}}'

# 种植
curl -X POST http://localhost:8080/api/v1/game/action \
  -H "Content-Type: application/json" \
  -d '{"type":"plant","params":{"seed":"potato"}}'

# 重置游戏
curl -X POST http://localhost:8080/api/v1/game/reset
```

### Python 示例

```python
import requests

BASE_URL = "http://localhost:8080/api/v1"

def get_state():
    return requests.get(f"{BASE_URL}/game/state").json()

def move(direction):
    return requests.post(f"{BASE_URL}/game/action", json={
        "type": "move",
        "params": {"direction": direction}
    }).json()

def use_tool(tool):
    return requests.post(f"{BASE_URL}/game/action", json={
        "type": "use_tool",
        "params": {"tool": tool}
    }).json()

# 使用示例
state = get_state()
print(f"Position: {state['player']['position']}")
print(f"Gold: {state['player']['gold']}")

# 移动并耕地
move("down")
use_tool("hoe")
```

---

*API版本: v1*
*最后更新: 2024*
