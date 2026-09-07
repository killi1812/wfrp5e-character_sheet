<script setup lang="ts">
import type { Talent, Language } from '../../constants/placeholders'

defineProps<{
  talents: Talent[]
  languages: Language[]
}>()

const emit = defineEmits<{
  (e: 'addTalent'): void
  (e: 'removeTalent', index: number): void
  (e: 'addLanguage'): void
  (e: 'removeLanguage', index: number): void
}>()
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Talents -->
    <v-col cols="12" md="8">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Talents</div>
          <v-btn color="primary" size="x-small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addTalent')">
            Add Talent
          </v-btn>
        </div>
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Talent Name</th>
              <th class="text-center">Times Taken</th>
              <th>Description / Effect</th>
              <th class="text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(talent, idx) in talents" :key="idx">
              <td><v-text-field v-model="talent.name" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 70px;"><v-text-field v-model.number="talent.rank" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td><v-text-field v-model="talent.desc" variant="plain" density="compact" hide-details /></td>
              <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeTalent', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>

    <!-- Languages -->
    <v-col cols="12" md="4">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Languages</div>
          <v-btn color="primary" size="x-small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addLanguage')">
            Add Language
          </v-btn>
        </div>
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Language</th>
              <th class="text-center">Adv</th>
              <th class="text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(lang, idx) in languages" :key="idx">
              <td><v-text-field v-model="lang.name" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 70px;"><v-text-field v-model.number="lang.adv" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeLanguage', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>
  </v-row>
</template>
