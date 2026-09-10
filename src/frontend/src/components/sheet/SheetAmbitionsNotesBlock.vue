<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <div class="d-flex flex-wrap gap-3 mb-4 ambitions-notes-row">
    <!-- Notes Card (70% on desktop) -->
    <div class="notes-col">
      <SectionCard title="Adventure & Campaign Notes" full-height>
        <v-textarea
          v-model="character.notes"
          label="Adventure & Campaign Notes"
          auto-grow
          :rows="12"
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Background, allies, enemies, quest notes..."
          class="notes-textarea"
        />
      </SectionCard>
    </div>

    <!-- Ambitions & Armour Points (30% on desktop) -->
    <div class="side-col d-flex flex-column gap-3">
      <!-- Ambitions Card -->
      <SectionCard title="Ambitions">
        <v-textarea
          v-model="character.personalAmbition"
          label="Personal Ambition"
          rows="2"
          auto-grow
          variant="outlined"
          density="comfortable"
          class="mb-2"
          hide-details
          placeholder="Personal goals..."
        />
        <v-textarea
          v-model="character.partyAmbition"
          label="Party Ambition"
          rows="2"
          auto-grow
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Party goals..."
        />
      </SectionCard>

      <!-- Armour Points (AP) Card (Relocated here from Vitals) -->
      <SectionCard title="Armour Points (AP)">
        <div class="ap-body-container py-1">
          <!-- 1. Head (Centered) -->
          <div class="d-flex justify-center mb-2">
            <v-tooltip text="Head protection (d100 roll 01-09)" location="top" :open-on-focus="false">
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
            <v-tooltip text="Left Arm protection (d100 roll 10-24)" location="top" :open-on-focus="false">
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

            <v-tooltip text="Right Arm protection (d100 roll 25-44)" location="top" :open-on-focus="false">
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
            <v-tooltip text="Body Torso protection (d100 roll 45-79)" location="top" :open-on-focus="false">
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
            <v-tooltip text="Left Leg protection (d100 roll 80-89)" location="top" :open-on-focus="false">
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

            <v-tooltip text="Right Leg protection (d100 roll 90-00)" location="top" :open-on-focus="false">
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
            <v-tooltip text="Shield AP Bonus" location="top" :open-on-focus="false">
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
    </div>
  </div>
</template>

<style scoped>
.gap-3 { gap: 12px; }

.ambitions-notes-row {
  display: flex;
}

.notes-col {
  flex: 0 0 calc(70% - 6px);
  max-width: calc(70% - 6px);
}

.side-col {
  flex: 0 0 calc(30% - 6px);
  max-width: calc(30% - 6px);
}

@media (max-width: 960px) {
  .notes-col, .side-col {
    flex: 0 0 100%;
    max-width: 100%;
  }
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
