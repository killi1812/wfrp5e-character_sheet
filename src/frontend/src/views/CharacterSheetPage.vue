<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  DEFAULT_CHARACTER,
  DEFAULT_BASIC_SKILLS,
  DEFAULT_ADVANCED_SKILLS,
  DEFAULT_LANGUAGES,
  NEW_ITEM_TEMPLATES,
  type CharacterModel,
  type Skill
} from '../constants/placeholders'
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

const computedWalk = computed(() => (Number(character.value.movement) || 0) * 2)
const computedRun = computed(() => (Number(character.value.movement) || 0) * 4)

const computedMaxEnc = computed(() => getCharBonus('S') + getCharBonus('T'))

const computedTotalEnc = computed(() => {
  let total = 0
  character.value.weapons.forEach(w => total += Number(w.enc) || 0)
  character.value.armour.forEach(a => total += Number(a.enc) || 0)
  character.value.trappings.forEach(t => total += (Number(t.enc) || 0) * (Number(t.qty) || 1))
  return total
})

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
      talents: Array.isArray(rawChar.talents) ? rawChar.talents : [],
      weapons: Array.isArray(rawChar.weapons) ? rawChar.weapons : [],
      armour: Array.isArray(rawChar.armour) ? rawChar.armour : [],
      trappings: Array.isArray(rawChar.trappings) ? rawChar.trappings : [],
      spells: Array.isArray(rawChar.spells) ? rawChar.spells : [],
      mutations: Array.isArray(rawChar.mutations) ? rawChar.mutations : [],
      advances2: Array.isArray(rawChar.advances2) && rawChar.advances2.length === 10 ? rawChar.advances2 : Array(10).fill(false),
      advances3: Array.isArray(rawChar.advances3) && rawChar.advances3.length === 10 ? rawChar.advances3 : Array(10).fill(false),
      advances4: Array.isArray(rawChar.advances4) && rawChar.advances4.length === 10 ? rawChar.advances4 : Array(10).fill(false),
    }

    basicSkills.value = Array.isArray(parsed.basicSkills) ? parsed.basicSkills : blank.basicSkills
    advancedSkills.value = Array.isArray(parsed.advancedSkills) ? parsed.advancedSkills : []
    languages.value = Array.isArray(parsed.languages) ? parsed.languages : []

    currentSheetUuid.value = undefined
    snackbarMessage.value = `Imported character: "${character.value.name || 'Unnamed'}"`
    saveSnackbar.value = true
  } catch (err: any) {
    console.error('Import error:', err)
    snackbarMessage.value = `Failed to import JSON: ${err?.message || 'Invalid file format'}`
    saveSnackbar.value = true
  }
}

// Add / Remove Row Handlers for Advanced Skills, Languages & items
function addAdvancedSkill() {
  advancedSkills.value.push(NEW_ITEM_TEMPLATES.advancedSkill())
}

function removeAdvancedSkill(index: number) {
  advancedSkills.value.splice(index, 1)
}

function addLanguage() {
  languages.value.push(NEW_ITEM_TEMPLATES.language())
}

function removeLanguage(index: number) {
  languages.value.splice(index, 1)
}

function addTalent() {
  character.value.talents.push(NEW_ITEM_TEMPLATES.talent())
}

function removeTalent(index: number) {
  character.value.talents.splice(index, 1)
}

function addWeapon() {
  character.value.weapons.push(NEW_ITEM_TEMPLATES.weapon())
}

function removeWeapon(index: number) {
  character.value.weapons.splice(index, 1)
}

function addArmour() {
  character.value.armour.push(NEW_ITEM_TEMPLATES.armour())
}

function removeArmour(index: number) {
  character.value.armour.splice(index, 1)
}

function addTrapping() {
  character.value.trappings.push(NEW_ITEM_TEMPLATES.trapping())
}

function removeTrapping(index: number) {
  character.value.trappings.splice(index, 1)
}

function addSpell() {
  character.value.spells.push(NEW_ITEM_TEMPLATES.spell())
}

function removeSpell(index: number) {
  character.value.spells.splice(index, 1)
}

function addMutation() {
  character.value.mutations.push(NEW_ITEM_TEMPLATES.mutation())
}

function removeMutation(index: number) {
  character.value.mutations.splice(index, 1)
}

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

    <SheetHeaderBlock :character="character" />
    <SheetCharacteristicsBlock :character="character" :get-char-current="getCharCurrent" />
    <SheetVitalsBlock
      :character="character"
      :computed-walk="computedWalk"
      :computed-run="computedRun"
      :computed-max-wounds="computedMaxWounds"
    />
    <SheetSkillsBlock
      :basic-skills="basicSkills"
      :advanced-skills="advancedSkills"
      :get-char-current="getCharCurrent"
      @add-advanced-skill="addAdvancedSkill"
      @remove-advanced-skill="removeAdvancedSkill"
    />
    <SheetTalentsBlock
      :talents="character.talents"
      :languages="languages"
      :get-char-current="getCharCurrent"
      @add-talent="addTalent"
      @remove-talent="removeTalent"
      @add-language="addLanguage"
      @remove-language="removeLanguage"
    />
    <SheetWeaponsBlock
      :weapons="character.weapons"
      :get-char-bonus="getCharBonus"
      @add-weapon="addWeapon"
      @remove-weapon="removeWeapon"
    />
    <SheetArmourTrappingsBlock
      :armour="character.armour"
      :trappings="character.trappings"
      :wealth="character.wealth"
      :computed-total-enc="computedTotalEnc"
      :computed-max-enc="computedMaxEnc"
      @add-armour="addArmour"
      @remove-armour="removeArmour"
      @add-trapping="addTrapping"
      @remove-trapping="removeTrapping"
    />
    <SheetSpellsMutationsBlock
      :character="character"
      @add-spell="addSpell"
      @remove-spell="removeSpell"
      @add-mutation="addMutation"
      @remove-mutation="removeMutation"
    />
    <SheetAmbitionsNotesBlock :character="character" />

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
.gap-2 { gap: 8px; }
</style>
