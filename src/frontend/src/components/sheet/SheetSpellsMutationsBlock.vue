<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'
import AddBtn from '../ui/AddBtn.vue'
import SinCounter from '../ui/SinCounter.vue'

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
    <!-- Spells & Prayers -->
    <v-col v-if="!character.spellsHidden" cols="12" md="8">
      <SectionCard title="Spells & Prayers" add-label="Add Spell / Prayer" full-height @add="emit('addSpell')">
        <!-- Sin Counter Widget row -->
        <div class="d-flex justify-end align-center mb-2">
          <SinCounter v-model="character.sin" />
        </div>

        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th style="width: 200px;">Name</th>
              <th class="text-center" style="width: 45px;">CN</th>
              <th style="width: 100px;">Range</th>
              <th style="width: 100px;">Target</th>
              <th style="width: 100px;">Duration</th>
              <th>Description / Effect</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(spell, idx) in character.spells" :key="idx">
              <td style="width: 200px; max-width: 200px;">
                <v-text-field v-model="spell.name" variant="plain" density="compact" hide-details placeholder="Spell / Prayer" />
              </td>
              <td style="width: 45px; max-width: 45px;">
                <v-text-field v-model.number="spell.cn" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
              </td>
              <td style="width: 100px; max-width: 100px;">
                <v-text-field v-model="spell.range" variant="plain" density="compact" hide-details placeholder="12 yd" />
              </td>
              <td style="width: 100px; max-width: 100px;">
                <v-text-field v-model="spell.target" variant="plain" density="compact" hide-details placeholder="1 Target" />
              </td>
              <td style="width: 100px; max-width: 100px;">
                <v-text-field v-model="spell.duration" variant="plain" density="compact" hide-details placeholder="Instant" />
              </td>
              <td>
                <v-textarea
                  v-model="spell.description"
                  variant="plain"
                  density="compact"
                  rows="1"
                  auto-grow
                  hide-details
                  placeholder="Effect and rules description..."
                />
              </td>
              <td class="text-right"><DeleteRowBtn @delete="emit('removeSpell', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </SectionCard>
    </v-col>

    <!-- Corruption & Mutations -->
    <v-col cols="12" :md="character.spellsHidden ? 6 : 4">
      <SectionCard title="Corruption & Mutations" full-height>
        <v-row dense class="mb-3">
          <v-col cols="6">
            <v-text-field v-model.number="character.corruption.current" label="Current Corruption" type="number" variant="outlined" density="compact" hide-details placeholder="0" />
          </v-col>
          <v-col cols="6">
            <v-text-field v-model.number="character.corruption.max" label="Max Threshold" type="number" variant="outlined" density="compact" hide-details placeholder="0" />
          </v-col>
        </v-row>

        <div class="d-flex justify-space-between align-center mb-2">
          <span class="text-caption font-weight-bold text-primary">Mutations:</span>
          <AddBtn label="Add Mutation" @click="emit('addMutation')" />
        </div>

        <div v-if="character.mutations.length === 0" class="text-caption text-medium-emphasis font-italic py-2 text-center">
          No mutations recorded. Click "+ Add Mutation" to add one.
        </div>

        <div v-for="(mut, idx) in character.mutations" :key="idx" class="mutation-card pa-2 rounded mb-2 bg-surface-variant border">
          <div class="d-flex justify-space-between align-center mb-1">
            <input
              v-model="mut.name"
              type="text"
              class="mutation-input font-weight-bold text-caption"
              placeholder="Mutation Name"
            />
            <DeleteRowBtn @delete="emit('removeMutation', idx)" />
          </div>
          <v-textarea
            v-model="mut.effect"
            variant="plain"
            density="compact"
            rows="1"
            auto-grow
            hide-details
            placeholder="Effect / Description"
            class="text-caption text-medium-emphasis"
          />
        </div>
      </SectionCard>
    </v-col>
  </v-row>
</template>

<style scoped>
.mutation-card {
  transition: border-color 0.2s;
}
.mutation-card:focus-within {
  border-color: rgb(var(--v-theme-primary)) !important;
}

.mutation-input {
  width: 100%;
  border: none;
  background: transparent;
  outline: none;
  color: currentColor;
}

:deep(td) {
  white-space: normal;
}
</style>
