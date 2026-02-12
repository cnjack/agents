import type { GameState } from '@/types/game'
import { Renderer } from './Renderer'

export type GameLoopCallback = (deltaTime: number) => void

export class GameLoop {
  private renderer: Renderer
  private lastTime = 0
  private running = false
  private animationFrameId: number | null = null
  private getState: () => GameState | null
  private onTick: GameLoopCallback | null = null

  constructor(
    canvas: HTMLCanvasElement,
    getState: () => GameState | null
  ) {
    this.renderer = new Renderer(canvas)
    this.getState = getState
  }

  start() {
    if (this.running) return

    this.running = true
    this.lastTime = performance.now()
    this.loop()
  }

  stop() {
    this.running = false
    if (this.animationFrameId !== null) {
      cancelAnimationFrame(this.animationFrameId)
      this.animationFrameId = null
    }
  }

  resize(width: number, height: number) {
    this.renderer.resize(width, height)
  }

  setOnTick(callback: GameLoopCallback) {
    this.onTick = callback
  }

  // Get current camera position
  getCameraPosition(): { x: number; y: number } {
    return this.renderer.getCameraPosition()
  }

  // Get hovered agent for rendering highlights
  setHoveredAgent(agentId: string | null) {
    this.renderer.setHoveredAgent(agentId)
  }

  private loop = () => {
    if (!this.running) return

    const currentTime = performance.now()
    const deltaTime = (currentTime - this.lastTime) / 1000
    this.lastTime = currentTime

    // Call tick callback
    if (this.onTick) {
      this.onTick(deltaTime)
    }

    // Render current state
    const state = this.getState()
    if (state) {
      this.renderer.render(state)
    }

    // Schedule next frame
    this.animationFrameId = requestAnimationFrame(this.loop)
  }
}
