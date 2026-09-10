<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import VitalsCard from './VitalsCard.vue'
import CareersCard from './CareersCard.vue'

defineProps<{
  character: CharacterModel
  computedWalk: number
  computedRun: number
  computedMaxWounds: number
  movementPenalty?: number
  travelFatigue?: number
}>()
</script>

<template>
  <div class="d-flex flex-wrap gap-3 mb-4 vitals-careers-row">
    <!-- Movement, Fate & Wounds (30% on desktop) -->
    <div class="vitals-col">
      <VitalsCard
        :character="character"
        :computed-walk="computedWalk"
        :computed-run="computedRun"
        :computed-max-wounds="computedMaxWounds"
        :movement-penalty="movementPenalty"
        :travel-fatigue="travelFatigue"
      />
    </div>

    <!-- Career & Advances (70% on desktop) -->
    <div class="careers-col">
      <CareersCard :character="character" />
    </div>
  </div>
</template>

<style scoped>
.gap-3 { gap: 12px; }

.vitals-careers-row {
  display: flex;
}

.vitals-col {
  flex: 0 0 calc(30% - 6px);
  max-width: calc(30% - 6px);
}

.careers-col {
  flex: 0 0 calc(70% - 6px);
  max-width: calc(70% - 6px);
}

@media (max-width: 960px) {
  .vitals-col, .careers-col {
    flex: 0 0 100%;
    max-width: 100%;
  }
}
</style>
