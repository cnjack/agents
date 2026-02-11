// Types matching the backend models (simplified for community-focused agent)

export interface Position {
  x: number
  y: number
}

export type Direction = 'up' | 'down' | 'left' | 'right'
export type Season = 'spring' | 'summer' | 'fall' | 'winter'
export type TileType = 'grass' | 'dirt' | 'tilled' | 'watered' | 'path' | 'water'
export type CropStage = 'seed' | 'sprout' | 'growing' | 'mature'

export interface Crop {
  stage: CropStage
  type?: string
  daysUntilMature?: number
}

export interface FarmTile {
  type: TileType
  crop?: { stage: CropStage }
}

export interface Farm {
  width: number
  height: number
  tiles: FarmTile[][]
}

export interface InventoryItem {
  id: string
  name: string
  type: string
  quantity: number
  price: number
}

export interface PlayerState {
  position: Position
  direction: Direction
  gold: number
  gifts: InventoryItem[]
  friendship: Record<string, number>  // NPC ID -> friendship points
}

export interface ScheduleEntry {
  hour: number
  minute: number
  location: string
  position: Position
}

export interface NPCState {
  id: string
  name: string
  position: Position
  direction: Direction
  friendship: number
  max_friendship: number
  location: string
  dialogue: string[]
  schedule?: ScheduleEntry[]
}

export interface Objective {
  id: string
  description: string
  type: string
  target: string
  count: number
  progress: number
  completed: boolean
}

export interface Reward {
  gold: number
  items: InventoryItem[]
}

export type QuestStatus = 'available' | 'active' | 'completed' | 'failed'

export interface QuestData {
  id: string
  name: string
  description: string
  objectives: Objective[]
  rewards: Reward
  status: QuestStatus
  giver: string
  deadline: number
}

export interface TimeState {
  day: number
  hour: number
  minute: number
  season: Season
  year: number
  paused: boolean
}

export interface GameState {
  player: PlayerState
  time: TimeState
  npcs: NPCState[]
  quests: QuestData[]
  map_width: number
  map_height: number
  game_over: boolean
  day_ended: boolean
  farm?: Farm
}

export interface Observation {
  player: PlayerState
  time: TimeState
  nearby_npcs: NPCState[]
  active_quests: QuestData[]
  messages: string[]
}

export interface Action {
  type: ActionType
  params: ActionParams
}

export type ActionType =
  | 'move'
  | 'talk'
  | 'give_gift'
  | 'accept_quest'
  | 'wait'
  | 'sleep'

export interface ActionParams {
  direction?: Direction
  target?: string
  npc?: string
  gift_item?: string
  quest_id?: string
  ticks?: number
}

export interface ActionResult {
  success: boolean
  message: string
  new_state?: GameState
  events?: GameEvent[]
}

export interface GameEvent {
  type: string
  data: unknown
  timestamp: number
}
