<script setup lang="ts">
defineProps<{
  modelValue: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()
</script>

<template>
  <div class="sin-counter-widget d-inline-flex align-center bg-surface px-2 py-1 rounded border">
    <span class="text-caption font-weight-bold text-secondary mr-2">Sin:</span>
    <v-btn
      icon="mdi-minus"
      size="x-small"
      variant="text"
      density="compact"
      :disabled="modelValue <= 0"
      @click="emit('update:modelValue', Math.max(0, (modelValue || 0) - 1))"
    />
    <input
      type="number"
      :value="modelValue"
      class="sin-input font-weight-black mx-1"
      placeholder="0"
      @input="emit('update:modelValue', Number(($event.target as HTMLInputElement).value) || 0)"
    />
    <v-btn
      icon="mdi-plus"
      size="x-small"
      variant="text"
      density="compact"
      @click="emit('update:modelValue', (modelValue || 0) + 1)"
    />
  </div>
</template>

<style scoped>
.sin-input {
  width: 32px;
  text-align: center;
  border: none;
  background: transparent;
  outline: none;
  color: currentColor;
  font-size: 0.9rem;
}
</style>
