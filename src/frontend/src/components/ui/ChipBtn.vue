<script setup lang="ts">
withDefaults(
  defineProps<{
    label?: string | number
    color?: string
    variant?: 'tonal' | 'elevated' | 'flat' | 'outlined' | 'text' | 'plain'
    size?: string
    icon?: string
    disabled?: boolean
  }>(),
  {
    label: '',
    color: 'primary',
    variant: 'tonal',
    size: 'default',
    icon: '',
    disabled: false,
  }
)

defineEmits<{
  (e: 'click', event: MouseEvent | KeyboardEvent): void
}>()
</script>

<template>
  <v-chip
    :color="color"
    :variant="variant"
    :size="size"
    :disabled="disabled"
    class="chip-btn font-weight-bold cursor-pointer"
    @click="$emit('click', $event)"
  >
    <v-icon v-if="icon" start :icon="icon" />
    <slot>{{ label }}</slot>
  </v-chip>
</template>

<style scoped>
.chip-btn {
  cursor: pointer;
  user-select: none;
  min-width: 48px;
  justify-content: center;
}
</style>
