<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

defineProps<{
  character: CharacterModel
  computedWalk: number
  computedRun: number
  computedMaxWounds: number
}>()
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Movement & Pools -->
    <v-col cols="12" md="6" lg="4">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg h-100 border">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-3">Movement & Fate/Resilience</div>
        <v-row dense class="mb-2">
          <v-col cols="4">
            <v-text-field v-model.number="character.movement" label="Move" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Walk distance = Move x 2 yards" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" :model-value="computedWalk" label="Walk (yd)" readonly variant="filled" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Run distance = Move x 4 yards" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" :model-value="computedRun" label="Run (yd)" readonly variant="filled" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
        </v-row>

        <v-row dense class="mb-2">
          <v-col cols="6">
            <v-tooltip text="Fate points allow surviving fatal blows" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.fate" label="Fate" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="6">
            <v-tooltip text="Fortune points reset daily to re-roll tests" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.fortune" label="Fortune" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
        </v-row>

        <v-row dense>
          <v-col cols="6">
            <v-tooltip text="Resilience score prevents corruption and mutation" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.resilience" label="Resilience" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="6">
            <v-tooltip text="Resolve score removes conditions in combat" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.resolve" label="Resolve" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
        </v-row>
      </v-card>
    </v-col>

    <!-- Wounds & Hardy -->
    <v-col cols="12" md="6" lg="4">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg h-100 border">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-3">Wounds & Hardy</div>
        <v-row dense class="mb-2">
          <v-col cols="6">
            <v-text-field v-model.number="character.wounds.current" label="Current Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="6">
            <v-tooltip text="Max Wounds = SB + (TB x 2) + WPB + Hardy" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" :model-value="computedMaxWounds" label="Max Wounds" readonly variant="filled" density="compact" hide-details class="font-weight-bold" />
              </template>
            </v-tooltip>
          </v-col>
        </v-row>

        <v-row dense>
          <v-col cols="12">
            <v-text-field v-model.number="character.wounds.hardy" label="Hardy Talent Bonus Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
        </v-row>
      </v-card>
    </v-col>

    <!-- Armour Points per Location -->
    <v-col cols="12" lg="4">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg h-100 border">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-2">Armour Points (AP)</div>
        <v-row dense>
          <v-col cols="4">
            <v-tooltip text="Head protection (d100 roll 01-09)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.head" label="Head (01-09)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Left Arm protection (d100 roll 10-24)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.leftArm" label="L. Arm (10-24)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Right Arm protection (d100 roll 25-44)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.rightArm" label="R. Arm (25-44)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Body Torso protection (d100 roll 45-79)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.body" label="Body (45-79)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Left Leg protection (d100 roll 80-89)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.leftLeg" label="L. Leg (80-89)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="4">
            <v-tooltip text="Right Leg protection (d100 roll 90-00)" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.rightLeg" label="R. Leg (90-00)" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
          <v-col cols="12" class="mt-1">
            <v-tooltip text="Shield AP Bonus" location="top">
              <template #activator="{ props: tProps }">
                <v-text-field v-bind="tProps" v-model.number="character.armourPoints.shield" label="Shield AP" type="number" variant="outlined" density="compact" hide-details />
              </template>
            </v-tooltip>
          </v-col>
        </v-row>
      </v-card>
    </v-col>
  </v-row>
</template>
