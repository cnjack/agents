<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGameStore } from '@/stores/gameStore'

const store = useGameStore()

const inventory = computed(() => store.player?.inventory || [])
const selectedItem = ref<string | null>(null)

const getItemIcon = (type: string): string => {
  const icons: Record<string, string> = {
    seed: '🌱',
    crop: '🌾',
    tool: '🔧',
    gift: '🎁',
    food: '🍎',
    resource: '📦',
  }
  return icons[type] || '📦'
}

const selectItem = (itemId: string) => {
  selectedItem.value = selectedItem.value === itemId ? null : itemId
}
</script>

<template>
  <div class="ui-panel">
    <div class="ui-panel-title">Inventory</div>

    <div class="inventory-grid">
      <div
        v-for="item in inventory"
        :key="item.id"
        class="inventory-slot"
        :class="{ selected: selectedItem === item.id }"
        @click="selectItem(item.id)"
      >
        <span class="text-xl">{{ getItemIcon(item.type) }}</span>
        <span class="inventory-count">{{ item.quantity }}</span>
      </div>

      <!-- Empty slots -->
      <div
        v-for="i in Math.max(0, 16 - inventory.length)"
        :key="`empty-${i}`"
        class="inventory-slot"
      />
    </div>

    <!-- Selected item details -->
    <div v-if="selectedItem" class="mt-3 p-2 bg-black/30 rounded">
      <div class="text-xs text-yellow-400">
        {{ inventory.find(i => i.id === selectedItem)?.name }}
      </div>
      <div class="text-xs text-gray-400 mt-1">
        Count: {{ inventory.find(i => i.id === selectedItem)?.quantity }}
      </div>
      <div class="text-xs text-green-400 mt-1">
        Value: {{ inventory.find(i => i.id === selectedItem)?.price }}g
      </div>
    </div>
  </div>
</template>
