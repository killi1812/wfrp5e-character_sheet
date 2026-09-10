<script setup lang="ts">
import { useSlots } from 'vue'

const props = withDefaults(
  defineProps<{
    label?: string
    size?: string
    color?: string
    variant?: 'tonal' | 'elevated' | 'flat' | 'outlined' | 'text' | 'plain'
    icon?: string
    density?: 'default' | 'comfortable' | 'compact'
  }>(),
  {
    label: '',
    size: 'small',
    color: 'primary',
    variant: 'tonal',
    icon: 'mdi-plus',
    density: 'default',
  }
)

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
  (e: 'add'): void
}>()

const slots = useSlots()
const hasLabel = () => Boolean(props.label || slots.default)

function handleClick(event: MouseEvent) {
  emit('click', event)
  emit('add')
}
</script>

<template>
  <v-btn
    v-if="hasLabel()"
    :color="color"
    :size="size"
    :variant="variant"
    :prepend-icon="icon"
    :density="density"
    class="font-weight-medium text-none"
    @click="handleClick"
  >
    <slot>{{ label }}</slot>
  </v-btn>
  <v-btn
    v-else
    :color="color"
    :size="size"
    :variant="variant"
    :icon="icon"
    :density="density"
    @click="handleClick"
  >
    <slot />
  </v-btn>
</template>
