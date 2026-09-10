// Centralized Placeholder Models & Initial Character Constants

export interface Characteristic {
  name: string
  initial: number
  advances: number
  hint: string
}

export interface Skill {
  name: string
  characteristic: string
  adv: number
  important?: boolean
}

export interface Talent {
  name: string
  desc: string
  important?: boolean
}

export interface Weapon {
  name: string
  group: string
  enc: number
  rangeReach: string
  damage: string
  qualities: string
  worn?: boolean
}

export interface ArmourItem {
  name: string
  locations: string
  enc: number
  ap: number
  qualities: string
  worn?: boolean
}

export interface TrappingItem {
  id?: string
  name: string
  category: string
  enc: number
  qty: number
  desc: string
  worn?: boolean
  isBag?: boolean
  bagSize?: number
  containedTrappings?: TrappingItem[]
}

export interface Spell {
  name: string
  cn: number
  range: string
  target: string
  duration: string
  description: string
}

export interface Mutation {
  name: string
  effect: string
}

export interface CareerEntry {
  class: string
  career: string
  status: string
  active: boolean
  advances2: boolean[] // 10 boxes
  advances3: boolean[] // 12 boxes
  advances4: boolean[] // 14 boxes
}

export interface MountAttack {
  name: string
  skillToRoll: string
  displayValue: number
  damage: string
  qualities: string
}

export interface MountTrait {
  name: string
  desc: string
}

export interface MountData {
  name: string
  characteristics: Record<string, Characteristic | null>
  attacks: MountAttack[]
  skills: Skill[]
  traits: MountTrait[]
  trappings: TrappingItem[]
}

export interface CharacterModel {
  name: string
  species: string
  appearance: string
  class: string
  career: string
  careers: CareerEntry[]
  advances2: boolean[]
  advances3: boolean[]
  advances4: boolean[]
  status: string
  movement: number
  xp: { current: number; spent: number }
  fate: number
  fateMax: number
  fortune: number
  fortuneMax: number
  personalAmbition: string
  partyAmbition: string
  characteristics: Record<string, Characteristic>
  wounds: { current: number; hardy: number }
  armourPoints: {
    head: number
    leftArm: number
    rightArm: number
    body: number
    leftLeg: number
    rightLeg: number
    shield: number
  }
  corruption: { current: number; max: number }
  sin: number
  mutations: Mutation[]
  wealth: { gc: number; ss: number; bp: number }
  talents: Talent[]
  weapons: Weapon[]
  armour: ArmourItem[]
  trappings: TrappingItem[]
  spells: Spell[]
  mount?: MountData
  spellsHidden: boolean
  mountHidden: boolean
  importantCharacteristics: string[]
  notes: string
}

export const DEFAULT_CHARACTER: CharacterModel = {
  name: '',
  species: '',
  appearance: '',
  class: '',
  career: '',
  careers: [
    {
      class: '',
      career: '',
      status: '',
      active: true,
      advances2: Array(10).fill(false),
      advances3: Array(12).fill(false),
      advances4: Array(14).fill(false),
    },
  ],
  advances2: Array(10).fill(false),
  advances3: Array(12).fill(false),
  advances4: Array(14).fill(false),
  status: '',
  movement: 4,

  xp: { current: 0, spent: 0 },

  fate: 0,
  fateMax: 0,
  fortune: 0,
  fortuneMax: 0,

  personalAmbition: '',
  partyAmbition: '',

  characteristics: {
    WS: { name: 'Weapon Skill', initial: 0, advances: 0, hint: 'Melee combat proficiency' },
    BS: { name: 'Ballistic Skill', initial: 0, advances: 0, hint: 'Ranged weapons proficiency' },
    S: { name: 'Strength', initial: 0, advances: 0, hint: 'Physical power and melee damage' },
    T: { name: 'Toughness', initial: 0, advances: 0, hint: 'Resilience to damage and disease' },
    I: { name: 'Initiative', initial: 0, advances: 0, hint: 'Reaction speed and perception' },
    Ag: { name: 'Agility', initial: 0, advances: 0, hint: 'Flexibility, dodge and balance' },
    Dex: { name: 'Dexterity', initial: 0, advances: 0, hint: 'Manual precision and fine craft' },
    Int: { name: 'Intelligence', initial: 0, advances: 0, hint: 'Reasoning, memory and lore' },
    WP: { name: 'Willpower', initial: 0, advances: 0, hint: 'Mental fortitude and spellcasting' },
    Fel: { name: 'Fellowship', initial: 0, advances: 0, hint: 'Social influence and charm' },
  },

  wounds: { current: 0, hardy: 0 },

  armourPoints: {
    head: 0,
    leftArm: 0,
    rightArm: 0,
    body: 0,
    leftLeg: 0,
    rightLeg: 0,
    shield: 0,
  },

  corruption: { current: 0, max: 0 },
  sin: 0,
  mutations: [],

  wealth: { gc: 0, ss: 0, bp: 0 },

  talents: [],
  weapons: [],
  armour: [],
  trappings: [],
  spells: [],
  mount: {
    name: '',
    characteristics: {
      WS: { name: 'Weapon Skill', initial: 30, advances: 0, hint: '' },
      BS: null,
      S: { name: 'Strength', initial: 45, advances: 0, hint: '' },
      T: { name: 'Toughness', initial: 45, advances: 0, hint: '' },
      I: { name: 'Initiative', initial: 30, advances: 0, hint: '' },
      Ag: { name: 'Agility', initial: 30, advances: 0, hint: '' },
      Dex: null,
      Int: { name: 'Intelligence', initial: 10, advances: 0, hint: '' },
      WP: { name: 'Willpower', initial: 30, advances: 0, hint: '' },
      Fel: { name: 'Fellowship', initial: 20, advances: 0, hint: '' },
    },
    attacks: [],
    skills: [],
    traits: [],
    trappings: [],
  },
  spellsHidden: false,
  mountHidden: true,
  importantCharacteristics: [],

  notes: '',
}

export const DEFAULT_BASIC_SKILLS: Skill[] = [
  { name: 'Art', characteristic: 'Dex', adv: 0 },
  { name: 'Athletics', characteristic: 'Ag', adv: 0 },
  { name: 'Bribery', characteristic: 'Fel', adv: 0 },
  { name: 'Charm', characteristic: 'Fel', adv: 0 },
  { name: 'Charm Animal', characteristic: 'WP', adv: 0 },
  { name: 'Climb', characteristic: 'S', adv: 0 },
  { name: 'Consume Alcohol', characteristic: 'T', adv: 0 },
  { name: 'Cool', characteristic: 'WP', adv: 0 },
  { name: 'Dodge', characteristic: 'Ag', adv: 0 },
  { name: 'Drive', characteristic: 'T', adv: 0 },
  { name: 'Endurance', characteristic: 'T', adv: 0 },
  { name: 'Entertain', characteristic: 'Fel', adv: 0 },
  { name: 'Gamble', characteristic: 'Int', adv: 0 },
  { name: 'Gossip', characteristic: 'Fel', adv: 0 },
  { name: 'Haggle', characteristic: 'Fel', adv: 0 },
  { name: 'Intimidate', characteristic: 'S', adv: 0 },
  { name: 'Intuition', characteristic: 'I', adv: 0 },
  { name: 'Leadership', characteristic: 'Fel', adv: 0 },
  { name: 'Melee (Basic)', characteristic: 'WS', adv: 0 },
  { name: 'Navigation', characteristic: 'Int', adv: 0 },
  { name: 'Outdoor Survival', characteristic: 'Int', adv: 0 },
  { name: 'Perception', characteristic: 'I', adv: 0 },
  { name: 'Ride', characteristic: 'Ag', adv: 0 },
  { name: 'Row', characteristic: 'S', adv: 0 },
  { name: 'Stealth', characteristic: 'Ag', adv: 0 },
]

export const DEFAULT_ADVANCED_SKILLS: Skill[] = []

export const DEFAULT_LANGUAGES: Skill[] = []

export const NEW_ITEM_TEMPLATES = {
  career: (): CareerEntry => ({
    class: '',
    career: '',
    status: '',
    active: false,
    advances2: Array(10).fill(false),
    advances3: Array(12).fill(false),
    advances4: Array(14).fill(false),
  }),
  advancedSkill: (): Skill => ({ name: '', characteristic: 'Int', adv: 0, important: false }),
  talent: (): Talent => ({ name: '', desc: '', important: false }),
  language: (): Skill => ({ name: '', characteristic: 'Int', adv: 0, important: false }),
  weapon: (): Weapon => ({ name: '', group: '', enc: 0, rangeReach: '', damage: '', qualities: '', worn: false }),
  armour: (): ArmourItem => ({ name: '', locations: '', enc: 0, ap: 0, qualities: '', worn: false }),
  trapping: (): TrappingItem => ({
    id: typeof crypto !== 'undefined' && crypto.randomUUID ? crypto.randomUUID() : `item-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    name: '',
    category: '',
    enc: 0,
    qty: 1,
    desc: '',
    worn: false,
    isBag: false,
    bagSize: 0,
    containedTrappings: [],
  }),
  bag: (): TrappingItem => ({
    id: typeof crypto !== 'undefined' && crypto.randomUUID ? crypto.randomUUID() : `bag-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    name: '',
    category: 'Containers',
    enc: 1,
    qty: 1,
    desc: '',
    worn: false,
    isBag: true,
    bagSize: 5,
    containedTrappings: [],
  }),
  spell: (): Spell => ({ name: '', cn: 0, range: '', target: '', duration: '', description: '' }),
  mutation: (): Mutation => ({ name: '', effect: '' }),
  mountAttack: (): MountAttack => ({ name: '', skillToRoll: 'WS', displayValue: 0, damage: '', qualities: '' }),
  mountTrait: (): MountTrait => ({ name: '', desc: '' }),
}

// Rich Mock Data Set for Demonstration & Testing
export const MOCK_CHARACTER: CharacterModel = {
  name: 'Gottfried von Altdorf',
  species: 'Human (Reiklander)',
  appearance: 'Tall, analytical eyes, crimson scholar robes',
  class: 'Academic',
  career: 'Wizard',
  careers: [
    {
      class: 'Academic',
      career: 'Wizard',
      status: 'Silver 3',
      active: true,
      advances2: [true, true, true, true, false, false, false, false, false, false],
      advances3: [false, false, false, false, false, false, false, false, false, false, false, false],
      advances4: [false, false, false, false, false, false, false, false, false, false, false, false, false, false],
    },
  ],
  advances2: [true, true, true, true, false, false, false, false, false, false],
  advances3: [false, false, false, false, false, false, false, false, false, false, false, false],
  advances4: [false, false, false, false, false, false, false, false, false, false, false, false, false, false],
  status: 'Silver 3',
  movement: 4,

  xp: { current: 150, spent: 1200 },

  fate: 3,
  fateMax: 3,
  fortune: 2,
  fortuneMax: 3,

  personalAmbition: 'Acquire an authentic Grimoire of Aqshy from Nuln',
  partyAmbition: 'Cleanse the sewers under Altdorf of mutant corruption',

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
  },

  wounds: { current: 12, hardy: 0 },

  armourPoints: {
    head: 0,
    leftArm: 1,
    rightArm: 1,
    body: 1,
    leftLeg: 0,
    rightLeg: 0,
    shield: 0,
  },

  corruption: { current: 1, max: 8 },
  sin: 1,
  mutations: [
    { name: 'Aethyric Glow', effect: 'Eyes faintly emit red sparks when channelling magic' },
  ],

  wealth: { gc: 14, ss: 28, bp: 56 },

  talents: [
    { name: 'Petty Magic', desc: 'Allows casting of simple cantrips.' },
    { name: 'Arcane Magic (Aqshy)', desc: 'Mastery of the Bright Order lore of Fire.' },
    { name: 'Aethyric Attunement', desc: '+10 to Channeling tests; mitigates miscasts.' },
    { name: 'Read/Write', desc: 'Can read and write Classical and Reikspiel fluently.' },
  ],

  weapons: [
    { name: 'Wizard Quarterstaff', group: 'Basic', enc: 2, rangeReach: 'Melee', damage: '+SB+4', qualities: 'Defensive, Two-Handed' },
    { name: 'Dagger of Aqshy', group: 'Basic', enc: 0, rangeReach: 'Melee', damage: '+SB+2', qualities: 'Fast' },
  ],

  armour: [
    { name: 'Leather Jack', locations: 'Body, Arms', enc: 2, ap: 1, qualities: 'Flexible' },
    { name: 'Wizard Robes (Quality)', locations: 'Body', enc: 1, ap: 0, qualities: '+1 Arcane Defense' },
  ],

  trappings: [
    { name: 'Grimoire of Aqshy', category: 'Trappings', enc: 1, qty: 1, desc: 'Contains arcane lore and spells' },
    { name: 'Healing Poultice', category: 'Consumables', enc: 1, qty: 4, desc: 'Restores +2 Wounds per use' },
    { name: 'Writing Kit & Parchment', category: 'Tools', enc: 1, qty: 1, desc: 'Inks, quills, wax seal' },
  ],

  spells: [
    { name: 'Dart (Aqshy)', cn: 0, range: '48 yards', target: '1 Target', duration: 'Instant', description: 'Fiery dart inflicting +6 damage.' },
    { name: 'Fireball', cn: 4, range: '24 yards', target: 'AoE (WPB yards)', duration: 'Instant', description: 'Explosive blast inflicting +8 damage and 1 Ablaze condition.' },
    { name: 'Cauterize', cn: 2, range: 'Touch', target: '1 Ally', duration: 'Instant', description: 'Stops bleeding immediately and restores 2 Wounds.' },
    { name: 'Crown of Flame', cn: 6, range: 'You', target: 'Self', duration: '10 Rounds', description: 'Surrounds caster in radiant fire (+2 AP, +10 Fear tests).' },
  ],

  mount: {
    name: 'Altdorf Courser',
    characteristics: {
      WS: { name: 'Weapon Skill', initial: 30, advances: 0, hint: '' },
      BS: null,
      S: { name: 'Strength', initial: 45, advances: 0, hint: '' },
      T: { name: 'Toughness', initial: 45, advances: 0, hint: '' },
      I: { name: 'Initiative', initial: 30, advances: 0, hint: '' },
      Ag: { name: 'Agility', initial: 32, advances: 0, hint: '' },
      Dex: null,
      Int: { name: 'Intelligence', initial: 10, advances: 0, hint: '' },
      WP: { name: 'Willpower', initial: 30, advances: 0, hint: '' },
      Fel: { name: 'Fellowship', initial: 20, advances: 0, hint: '' },
    },
    attacks: [
      { name: 'Bite', skillToRoll: 'WS', displayValue: 30, damage: '+SB+2', qualities: '' },
      { name: 'Kick', skillToRoll: 'WS', displayValue: 30, damage: '+SB+4', qualities: '' },
    ],
    skills: [
      { name: 'Athletics', characteristic: 'Ag', adv: 10 },
      { name: 'Endurance', characteristic: 'T', adv: 15 },
    ],
    traits: [
      { name: 'Size (Large)', desc: 'Larger than humanoid' },
      { name: 'Trained (Mount)', desc: 'Trained to carry a rider in travel and skirmish' },
    ],
    trappings: [
      { id: 'mount-saddle', name: 'Riding Saddle & Bridle', category: 'Gear', enc: 2, qty: 1, desc: 'Quality leather harness' },
      { id: 'mount-saddlebags', name: 'Saddlebags', category: 'Containers', enc: 1, qty: 1, desc: 'Large leather pouches', isBag: true, bagSize: 8, containedTrappings: [
        { id: 'mount-feed', name: 'Oats & Grain (1 week)', category: 'Food', enc: 2, qty: 1, desc: 'Dry fodder' },
      ] },
    ],
  },
  spellsHidden: false,
  mountHidden: true,
  importantCharacteristics: ['WP', 'Int'],

  notes: 'Trained at the Colleges of Magic in Altdorf under Master Thaddeus. Seeking rare alchemical reagents in the Reikland to advance to Master Wizard.',
}

export const MOCK_BASIC_SKILLS: Skill[] = [
  { name: 'Art', characteristic: 'Dex', adv: 0 },
  { name: 'Athletics', characteristic: 'Ag', adv: 5 },
  { name: 'Bribery', characteristic: 'Fel', adv: 0 },
  { name: 'Charm', characteristic: 'Fel', adv: 5 },
  { name: 'Charm Animal', characteristic: 'WP', adv: 0 },
  { name: 'Climb', characteristic: 'S', adv: 0 },
  { name: 'Consume Alcohol', characteristic: 'T', adv: 5 },
  { name: 'Cool', characteristic: 'WP', adv: 8 },
  { name: 'Dodge', characteristic: 'Ag', adv: 5 },
  { name: 'Drive', characteristic: 'T', adv: 0 },
  { name: 'Endurance', characteristic: 'T', adv: 5 },
  { name: 'Entertain', characteristic: 'Fel', adv: 0 },
  { name: 'Gamble', characteristic: 'Int', adv: 0 },
  { name: 'Gossip', characteristic: 'Fel', adv: 4 },
  { name: 'Haggle', characteristic: 'Fel', adv: 2 },
  { name: 'Intimidate', characteristic: 'S', adv: 0 },
  { name: 'Intuition', characteristic: 'I', adv: 8 },
  { name: 'Leadership', characteristic: 'Fel', adv: 0 },
  { name: 'Melee (Basic)', characteristic: 'WS', adv: 5 },
  { name: 'Navigation', characteristic: 'Int', adv: 0 },
  { name: 'Outdoor Survival', characteristic: 'Int', adv: 0 },
  { name: 'Perception', characteristic: 'I', adv: 10 },
  { name: 'Ride', characteristic: 'Ag', adv: 0 },
  { name: 'Row', characteristic: 'S', adv: 0 },
  { name: 'Stealth', characteristic: 'Ag', adv: 4 },
]

export const MOCK_ADVANCED_SKILLS: Skill[] = [
  { name: 'Channeling (Aqshy)', characteristic: 'WP', adv: 14 },
  { name: 'Lore (Arcane)', characteristic: 'Int', adv: 12 },
  { name: 'Language (Magick)', characteristic: 'Int', adv: 15 },
]

export const MOCK_LANGUAGES: Skill[] = [
  { name: 'Reikspiel (Native)', characteristic: 'Int', adv: 0 },
  { name: 'Classical', characteristic: 'Int', adv: 10 },
  { name: 'Eltharin', characteristic: 'Int', adv: 5 },
]
