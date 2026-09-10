<script setup lang="ts">
import type { TrappingItem } from '../../constants/placeholders'
import AddBtn from '../ui/AddBtn.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

defineProps<{
  bag: TrappingItem
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'addItem'): void
  (e: 'removeItem', subIndex: number): void
  (e: 'dropIntoBag', event: DragEvent): void
  (e: 'dragStartItem', payload: { event: DragEvent; item: TrappingItem }): void
}>()
</script>

<template>
  <tr v-if="bag.isBag && isOpen" class="bag-contents-row">
    <td colspan="6" class="pa-2 bg-surface">
      <div
        class="bag-inner-box pa-2 rounded border"
        @dragover.prevent
        @drop.stop="emit('dropIntoBag', $event)"
      >
        <div class="d-flex justify-space-between align-center mb-1">
          <span class="text-caption font-weight-bold text-secondary">
            Contents of {{ bag.name || 'Bag' }} (Excluded from personal Enc)
          </span>
          <AddBtn label="Add item to bag" color="secondary" variant="text" @click="emit('addItem')" />
        </div>

        <div
          v-if="!bag.containedTrappings || bag.containedTrappings.length === 0"
          class="text-caption text-medium-emphasis font-italic py-2 text-center"
        >
          Bag is empty. Drag items here or click "+ Add item to bag".
        </div>

        <v-table v-else density="compact" class="bg-transparent text-caption">
          <tbody>
            <tr
              v-for="(sub, subIdx) in bag.containedTrappings"
              :key="sub.id || subIdx"
              draggable="true"
              class="sub-item-row"
              @dragstart="emit('dragStartItem', { event: $event, item: sub })"
            >
              <td style="width: 24px;" class="drag-handle text-center">
                <v-icon icon="mdi-drag-vertical" size="small" color="medium-emphasis" />
              </td>
              <td>
                <v-text-field
                  v-model="sub.name"
                  variant="plain"
                  density="compact"
                  hide-details
                  placeholder="Item Name"
                />
              </td>
              <td style="width: 50px;">
                <v-text-field
                  v-model.number="sub.enc"
                  type="number"
                  variant="plain"
                  density="compact"
                  hide-details
                  class="text-center"
                  placeholder="0"
                />
              </td>
              <td style="width: 50px;">
                <v-text-field
                  v-model.number="sub.qty"
                  type="number"
                  variant="plain"
                  density="compact"
                  hide-details
                  class="text-center"
                  placeholder="1"
                />
              </td>
              <td style="width: 40px;" class="text-center">
                <v-checkbox-btn v-model="sub.worn" density="compact" hide-details color="primary" />
              </td>
              <td style="width: 36px;" class="text-right">
                <DeleteRowBtn @delete="emit('removeItem', subIdx)" />
              </td>
            </tr>
          </tbody>
        </v-table>
      </div>
    </td>
  </tr>
</template>

<style scoped>
.bag-inner-box {
  background: rgba(var(--v-theme-surface-variant), 0.35);
  border-style: dashed !important;
}

.sub-item-row {
  cursor: grab;
}
.sub-item-row:active {
  cursor: grabbing;
}
.drag-handle {
  cursor: grab;
}
</style>
