<script setup lang="ts">
import type { TrappingItem } from '../../../constants/placeholders'
import AddBtn from '../../ui/AddBtn.vue'
import DeleteRowBtn from '../../ui/DeleteRowBtn.vue'

defineProps<{
  bag: TrappingItem
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'addItem'): void
  (e: 'removeItem', subIndex: number): void
}>()
</script>

<template>
  <tr v-if="bag.isBag && isOpen" class="bag-contents-row">
    <td colspan="6" class="pa-2 bg-surface">
      <div class="bag-inner-box pa-2 rounded border">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="d-flex align-center">
          </div>
          <AddBtn label="Add Item" color="secondary" @click="emit('addItem')" />
        </div>

        <div v-if="!bag.containedTrappings || bag.containedTrappings.length === 0"
          class="text-caption text-medium-emphasis font-italic py-3 text-center">
          Bag is empty. Click "+ Add Item" to put items in this bag.
        </div>

        <v-table v-else density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th class="text-left font-weight-bold">Item Name</th>
              <th class="text-center font-weight-bold" style="width: 60px;">Enc</th>
              <th class="text-center font-weight-bold" style="width: 60px;">Qty</th>
              <th class="text-center font-weight-bold" style="width: 45px;">Worn</th>
              <th class="text-right font-weight-bold" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(sub, subIdx) in bag.containedTrappings" :key="sub.id || subIdx" class="sub-item-row">
              <td>
                <v-text-field v-model="sub.name" variant="plain" density="compact" hide-details
                  placeholder="Item Name" />
              </td>
              <td style="width: 60px;">
                <v-text-field v-model.number="sub.enc" type="number" variant="plain" density="compact" hide-details
                  class="text-center" placeholder="0" />
              </td>
              <td style="width: 60px;">
                <v-text-field v-model.number="sub.qty" type="number" variant="plain" density="compact" hide-details
                  class="text-center" placeholder="1" />
              </td>
              <td style="width: 45px;" class="text-center">
                <v-checkbox-btn v-model="sub.worn" density="compact" hide-details color="primary" />
              </td>
              <td style="width: 40px;" class="text-right">
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
  border-style: dashed;
}
</style>
