<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <v-row dense>
    <v-col cols="12" md="6">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-2">Ambitions</div>
        <v-textarea v-model="character.personalAmbition" label="Short-Term Personal Ambition" rows="2"
          variant="outlined" density="comfortable" class="mb-2" hide-details />
        <v-textarea v-model="character.partyAmbition" label="Long-Term Party Ambition" rows="2" variant="outlined"
          density="comfortable" hide-details />
      </v-card>
    </v-col>
    <v-col cols="12" md="6">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="text-subtitle-2 font-weight-black text-uppercase text-primary mb-2">Character Background & Notes
        </div>
        <v-textarea v-model="character.notes" label="Background & Adventure Notes" rows="5" variant="outlined"
          density="comfortable" hide-details />
      </v-card>
    </v-col>
  </v-row>
</template>
