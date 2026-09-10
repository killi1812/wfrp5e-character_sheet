<script setup lang="ts">
defineProps<{
  modelValue?: string | number
  label?: string
  tooltip?: string
  type?: 'text' | 'number'
  variant?: 'outlined' | 'filled' | 'plain'
  readonly?: boolean
  tooltipLocation?: 'top' | 'bottom' | 'left' | 'right'
  fieldClass?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string | number): void
}>()
</script>

<template>
  <!-- With tooltip -->
  <v-tooltip
    v-if="tooltip"
    :text="tooltip"
    :location="tooltipLocation || 'top'"
    :open-on-focus="false"
    :open-on-click="false"
    :open-delay="150"
    :close-delay="50"
  >
    <template #activator="{ props: tProps }">
      <v-text-field
        v-bind="tProps"
        :model-value="modelValue"
        :label="label"
        :type="type || 'text'"
        :variant="variant || 'outlined'"
        :readonly="readonly || false"
        density="compact"
        hide-details
        :class="fieldClass"
        @update:model-value="$emit('update:modelValue', $event)"
      />
    </template>
  </v-tooltip>

  <!-- Without tooltip -->
  <v-text-field
    v-else
    :model-value="modelValue"
    :label="label"
    :type="type || 'text'"
    :variant="variant || 'outlined'"
    :readonly="readonly || false"
    density="compact"
    hide-details
    :class="fieldClass"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>
