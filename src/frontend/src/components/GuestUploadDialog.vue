<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
  characterName: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm-upload'): void
}>()

function close() {
  emit('update:modelValue', false)
}

function handleConfirm() {
  emit('confirm-upload')
  close()
}
</script>

<template>
  <v-dialog :model-value="modelValue" max-width="500" transition="dialog-bottom-transition" @update:model-value="emit('update:modelValue', $event)">
    <v-card color="surface" class="pa-6 rounded-xl elevation-24 border">
      <div class="d-flex align-center mb-3">
        <v-icon icon="mdi-cloud-upload" color="primary" class="mr-2" size="large" />
        <h2 class="text-h6 font-weight-black mb-0 text-on-surface">Offline Character Found</h2>
      </div>

      <p class="text-body-2 text-medium-emphasis mb-4">
        You have an offline character saved in this browser
        <strong v-if="characterName">("{{ characterName }}")</strong>.
        Would you like to upload it and sync it to your cloud account?
      </p>

      <div class="d-flex justify-end gap-2">
        <v-btn variant="text" @click="close">
          Keep Offline
        </v-btn>
        <v-btn color="primary" variant="flat" prepend-icon="mdi-cloud-upload" @click="handleConfirm">
          Upload to Account
        </v-btn>
      </div>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>
