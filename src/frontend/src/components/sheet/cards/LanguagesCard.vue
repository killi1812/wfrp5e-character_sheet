<script setup lang="ts">
import type { Skill } from '../../../constants/placeholders'
import SectionCard from '../../ui/SectionCard.vue'
import DeleteRowBtn from '../../ui/DeleteRowBtn.vue'
import ImportantBtn from '../../ui/ImportantBtn.vue'

const props = defineProps<{
  languages: Skill[]
  getCharCurrent: (code: string) => number
}>()

const emit = defineEmits<{
  (e: 'addLanguage'): void
  (e: 'removeLanguage', index: number): void
}>()

const getLanguageTotal = (lang: Skill) => {
  const base = props.getCharCurrent('Int')
  return base + (Number(lang.adv) || 0)
}
</script>

<template>
  <SectionCard title="Languages" add-label="Add Language" full-height @add="emit('addLanguage')">
    <v-table density="compact" class="bg-transparent text-caption">
      <thead>
        <tr>
          <th>Language</th>
          <th class="text-center" style="width: 70px;">Adv</th>
          <th class="text-center" style="width: 60px;">Total</th>
          <th class="text-right" style="width: 40px;">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(lang, idx) in languages" :key="idx" :class="{ 'important-row': lang.important }">
          <td>
            <div class="d-flex align-center important-trigger-cell">
              <ImportantBtn
                :active="lang.important"
                class="mr-1 flex-shrink-0"
                @toggle="lang.important = !lang.important"
              />
              <v-text-field v-model="lang.name" variant="plain" density="compact" hide-details placeholder="Language Name" />
            </div>
          </td>
          <td class="text-center">
            <input
              v-model.number="lang.adv"
              type="number"
              class="skill-num-input font-weight-medium"
              placeholder="0"
            />
          </td>
          <td class="text-center font-weight-black">
            <v-tooltip text="Language Total = Int + Advances" location="right">
              <template #activator="{ props: tProps }">
                <span v-bind="tProps" class="text-high-emphasis">
                  {{ getLanguageTotal(lang) }}
                </span>
              </template>
            </v-tooltip>
          </td>
          <td class="text-right"><DeleteRowBtn @delete="emit('removeLanguage', idx)" /></td>
        </tr>
      </tbody>
    </v-table>
  </SectionCard>
</template>

<style scoped>
.skill-num-input {
  width: 50px;
  text-align: center;
  background: rgba(var(--v-theme-surface), 0.6);
  color: currentColor;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 2px 4px;
  font-size: 0.85rem;
  outline: none;
}

.skill-num-input:focus {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-surface));
}

:deep(td) {
  white-space: normal;
}
</style>
