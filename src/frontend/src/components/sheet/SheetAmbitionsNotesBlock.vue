<script setup lang="ts">
import type { CharacterModel } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import ArmourPointsDiagram from './ArmourPointsDiagram.vue'

defineProps<{
  character: CharacterModel
}>()
</script>

<template>
  <div class="d-flex flex-wrap gap-3 mb-4 ambitions-notes-row">
    <!-- Notes Card (70% on desktop) -->
    <div class="notes-col">
      <SectionCard title="Adventure & Campaign Notes" full-height>
        <v-textarea
          v-model="character.notes"
          label="Adventure & Campaign Notes"
          auto-grow
          :rows="12"
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Background, allies, enemies, quest notes..."
          class="notes-textarea"
        />
      </SectionCard>
    </div>

    <!-- Ambitions & Armour Points (30% on desktop) -->
    <div class="side-col d-flex flex-column gap-3">
      <!-- Ambitions Card -->
      <SectionCard title="Ambitions">
        <v-textarea
          v-model="character.personalAmbition"
          label="Personal Ambition"
          rows="2"
          auto-grow
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
          auto-grow
          variant="outlined"
          density="comfortable"
          hide-details
          placeholder="Party goals..."
        />
      </SectionCard>

      <!-- Armour Points (AP) Diagram -->
      <ArmourPointsDiagram :armour-points="character.armourPoints" />
    </div>
  </div>
</template>

<style scoped>
.gap-3 { gap: 12px; }

.ambitions-notes-row {
  display: flex;
}

.notes-col {
  flex: 0 0 calc(70% - 6px);
  max-width: calc(70% - 6px);
}

.side-col {
  flex: 0 0 calc(30% - 6px);
  max-width: calc(30% - 6px);
}

@media (max-width: 960px) {
  .notes-col, .side-col {
    flex: 0 0 100%;
    max-width: 100%;
  }
}
</style>
