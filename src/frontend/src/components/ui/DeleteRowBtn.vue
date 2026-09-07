<script setup lang="ts">
defineProps<{
  size?: string
}>()

defineEmits<{
  (e: 'delete'): void
}>()
</script>

<template>
  <v-btn
    icon="mdi-delete"
    :size="size || 'x-small'"
    color="error"
    variant="text"
    @click="$emit('delete')"
  />
</template>
