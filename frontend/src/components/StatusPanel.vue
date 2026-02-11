<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/gameStore'

const store = useGameStore()

const energyPercent = computed(() => store.energyPercent)
const gold = computed(() => store.player?.gold || 0)
</script>

<template>
  <div class="ui-panel">
    <div class="ui-panel-title">Status</div>

    <!-- Energy Bar -->
    <div class="mb-3">
      <div class="text-xs text-green-400 mb-1">Energy</div>
      <div class="status-bar">
        <div
          class="status-bar-fill energy-bar"
          :style="{ width: `${energyPercent}%` }"
        />
      </div>
      <div class="text-xs text-gray-400 mt-1">
        {{ store.player?.energy }}/{{ store.player?.max_energy }}
      </div>
    </div>

    <!-- Gold -->
    <div class="flex items-center gap-2">
      <span class="text-yellow-400 text-lg">💰</span>
      <span class="text-yellow-400 text-sm pixel-text">{{ gold }}</span>
    </div>

    <!-- Position -->
    <div class="mt-3 text-xs text-gray-500">
      Position: ({{ store.player?.position.x }}, {{ store.player?.position.y }})
    </div>
  </div>
</template>
