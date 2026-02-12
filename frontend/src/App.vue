<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useGameStore } from './stores/gameStore'
import { connectWebSocket, disconnectWebSocket } from './api/websocket'
import { GameLoop } from './game/engine/GameLoop'
import { InputHandler } from './game/engine/InputHandler'
import GameCanvas from './components/GameCanvas.vue'
import StatusPanel from './components/StatusPanel.vue'
import QuestPanel from './components/QuestPanel.vue'
import TimeDisplay from './components/TimeDisplay.vue'
import NPCDialog from './components/NPCDialog.vue'

const store = useGameStore()
const gameLoop = ref<GameLoop | null>(null)
const inputHandler = ref<InputHandler | null>(null)
const statePollInterval = ref<number | null>(null)

// NPC dialog state
const dialogVisible = ref(false)
const dialogNpcId = ref('')

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

// Handle NPC click from canvas
const handleNpcClick = (npcId: string) => {
  dialogNpcId.value = npcId
  dialogVisible.value = true
  // Trigger the talk action
  store.talk(npcId)
}

// Handle dialog close
const handleDialogClose = () => {
  dialogVisible.value = false
}
</script>

<template>
  <div class="game-container">
    <!-- Main Game Area -->
    <div class="game-viewport">
      <GameCanvas @ready="handleGameReady" @npc-click="handleNpcClick" />
    </div>

    <!-- UI Sidebar -->
    <div class="ui-sidebar">
      <TimeDisplay />
      <StatusPanel />
      <QuestPanel />
    </div>

    <!-- NPC Dialog -->
    <NPCDialog
      :npc-id="dialogNpcId"
      :visible="dialogVisible"
      @close="handleDialogClose"
    />
  </div>
</template>
