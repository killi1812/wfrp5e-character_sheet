<script setup lang="ts">
import type { Talent, Skill } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

const props = defineProps<{
  talents: Talent[]
  languages: Skill[]
  getCharCurrent: (code: string) => number
}>()

const emit = defineEmits<{
  (e: 'addTalent'): void
  (e: 'removeTalent', index: number): void
  (e: 'addLanguage'): void
  (e: 'removeLanguage', index: number): void
}>()

const getLanguageTotal = (lang: Skill) => {
  const base = props.getCharCurrent('Int')
  return base + (Number(lang.adv) || 0)
}
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Talents -->
    <v-col cols="12" md="8">
      <SectionCard title="Talents" add-label="Add Talent" full-height @add="emit('addTalent')">
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th style="width: 200px;">Talent Name</th>
              <th>Description / Effect</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(talent, idx) in talents" :key="idx">
              <td style="width: 200px; max-width: 200px;"><v-text-field v-model="talent.name" variant="plain" density="compact" hide-details placeholder="Talent Name" /></td>
              <td><v-text-field v-model="talent.desc" variant="plain" density="compact" hide-details placeholder="Effect description" /></td>
              <td class="text-right"><DeleteRowBtn @delete="emit('removeTalent', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </SectionCard>
    </v-col>

    <!-- Languages (Skills with Int) -->
    <v-col cols="12" md="4">
      <SectionCard title="Languages" add-label="Add Language" full-height @add="emit('addLanguage')">
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Language</th>
              <th class="text-center" style="width: 60px;">Adv</th>
              <th class="text-right" style="width: 60px;">Total</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(lang, idx) in languages" :key="idx">
              <td><v-text-field v-model="lang.name" variant="plain" density="compact" hide-details placeholder="Language Name" /></td>
              <td style="max-width: 60px;"><v-text-field v-model.number="lang.adv" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
              <td class="text-right font-weight-black">
                <v-tooltip text="Language Total = Int + Advances" location="right">
                  <template #activator="{ props: tProps }">
                    <span v-bind="tProps" class="px-2 py-1 bg-surface rounded text-high-emphasis border font-weight-bold">
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
    </v-col>
  </v-row>
</template>
