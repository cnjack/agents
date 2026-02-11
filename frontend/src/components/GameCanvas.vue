<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useGameStore } from '@/stores/gameStore'
import { GameLoop } from '@/game/engine/GameLoop'

const emit = defineEmits<{
  (e: 'ready', gameLoop: GameLoop): void
}>()

const store = useGameStore()
const canvasRef = ref<HTMLCanvasElement | null>(null)
const gameLoop = ref<GameLoop | null>(null)

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

onMounted(() => {
  if (!canvasRef.value) return

  // Create game loop
  gameLoop.value = new GameLoop(canvasRef.value, () => store.gameState)
  gameLoop.value.start()

  // Initial resize
  resizeCanvas()

  // Watch for resize
  window.addEventListener('resize', resizeCanvas)

  // Notify parent
  emit('ready', gameLoop.value)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeCanvas)
  if (gameLoop.value) {
    gameLoop.value.stop()
  }
})
</script>

<template>
  <canvas
    ref="canvasRef"
    class="game-canvas"
  />
</template>

<style scoped>
.game-canvas {
  display: block;
  background: #1a1a2e;
}
</style>
