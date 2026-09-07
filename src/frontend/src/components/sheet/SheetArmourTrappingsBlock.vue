<script setup lang="ts">
import type { ArmourItem, TrappingItem } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

defineProps<{
  armour: ArmourItem[]
  trappings: TrappingItem[]
  wealth: { gc: number; ss: number; bp: number }
  computedTotalEnc: number
  computedMaxEnc: number
}>()

const emit = defineEmits<{
  (e: 'addArmour'): void
  (e: 'removeArmour', index: number): void
  (e: 'addTrapping'): void
  (e: 'removeTrapping', index: number): void
}>()
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Armour Equipment -->
    <v-col cols="12" md="6">
      <SectionCard title="Armour Gear" add-label="Add Armour" full-height @add="emit('addArmour')">
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th style="min-width: 150px;">Name</th>
              <th>Locations</th>
              <th class="text-center" style="width: 55px;">Enc</th>
              <th class="text-center" style="width: 55px;">AP</th>
              <th>Qualities</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(arm, idx) in armour" :key="idx">
              <td style="min-width: 150px;">
                <v-text-field v-model="arm.name" variant="plain" density="compact" hide-details placeholder="Armour Name" />
              </td>
              <td>
                <v-text-field v-model="arm.locations" variant="plain" density="compact" hide-details placeholder="Body, Arms" />
              </td>
              <td style="max-width: 55px;">
                <v-text-field v-model.number="arm.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
              </td>
              <td style="max-width: 55px;">
                <v-text-field v-model.number="arm.ap" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
              </td>
              <td>
                <v-text-field v-model="arm.qualities" variant="plain" density="compact" hide-details placeholder="Flexible" />
              </td>
              <td class="text-right">
                <DeleteRowBtn @delete="emit('removeArmour', idx)" />
              </td>
            </tr>
          </tbody>
        </v-table>
      </SectionCard>
    </v-col>

    <!-- Trappings & Wealth -->
    <v-col cols="12" md="6">
      <SectionCard title="Trappings & Wealth" add-label="Add Item" full-height @add="emit('addTrapping')">
        <!-- Encumbrance Summary (without the progress bar line) -->
        <div class="d-flex justify-space-between align-center text-caption font-weight-bold mb-2 text-primary">
          <span>Encumbrance: {{ computedTotalEnc }} / {{ computedMaxEnc }} Max</span>
        </div>

        <!-- Trappings table on top -->
        <v-table density="compact" class="bg-transparent text-caption mb-3">
          <thead>
            <tr>
              <th>Item Name</th>
              <th class="text-center" style="width: 55px;">Enc</th>
              <th class="text-center" style="width: 55px;">Qty</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(t, idx) in trappings" :key="idx">
              <td><v-text-field v-model="t.name" variant="plain" density="compact" hide-details placeholder="Item Name" /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="t.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="t.qty" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="1" /></td>
              <td class="text-right"><DeleteRowBtn @delete="emit('removeTrapping', idx)" /></td>
            </tr>
          </tbody>
        </v-table>

        <!-- Wealth moved down with bigger font and better styling -->
        <div class="pa-3 bg-surface-variant rounded-lg border">
          <div class="text-caption font-weight-bold text-primary mb-2">Money / Wealth</div>
          <v-row dense>
            <v-col cols="4">
              <div class="wealth-box text-center pa-2 rounded border bg-surface">
                <div class="text-caption text-medium-emphasis font-weight-bold mb-1">GC (Gold)</div>
                <input v-model.number="wealth.gc" type="number" class="wealth-input font-weight-bold" placeholder="0" />
              </div>
            </v-col>
            <v-col cols="4">
              <div class="wealth-box text-center pa-2 rounded border bg-surface">
                <div class="text-caption text-medium-emphasis font-weight-bold mb-1">SS (Silver)</div>
                <input v-model.number="wealth.ss" type="number" class="wealth-input font-weight-bold" placeholder="0" />
              </div>
            </v-col>
            <v-col cols="4">
              <div class="wealth-box text-center pa-2 rounded border bg-surface">
                <div class="text-caption text-medium-emphasis font-weight-bold mb-1">BP (Brass)</div>
                <input v-model.number="wealth.bp" type="number" class="wealth-input font-weight-bold" placeholder="0" />
              </div>
            </v-col>
          </v-row>
        </div>
      </SectionCard>
    </v-col>
  </v-row>
</template>

<style scoped>
.wealth-box {
  transition: border-color 0.2s;
}
.wealth-box:focus-within {
  border-color: rgb(var(--v-theme-primary)) !important;
}

.wealth-input {
  font-size: 1.25rem;
  width: 100%;
  text-align: center;
  border: none;
  background: transparent;
  outline: none;
  color: currentColor;
}
</style>
