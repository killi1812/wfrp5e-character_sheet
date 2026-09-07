<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

defineProps<{
  character: CharacterModel
}>()

const emit = defineEmits<{
  (e: 'addSpell'): void
  (e: 'removeSpell', index: number): void
  (e: 'addMutation'): void
  (e: 'removeMutation', index: number): void
}>()
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Spells / Prayers -->
    <v-col cols="12" md="8">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Spells & Incantations</div>
          <v-btn color="primary" size="x-small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addSpell')">
            Add Spell
          </v-btn>
        </div>
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Spell Name</th>
              <th class="text-center">CN</th>
              <th>Range</th>
              <th>Target</th>
              <th>Duration</th>
              <th class="text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(spell, idx) in character.spells" :key="idx">
              <td><v-text-field v-model="spell.name" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 50px;"><v-text-field v-model.number="spell.cn" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td style="max-width: 90px;"><v-text-field v-model="spell.range" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 90px;"><v-text-field v-model="spell.target" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 90px;"><v-text-field v-model="spell.duration" variant="plain" density="compact" hide-details /></td>
              <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeSpell', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>

    <!-- Corruption & Mutations -->
    <v-col cols="12" md="4">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-2">Corruption & Mutations</div>
        <v-row dense class="mb-3">
          <v-col cols="6">
            <v-text-field v-model.number="character.corruption.current" label="Corruption" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
          <v-col cols="6">
            <v-text-field v-model.number="character.sin" label="Sin Points" type="number" variant="outlined" density="compact" hide-details />
          </v-col>
        </v-row>

        <div class="d-flex justify-space-between align-center mb-1">
          <span class="text-caption font-weight-bold">Mutations:</span>
          <v-btn icon="mdi-plus" size="x-small" color="primary" variant="text" @click="emit('addMutation')" />
        </div>

        <v-list density="compact" class="bg-transparent pa-0">
          <v-list-item v-for="(mut, idx) in character.mutations" :key="idx" class="px-2 py-1 rounded mb-1 bg-surface-variant border">
            <div class="d-flex justify-space-between align-center">
              <div class="text-caption font-weight-bold">{{ mut.name }}</div>
              <v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeMutation', idx)" />
            </div>
            <div class="text-caption text-medium-emphasis">{{ mut.effect }}</div>
          </v-list-item>
        </v-list>
      </v-card>
    </v-col>
  </v-row>
</template>
