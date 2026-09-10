// API Service / Placeholder Client for Character Sheets
import {
  type CharacterModel,
  type Skill,
  type TrappingItem,
  DEFAULT_CHARACTER,
  DEFAULT_BASIC_SKILLS,
  DEFAULT_ADVANCED_SKILLS,
  DEFAULT_LANGUAGES,
  MOCK_CHARACTER,
  MOCK_BASIC_SKILLS,
  MOCK_ADVANCED_SKILLS,
  MOCK_LANGUAGES,
} from '../constants/placeholders'

export interface CharacterSheetPayload {
  character: CharacterModel
  basicSkills: Skill[]
  advancedSkills: Skill[]
  languages: Skill[]
}

export const GUEST_STORAGE_KEY = 'wfrp_guest_character_sheet'

export interface GuestSheetStorage {
  version: 1
  updatedAt: string
  payload: CharacterSheetPayload
}

export interface SaveCharacterResult {
  success: boolean
  uuid: string
  mode: 'cloud' | 'local'
}

export interface CharacterSheetSummary {
  uuid: string
  name: string
  species: string
  career: string
  updatedAt?: string
}

class CharacterApiService {
  private baseUrl = '/api/character-sheets'

  private getAuthHeader(): Record<string, string> {
    const token = localStorage.getItem('auth_token')
    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  /**
   * Check whether user currently holds an active auth token
   */
  isAuthenticated(): boolean {
    return Boolean(localStorage.getItem('auth_token'))
  }

  /**
   * Load guest character sheet from browser localStorage
   */
  loadGuestSheet(): CharacterSheetPayload | null {
    try {
      const raw = localStorage.getItem(GUEST_STORAGE_KEY)
      if (!raw) return null
      const parsed = JSON.parse(raw) as GuestSheetStorage
      if (parsed && parsed.payload) {
        return parsed.payload
      }
    } catch (err) {
      console.warn('Failed to parse guest character sheet from localStorage:', err)
    }
    return null
  }

  /**
   * Save guest character sheet into browser localStorage (Guest mode only)
   */
  saveGuestSheet(payload: CharacterSheetPayload): void {
    const data: GuestSheetStorage = {
      version: 1,
      updatedAt: new Date().toISOString(),
      payload,
    }
    localStorage.setItem(GUEST_STORAGE_KEY, JSON.stringify(data))
  }

  /**
   * Purge guest character sheet from browser localStorage
   */
  clearGuestSheet(): void {
    localStorage.removeItem(GUEST_STORAGE_KEY)
  }

  /**
   * Load mock demo character dataset
   */
  getMockData(): CharacterSheetPayload {
    return {
      character: JSON.parse(JSON.stringify(MOCK_CHARACTER)),
      basicSkills: JSON.parse(JSON.stringify(MOCK_BASIC_SKILLS)),
      advancedSkills: JSON.parse(JSON.stringify(MOCK_ADVANCED_SKILLS)),
      languages: JSON.parse(JSON.stringify(MOCK_LANGUAGES)),
    }
  }

  /**
   * Get clean empty character dataset (for new sheets)
   */
  getBlankData(): CharacterSheetPayload {
    return {
      character: JSON.parse(JSON.stringify(DEFAULT_CHARACTER)),
      basicSkills: JSON.parse(JSON.stringify(DEFAULT_BASIC_SKILLS)),
      advancedSkills: JSON.parse(JSON.stringify(DEFAULT_ADVANCED_SKILLS)),
      languages: JSON.parse(JSON.stringify(DEFAULT_LANGUAGES)),
    }
  }

  /**
   * Fetch a character sheet by UUID, or returns guest / mock / blank data
   */
  async getCharacter(uuid?: string): Promise<CharacterSheetPayload> {
    if (this.isAuthenticated() && uuid && uuid !== 'demo' && uuid !== 'mock') {
      try {
        const response = await fetch(`${this.baseUrl}/${uuid}`, {
          headers: {
            'Content-Type': 'application/json',
            ...this.getAuthHeader(),
          },
        })
        if (response.ok) {
          const data = await response.json()
          return this.adaptFromBackend(data)
        }
      } catch (err) {
        console.warn('API fetch failed:', err)
      }
    }

    // When not logged in, retrieve guest sheet if present
    if (!this.isAuthenticated()) {
      const guest = this.loadGuestSheet()
      if (guest) {
        return guest
      }
    }

    // Default to mock data when explicitly requested
    if (uuid === 'demo' || uuid === 'mock') {
      return this.getMockData()
    }

    return this.getBlankData()
  }

  /**
   * Save character sheet.
   * If logged in: Saves strictly to backend API, and purges localStorage.
   * If not logged in: Saves to browser localStorage (Guest Mode).
   */
  async saveCharacter(payload: CharacterSheetPayload, uuid?: string): Promise<SaveCharacterResult> {
    if (!this.isAuthenticated()) {
      this.saveGuestSheet(payload)
      return { success: true, uuid: 'guest', mode: 'local' }
    }

    // Authenticated user: Save strictly to cloud API, never cache in localStorage
    const isExistingCloudSheet = Boolean(uuid && uuid !== 'guest' && uuid !== 'mock-gottfried')
    const targetUuid = isExistingCloudSheet ? uuid! : this.generateUUID()
    const endpoint = isExistingCloudSheet ? `${this.baseUrl}/${targetUuid}` : `${this.baseUrl}/`
    const method = isExistingCloudSheet ? 'PUT' : 'POST'

    const response = await fetch(endpoint, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...this.getAuthHeader(),
      },
      body: JSON.stringify(this.adaptToBackend(payload, targetUuid)),
    })

    if (!response.ok) {
      throw new Error(`Failed to save character to server: ${response.statusText}`)
    }

    const result = await response.json()
    // Remove local storage data to keep sheet solely on the server
    this.clearGuestSheet()

    return { success: true, uuid: result.uuid || targetUuid, mode: 'cloud' }
  }

  /**
   * Delete character sheet from backend
   */
  async deleteCharacter(uuid: string): Promise<boolean> {
    if (!this.isAuthenticated() || !uuid || uuid === 'guest' || uuid === 'mock-gottfried') {
      return false
    }

    try {
      const response = await fetch(`${this.baseUrl}/${uuid}`, {
        method: 'DELETE',
        headers: {
          ...this.getAuthHeader(),
        },
      })
      return response.ok
    } catch (err) {
      console.warn('Backend delete failed:', err)
      return false
    }
  }

  /**
   * List available character sheets
   */
  async listCharacters(): Promise<CharacterSheetSummary[]> {
    if (this.isAuthenticated()) {
      try {
        const response = await fetch(`${this.baseUrl}/`, {
          headers: {
            'Content-Type': 'application/json',
            ...this.getAuthHeader(),
          },
        })
        if (response.ok) {
          const list = await response.json()
          if (Array.isArray(list)) {
            return list.map((item: any) => ({
              uuid: item.uuid || '',
              name: item.name || 'Unnamed Character',
              species: item.species || '',
              career: item.career || '',
              updatedAt: item.updatedAt || '',
            }))
          }
        }
      } catch (err) {
        console.warn('Backend list failed, returning empty:', err)
      }
      return []
    }

    // Guest mode: check if local guest sheet exists
    const guest = this.loadGuestSheet()
    if (guest && guest.character.name) {
      return [
        {
          uuid: 'guest',
          name: guest.character.name,
          species: guest.character.species,
          career: guest.character.career,
          updatedAt: 'Browser Storage',
        },
      ]
    }

    return [
      {
        uuid: 'mock-gottfried',
        name: MOCK_CHARACTER.name,
        species: MOCK_CHARACTER.species,
        career: MOCK_CHARACTER.career,
        updatedAt: new Date().toISOString(),
      },
    ]
  }

  /**
   * Helper: Generate a valid RFC4122 v4 UUID
   */
  generateUUID(): string {
    if (typeof crypto !== 'undefined' && crypto.randomUUID) {
      return crypto.randomUUID()
    }
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
  }

  /**
   * Adapter: map backend schema to frontend model
   */
  private adaptFromBackend(backendData: any): CharacterSheetPayload {
    const blank = this.getBlankData()
    if (!backendData) return blank

    // Map Characteristics from backend
    const characteristics = { ...blank.character.characteristics }
    if (backendData.characteristics) {
      for (const key of Object.keys(characteristics)) {
        const lower = key.toLowerCase()
        const backendStat = backendData.characteristics[lower]
        if (backendStat) {
          characteristics[key] = {
            ...characteristics[key],
            initial: Number(backendStat.initial) || 0,
            advances: Number(backendStat.advances) || 0,
          }
        }
      }
    }

    const mapTrappingFromBackend = (t: any): any => ({
      id: t.id || (typeof crypto !== 'undefined' && crypto.randomUUID ? crypto.randomUUID() : `item-${Math.random().toString(36).slice(2, 7)}`),
      name: t.name || '',
      category: t.category || '',
      enc: Number(t.enc) || 0,
      qty: Number(t.qty) || 1,
      desc: t.description || t.desc || '',
      worn: Boolean(t.worn),
      isBag: Boolean(t.isBag),
      bagSize: Number(t.bagSize) || 0,
      containedTrappings: Array.isArray(t.containedTrappings) ? t.containedTrappings.map(mapTrappingFromBackend) : [],
    })

    const careers: any[] = Array.isArray(backendData.careers) && backendData.careers.length > 0
      ? backendData.careers.map((c: any) => ({
          class: c.class || '',
          career: c.career || '',
          status: c.status || '',
          active: Boolean(c.active),
          advances2: Array.isArray(c.advances2) && c.advances2.length === 10 ? c.advances2 : Array(10).fill(false),
          advances3: Array.isArray(c.advances3) && c.advances3.length === 12 ? c.advances3 : Array(12).fill(false),
          advances4: Array.isArray(c.advances4) && c.advances4.length === 14 ? c.advances4 : Array(14).fill(false),
        }))
      : [
          {
            class: backendData.class || '',
            career: backendData.career || '',
            status: backendData.status || '',
            active: true,
            advances2: Array.isArray(backendData.advances2) && backendData.advances2.length === 10 ? backendData.advances2 : Array(10).fill(false),
            advances3: Array.isArray(backendData.advances3) && backendData.advances3.length === 12 ? backendData.advances3 : Array(12).fill(false),
            advances4: Array.isArray(backendData.advances4) && backendData.advances4.length === 14 ? backendData.advances4 : Array(14).fill(false),
          },
        ]

    const importantTalentsSet = new Set(backendData.importantTalents || [])
    const importantSkillsSet = new Set(backendData.importantSkills || [])

    // Map Mount
    let mount = blank.character.mount
    if (backendData.mount) {
      const bm = backendData.mount
      const mChars: any = { ...blank.character.mount?.characteristics }
      if (bm.characteristics) {
        for (const code of ['WS', 'BS', 'S', 'T', 'I', 'Ag', 'Dex', 'Int', 'WP', 'Fel']) {
          const lower = code.toLowerCase()
          const bStat = bm.characteristics[lower]
          if (bStat) {
            mChars[code] = {
              name: code,
              initial: Number(bStat.initial) || 0,
              advances: Number(bStat.advances) || 0,
              hint: '',
            }
          } else {
            mChars[code] = null
          }
        }
      }
      mount = {
        name: bm.name || '',
        characteristics: mChars,
        attacks: Array.isArray(bm.attacks)
          ? bm.attacks.map((a: any) => ({
              name: a.name || '',
              skillToRoll: a.skillToRoll || '',
              displayValue: Number(a.displayValue) || 0,
              damage: a.damage || '',
              qualities: a.qualities || '',
            }))
          : [],
        skills: Array.isArray(bm.skills)
          ? bm.skills.map((s: any) => ({
              name: s.name || '',
              characteristic: s.characteristic || 'WS',
              adv: Number(s.adv) || 0,
            }))
          : [],
        traits: Array.isArray(bm.traits)
          ? bm.traits.map((tr: any) => ({
              name: tr.name || '',
              desc: tr.description || tr.desc || '',
            }))
          : [],
        trappings: Array.isArray(bm.trappings) ? bm.trappings.map(mapTrappingFromBackend) : [],
      }
    }

    const character: CharacterModel = {
      ...blank.character,
      name: backendData.name || '',
      species: backendData.species || '',
      appearance: backendData.appearance || '',
      class: backendData.class || careers.find((c) => c.active)?.class || '',
      career: backendData.career || careers.find((c) => c.active)?.career || '',
      careers,
      status: backendData.status || careers.find((c) => c.active)?.status || '',
      movement: backendData.movement || 4,
      xp: {
        current: backendData.xpCurrent || 0,
        spent: backendData.xpSpent || 0,
      },
      fate: backendData.fate || 0,
      fateMax: backendData.fateMax !== undefined ? Number(backendData.fateMax) : (backendData.fate || 0),
      fortune: backendData.fortune || 0,
      fortuneMax: backendData.fortuneMax !== undefined ? Number(backendData.fortuneMax) : (backendData.fortune || 0),
      personalAmbition: backendData.personalAmbition || '',
      partyAmbition: backendData.partyAmbition || '',
      notes: backendData.notes || '',
      sin: Number(backendData.sin) || 0,
      characteristics,
      wounds: {
        current: backendData.wounds?.current || 0,
        hardy: backendData.wounds?.hardy || 0,
      },
      corruption: {
        current: backendData.corruptionPoints || 0,
        max: 8,
      },
      wealth: {
        gc: backendData.wealth?.gc || 0,
        ss: backendData.wealth?.ss || 0,
        bp: backendData.wealth?.bp || 0,
      },
      armourPoints: {
        head: backendData.armourPoints?.head || 0,
        leftArm: backendData.armourPoints?.primaryArm || 0,
        rightArm: backendData.armourPoints?.secondaryArm || 0,
        body: backendData.armourPoints?.body || 0,
        leftLeg: backendData.armourPoints?.primaryLeg || 0,
        rightLeg: backendData.armourPoints?.secondaryLeg || 0,
        shield: backendData.armourPoints?.shield || 0,
      },
      advances2: Array.isArray(backendData.advances2) && backendData.advances2.length === 10 ? backendData.advances2 : Array(10).fill(false),
      advances3: Array.isArray(backendData.advances3) && backendData.advances3.length === 12 ? backendData.advances3 : Array(12).fill(false),
      advances4: Array.isArray(backendData.advances4) && backendData.advances4.length === 14 ? backendData.advances4 : Array(14).fill(false),
      talents: Array.isArray(backendData.talents)
        ? backendData.talents.map((t: any) => ({
            name: t.name || '',
            desc: t.description || t.desc || '',
            important: importantTalentsSet.has(t.name),
          }))
        : [],
      weapons: Array.isArray(backendData.weapons)
        ? backendData.weapons.map((w: any) => ({
            name: w.name || '',
            group: w.group || '',
            enc: Number(w.enc) || 0,
            rangeReach: w.rangeReach || '',
            damage: w.damage || '',
            qualities: w.qualities || '',
            worn: Boolean(w.worn),
          }))
        : [],
      armour: Array.isArray(backendData.armour)
        ? backendData.armour.map((a: any) => ({
            name: a.name || '',
            locations: a.locations || '',
            enc: Number(a.enc) || 0,
            ap: Number(a.ap) || 0,
            qualities: a.qualities || '',
            worn: Boolean(a.worn),
          }))
        : [],
      trappings: Array.isArray(backendData.trappings)
        ? backendData.trappings.map(mapTrappingFromBackend)
        : [],
      spells: Array.isArray(backendData.spellsAndPrayers)
        ? backendData.spellsAndPrayers.map((s: any) => ({
            name: s.name || '',
            cn: Number(s.cn) || 0,
            range: s.range || '',
            target: s.target || '',
            duration: s.duration || '',
            description: s.description || '',
          }))
        : [],
      mutations: Array.isArray(backendData.mutations)
        ? backendData.mutations.map((m: any) => ({
            name: m.name || '',
            effect: m.effect || '',
          }))
        : [],
      mount,
      spellsHidden: Boolean(backendData.spellsHidden),
      mountHidden: backendData.mountHidden !== undefined ? Boolean(backendData.mountHidden) : true,
      importantCharacteristics: Array.isArray(backendData.importantCharacteristics) ? backendData.importantCharacteristics : [],
    }

    return {
      character,
      basicSkills: backendData.skills?.filter((s: any) => s.type === 'Basic').map((s: any) => ({
        name: s.name || '',
        characteristic: s.characteristic || 'WS',
        adv: Number(s.adv) || 0,
        important: importantSkillsSet.has(s.name),
      })) || blank.basicSkills,
      advancedSkills: backendData.skills?.filter((s: any) => s.type === 'Advanced').map((s: any) => ({
        name: s.name || '',
        characteristic: s.characteristic || 'WS',
        adv: Number(s.adv) || 0,
        important: importantSkillsSet.has(s.name),
      })) || [],
      languages: backendData.languages?.map((l: any) => ({
        name: l.name || '',
        characteristic: 'Int',
        adv: Number(l.adv) || 0,
        important: importantSkillsSet.has(l.name),
      })) || [],
    }
  }

  /**
   * Adapter: map frontend model to backend schema
   */
  private adaptToBackend(payload: CharacterSheetPayload, uuid: string): any {
    const { character, basicSkills, advancedSkills, languages } = payload

    const chars = character.characteristics || {}
    const mapStat = (code: string) => {
      const key = Object.keys(chars).find((k) => k.toLowerCase() === code.toLowerCase())
      const c = key ? chars[key] : undefined
      const initial = Number(c?.initial) || 0
      const advances = Number(c?.advances) || 0
      return {
        initial,
        advances,
        current: initial + advances,
      }
    }

    const serializeTrapping = (t: TrappingItem): any => ({
      name: t.name || '',
      category: t.category || '',
      enc: Number(t.enc) || 0,
      description: t.desc || '',
      worn: Boolean(t.worn),
      isBag: Boolean(t.isBag),
      bagSize: Number(t.bagSize) || 0,
      containedTrappings: Array.isArray(t.containedTrappings) ? t.containedTrappings.map(serializeTrapping) : [],
    })

    const activeCareer = character.careers.find((c) => c.active) || character.careers[0]

    return {
      uuid,
      name: character.name,
      species: character.species,
      appearance: character.appearance,
      class: activeCareer?.class || character.class || '',
      career: activeCareer?.career || character.career || '',
      careers: character.careers.map((c) => ({
        class: c.class || '',
        career: c.career || '',
        status: c.status || '',
        active: Boolean(c.active),
        advances2: c.advances2,
        advances3: c.advances3,
        advances4: c.advances4,
      })),
      status: activeCareer?.status || character.status || '',
      movement: Number(character.movement) || 4,
      advances2: activeCareer?.advances2 || character.advances2,
      advances3: activeCareer?.advances3 || character.advances3,
      advances4: activeCareer?.advances4 || character.advances4,
      xpCurrent: Number(character.xp?.current) || 0,
      xpSpent: Number(character.xp?.spent) || 0,
      xpTotal: (Number(character.xp?.current) || 0) + (Number(character.xp?.spent) || 0),
      characteristics: {
        ws: mapStat('WS'),
        bs: mapStat('BS'),
        s: mapStat('S'),
        t: mapStat('T'),
        i: mapStat('I'),
        ag: mapStat('Ag'),
        dex: mapStat('Dex'),
        int: mapStat('Int'),
        wp: mapStat('WP'),
        fel: mapStat('Fel'),
      },
      fate: Number(character.fate) || 0,
      fateMax: Number(character.fateMax) || 0,
      fortune: Number(character.fortune) || 0,
      fortuneMax: Number(character.fortuneMax) || 0,
      personalAmbition: character.personalAmbition || '',
      partyAmbition: character.partyAmbition || '',
      notes: character.notes || '',
      sin: Number(character.sin) || 0,
      corruptionPoints: Number(character.corruption?.current) || 0,
      wounds: {
        current: Number(character.wounds?.current) || 0,
        hardy: Number(character.wounds?.hardy) || 0,
      },
      wealth: {
        gc: Number(character.wealth?.gc) || 0,
        ss: Number(character.wealth?.ss) || 0,
        bp: Number(character.wealth?.bp) || 0,
      },
      armourPoints: {
        head: Number(character.armourPoints?.head) || 0,
        primaryArm: Number(character.armourPoints?.leftArm) || 0,
        secondaryArm: Number(character.armourPoints?.rightArm) || 0,
        body: Number(character.armourPoints?.body) || 0,
        primaryLeg: Number(character.armourPoints?.leftLeg) || 0,
        secondaryLeg: Number(character.armourPoints?.rightLeg) || 0,
        shield: Number(character.armourPoints?.shield) || 0,
      },
      mutations: character.mutations.map((m) => ({ name: m.name || '', effect: m.effect || '' })),
      talents: character.talents.map((t) => ({ name: t.name || '', description: t.desc || '', page: '' })),
      importantTalents: character.talents.filter((t) => t.important).map((t) => t.name),
      importantSkills: [
        ...basicSkills.filter((s) => s.important).map((s) => s.name),
        ...advancedSkills.filter((s) => s.important).map((s) => s.name),
        ...languages.filter((l) => l.important).map((l) => l.name),
      ],
      importantCharacteristics: character.importantCharacteristics || [],
      weapons: character.weapons.map((w) => ({
        name: w.name || '',
        group: w.group || '',
        enc: Number(w.enc) || 0,
        rangeReach: w.rangeReach || '',
        damage: w.damage || '',
        qualities: w.qualities || '',
        worn: Boolean(w.worn),
      })),
      armour: character.armour.map((a) => ({
        name: a.name || '',
        locations: a.locations || '',
        enc: Number(a.enc) || 0,
        ap: Number(a.ap) || 0,
        qualities: a.qualities || '',
        worn: Boolean(a.worn),
      })),
      trappings: character.trappings.map(serializeTrapping),
      spellsAndPrayers: character.spells.map((s) => ({
        name: s.name || '',
        cn: Number(s.cn) || 0,
        range: s.range || '',
        target: s.target || '',
        duration: s.duration || '',
        description: s.description || '',
        sin: 0,
      })),
      skills: [
        ...basicSkills.map((s) => ({ name: s.name, characteristic: s.characteristic, adv: Number(s.adv) || 0, total: 0, type: 'Basic' })),
        ...advancedSkills.map((s) => ({ name: s.name, characteristic: s.characteristic, adv: Number(s.adv) || 0, total: 0, type: 'Advanced' })),
      ],
      languages: languages.map((l) => ({ name: l.name, int: 0, adv: Number(l.adv) || 0, skill: 0 })),
      spellsHidden: Boolean(character.spellsHidden),
      mountHidden: Boolean(character.mountHidden),
      mount: character.mount ? {
        name: character.mount.name || '',
        characteristics: {
          ws: character.mount.characteristics.WS ? { initial: Number(character.mount.characteristics.WS.initial) || 0, current: (Number(character.mount.characteristics.WS.initial) || 0) + (Number(character.mount.characteristics.WS.advances) || 0) } : null,
          bs: character.mount.characteristics.BS ? { initial: Number(character.mount.characteristics.BS.initial) || 0, current: (Number(character.mount.characteristics.BS.initial) || 0) + (Number(character.mount.characteristics.BS.advances) || 0) } : null,
          s: character.mount.characteristics.S ? { initial: Number(character.mount.characteristics.S.initial) || 0, current: (Number(character.mount.characteristics.S.initial) || 0) + (Number(character.mount.characteristics.S.advances) || 0) } : null,
          t: character.mount.characteristics.T ? { initial: Number(character.mount.characteristics.T.initial) || 0, current: (Number(character.mount.characteristics.T.initial) || 0) + (Number(character.mount.characteristics.T.advances) || 0) } : null,
          i: character.mount.characteristics.I ? { initial: Number(character.mount.characteristics.I.initial) || 0, current: (Number(character.mount.characteristics.I.initial) || 0) + (Number(character.mount.characteristics.I.advances) || 0) } : null,
          ag: character.mount.characteristics.Ag ? { initial: Number(character.mount.characteristics.Ag.initial) || 0, current: (Number(character.mount.characteristics.Ag.initial) || 0) + (Number(character.mount.characteristics.Ag.advances) || 0) } : null,
          dex: character.mount.characteristics.Dex ? { initial: Number(character.mount.characteristics.Dex.initial) || 0, current: (Number(character.mount.characteristics.Dex.initial) || 0) + (Number(character.mount.characteristics.Dex.advances) || 0) } : null,
          int: character.mount.characteristics.Int ? { initial: Number(character.mount.characteristics.Int.initial) || 0, current: (Number(character.mount.characteristics.Int.initial) || 0) + (Number(character.mount.characteristics.Int.advances) || 0) } : null,
          wp: character.mount.characteristics.WP ? { initial: Number(character.mount.characteristics.WP.initial) || 0, current: (Number(character.mount.characteristics.WP.initial) || 0) + (Number(character.mount.characteristics.WP.advances) || 0) } : null,
          fel: character.mount.characteristics.Fel ? { initial: Number(character.mount.characteristics.Fel.initial) || 0, current: (Number(character.mount.characteristics.Fel.initial) || 0) + (Number(character.mount.characteristics.Fel.advances) || 0) } : null,
        },
        attacks: character.mount.attacks.map((a) => ({
          name: a.name || '',
          skillToRoll: a.skillToRoll || '',
          displayValue: Number(a.displayValue) || 0,
          damage: a.damage || '',
          qualities: a.qualities || '',
        })),
        skills: character.mount.skills.map((s) => ({
          name: s.name || '',
          characteristic: s.characteristic || 'WS',
          adv: Number(s.adv) || 0,
          total: 0,
          type: 'Advanced',
        })),
        traits: character.mount.traits.map((tr) => ({
          name: tr.name || '',
          description: tr.desc || '',
        })),
        trappings: character.mount.trappings.map(serializeTrapping),
      } : null,
    }
  }
}

export const characterApi = new CharacterApiService()
