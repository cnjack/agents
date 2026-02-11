<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import { useGameStore } from '@/stores/gameStore'

const store = useGameStore()

const props = defineProps<{
  npcId: string
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'sendMessage', message: string): void
}>()

// NPC data
const npc = computed(() => {
  return store.npcs.find(n => n.id === props.npcId)
})

// NPC personality data (matching NPC_PROFILES.md)
const npcPersonalities: Record<string, {
  role: string
  traits: string[]
  speech_style: string
  mood_icons: Record<string, string>
}> = {
  lewis: {
    role: 'Mayor',
    traits: ['friendly', 'responsible', 'nostalgic'],
    speech_style: '热情但略带官腔',
    mood_icons: { happy: '😊', neutral: '😐', sad: '😔', excited: '🤩', angry: '😤' }
  },
  pierre: {
    role: 'Shopkeeper',
    traits: ['business_minded', 'friendly', 'ambitious'],
    speech_style: '热情好客，善于推销',
    mood_icons: { happy: '😊', neutral: '😐', sad: '😔', excited: '🤑', angry: '😤' }
  },
  robin: {
    role: 'Carpenter',
    traits: ['hardworking', 'creative', 'direct'],
    speech_style: '直接、热情、专业',
    mood_icons: { happy: '😊', neutral: '😐', sad: '😔', excited: '🛠️', angry: '😤' }
  },
  haley: {
    role: 'Villager',
    traits: ['photography_lover', 'gradually_warming'],
    speech_style: '初期傲慢冷淡，逐渐温柔',
    mood_icons: { happy: '📸', neutral: '😒', sad: '😔', excited: '✨', angry: '😤' }
  },
  willy: {
    role: 'Fisherman',
    traits: ['friendly', 'wise', 'patient'],
    speech_style: '温和、缓慢、充满智慧',
    mood_icons: { happy: '🐟', neutral: '😐', sad: '😔', excited: '🎣', angry: '😤' }
  }
}

const personality = computed(() => {
  return npcPersonalities[props.npcId] || {
    role: 'Villager',
    traits: ['friendly'],
    speech_style: '',
    mood_icons: { happy: '😊', neutral: '😐', sad: '😔', excited: '✨', angry: '😤' }
  }
})

// Current dialogue and AI response state
const currentDialogue = ref<string>('')
const displayedText = ref<string>('')
const isTyping = ref(false)
const playerInput = ref<string>('')
const inputFocused = ref(false)
const currentMood = ref<string>('neutral')

// Chat history for the session
interface ChatMessage {
  type: 'npc' | 'player'
  text: string
  timestamp: number
}
const chatHistory = ref<ChatMessage[]>([])

// Simulated mood based on conversation
const moodValue = computed(() => {
  const baseValue = 50
  const friendshipBonus = (npc.value?.friendship || 0) / 10
  return Math.min(100, Math.max(0, baseValue + friendshipBonus))
})

// Calculate hearts display
const heartsDisplay = computed(() => {
  if (!npc.value) return []
  const maxHearts = Math.floor(npc.value.max_friendship / 250)
  const filledHearts = Math.floor(npc.value.friendship / 250)
  return Array(maxHearts).fill(0).map((_, i) => ({
    filled: i < filledHearts,
    index: i
  }))
})

// Heart fill percentage for partial hearts
const heartFillPercent = computed(() => {
  if (!npc.value) return 0
  const remainder = npc.value.friendship % 250
  return (remainder / 250) * 100
})

// Typewriter effect for dialogue
const typewriterEffect = async (text: string) => {
  isTyping.value = true
  displayedText.value = ''

  for (let i = 0; i < text.length; i++) {
    displayedText.value += text[i]
    await new Promise(resolve => setTimeout(resolve, 30))
  }

  isTyping.value = false
}

// Show initial dialogue
const showDialogue = () => {
  if (npc.value && npc.value.dialogue.length > 0) {
    currentDialogue.value = npc.value.dialogue[
      Math.floor(Math.random() * npc.value.dialogue.length)
    ]
    chatHistory.value = []
    chatHistory.value.push({
      type: 'npc',
      text: currentDialogue.value,
      timestamp: Date.now()
    })
    typewriterEffect(currentDialogue.value)
  }
}

// Send player message
const sendMessage = async () => {
  if (!playerInput.value.trim()) return

  const message = playerInput.value.trim()
  chatHistory.value.push({
    type: 'player',
    text: message,
    timestamp: Date.now()
  })
  playerInput.value = ''

  // Emit to parent for AI processing
  emit('sendMessage', message)

  // Simulate AI response (placeholder - will be replaced by actual AI integration)
  await simulateAIResponse(message)
}

// Simulate AI response (placeholder)
const simulateAIResponse = async (_message: string) => {
  // This will be replaced by actual AI integration
  const responses = [
    "That's interesting! Tell me more.",
    "I see what you mean.",
    "That's a great point!",
    "Hmm, let me think about that...",
    "I appreciate you sharing that with me."
  ]
  const response = responses[Math.floor(Math.random() * responses.length)]

  // Update mood based on conversation
  const moods = ['happy', 'neutral', 'excited', 'neutral', 'happy']
  currentMood.value = moods[Math.floor(Math.random() * moods.length)]

  chatHistory.value.push({
    type: 'npc',
    text: response,
    timestamp: Date.now()
  })
  typewriterEffect(response)
}

// Close dialog
const close = () => {
  emit('close')
}

// Get mood icon
const getMoodIcon = (mood: string) => {
  return personality.value.mood_icons[mood] || '😐'
}

// Get mood color class
const getMoodColorClass = (mood: string) => {
  const colors: Record<string, string> = {
    happy: 'mood-happy',
    neutral: 'mood-neutral',
    sad: 'mood-sad',
    excited: 'mood-excited',
    angry: 'mood-angry'
  }
  return colors[mood] || 'mood-neutral'
}

// Watch for visibility changes
watch(() => props.visible, (visible) => {
  if (visible) {
    showDialogue()
    currentMood.value = 'neutral'
    nextTick(() => {
      const inputEl = document.querySelector('.player-input-field') as HTMLInputElement
      if (inputEl) inputEl.focus()
    })
  }
})

// Handle keyboard input
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  } else if (e.key === 'Escape') {
    close()
  }
}
</script>

<template>
  <Transition name="dialog-fade">
    <div v-if="visible && npc" class="dialog-overlay" @click.self="close">
      <div class="dialog-container">
        <!-- NPC Portrait and Info Header -->
        <div class="dialog-header">
          <div class="npc-portrait">
            <div class="portrait-frame">
              <div class="portrait-placeholder">
                {{ npc.name.charAt(0) }}
              </div>
            </div>
            <!-- Mood indicator -->
            <div class="mood-indicator" :class="getMoodColorClass(currentMood)">
              <span class="mood-icon">{{ getMoodIcon(currentMood) }}</span>
              <span class="mood-label">{{ currentMood }}</span>
            </div>
          </div>

          <div class="npc-info">
            <div class="npc-name">{{ npc.name }}</div>
            <div class="npc-role">{{ personality.role }}</div>
            <div class="npc-traits">
              <span
                v-for="trait in personality.traits.slice(0, 3)"
                :key="trait"
                class="trait-tag"
              >
                {{ trait }}
              </span>
            </div>
          </div>

          <!-- Hearts Display -->
          <div class="friendship-display">
            <div class="hearts-container">
              <div
                v-for="heart in heartsDisplay"
                :key="heart.index"
                class="heart-wrapper"
              >
                <span class="heart empty">♡</span>
                <span
                  class="heart filled"
                  :style="{ opacity: heart.filled ? 1 : 0 }"
                >❤</span>
                <!-- Partial heart fill -->
                <span
                  v-if="!heart.filled && heart.index === heartsDisplay.filter(h => h.filled).length"
                  class="heart partial"
                  :style="{ clipPath: `inset(0 ${100 - heartFillPercent}% 0 0)` }"
                >❤</span>
              </div>
            </div>
            <div class="friendship-label">
              {{ Math.floor(npc.friendship / 250) }}/{{ Math.floor(npc.max_friendship / 250) }} Hearts
            </div>
          </div>

          <!-- Close button -->
          <button class="close-btn" @click="close" title="Close (Esc)">
            ✕
          </button>
        </div>

        <!-- Chat History -->
        <div class="chat-container" ref="chatContainer">
          <div class="chat-messages">
            <div
              v-for="(msg, index) in chatHistory"
              :key="index"
              :class="['message', msg.type]"
            >
              <div class="message-sender">
                {{ msg.type === 'npc' ? npc.name : 'You' }}
              </div>
              <div class="message-bubble">
                <span v-if="msg.type === 'npc' && index === chatHistory.length - 1 && isTyping">
                  {{ displayedText }}<span class="typing-cursor">|</span>
                </span>
                <span v-else>{{ msg.text }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Speech Bubble (Current NPC dialogue with typewriter) -->
        <div v-if="displayedText && chatHistory.length <= 1" class="speech-bubble-container">
          <div class="speech-bubble">
            <div class="bubble-content">
              {{ displayedText }}<span v-if="isTyping" class="typing-cursor">|</span>
            </div>
            <div class="bubble-tail"></div>
          </div>
        </div>

        <!-- Player Input -->
        <div class="input-container">
          <div class="input-wrapper" :class="{ focused: inputFocused }">
            <input
              v-model="playerInput"
              type="text"
              class="player-input-field"
              placeholder="Type your message... (Enter to send)"
              @focus="inputFocused = true"
              @blur="inputFocused = false"
              @keydown="handleKeydown"
              :disabled="isTyping"
            />
            <button
              class="send-btn"
              @click="sendMessage"
              :disabled="!playerInput.trim() || isTyping"
            >
              <span class="send-icon">➤</span>
            </button>
          </div>
          <div class="input-hints">
            <span class="hint">Press <kbd>Enter</kbd> to send</span>
            <span class="hint">Press <kbd>Esc</kbd> to close</span>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="quick-actions">
          <button class="quick-action-btn" @click="playerInput = 'Hello!'">
            👋 Hello
          </button>
          <button class="quick-action-btn" @click="playerInput = 'How are you?'">
            💭 How are you?
          </button>
          <button class="quick-action-btn" @click="playerInput = 'Goodbye!'">
            👋 Goodbye
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* Dialog Overlay */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

/* Dialog Container */
.dialog-container {
  background: linear-gradient(180deg, #2d1b0e 0%, #1a0f08 100%);
  border: 4px solid #8b6914;
  border-radius: 8px;
  width: 100%;
  max-width: 700px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow:
    0 0 0 2px #4a3728,
    0 10px 40px rgba(0, 0, 0, 0.5),
    inset 0 1px 0 rgba(255, 215, 0, 0.2);
}

/* Dialog Header */
.dialog-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 2px solid #4a3728;
  position: relative;
}

/* NPC Portrait */
.npc-portrait {
  position: relative;
  flex-shrink: 0;
}

.portrait-frame {
  width: 80px;
  height: 80px;
  border: 3px solid #ffd700;
  border-radius: 8px;
  background: linear-gradient(135deg, #4a3728 0%, #2d1b0e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    inset 0 2px 4px rgba(0, 0, 0, 0.3),
    0 2px 8px rgba(0, 0, 0, 0.3);
}

.portrait-placeholder {
  font-family: 'Press Start 2P', monospace;
  font-size: 32px;
  color: #ffd700;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.5);
}

/* Mood Indicator */
.mood-indicator {
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 10px;
  font-family: 'Press Start 2P', monospace;
  border: 2px solid;
}

.mood-happy {
  background: rgba(76, 175, 80, 0.2);
  border-color: #4caf50;
  color: #4caf50;
}

.mood-neutral {
  background: rgba(158, 158, 158, 0.2);
  border-color: #9e9e9e;
  color: #9e9e9e;
}

.mood-sad {
  background: rgba(33, 150, 243, 0.2);
  border-color: #2196f3;
  color: #2196f3;
}

.mood-excited {
  background: rgba(255, 193, 7, 0.2);
  border-color: #ffc107;
  color: #ffc107;
}

.mood-angry {
  background: rgba(244, 67, 54, 0.2);
  border-color: #f44336;
  color: #f44336;
}

.mood-icon {
  font-size: 14px;
}

.mood-label {
  text-transform: uppercase;
}

/* NPC Info */
.npc-info {
  flex: 1;
}

.npc-name {
  font-family: 'Press Start 2P', monospace;
  font-size: 14px;
  color: #ffd700;
  margin-bottom: 4px;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.5);
}

.npc-role {
  font-size: 10px;
  color: #aaa;
  margin-bottom: 8px;
}

.npc-traits {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.trait-tag {
  font-size: 8px;
  padding: 2px 6px;
  background: rgba(255, 215, 0, 0.1);
  border: 1px solid #8b6914;
  border-radius: 4px;
  color: #cd853f;
}

/* Friendship Display */
.friendship-display {
  text-align: center;
  padding: 8px;
}

.hearts-container {
  display: flex;
  gap: 2px;
  justify-content: center;
  margin-bottom: 4px;
}

.heart-wrapper {
  position: relative;
  width: 20px;
  height: 20px;
}

.heart {
  position: absolute;
  top: 0;
  left: 0;
  font-size: 18px;
  transition: all 0.3s ease;
}

.heart.empty {
  color: #333;
}

.heart.filled {
  color: #ff4757;
  text-shadow: 0 0 8px rgba(255, 71, 87, 0.5);
}

.heart.partial {
  color: #ff4757;
}

.friendship-label {
  font-size: 8px;
  color: #888;
  font-family: 'Press Start 2P', monospace;
}

/* Close Button */
.close-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  background: rgba(244, 67, 54, 0.2);
  border: 2px solid #f44336;
  border-radius: 4px;
  color: #f44336;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f44336;
  color: white;
}

/* Chat Container */
.chat-container {
  flex: 1;
  min-height: 200px;
  max-height: 300px;
  overflow-y: auto;
  padding: 16px;
  background: rgba(0, 0, 0, 0.2);
}

.chat-messages {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message {
  display: flex;
  flex-direction: column;
  max-width: 85%;
}

.message.npc {
  align-self: flex-start;
}

.message.player {
  align-self: flex-end;
}

.message-sender {
  font-size: 8px;
  color: #888;
  margin-bottom: 4px;
  font-family: 'Press Start 2P', monospace;
}

.message.player .message-sender {
  text-align: right;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 10px;
  line-height: 1.6;
  font-family: 'Press Start 2P', monospace;
}

.message.npc .message-bubble {
  background: linear-gradient(135deg, #4a3728 0%, #3d2a1a 100%);
  border: 2px solid #8b6914;
  color: white;
  border-bottom-left-radius: 4px;
}

.message.player .message-bubble {
  background: linear-gradient(135deg, #2d5016 0%, #1a3a0a 100%);
  border: 2px solid #4a7c23;
  color: white;
  border-bottom-right-radius: 4px;
}

.typing-cursor {
  animation: blink 0.8s infinite;
  color: #ffd700;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* Speech Bubble (for single dialogue) */
.speech-bubble-container {
  padding: 0 16px;
}

.speech-bubble {
  position: relative;
  background: rgba(0, 0, 0, 0.9);
  border: 3px solid #ffd700;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 8px;
}

.bubble-content {
  font-family: 'Press Start 2P', monospace;
  font-size: 10px;
  line-height: 2;
  color: white;
}

.bubble-tail {
  position: absolute;
  bottom: -12px;
  left: 24px;
  width: 0;
  height: 0;
  border-left: 12px solid transparent;
  border-right: 12px solid transparent;
  border-top: 12px solid #ffd700;
}

.bubble-tail::after {
  content: '';
  position: absolute;
  top: -14px;
  left: -10px;
  width: 0;
  height: 0;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-top: 10px solid rgba(0, 0, 0, 0.9);
}

/* Input Container */
.input-container {
  padding: 16px;
  background: rgba(0, 0, 0, 0.3);
  border-top: 2px solid #4a3728;
}

.input-wrapper {
  display: flex;
  gap: 8px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.4);
  border: 2px solid #4a3728;
  border-radius: 8px;
  transition: all 0.2s;
}

.input-wrapper.focused {
  border-color: #ffd700;
  box-shadow: 0 0 8px rgba(255, 215, 0, 0.3);
}

.player-input-field {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: white;
  font-family: 'Press Start 2P', monospace;
  font-size: 10px;
}

.player-input-field::placeholder {
  color: #666;
}

.player-input-field:disabled {
  opacity: 0.5;
}

.send-btn {
  width: 40px;
  height: 32px;
  background: linear-gradient(180deg, #5a8f3a 0%, #3d6b28 100%);
  border: 2px solid #2d5018;
  border-radius: 4px;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: linear-gradient(180deg, #6aa04a 0%, #4d7b38 100%);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-icon {
  font-size: 14px;
}

.input-hints {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  justify-content: center;
}

.hint {
  font-size: 8px;
  color: #666;
}

.hint kbd {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid #444;
  border-radius: 3px;
  padding: 1px 4px;
  font-size: 8px;
}

/* Quick Actions */
.quick-actions {
  display: flex;
  gap: 8px;
  padding: 0 16px 16px;
  justify-content: center;
}

.quick-action-btn {
  padding: 8px 16px;
  background: rgba(255, 215, 0, 0.1);
  border: 2px solid #8b6914;
  border-radius: 4px;
  color: #cd853f;
  font-family: 'Press Start 2P', monospace;
  font-size: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-action-btn:hover {
  background: rgba(255, 215, 0, 0.2);
  border-color: #ffd700;
  color: #ffd700;
}

/* Transition */
.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: all 0.3s ease;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}

.dialog-fade-enter-from .dialog-container,
.dialog-fade-leave-to .dialog-container {
  transform: translateY(20px) scale(0.95);
}

/* Scrollbar */
.chat-container::-webkit-scrollbar {
  width: 8px;
}

.chat-container::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
}

.chat-container::-webkit-scrollbar-thumb {
  background: #4a3728;
  border-radius: 4px;
}

.chat-container::-webkit-scrollbar-thumb:hover {
  background: #8b6914;
}
</style>
