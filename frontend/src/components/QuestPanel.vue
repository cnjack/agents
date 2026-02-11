<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/gameStore'
import type { QuestData, Objective } from '@/types/game'

const store = useGameStore()

const activeQuests = computed(() => store.activeQuests)
const availableQuests = computed(() => store.availableQuests)

const getObjectiveProgress = (objective: Objective): string => {
  return `${objective.progress}/${objective.count}`
}

const isObjectiveComplete = (objective: Objective): boolean => {
  return objective.progress >= objective.count
}

const acceptQuest = async (questId: string) => {
  await store.acceptQuest(questId)
}

const getQuestGiverName = (giver: string): string => {
  const names: Record<string, string> = {
    lewis: 'Lewis',
    pierre: 'Pierre',
    robin: 'Robin',
    haley: 'Haley',
    willy: 'Willy',
  }
  return names[giver] || giver
}
</script>

<template>
  <div class="ui-panel">
    <div class="ui-panel-title">Quests</div>

    <!-- Active Quests -->
    <div v-if="activeQuests.length > 0" class="mb-4">
      <div class="text-xs text-green-400 mb-2">Active</div>
      <div
        v-for="quest in activeQuests"
        :key="quest.id"
        class="quest-item"
      >
        <div class="quest-name">{{ quest.name }}</div>
        <div class="quest-description">{{ quest.description }}</div>
        <div
          v-for="obj in quest.objectives"
          :key="obj.id"
          class="quest-progress"
          :class="{ 'text-yellow-400': isObjectiveComplete(obj) }"
        >
          {{ isObjectiveComplete(obj) ? '✓' : '○' }}
          {{ obj.description }}: {{ getObjectiveProgress(obj) }}
        </div>
      </div>
    </div>

    <!-- Available Quests -->
    <div v-if="availableQuests.length > 0">
      <div class="text-xs text-yellow-400 mb-2">Available</div>
      <div
        v-for="quest in availableQuests"
        :key="quest.id"
        class="quest-item cursor-pointer hover:bg-black/20"
        @click="acceptQuest(quest.id)"
      >
        <div class="quest-name">{{ quest.name }}</div>
        <div class="quest-description">{{ quest.description }}</div>
        <div class="text-xs text-gray-500 mt-1">
          From: {{ getQuestGiverName(quest.giver) }}
        </div>
        <div class="text-xs text-green-400 mt-1">
          Reward: {{ quest.rewards.gold }}g
        </div>
      </div>
    </div>

    <!-- No Quests -->
    <div v-if="activeQuests.length === 0 && availableQuests.length === 0" class="text-center text-gray-500 text-xs py-4">
      No quests available
    </div>
  </div>
</template>
