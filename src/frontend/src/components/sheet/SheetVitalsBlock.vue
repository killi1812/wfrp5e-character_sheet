<script setup lang="ts">
import { computed } from 'vue'
import type { CharacterModel, CareerEntry } from '../../constants/placeholders'
import { NEW_ITEM_TEMPLATES } from '../../constants/placeholders'
import TooltipField from '../ui/TooltipField.vue'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

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

function addCareer() {
  if (!props.character.careers) {
    props.character.careers = []
  }
  const isFirst = props.character.careers.length === 0
  const newC = NEW_ITEM_TEMPLATES.career()
  newC.active = isFirst
  props.character.careers.push(newC)
  if (isFirst) {
    syncActiveCareer(newC)
  }
}

function removeCareer(index: number) {
  const wasActive = props.character.careers[index].active
  props.character.careers.splice(index, 1)
  if (props.character.careers.length > 0) {
    if (wasActive) {
      props.character.careers[0].active = true
      syncActiveCareer(props.character.careers[0])
    }
  } else {
    props.character.career = ''
    props.character.class = ''
    props.character.status = ''
  }
}

function setActiveCareer(index: number) {
  props.character.careers.forEach((c, idx) => {
    c.active = idx === index
  })
  syncActiveCareer(props.character.careers[index])
}

function syncActiveCareer(c: CareerEntry) {
  props.character.career = c.career
  props.character.class = c.class
  props.character.status = c.status
}
</script>

<template>
  <div class="d-flex flex-wrap gap-3 mb-4 vitals-careers-row">
    <!-- Movement, Fate & Wounds (30% on desktop) -->
    <div class="vitals-col">
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
    </div>

    <!-- Career & Advances (70% on desktop) -->
    <div class="careers-col">
      <SectionCard title="Careers & Advances" full-height add-label="Add Career" @add="addCareer">
        <div v-if="character.careers.length === 0" class="text-caption text-medium-emphasis font-italic py-4 text-center">
          No careers added. Click "+ Add Career" to add one.
        </div>

        <div v-for="(c, idx) in character.careers" :key="idx" class="career-row-card mb-3 pa-3 rounded-lg border" :class="{ 'career-active': c.active }">
          <!-- Career Header info: Active toggle, Class, Career, Status, Delete -->
          <div class="d-flex flex-wrap align-center gap-2 mb-2">
            <!-- Active / Important marker -->
            <v-tooltip :text="c.active ? 'Current Active & Important Career' : 'Click to set as Active Career'" location="top" :open-on-focus="false">
              <template #activator="{ props: tProps }">
                <button
                  v-bind="tProps"
                  type="button"
                  class="career-active-btn d-flex align-center px-2 py-1 rounded"
                  :class="{ active: c.active }"
                  @click="setActiveCareer(idx)"
                >
                  <v-icon :icon="c.active ? 'mdi-star' : 'mdi-star-outline'" size="small" class="mr-1" />
                  <span class="text-caption font-weight-bold">{{ c.active ? 'Active' : 'Inactive' }}</span>
                </button>
              </template>
            </v-tooltip>

            <div style="width: 140px;">
              <v-text-field
                v-model="c.class"
                label="Class"
                variant="outlined"
                density="compact"
                hide-details
                placeholder="Academic..."
                @update:model-value="c.active && syncActiveCareer(c)"
              />
            </div>
            <div class="flex-grow-1" style="min-width: 150px;">
              <v-text-field
                v-model="c.career"
                label="Career"
                variant="outlined"
                density="compact"
                hide-details
                placeholder="Wizard..."
                @update:model-value="c.active && syncActiveCareer(c)"
              />
            </div>
            <div style="width: 150px;">
              <v-text-field
                v-model="c.status"
                label="Status"
                variant="outlined"
                density="compact"
                hide-details
                placeholder="Silver 3"
                @update:model-value="c.active && syncActiveCareer(c)"
              />
            </div>

            <DeleteRowBtn @delete="removeCareer(idx)" />
          </div>

          <!-- Tier Advances (Rank 2: 10 boxes in 2 rows of 5, Rank 3: 12 boxes in 2 rows of 6, Rank 4: 14 boxes in 2 rows of 7) -->
          <v-row dense class="pt-2 border-t">
            <!-- Tier 2 (10 boxes: 2 rows of 5) -->
            <v-col cols="12" md="4" class="tier-col">
              <div class="tier-card pa-2 rounded border bg-surface h-100">
                <div class="d-flex justify-space-between align-center mb-2">
                  <span class="text-caption font-weight-bold text-primary">Tier 2 Advances</span>
                  <v-chip size="x-small" variant="tonal" color="primary" class="font-weight-bold">
                    {{ c.advances2.filter(Boolean).length }}/10
                  </v-chip>
                </div>
                <div class="tier-boxes-grid grid-5">
                  <v-tooltip v-for="i in 10" :key="'adv2-' + i" :text="`Tier 2 Advance ${i}`" location="top" :open-on-focus="false">
                    <template #activator="{ props: tProps }">
                      <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: c.advances2[i - 1] }">
                        <input
                          type="checkbox"
                          :checked="c.advances2[i - 1]"
                          class="career-checkbox-input"
                          @change="c.advances2[i - 1] = !c.advances2[i - 1]"
                        />
                        <span class="career-checkbox-box"></span>
                      </label>
                    </template>
                  </v-tooltip>
                </div>
              </div>
            </v-col>

            <!-- Tier 3 (12 boxes: 2 rows of 6) -->
            <v-col cols="12" md="4" class="tier-col">
              <div class="tier-card pa-2 rounded border bg-surface h-100">
                <div class="d-flex justify-space-between align-center mb-2">
                  <span class="text-caption font-weight-bold text-primary">Tier 3 Advances</span>
                  <v-chip size="x-small" variant="tonal" color="primary" class="font-weight-bold">
                    {{ c.advances3.filter(Boolean).length }}/12
                  </v-chip>
                </div>
                <div class="tier-boxes-grid grid-6">
                  <v-tooltip v-for="i in 12" :key="'adv3-' + i" :text="`Tier 3 Advance ${i}`" location="top" :open-on-focus="false">
                    <template #activator="{ props: tProps }">
                      <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: c.advances3[i - 1] }">
                        <input
                          type="checkbox"
                          :checked="c.advances3[i - 1]"
                          class="career-checkbox-input"
                          @change="c.advances3[i - 1] = !c.advances3[i - 1]"
                        />
                        <span class="career-checkbox-box"></span>
                      </label>
                    </template>
                  </v-tooltip>
                </div>
              </div>
            </v-col>

            <!-- Tier 4 (14 boxes: 2 rows of 7) -->
            <v-col cols="12" md="4" class="tier-col">
              <div class="tier-card pa-2 rounded border bg-surface h-100">
                <div class="d-flex justify-space-between align-center mb-2">
                  <span class="text-caption font-weight-bold text-primary">Tier 4 Advances</span>
                  <v-chip size="x-small" variant="tonal" color="primary" class="font-weight-bold">
                    {{ c.advances4.filter(Boolean).length }}/14
                  </v-chip>
                </div>
                <div class="tier-boxes-grid grid-7">
                  <v-tooltip v-for="i in 14" :key="'adv4-' + i" :text="`Tier 4 Advance ${i}`" location="top" :open-on-focus="false">
                    <template #activator="{ props: tProps }">
                      <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: c.advances4[i - 1] }">
                        <input
                          type="checkbox"
                          :checked="c.advances4[i - 1]"
                          class="career-checkbox-input"
                          @change="c.advances4[i - 1] = !c.advances4[i - 1]"
                        />
                        <span class="career-checkbox-box"></span>
                      </label>
                    </template>
                  </v-tooltip>
                </div>
              </div>
            </v-col>
          </v-row>
        </div>
      </SectionCard>
    </div>
  </div>
</template>

<style scoped>
.gap-1 { gap: 4px; }
.gap-2 { gap: 8px; }
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

.border-t {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.career-row-card {
  background: rgba(var(--v-theme-surface-variant), 0.3);
  transition: all 0.2s ease;
}

.career-row-card.career-active {
  border-color: rgb(var(--v-theme-primary)) !important;
  box-shadow: 0 0 8px rgba(var(--v-theme-primary), 0.2);
  background: rgba(var(--v-theme-primary), 0.05);
}

.career-active-btn {
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background: rgba(var(--v-theme-surface), 0.8);
  color: rgba(var(--v-theme-on-surface), 0.7);
  cursor: pointer;
  transition: all 0.15s ease;
}

.career-active-btn.active {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-on-primary));
}

.tier-card {
  background: rgba(var(--v-theme-surface), 0.7);
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
}

.tier-boxes-grid {
  display: grid;
  gap: 4px;
}

.grid-5 {
  grid-template-columns: repeat(5, 1fr);
}

.grid-6 {
  grid-template-columns: repeat(6, 1fr);
}

.grid-7 {
  grid-template-columns: repeat(7, 1fr);
}

.career-checkbox-label {
  cursor: pointer;
  user-select: none;
  display: flex;
  justify-content: center;
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
  width: 22px;
  height: 22px;
  border: 1.5px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  background: rgba(var(--v-theme-surface), 0.7);
  transition: all 0.15s ease;
}

.career-checkbox-label.checked .career-checkbox-box {
  background: rgb(var(--v-theme-primary));
  border-color: rgb(var(--v-theme-primary));
}

.career-checkbox-label:hover .career-checkbox-box {
  border-color: rgb(var(--v-theme-primary));
}
</style>
