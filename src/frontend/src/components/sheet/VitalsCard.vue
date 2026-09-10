<script setup lang="ts">
import { computed } from 'vue'
import type { CharacterModel } from '../../constants/placeholders'
import TooltipField from '../ui/TooltipField.vue'
import SectionCard from '../ui/SectionCard.vue'

const props = defineProps<{
  character: CharacterModel
  computedWalk: number
  computedRun: number
  computedMaxWounds: number
  movementPenalty?: number
  travelFatigue?: number
}>()

const effectiveMovement = computed(() => {
  const base = Number(props.character.movement) || 0
  const penalty = props.movementPenalty || 0
  return Math.max(0, base - penalty)
})
</script>

<template>
  <SectionCard title="Movement, Fate & Wounds" full-height>
    <!-- Movement Row -->
    <div class="mb-3 pa-2 rounded border bg-surface-variant">
      <div class="text-caption font-weight-bold text-primary mb-1">Movement</div>
      <v-row dense class="mb-1">
        <v-col cols="4">
          <v-text-field
            :model-value="effectiveMovement"
            label="Move"
            type="number"
            variant="outlined"
            density="compact"
            hide-details
            @update:model-value="character.movement = Number($event) || 0"
          />
        </v-col>
        <v-col cols="4">
          <TooltipField :model-value="computedWalk" label="Walk (yd)" readonly variant="filled" tooltip="Walk distance = Move x 2 yards" />
        </v-col>
        <v-col cols="4">
          <TooltipField :model-value="computedRun" label="Run (yd)" readonly variant="filled" tooltip="Run distance = Move x 4 yards" />
        </v-col>
      </v-row>

      <!-- Encumbrance penalties under Movement -->
      <div v-if="(movementPenalty || 0) > 0" class="text-caption text-error font-weight-bold mt-1 px-1">
        Base Move: {{ character.movement }} (-{{ movementPenalty }} Encumbrance)
      </div>
      <div v-if="(travelFatigue || 0) > 0" class="text-caption text-warning font-weight-black mt-1 px-1 d-flex align-center">
        <v-icon icon="mdi-alert" size="small" class="mr-1" />
        Travel Fatigue: +{{ travelFatigue }}
      </div>
    </div>

    <!-- Fate & Fortune (Single row with tight spacing and subtle separator) -->
    <div class="mb-3 pa-2 rounded border bg-surface-variant">
      <div class="text-caption font-weight-bold text-primary mb-1">Fate & Fortune</div>
      <div class="d-flex align-center" style="gap: 4px;">
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fate" label="Fate" type="number" tooltip="Current Fate points (survive fatal blows)" />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fateMax" label="Max Fate" type="number" tooltip="Max Fate pool" />
        </div>
        <v-divider vertical class="mx-1 align-self-stretch" style="opacity: 0.3;" />
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fortune" label="Fortune" type="number" tooltip="Current Fortune points (re-roll tests)" />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fortuneMax" label="Max Fortune" type="number" tooltip="Max Fortune pool" />
        </div>
      </div>
    </div>

    <!-- Wounds (Single row with tight spacing and subtle separator for Hardy) -->
    <div class="mb-3 pa-2 rounded border bg-surface-variant">
      <div class="text-caption font-weight-bold text-primary mb-1">Wounds</div>
      <div class="d-flex align-center" style="gap: 4px;">
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField
            v-model.number="character.wounds.current"
            label="Current"
            type="number"
            tooltip="Current Wounds"
          />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField
            :model-value="computedMaxWounds"
            label="Max"
            readonly
            variant="filled"
            tooltip="Max Wounds = SB + (TB x 2) + WPB + Hardy"
            field-class="font-weight-bold"
          />
        </div>
        <v-divider vertical class="mx-1 align-self-stretch" style="opacity: 0.3;" />
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField
            v-model.number="character.wounds.hardy"
            label="Hardy"
            type="number"
            tooltip="Hardy talent bonus wounds"
          />
        </div>
      </div>
    </div>

    <!-- Experience (XP) Row -->
    <div class="d-flex gap-1 align-center pt-2 border-t">
      <TooltipField v-model.number="character.xp.current" label="Current XP" type="number" tooltip="Experience points available to spend" />
      <TooltipField v-model.number="character.xp.spent" label="Spent XP" type="number" tooltip="Total experience points spent" />
      <v-tooltip text="Total XP earned over career (Current + Spent)" location="top" :open-on-focus="false">
        <template #activator="{ props: tProps }">
          <v-chip v-bind="tProps" color="primary" variant="tonal" class="font-weight-bold ml-1">
            {{ (character.xp.current || 0) + (character.xp.spent || 0) }} Total
          </v-chip>
        </template>
      </v-tooltip>
    </div>
  </SectionCard>
</template>

<style scoped>
.gap-1 { gap: 4px; }
</style>
