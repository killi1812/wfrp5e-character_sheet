<script setup lang="ts">
import AddBtn from './AddBtn.vue'

defineProps<{
  title: string
  addLabel?: string
  addSize?: string
  elevation?: number
  color?: string
  cardVariant?: 'elevated' | 'flat' | 'tonal' | 'outlined' | 'text' | 'plain'
  fullHeight?: boolean
  titleClass?: string
  mb?: string
}>()

defineEmits<{
  (e: 'add'): void
}>()
</script>

<template>
  <v-card
    :color="color || 'surface'"
    :elevation="cardVariant === 'outlined' ? 0 : (elevation ?? 2)"
    :variant="cardVariant"
    class="pa-4 rounded-lg border"
    :class="{ 'h-100': fullHeight, [mb || 'mb-0']: true }"
  >
    <div
      v-if="title"
      class="d-flex justify-space-between align-center"
      :class="addLabel ? 'mb-2' : 'mb-3'"
    >
      <div :class="[titleClass || 'text-subtitle-2', 'font-weight-black text-uppercase text-primary']">
        {{ title }}
      </div>
      <div class="d-flex align-center gap-1">
        <slot name="actions">
          <AddBtn
            v-if="addLabel"
            :label="addLabel"
            :size="addSize || 'small'"
            @click="$emit('add')"
          />
        </slot>
      </div>
    </div>
    <slot />
  </v-card>
</template>
