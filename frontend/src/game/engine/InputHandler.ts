import { useGameStore } from '@/stores/gameStore'
import type { Direction } from '@/types/game'

export type KeyHandler = (key: string) => void

export class InputHandler {
  private store: ReturnType<typeof useGameStore>
  private keysPressed: Set<string> = new Set()
  private moveDelay = 120 // ms between moves
  private lastMoveTime = 0
  private moveInterval: number | null = null

  constructor(store: ReturnType<typeof useGameStore>) {
    this.store = store
  }

  attach() {
    window.addEventListener('keydown', this.handleKeyDown)
    window.addEventListener('keyup', this.handleKeyUp)
  }

  detach() {
    window.removeEventListener('keydown', this.handleKeyDown)
    window.removeEventListener('keyup', this.handleKeyUp)
    if (this.moveInterval !== null) {
      clearInterval(this.moveInterval)
      this.moveInterval = null
    }
  }

  private handleKeyDown = async (event: KeyboardEvent) => {
    // Ignore if already pressed (prevent repeat)
    if (this.keysPressed.has(event.code)) return
    this.keysPressed.add(event.code)

    const direction = this.getDirectionFromKey(event.code)

    if (direction) {
      event.preventDefault()
      // Execute first move immediately
      await this.store.move(direction)
      this.lastMoveTime = Date.now()

      // Start continuous movement
      if (this.moveInterval !== null) {
        clearInterval(this.moveInterval)
      }

      this.moveInterval = window.setInterval(async () => {
        const now = Date.now()
        if (now - this.lastMoveTime >= this.moveDelay) {
          // Check if key is still pressed
          if (this.keysPressed.has(event.code)) {
            await this.store.move(direction)
            this.lastMoveTime = now
          }
        }
      }, this.moveDelay / 2)

      return
    }

    // Tool selection (1-5)
    if (event.code >= 'Digit1' && event.code <= 'Digit5') {
      event.preventDefault()
      const toolIndex = parseInt(event.code.replace('Digit', '')) - 1
      const tools = this.store.player?.tools
      if (tools && tools[toolIndex]) {
        console.log('Select tool:', tools[toolIndex])
      }
      return
    }

    // Action keys
    switch (event.code) {
      case 'Space':
        event.preventDefault()
        await this.handleAction()
        break

      case 'KeyE':
        event.preventDefault()
        await this.handleInteract()
        break

      case 'KeyQ':
        event.preventDefault()
        await this.store.harvest()
        break

      case 'Escape':
        // Close dialogs, etc.
        break
    }
  }

  private handleKeyUp = (event: KeyboardEvent) => {
    this.keysPressed.delete(event.code)

    // Stop continuous movement when movement key is released
    const direction = this.getDirectionFromKey(event.code)
    if (direction && this.moveInterval !== null) {
      clearInterval(this.moveInterval)
      this.moveInterval = null
    }
  }

  private getDirectionFromKey(code: string): Direction | null {
    switch (code) {
      case 'ArrowUp':
      case 'KeyW':
        return 'up'
      case 'ArrowDown':
      case 'KeyS':
        return 'down'
      case 'ArrowLeft':
      case 'KeyA':
        return 'left'
      case 'ArrowRight':
      case 'KeyD':
        return 'right'
      default:
        return null
    }
  }

  private async handleAction() {
    const player = this.store.player
    if (!player) return

    const tool = player.active_tool
    if (tool) {
      await this.store.useTool(tool)
    }
  }

  private async handleInteract() {
    const player = this.store.player
    const npcs = this.store.npcs
    if (!player) return

    // Check for nearby NPC
    const facingPos = this.getFacingPosition(player.position.x, player.position.y, player.direction)
    const nearbyNpc = npcs.find(
      npc => npc.position.x === facingPos.x && npc.position.y === facingPos.y
    )

    if (nearbyNpc) {
      await this.store.talk(nearbyNpc.id)
    }
  }

  private getFacingPosition(x: number, y: number, direction: Direction): { x: number; y: number } {
    switch (direction) {
      case 'up': return { x, y: y - 1 }
      case 'down': return { x, y: y + 1 }
      case 'left': return { x: x - 1, y }
      case 'right': return { x: x + 1, y }
    }
  }

  isKeyPressed(code: string): boolean {
    return this.keysPressed.has(code)
  }
}
