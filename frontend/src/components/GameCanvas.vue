<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useGameStore } from '@/stores/gameStore'
import { GameLoop } from '@/game/engine/GameLoop'
import { MouseHandler, type HoveredAgent } from '@/game/engine/MouseHandler'

const emit = defineEmits<{
  (e: 'ready', gameLoop: GameLoop): void
  (e: 'npc-click', npcId: string): void
}>()

const store = useGameStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const gameLoop = ref<GameLoop | null>(null)
const mouseHandler = ref<MouseHandler | null>(null)

// Track hovered agent for cursor changes
const hoveredAgent = ref<HoveredAgent | null>(null)

// Computed cursor style
const cursorStyle = computed(() => {
  return hoveredAgent.value ? 'pointer' : 'default'
})

const resizeCanvas = () => {
  if (!canvasRef.value) return

  const container = canvasRef.value.parentElement
  if (!container) return

  // Set canvas size to fit container with some padding
  const width = Math.floor(container.clientWidth * 0.9)
  const height = Math.floor(container.clientHeight * 0.9)

  canvasRef.value.width = width
  canvasRef.value.height = height

  if (gameLoop.value) {
    gameLoop.value.resize(width, height)
  }
}

// Handle mouse move - update hover state
const handleMouseMove = (event: MouseEvent) => {
  if (!canvasRef.value || !mouseHandler.value || !gameLoop.value || !store.gameState) return

  const rect = canvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  // Get camera position from game loop
  const camera = gameLoop.value.getCameraPosition()

  const hovered = mouseHandler.value.updateHover(x, y, camera.x, camera.y, store.gameState)
  hoveredAgent.value = hovered

  // Update renderer with hovered agent for highlight effect
  if (hovered) {
    gameLoop.value.setHoveredAgent(hovered.id)
  } else {
    gameLoop.value.setHoveredAgent(null)
  }
}

// Handle mouse click - emit npc-click event
const handleClick = (event: MouseEvent) => {
  if (!hoveredAgent.value) return

  if (hoveredAgent.value.type === 'npc') {
    emit('npc-click', hoveredAgent.value.id)
  }
}

// Prevent context menu on right-click
const handleContextMenu = (event: MouseEvent) => {
  event.preventDefault()
  return false
}

// Handle mouse leave - clear hover state
const handleMouseLeave = () => {
  hoveredAgent.value = null
  if (gameLoop.value) {
    gameLoop.value.setHoveredAgent(null)
  }
}

onMounted(() => {
  if (!canvasRef.value) return

  // Create game loop
  gameLoop.value = new GameLoop(canvasRef.value, () => store.gameState)
  gameLoop.value.start()

  // Create mouse handler
  mouseHandler.value = new MouseHandler(canvasRef.value)

  // Initial resize
  resizeCanvas()

  // Watch for resize
  window.addEventListener('resize', resizeCanvas)

  // Add mouse event listeners
  canvasRef.value.addEventListener('mousemove', handleMouseMove)
  canvasRef.value.addEventListener('mousedown', handleClick)
  canvasRef.value.addEventListener('mouseleave', handleMouseLeave)
  canvasRef.value.addEventListener('contextmenu', handleContextMenu)

  // Notify parent
  emit('ready', gameLoop.value)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeCanvas)
  if (gameLoop.value) {
    gameLoop.value.stop()
  }
  if (canvasRef.value) {
    canvasRef.value.removeEventListener('mousemove', handleMouseMove)
    canvasRef.value.removeEventListener('mousedown', handleClick)
    canvasRef.value.removeEventListener('mouseleave', handleMouseLeave)
    canvasRef.value.removeEventListener('contextmenu', handleContextMenu)
  }
})

// Expose mouse handler for parent access
defineExpose({
  mouseHandler
})
</script>

<template>
  <canvas
    ref="canvasRef"
    class="game-canvas"
    :style="{ cursor: cursorStyle }"
  />
</template>

<style scoped>
.game-canvas {
  display: block;
  background: #1a1a2e;
}
</style>
