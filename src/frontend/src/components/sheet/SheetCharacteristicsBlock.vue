<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import CharacteristicCard from '../ui/CharacteristicCard.vue'

const props = defineProps<{
  character: CharacterModel
  getCharCurrent: (code: string) => number
  agilityPenalty?: number
}>()

function isImportant(code: string): boolean {
  return props.character.importantCharacteristics?.includes(code) || false
}

function toggleImportant(code: string) {
  if (!props.character.importantCharacteristics) {
    props.character.importantCharacteristics = []
  }
  const idx = props.character.importantCharacteristics.indexOf(code)
  if (idx >= 0) {
    props.character.importantCharacteristics.splice(idx, 1)
  } else {
    props.character.importantCharacteristics.push(code)
  }
}

function getEffectiveScore(code: string): number {
  const base = props.getCharCurrent(code)
  if (code === 'Ag' && props.agilityPenalty) {
    return Math.max(0, base - props.agilityPenalty)
  }
  return base
}
</script>

<template>
  <v-card color="surface" elevation="2" class="mb-4 pa-4 rounded-lg border">
    <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-3">
      Characteristics Grid
    </div>
    <v-row dense>
      <v-col
        v-for="(stat, code) in character.characteristics"
        :key="code"
        cols="6"
        sm="4"
        md="2"
        lg="1"
        class="flex-grow-1"
      >
        <CharacteristicCard
          :code="code as string"
          :name="stat.name"
          :hint="stat.hint"
          :initial="stat.initial"
          :advances="stat.advances"
          :current="getEffectiveScore(code as string)"
          show-important
          :is-important="isImportant(code as string)"
          :penalty="code === 'Ag' ? (agilityPenalty || 0) : 0"
          :penalty-tooltip="code === 'Ag' && agilityPenalty ? `Base: ${getCharCurrent(code as string)} (Penalized -${agilityPenalty} by Encumbrance)` : ''"
          @update:initial="stat.initial = $event"
          @update:advances="stat.advances = $event"
          @toggle-important="toggleImportant(code as string)"
        />
      </v-col>
    </v-row>
  </v-card>
</template>
