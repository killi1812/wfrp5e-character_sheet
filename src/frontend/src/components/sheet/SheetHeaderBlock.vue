<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import TooltipField from '../ui/TooltipField.vue'

defineProps<{
  character: CharacterModel
  isDirty?: boolean
}>()

defineEmits<{
  (e: 'openMenu'): void
}>()
</script>

<template>
  <div class="d-flex flex-wrap gap-3 mb-4 header-actions-row">
    <!-- Header Block (80% on desktop) -->
    <v-card color="surface" elevation="2" class="pa-4 rounded-lg border flex-grow-1 header-card">
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

    <!-- Actions Block (20% on desktop) -->
    <v-card color="surface" elevation="2" class="pa-3 rounded-lg border d-flex flex-column justify-center align-center actions-card">
      <div class="d-flex align-center justify-space-between w-100 mb-2">
        <v-chip
          :color="isDirty ? 'warning' : 'success'"
          size="small"
          variant="tonal"
          class="font-weight-bold"
        >
          <v-icon start :icon="isDirty ? 'mdi-alert-circle-outline' : 'mdi-check-circle-outline'" size="small" />
          {{ isDirty ? 'Not saved' : 'Saved' }}
        </v-chip>

        <!-- Kebab Menu Trigger -->
        <v-tooltip text="Sheet Menu & Settings" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <v-btn
              v-bind="tProps"
              icon="mdi-dots-vertical"
              size="small"
              color="primary"
              variant="text"
              @click="$emit('openMenu')"
            />
          </template>
        </v-tooltip>
      </div>

      <!-- Quick Toggles (Spells & Mount) -->
      <div class="d-flex align-center justify-space-around w-100 pt-1 border-t">
        <v-tooltip :text="character.spellsHidden ? 'Show Spells & Prayers' : 'Hide Spells & Prayers'" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <v-btn
              v-bind="tProps"
              :icon="character.spellsHidden ? 'mdi-auto-fix' : 'mdi-auto-fix'"
              size="small"
              :color="character.spellsHidden ? 'default' : 'primary'"
              :variant="character.spellsHidden ? 'outlined' : 'tonal'"
              @click="character.spellsHidden = !character.spellsHidden"
            />
          </template>
        </v-tooltip>

        <v-tooltip :text="character.mountHidden ? 'Show Mount Section' : 'Hide Mount Section'" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <v-btn
              v-bind="tProps"
              icon="mdi-horse"
              size="small"
              :color="character.mountHidden ? 'default' : 'primary'"
              :variant="character.mountHidden ? 'outlined' : 'tonal'"
              @click="character.mountHidden = !character.mountHidden"
            />
          </template>
        </v-tooltip>
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.header-actions-row {
  display: flex;
  gap: 12px;
}

.header-card {
  flex: 0 0 calc(80% - 6px);
  max-width: calc(80% - 6px);
}

.actions-card {
  flex: 0 0 calc(20% - 6px);
  max-width: calc(20% - 6px);
  min-width: 170px;
}

.border-t {
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

@media (max-width: 960px) {
  .header-card, .actions-card {
    flex: 0 0 100%;
    max-width: 100%;
  }
}
</style>
