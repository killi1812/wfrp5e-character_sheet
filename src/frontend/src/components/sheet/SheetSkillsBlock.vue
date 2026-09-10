<script setup lang="ts">
import type { Skill } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

const props = defineProps<{
  basicSkills: Skill[]
  advancedSkills: Skill[]
  getCharCurrent: (code: string) => number
}>()

const emit = defineEmits<{
  (e: 'addAdvancedSkill'): void
  (e: 'removeAdvancedSkill', index: number): void
}>()

const getSkillTotal = (skill: { characteristic: string; adv: number }) => {
  const base = props.getCharCurrent(skill.characteristic)
  return base + (Number(skill.adv) || 0)
}
</script>

<template>
  <v-row dense class="mb-4">
    <!-- BASIC SKILLS -->
    <v-col cols="12" md="6">
      <SectionCard title="Basic Skills" full-height>
        <v-table density="compact" class="bg-transparent text-caption skill-table">
          <thead>
            <tr>
              <th class="text-left font-weight-bold">Skill</th>
              <th class="text-center font-weight-bold" style="width: 70px;">Char</th>
              <th class="text-center font-weight-bold" style="width: 70px;">Adv</th>
              <th class="text-right font-weight-bold" style="width: 70px;">Total</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="skill in basicSkills" :key="skill.name" :class="{ 'skill-important-row': skill.important }">
              <td>
                <div class="d-flex align-center justify-space-between position-relative skill-name-cell">
                  <span class="font-weight-medium text-high-emphasis text-wrap mr-1">{{ skill.name }}</span>
                  <v-btn
                    icon
                    size="x-small"
                    variant="text"
                    density="compact"
                    class="skill-important-btn"
                    :class="{ 'btn-active': skill.important }"
                    :title="skill.important ? 'Marked as important' : 'Mark as important'"
                    @click="skill.important = !skill.important"
                  >
                    <v-icon size="14">{{ skill.important ? 'mdi-alert-circle' : 'mdi-alert-circle-outline' }}</v-icon>
                  </v-btn>
                </div>
              </td>
              <td class="text-center text-primary font-weight-bold">{{ skill.characteristic }}</td>
              <td class="text-center">
                <input v-model.number="skill.adv" type="number" class="skill-num-input font-weight-medium" placeholder="0" />
              </td>
              <td class="text-right font-weight-black">
                <v-tooltip text="Skill Total = Characteristic + Advances" location="right">
                  <template #activator="{ props: tProps }">
                    <span v-bind="tProps" class="skill-total-badge text-high-emphasis">
                      {{ getSkillTotal(skill) }}
                    </span>
                  </template>
                </v-tooltip>
              </td>
            </tr>
          </tbody>
        </v-table>
      </SectionCard>
    </v-col>

    <!-- ADVANCED SKILLS -->
    <v-col cols="12" md="6">
      <SectionCard title="Advanced Skills" add-label="Add Skill" full-height @add="emit('addAdvancedSkill')">
        <v-table density="compact" class="bg-transparent text-caption skill-table">
          <thead>
            <tr>
              <th class="text-left font-weight-bold">Skill</th>
              <th class="text-center font-weight-bold" style="width: 70px;">Char</th>
              <th class="text-center font-weight-bold" style="width: 70px;">Adv</th>
              <th class="text-right font-weight-bold" style="width: 70px;">Total</th>
              <th style="width: 40px;"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(skill, index) in advancedSkills" :key="index" :class="{ 'skill-important-row': skill.important }">
              <td>
                <div class="d-flex align-center justify-space-between position-relative skill-name-cell">
                  <input v-model="skill.name" type="text" class="skill-text-input font-weight-medium" placeholder="Skill Name" />
                  <v-btn
                    icon
                    size="x-small"
                    variant="text"
                    density="compact"
                    class="skill-important-btn ml-1"
                    :class="{ 'btn-active': skill.important }"
                    :title="skill.important ? 'Marked as important' : 'Mark as important'"
                    @click="skill.important = !skill.important"
                  >
                    <v-icon size="14">{{ skill.important ? 'mdi-alert-circle' : 'mdi-alert-circle-outline' }}</v-icon>
                  </v-btn>
                </div>
              </td>
              <td class="text-center">
                <input v-model="skill.characteristic" type="text" class="skill-char-input font-weight-bold text-primary text-uppercase" placeholder="Int" />
              </td>
              <td class="text-center">
                <input v-model.number="skill.adv" type="number" class="skill-num-input font-weight-medium" placeholder="0" />
              </td>
              <td class="text-right font-weight-black">
                <v-tooltip text="Skill Total = Characteristic + Advances" location="right">
                  <template #activator="{ props: tProps }">
                    <span v-bind="tProps" class="skill-total-badge text-high-emphasis">
                      {{ getSkillTotal(skill) }}
                    </span>
                  </template>
                </v-tooltip>
              </td>
              <td class="text-right">
                <DeleteRowBtn @delete="emit('removeAdvancedSkill', index)" />
              </td>
            </tr>
          </tbody>
        </v-table>
      </SectionCard>
    </v-col>
  </v-row>
</template>

<style scoped>
.skill-num-input {
  width: 50px;
  text-align: center;
  background: rgba(var(--v-theme-surface), 0.6);
  color: currentColor;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 2px 4px;
  font-size: 0.85rem;
  outline: none;
}

.skill-char-input {
  width: 50px;
  text-align: center;
  background: rgba(var(--v-theme-surface), 0.6);
  color: currentColor;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 2px 4px;
  font-size: 0.85rem;
  outline: none;
}

.skill-text-input {
  width: 100%;
  background: rgba(var(--v-theme-surface), 0.6);
  color: currentColor;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 0.85rem;
  outline: none;
}

.skill-num-input:focus,
.skill-char-input:focus,
.skill-text-input:focus {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-surface));
}

.skill-total-badge {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(var(--v-theme-surface), 0.8);
  color: rgb(var(--v-theme-on-surface)) !important;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  font-weight: 700;
}

.skill-table :deep(td) {
  padding: 4px 8px !important;
  height: auto !important;
  white-space: normal;
}

.skill-important-row {
  background-color: rgba(var(--v-theme-primary), 0.12) !important;
}

.skill-name-cell .skill-important-btn {
  opacity: 0;
  transition: opacity 0.2s ease-in-out;
  color: rgba(var(--v-theme-on-surface), 0.6);
}

.skill-name-cell:hover .skill-important-btn,
.skill-important-btn.btn-active {
  opacity: 1;
}

.skill-important-btn.btn-active {
  color: rgb(var(--v-theme-primary)) !important;
}
</style>
