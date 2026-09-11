<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  DEFAULT_CHARACTER,
  DEFAULT_BASIC_SKILLS,
  DEFAULT_ADVANCED_SKILLS,
  DEFAULT_LANGUAGES,
  NEW_ITEM_TEMPLATES,
  type CharacterModel,
  type Skill,
  type MountData
} from '../constants/placeholders'
import { ensureBoolArray } from '../utils/arrays'
import { useEncumbrance } from '../composables/useEncumbrance'
import { useListManager } from '../composables/useListManager'
import { characterApi } from '../services/characterApi'

import SheetHeaderBlock from '../components/sheet/SheetHeaderBlock.vue'
import SheetCharacteristicsBlock from '../components/sheet/SheetCharacteristicsBlock.vue'
import SheetVitalsBlock from '../components/sheet/SheetVitalsBlock.vue'
import SheetSkillsBlock from '../components/sheet/SheetSkillsBlock.vue'
import SheetTalentsBlock from '../components/sheet/SheetTalentsBlock.vue'
import SheetWeaponsBlock from '../components/sheet/SheetWeaponsBlock.vue'
import SheetArmourTrappingsBlock from '../components/sheet/SheetArmourTrappingsBlock.vue'
import SheetSpellsMutationsBlock from '../components/sheet/SheetSpellsMutationsBlock.vue'
import SheetAmbitionsNotesBlock from '../components/sheet/SheetAmbitionsNotesBlock.vue'
import SheetMountBlock from '../components/sheet/SheetMountBlock.vue'

const emit = defineEmits<{
  (e: 'openMenu'): void
}>()

// WFRP 5e Character Sheet Reactive Model
const character = ref<CharacterModel>(JSON.parse(JSON.stringify(DEFAULT_CHARACTER)))

// Basic & Advanced Skills & Languages
const basicSkills = ref<Skill[]>(JSON.parse(JSON.stringify(DEFAULT_BASIC_SKILLS)))
const advancedSkills = ref<Skill[]>(JSON.parse(JSON.stringify(DEFAULT_ADVANCED_SKILLS)))
const languages = ref<Skill[]>(JSON.parse(JSON.stringify(DEFAULT_LANGUAGES)))

// API & Mock State
const isSaving = ref(false)
const saveSnackbar = ref(false)
const snackbarMessage = ref('')
const currentSheetUuid = ref<string | undefined>(undefined)

// Snapshot for dirty tracking
const savedSnapshot = ref<string>('')

function createSnapshot(): string {
  return JSON.stringify({
    character: character.value,
    basicSkills: basicSkills.value,
    advancedSkills: advancedSkills.value,
    languages: languages.value,
  })
}

function updateSavedSnapshot() {
  savedSnapshot.value = createSnapshot()
}

const isDirty = computed(() => {
  if (!savedSnapshot.value) return false
  return createSnapshot() !== savedSnapshot.value
})

onMounted(async () => {
  // Restore guest character sheet from browser localStorage if unauthenticated
  if (!characterApi.isAuthenticated()) {
    const guestData = characterApi.loadGuestSheet()
    if (guestData) {
      character.value = guestData.character
      basicSkills.value = guestData.basicSkills
      advancedSkills.value = guestData.advancedSkills
      languages.value = guestData.languages
    }
  } else {
    await loadUserCharacter()
  }
  updateSavedSnapshot()
})

// Formulas
const getCharCurrent = (code: string) => {
  if (!code) return 0
  const normalizedCode = Object.keys(character.value.characteristics).find(
    k => k.toLowerCase() === code.trim().toLowerCase()
  )
  if (!normalizedCode) return 0
  const c = character.value.characteristics[normalizedCode]
  return (Number(c.initial) || 0) + (Number(c.advances) || 0)
}

const getCharBonus = (code: string) => {
  return Math.floor(getCharCurrent(code) / 10)
}

const computedMaxWounds = computed(() => {
  const sb = getCharBonus('S')
  const tb = getCharBonus('T')
  const wpb = getCharBonus('WP')
  const hardy = Number(character.value.wounds.hardy) || 0
  return sb + (tb * 2) + wpb + hardy
})

const {
  computedMaxEnc,
  encBreakdown,
  computedTotalEnc,
  movementPenalty,
  computedWalk,
  computedRun,
  agilityPenalty,
  travelFatigue,
} = useEncumbrance(character, getCharBonus, getCharCurrent)

// API / Mock Actions
async function saveSheet() {
  isSaving.value = true
  try {
    const res = await characterApi.saveCharacter({
      character: character.value,
      basicSkills: basicSkills.value,
      advancedSkills: advancedSkills.value,
      languages: languages.value,
    }, currentSheetUuid.value)

    if (res.uuid && res.uuid !== 'guest') {
      currentSheetUuid.value = res.uuid
    }

    if (res.mode === 'cloud') {
      snackbarMessage.value = 'Character sheet saved to cloud account!'
    } else {
      snackbarMessage.value = 'Character sheet saved to browser storage!'
    }
    saveSnackbar.value = true
    updateSavedSnapshot()
  } catch {
    snackbarMessage.value = 'Error saving sheet'
    saveSnackbar.value = true
  } finally {
    isSaving.value = false
  }
}

function loadMockData() {
  const mock = characterApi.getMockData()
  character.value = mock.character
  basicSkills.value = mock.basicSkills
  advancedSkills.value = mock.advancedSkills
  languages.value = mock.languages
  currentSheetUuid.value = 'mock-gottfried'
  snackbarMessage.value = 'Loaded demo character (Gottfried von Altdorf)!'
  saveSnackbar.value = true
  updateSavedSnapshot()
}

function newBlankSheet() {
  const blank = characterApi.getBlankData()
  character.value = blank.character
  basicSkills.value = blank.basicSkills
  advancedSkills.value = blank.advancedSkills
  languages.value = blank.languages
  currentSheetUuid.value = undefined
  if (!characterApi.isAuthenticated()) {
    characterApi.clearGuestSheet()
  }
  snackbarMessage.value = 'New blank character sheet ready.'
  saveSnackbar.value = true
  updateSavedSnapshot()
}

async function loadUserCharacter() {
  if (!characterApi.isAuthenticated()) return
  try {
    const sheets = await characterApi.listCharacters()
    if (sheets.length > 0 && sheets[0].uuid && sheets[0].uuid !== 'mock-gottfried' && sheets[0].uuid !== 'guest') {
      await loadCharacterSheet(sheets[0].uuid)
    }
  } catch (err) {
    console.warn('Failed to load user character:', err)
  }
}

async function loadCharacterSheet(uuid: string) {
  isSaving.value = true
  try {
    const data = await characterApi.getCharacter(uuid)
    character.value = data.character
    basicSkills.value = data.basicSkills
    advancedSkills.value = data.advancedSkills
    languages.value = data.languages
    currentSheetUuid.value = (uuid !== 'guest' && uuid !== 'mock' && uuid !== 'demo' && uuid !== 'mock-gottfried') ? uuid : undefined
    snackbarMessage.value = `Loaded ${data.character.name || 'character sheet'}!`
    saveSnackbar.value = true
    updateSavedSnapshot()
  } catch {
    snackbarMessage.value = 'Failed to load character sheet'
    saveSnackbar.value = true
  } finally {
    isSaving.value = false
  }
}

function exportJson() {
  const payload = {
    version: 1,
    exportedAt: new Date().toISOString(),
    character: character.value,
    basicSkills: basicSkills.value,
    advancedSkills: advancedSkills.value,
    languages: languages.value,
  }

  const dateStr = new Date().toISOString().split('T')[0]
  const rawName = character.value.name?.trim() || 'character'
  const safeName = rawName.replace(/[^\w\s.-]/gi, '').trim().replace(/\s+/g, '-') || 'character'
  const filename = `${dateStr}-${safeName}.json`

  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)

  snackbarMessage.value = `Downloaded ${filename}`
  saveSnackbar.value = true
}

async function importJson(file: File) {
  try {
    const text = await file.text()
    const parsed = JSON.parse(text)
    const blank = characterApi.getBlankData()

    const rawChar = parsed.character || parsed
    if (!rawChar || typeof rawChar !== 'object') {
      throw new Error('Invalid JSON character sheet format')
    }

    const defaultMount = JSON.parse(JSON.stringify(DEFAULT_CHARACTER.mount))

    character.value = {
      ...blank.character,
      ...rawChar,
      characteristics: {
        ...blank.character.characteristics,
        ...(rawChar.characteristics || {}),
      },
      wounds: { ...blank.character.wounds, ...(rawChar.wounds || {}) },
      armourPoints: { ...blank.character.armourPoints, ...(rawChar.armourPoints || {}) },
      wealth: { ...blank.character.wealth, ...(rawChar.wealth || {}) },
      careers: Array.isArray(rawChar.careers)
        ? rawChar.careers
        : (rawChar.career ? [{
          class: rawChar.class || '',
          career: rawChar.career,
          status: rawChar.status || '',
          active: true,
          advances2: Array(10).fill(false),
          advances3: Array(12).fill(false),
          advances4: Array(14).fill(false),
        }] : []),
      talents: Array.isArray(rawChar.talents) ? rawChar.talents : [],
      weapons: Array.isArray(rawChar.weapons) ? rawChar.weapons : [],
      armour: Array.isArray(rawChar.armour) ? rawChar.armour : [],
      trappings: Array.isArray(rawChar.trappings) ? rawChar.trappings : [],
      spells: Array.isArray(rawChar.spells) ? rawChar.spells : [],
      mutations: Array.isArray(rawChar.mutations) ? rawChar.mutations : [],
      mount: rawChar.mount || defaultMount,
      spellsHidden: Boolean(rawChar.spellsHidden),
      mountHidden: rawChar.mountHidden !== undefined ? Boolean(rawChar.mountHidden) : true,
      importantCharacteristics: Array.isArray(rawChar.importantCharacteristics) ? rawChar.importantCharacteristics : [],
      advances2: ensureBoolArray(rawChar.advances2, 10),
      advances3: ensureBoolArray(rawChar.advances3, 12),
      advances4: ensureBoolArray(rawChar.advances4, 14),
    }

    basicSkills.value = Array.isArray(parsed.basicSkills) ? parsed.basicSkills : blank.basicSkills
    advancedSkills.value = Array.isArray(parsed.advancedSkills) ? parsed.advancedSkills : []
    languages.value = Array.isArray(parsed.languages) ? parsed.languages : []

    currentSheetUuid.value = undefined
    snackbarMessage.value = `Imported character: "${character.value.name || 'Unnamed'}"`
    saveSnackbar.value = true
    updateSavedSnapshot()
  } catch (err: any) {
    console.error('Import error:', err)
    snackbarMessage.value = `Failed to import JSON: ${err?.message || 'Invalid file format'}`
    saveSnackbar.value = true
  }
}

// Add / Remove Row Handlers for Advanced Skills, Languages & items
const { add: addAdvancedSkill, remove: removeAdvancedSkill } = useListManager(
  () => advancedSkills.value,
  NEW_ITEM_TEMPLATES.advancedSkill
)
const { add: addLanguage, remove: removeLanguage } = useListManager(
  () => languages.value,
  NEW_ITEM_TEMPLATES.language
)
const { add: addTalent, remove: removeTalent } = useListManager(
  () => character.value.talents,
  NEW_ITEM_TEMPLATES.talent
)
const { add: addWeapon, remove: removeWeapon } = useListManager(
  () => character.value.weapons,
  NEW_ITEM_TEMPLATES.weapon
)
const { add: addArmour, remove: removeArmour } = useListManager(
  () => character.value.armour,
  NEW_ITEM_TEMPLATES.armour
)
const { add: addTrapping, remove: removeTrapping } = useListManager(
  () => character.value.trappings,
  NEW_ITEM_TEMPLATES.trapping
)
const addBag = () => character.value.trappings.push(NEW_ITEM_TEMPLATES.bag())
const { add: addSpell, remove: removeSpell } = useListManager(
  () => character.value.spells,
  NEW_ITEM_TEMPLATES.spell
)
const { add: addMutation, remove: removeMutation } = useListManager(
  () => character.value.mutations,
  NEW_ITEM_TEMPLATES.mutation
)

// Mount Row Handlers
function getMount(): MountData {
  if (!character.value.mount) {
    character.value.mount = JSON.parse(JSON.stringify(DEFAULT_CHARACTER.mount))
  }
  return character.value.mount!
}

const { add: addMountAttack, remove: removeMountAttack } = useListManager(
  () => {
    const m = getMount()
    if (!m.attacks) m.attacks = []
    return m.attacks
  },
  NEW_ITEM_TEMPLATES.mountAttack
)
const { add: addMountSkill, remove: removeMountSkill } = useListManager(
  () => {
    const m = getMount()
    if (!m.skills) m.skills = []
    return m.skills
  },
  NEW_ITEM_TEMPLATES.advancedSkill
)
const { add: addMountTrait, remove: removeMountTrait } = useListManager(
  () => {
    const m = getMount()
    if (!m.traits) m.traits = []
    return m.traits
  },
  NEW_ITEM_TEMPLATES.mountTrait
)
const { add: addMountTrapping, remove: removeMountTrapping } = useListManager(
  () => {
    const m = getMount()
    if (!m.trappings) m.trappings = []
    return m.trappings
  },
  NEW_ITEM_TEMPLATES.trapping
)

defineExpose({
  saveSheet,
  loadUserCharacter,
  loadCharacterSheet,
  newBlankSheet,
  loadMockData,
  exportJson,
  importJson,
  currentSheetUuid,
})
</script>

<template>
  <main class="wfrp-sheet-body pa-3 pa-md-5">

    <SheetHeaderBlock :character="character" :is-dirty="isDirty" @open-menu="emit('openMenu')" />
    <SheetCharacteristicsBlock :character="character" :get-char-current="getCharCurrent"
      :agility-penalty="agilityPenalty" />
    <SheetVitalsBlock :character="character" :computed-walk="computedWalk" :computed-run="computedRun"
      :computed-max-wounds="computedMaxWounds" :movement-penalty="movementPenalty" :travel-fatigue="travelFatigue" />
    <SheetSkillsBlock :basic-skills="basicSkills" :advanced-skills="advancedSkills" :get-char-current="getCharCurrent"
      @add-advanced-skill="addAdvancedSkill" @remove-advanced-skill="removeAdvancedSkill" />
    <SheetTalentsBlock :talents="character.talents" :languages="languages" :get-char-current="getCharCurrent"
      @add-talent="addTalent" @remove-talent="removeTalent" @add-language="addLanguage"
      @remove-language="removeLanguage" />
    <SheetWeaponsBlock :weapons="character.weapons" :get-char-bonus="getCharBonus" @add-weapon="addWeapon"
      @remove-weapon="removeWeapon" />
    <SheetArmourTrappingsBlock :armour="character.armour" :trappings="character.trappings" :wealth="character.wealth"
      :computed-total-enc="computedTotalEnc" :computed-max-enc="computedMaxEnc" :enc-breakdown="encBreakdown"
      @add-armour="addArmour" @remove-armour="removeArmour" @add-trapping="addTrapping" @add-bag="addBag"
      @remove-trapping="removeTrapping" />
    <SheetSpellsMutationsBlock :character="character" @add-spell="addSpell" @remove-spell="removeSpell"
      @add-mutation="addMutation" @remove-mutation="removeMutation" />

    <SheetAmbitionsNotesBlock :character="character" />
    <SheetMountBlock v-if="!character.mountHidden" :mount="getMount()" @add-attack="addMountAttack"
      @remove-attack="removeMountAttack" @add-skill="addMountSkill" @remove-skill="removeMountSkill"
      @add-trait="addMountTrait" @remove-trait="removeMountTrait" @add-trapping="addMountTrapping"
      @remove-trapping="removeMountTrapping" />

    <!-- Feedback Snackbar -->
    <v-snackbar v-model="saveSnackbar" :timeout="2500" color="surface" location="bottom end" class="border">
      <div class="d-flex align-center">
        <v-icon icon="mdi-check-circle" color="success" class="mr-2" />
        <span class="text-body-2 font-weight-bold">{{ snackbarMessage }}</span>
      </div>
    </v-snackbar>
  </main>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>
