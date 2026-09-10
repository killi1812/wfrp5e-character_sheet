<script setup lang="ts">
import type { CharacterModel, CareerEntry } from '../../constants/placeholders'
import { NEW_ITEM_TEMPLATES } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'
import CareerTierCard from './CareerTierCard.vue'

const props = defineProps<{
  character: CharacterModel
}>()

function addCareer() {
  if (!props.character.careers) {
    props.character.careers = []
  }
  const isFirst = props.character.careers.length === 0
  const newC = NEW_ITEM_TEMPLATES.career()
  newC.active = isFirst
  props.character.careers.push(newC)
  if (isFirst) {
    syncActiveCareer(newC)
  }
}

function removeCareer(index: number) {
  const wasActive = props.character.careers[index].active
  props.character.careers.splice(index, 1)
  if (props.character.careers.length > 0) {
    if (wasActive) {
      props.character.careers[0].active = true
      syncActiveCareer(props.character.careers[0])
    }
  } else {
    props.character.career = ''
    props.character.class = ''
    props.character.status = ''
  }
}

function setActiveCareer(index: number) {
  props.character.careers.forEach((c, idx) => {
    c.active = idx === index
  })
  syncActiveCareer(props.character.careers[index])
}

function syncActiveCareer(c: CareerEntry) {
  props.character.career = c.career
  props.character.class = c.class
  props.character.status = c.status
}
</script>

<template>
  <SectionCard title="Careers & Advances" full-height add-label="Add Career" @add="addCareer">
    <div v-if="character.careers.length === 0" class="text-caption text-medium-emphasis font-italic py-4 text-center">
      No careers added. Click "+ Add Career" to add one.
    </div>

    <div
      v-for="(c, idx) in character.careers"
      :key="idx"
      class="career-row-card mb-3 pa-3 rounded-lg border"
      :class="{ 'career-active': c.active }"
    >
      <!-- Career Header info: Active toggle, Class, Career, Status, Delete -->
      <div class="d-flex flex-wrap align-center gap-2 mb-2">
        <!-- Active / Important marker -->
        <v-tooltip :text="c.active ? 'Current Active & Important Career' : 'Click to set as Active Career'" location="top" :open-on-focus="false">
          <template #activator="{ props: tProps }">
            <button
              v-bind="tProps"
              type="button"
              class="career-active-btn d-flex align-center px-2 py-1 rounded"
              :class="{ active: c.active }"
              @click="setActiveCareer(idx)"
            >
              <v-icon :icon="c.active ? 'mdi-star' : 'mdi-star-outline'" size="small" class="mr-1" />
              <span class="text-caption font-weight-bold">{{ c.active ? 'Active' : 'Inactive' }}</span>
            </button>
          </template>
        </v-tooltip>

        <div style="width: 140px;">
          <v-text-field
            v-model="c.class"
            label="Class"
            variant="outlined"
            density="compact"
            hide-details
            placeholder="Academic..."
            @update:model-value="c.active && syncActiveCareer(c)"
          />
        </div>
        <div class="flex-grow-1" style="min-width: 150px;">
          <v-text-field
            v-model="c.career"
            label="Career"
            variant="outlined"
            density="compact"
            hide-details
            placeholder="Wizard..."
            @update:model-value="c.active && syncActiveCareer(c)"
          />
        </div>
        <div style="width: 150px;">
          <v-text-field
            v-model="c.status"
            label="Status"
            variant="outlined"
            density="compact"
            hide-details
            placeholder="Silver 3"
            @update:model-value="c.active && syncActiveCareer(c)"
          />
        </div>

        <DeleteRowBtn @delete="removeCareer(idx)" />
      </div>

      <!-- Tier Advances (Tier 2, Tier 3, Tier 4 using reusable CareerTierCard) -->
      <v-row dense class="pt-2 border-t">
        <v-col cols="12" md="4" class="tier-col">
          <CareerTierCard :tier="2" :total="10" :cols="5" :advances="c.advances2" />
        </v-col>
        <v-col cols="12" md="4" class="tier-col">
          <CareerTierCard :tier="3" :total="12" :cols="6" :advances="c.advances3" />
        </v-col>
        <v-col cols="12" md="4" class="tier-col">
          <CareerTierCard :tier="4" :total="14" :cols="7" :advances="c.advances4" />
        </v-col>
      </v-row>
    </div>
  </SectionCard>
</template>

<style scoped>
.gap-2 { gap: 8px; }

.career-row-card {
  background: rgba(var(--v-theme-surface), 0.7);
  transition: all 0.2s ease;
}
.career-row-card.career-active {
  border-color: rgb(var(--v-theme-primary)) !important;
  box-shadow: 0 0 10px rgba(var(--v-theme-primary), 0.2);
}

.career-active-btn {
  background: transparent;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  color: rgba(var(--v-theme-on-surface), 0.6);
  cursor: pointer;
  transition: all 0.2s ease;
}
.career-active-btn:hover {
  border-color: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-primary));
}
.career-active-btn.active {
  background: rgba(var(--v-theme-primary), 0.15);
  border-color: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-primary));
}
</style>
