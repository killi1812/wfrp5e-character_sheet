<script setup lang="ts">
import type { Skill } from '../../constants/placeholders'

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
  <v-card color="surface" elevation="2" class="mb-4 pa-4 rounded-lg border">
    <div class="d-flex justify-space-between align-center mb-3">
      <div class="text-subtitle-1 font-weight-black text-uppercase text-primary">Skills System</div>
      <v-btn color="primary" size="small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addAdvancedSkill')">
        Add Advanced Skill
      </v-btn>
    </div>

    <v-row dense>
      <!-- BASIC SKILLS (READ-ONLY LIST, NON-DELETABLE) -->
      <v-col cols="12" md="6">
        <v-card color="surface-variant" variant="outlined" class="pa-3 rounded-lg h-100 border">
          <div class="text-subtitle-2 font-weight-bold text-primary mb-2 border-bottom pb-1">
            Printed Basic Skills (Standard)
          </div>
          <v-table density="compact" class="bg-transparent text-caption skill-table">
            <thead>
              <tr>
                <th class="text-left font-weight-bold">Skill Name</th>
                <th class="text-center font-weight-bold">Char</th>
                <th class="text-center font-weight-bold" style="width: 70px;">Adv</th>
                <th class="text-right font-weight-bold" style="width: 70px;">Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="skill in basicSkills" :key="skill.name">
                <td class="font-weight-medium text-high-emphasis">{{ skill.name }}</td>
                <td class="text-center text-primary font-weight-bold">{{ skill.characteristic }}</td>
                <td class="text-center">
                  <input
                    v-model.number="skill.adv"
                    type="number"
                    class="skill-num-input font-weight-medium"
                  />
                </td>
                <td class="text-right font-weight-black">
                  <v-tooltip text="Skill Total = Characteristic Current + Advances" location="right">
                    <template #activator="{ props: tProps }">
                      <span v-bind="tProps" class="px-2 py-1 bg-surface rounded text-primary border font-weight-bold">
                        {{ getSkillTotal(skill) }}
                      </span>
                    </template>
                  </v-tooltip>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>

      <!-- ADVANCED SKILLS (ADD/REMOVEABLE LIST) -->
      <v-col cols="12" md="6">
        <v-card color="surface-variant" variant="outlined" class="pa-3 rounded-lg h-100 border">
          <div class="text-subtitle-2 font-weight-bold text-secondary mb-2 border-bottom pb-1">
            Advanced / Grouped Skills (Custom & Learned)
          </div>
          <v-table density="compact" class="bg-transparent text-caption skill-table">
            <thead>
              <tr>
                <th class="text-left font-weight-bold">Skill Name</th>
                <th class="text-center font-weight-bold" style="width: 70px;">Char</th>
                <th class="text-center font-weight-bold" style="width: 70px;">Adv</th>
                <th class="text-right font-weight-bold" style="width: 70px;">Total</th>
                <th class="text-right font-weight-bold" style="width: 40px;"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(skill, index) in advancedSkills" :key="index">
                <td>
                  <input
                    v-model="skill.name"
                    type="text"
                    class="skill-text-input font-weight-medium"
                    placeholder="Skill Name"
                  />
                </td>
                <td class="text-center">
                  <input
                    v-model="skill.characteristic"
                    type="text"
                    class="skill-char-input font-weight-bold text-secondary text-uppercase"
                    placeholder="Char"
                  />
                </td>
                <td class="text-center">
                  <input
                    v-model.number="skill.adv"
                    type="number"
                    class="skill-num-input font-weight-medium"
                  />
                </td>
                <td class="text-right font-weight-black">
                  <v-tooltip text="Skill Total = Characteristic Current + Advances" location="right">
                    <template #activator="{ props: tProps }">
                      <span v-bind="tProps" class="px-2 py-1 bg-surface rounded text-secondary border font-weight-bold">
                        {{ getSkillTotal(skill) }}
                      </span>
                    </template>
                  </v-tooltip>
                </td>
                <td class="text-right">
                  <v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeAdvancedSkill', index)" />
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
    </v-row>
  </v-card>
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

.skill-table :deep(td) {
  padding: 4px 8px !important;
  height: auto !important;
}
</style>
