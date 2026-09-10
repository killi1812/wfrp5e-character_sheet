<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    active?: boolean
    tooltipActive?: string
    tooltipInactive?: string
    asChar?: boolean
    size?: string
    iconSize?: number | string
  }>(),
  {
    active: false,
    tooltipActive: 'Unmark as Important',
    tooltipInactive: 'Mark as Important',
    asChar: false,
    size: 'x-small',
    iconSize: 14,
  }
)

const emit = defineEmits<{
  (e: 'update:active', value: boolean): void
  (e: 'toggle'): void
}>()

function handleClick(event: MouseEvent) {
  event.stopPropagation()
  const next = !props.active
  emit('update:active', next)
  emit('toggle')
}
</script>

<template>
  <v-tooltip :text="active ? tooltipActive : tooltipInactive" location="top" :open-on-focus="false">
    <template #activator="{ props: tProps }">
      <!-- Char style '!' button (e.g. Characteristics header) -->
      <button
        v-if="asChar"
        v-bind="tProps"
        type="button"
        class="important-char-btn"
        :class="{ active }"
        @click="handleClick"
      >
        !
      </button>

      <!-- Icon button (e.g. Skills, Talents, Languages) -->
      <v-btn
        v-else
        v-bind="tProps"
        icon
        :size="size"
        variant="text"
        density="compact"
        class="important-icon-btn"
        :class="{ active }"
        @click="handleClick"
      >
        <v-icon :size="iconSize">
          {{ active ? 'mdi-alert-circle' : 'mdi-alert-circle-outline' }}
        </v-icon>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<style scoped>
/* Character '!' style */
.important-char-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1px solid rgba(var(--v-theme-primary), 0.3);
  background: transparent;
  color: rgb(var(--v-theme-primary));
  font-size: 10px;
  font-weight: 900;
  line-height: 1;
  cursor: pointer;
  padding: 0;
  opacity: 0.4;
  transition: all 0.2s ease;
}

.important-char-btn:hover {
  opacity: 1;
  background: rgba(var(--v-theme-primary), 0.15);
  border-color: rgb(var(--v-theme-primary));
}

.important-char-btn.active {
  opacity: 1;
  background: rgb(var(--v-theme-warning));
  color: rgb(var(--v-theme-on-warning));
  border-color: rgb(var(--v-theme-warning));
}

/* Icon style */
.important-icon-btn {
  opacity: 0.35;
  transition: opacity 0.2s ease, color 0.2s ease;
  color: rgb(var(--v-theme-on-surface));
}

.important-icon-btn:hover {
  opacity: 1;
  color: rgb(var(--v-theme-warning));
}

.important-icon-btn.active {
  opacity: 1;
  color: rgb(var(--v-theme-warning)) !important;
}
</style>
