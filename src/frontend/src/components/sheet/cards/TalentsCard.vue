<script setup lang="ts">
import type { Talent } from '../../../constants/placeholders'
import SectionCard from '../../ui/SectionCard.vue'
import DeleteRowBtn from '../../ui/DeleteRowBtn.vue'
import ImportantBtn from '../../ui/ImportantBtn.vue'

defineProps<{
  talents: Talent[]
}>()

const emit = defineEmits<{
  (e: 'addTalent'): void
  (e: 'removeTalent', index: number): void
}>()
</script>

<template>
  <SectionCard title="Talents" add-label="Add Talent" full-height @add="emit('addTalent')">
    <v-table density="compact" class="bg-transparent text-caption">
      <thead>
        <tr>
          <th style="width: 220px;">Talent Name</th>
          <th>Description / Effect</th>
          <th class="text-right" style="width: 40px;">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(talent, idx) in talents" :key="idx" :class="{ 'important-row': talent.important }">
          <td style="width: 220px; max-width: 220px;">
            <div class="d-flex align-center important-trigger-cell">
              <ImportantBtn
                :active="talent.important"
                class="mr-1 flex-shrink-0"
                @toggle="talent.important = !talent.important"
              />
              <v-text-field v-model="talent.name" variant="plain" density="compact" hide-details placeholder="Talent Name" />
            </div>
          </td>
          <td>
            <v-textarea
              v-model="talent.desc"
              variant="plain"
              density="compact"
              rows="1"
              auto-grow
              hide-details
              placeholder="Effect description"
            />
          </td>
          <td class="text-right"><DeleteRowBtn @delete="emit('removeTalent', idx)" /></td>
        </tr>
      </tbody>
    </v-table>
  </SectionCard>
</template>

<style scoped>
:deep(td) {
  white-space: normal;
}
</style>
