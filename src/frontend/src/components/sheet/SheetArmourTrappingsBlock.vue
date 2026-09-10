<script setup lang="ts">
import { ref, computed } from 'vue'
import type { ArmourItem, TrappingItem } from '../../constants/placeholders'
import { NEW_ITEM_TEMPLATES } from '../../constants/placeholders'
import SectionCard from '../ui/SectionCard.vue'
import DeleteRowBtn from '../ui/DeleteRowBtn.vue'

const props = defineProps<{
  armour: ArmourItem[]
  trappings: TrappingItem[]
  wealth: { gc: number; ss: number; bp: number }
  computedTotalEnc: number
  computedMaxEnc: number
  encBreakdown?: {
    weapons: number
    armour: number
    trappings: number
    coins: number
    totalCoins: number
  }
}>()

const emit = defineEmits<{
  (e: 'addArmour'): void
  (e: 'removeArmour', index: number): void
  (e: 'addTrapping'): void
  (e: 'addBag'): void
  (e: 'removeTrapping', index: number): void
}>()

// Track open/collapsed state of bags
const openBags = ref<Record<string, boolean>>({})

function toggleBag(bagId: string) {
  openBags.value[bagId] = !openBags.value[bagId]
}

function getContainedEnc(bag: TrappingItem): number {
  if (!bag.containedTrappings) return 0
  return bag.containedTrappings.reduce((sum, item) => {
    const enc = item.worn ? Math.max(0, (Number(item.enc) || 0) - 1) : (Number(item.enc) || 0)
    return sum + enc * (Number(item.qty) || 1)
  }, 0)
}

function addItemToBag(bag: TrappingItem) {
  if (!bag.containedTrappings) {
    bag.containedTrappings = []
  }
  bag.containedTrappings.push(NEW_ITEM_TEMPLATES.trapping())
  if (bag.id) {
    openBags.value[bag.id] = true
  }
}

function removeItemFromBag(bag: TrappingItem, index: number) {
  bag.containedTrappings?.splice(index, 1)
}

// Drag & Drop Handling
const draggedItemId = ref<string | null>(null)
const draggedFromBagId = ref<string | null>(null)

function onDragStart(e: DragEvent, item: TrappingItem, fromBagId: string | null = null) {
  draggedItemId.value = item.id || null
  draggedFromBagId.value = fromBagId
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', JSON.stringify({ itemId: item.id, fromBagId }))
  }
}

function onDropIntoBag(e: DragEvent, targetBag: TrappingItem) {
  e.preventDefault()
  if (!targetBag.containedTrappings) {
    targetBag.containedTrappings = []
  }

  const itemId = draggedItemId.value
  const fromBagId = draggedFromBagId.value
  if (!itemId || targetBag.id === itemId) return

  let foundItem: TrappingItem | null = null

  // Remove from source
  if (fromBagId) {
    const sourceBag = props.trappings.find((t) => t.id === fromBagId)
    if (sourceBag && sourceBag.containedTrappings) {
      const idx = sourceBag.containedTrappings.findIndex((t) => t.id === itemId)
      if (idx >= 0) {
        foundItem = sourceBag.containedTrappings.splice(idx, 1)[0]
      }
    }
  } else {
    const idx = props.trappings.findIndex((t) => t.id === itemId)
    if (idx >= 0) {
      foundItem = props.trappings.splice(idx, 1)[0]
    }
  }

  // Add to target bag
  if (foundItem) {
    targetBag.containedTrappings.push(foundItem)
    if (targetBag.id) {
      openBags.value[targetBag.id] = true
    }
  }

  draggedItemId.value = null
  draggedFromBagId.value = null
}

function onDropToRoot(e: DragEvent) {
  e.preventDefault()
  const itemId = draggedItemId.value
  const fromBagId = draggedFromBagId.value
  if (!itemId || !fromBagId) return

  const sourceBag = props.trappings.find((t) => t.id === fromBagId)
  if (sourceBag && sourceBag.containedTrappings) {
    const idx = sourceBag.containedTrappings.findIndex((t) => t.id === itemId)
    if (idx >= 0) {
      const item = sourceBag.containedTrappings.splice(idx, 1)[0]
      props.trappings.push(item)
    }
  }

  draggedItemId.value = null
  draggedFromBagId.value = null
}

const totalCoins = computed(() => (Number(props.wealth.gc) || 0) + (Number(props.wealth.ss) || 0) + (Number(props.wealth.bp) || 0))
const coinsEnc = computed(() => Math.floor(totalCoins.value / 200))
</script>

<template>
  <v-row dense class="mb-4">
    <!-- Armour Equipment -->
    <v-col cols="12" md="6">
      <SectionCard title="Armour Gear" add-label="Add Armour" full-height @add="emit('addArmour')">
        <v-table density="compact" class="bg-transparent text-caption">
          <thead>
            <tr>
              <th style="min-width: 140px;">Name</th>
              <th>Locations</th>
              <th class="text-center" style="width: 50px;">Enc</th>
              <th class="text-center" style="width: 50px;">AP</th>
              <th>Qualities</th>
              <th class="text-center" style="width: 45px;">Worn</th>
              <th class="text-right" style="width: 40px;">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(arm, idx) in armour" :key="idx">
              <td style="min-width: 140px;">
                <v-text-field v-model="arm.name" variant="plain" density="compact" hide-details placeholder="Armour Name" />
              </td>
              <td>
                <v-text-field v-model="arm.locations" variant="plain" density="compact" hide-details placeholder="Body, Arms" />
              </td>
              <td style="max-width: 50px;">
                <v-text-field v-model.number="arm.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
              </td>
              <td style="max-width: 50px;">
                <v-text-field v-model.number="arm.ap" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
              </td>
              <td>
                <v-textarea
                  v-model="arm.qualities"
                  variant="plain"
                  density="compact"
                  rows="1"
                  auto-grow
                  hide-details
                  placeholder="Flexible"
                />
              </td>
              <td class="text-center">
                <v-checkbox-btn v-model="arm.worn" density="compact" hide-details color="primary" />
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
      <SectionCard title="Trappings & Wealth" full-height>
        <template #actions>
          <v-btn size="x-small" color="primary" variant="tonal" prepend-icon="mdi-plus" class="mr-1" @click="emit('addTrapping')">
            Add Item
          </v-btn>
          <v-btn size="x-small" color="primary" variant="tonal" prepend-icon="mdi-plus" @click="emit('addBag')">
            Add Bag
          </v-btn>
        </template>

        <!-- Encumbrance Summary with Hover Breakdown -->
        <div class="d-flex justify-space-between align-center text-caption font-weight-bold mb-3 text-primary pa-2 bg-surface-variant rounded border">
          <v-tooltip location="top" :open-on-focus="false">
            <template #activator="{ props: tProps }">
              <span v-bind="tProps" class="cursor-pointer d-flex align-center">
                <v-icon icon="mdi-weight" size="small" class="mr-1" />
                Encumbrance: {{ computedTotalEnc }} / {{ computedMaxEnc }} Max
                <v-icon icon="mdi-information-outline" size="x-small" class="ml-1 text-medium-emphasis" />
              </span>
            </template>
            <div class="text-caption">
              <div class="font-weight-bold mb-1 border-bottom pb-1">Encumbrance Breakdown</div>
              <div>Weapons: {{ encBreakdown?.weapons || 0 }} Enc</div>
              <div>Armour: {{ encBreakdown?.armour || 0 }} Enc</div>
              <div>Trappings: {{ encBreakdown?.trappings || 0 }} Enc</div>
              <div>Coins: {{ coinsEnc }} Enc ({{ totalCoins }} coins)</div>
              <div class="font-weight-bold border-t mt-1 pt-1">
                Total: {{ computedTotalEnc }} / {{ computedMaxEnc }} Enc
              </div>
            </div>
          </v-tooltip>

          <span class="text-caption text-medium-emphasis">Drag items to move into bags</span>
        </div>

        <!-- Root Trappings Drop Zone -->
        <div class="root-trappings-container mb-3" @dragover.prevent @drop="onDropToRoot">
          <v-table density="compact" class="bg-transparent text-caption">
            <thead>
              <tr>
                <th style="width: 28px;"></th>
                <th>Item Name</th>
                <th class="text-center" style="width: 60px;">Enc</th>
                <th class="text-center" style="width: 70px;">Qty/Size</th>
                <th class="text-center" style="width: 45px;">Worn</th>
                <th class="text-right" style="width: 40px;">Action</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(t, idx) in trappings" :key="t.id || idx">
                <!-- REGULAR ITEM ROW -->
                <tr
                  v-if="!t.isBag"
                  draggable="true"
                  class="trapping-item-row"
                  @dragstart="onDragStart($event, t, null)"
                >
                  <td class="text-center pa-0 drag-handle" style="cursor: grab;">
                    <v-icon icon="mdi-drag-vertical" size="small" color="medium-emphasis" />
                  </td>
                  <td>
                    <v-text-field v-model="t.name" variant="plain" density="compact" hide-details placeholder="Item Name" />
                  </td>
                  <td style="max-width: 60px;">
                    <v-text-field v-model.number="t.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
                  </td>
                  <td style="max-width: 70px;">
                    <v-text-field v-model.number="t.qty" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="1" />
                  </td>
                  <td class="text-center">
                    <v-checkbox-btn v-model="t.worn" density="compact" hide-details color="primary" />
                  </td>
                  <td class="text-right">
                    <DeleteRowBtn @delete="emit('removeTrapping', idx)" />
                  </td>
                </tr>

                <!-- BAG / CONTAINER ROW -->
                <tr
                  v-else
                  class="bag-row bg-surface-variant"
                  @dragover.prevent
                  @drop.stop="onDropIntoBag($event, t)"
                >
                  <td class="text-center pa-0">
                    <v-btn
                      icon
                      size="x-small"
                      variant="text"
                      @click="toggleBag(t.id || String(idx))"
                    >
                      <v-icon :icon="openBags[t.id || String(idx)] ? 'mdi-chevron-down' : 'mdi-chevron-right'" />
                    </v-btn>
                  </td>
                  <td>
                    <div class="d-flex align-center">
                      <v-icon icon="mdi-bag-personal" size="small" color="secondary" class="mr-1" />
                      <v-text-field v-model="t.name" variant="plain" density="compact" hide-details placeholder="Bag Name" class="font-weight-bold" />
                    </div>
                  </td>
                  <td style="max-width: 60px;">
                    <v-text-field v-model.number="t.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="1" />
                  </td>
                  <td style="min-width: 100px;">
                    <div class="d-flex align-center gap-1 justify-center">
                      <v-text-field v-model.number="t.bagSize" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="5" style="max-width: 45px;" />
                      <span class="text-caption text-medium-emphasis font-weight-bold">
                        ({{ getContainedEnc(t) }}/{{ t.bagSize || 0 }} Enc)
                      </span>
                    </div>
                  </td>
                  <td class="text-center">
                    <v-checkbox-btn v-model="t.worn" density="compact" hide-details color="primary" />
                  </td>
                  <td class="text-right">
                    <DeleteRowBtn @delete="emit('removeTrapping', idx)" />
                  </td>
                </tr>

                <!-- BAG CONTENTS ACCORDION SECTION -->
                <tr v-if="t.isBag && openBags[t.id || String(idx)]" class="bag-contents-row">
                  <td colspan="6" class="pa-2 bg-surface">
                    <div class="bag-inner-box pa-2 rounded border" @dragover.prevent @drop.stop="onDropIntoBag($event, t)">
                      <div class="d-flex justify-space-between align-center mb-1">
                        <span class="text-caption font-weight-bold text-secondary">
                          Contents of {{ t.name || 'Bag' }} (Excluded from personal Enc)
                        </span>
                        <v-btn size="x-small" color="secondary" variant="text" prepend-icon="mdi-plus" @click="addItemToBag(t)">
                          Add item to bag
                        </v-btn>
                      </div>

                      <div v-if="!t.containedTrappings || t.containedTrappings.length === 0" class="text-caption text-medium-emphasis font-italic py-2 text-center">
                        Bag is empty. Drag items here or click "+ Add item to bag".
                      </div>

                      <v-table v-else density="compact" class="bg-transparent text-caption">
                        <tbody>
                          <tr
                            v-for="(sub, subIdx) in t.containedTrappings"
                            :key="sub.id || subIdx"
                            draggable="true"
                            class="sub-item-row"
                            @dragstart="onDragStart($event, sub, t.id || null)"
                          >
                            <td style="width: 24px;" class="drag-handle text-center">
                              <v-icon icon="mdi-drag-vertical" size="small" color="medium-emphasis" />
                            </td>
                            <td>
                              <v-text-field v-model="sub.name" variant="plain" density="compact" hide-details placeholder="Item Name" />
                            </td>
                            <td style="width: 50px;">
                              <v-text-field v-model.number="sub.enc" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="0" />
                            </td>
                            <td style="width: 50px;">
                              <v-text-field v-model.number="sub.qty" type="number" variant="plain" density="compact" hide-details class="text-center" placeholder="1" />
                            </td>
                            <td style="width: 40px;" class="text-center">
                              <v-checkbox-btn v-model="sub.worn" density="compact" hide-details color="primary" />
                            </td>
                            <td style="width: 36px;" class="text-right">
                              <DeleteRowBtn @delete="removeItemFromBag(t, subIdx)" />
                            </td>
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </v-table>
        </div>

        <!-- Wealth moved down with bigger font and coins enc count -->
        <div class="pa-3 bg-surface-variant rounded-lg border">
          <div class="d-flex justify-space-between align-center mb-2">
            <span class="text-caption font-weight-bold text-primary">Money / Wealth</span>
            <span class="text-caption text-medium-emphasis">
              Total: {{ totalCoins }} coins (+{{ coinsEnc }} Enc, 1 per 200 coins)
            </span>
          </div>
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
.gap-1 { gap: 4px; }

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

.drag-handle {
  cursor: grab;
}

.trapping-item-row:hover,
.sub-item-row:hover {
  background: rgba(var(--v-theme-surface-variant), 0.2);
}

.bag-row {
  border-left: 3px solid rgb(var(--v-theme-secondary));
}

.bag-inner-box {
  background: rgba(var(--v-theme-surface-variant), 0.2);
}

.cursor-pointer {
  cursor: pointer;
}
</style>
