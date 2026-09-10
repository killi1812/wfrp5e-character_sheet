<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

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
        <v-card
          color="surface-variant"
          variant="outlined"
          class="pa-2 rounded text-center char-card position-relative"
          :class="{ 'char-important': isImportant(code as string) }"
        >
          <!-- Characteristic Header with Hover Important Button -->
          <div class="char-header position-relative border-bottom pb-1 mb-2 d-flex align-center justify-center">
            <v-tooltip :text="`${stat.name}${stat.hint ? ' — ' + stat.hint : ''}`" location="top" :open-on-focus="false">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="text-subtitle-2 font-weight-black text-primary cursor-pointer px-1">
                  {{ code }}
                </div>
              </template>
            </v-tooltip>

            <!-- Important '!' hover toggle button -->
            <v-tooltip :text="isImportant(code as string) ? 'Unmark as Important' : 'Mark as Important'" location="top" :open-on-focus="false">
              <template #activator="{ props: tProps }">
                <button
                  v-bind="tProps"
                  type="button"
                  class="important-toggle-btn"
                  :class="{ active: isImportant(code as string) }"
                  @click.stop="toggleImportant(code as string)"
                >
                  !
                </button>
              </template>
            </v-tooltip>
          </div>

          <!-- Initial Field -->
          <div class="d-flex align-center justify-space-between mb-1 px-1">
            <span class="text-caption text-high-emphasis font-weight-medium" style="font-size: 0.75rem;">Init</span>
            <input
              v-model.number="stat.initial"
              type="number"
              class="char-num-input"
              placeholder="0"
            />
          </div>

          <!-- Advances Field -->
          <div class="d-flex align-center justify-space-between mb-2 px-1">
            <span class="text-caption text-high-emphasis font-weight-medium" style="font-size: 0.75rem;">Adv</span>
            <input
              v-model.number="stat.advances"
              type="number"
              class="char-num-input adv-input"
              placeholder="0"
            />
          </div>

          <!-- Current Total Badge (with Encumbrance Agility Penalty if applicable) -->
          <v-tooltip
            :text="code === 'Ag' && agilityPenalty ? `Base: ${getCharCurrent(code as string)} (Penalized -${agilityPenalty} by Encumbrance)` : 'Current Score = Initial + Advances'"
            location="bottom"
            :open-on-focus="false"
          >
            <template #activator="{ props: tProps }">
              <div
                v-bind="tProps"
                class="stat-total-badge text-on-primary font-weight-black rounded py-1"
                :class="code === 'Ag' && agilityPenalty ? 'bg-error' : 'bg-primary'"
              >
                {{ getEffectiveScore(code as string) }}
                <span v-if="code === 'Ag' && agilityPenalty" class="text-caption font-weight-bold ml-1">
                  (-{{ agilityPenalty }})
                </span>
              </div>
            </template>
          </v-tooltip>
        </v-card>
      </v-col>
    </v-row>
  </v-card>
</template>

<style scoped>
.char-card {
  transition: border-color 0.2s;
}
.char-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.5) !important;
}

.char-num-input {
  width: 44px;
  text-align: center;
  background: rgba(var(--v-theme-surface), 0.9);
  color: rgb(var(--v-theme-on-surface)) !important;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 2px 4px;
  font-size: 0.9rem;
  font-weight: 700;
  outline: none;
  transition: all 0.15s ease;
}

.char-num-input::placeholder {
  color: rgba(var(--v-theme-on-surface), 0.4);
}

.char-num-input:focus {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 0 0 1.5px rgba(var(--v-theme-primary), 0.3);
}

.char-num-input.adv-input {
  color: rgb(var(--v-theme-primary)) !important;
}

.stat-total-badge {
  font-size: 1rem;
  letter-spacing: 0.5px;
  line-height: 1.2;
}

.cursor-pointer {
  cursor: pointer;
}

.important-toggle-btn {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background: rgba(var(--v-theme-surface), 0.9);
  color: rgba(var(--v-theme-on-surface), 0.6);
  font-size: 0.7rem;
  font-weight: 900;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: all 0.15s ease;
}

.char-card:hover .important-toggle-btn,
.important-toggle-btn.active {
  opacity: 1;
}

.important-toggle-btn.active {
  background: rgb(var(--v-theme-warning));
  color: rgb(var(--v-theme-on-warning));
  border-color: rgb(var(--v-theme-warning));
}

.char-card.char-important {
  border-color: rgb(var(--v-theme-warning)) !important;
  box-shadow: 0 0 8px rgba(var(--v-theme-warning), 0.35);
  background: rgba(var(--v-theme-warning), 0.05) !important;
}
</style>
