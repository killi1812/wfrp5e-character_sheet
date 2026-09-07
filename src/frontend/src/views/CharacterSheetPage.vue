<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  DEFAULT_CHARACTER,
  DEFAULT_BASIC_SKILLS,
  DEFAULT_ADVANCED_SKILLS,
  NEW_ITEM_TEMPLATES,
  type CharacterModel,
  type Skill
} from '../constants/placeholders'

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

// Basic & Advanced Skills
const basicSkills = ref<Skill[]>(JSON.parse(JSON.stringify(DEFAULT_BASIC_SKILLS)))
const advancedSkills = ref<Skill[]>(JSON.parse(JSON.stringify(DEFAULT_ADVANCED_SKILLS)))

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

// Add / Remove Row Handlers for Advanced Skills & items
function addAdvancedSkill() {
  advancedSkills.value.push(NEW_ITEM_TEMPLATES.advancedSkill())
}

function removeAdvancedSkill(index: number) {
  advancedSkills.value.splice(index, 1)
}

function addTalent() {
  character.value.talents.push(NEW_ITEM_TEMPLATES.talent())
}

function removeTalent(index: number) {
  character.value.talents.splice(index, 1)
}

function addLanguage() {
  character.value.languages.push(NEW_ITEM_TEMPLATES.language())
}

function removeLanguage(index: number) {
  character.value.languages.splice(index, 1)
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
      :languages="character.languages"
      @add-talent="addTalent"
      @remove-talent="removeTalent"
      @add-language="addLanguage"
      @remove-language="removeLanguage"
    />
    <SheetWeaponsBlock
      :weapons="character.weapons"
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
  </main>
</template>
