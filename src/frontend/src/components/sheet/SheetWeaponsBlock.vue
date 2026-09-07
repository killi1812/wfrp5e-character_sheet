<script setup lang="ts">
import type { Weapon } from '../../constants/placeholders'

defineProps<{
  weapons: Weapon[]
}>()

const emit = defineEmits<{
  (e: 'addWeapon'): void
  (e: 'removeWeapon', index: number): void
}>()
</script>

<template>
  <v-card color="surface" elevation="2" class="mb-4 pa-4 rounded-lg border">
    <div class="d-flex justify-space-between align-center mb-2">
      <div class="text-subtitle-2 font-weight-black text-uppercase text-primary">Weapons</div>
      <v-btn color="primary" size="small" prepend-icon="mdi-plus" variant="tonal" @click="emit('addWeapon')">
        Add Weapon
      </v-btn>
    </div>
    <v-table density="compact" class="bg-transparent text-caption">
      <thead>
        <tr>
          <th>Name</th>
          <th>Group</th>
          <th class="text-center">Enc</th>
          <th>Range/Reach</th>
          <th>Damage</th>
          <th>Qualities / Flaws</th>
          <th class="text-right">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(w, idx) in weapons" :key="idx">
          <td><v-text-field v-model="w.name" variant="plain" density="compact" hide-details /></td>
          <td style="max-width: 90px;"><v-text-field v-model="w.group" variant="plain" density="compact" hide-details /></td>
          <td style="max-width: 60px;"><v-text-field v-model.number="w.enc" type="number" variant="plain" density="compact" hide-details class="text-center" /></td>
          <td style="max-width: 100px;"><v-text-field v-model="w.rangeReach" variant="plain" density="compact" hide-details /></td>
          <td style="max-width: 90px;"><v-text-field v-model="w.damage" variant="plain" density="compact" hide-details /></td>
          <td><v-text-field v-model="w.qualities" variant="plain" density="compact" hide-details /></td>
          <td class="text-right"><v-btn icon="mdi-delete" size="x-small" color="error" variant="text" @click="emit('removeWeapon', idx)" /></td>
        </tr>
      </tbody>
    </v-table>
  </v-card>
</template>
