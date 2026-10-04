<script setup lang="ts">
import type { Talent, Skill } from '../../../constants/placeholders'
import TalentsCard from '../cards/TalentsCard.vue'
import LanguagesCard from '../cards/LanguagesCard.vue'

defineProps<{
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
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Talents -->
    <v-col cols="12" md="8">
      <TalentsCard
        :talents="talents"
        @add-talent="emit('addTalent')"
        @remove-talent="(idx) => emit('removeTalent', idx)"
      />
    </v-col>

    <!-- Languages (Skills with Int) -->
    <v-col cols="12" md="4">
      <LanguagesCard
        :languages="languages"
        :get-char-current="getCharCurrent"
        @add-language="emit('addLanguage')"
        @remove-language="(idx) => emit('removeLanguage', idx)"
      />
    </v-col>
  </v-row>
</template>
