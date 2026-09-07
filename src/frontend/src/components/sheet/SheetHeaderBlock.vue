<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import TooltipField from '../ui/TooltipField.vue'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <v-card color="surface" elevation="2" class="mb-4 pa-4 rounded-lg border">
    <v-row dense align="center">
      <v-col cols="12" md="4">
        <TooltipField v-model="character.name" label="Name" tooltip="Character Name" field-class="font-weight-bold" />
      </v-col>
      <v-col cols="12" sm="6" md="4">
        <v-text-field v-model="character.species" label="Species" variant="outlined" density="compact" hide-details placeholder="Human, Elf, Dwarf..." />
      </v-col>
      <v-col cols="12" sm="6" md="4">
        <v-text-field v-model="character.appearance" label="Appearance" variant="outlined" density="compact" hide-details placeholder="Height, build, features..." />
      </v-col>
    </v-row>
  </v-card>
</template>
