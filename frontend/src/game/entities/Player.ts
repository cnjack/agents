import type { Position, Direction, InventoryItem } from '@/types/game'

export class Player {
  position: Position
  direction: Direction
  gold: number
  gifts: InventoryItem[]
  friendship: Record<string, number>

  constructor(data: {
    position: Position
    direction: Direction
    gold: number
    gifts: InventoryItem[]
    friendship: Record<string, number>
  }) {
    this.position = data.position
    this.direction = data.direction
    this.gold = data.gold
    this.gifts = data.gifts
    this.friendship = data.friendship || {}
  }

  getFacingPosition(): Position {
    const pos = { ...this.position }
    switch (this.direction) {
      case 'up': pos.y--; break
      case 'down': pos.y++; break
      case 'left': pos.x--; break
      case 'right': pos.x++; break
    }
    return pos
  }

  hasGift(giftId: string): boolean {
    return this.gifts.some(gift => gift.id === giftId && gift.quantity > 0)
  }

  getGiftCount(giftId: string): number {
    const gift = this.gifts.find(g => g.id === giftId)
    return gift ? gift.quantity : 0
  }

  getFriendship(npcId: string): number {
    return this.friendship[npcId] || 0
  }
}
