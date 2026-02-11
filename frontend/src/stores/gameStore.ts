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
    resetGame,
  }
})
