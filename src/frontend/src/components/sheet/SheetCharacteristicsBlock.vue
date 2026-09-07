<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

defineProps<{
  character: CharacterModel
  getCharCurrent: (code: string) => number
}>()
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
        <v-card color="surface-variant" variant="outlined" class="pa-2 rounded text-center char-card">
          <v-tooltip :text="`${stat.name}${stat.hint ? ' — ' + stat.hint : ''}`" location="top">
            <template #activator="{ props: tProps }">
              <div v-bind="tProps" class="text-subtitle-2 font-weight-black text-primary border-bottom pb-1 mb-2 cursor-pointer">
                {{ code }}
              </div>
            </template>
          </v-tooltip>

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

          <!-- Current Total Badge -->
          <v-tooltip text="Current Score = Initial + Advances" location="bottom">
            <template #activator="{ props: tProps }">
              <div v-bind="tProps" class="stat-total-badge bg-primary text-on-primary font-weight-black rounded py-1">
                {{ getCharCurrent(code as string) }}
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
</style>
