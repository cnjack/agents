import { useGameStore } from '@/stores/gameStore'
import type { GameState } from '@/types/game'

let ws: WebSocket | null = null
let reconnectAttempts = 0
const maxReconnectAttempts = 5

export function connectWebSocket() {
  const store = useGameStore()

  if (ws) {
    ws.close()
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/v1/ws`

  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    console.log('WebSocket connected')
    store.setConnected(true)
    reconnectAttempts = 0
  }

  ws.onclose = () => {
    console.log('WebSocket disconnected')
    store.setConnected(false)

    // Attempt reconnect
    if (reconnectAttempts < maxReconnectAttempts) {
      reconnectAttempts++
      setTimeout(() => connectWebSocket(), 1000 * reconnectAttempts)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    store.setError('WebSocket connection error')
  }

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)

      if (data.type === 'state' && data.data) {
        store.updateState(data.data as GameState)
      } else if (data.type === 'action_result') {
        if (data.state) {
          store.updateState(data.state as GameState)
        }
        if (data.message) {
          store.addMessage(data.message)
        }
      } else if (data.type === 'npc_movement') {
        // Handle NPC movement update for smooth animation
        store.updateNPCMovement(data.data)
      } else {
        // Direct state update
        store.updateState(data as GameState)
      }
    } catch (err) {
      console.error('Failed to parse WebSocket message:', err)
    }
  }
}

export function disconnectWebSocket() {
  if (ws) {
    ws.close()
    ws = null
  }
}

export function sendAction(action: unknown) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(action))
  } else {
    console.warn('WebSocket not connected')
  }
}

export function isWebSocketConnected(): boolean {
  return ws !== null && ws.readyState === WebSocket.OPEN
}
