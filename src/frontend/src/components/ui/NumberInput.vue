<script setup lang="ts">
const model = defineModel<number>()

withDefaults(
  defineProps<{
    placeholder?: string | number
    min?: number
    max?: number
  }>(),
  {
    placeholder: '0',
  }
)
</script>

<template>
  <input
    v-model.number="model"
    type="number"
    class="skill-num-input font-weight-medium"
    :placeholder="String(placeholder)"
    :min="min"
    :max="max"
  />
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

.skill-num-input:focus {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-surface));
}
</style>
