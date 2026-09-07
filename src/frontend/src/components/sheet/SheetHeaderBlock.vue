<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <v-card color="surface" elevation="2" class="mb-4 pa-4 rounded-lg border">
    <v-row dense align="center">
      <v-col cols="12" md="3">
        <v-tooltip text="Character Name" location="top">
          <template #activator="{ props: tProps }">
            <v-text-field v-bind="tProps" v-model="character.name" label="Name" variant="outlined" density="compact" hide-details class="font-weight-bold" />
          </template>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.species" label="Species" variant="outlined" density="compact" hide-details />
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.appearance" label="Appearance" variant="outlined" density="compact" hide-details />
      </v-col>
      <v-col cols="12" sm="6" md="2">
        <v-text-field v-model="character.class" label="Class" variant="outlined" density="compact" hide-details />
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-text-field v-model="character.career" label="Career" variant="outlined" density="compact" hide-details />
      </v-col>

      <v-col cols="6" sm="4" md="2">
        <v-text-field v-model.number="character.careerLevel" label="Career Level" type="number" variant="outlined" density="compact" hide-details />
      </v-col>
      <v-col cols="6" sm="4" md="2">
        <v-tooltip text="Current Career Tier Level (e.g. 1, 2, 3, 4)" location="top">
          <template #activator="{ props: tProps }">
            <v-text-field
              v-bind="tProps"
              v-model.number="character.careerTier"
              label="Career Tier"
              type="number"
              variant="outlined"
              density="compact"
              hide-details
            />
          </template>
        </v-tooltip>
      </v-col>
      <v-col cols="12" sm="4" md="4">
        <v-text-field v-model="character.careerPath" label="Career Path" variant="outlined" density="compact" hide-details />
      </v-col>
      <v-col cols="6" sm="6" md="2">
        <v-text-field v-model="character.status" label="Status" variant="outlined" density="compact" hide-details />
      </v-col>

      <!-- Experience (XP) -->
      <v-col cols="12" md="2">
        <div class="d-flex gap-1 align-center">
          <v-tooltip text="Experience points currently available to spend" location="top">
            <template #activator="{ props: tProps }">
              <v-text-field v-bind="tProps" v-model.number="character.xp.current" label="Current XP" type="number" variant="outlined" density="compact" hide-details />
            </template>
          </v-tooltip>
          <v-tooltip text="Total experience points spent on advances" location="top">
            <template #activator="{ props: tProps }">
              <v-text-field v-bind="tProps" v-model.number="character.xp.spent" label="Spent XP" type="number" variant="outlined" density="compact" hide-details />
            </template>
          </v-tooltip>
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
