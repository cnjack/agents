import type { Position, Direction, ScheduleEntry } from '@/types/game'

export class NPC {
  id: string
  name: string
  position: Position
  direction: Direction
  friendship: number
  maxFriendship: number
  location: string
  dialogue: string[]
  schedule: ScheduleEntry[]

  constructor(data: {
    id: string
    name: string
    position: Position
    direction: Direction
    friendship: number
    max_friendship: number
    location: string
    dialogue: string[]
    schedule?: ScheduleEntry[]
  }) {
    this.id = data.id
    this.name = data.name
    this.position = data.position
    this.direction = data.direction
    this.friendship = data.friendship
    this.maxFriendship = data.max_friendship
    this.location = data.location
    this.dialogue = data.dialogue
    this.schedule = data.schedule || []
  }

  getHearts(): number {
    return Math.floor(this.friendship / 250)
  }

  getMaxHearts(): number {
    return Math.floor(this.maxFriendship / 250)
  }

  getFriendshipLevel(): string {
    const hearts = this.getHearts()
    if (hearts >= 10) return 'Best Friend'
    if (hearts >= 8) return 'Close Friend'
    if (hearts >= 6) return 'Friend'
    if (hearts >= 4) return 'Friendly'
    if (hearts >= 2) return 'Acquaintance'
    if (hearts >= 1) return 'Newcomer'
    return 'Stranger'
  }

  isAtPosition(x: number, y: number): boolean {
    return this.position.x === x && this.position.y === y
  }

  getDialogue(): string {
    if (this.dialogue.length === 0) return '...'
    return this.dialogue[Math.floor(Math.random() * this.dialogue.length)]
  }
}
