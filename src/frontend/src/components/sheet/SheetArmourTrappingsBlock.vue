<script setup lang="ts">
import type { ArmourItem, TrappingItem } from '../../constants/placeholders'

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
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Armour Gear</div>
          <v-btn color="primary" size="x-small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addArmour')">
            Add Armour
          </v-btn>
        </div>
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Name</th>
              <th>Locations</th>
              <th class="text-center">Enc</th>
              <th class="text-center">AP</th>
              <th>Qualities</th>
              <th class="text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(arm, idx) in armour" :key="idx">
              <td><v-text-field v-model="arm.name" variant="plain" density="compact" hide-details /></td>
              <td><v-text-field v-model="arm.locations" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="arm.enc" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="arm.ap" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td><v-text-field v-model="arm.qualities" variant="plain" density="compact" hide-details /></td>
              <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeArmour', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>

    <!-- Trappings & Wealth -->
    <v-col cols="12" md="6">
      <v-card color="surface" elevation="2" class="pa-4 rounded-lg border h-100">
        <div class="d-flex justify-space-between align-center mb-2">
          <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Trappings & Wealth</div>
          <v-btn color="primary" size="x-small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addTrapping')">
            Add Item
          </v-btn>
        </div>

        <!-- Wealth Counter Row -->
        <div class="pa-2 bg-surface-variant rounded-lg mb-3 d-flex align-center justify-space-between text-caption border">
          <span class="font-weight-bold text-primary">Money:</span>
          <div class="d-flex gap-2 align-center">
            <v-text-field v-model.number="wealth.gc" label="GC (Gold)" type="number" variant="plain" density="compact" hide-details style="width: 75px;" />
            <v-text-field v-model.number="wealth.ss" label="SS (Silver)" type="number" variant="plain" density="compact" hide-details style="width: 75px;" />
            <v-text-field v-model.number="wealth.bp" label="BP (Brass)" type="number" variant="plain" density="compact" hide-details style="width: 75px;" />
          </div>
        </div>

        <!-- Encumbrance Summary -->
        <div class="d-flex justify-space-between text-caption font-weight-bold mb-2 text-primary">
          <span>Encumbrance: {{ computedTotalEnc }} / {{ computedMaxEnc }} Max</span>
          <v-progress-linear :model-value="(computedTotalEnc / (computedMaxEnc || 1)) * 100" color="primary" height="6" rounded />
        </div>

        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th>Item Name</th>
              <th class="text-center">Enc</th>
              <th class="text-center">Qty</th>
              <th class="text-right">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(t, idx) in trappings" :key="idx">
              <td><v-text-field v-model="t.name" variant="plain" density="compact" hide-details /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="t.enc" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td style="max-width: 55px;"><v-text-field v-model.number="t.qty" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
              <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeTrapping', idx)" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
.gap-2 { gap: 8px; }
</style>
