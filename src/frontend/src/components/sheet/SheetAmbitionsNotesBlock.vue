<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <v-row dense>
    <v-col cols="12" md="6">
      <SectionCard title="Ambitions" full-height>
        <v-textarea
          v-model="character.personalAmbition"
          label="Personal Ambition"
          rows="2"
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
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Party goals..."
        />
      </SectionCard>
    </v-col>
    <v-col cols="12" md="6">
      <SectionCard title="Notes" full-height>
        <v-textarea
          v-model="character.notes"
          label="Adventure & Campaign Notes"
          rows="5"
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Background, allies, enemies, quest notes..."
        />
      </SectionCard>
    </v-col>
  </v-row>
</template>
