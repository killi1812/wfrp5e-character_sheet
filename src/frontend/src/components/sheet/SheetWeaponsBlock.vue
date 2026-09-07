<script setup lang="ts">
import type { Weapon } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

const props = defineProps<{
  weapons: Weapon[]
  getCharBonus: (code: string) => number
}>()

const emit = defineEmits<{
  (e: 'addWeapon'): void
  (e: 'removeWeapon', index: number): void
}>()

/** Parse damage like "+SB+4" and compute total if SB is present */
function formatDamage(damage: string): string {
  if (!damage) return ''
  const sbMatch = damage.match(/\+?\s*SB\s*\+?\s*(\d+)/i)
  if (sbMatch) {
    const bonus = Number(sbMatch[1]) || 0
    const sb = props.getCharBonus('S')
    return `${damage} (${sb + bonus})`
  }
  return damage
}
</script>

<template>
  <SectionCard title="Weapons" add-label="Add Weapon" add-size="small" mb="mb-4" @add="emit('addWeapon')">
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
          <td><v-text-field v-model="w.name" variant="plain" density="compact" hide-details placeholder="Weapon Name" /></td>
          <td style="max-width: 90px;"><v-text-field v-model="w.group" variant="plain" density="compact" hide-details placeholder="Basic" /></td>
          <td style="max-width: 60px;"><v-text-field v-model.number="w.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
          <td style="max-width: 100px;"><v-text-field v-model="w.rangeReach" variant="plain" density="compact" hide-details placeholder="Melee" /></td>
          <td style="max-width: 120px;">
            <v-tooltip :text="'Computed: ' + formatDamage(w.damage)" location="top" :disabled="!w.damage">
              <template #activator="{ props: tProps }">
                <div v-bind="tProps" class="d-flex align-center">
                  <v-text-field v-model="w.damage" variant="plain" density="compact" hide-details placeholder="+SB+4" />
                  <span v-if="w.damage && w.damage.toLowerCase().includes('sb')" class="text-caption text-primary font-weight-bold ml-1 text-no-wrap">
                    ({{ formatDamage(w.damage).match(/\((\d+)\)/)?.[1] }})
                  </span>
                </div>
              </template>
            </v-tooltip>
          </td>
          <td><v-text-field v-model="w.qualities" variant="plain" density="compact" hide-details placeholder="Qualities..." /></td>
          <td class="text-right"><DeleteRowBtn @delete="emit('removeWeapon', idx)" /></td>
        </tr>
      </tbody>
    </v-table>
  </SectionCard>
</template>
