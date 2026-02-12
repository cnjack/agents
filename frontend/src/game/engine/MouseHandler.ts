import type { Position, NPCState, PlayerState, GameState } from '@/types/game'

const TILE_SIZE = 48
const INTERACTION_RANGE = 3 // Maximum tiles for interaction

export interface MouseEvent {
  screenX: number
  screenY: number
  gameX: number
  gameY: number
}

export interface HoveredAgent {
  id: string
  type: 'npc' | 'player'
  distance: number
  position: Position
}

export class MouseHandler {
  private canvas: HTMLCanvasElement
  private hoveredAgent: HoveredAgent | null = null
  private mousePosition: MouseEvent | null = null

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas
  }

  // Convert screen coordinates to game world coordinates
  screenToGame(screenX: number, screenY: number, cameraX: number, cameraY: number): Position {
    return {
      x: Math.floor((screenX + cameraX) / TILE_SIZE),
      y: Math.floor((screenY + cameraY) / TILE_SIZE)
    }
  }

  // Get the agent (NPC or player) at the clicked position
  getAgentAtPosition(
    gameX: number,
    gameY: number,
    state: GameState
  ): HoveredAgent | null {
    // Check NPCs first (player renders on top, but we check NPCs first for interaction)
    for (const npc of state.npcs) {
      if (npc.position.x === gameX && npc.position.y === gameY) {
        return {
          id: npc.id,
          type: 'npc',
          distance: this.getDistance(state.player.position, npc.position),
          position: npc.position
        }
      }
    }

    // Check player
    if (state.player.position.x === gameX && state.player.position.y === gameY) {
      return {
        id: 'player',
        type: 'player',
        distance: 0,
        position: state.player.position
      }
    }

    return null
  }

  // Get all agents within interaction range of the player
  getAgentsInRange(state: GameState, range: number = INTERACTION_RANGE): HoveredAgent[] {
    const agents: HoveredAgent[] = []
    const playerPos = state.player.position

    for (const npc of state.npcs) {
      const distance = this.getDistance(playerPos, npc.position)
      if (distance <= range) {
        agents.push({
          id: npc.id,
          type: 'npc',
          distance,
          position: npc.position
        })
      }
    }

    return agents.sort((a, b) => a.distance - b.distance)
  }

  // Check if an agent is within interaction range
  isInRange(agent: HoveredAgent, playerPos: Position, range: number = INTERACTION_RANGE): boolean {
    return this.getDistance(playerPos, agent.position) <= range
  }

  // Calculate Manhattan distance between two positions
  private getDistance(pos1: Position, pos2: Position): number {
    return Math.abs(pos1.x - pos2.x) + Math.abs(pos1.y - pos2.y)
  }

  // Update the currently hovered agent
  updateHover(screenX: number, screenY: number, cameraX: number, cameraY: number, state: GameState): HoveredAgent | null {
    const gamePos = this.screenToGame(screenX, screenY, cameraX, cameraY)
    this.mousePosition = { screenX, screenY, gameX: gamePos.x, gameY: gamePos.y }

    const agent = this.getAgentAtPosition(gamePos.x, gamePos.y, state)
    if (agent && this.isInRange(agent, state.player.position)) {
      this.hoveredAgent = agent
    } else {
      this.hoveredAgent = null
    }

    return this.hoveredAgent
  }

  // Get the current hovered agent
  getHoveredAgent(): HoveredAgent | null {
    return this.hoveredAgent
  }

  // Get current mouse position
  getMousePosition(): MouseEvent | null {
    return this.mousePosition
  }

  // Get screen coordinates for a game position (for rendering highlights)
  gameToScreen(gameX: number, gameY: number, cameraX: number, cameraY: number): { x: number; y: number } {
    return {
      x: gameX * TILE_SIZE - cameraX,
      y: gameY * TILE_SIZE - cameraY
    }
  }
}
