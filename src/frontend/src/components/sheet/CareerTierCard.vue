<script setup lang="ts">
const props = defineProps<{
  tier: number
  total: number
  cols: number
  advances: boolean[]
}>()

function toggle(index: number) {
  props.advances[index] = !props.advances[index]
}
</script>

<template>
  <div class="tier-card pa-2 rounded border bg-surface h-100">
    <div class="d-flex justify-space-between align-center mb-2">
      <span class="text-caption font-weight-bold text-primary">Tier {{ tier }} Advances</span>
      <v-chip size="x-small" variant="tonal" color="primary" class="font-weight-bold">
        {{ advances.filter(Boolean).length }}/{{ total }}
      </v-chip>
    </div>
    <div class="tier-boxes-grid" :class="`grid-${cols}`">
      <v-tooltip
        v-for="i in total"
        :key="`adv${tier}-${i}`"
        :text="`Tier ${tier} Advance ${i}`"
        location="top"
        :open-on-focus="false"
      >
        <template #activator="{ props: tProps }">
          <label v-bind="tProps" class="career-checkbox-label" :class="{ checked: advances[i - 1] }">
            <input
              type="checkbox"
              :checked="advances[i - 1]"
              class="career-checkbox-input"
              @change="toggle(i - 1)"
            />
            <span class="career-checkbox-box"></span>
          </label>
        </template>
      </v-tooltip>
    </div>
  </div>
</template>

<style scoped>
.tier-card {
  transition: border-color 0.2s;
}
.tier-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.35);
}

.tier-boxes-grid {
  display: grid;
  gap: 3px;
  justify-items: center;
}
.grid-5 { grid-template-columns: repeat(5, 1fr); }
.grid-6 { grid-template-columns: repeat(6, 1fr); }
.grid-7 { grid-template-columns: repeat(7, 1fr); }

.career-checkbox-label {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 1px;
}

.career-checkbox-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.career-checkbox-box {
  display: block;
  width: 14px;
  height: 14px;
  border: 1.5px solid rgba(var(--v-theme-on-surface), 0.35);
  border-radius: 2px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  transition: all 0.15s ease;
}

.career-checkbox-label:hover .career-checkbox-box {
  border-color: rgb(var(--v-theme-primary));
  background: rgba(var(--v-theme-primary), 0.1);
}

.career-checkbox-label.checked .career-checkbox-box {
  background: rgb(var(--v-theme-primary));
  border-color: rgb(var(--v-theme-primary));
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.2);
}
</style>
