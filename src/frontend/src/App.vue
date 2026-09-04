<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import AuthDialog from './components/AuthDialog.vue'
import AdminPage from './views/AdminPage.vue'

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'wfrpDark')

function toggleTheme() {
  theme.global.name.value = isDark.value ? 'wfrpLight' : 'wfrpDark'
}

// Page Navigation View State ('sheet' | 'admin')
const currentPage = ref<'sheet' | 'admin'>('sheet')

onMounted(() => {
  if (window.location.hash === '#/admin') {
    currentPage.value = 'admin'
  }
  window.addEventListener('hashchange', () => {
    currentPage.value = window.location.hash === '#/admin' ? 'admin' : 'sheet'
  })
})

function navigateTo(page: 'sheet' | 'admin') {
  currentPage.value = page
  window.location.hash = page === 'admin' ? '#/admin' : '#/'
  showKebabOverlay.value = false
}

// User & Auth State
const currentUser = ref<{ username: string; role: string; token: string } | null>(null)
const showAuthDialog = ref(false)
const showKebabOverlay = ref(false)

function onLoginSuccess(user: { username: string; role: string; token: string }) {
  currentUser.value = user
}

function logout() {
  currentUser.value = null
}

const isAdmin = computed(() => currentUser.value?.role === 'admin')

function printPage() {
  showKebabOverlay.value = false
  window.print()
}

// WFRP 5e Character Sheet Reactive Model (Faithful to PDF)
const character = ref({
  name: 'Gottfried von Altdorf',
  species: 'Human (Reiklander)',
  appearance: 'Tall, slender, analytical eyes, scholar robes',
  class: 'Academic',
  career: 'Wizard',
  careerLevel: 2,
  careerTier: 2,
  careerPath: 'Apprentice -> Journeyman -> Master -> Wizard Lord',
  status: 'Silver 3',
  movement: 4,

  // XP
  xp: { current: 150, spent: 1200 },

  // Fate, Fortune, Resilience, Resolve
  fate: 3,
  fortune: 2,
  resilience: 2,
  resolve: 2,

  // Ambitions
  personalAmbition: 'Acquire an authentic Grimoire of Aqshy from Nuln',
  partyAmbition: 'Cleanse the sewers under Altdorf of mutant corruption',

  // Characteristics
  characteristics: {
    WS: { name: 'Weapon Skill', initial: 32, advances: 10, hint: 'Melee combat proficiency' },
    BS: { name: 'Ballistic Skill', initial: 28, advances: 5, hint: 'Ranged weapons proficiency' },
    S: { name: 'Strength', initial: 30, advances: 0, hint: 'Physical power and melee damage' },
    T: { name: 'Toughness', initial: 35, advances: 5, hint: 'Resilience to damage and disease' },
    I: { name: 'Initiative', initial: 41, advances: 10, hint: 'Reaction speed and perception' },
    Ag: { name: 'Agility', initial: 34, advances: 5, hint: 'Flexibility, dodge and balance' },
    Dex: { name: 'Dexterity', initial: 38, advances: 5, hint: 'Manual precision and fine craft' },
    Int: { name: 'Intelligence', initial: 45, advances: 15, hint: 'Reasoning, memory and lore' },
    WP: { name: 'Willpower', initial: 42, advances: 15, hint: 'Mental fortitude and spellcasting' },
    Fel: { name: 'Fellowship', initial: 31, advances: 5, hint: 'Social influence and charm' },
  } as Record<string, { name: string; initial: number; advances: number; hint: string }>,

  // Wounds & Hardy
  wounds: { current: 12, hardy: 0 },

  // Armour Points per Hit Location (Matching PDF layout)
  armourPoints: {
    head: 0,
    leftArm: 1,
    rightArm: 1,
    body: 1,
    leftLeg: 0,
    rightLeg: 0,
    shield: 0,
  },

  // Corruption, Sin & Mutations
  corruption: { current: 1, max: 8 },
  sin: 0,
  mutations: [
    { name: 'Aethyric Glow', effect: 'Eyes faintly emit red sparks when channelled' },
  ],

  // Wealth
  wealth: { gc: 14, ss: 28, bp: 56 },

  // Encumbrance breakdown (Matching PDF)
  encumbrance: {
    weapons: 2,
    armour: 3,
    trappings: 5,
    other: 1,
  },

  // Basic & Advanced Skills (Exact PDF List)
  skills: [
    { name: 'Art (Painting)', characteristic: 'Dex', adv: 0, type: 'Basic' },
    { name: 'Athletics', characteristic: 'Ag', adv: 5, type: 'Basic' },
    { name: 'Bribery', characteristic: 'Fel', adv: 0, type: 'Basic' },
    { name: 'Charm', characteristic: 'Fel', adv: 5, type: 'Basic' },
    { name: 'Charm Animal', characteristic: 'WP', adv: 0, type: 'Basic' },
    { name: 'Climb', characteristic: 'S', adv: 0, type: 'Basic' },
    { name: 'Consume Alcohol', characteristic: 'T', adv: 5, type: 'Basic' },
    { name: 'Cool', characteristic: 'WP', adv: 8, type: 'Basic' },
    { name: 'Dodge', characteristic: 'Ag', adv: 5, type: 'Basic' },
    { name: 'Drive', characteristic: 'T', adv: 0, type: 'Basic' },
    { name: 'Endurance', characteristic: 'T', adv: 5, type: 'Basic' },
    { name: 'Entertain (Storytelling)', characteristic: 'Fel', adv: 0, type: 'Basic' },
    { name: 'Gamble', characteristic: 'Int', adv: 0, type: 'Basic' },
    { name: 'Gossip', characteristic: 'Fel', adv: 4, type: 'Basic' },
    { name: 'Haggle', characteristic: 'Fel', adv: 2, type: 'Basic' },
    { name: 'Intimidate', characteristic: 'S', adv: 0, type: 'Basic' },
    { name: 'Intuition', characteristic: 'I', adv: 8, type: 'Basic' },
    { name: 'Leadership', characteristic: 'Fel', adv: 0, type: 'Basic' },
    { name: 'Melee (Basic)', characteristic: 'WS', adv: 5, type: 'Basic' },
    { name: 'Melee (Polearm)', characteristic: 'WS', adv: 5, type: 'Basic' },
    { name: 'Navigation', characteristic: 'Int', adv: 0, type: 'Basic' },
    { name: 'Outdoor Survival', characteristic: 'Int', adv: 0, type: 'Basic' },
    { name: 'Perception', characteristic: 'I', adv: 10, type: 'Basic' },
    { name: 'Ride (Horse)', characteristic: 'Ag', adv: 0, type: 'Basic' },
    { name: 'Row', characteristic: 'S', adv: 0, type: 'Basic' },
    { name: 'Stealth (Urban)', characteristic: 'Ag', adv: 4, type: 'Basic' },
    { name: 'Language (Magick)', characteristic: 'Int', adv: 15, type: 'Advanced' },
    { name: 'Channeling (Aqshy)', characteristic: 'WP', adv: 14, type: 'Advanced' },
    { name: 'Lore (Arcane)', characteristic: 'Int', adv: 12, type: 'Advanced' },
  ],

  // Known Languages
  languages: [
    { name: 'Reikspiel (Native)', adv: 0 },
    { name: 'Classical', adv: 10 },
    { name: 'Eltharin', adv: 5 },
  ],

  // Talents
  talents: [
    { name: 'Petty Magic', rank: 1, desc: 'Allows casting of simple petty cantrips.' },
    { name: 'Arcane Magic (Aqshy)', rank: 1, desc: 'Unlocks mastery of the Bright Order lore of Fire.' },
    { name: 'Aethyric Attunement', rank: 2, desc: '+10 to Channeling tests and mitigates Minor Miscasts.' },
    { name: 'Read/Write', rank: 1, desc: 'Can read and write Classical and Reikspiel fluently.' },
  ],

  // Weapons
  weapons: [
    { name: 'Wizard Quarterstaff', group: 'Basic', enc: 2, rangeReach: 'Melee', damage: '+SB+4', qualities: 'Defensive, Two-Handed' },
    { name: 'Dagger of Aqshy', group: 'Basic', enc: 0, rangeReach: 'Melee', damage: '+SB+2', qualities: 'Fast' },
  ],

  // Armour
  armour: [
    { name: 'Leather Jack', locations: 'Body, Arms', enc: 2, ap: 1, qualities: 'Flexible' },
    { name: 'Wizard Robes (Quality)', locations: 'Body', enc: 1, ap: 0, qualities: '+1 Arcane Defense' },
  ],

  // Trappings
  trappings: [
    { name: 'Grimoire of Aqshy', category: 'Trappings', enc: 1, qty: 1, desc: 'Contains arcane lore and spells' },
    { name: 'Healing Poultice', category: 'Consumables', enc: 1, qty: 4, desc: 'Restores +2 Wounds per use' },
    { name: 'Writing Kit & Parchment', category: 'Tools', enc: 1, qty: 1, desc: 'Inks, quills, wax seal' },
  ],

  // Spells and Prayers
  spells: [
    { name: 'Dart (Aqshy)', cn: 0, range: '48 yards', target: '1 Target', duration: 'Instant', description: 'Fires a fiery dart inflicting +6 damage.' },
    { name: 'Fireball', cn: 4, range: '24 yards', target: 'AoE (Willpower Bonus yards)', duration: 'Instant', description: 'Explosive blast inflicting +8 damage and 1 Ablaze condition.' },
    { name: 'Cauterize', cn: 2, range: 'Touch', target: '1 Ally', duration: 'Instant', description: 'Stops bleeding immediately and restores 2 Wounds.' },
    { name: 'Crown of Flame', cn: 6, range: 'You', target: 'Self', duration: '10 Rounds', description: 'Surrounds caster in radiant fire (+2 AP, +10 Fear tests).' },
  ],

  // Notes
  notes: 'Trained at the Colleges of Magic in Altdorf under Master Thaddeus. Seeking rare alchemical reagents in the Reikland to advance to Master Wizard.',
})

// Dynamically calculated WFRP 5e formulas
const getCharCurrent = (code: string) => {
  const c = character.value.characteristics[code]
  if (!c) return 0
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

const getSkillTotal = (skill: { characteristic: string; adv: number }) => {
  const base = getCharCurrent(skill.characteristic)
  return base + (Number(skill.adv) || 0)
}

// Add/Remove Rows
function addSkill() {
  character.value.skills.push({ name: 'New Skill', characteristic: 'Int', adv: 0, type: 'Advanced' })
}

function removeSkill(index: number) {
  character.value.skills.splice(index, 1)
}

function addTalent() {
  character.value.talents.push({ name: 'New Talent', rank: 1, desc: 'Talent description' })
}

function removeTalent(index: number) {
  character.value.talents.splice(index, 1)
}

function addLanguage() {
  character.value.languages.push({ name: 'New Language', adv: 0 })
}

function removeLanguage(index: number) {
  character.value.languages.splice(index, 1)
}

function addWeapon() {
  character.value.weapons.push({ name: 'New Weapon', group: 'Basic', enc: 1, rangeReach: 'Melee', damage: '+SB+3', qualities: '' })
}

function removeWeapon(index: number) {
  character.value.weapons.splice(index, 1)
}

function addArmour() {
  character.value.armour.push({ name: 'New Armor', locations: 'Body', enc: 1, ap: 1, qualities: '' })
}

function removeArmour(index: number) {
  character.value.armour.splice(index, 1)
}

function addTrapping() {
  character.value.trappings.push({ name: 'New Item', category: 'Trappings', enc: 1, qty: 1, desc: '' })
}

function removeTrapping(index: number) {
  character.value.trappings.splice(index, 1)
}

function addSpell() {
  character.value.spells.push({ name: 'New Spell', cn: 2, range: '12 yards', target: '1 Target', duration: 'Instant', description: '' })
}

function removeSpell(index: number) {
  character.value.spells.splice(index, 1)
}

function addMutation() {
  character.value.mutations.push({ name: 'New Mutation', effect: 'Effect details' })
}

function removeMutation(index: number) {
  character.value.mutations.splice(index, 1)
}
</script>

<template>
  <v-app class="wfrp-full-app">
    <!-- Floating Kebab Button -->
    <div class="kebab-fixed-pos">
      <v-tooltip text="System Options & Settings" location="left">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="tooltipProps"
            icon="mdi-dots-vertical"
            color="primary"
            elevation="8"
            size="large"
            class="kebab-fab"
            @click="showKebabOverlay = true"
          />
        </template>
      </v-tooltip>
    </div>

    <!-- Kebab POPOUT OVERLAY MODAL (Full Dialog, No dropdown menu) -->
    <v-dialog v-model="showKebabOverlay" max-width="500" transition="dialog-bottom-transition">
      <v-card color="surface" class="pa-6 rounded-xl elevation-24 border-primary">
        <div class="d-flex justify-space-between align-center mb-4">
          <div class="d-flex align-center">
            <v-icon icon="mdi-shield-cog" color="primary" class="mr-2" size="large" />
            <h2 class="text-h5 font-weight-black mb-0">Character Sheet Menu</h2>
          </div>
          <v-btn icon="mdi-close" variant="text" size="small" @click="showKebabOverlay = false" />
        </div>

        <!-- User Status Banner -->
        <div v-if="currentUser" class="pa-4 bg-surface-variant rounded-lg mb-4 d-flex align-center justify-space-between border">
          <div>
            <div class="text-subtitle-1 font-weight-bold">{{ currentUser.username }}</div>
            <v-chip size="x-small" :color="isAdmin ? 'secondary' : 'info'" class="mt-1" variant="flat">
              {{ currentUser.role.toUpperCase() }}
            </v-chip>
          </div>
          <v-avatar color="primary" size="40">
            <v-icon icon="mdi-account" color="surface" />
          </v-avatar>
        </div>

        <v-list class="bg-transparent pa-0">
          <v-list-item
            :prepend-icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
            :title="isDark ? 'Light Theme Mode' : 'Dark Theme Mode'"
            class="mb-2 rounded-lg bg-surface-variant border"
            @click="toggleTheme"
          />
          <v-list-item
            prepend-icon="mdi-printer"
            title="Print / Export PDF Sheet"
            class="mb-2 rounded-lg bg-surface-variant border"
            @click="printPage"
          />

          <!-- Admin Panel Button (If logged in as Admin) -->
          <v-list-item
            v-if="isAdmin"
            prepend-icon="mdi-shield-crown"
            title="Admin Control Panel"
            subtitle="User & Database Management"
            color="secondary"
            class="mb-2 rounded-lg bg-secondary text-white font-weight-bold"
            @click="navigateTo('admin')"
          />

          <!-- Auth Button -->
          <v-list-item
            v-if="!currentUser"
            prepend-icon="mdi-login"
            title="Sign In / Register"
            class="mt-4 rounded-lg bg-primary text-surface font-weight-bold"
            @click="showKebabOverlay = false; showAuthDialog = true"
          />
          <v-list-item
            v-else
            prepend-icon="mdi-logout"
            title="Log Out"
            class="mt-4 rounded-lg bg-error text-white font-weight-bold"
            @click="logout(); showKebabOverlay = false"
          />
        </v-list>
      </v-card>
    </v-dialog>

    <!-- PAGE CONDITIONAL RENDERING -->
    <template v-if="currentPage === 'sheet'">
      <!-- SNAP-SCROLL CONTAINER: PAGE 1 & PAGE 2 (End-to-End Width, No Centering) -->
      <div class="snap-scroll-wrapper">
      
      <!-- PAGE 1: CHARACTER OVERVIEW & SKILLS (Exact PDF Page 1) -->
      <section id="page-1" class="sheet-page">
        <div class="page-indicator-badge">PAGE 1 - CHARACTER & SKILLS</div>
        
        <!-- HEADER BLOCK -->
        <v-card color="surface" elevation="2" class="mb-4 pa-3 rounded-lg border-gold">
          <v-row dense align="center">
            <v-col cols="12" md="3">
              <v-tooltip text="Character Name" location="top">
                <template #activator="{ props: tProps }">
                  <v-text-field v-bind="tProps" v-model="character.name" label="Name" variant="outlined" density="compact" hide-details class="font-weight-bold" />
                </template>
              </v-tooltip>
            </v-col>
            <v-col cols="12" sm="6" md="2">
              <v-text-field v-model="character.species" label="Species" variant="outlined" density="compact" hide-details />
            </v-col>
            <v-col cols="12" sm="6" md="2">
              <v-text-field v-model="character.appearance" label="Appearance" variant="outlined" density="compact" hide-details />
            </v-col>
            <v-col cols="12" sm="6" md="2">
              <v-text-field v-model="character.class" label="Class" variant="outlined" density="compact" hide-details />
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-text-field v-model="character.career" label="Career" variant="outlined" density="compact" hide-details />
            </v-col>
          </v-row>

          <v-row dense class="mt-2" align="center">
            <v-col cols="12" sm="4" md="2">
              <v-text-field v-model="character.status" label="Status" variant="outlined" density="compact" hide-details />
            </v-col>
            
            <!-- Career Tier Number Input -->
            <v-col cols="12" sm="4" md="2">
              <v-tooltip text="Current Career Tier (e.g. 1, 2, 3, 4)" location="top">
                <template #activator="{ props: tProps }">
                  <v-text-field
                    v-bind="tProps"
                    v-model.number="character.careerTier"
                    label="Career Tier"
                    type="number"
                    variant="outlined"
                    density="compact"
                    hide-details
                    class="no-spinner"
                  />
                </template>
              </v-tooltip>
            </v-col>

            <!-- XP Tracker -->
            <v-col cols="12" md="5">
              <div class="d-flex align-center gap-2">
                <v-tooltip text="Experience points currently available to spend" location="top">
                  <template #activator="{ props: tProps }">
                    <v-text-field v-bind="tProps" v-model.number="character.xp.current" label="XP Current" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                  </template>
                </v-tooltip>
                <v-tooltip text="Total experience points spent on advances" location="top">
                  <template #activator="{ props: tProps }">
                    <v-text-field v-bind="tProps" v-model.number="character.xp.spent" label="XP Spent" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                  </template>
                </v-tooltip>
                <v-tooltip text="Total XP earned over career (Current + Spent)" location="top">
                  <template #activator="{ props: tProps }">
                    <v-text-field v-bind="tProps" :model-value="character.xp.current + character.xp.spent" label="XP Total" type="number" variant="outlined" density="compact" hide-details readonly class="bg-surface-variant no-spinner" />
                  </template>
                </v-tooltip>
              </div>
            </v-col>
          </v-row>
        </v-card>

        <!-- CHARACTERISTICS BAR (WS, BS, S, T, I, Ag, Dex, INT, WP, FEL) -->
        <section class="mb-4">
          <div class="section-title mb-2 d-flex align-center">
            <v-icon icon="mdi-chart-bar" color="primary" class="mr-2" />
            <span>CHARACTERISTICS</span>
          </div>

          <div class="characteristics-row">
            <div
              v-for="(stat, code) in character.characteristics"
              :key="code"
              class="stat-card border rounded-lg pa-2 bg-surface text-center"
            >
              <v-tooltip :text="stat.hint" location="top">
                <template #activator="{ props: tProps }">
                  <div v-bind="tProps" class="stat-code text-h6 font-weight-black text-primary cursor-pointer">{{ code }}</div>
                </template>
              </v-tooltip>
              <div class="text-caption text-truncate text-medium-emphasis mb-1">{{ stat.name }}</div>
              <v-divider class="my-1" />
              <div class="d-flex justify-space-between align-center mb-1 text-caption">
                <span>Init:</span>
                <input v-model.number="stat.initial" type="number" class="clean-num-input" />
              </div>
              <div class="d-flex justify-space-between align-center mb-1 text-caption">
                <span>Adv:</span>
                <input v-model.number="stat.advances" type="number" class="clean-num-input text-success font-weight-bold" />
              </div>
              <v-divider class="my-1" />
              <v-tooltip text="Current Score = Initial + Advances" location="bottom">
                <template #activator="{ props: tProps }">
                  <div v-bind="tProps" class="stat-total text-h5 font-weight-black text-on-surface">
                    {{ getCharCurrent(code) }}
                  </div>
                </template>
              </v-tooltip>
            </div>
          </div>
        </section>

        <!-- MOVEMENT, FATE, RESILIENCE & WOUNDS -->
        <v-row dense class="mb-4">
          <!-- Movement -->
          <v-col cols="12" sm="6" md="3">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-primary text-uppercase mb-2">Movement</div>
              <v-row dense align="center">
                <v-col cols="4">
                  <v-text-field v-model.number="character.movement" label="Move (M)" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                </v-col>
                <v-col cols="4">
                  <v-tooltip text="Walk distance = Move x 2 yards" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" :model-value="computedWalk" label="Walk (Mx2)" readonly variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
                <v-col cols="4">
                  <v-tooltip text="Run distance = Move x 4 yards" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" :model-value="computedRun" label="Run (Mx4)" readonly variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>

          <!-- Fate & Fortune -->
          <v-col cols="12" sm="6" md="3">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-warning text-uppercase mb-2">Fate & Fortune</div>
              <v-row dense>
                <v-col cols="6">
                  <v-tooltip text="Fate points allow surviving fatal blows" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" v-model.number="character.fate" label="Fate" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
                <v-col cols="6">
                  <v-tooltip text="Fortune points reset daily to re-roll tests" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" v-model.number="character.fortune" label="Fortune" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>

          <!-- Resilience & Resolve -->
          <v-col cols="12" sm="6" md="3">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-info text-uppercase mb-2">Resilience & Resolve</div>
              <v-row dense>
                <v-col cols="6">
                  <v-tooltip text="Resilience score prevents corruption and mutation" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" v-model.number="character.resilience" label="Resilience" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
                <v-col cols="6">
                  <v-tooltip text="Resolve score removes conditions in combat" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" v-model.number="character.resolve" label="Resolve" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>

          <!-- Wounds -->
          <v-col cols="12" sm="6" md="3">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-error text-uppercase mb-2">Current Wounds</div>
              <v-row dense align="center">
                <v-col cols="6">
                  <v-text-field v-model.number="character.wounds.current" label="Current" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                </v-col>
                <v-col cols="6">
                  <v-tooltip text="Max Wounds = SB + (TB x 2) + WPB + Hardy" location="top">
                    <template #activator="{ props: tProps }">
                      <v-text-field v-bind="tProps" :model-value="computedMaxWounds" label="Max Wounds" type="number" readonly variant="outlined" density="compact" hide-details class="no-spinner" />
                    </template>
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>
        </v-row>

        <!-- ARMOUR POINTS HIT LOCATIONS & AMBITIONS -->
        <v-row dense class="mb-4">
          <v-col cols="12" lg="7">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-primary text-uppercase mb-2">Armour Points per Hit Location</div>
              <v-row dense class="text-center">
                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Head protection (d100 roll 01-09)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">Head</div>
                        <div class="text-caption text-medium-emphasis">01-09</div>
                        <input v-model.number="character.armourPoints.head" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Left Arm protection (d100 roll 10-24)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">L. Arm</div>
                        <div class="text-caption text-medium-emphasis">10-24</div>
                        <input v-model.number="character.armourPoints.leftArm" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Right Arm protection (d100 roll 25-44)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">R. Arm</div>
                        <div class="text-caption text-medium-emphasis">25-44</div>
                        <input v-model.number="character.armourPoints.rightArm" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Body Torso protection (d100 roll 45-79)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">Body</div>
                        <div class="text-caption text-medium-emphasis">45-79</div>
                        <input v-model.number="character.armourPoints.body" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Left Leg protection (d100 roll 80-89)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">L. Leg</div>
                        <div class="text-caption text-medium-emphasis">80-89</div>
                        <input v-model.number="character.armourPoints.leftLeg" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Right Leg protection (d100 roll 90-00)" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">R. Leg</div>
                        <div class="text-caption text-medium-emphasis">90-00</div>
                        <input v-model.number="character.armourPoints.rightLeg" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>

                <v-col cols="4" sm="3" md="1.7">
                  <v-tooltip text="Shield AP Bonus" location="top">
                    <template #activator="{ props: tProps }">
                      <div v-bind="tProps" class="pa-2 border rounded bg-surface-variant">
                        <div class="text-caption font-weight-bold">Shield</div>
                        <div class="text-caption text-medium-emphasis">Block</div>
                        <input v-model.number="character.armourPoints.shield" type="number" class="clean-num-input mt-1" />
                      </div>
                    </template>
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-card>
          </v-col>

          <v-col cols="12" lg="5">
            <v-card color="surface" elevation="2" class="pa-3 rounded-lg h-100 border">
              <div class="text-caption font-weight-bold text-primary text-uppercase mb-2">Ambitions</div>
              <v-text-field v-model="character.personalAmbition" label="Personal Ambition" variant="outlined" density="compact" class="mb-2" hide-details />
              <v-text-field v-model="character.partyAmbition" label="Party Ambition" variant="outlined" density="compact" hide-details />
            </v-card>
          </v-col>
        </v-row>

        <!-- SKILLS TABLE (2 Columns matching PDF) -->
        <section class="mb-4">
          <v-card color="surface" elevation="2" class="rounded-lg border">
            <v-card-title class="pa-3 d-flex justify-space-between align-center text-subtitle-1 font-weight-bold">
              <div class="d-flex align-center">
                <v-icon icon="mdi-book-open-variant" color="primary" class="mr-2" />
                <span>SKILLS LIST (Basic & Advanced)</span>
              </div>
              <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addSkill" />
            </v-card-title>

            <v-card-text class="pa-3 pt-0">
              <v-row dense>
                <v-col
                  v-for="(skill, idx) in character.skills"
                  :key="idx"
                  cols="12"
                  md="6"
                >
                  <div class="d-flex align-center gap-1 border-bottom py-1">
                    <input v-model="skill.name" type="text" class="inline-text-input font-weight-bold flex-grow-1" placeholder="Skill Name" />
                    <select v-model="skill.characteristic" class="inline-select">
                      <option v-for="c in Object.keys(character.characteristics)" :key="c" :value="c">{{ c }}</option>
                    </select>
                    <input v-model.number="skill.adv" type="number" class="clean-num-input text-success font-weight-bold" placeholder="Adv" />
                    <v-tooltip text="Skill Total = Characteristic Current + Advances" location="right">
                      <template #activator="{ props: tProps }">
                        <div v-bind="tProps" class="skill-total-tag">{{ getSkillTotal(skill) }}</div>
                      </template>
                    </v-tooltip>
                    <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeSkill(idx)" />
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </section>

        <!-- TALENTS & LANGUAGES -->
        <v-row dense>
          <!-- Talents -->
          <v-col cols="12" md="8">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="d-flex justify-space-between align-center mb-2">
                <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                  <v-icon icon="mdi-star" color="primary" class="mr-2" />
                  TALENTS
                </div>
                <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addTalent" />
              </div>

              <div v-for="(talent, idx) in character.talents" :key="idx" class="border rounded pa-2 mb-2 bg-surface-variant">
                <v-row dense align="center">
                  <v-col cols="7">
                    <input v-model="talent.name" type="text" class="inline-text-input font-weight-bold" placeholder="Talent Name" />
                  </v-col>
                  <v-col cols="3">
                    <input v-model.number="talent.rank" type="number" class="clean-num-input" placeholder="Rank" />
                  </v-col>
                  <v-col cols="2" class="text-right">
                    <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeTalent(idx)" />
                  </v-col>
                  <v-col cols="12">
                    <input v-model="talent.desc" type="text" class="inline-text-input text-caption text-medium-emphasis" placeholder="Description & Effects" />
                  </v-col>
                </v-row>
              </div>
            </v-card>
          </v-col>

          <!-- Known Languages -->
          <v-col cols="12" md="4">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="d-flex justify-space-between align-center mb-2">
                <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                  <v-icon icon="mdi-translate" color="primary" class="mr-2" />
                  KNOWN LANGUAGES
                </div>
                <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addLanguage" />
              </div>

              <div v-for="(lang, idx) in character.languages" :key="idx" class="d-flex align-center gap-1 border-bottom py-1">
                <input v-model="lang.name" type="text" class="inline-text-input font-weight-bold flex-grow-1" placeholder="Language Name" />
                <input v-model.number="lang.adv" type="number" class="clean-num-input text-success" placeholder="Adv" />
                <div class="skill-total-tag">{{ getCharCurrent('Int') + (Number(lang.adv) || 0) }}</div>
                <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeLanguage(idx)" />
              </div>
            </v-card>
          </v-col>
        </v-row>
      </section>

      <!-- PAGE 2: WEAPONS, ARMOUR, TRAPPINGS & MAGIC (Exact PDF Page 2) -->
      <section id="page-2" class="sheet-page">
        <div class="page-indicator-badge">PAGE 2 - WEAPONS, INVENTORY & MAGIC</div>

        <!-- WEAPONS & ARMOUR -->
        <v-row dense class="mb-4">
          <!-- Weapons -->
          <v-col cols="12" lg="6">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="d-flex justify-space-between align-center mb-2">
                <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                  <v-icon icon="mdi-sword-cross" color="primary" class="mr-2" />
                  WEAPONS
                </div>
                <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addWeapon" />
              </div>

              <div v-for="(w, idx) in character.weapons" :key="idx" class="border rounded pa-2 mb-2 bg-surface-variant">
                <v-row dense align="center">
                  <v-col cols="4">
                    <input v-model="w.name" type="text" class="inline-text-input font-weight-bold" placeholder="Weapon Name" />
                  </v-col>
                  <v-col cols="2">
                    <input v-model="w.group" type="text" class="inline-text-input" placeholder="Group" />
                  </v-col>
                  <v-col cols="2">
                    <input v-model.number="w.enc" type="number" class="clean-num-input" placeholder="Enc" />
                  </v-col>
                  <v-col cols="3">
                    <input v-model="w.damage" type="text" class="inline-text-input text-error font-weight-bold" placeholder="Damage" />
                  </v-col>
                  <v-col cols="1" class="text-right">
                    <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeWeapon(idx)" />
                  </v-col>
                  <v-col cols="6">
                    <input v-model="w.rangeReach" type="text" class="inline-text-input text-caption" placeholder="Range / Reach" />
                  </v-col>
                  <v-col cols="6">
                    <input v-model="w.qualities" type="text" class="inline-text-input text-caption" placeholder="Qualities" />
                  </v-col>
                </v-row>
              </div>
            </v-card>
          </v-col>

          <!-- Armour -->
          <v-col cols="12" lg="6">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="d-flex justify-space-between align-center mb-2">
                <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                  <v-icon icon="mdi-shield" color="primary" class="mr-2" />
                  ARMOUR
                </div>
                <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addArmour" />
              </div>

              <div v-for="(a, idx) in character.armour" :key="idx" class="border rounded pa-2 mb-2 bg-surface-variant">
                <v-row dense align="center">
                  <v-col cols="4">
                    <input v-model="a.name" type="text" class="inline-text-input font-weight-bold" placeholder="Armour Name" />
                  </v-col>
                  <v-col cols="3">
                    <input v-model="a.locations" type="text" class="inline-text-input" placeholder="Locations" />
                  </v-col>
                  <v-col cols="2">
                    <input v-model.number="a.enc" type="number" class="clean-num-input" placeholder="Enc" />
                  </v-col>
                  <v-col cols="2">
                    <input v-model.number="a.ap" type="number" class="clean-num-input text-success font-weight-bold" placeholder="AP" />
                  </v-col>
                  <v-col cols="1" class="text-right">
                    <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeArmour(idx)" />
                  </v-col>
                  <v-col cols="12">
                    <input v-model="a.qualities" type="text" class="inline-text-input text-caption" placeholder="Qualities" />
                  </v-col>
                </v-row>
              </div>
            </v-card>
          </v-col>
        </v-row>

        <!-- TRAPPINGS, WEALTH & ENCUMBRANCE -->
        <section class="mb-4">
          <v-card color="surface" elevation="2" class="rounded-lg pa-4 border">
            <div class="d-flex justify-space-between align-center mb-3">
              <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                <v-icon icon="mdi-bag-personal" color="primary" class="mr-2" />
                TRAPPINGS & WEALTH
              </div>
              <div class="d-flex align-center gap-3">
                <v-tooltip text="Max Encumbrance = Strength Bonus + Toughness Bonus" location="top">
                  <template #activator="{ props: tProps }">
                    <span v-bind="tProps" class="text-subtitle-2 font-weight-bold">
                      Encumbrance: {{ computedTotalEnc }} / {{ computedMaxEnc }}
                    </span>
                  </template>
                </v-tooltip>
                <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addTrapping" />
              </div>
            </div>

            <!-- Currency Row (GC, SS, BP/D) -->
            <v-row dense class="mb-3">
              <v-col cols="4">
                <div class="pa-2 border rounded text-center bg-surface-variant">
                  <div class="text-caption font-weight-bold text-warning">Gold Crowns (GC)</div>
                  <input v-model.number="character.wealth.gc" type="number" class="clean-num-input text-h6 text-warning" />
                </div>
              </v-col>
              <v-col cols="4">
                <div class="pa-2 border rounded text-center bg-surface-variant">
                  <div class="text-caption font-weight-bold text-info">Silver Shillings (SS)</div>
                  <input v-model.number="character.wealth.ss" type="number" class="clean-num-input text-h6 text-info" />
                </div>
              </v-col>
              <v-col cols="4">
                <div class="pa-2 border rounded text-center bg-surface-variant">
                  <div class="text-caption font-weight-bold text-secondary">Brass Pennies (D/BP)</div>
                  <input v-model.number="character.wealth.bp" type="number" class="clean-num-input text-h6 text-secondary" />
                </div>
              </v-col>
            </v-row>

            <!-- Trappings Items List -->
            <div v-for="(t, idx) in character.trappings" :key="idx" class="d-flex align-center gap-2 border-bottom py-2">
              <input v-model="t.name" type="text" class="inline-text-input font-weight-bold flex-grow-1" placeholder="Item Name" />
              <input v-model="t.category" type="text" class="inline-text-input" style="width: 110px;" placeholder="Category" />
              <input v-model.number="t.qty" type="number" class="clean-num-input" placeholder="Qty" />
              <input v-model.number="t.enc" type="number" class="clean-num-input" placeholder="Enc" />
              <input v-model="t.desc" type="text" class="inline-text-input text-caption flex-grow-1" placeholder="Description" />
              <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeTrapping(idx)" />
            </div>
          </v-card>
        </section>

        <!-- SPELLS AND PRAYERS -->
        <section class="mb-4">
          <v-card color="surface" elevation="2" class="rounded-lg pa-4 border">
            <div class="d-flex justify-space-between align-center mb-3">
              <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center">
                <v-icon icon="mdi-wizard-hat" color="primary" class="mr-2" />
                SPELLS AND PRAYERS
              </div>
              <v-btn icon="mdi-plus" size="small" color="primary" variant="tonal" @click="addSpell" />
            </div>

            <div v-for="(s, idx) in character.spells" :key="idx" class="border rounded pa-2 mb-2 bg-surface-variant">
              <v-row dense align="center">
                <v-col cols="4">
                  <input v-model="s.name" type="text" class="inline-text-input font-weight-bold text-primary" placeholder="Spell Name" />
                </v-col>
                <v-col cols="2">
                  <v-tooltip text="Casting Number required for test" location="top">
                    <template #activator="{ props: tProps }">
                      <input v-bind="tProps" v-model.number="s.cn" type="number" class="clean-num-input text-warning font-weight-bold" placeholder="CN" />
                    </template>
                  </v-tooltip>
                </v-col>
                <v-col cols="2">
                  <input v-model="s.range" type="text" class="inline-text-input text-caption" placeholder="Range" />
                </v-col>
                <v-col cols="3">
                  <input v-model="s.target" type="text" class="inline-text-input text-caption" placeholder="Target" />
                </v-col>
                <v-col cols="1" class="text-right">
                  <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeSpell(idx)" />
                </v-col>
                <v-col cols="12">
                  <input v-model="s.description" type="text" class="inline-text-input text-caption" placeholder="Spell Description & Effect" />
                </v-col>
              </v-row>
            </div>
          </v-card>
        </section>

        <!-- CORRUPTION, MUTATIONS & NOTES -->
        <v-row dense>
          <v-col cols="12" md="6">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="d-flex justify-space-between align-center mb-2">
                <div class="text-subtitle-1 font-weight-bold text-secondary d-flex align-center">
                  <v-icon icon="mdi-biohazard" color="secondary" class="mr-2" />
                  CORRUPTION & MUTATIONS
                </div>
                <v-btn icon="mdi-plus" size="small" color="secondary" variant="tonal" @click="addMutation" />
              </div>

              <v-row dense class="mb-2">
                <v-col cols="4">
                  <v-text-field v-model.number="character.corruption.current" label="Corruption" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                </v-col>
                <v-col cols="4">
                  <v-text-field v-model.number="character.corruption.max" label="Max (SB+TB)" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                </v-col>
                <v-col cols="4">
                  <v-text-field v-model.number="character.sin" label="Sin Score" type="number" variant="outlined" density="compact" hide-details class="no-spinner" />
                </v-col>
              </v-row>

              <div v-for="(m, idx) in character.mutations" :key="idx" class="d-flex align-center gap-1 border-bottom py-1">
                <input v-model="m.name" type="text" class="inline-text-input font-weight-bold text-secondary" placeholder="Mutation" />
                <input v-model="m.effect" type="text" class="inline-text-input text-caption flex-grow-1" placeholder="Effect" />
                <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="removeMutation(idx)" />
              </div>
            </v-card>
          </v-col>

          <!-- Notes & Journal -->
          <v-col cols="12" md="6">
            <v-card color="surface" elevation="2" class="rounded-lg pa-3 border h-100">
              <div class="text-subtitle-1 font-weight-bold text-primary d-flex align-center mb-2">
                <v-icon icon="mdi-notebook" color="primary" class="mr-2" />
                NOTES & JOURNAL
              </div>
              <v-textarea
                v-model="character.notes"
                variant="outlined"
                rows="5"
                auto-grow
                density="comfortable"
                hide-details
              />
            </v-card>
          </v-col>
        </v-row>
      </section>

    </div>
    </template>

    <!-- STANDALONE ADMIN PAGE VIEW -->
    <template v-else-if="currentPage === 'admin'">
      <AdminPage :current-user="currentUser" @back="navigateTo('sheet')" />
    </template>

    <!-- Dialogs -->
    <AuthDialog v-model="showAuthDialog" @login-success="onLoginSuccess" />
  </v-app>
</template>

<style>
/* Global CSS Rules for full stretch, scroll snap and hiding number spinners */

/* Hide HTML & Chrome/Edge/Safari/Firefox number spinner arrows */
input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type='number'] {
  -moz-appearance: textfield;
  appearance: textfield;
}

.no-spinner input[type='number'] {
  -moz-appearance: textfield;
  appearance: textfield;
}

.wfrp-full-app {
  min-height: 100vh;
  width: 100% !important;
  max-width: 100% !important;
  margin: 0 !important;
  padding: 0 !important;
}

/* Snap Scroll Container between Page 1 and Page 2 */
.snap-scroll-wrapper {
  height: 100vh;
  width: 100%;
  overflow-y: auto;
  scroll-snap-type: y mandatory;
  scroll-behavior: smooth;
}

.sheet-page {
  min-height: 100vh;
  width: 100% !important;
  max-width: 100% !important;
  scroll-snap-align: start;
  scroll-snap-stop: always;
  padding: 16px 24px;
  box-sizing: border-box;
  position: relative;
}

.page-indicator-badge {
  font-size: 0.75rem;
  font-weight: 900;
  letter-spacing: 1px;
  color: #C99700;
  background: rgba(201, 151, 0, 0.12);
  border: 1px solid rgba(201, 151, 0, 0.3);
  padding: 4px 10px;
  border-radius: 4px;
  display: inline-block;
  margin-bottom: 12px;
}

/* Floating Kebab Action Button */
.kebab-fixed-pos {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 3000;
}

.kebab-fab {
  border-radius: 50%;
}

.section-title {
  font-size: 1.15rem;
  font-weight: 900;
  letter-spacing: 1px;
  color: var(--v-theme-primary);
}

/* 10 Characteristics Horizontal Grid */
.characteristics-row {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 8px;
}

@media (max-width: 1200px) {
  .characteristics-row {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (max-width: 600px) {
  .characteristics-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Clean Custom Number Input (No spinner arrows) */
.clean-num-input {
  width: 44px;
  text-align: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 4px;
  color: inherit;
  font-size: 0.9rem;
  padding: 2px 0;
  outline: none;
  transition: border-color 0.2s ease;
}

.clean-num-input:focus {
  border-color: #C99700;
  box-shadow: 0 0 6px rgba(201, 151, 0, 0.4);
}

/* Inline Text Input for tables */
.inline-text-input {
  background: transparent;
  border: none;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.2);
  color: inherit;
  font-size: 0.9rem;
  padding: 2px 4px;
  outline: none;
  width: 100%;
}

.inline-text-input:focus {
  border-bottom-color: #C99700;
}

.inline-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 4px;
  color: inherit;
  font-size: 0.85rem;
  padding: 2px 4px;
  outline: none;
}

.skill-total-tag {
  width: 36px;
  text-align: center;
  font-weight: 900;
  font-size: 0.95rem;
  color: #C99700;
}

.gap-1 { gap: 4px; }
.gap-2 { gap: 8px; }
.gap-3 { gap: 12px; }

.border-bottom {
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.border-gold {
  border: 1px solid rgba(201, 151, 0, 0.3) !important;
}

@media print {
  .kebab-fixed-pos,
  .page-indicator-badge {
    display: none !important;
  }

  .sheet-page {
    page-break-after: always;
    height: auto;
    min-height: auto;
  }
}
</style>
