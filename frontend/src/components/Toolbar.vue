<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGameStore } from '@/stores/gameStore'

const store = useGameStore()

const tools = computed(() => store.player?.tools || [])
const activeTool = computed(() => store.player?.active_tool || '')

const toolIcons: Record<string, string> = {
  hoe: '⛏️',
  watering_can: '💧',
  axe: '🪓',
  pickaxe: '⛏️',
  scythe: '🔪',
  fishing_rod: '🎣',
}

const toolNames: Record<string, string> = {
  hoe: 'Hoe',
  watering_can: 'Watering Can',
  axe: 'Axe',
  pickaxe: 'Pickaxe',
  scythe: 'Scythe',
  fishing_rod: 'Fishing Rod',
}

const selectTool = async (tool: string) => {
  // Update active tool via action
  // This would be a store action in a full implementation
  console.log('Selected tool:', tool)
}

const handleToolAction = async () => {
  if (activeTool.value) {
    await store.useTool(activeTool.value)
  }
}
</script>

<template>
  <div class="toolbar">
    <div
      v-for="(tool, index) in tools"
      :key="tool"
      class="tool-slot"
      :class="{ active: tool === activeTool }"
      @click="selectTool(tool)"
    >
      <span class="tool-icon">{{ toolIcons[tool] || '🔧' }}</span>
    </div>
  </div>
</template>
