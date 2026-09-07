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
      <v-col cols="12" md="3">
        <TooltipField v-model="character.name" label="Name" tooltip="Character Name" field-class="font-weight-bold" />
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.species" label="Species" variant="outlined" density="compact" hide-details placeholder="Human, Elf, Dwarf..." />
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.appearance" label="Appearance" variant="outlined" density="compact" hide-details placeholder="Height, build, features..." />
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.class" label="Class" variant="outlined" density="compact" hide-details placeholder="Academic, Warrior..." />
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-text-field v-model="character.career" label="Career" variant="outlined" density="compact" hide-details placeholder="Wizard, Soldier..." />
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-text-field v-model="character.status" label="Status" variant="outlined" density="compact" hide-details placeholder="Silver 3" />
      </v-col>

      <!-- Experience (XP) -->
      <v-col cols="12" sm="6" md="4">
        <div class="d-flex gap-1 align-center">
          <TooltipField v-model.number="character.xp.current" label="Current XP" type="number" tooltip="Experience points currently available to spend" />
          <TooltipField v-model.number="character.xp.spent" label="Spent XP" type="number" tooltip="Total experience points spent on advances" />
          <v-tooltip text="Total XP earned over career (Current + Spent)" location="top">
            <template #activator="{ props: tProps }">
              <v-chip v-bind="tProps" color="primary" variant="tonal" class="font-weight-bold ml-1">
                {{ (character.xp.current || 0) + (character.xp.spent || 0) }} Total
              </v-chip>
            </template>
          </v-tooltip>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>

<style scoped>
.gap-1 { gap: 4px; }
</style>
