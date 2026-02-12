# 游戏交互与移动系统优化方案

## 问题分析

### 1. 当前对话交互流程
```
前端 InputHandler.handleInteract()
    ↓
gameStore.talk(npcId)
    ↓
executeAction({type: 'talk', params: {npc}})
    ↓
POST /api/v1/game/action
    ↓
后端 handler.handleTalk()
    ↓
返回 dialogue + state
    ↓
前端 updateState() 直接替换状态
```

### 2. 当前移动问题
- NPC位置根据Schedule直接切换（闪现）
- 玩家移动等待服务器响应后才更新
- 网络延迟导致移动卡顿
- 无插值动画

### 3. 网络波动问题
- WebSocket心跳30秒间隔过长
- 状态更新无缓冲
- 无客户端预测
- 无延迟补偿

---

## 优化方案

### 一、寻路与平滑移动系统

#### 1.1 NPC移动架构

```typescript
// frontend/src/game/entities/MovingEntity.ts

interface MovementState {
  // 位置状态
  currentPos: { x: number; y: number }      // 当前像素位置
  targetPos: { x: number; y: number }        // 目标像素位置
  tilePos: { x: number; y: number }          // 当前格子位置

  // 移动状态
  isMoving: boolean
  direction: Direction
  speed: number                              // 像素/秒

  // 寻路
  path: Array<{ x: number; y: number }>      // 路径点列表
  currentPathIndex: number

  // 动画
  animFrame: number
  animTimer: number
}

class MovingEntity {
  protected movement: MovementState

  // 平滑移动更新
  update(deltaTime: number) {
    if (!this.movement.isMoving) return

    const dx = this.movement.targetPos.x - this.movement.currentPos.x
    const dy = this.movement.targetPos.y - this.movement.currentPos.y
    const distance = Math.sqrt(dx * dx + dy * dy)

    if (distance < 2) {
      // 到达目标点
      this.movement.currentPos = { ...this.movement.targetPos }
      this.moveToNextPathPoint()
    } else {
      // 继续移动
      const moveDistance = this.movement.speed * deltaTime
      const ratio = Math.min(moveDistance / distance, 1)

      this.movement.currentPos.x += dx * ratio
      this.movement.currentPos.y += dy * ratio
    }

    // 更新动画帧
    this.updateAnimation(deltaTime)
  }

  // 移动到下一个路径点
  private moveToNextPathPoint() {
    this.movement.currentPathIndex++

    if (this.movement.currentPathIndex >= this.movement.path.length) {
      // 路径完成
      this.movement.isMoving = false
      this.movement.path = []
      return
    }

    // 设置下一个目标
    const nextTile = this.movement.path[this.movement.currentPathIndex]
    this.movement.targetPos = {
      x: nextTile.x * TILE_SIZE + TILE_SIZE / 2,
      y: nextTile.y * TILE_SIZE + TILE_SIZE / 2
    }
    this.updateDirection()
  }

  // 设置新路径
  setPath(path: Array<{ x: number; y: number }>) {
    if (path.length === 0) return

    this.movement.path = path
    this.movement.currentPathIndex = 0
    this.movement.isMoving = true

    // 立即设置第一个目标
    const firstTile = path[0]
    this.movement.targetPos = {
      x: firstTile.x * TILE_SIZE + TILE_SIZE / 2,
      y: firstTile.y * TILE_SIZE + TILE_SIZE / 2
    }
    this.updateDirection()
  }
}
```

#### 1.2 A*寻路算法

```typescript
// frontend/src/game/pathfinding/AStar.ts

interface PathNode {
  x: number
  y: number
  g: number  // 从起点到当前节点的成本
  h: number  // 从当前节点到终点的估算成本
  f: number  // g + h
  parent: PathNode | null
}

export class Pathfinder {
  private width: number
  private height: number
  private walkableMap: boolean[][]

  constructor(width: number, height: number) {
    this.width = width
    this.height = height
    this.walkableMap = []
  }

  // 更新可行走地图
  updateWalkableMap(tiles: Tile[][], npcs: NPCState[]) {
    this.walkableMap = []

    for (let y = 0; y < this.height; y++) {
      this.walkableMap[y] = []
      for (let x = 0; x < this.width; x++) {
        const tile = tiles[y]?.[x]
        this.walkableMap[y][x] = this.isTileWalkable(tile)
      }
    }

    // 标记NPC占用位置
    for (const npc of npcs) {
      if (npc.position.x >= 0 && npc.position.x < this.width &&
          npc.position.y >= 0 && npc.position.y < this.height) {
        this.walkableMap[npc.position.y][npc.position.x] = false
      }
    }
  }

  // A*寻路
  findPath(
    startX: number, startY: number,
    endX: number, endY: number
  ): Array<{ x: number; y: number }> {

    if (!this.isWalkable(endX, endY)) {
      return []
    }

    const openList: PathNode[] = []
    const closedSet = new Set<string>()

    const startNode: PathNode = {
      x: startX, y: startY,
      g: 0, h: this.heuristic(startX, startY, endX, endY),
      f: 0, parent: null
    }
    startNode.f = startNode.g + startNode.h

    openList.push(startNode)

    while (openList.length > 0) {
      // 找f值最小的节点
      openList.sort((a, b) => a.f - b.f)
      const current = openList.shift()!

      // 到达终点
      if (current.x === endX && current.y === endY) {
        return this.reconstructPath(current)
      }

      closedSet.add(`${current.x},${current.y}`)

      // 检查相邻节点
      const neighbors = [
        { x: current.x - 1, y: current.y },
        { x: current.x + 1, y: current.y },
        { x: current.x, y: current.y - 1 },
        { x: current.x, y: current.y + 1 }
      ]

      for (const neighbor of neighbors) {
        const key = `${neighbor.x},${neighbor.y}`

        if (closedSet.has(key)) continue
        if (!this.isWalkable(neighbor.x, neighbor.y)) continue

        const g = current.g + 1
        const h = this.heuristic(neighbor.x, neighbor.y, endX, endY)

        const existingNode = openList.find(n => n.x === neighbor.x && n.y === neighbor.y)

        if (!existingNode) {
          openList.push({
            x: neighbor.x, y: neighbor.y,
            g, h, f: g + h,
            parent: current
          })
        } else if (g < existingNode.g) {
          existingNode.g = g
          existingNode.f = g + existingNode.h
          existingNode.parent = current
        }
      }
    }

    return []  // 无路径
  }

  private heuristic(x1: number, y1: number, x2: number, y2: number): number {
    // 曼哈顿距离
    return Math.abs(x1 - x2) + Math.abs(y1 - y2)
  }

  private isWalkable(x: number, y: number): boolean {
    if (x < 0 || x >= this.width || y < 0 || y >= this.height) {
      return false
    }
    return this.walkableMap[y]?.[x] ?? false
  }

  private reconstructPath(node: PathNode): Array<{ x: number; y: number }> {
    const path: Array<{ x: number; y: number }> = []
    let current: PathNode | null = node

    while (current) {
      path.unshift({ x: current.x, y: current.y })
      current = current.parent
    }

    return path.slice(1)  // 移除起点
  }
}
```

#### 1.3 NPC移动管理器

```typescript
// frontend/src/game/entities/NPCMovementManager.ts

export class NPCMovementManager {
  private npcs: Map<string, MovingNPC> = new Map()
  private pathfinder: Pathfinder

  constructor(pathfinder: Pathfinder) {
    this.pathfinder = pathfinder
  }

  // 更新NPC位置（来自服务器）
  updateNPCPosition(npcId: string, targetTile: { x: number; y: number }) {
    const npc = this.npcs.get(npcId)
    if (!npc) return

    // 如果NPC已在目标位置附近，不移动
    const currentTile = npc.getTilePosition()
    if (currentTile.x === targetTile.x && currentTile.y === targetTile.y) {
      return
    }

    // 计算路径
    const path = this.pathfinder.findPath(
      currentTile.x, currentTile.y,
      targetTile.x, targetTile.y
    )

    if (path.length > 0) {
      npc.setPath(path)
    }
  }

  // 批量更新（处理网络延迟）
  batchUpdate(updates: Array<{ id: string; pos: { x: number; y: number } }>) {
    for (const update of updates) {
      this.updateNPCPosition(update.id, update.pos)
    }
  }

  // 帧更新
  update(deltaTime: number) {
    for (const npc of this.npcs.values()) {
      npc.update(deltaTime)
    }
  }

  // 获取NPC渲染位置（平滑后）
  getRenderPosition(npcId: string): { x: number; y: number } {
    const npc = this.npcs.get(npcId)
    return npc ? npc.getRenderPosition() : { x: 0, y: 0 }
  }
}
```

---

### 二、客户端预测与延迟补偿

#### 2.1 客户端预测系统

```typescript
// frontend/src/game/network/ClientPrediction.ts

interface PredictedAction {
  id: number
  type: string
  params: any
  timestamp: number
  predictedState: Partial<GameState>
  confirmed: boolean
}

export class ClientPrediction {
  private pendingActions: PredictedAction[] = []
  private actionIdCounter = 0
  private maxPendingActions = 10

  // 预测执行动作
  predictAction(action: Action, currentState: GameState): PredictedAction {
    const predictedAction: PredictedAction = {
      id: ++this.actionIdCounter,
      type: action.type,
      params: action.params,
      timestamp: Date.now(),
      predictedState: this.simulateAction(action, currentState),
      confirmed: false
    }

    this.pendingActions.push(predictedAction)

    // 清理过期的预测
    while (this.pendingActions.length > this.maxPendingActions) {
      this.pendingActions.shift()
    }

    return predictedAction
  }

  // 服务器状态确认
  onServerState(serverState: GameState) {
    // 移除已确认的动作
    this.pendingActions = this.pendingActions.filter(action => {
      if (action.timestamp < serverState.timestamp) {
        return false
      }
      return true
    })

    // 检测预测错误，需要修正
    for (const action of this.pendingActions) {
      if (!this.compareStates(action.predictedState, serverState)) {
        // 预测错误，需要回滚
        return {
          needsCorrection: true,
          correctState: serverState,
          pendingActions: this.pendingActions
        }
      }
    }

    return { needsCorrection: false }
  }

  // 本地模拟动作
  private simulateAction(action: Action, state: GameState): Partial<GameState> {
    const predicted = JSON.parse(JSON.stringify(state))

    switch (action.type) {
      case 'move':
        const dir = action.params.direction
        const player = predicted.player

        switch (dir) {
          case 'up': player.position.y--; break
          case 'down': player.position.y++; break
          case 'left': player.position.x--; break
          case 'right': player.position.x++; break
        }
        player.direction = dir
        break

      // 其他动作模拟...
    }

    return predicted
  }

  // 状态比较
  private compareStates(predicted: any, actual: any): boolean {
    return predicted.player?.position?.x === actual.player?.position?.x &&
           predicted.player?.position?.y === actual.player?.position?.y
  }
}
```

#### 2.2 状态插值系统

```typescript
// frontend/src/game/network/StateInterpolator.ts

interface StateSnapshot {
  timestamp: number
  state: GameState
}

export class StateInterpolator {
  private snapshots: StateSnapshot[] = []
  private maxSnapshots = 10
  private interpolationDelay = 100  // 毫秒，用于平滑

  // 添加服务器状态快照
  addSnapshot(state: GameState) {
    this.snapshots.push({
      timestamp: Date.now(),
      state: state
    })

    while (this.snapshots.length > this.maxSnapshots) {
      this.snapshots.shift()
    }
  }

  // 获取插值后的状态
  getInterpolatedState(): GameState | null {
    if (this.snapshots.length < 2) {
      return this.snapshots[0]?.state || null
    }

    const renderTime = Date.now() - this.interpolationDelay

    // 找到两个快照之间的位置
    for (let i = 0; i < this.snapshots.length - 1; i++) {
      const snap1 = this.snapshots[i]
      const snap2 = this.snapshots[i + 1]

      if (renderTime >= snap1.timestamp && renderTime <= snap2.timestamp) {
        // 计算插值比例
        const t = (renderTime - snap1.timestamp) /
                  (snap2.timestamp - snap1.timestamp)

        return this.interpolate(snap1.state, snap2.state, t)
      }
    }

    // 返回最新状态
    return this.snapshots[this.snapshots.length - 1].state
  }

  // 状态插值
  private interpolate(state1: GameState, state2: GameState, t: number): GameState {
    return {
      ...state2,
      player: {
        ...state2.player,
        position: {
          x: this.lerp(state1.player.position.x, state2.player.position.x, t),
          y: this.lerp(state1.player.position.y, state2.player.position.y, t)
        }
      },
      npcs: state2.npcs.map(npc2 => {
        const npc1 = state1.npcs.find(n => n.id === npc2.id)
        if (!npc1) return npc2

        return {
          ...npc2,
          position: {
            x: this.lerp(npc1.position.x, npc2.position.x, t),
            y: this.lerp(npc1.position.y, npc2.position.y, t)
          }
        }
      })
    }
  }

  private lerp(a: number, b: number, t: number): number {
    return a + (b - a) * t
  }
}
```

---

### 三、网络优化

#### 3.1 WebSocket优化

```typescript
// frontend/src/api/OptimizedWebSocket.ts

interface WebSocketConfig {
  heartbeatInterval: number     // 心跳间隔
  reconnectDelay: number        // 重连延迟
  maxReconnectAttempts: number  // 最大重连次数
  messageQueueSize: number      // 消息队列大小
}

export class OptimizedWebSocket {
  private ws: WebSocket | null = null
  private config: WebSocketConfig = {
    heartbeatInterval: 5000,    // 5秒心跳
    reconnectDelay: 1000,
    maxReconnectAttempts: 5,
    messageQueueSize: 100
  }

  private messageQueue: Array<{ type: string; data: any }> = []
  private heartbeatTimer: number | null = null
  private lastPongTime: number = 0
  private isConnecting = false

  // 连接
  connect(url: string) {
    if (this.isConnecting || this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    this.isConnecting = true
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this.isConnecting = false
      this.lastPongTime = Date.now()
      this.startHeartbeat()
      this.flushMessageQueue()
    }

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)

      if (data.type === 'pong') {
        this.lastPongTime = Date.now()
        return
      }

      this.handleMessage(data)
    }

    this.ws.onclose = () => {
      this.stopHeartbeat()
      this.scheduleReconnect(url)
    }
  }

  // 心跳机制
  private startHeartbeat() {
    this.heartbeatTimer = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        // 检测服务器是否响应
        if (Date.now() - this.lastPongTime > this.config.heartbeatInterval * 3) {
          console.warn('Server unresponsive, reconnecting...')
          this.ws.close()
          return
        }

        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, this.config.heartbeatInterval)
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  // 消息队列
  send(type: string, data: any) {
    const message = { type, data, timestamp: Date.now() }

    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      // 连接断开时加入队列
      if (this.messageQueue.length < this.config.messageQueueSize) {
        this.messageQueue.push(message)
      }
    }
  }

  private flushMessageQueue() {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift()!
      this.ws?.send(JSON.stringify(message))
    }
  }

  // 重连
  private scheduleReconnect(url: string) {
    // 实现重连逻辑...
  }

  // 消息处理
  private handleMessage(data: any) {
    // 分发到处理器...
  }
}
```

#### 3.2 后端WebSocket优化

```go
// backend/internal/websocket/optimized_hub.go

package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type OptimizedHub struct {
    clients       map[*Client]bool
    broadcast     chan []byte
    register      chan *Client
    unregister    chan *Client
    mu            sync.RWMutex

    // 心跳配置
    pingInterval  time.Duration
    pongWait      time.Duration

    // 状态广播节流
    lastBroadcast time.Time
    minInterval   time.Duration  // 最小广播间隔
    pendingState  []byte
}

func NewOptimizedHub() *OptimizedHub {
    return &OptimizedHub{
        clients:      make(map[*Client]bool),
        broadcast:    make(chan []byte, 256),
        register:     make(chan *Client),
        unregister:   make(chan *Client),
        pingInterval: 5 * time.Second,
        pongWait:     15 * time.Second,
        minInterval:  50 * time.Millisecond,  // 20 FPS 最大
    }
}

func (h *OptimizedHub) Run() {
    ticker := time.NewTicker(h.minInterval)
    defer ticker.Stop()

    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            delete(h.clients, client)
            h.mu.Unlock()

        case <-ticker.C:
            // 节流广播待处理状态
            if h.pendingState != nil {
                h.broadcastState(h.pendingState)
                h.pendingState = nil
            }
        }
    }
}

// 节流的状态广播
func (h *OptimizedHub) BroadcastStateThrottled(state []byte) {
    h.pendingState = state
}

func (h *OptimizedHub) broadcastState(state []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    for client := range h.clients {
        select {
        case client.send <- state:
        default:
            // 客户端缓冲区满，跳过
        }
    }
}

// 增强的客户端结构
type OptimizedClient struct {
    hub         *OptimizedHub
    conn        *websocket.Conn
    send        chan []byte
    lastPing    time.Time
    latency     time.Duration  // 客户端延迟
}

func (c *OptimizedClient) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(c.hub.pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.latency = time.Since(c.lastPing)
        c.conn.SetReadDeadline(time.Now().Add(c.hub.pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        var msg map[string]interface{}
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        // 处理心跳
        if msg["type"] == "ping" {
            c.lastPing = time.Now()
            c.send <- []byte(`{"type":"pong","timestamp":` +
                string(time.Now().UnixMilli()) + `}`)
            continue
        }

        // 处理其他消息...
    }
}
```

---

### 四、对话系统优化

#### 4.1 对话状态机

```typescript
// frontend/src/game/dialogue/DialogueManager.ts

type DialogueState = 'idle' | 'approaching' | 'talking' | 'ending'

interface DialogueSession {
  npcId: string
  npcName: string
  state: DialogueState
  currentDialogue: string
  options: string[]
  friendship: number
  startTime: number
}

export class DialogueManager {
  private session: DialogueSession | null = null
  private state: DialogueState = 'idle'

  // 开始对话
  async startDialogue(npcId: string, gameStore: GameStore): Promise<boolean> {
    if (this.state !== 'idle') return false

    const npcs = gameStore.npcs
    const player = gameStore.player
    if (!player) return false

    // 找到NPC
    const npc = npcs.find(n => n.id === npcId)
    if (!npc) return false

    // 检查距离
    const distance = this.getDistance(player.position, npc.position)
    if (distance > 3) {
      // 太远，需要走近
      this.state = 'approaching'
      await this.approachNPC(npcId, gameStore)
    }

    // 开始对话
    this.state = 'talking'
    this.session = {
      npcId,
      npcName: npc.name,
      state: 'talking',
      currentDialogue: '',
      options: [],
      friendship: npc.friendship,
      startTime: Date.now()
    }

    // 请求对话内容
    const result = await gameStore.talk(npcId)
    if (result && this.session) {
      this.session.currentDialogue = gameStore.messages[gameStore.messages.length - 1] || '...'
    }

    return true
  }

  // 自动走向NPC
  private async approachNPC(npcId: string, gameStore: GameStore): Promise<void> {
    // 实现自动寻路接近NPC
  }

  // 结束对话
  endDialogue() {
    if (this.session) {
      this.state = 'ending'
      // 可以在这里添加结束动画
      setTimeout(() => {
        this.state = 'idle'
        this.session = null
      }, 300)
    }
  }

  // 获取当前状态
  getState(): DialogueState {
    return this.state
  }

  // 获取当前会话
  getSession(): DialogueSession | null {
    return this.session
  }

  private getDistance(pos1: { x: number; y: number }, pos2: { x: number; y: number }): number {
    return Math.abs(pos1.x - pos2.x) + Math.abs(pos1.y - pos2.y)
  }
}
```

#### 4.2 后端对话处理优化

```go
// backend/internal/services/ai/dialogue_optimized.go

package ai

import (
    "context"
    "encoding/json"
    "sync"
    "time"

    "stardew-agent/internal/models"
)

type DialogueCache struct {
    mu      sync.RWMutex
    entries map[string]*CachedDialogue
    ttl     time.Duration
}

type CachedDialogue struct {
    dialogue   string
    expiresAt  time.Time
}

// 对话请求上下文
type DialogueContext struct {
    NPCID          string
    NPCName        string
    PlayerName     string
    Friendship     int
    TimeOfDay      string
    Weather        string
    RecentEvents   []string
    ConversationID string  // 用于连续对话
}

// 流式对话响应
type StreamDialogueResponse struct {
    Text      string `json:"text"`
    Done      bool   `json:"done"`
    Emotion   string `json:"emotion,omitempty"`
    Action    string `json:"action,omitempty"`
}

// 生成对话（支持流式）
func (s *DialogueService) GenerateDialogueStream(
    ctx context.Context,
    req *DialogueContext,
) (<-chan StreamDialogueResponse, error) {

    responseChan := make(chan StreamDialogueResponse, 10)

    go func() {
        defer close(responseChan)

        // 构建系统提示
        systemPrompt := s.buildSystemPrompt(req)
        userPrompt := s.buildUserPrompt(req)

        // 调用LLM（流式）
        stream, err := s.llmClient.GenerateStream(ctx, systemPrompt, userPrompt)
        if err != nil {
            responseChan <- StreamDialogueResponse{
                Text: "...",
                Done: true,
            }
            return
        }

        // 处理流式响应
        for chunk := range stream {
            responseChan <- StreamDialogueResponse{
                Text: chunk,
                Done: false,
            }
        }

        responseChan <- StreamDialogueResponse{
            Done: true,
        }
    }()

    return responseChan, nil
}

// 构建系统提示
func (s *DialogueService) buildSystemPrompt(req *DialogueContext) string {
    profile, ok := s.profiles[req.NPCID]
    if !ok {
        return "你是一个友好的NPC。"
    }

    // 基于友谊度调整对话深度
    friendshipLevel := req.Friendship / 250

    prompt := "你是" + profile.Name + "，"
    prompt += "角色定位：" + profile.Role + "。\n"
    prompt += "性格特点：" + joinTraits(profile.Traits) + "。\n"

    if friendshipLevel >= 4 {
        prompt += "你们是好朋友，可以分享更私人的话题。\n"
        if len(profile.Secrets) > 0 {
            prompt += "你可以暗示一些秘密：" + profile.Secrets[0] + "\n"
        }
    }

    prompt += "说话风格：" + profile.SpeechStyle + "\n"
    prompt += "保持回复简短（2-3句话），符合角色性格。"

    return prompt
}
```

---

### 五、完整集成示例

#### 5.1 更新后的GameLoop

```typescript
// frontend/src/game/engine/OptimizedGameLoop.ts

export class OptimizedGameLoop {
  private renderer: Renderer
  private lastTime = 0
  private running = false

  // 新增组件
  private prediction: ClientPrediction
  private interpolator: StateInterpolator
  private npcMovement: NPCMovementManager
  private dialogueManager: DialogueManager

  constructor(canvas: HTMLCanvasElement) {
    this.renderer = new Renderer(canvas)
    this.prediction = new ClientPrediction()
    this.interpolator = new StateInterpolator()
    this.npcMovement = new NPCMovementManager(new Pathfinder(48, 48))
    this.dialogueManager = new DialogueManager()
  }

  private loop = () => {
    if (!this.running) return

    const currentTime = performance.now()
    const deltaTime = (currentTime - this.lastTime) / 1000
    this.lastTime = currentTime

    // 1. 更新NPC移动
    this.npcMovement.update(deltaTime)

    // 2. 获取插值状态
    const interpolatedState = this.interpolator.getInterpolatedState()

    // 3. 应用客户端预测
    const finalState = this.applyPrediction(interpolatedState)

    // 4. 渲染
    if (finalState) {
      this.renderer.render(finalState, this.npcMovement)
    }

    requestAnimationFrame(this.loop)
  }

  // WebSocket消息处理
  onServerMessage(data: any) {
    switch (data.type) {
      case 'state':
        this.interpolator.addSnapshot(data.data)
        break

      case 'npc_move':
        this.npcMovement.updateNPCPosition(data.npcId, data.position)
        break

      case 'action_result':
        // 确认预测
        this.prediction.onServerState(data.state)
        break
    }
  }

  // 玩家输入处理
  async handleMove(direction: Direction) {
    // 1. 立即预测并更新本地状态
    const action = { type: 'move', params: { direction } }
    this.prediction.predictAction(action, this.getCurrentState())

    // 2. 发送到服务器
    await this.gameStore.move(direction)
  }
}
```

---

## 实施优先级

| 优先级 | 任务 | 预估时间 |
|--------|------|----------|
| P0 | 客户端预测系统 | 2h |
| P0 | 状态插值系统 | 2h |
| P1 | NPC平滑移动 | 3h |
| P1 | A*寻路算法 | 2h |
| P1 | WebSocket心跳优化 | 1h |
| P2 | 对话管理器 | 2h |
| P2 | 延迟补偿显示 | 2h |

---

*文档版本: 1.0*
*创建者: Game Engine Developer*
*日期: 2026-02-12*
