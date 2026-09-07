<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import TooltipField from '../ui/TooltipField.vue'
import SectionCard from '../ui/SectionCard.vue'

defineProps<{
  character: CharacterModel
  computedWalk: number
  computedRun: number
  computedMaxWounds: number
}>()
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Movement, Fate & Wounds (Merged into same card) -->
    <v-col cols="12" md="4">
      <SectionCard title="Movement, Fate & Wounds" full-height>
        <v-row dense class="mb-2">
          <v-col cols="4">
            <v-text-field v-model.number="character.movement" label="Move" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="4">
            <TooltipField :model-value="computedWalk" label="Walk (yd)" readonly variant="filled" tooltip="Walk distance = Move x 2 yards" />
          </v-col>
          <v-col cols="4">
            <TooltipField :model-value="computedRun" label="Run (yd)" readonly variant="filled" tooltip="Run distance = Move x 4 yards" />
          </v-col>
        </v-row>

        <v-row dense class="mb-2">
          <v-col cols="6">
            <TooltipField v-model.number="character.fate" label="Fate" type="number" tooltip="Fate points allow surviving fatal blows" />
          </v-col>
          <v-col cols="6">
            <TooltipField v-model.number="character.fortune" label="Fortune" type="number" tooltip="Fortune points reset daily to re-roll tests" />
          </v-col>
        </v-row>

        <v-row dense class="mb-2">
          <v-col cols="6">
            <v-text-field v-model.number="character.wounds.current" label="Current Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="6">
            <TooltipField :model-value="computedMaxWounds" label="Max Wounds" readonly variant="filled" tooltip="Max Wounds = SB + (TB x 2) + WPB + Hardy" field-class="font-weight-bold" />
          </v-col>
        </v-row>

        <v-row dense>
          <v-col cols="12">
            <v-text-field v-model.number="character.wounds.hardy" label="Hardy Talent Bonus Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
        </v-row>
      </SectionCard>
    </v-col>

    <!-- Career Advances (Moved next to Wounds: 3 sets of 10 checkboxes labeled 2 to 4) -->
    <v-col cols="12" md="4">
      <SectionCard title="Career Advances" full-height>
        <!-- Set 2: Tier 2 (advances2) -->
        <div class="mb-3">
          <div class="d-flex justify-space-between align-center mb-1">
            <span class="text-caption font-weight-bold text-primary">Tier 2 Advances</span>
            <span class="text-caption text-medium-emphasis">{{ character.advances2.filter(Boolean).length }}/10</span>
          </div>
          <div class="d-flex align-center flex-wrap gap-1">
            <v-tooltip v-for="i in 10" :key="'adv2-' + i" :text="`Tier 2 Advance ${i}`" location="top">
              <template #activator="{ props: tProps }">
                <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: character.advances2[i - 1] }">
                  <input
                    type="checkbox"
                    :checked="character.advances2[i - 1]"
                    class="career-checkbox-input"
                    @change="character.advances2[i - 1] = !character.advances2[i - 1]"
                  />
                  <span class="career-checkbox-box">2</span>
                </label>
              </template>
            </v-tooltip>
          </div>
        </div>

        <!-- Set 3: Tier 3 (advances3) -->
        <div class="mb-3">
          <div class="d-flex justify-space-between align-center mb-1">
            <span class="text-caption font-weight-bold text-primary">Tier 3 Advances</span>
            <span class="text-caption text-medium-emphasis">{{ character.advances3.filter(Boolean).length }}/10</span>
          </div>
          <div class="d-flex align-center flex-wrap gap-1">
            <v-tooltip v-for="i in 10" :key="'adv3-' + i" :text="`Tier 3 Advance ${i}`" location="top">
              <template #activator="{ props: tProps }">
                <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: character.advances3[i - 1] }">
                  <input
                    type="checkbox"
                    :checked="character.advances3[i - 1]"
                    class="career-checkbox-input"
                    @change="character.advances3[i - 1] = !character.advances3[i - 1]"
                  />
                  <span class="career-checkbox-box">3</span>
                </label>
              </template>
            </v-tooltip>
          </div>
        </div>

        <!-- Set 4: Tier 4 (advances4) -->
        <div>
          <div class="d-flex justify-space-between align-center mb-1">
            <span class="text-caption font-weight-bold text-primary">Tier 4 Advances</span>
            <span class="text-caption text-medium-emphasis">{{ character.advances4.filter(Boolean).length }}/10</span>
          </div>
          <div class="d-flex align-center flex-wrap gap-1">
            <v-tooltip v-for="i in 10" :key="'adv4-' + i" :text="`Tier 4 Advance ${i}`" location="top">
              <template #activator="{ props: tProps }">
                <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: character.advances4[i - 1] }">
                  <input
                    type="checkbox"
                    :checked="character.advances4[i - 1]"
                    class="career-checkbox-input"
                    @change="character.advances4[i - 1] = !character.advances4[i - 1]"
                  />
                  <span class="career-checkbox-box">4</span>
                </label>
              </template>
            </v-tooltip>
          </div>
        </div>
      </SectionCard>
    </v-col>

    <!-- Armour Points per Location -->
    <v-col cols="12" md="4">
      <SectionCard title="Armour Points (AP)" full-height>
        <v-row dense>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.head" label="Head (01-09)" type="number" tooltip="Head protection (d100 roll 01-09)" />
          </v-col>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.leftArm" label="L. Arm (10-24)" type="number" tooltip="Left Arm protection (d100 roll 10-24)" />
          </v-col>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.rightArm" label="R. Arm (25-44)" type="number" tooltip="Right Arm protection (d100 roll 25-44)" />
          </v-col>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.body" label="Body (45-79)" type="number" tooltip="Body Torso protection (d100 roll 45-79)" />
          </v-col>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.leftLeg" label="L. Leg (80-89)" type="number" tooltip="Left Leg protection (d100 roll 80-89)" />
          </v-col>
          <v-col cols="4">
            <TooltipField v-model.number="character.armourPoints.rightLeg" label="R. Leg (90-00)" type="number" tooltip="Right Leg protection (d100 roll 90-00)" />
          </v-col>
          <v-col cols="12" class="mt-1">
            <TooltipField v-model.number="character.armourPoints.shield" label="Shield AP" type="number" tooltip="Shield AP Bonus" />
          </v-col>
        </v-row>
      </SectionCard>
    </v-col>
  </v-row>
</template>

<style scoped>
.gap-1 { gap: 4px; }

.career-checkbox-label {
  cursor: pointer;
  user-select: none;
}

.career-checkbox-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.career-checkbox-box {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 2px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 700;
  color: rgba(var(--v-theme-on-surface), 0.6);
  background: rgba(var(--v-theme-surface-variant), 0.4);
  transition: all 0.15s ease;
}

.career-checkbox-label.checked .career-checkbox-box {
  background: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-on-primary));
  border-color: rgb(var(--v-theme-primary));
}

.career-checkbox-label:hover .career-checkbox-box {
  border-color: rgb(var(--v-theme-primary));
}
</style>
