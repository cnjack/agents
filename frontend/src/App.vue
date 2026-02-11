<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useGameStore } from './stores/gameStore'
import { connectWebSocket, disconnectWebSocket } from './api/websocket'
import { GameLoop } from './game/engine/GameLoop'
import { InputHandler } from './game/engine/InputHandler'
import GameCanvas from './components/GameCanvas.vue'
import StatusPanel from './components/StatusPanel.vue'
import Inventory from './components/Inventory.vue'
import Toolbar from './components/Toolbar.vue'
import QuestPanel from './components/QuestPanel.vue'
import TimeDisplay from './components/TimeDisplay.vue'

const store = useGameStore()
const gameLoop = ref<GameLoop | null>(null)
const inputHandler = ref<InputHandler | null>(null)
const statePollInterval = ref<number | null>(null)

onMounted(async () => {
  // Fetch initial state
  await store.fetchState()

  // Connect WebSocket for real-time updates
  connectWebSocket()

  // Setup input handler
  inputHandler.value = new InputHandler(store)
  inputHandler.value.attach()

  // Poll for state updates every 500ms (for time updates)
  statePollInterval.value = window.setInterval(async () => {
    // Only poll if WebSocket is not connected
    // WebSocket will handle updates when connected
    await store.fetchState()
  }, 500)
})

onUnmounted(() => {
  disconnectWebSocket()

  if (statePollInterval.value !== null) {
    clearInterval(statePollInterval.value)
  }

  if (inputHandler.value) {
    inputHandler.value.detach()
  }

  if (gameLoop.value) {
    gameLoop.value.stop()
  }
})

const handleGameReady = (loop: GameLoop) => {
  gameLoop.value = loop
}
</script>

<template>
  <div class="game-container">
    <!-- Main Game Area -->
    <div class="game-viewport">
      <GameCanvas @ready="handleGameReady" />
      <Toolbar />
    </div>

    <!-- UI Sidebar -->
    <div class="ui-sidebar">
      <TimeDisplay />
      <StatusPanel />
      <Inventory />
      <QuestPanel />
    </div>
  </div>
</template>
