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
}

export interface Talent {
  name: string
  rank: number
  desc: string
}

export interface Language {
  name: string
  adv: number
}

export interface Weapon {
  name: string
  group: string
  enc: number
  rangeReach: string
  damage: string
  qualities: string
}

export interface ArmourItem {
  name: string
  locations: string
  enc: number
  ap: number
  qualities: string
}

export interface TrappingItem {
  name: string
  category: string
  enc: number
  qty: number
  desc: string
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

export interface CharacterModel {
  name: string
  species: string
  appearance: string
  class: string
  career: string
  careerLevel: number
  careerTier: number
  careerPath: string
  status: string
  movement: number
  xp: { current: number; spent: number }
  fate: number
  fortune: number
  resilience: number
  resolve: number
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
  languages: Language[]
  talents: Talent[]
  weapons: Weapon[]
  armour: ArmourItem[]
  trappings: TrappingItem[]
  spells: Spell[]
  notes: string
}

export const DEFAULT_CHARACTER: CharacterModel = {
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

  xp: { current: 150, spent: 1200 },

  fate: 3,
  fortune: 2,
  resilience: 2,
  resolve: 2,

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
  sin: 0,
  mutations: [
    { name: 'Aethyric Glow', effect: 'Eyes faintly emit red sparks when channelled' },
  ],

  wealth: { gc: 14, ss: 28, bp: 56 },

  languages: [
    { name: 'Reikspiel (Native)', adv: 0 },
    { name: 'Classical', adv: 10 },
    { name: 'Eltharin', adv: 5 },
  ],

  talents: [
    { name: 'Petty Magic', rank: 1, desc: 'Allows casting of simple petty cantrips.' },
    { name: 'Arcane Magic (Aqshy)', rank: 1, desc: 'Unlocks mastery of the Bright Order lore of Fire.' },
    { name: 'Aethyric Attunement', rank: 2, desc: '+10 to Channeling tests and mitigates Minor Miscasts.' },
    { name: 'Read/Write', rank: 1, desc: 'Can read and write Classical and Reikspiel fluently.' },
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
    { name: 'Dart (Aqshy)', cn: 0, range: '48 yards', target: '1 Target', duration: 'Instant', description: 'Fires a fiery dart inflicting +6 damage.' },
    { name: 'Fireball', cn: 4, range: '24 yards', target: 'AoE (Willpower Bonus yards)', duration: 'Instant', description: 'Explosive blast inflicting +8 damage and 1 Ablaze condition.' },
    { name: 'Cauterize', cn: 2, range: 'Touch', target: '1 Ally', duration: 'Instant', description: 'Stops bleeding immediately and restores 2 Wounds.' },
    { name: 'Crown of Flame', cn: 6, range: 'You', target: 'Self', duration: '10 Rounds', description: 'Surrounds caster in radiant fire (+2 AP, +10 Fear tests).' },
  ],

  notes: 'Trained at the Colleges of Magic in Altdorf under Master Thaddeus. Seeking rare alchemical reagents in the Reikland to advance to Master Wizard.',
}

export const DEFAULT_BASIC_SKILLS: Skill[] = [
  { name: 'Art (Painting)', characteristic: 'Dex', adv: 0 },
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
  { name: 'Entertain (Storytelling)', characteristic: 'Fel', adv: 0 },
  { name: 'Gamble', characteristic: 'Int', adv: 0 },
  { name: 'Gossip', characteristic: 'Fel', adv: 4 },
  { name: 'Haggle', characteristic: 'Fel', adv: 2 },
  { name: 'Intimidate', characteristic: 'S', adv: 0 },
  { name: 'Intuition', characteristic: 'I', adv: 8 },
  { name: 'Leadership', characteristic: 'Fel', adv: 0 },
  { name: 'Melee (Basic)', characteristic: 'WS', adv: 5 },
  { name: 'Melee (Polearm)', characteristic: 'WS', adv: 5 },
  { name: 'Navigation', characteristic: 'Int', adv: 0 },
  { name: 'Outdoor Survival', characteristic: 'Int', adv: 0 },
  { name: 'Perception', characteristic: 'I', adv: 10 },
  { name: 'Ride (Horse)', characteristic: 'Ag', adv: 0 },
  { name: 'Row', characteristic: 'S', adv: 0 },
  { name: 'Stealth (Urban)', characteristic: 'Ag', adv: 4 },
]

export const DEFAULT_ADVANCED_SKILLS: Skill[] = [
  { name: 'Language (Magick)', characteristic: 'Int', adv: 15 },
  { name: 'Channeling (Aqshy)', characteristic: 'WP', adv: 14 },
  { name: 'Lore (Arcane)', characteristic: 'Int', adv: 12 },
]

export const NEW_ITEM_TEMPLATES = {
  advancedSkill: (): Skill => ({ name: 'New Advanced Skill', characteristic: 'Int', adv: 0 }),
  talent: (): Talent => ({ name: 'New Talent', rank: 1, desc: 'Talent description' }),
  language: (): Language => ({ name: 'New Language', adv: 0 }),
  weapon: (): Weapon => ({ name: 'New Weapon', group: 'Basic', enc: 1, rangeReach: 'Melee', damage: '+SB+3', qualities: '' }),
  armour: (): ArmourItem => ({ name: 'New Armor', locations: 'Body', enc: 1, ap: 1, qualities: '' }),
  trapping: (): TrappingItem => ({ name: 'New Item', category: 'Trappings', enc: 1, qty: 1, desc: '' }),
  spell: (): Spell => ({ name: 'New Spell', cn: 2, range: '12 yards', target: '1 Target', duration: 'Instant', description: '' }),
  mutation: (): Mutation => ({ name: 'New Mutation', effect: 'Effect details' }),
}
