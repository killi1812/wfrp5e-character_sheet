<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { CharacterModel } from '../../../constants/placeholders'
import TooltipField from '../../ui/TooltipField.vue'
import SectionCard from '../../ui/SectionCard.vue'
import ChipBtn from '../../ui/ChipBtn.vue'

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

// Momentum session storage state
const MOMENTUM_KEY = 'wfrp_momentum'
function readMomentum(): boolean {
  try {
    return sessionStorage.getItem(MOMENTUM_KEY) === '1'
  } catch {
    return false
  }
}
const momentum = ref(readMomentum())
watch(momentum, (val) => {
  try {
    sessionStorage.setItem(MOMENTUM_KEY, val ? '1' : '0')
  } catch {
    // sessionStorage blocked or unavailable
  }
})

// XP Rework
const xpTotal = computed(() => (Number(props.character.xp?.current) || 0) + (Number(props.character.xp?.spent) || 0))
const xpRemaining = computed(() => Number(props.character.xp?.current) || 0)

const clampInt = (v: number) => Math.max(0, Math.trunc(Number(v) || 0))

function ensureXp() {
  if (!props.character.xp) {
    props.character.xp = { current: 0, spent: 0 }
  }
}

function setXpTotal(v: number) {
  ensureXp()
  props.character.xp.current = clampInt(v) - (Number(props.character.xp.spent) || 0)
}

function setXpSpent(v: number) {
  ensureXp()
  const total = xpTotal.value
  props.character.xp.spent = clampInt(v)
  props.character.xp.current = total - clampInt(v)
}

const adjustXp = (target: 'total' | 'spent', delta: number) => {
  ensureXp()
  if (target === 'total') {
    setXpTotal(xpTotal.value + delta)
  } else {
    setXpSpent((Number(props.character.xp.spent) || 0) + delta)
  }
}
</script>

<template>
  <SectionCard title="Movement, Fate & Wounds" full-height>
    <!-- Movement Row -->
    <div class="mb-3">
      <div class="text-caption font-weight-bold text-primary mb-1">Movement</div>
      <v-row dense class="mb-1">
        <v-col cols="4">
          <v-text-field :model-value="effectiveMovement" label="Move" type="number" variant="outlined" density="compact"
            hide-details @update:model-value="character.movement = Number($event) || 0" />
        </v-col>
        <v-col cols="4">
          <TooltipField :model-value="computedWalk" label="Walk (yd)" readonly variant="filled"
            tooltip="Walk distance = Move x 2 yards" />
        </v-col>
        <v-col cols="4">
          <TooltipField :model-value="computedRun" label="Run (yd)" readonly variant="filled"
            tooltip="Run distance = Move x 4 yards" />
        </v-col>
      </v-row>

      <!-- Encumbrance penalties under Movement -->
      <div v-if="(movementPenalty || 0) > 0" class="text-caption text-error font-weight-bold mt-1 px-1">
        Base Move: {{ character.movement }} (-{{ movementPenalty }} Encumbrance)
      </div>
      <div v-if="(travelFatigue || 0) > 0"
        class="text-caption text-warning font-weight-black mt-1 px-1 d-flex align-center">
        <v-icon icon="mdi-alert" size="small" class="mr-1" />
        Travel Fatigue: +{{ travelFatigue }}
      </div>
    </div>

    <!-- Momentum Section -->
    <div class="mb-3">
      <div class="d-flex align-center justify-space-between mb-1">
        <div class="text-caption font-weight-bold text-primary">Momentum</div>
        <v-tooltip text="Momentum is tracked in this browser tab only — it is not saved to the server" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <v-icon v-bind="tProps" icon="mdi-help-circle-outline" size="small" color="on-surface-variant" class="cursor-pointer" />
          </template>
        </v-tooltip>
      </div>
      <v-btn block class="momentum-btn" :class="{ 'momentum-btn--on': momentum }"
        color="secondary" :variant="momentum ? 'flat' : 'outlined'" :aria-pressed="momentum"
        @click="momentum = !momentum">
        <v-icon start size="28" class="momentum-icon">mdi-lightning-bolt</v-icon>
        <span class="momentum-title">Momentum</span>
        <span class="momentum-state">{{ momentum ? 'Active' : 'Off' }}</span>
      </v-btn>
    </div>

    <!-- Fate & Fortune (Single row with tight spacing and subtle separator) -->
    <div class="mb-3">
      <div class="text-caption font-weight-bold text-primary mb-1">Fate & Fortune</div>
      <div class="d-flex align-center" style="gap: 4px;">
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fate" label="Fate" type="number"
            tooltip="Current Fate points (survive fatal blows)" />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fateMax" label="Max Fate" type="number" tooltip="Max Fate pool" />
        </div>
        <v-divider vertical class="mx-1 align-self-stretch" style="opacity: 0.3;" />
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fortune" label="Fortune" type="number"
            tooltip="Current Fortune points (re-roll tests)" />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.fortuneMax" label="Max Fortune" type="number"
            tooltip="Max Fortune pool" />
        </div>
      </div>
    </div>

    <!-- Wounds (Single row with tight spacing and subtle separator for Hardy) -->
    <div class="mb-3">
      <div class="text-caption font-weight-bold text-primary mb-1">Wounds</div>
      <div class="d-flex align-center" style="gap: 4px;">
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.wounds.current" label="Current" type="number"
            tooltip="Current Wounds" />
        </div>
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField :model-value="computedMaxWounds" label="Max" readonly variant="filled"
            tooltip="Max Wounds = SB + (TB x 2) + WPB + Hardy" field-class="font-weight-bold" />
        </div>
        <v-divider vertical class="mx-1 align-self-stretch" style="opacity: 0.3;" />
        <div style="flex: 1 1 0; min-width: 0;">
          <TooltipField v-model.number="character.wounds.hardy" label="Hardy" type="number"
            tooltip="Hardy talent bonus wounds" />
        </div>
      </div>
    </div>

    <!-- Experience (XP) Section -->
    <div class="pt-2 border-t">
      <div class="d-flex align-center justify-space-between mb-2">
        <div class="text-caption font-weight-bold text-primary">Experience (XP)</div>
        <v-tooltip text="Remaining XP available to spend (Total - Spent)" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <v-chip v-bind="tProps" :color="xpRemaining < 0 ? 'error' : 'primary'" variant="tonal" class="font-weight-bold">
              Remaining: {{ xpRemaining }} XP
            </v-chip>
          </template>
        </v-tooltip>
      </div>

      <v-row dense>
        <!-- Total XP Column -->
        <v-col cols="6">
          <TooltipField :model-value="xpTotal" label="Total XP" type="number"
            tooltip="Total experience points earned"
            @update:model-value="setXpTotal(Number($event) || 0)" />
          <div class="d-flex flex-wrap gap-1 mt-1 justify-center">
            <ChipBtn color="success" size="default" label="+50" @click="adjustXp('total', 50)" />
            <ChipBtn color="success" size="default" label="+100" @click="adjustXp('total', 100)" />
            <ChipBtn color="success" size="default" label="+125" @click="adjustXp('total', 125)" />
          </div>
          <div class="d-flex flex-wrap gap-1 mt-1 justify-center">
            <ChipBtn color="error" size="default" label="-50" @click="adjustXp('total', -50)" />
            <ChipBtn color="error" size="default" label="-100" @click="adjustXp('total', -100)" />
            <ChipBtn color="error" size="default" label="-125" @click="adjustXp('total', -125)" />
          </div>
        </v-col>

        <!-- Spent XP Column -->
        <v-col cols="6">
          <TooltipField :model-value="character.xp?.spent || 0" label="Spent XP" type="number"
            tooltip="Total experience points spent"
            @update:model-value="setXpSpent(Number($event) || 0)" />
          <div class="d-flex flex-wrap gap-1 mt-1 justify-center">
            <ChipBtn color="success" size="default" label="+50" @click="adjustXp('spent', 50)" />
            <ChipBtn color="success" size="default" label="+100" @click="adjustXp('spent', 100)" />
            <ChipBtn color="success" size="default" label="+125" @click="adjustXp('spent', 125)" />
          </div>
          <div class="d-flex flex-wrap gap-1 mt-1 justify-center">
            <ChipBtn color="error" size="default" label="-50" @click="adjustXp('spent', -50)" />
            <ChipBtn color="error" size="default" label="-100" @click="adjustXp('spent', -100)" />
            <ChipBtn color="error" size="default" label="-125" @click="adjustXp('spent', -125)" />
          </div>
        </v-col>
      </v-row>
    </div>
  </SectionCard>
</template>

<style scoped>
.gap-1 {
  gap: 4px;
}

.border-t {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.cursor-pointer {
  cursor: pointer;
}

.momentum-btn {
  height: 56px !important;
  font-weight: 900;
  letter-spacing: 0.12em;
  border-width: 2px;
  text-transform: uppercase;
}

.momentum-btn:not(.momentum-btn--on) {
  border-style: dashed;
  opacity: 0.75;
}

.momentum-btn--on {
  box-shadow: 0 0 0 2px rgba(var(--v-theme-primary), 0.9),
              0 0 22px rgba(var(--v-theme-secondary), 0.65);
}

.momentum-btn--on .momentum-icon {
  animation: momentum-pulse 1.6s ease-in-out infinite;
}

.momentum-title {
  font-size: 0.95rem;
}

.momentum-state {
  margin-left: auto;
  font-size: 0.75rem;
  opacity: 0.9;
}

@keyframes momentum-pulse {
  50% {
    transform: scale(1.25);
  }
}

@media (prefers-reduced-motion: reduce) {
  .momentum-btn--on .momentum-icon {
    animation: none;
  }
}
</style>
