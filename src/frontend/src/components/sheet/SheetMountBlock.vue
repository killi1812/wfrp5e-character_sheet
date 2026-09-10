<script setup lang="ts">
import { computed } from 'vue'
import {
  type MountData,
  type Skill,
  STAT_KEYS,
  STAT_NAMES,
  type StatKey,
} from '../../constants/placeholders'
import { formatDamage } from '../../utils/damage'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'
import CharacteristicCard from '../ui/CharacteristicCard.vue'

const props = defineProps<{
  mount: MountData
}>()

const emit = defineEmits<{
  (e: 'addAttack'): void
  (e: 'removeAttack', index: number): void
  (e: 'addSkill'): void
  (e: 'removeSkill', index: number): void
  (e: 'addTrait'): void
  (e: 'removeTrait', index: number): void
  (e: 'addTrapping'): void
  (e: 'removeTrapping', index: number): void
}>()

function getMountStatCurrent(code: string): number | null {
  const stat = props.mount.characteristics?.[code]
  if (!stat) return null
  return (Number(stat.initial) || 0) + (Number(stat.advances) || 0)
}

function getMountStatBonus(code: string): number {
  const cur = getMountStatCurrent(code)
  if (cur === null) return 0
  return Math.floor(cur / 10)
}

function enableStat(code: string) {
  if (!props.mount.characteristics) {
    props.mount.characteristics = {}
  }
  props.mount.characteristics[code] = {
    name: STAT_NAMES[code as StatKey] || code,
    initial: 20,
    advances: 0,
    hint: '',
  }
}

function clearStat(code: string) {
  if (props.mount.characteristics) {
    props.mount.characteristics[code] = null
  }
}

const mountMaxEnc = computed(() => {
  const sb = getMountStatBonus('S')
  const tb = getMountStatBonus('T')
  return sb + tb
})

const mountTotalEnc = computed(() => {
  let total = 0
  if (props.mount.trappings) {
    props.mount.trappings.forEach((t) => {
      total += (Number(t.enc) || 0) * (Number(t.qty) || 1)
    })
  }
  return total
})

function getMountSkillTotal(skill: Skill): number {
  const base = getMountStatCurrent(skill.characteristic)
  if (base === null) return Number(skill.adv) || 0
  return base + (Number(skill.adv) || 0)
}
</script>

<template>
  <SectionCard title="Mount / Steed" full-height mb="mb-4">
    <div class="mb-4">
      <v-row dense align="center">
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            v-model="mount.name"
            label="Mount Name / Breed"
            placeholder="e.g. Warhorse, Elven Steed"
            variant="outlined"
            density="compact"
            hide-details
          />
        </v-col>
        <v-col cols="12" sm="6" md="8" class="d-flex justify-sm-end align-center">
          <!-- Mount Encumbrance Indicator -->
          <v-chip
            :color="mountTotalEnc > mountMaxEnc ? 'error' : 'secondary'"
            variant="tonal"
            class="font-weight-bold"
          >
            <v-icon start icon="mdi-weight" />
            Mount Encumbrance: {{ mountTotalEnc }} / {{ mountMaxEnc }}
          </v-chip>
        </v-col>
      </v-row>
    </div>

    <!-- Mount Characteristics Grid -->
    <div class="mb-4">
      <div class="text-caption font-weight-bold text-primary text-uppercase mb-2">
        Characteristics
      </div>
      <v-row dense>
        <v-col
          v-for="code in STAT_KEYS"
          :key="code"
          cols="6"
          sm="4"
          md="2"
          lg="1"
          class="flex-grow-1"
        >
          <CharacteristicCard
            :code="code"
            :name="STAT_NAMES[code]"
            :initial="mount.characteristics?.[code]?.initial"
            :advances="mount.characteristics?.[code]?.advances"
            :current="getMountStatCurrent(code) ?? 0"
            :nullable="code === 'BS' || code === 'Dex'"
            :is-null="!mount.characteristics?.[code]"
            :use-badge="false"
            @update:initial="mount.characteristics[code] && (mount.characteristics[code]!.initial = $event)"
            @update:advances="mount.characteristics[code] && (mount.characteristics[code]!.advances = $event)"
            @toggle-nullable="mount.characteristics?.[code] ? clearStat(code) : enableStat(code)"
          />
        </v-col>
      </v-row>
    </div>

    <!-- Mount Attacks & Skills Row -->
    <v-row dense class="mb-4">
      <!-- Attacks -->
      <v-col cols="12" md="6">
        <SectionCard title="Mount Attacks" add-label="Add Attack" full-height @add="emit('addAttack')">
          <v-table density="compact" class="bg-transparent text-caption">
            <thead>
              <tr>
                <th>Attack</th>
                <th style="width: 70px;">Skill</th>
                <th class="text-center" style="width: 50px;">Val</th>
                <th>Damage</th>
                <th>Qualities</th>
                <th class="text-right" style="width: 40px;">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(att, idx) in mount.attacks" :key="idx">
                <td><v-text-field v-model="att.name" variant="plain" density="compact" hide-details placeholder="Bite / Trample" /></td>
                <td style="max-width: 70px;"><v-text-field v-model="att.skillToRoll" variant="plain" density="compact" hide-details placeholder="WS" /></td>
                <td style="max-width: 50px;"><v-text-field v-model.number="att.displayValue" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
                <td style="max-width: 110px;">
                  <v-tooltip :text="'Computed: ' + formatDamage(att.damage, getMountStatBonus('S'))" location="top" :disabled="!att.damage">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="d-flex align-center">
                        <v-text-field v-model="att.damage" variant="plain" density="compact" hide-details placeholder="+SB+4" />
                        <span v-if="att.damage && att.damage.toLowerCase().includes('sb')" class="text-caption text-primary font-weight-bold ml-1 text-no-wrap">
                          ({{ formatDamage(att.damage, getMountStatBonus('S')).match(/\((\d+)\)/)?.[1] }})
                        </span>
                      </div>
                    </template>
                  </v-tooltip>
                </td>
                <td>
                  <v-textarea
                    v-model="att.qualities"
                    variant="plain"
                    density="compact"
                    rows="1"
                    auto-grow
                    hide-details
                    placeholder="Qualities..."
                  />
                </td>
                <td class="text-right"><DeleteRowBtn @delete="emit('removeAttack', idx)" /></td>
              </tr>
            </tbody>
          </v-table>
        </SectionCard>
      </v-col>

      <!-- Advanced Skills -->
      <v-col cols="12" md="6">
        <SectionCard title="Mount Skills" add-label="Add Skill" full-height @add="emit('addSkill')">
          <v-table density="compact" class="bg-transparent text-caption">
            <thead>
              <tr>
                <th>Skill</th>
                <th class="text-center" style="width: 60px;">Char</th>
                <th class="text-center" style="width: 60px;">Adv</th>
                <th class="text-right" style="width: 60px;">Total</th>
                <th class="text-right" style="width: 40px;">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(sk, idx) in mount.skills" :key="idx">
                <td><v-text-field v-model="sk.name" variant="plain" density="compact" hide-details placeholder="Skill Name" /></td>
                <td style="max-width: 60px;"><v-text-field v-model="sk.characteristic" variant="plain" density="compact" hide-details class="text-center text-uppercase" placeholder="Ag" /></td>
                <td style="max-width: 60px;"><v-text-field v-model.number="sk.adv" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
                <td class="text-right font-weight-bold">
                  {{ getMountSkillTotal(sk) }}
                </td>
                <td class="text-right"><DeleteRowBtn @delete="emit('removeSkill', idx)" /></td>
              </tr>
            </tbody>
          </v-table>
        </SectionCard>
      </v-col>
    </v-row>

    <!-- Mount Traits & Trappings Row -->
    <v-row dense>
      <!-- Traits -->
      <v-col cols="12" md="6">
        <SectionCard title="Mount Traits" add-label="Add Trait" full-height @add="emit('addTrait')">
          <v-table density="compact" class="bg-transparent text-caption">
            <thead>
              <tr>
                <th style="width: 140px;">Trait</th>
                <th>Description</th>
                <th class="text-right" style="width: 40px;">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(tr, idx) in mount.traits" :key="idx">
                <td style="width: 140px; max-width: 140px;"><v-text-field v-model="tr.name" variant="plain" density="compact" hide-details placeholder="e.g. Quadruped" /></td>
                <td>
                  <v-textarea
                    v-model="tr.desc"
                    variant="plain"
                    density="compact"
                    rows="1"
                    auto-grow
                    hide-details
                    placeholder="Description..."
                  />
                </td>
                <td class="text-right"><DeleteRowBtn @delete="emit('removeTrait', idx)" /></td>
              </tr>
            </tbody>
          </v-table>
        </SectionCard>
      </v-col>

      <!-- Trappings -->
      <v-col cols="12" md="6">
        <SectionCard title="Mount Trappings" add-label="Add Trapping" full-height @add="emit('addTrapping')">
          <v-table density="compact" class="bg-transparent text-caption">
            <thead>
              <tr>
                <th>Item</th>
                <th style="width: 90px;">Category</th>
                <th class="text-center" style="width: 50px;">Enc</th>
                <th class="text-center" style="width: 50px;">Qty</th>
                <th>Description</th>
                <th class="text-right" style="width: 40px;">Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(trap, idx) in mount.trappings" :key="idx">
                <td><v-text-field v-model="trap.name" variant="plain" density="compact" hide-details placeholder="Saddle / Harness" /></td>
                <td style="max-width: 90px;"><v-text-field v-model="trap.category" variant="plain" density="compact" hide-details placeholder="Tack" /></td>
                <td style="max-width: 50px;"><v-text-field v-model.number="trap.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" /></td>
                <td style="max-width: 50px;"><v-text-field v-model.number="trap.qty" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="1" /></td>
                <td>
                  <v-textarea
                    v-model="trap.desc"
                    variant="plain"
                    density="compact"
                    rows="1"
                    auto-grow
                    hide-details
                    placeholder="Description..."
                  />
                </td>
                <td class="text-right"><DeleteRowBtn @delete="emit('removeTrapping', idx)" /></td>
              </tr>
            </tbody>
          </v-table>
        </SectionCard>
      </v-col>
    </v-row>
  </SectionCard>
</template>

<style scoped>
:deep(td) {
  white-space: normal;
}
</style>
