import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GameState, Action, Direction } from '@/types/game'

export const useGameStore = defineStore('game', () => {
  // State
  const gameState = ref<GameState | null>(null)
  const isConnected = ref(false)
  const isLoading = ref(true)
  const error = ref<string | null>(null)
  const messages = ref<string[]>([])

  // Getters
  const player = computed(() => gameState.value?.player)
  const time = computed(() => gameState.value?.time)
  const npcs = computed(() => gameState.value?.npcs || [])
  const quests = computed(() => gameState.value?.quests || [])

  const activeQuests = computed(() =>
    quests.value.filter(q => q.status === 'active')
  )

  const availableQuests = computed(() =>
    quests.value.filter(q => q.status === 'available')
  )

  const completedQuests = computed(() =>
    quests.value.filter(q => q.status === 'completed')
  )

  const timeString = computed(() => {
    if (!time.value) return ''
    const hour = time.value.hour > 12 ? time.value.hour - 12 : time.value.hour
    const minute = time.value.minute.toString().padStart(2, '0')
    const period = time.value.hour >= 12 && time.value.hour < 24 ? 'PM' : 'AM'
    return `${hour}:${minute} ${period}`
  })

  const dateString = computed(() => {
    if (!time.value) return ''
    const season = time.value.season.charAt(0).toUpperCase() + time.value.season.slice(1)
    return `${season} ${time.value.day}, Year ${time.value.year}`
  })

  // Actions
  function updateState(state: GameState) {
    gameState.value = state
    isLoading.value = false
  }

  // Update single NPC position for smooth movement animation
  function updateNPCMovement(movement: {
    npc_id: string
    old_pos: { x: number; y: number }
    new_pos: { x: number; y: number }
    direction: string
    is_moving: boolean
    path_ended?: boolean
  }) {
    if (!gameState.value) return

    const npcIndex = gameState.value.npcs.findIndex(n => n.id === movement.npc_id)
    if (npcIndex === -1) return

    const npc = gameState.value.npcs[npcIndex]
    // Store last position for interpolation
    npc.last_position = { x: movement.old_pos.x, y: movement.old_pos.y }
    npc.position = { x: movement.new_pos.x, y: movement.new_pos.y }
    npc.direction = movement.direction as Direction
    npc.is_moving = movement.is_moving
  }

  function setConnected(connected: boolean) {
    isConnected.value = connected
  }

  function setError(err: string | null) {
    error.value = err
  }

  function addMessage(msg: string) {
    messages.value.push(msg)
    if (messages.value.length > 50) {
      messages.value.shift()
    }
  }

  function clearMessages() {
    messages.value = []
  }

  async function fetchState() {
    try {
      const response = await fetch('/api/v1/game/state')
      if (!response.ok) throw new Error('Failed to fetch state')
      const state = await response.json()
      updateState(state)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    }
  }

  async function executeAction(action: Action): Promise<boolean> {
    try {
      const response = await fetch('/api/v1/game/action', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(action),
      })
      const result = await response.json()

      if (result.success && result.new_state) {
        updateState(result.new_state)
      }

      if (result.message) {
        addMessage(result.message)
      }

      return result.success
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Action failed')
      return false
    }
  }

  async function move(direction: Direction): Promise<boolean> {
    return executeAction({ type: 'move', params: { direction } })
  }

  async function talk(npc: string): Promise<boolean> {
    return executeAction({ type: 'talk', params: { npc } })
  }

  async function giveGift(npc: string, giftItem: string): Promise<boolean> {
    return executeAction({ type: 'give_gift', params: { npc, gift_item: giftItem } })
  }

  async function acceptQuest(questId: string): Promise<boolean> {
    return executeAction({ type: 'accept_quest', params: { quest_id: questId } })
  }

  async function wait(ticks: number): Promise<boolean> {
    return executeAction({ type: 'wait', params: { ticks } })
  }

  async function sleep(): Promise<boolean> {
    return executeAction({ type: 'sleep', params: {} })
  }

  async function useTool(tool: string): Promise<boolean> {
    return executeAction({ type: 'use_tool', params: { tool } })
  }

  async function harvest(): Promise<boolean> {
    return executeAction({ type: 'harvest', params: {} })
  }

  async function plant(seed: string): Promise<boolean> {
    return executeAction({ type: 'plant', params: { seed } })
  }

  async function buy(item: string, quantity: number): Promise<boolean> {
    return executeAction({ type: 'buy', params: { item, quantity } })
  }

  async function sell(item: string, quantity: number): Promise<boolean> {
    return executeAction({ type: 'sell', params: { item, quantity } })
  }

  async function resetGame(): Promise<void> {
    try {
      const response = await fetch('/api/v1/game/reset', { method: 'POST' })
      const result = await response.json()
      if (result.state) {
        updateState(result.state)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Reset failed')
    }
  }

  return {
    // State
    gameState,
    isConnected,
    isLoading,
    error,
    messages,
    // Getters
    player,
    time,
    npcs,
    quests,
    activeQuests,
    availableQuests,
    completedQuests,
    timeString,
    dateString,
    // Actions
    updateState,
    updateNPCMovement,
    setConnected,
    setError,
    addMessage,
    clearMessages,
    fetchState,
    executeAction,
    move,
    talk,
    giveGift,
    acceptQuest,
    wait,
    sleep,
    useTool,
    harvest,
    plant,
    buy,
    sell,
    resetGame,
  }
})
