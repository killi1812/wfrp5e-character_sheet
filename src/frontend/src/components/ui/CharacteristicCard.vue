<script setup lang="ts">
import ImportantBtn from './ImportantBtn.vue'

const props = withDefaults(
  defineProps<{
    code: string
    name?: string
    hint?: string
    initial?: number | null
    advances?: number | null
    current: number
    isImportant?: boolean
    showImportant?: boolean
    nullable?: boolean
    isNull?: boolean
    penalty?: number
    penaltyTooltip?: string
    useBadge?: boolean
  }>(),
  {
    name: '',
    hint: '',
    initial: 0,
    advances: 0,
    isImportant: false,
    showImportant: false,
    nullable: false,
    isNull: false,
    penalty: 0,
    penaltyTooltip: '',
    useBadge: true,
  }
)

const emit = defineEmits<{
  (e: 'update:initial', val: number): void
  (e: 'update:advances', val: number): void
  (e: 'toggleImportant'): void
  (e: 'toggleNullable'): void
}>()
</script>

<template>
  <v-card color="surface-variant" variant="outlined" class="pa-2 rounded text-center char-card position-relative"
    :class="{ 'char-important': isImportant }">
    <!-- Header with Code, Tooltip, and Actions (Important or Nullable Toggle) -->
    <div class="char-header position-relative border-bottom pb-1 mb-2 d-flex align-center"
      :class="nullable ? 'justify-space-between' : 'justify-center'">
      <v-tooltip :text="name ? `${name}${hint ? ' — ' + hint : ''}` : code" location="top" :open-on-focus="false">
        <template #activator="{ props: tProps }">
          <div v-bind="tProps" class="text-subtitle-2 font-weight-black text-primary cursor-pointer px-1">
            {{ code }}
          </div>
        </template>
      </v-tooltip>

      <!-- Important toggle button (!) -->
      <ImportantBtn v-if="showImportant" as-char :active="isImportant" class="important-toggle-btn"
        @toggle="emit('toggleImportant')" />

    </div>

    <!-- Disabled / Null State -->
    <div v-if="isNull" class="py-3 text-medium-emphasis font-italic text-caption">
      —
    </div>

    <!-- Active State: Initial, Advances, Current Total -->
    <template v-else>
      <!-- Initial Field -->
      <div class="d-flex align-center justify-space-between mb-1 px-1">
        <span class="text-caption text-high-emphasis font-weight-medium" style="font-size: 0.75rem;">Init</span>
        <input :value="initial ?? 0" type="number" class="char-num-input" placeholder="0"
          @input="emit('update:initial', Number(($event.target as HTMLInputElement).value) || 0)" />
      </div>

      <!-- Advances Field -->
      <div class="d-flex align-center justify-space-between mb-2 px-1">
        <span class="text-caption text-high-emphasis font-weight-medium" style="font-size: 0.75rem;">Adv</span>
        <input :value="advances ?? 0" type="number" class="char-num-input adv-input" placeholder="0"
          @input="emit('update:advances', Number(($event.target as HTMLInputElement).value) || 0)" />
      </div>

      <!-- Total Score Display -->
      <v-tooltip v-if="penalty > 0 || penaltyTooltip"
        :text="penaltyTooltip || `Base: ${current + penalty} (Penalized -${penalty})`" location="bottom"
        :open-on-focus="false">
        <template #activator="{ props: tProps }">
          <div v-if="useBadge" v-bind="tProps" class="stat-total-badge text-on-primary font-weight-black rounded py-1"
            :class="penalty > 0 ? 'bg-error' : 'bg-primary'">
            {{ current }}
            <span v-if="penalty > 0" class="text-caption font-weight-bold ml-1">
              (-{{ penalty }})
            </span>
          </div>
          <div v-else v-bind="tProps" class="border-top pt-1 text-center font-weight-black text-body-2"
            :class="penalty > 0 ? 'text-error' : 'text-high-emphasis'">
            {{ current }} (-{{ penalty }})
          </div>
        </template>
      </v-tooltip>

      <template v-else>
        <div v-if="useBadge" class="stat-total-badge text-on-primary font-weight-black rounded py-1 bg-primary">
          {{ current }}
        </div>
        <div v-else class="border-top pt-1 text-center font-weight-black text-body-2 text-high-emphasis">
          {{ current }}
        </div>
      </template>
    </template>
  </v-card>
</template>

<style scoped>
.char-card {
  transition: border-color 0.2s, box-shadow 0.2s;
}

.char-card:hover {
  border-color: rgba(var(--v-theme-primary), 0.5) !important;
}

.char-num-input {
  width: 44px;
  text-align: center;
  background: rgba(var(--v-theme-surface), 0.9);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.2);
  border-radius: 4px;
  color: rgb(var(--v-theme-on-surface));
  font-size: 0.8rem;
  font-weight: 700;
  padding: 1px 2px;
  outline: none;
}

.char-num-input:focus {
  border-color: rgb(var(--v-theme-primary));
  box-shadow: 0 0 0 1px rgb(var(--v-theme-primary));
}

.adv-input {
  background: rgba(var(--v-theme-primary), 0.08);
  border-color: rgba(var(--v-theme-primary), 0.35);
  color: rgb(var(--v-theme-primary));
}

.stat-total-badge {
  font-size: 0.95rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  transition: background-color 0.2s;
}

.char-important {
  border-color: rgb(var(--v-theme-warning));
  box-shadow: 0 0 8px rgba(var(--v-theme-warning), 0.35);
  background: rgba(var(--v-theme-warning), 0.05);
}

.important-toggle-btn {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

.char-card:hover :deep(.important-char-btn),
:deep(.important-char-btn.active) {
  opacity: 1;
}

:deep(.important-char-btn:not(.active)) {
  opacity: 0;
}
</style>
