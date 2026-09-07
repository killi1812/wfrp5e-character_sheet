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
    <!-- Movement, Fate, Wounds & Experience (4 cols) -->
    <v-col cols="12" md="4">
      <SectionCard title="Movement, Fate & Wounds" full-height>
        <!-- Movement Row -->
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

        <!-- Fate & Fortune Row -->
        <v-row dense class="mb-2">
          <v-col cols="6">
            <TooltipField v-model.number="character.fate" label="Fate" type="number" tooltip="Fate points allow surviving fatal blows" />
          </v-col>
          <v-col cols="6">
            <TooltipField v-model.number="character.fortune" label="Fortune" type="number" tooltip="Fortune points reset daily to re-roll tests" />
          </v-col>
        </v-row>

        <!-- Wounds Row -->
        <v-row dense class="mb-2">
          <v-col cols="6">
            <v-text-field v-model.number="character.wounds.current" label="Current Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="6">
            <TooltipField :model-value="computedMaxWounds" label="Max Wounds" readonly variant="filled" tooltip="Max Wounds = SB + (TB x 2) + WPB + Hardy" field-class="font-weight-bold" />
          </v-col>
        </v-row>

        <!-- Hardy Bonus Row -->
        <v-row dense class="mb-2">
          <v-col cols="12">
            <v-text-field v-model.number="character.wounds.hardy" label="Hardy Talent Bonus Wounds" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
        </v-row>

        <!-- Experience (XP) Row (Moved into Movement/Vitals box) -->
        <div class="d-flex gap-1 align-center mt-2 pt-2 border-t">
          <TooltipField v-model.number="character.xp.current" label="Current XP" type="number" tooltip="Experience points currently available to spend" />
          <TooltipField v-model.number="character.xp.spent" label="Spent XP" type="number" tooltip="Total experience points spent on advances" />
          <v-tooltip text="Total XP earned over career (Current + Spent)" location="top">
            <template #activator="{ props: tProps }">
              <v-chip v-bind="tProps" color="primary" variant="tonal" class="font-weight-bold ml-1">
                {{ (character.xp.current || 0) + (character.xp.spent || 0) }} Total
              </v-chip>
            </template>
          </v-tooltip>
        </div>
      </SectionCard>
    </v-col>

    <!-- Career & Advances (5 cols - expanded space) -->
    <v-col cols="12" md="5">
      <SectionCard title="Career & Advances" full-height>
        <!-- Class, Career, Status (Moved into Career box) -->
        <v-row dense class="mb-3">
          <v-col cols="4">
            <v-text-field v-model="character.class" label="Class" variant="outlined" density="compact" hide-details placeholder="Academic..." />
          </v-col>
          <v-col cols="5">
            <v-text-field v-model="character.career" label="Career" variant="outlined" density="compact" hide-details placeholder="Wizard..." />
          </v-col>
          <v-col cols="3">
            <v-text-field v-model="character.status" label="Status" variant="outlined" density="compact" hide-details placeholder="Silver 3" />
          </v-col>
        </v-row>

        <!-- Set 2: Tier 2 Advances (advances2) -->
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

        <!-- Set 3: Tier 3 Advances (advances3) -->
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

        <!-- Set 4: Tier 4 Advances (advances4) -->
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

    <!-- Armour Points per Location (3 cols - reduced space) -->
    <v-col cols="12" md="3">
      <SectionCard title="Armour Points (AP)" full-height>
        <!-- Anatomical Compact Body Layout -->
        <div class="ap-body-container py-1">
          <!-- 1. Head (Centered) -->
          <div class="d-flex justify-center mb-2">
            <v-tooltip text="Head protection (d100 roll 01-09)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">Head</span>
                  <span class="ap-roll">01-09</span>
                  <input
                    v-model.number="character.armourPoints.head"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>
          </div>

          <!-- 2. Left Arm & Right Arm -->
          <div class="d-flex justify-space-between align-center mb-2 px-2">
            <v-tooltip text="Left Arm protection (d100 roll 10-24)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">L. Arm</span>
                  <span class="ap-roll">10-24</span>
                  <input
                    v-model.number="character.armourPoints.leftArm"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>

            <v-tooltip text="Right Arm protection (d100 roll 25-44)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">R. Arm</span>
                  <span class="ap-roll">25-44</span>
                  <input
                    v-model.number="character.armourPoints.rightArm"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>
          </div>

          <!-- 3. Body (Centered) -->
          <div class="d-flex justify-center mb-2">
            <v-tooltip text="Body Torso protection (d100 roll 45-79)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">Body</span>
                  <span class="ap-roll">45-79</span>
                  <input
                    v-model.number="character.armourPoints.body"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>
          </div>

          <!-- 4. Left Leg & Right Leg -->
          <div class="d-flex justify-space-between align-center mb-2 px-2">
            <v-tooltip text="Left Leg protection (d100 roll 80-89)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">L. Leg</span>
                  <span class="ap-roll">80-89</span>
                  <input
                    v-model.number="character.armourPoints.leftLeg"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>

            <v-tooltip text="Right Leg protection (d100 roll 90-00)" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">R. Leg</span>
                  <span class="ap-roll">90-00</span>
                  <input
                    v-model.number="character.armourPoints.rightLeg"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>
          </div>

          <!-- 5. Shield (Centered) -->
          <div class="d-flex justify-center">
            <v-tooltip text="Shield AP Bonus" location="top">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="ap-node">
                  <span class="ap-name">Shield</span>
                  <span class="ap-roll">Bonus</span>
                  <input
                    v-model.number="character.armourPoints.shield"
                    type="number"
                    class="ap-field"
                    placeholder="0"
                  />
                </div>
              </template>
            </v-tooltip>
          </div>
        </div>
      </SectionCard>
    </v-col>
  </v-row>
</template>

<style scoped>
.gap-1 { gap: 4px; }

.border-t {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

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

.ap-node {
  display: flex;
  flex-direction: column;
  align-items: center;
  user-select: none;
}

.ap-name {
  font-size: 0.72rem;
  font-weight: 700;
  color: rgb(var(--v-theme-primary));
  line-height: 1.1;
}

.ap-roll {
  font-size: 0.65rem;
  color: rgba(var(--v-theme-on-surface), 0.6);
  line-height: 1;
  margin-bottom: 2px;
}

.ap-field {
  width: 44px;
  height: 32px;
  text-align: center;
  font-size: 1.1rem;
  font-weight: 800;
  border: 1.5px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
  background: rgba(var(--v-theme-surface-variant), 0.4);
  color: currentColor;
  outline: none;
  transition: all 0.15s ease;
}

.ap-field:focus {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 0 0 2px rgba(var(--v-theme-primary), 0.25);
}
</style>
