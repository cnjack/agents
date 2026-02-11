<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/gameStore'

const store = useGameStore()

const timeString = computed(() => store.timeString)
const dateString = computed(() => store.dateString)

const isNight = computed(() => {
  const hour = store.time?.hour || 0
  return hour >= 18 || hour < 6
})

const dayPhase = computed(() => {
  const hour = store.time?.hour || 0
  if (hour >= 6 && hour < 12) return '🌅 Morning'
  if (hour >= 12 && hour < 18) return '☀️ Afternoon'
  if (hour >= 18 && hour < 21) return '🌆 Evening'
  return '🌙 Night'
})

const handleSleep = async () => {
  if (confirm('End the day and sleep?')) {
    await store.sleep()
  }
}
</script>

<template>
  <div class="time-display">
    <div class="time-clock">{{ timeString }}</div>
    <div class="time-date">{{ dateString }}</div>
    <div class="text-xs mt-2" :class="isNight ? 'text-blue-300' : 'text-yellow-300'">
      {{ dayPhase }}
    </div>
    <button
      class="pixel-btn mt-3 w-full text-xs"
      @click="handleSleep"
    >
      💤 Sleep
    </button>
  </div>
</template>
