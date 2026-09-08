// API Service / Placeholder Client for Character Sheets
import {
  type CharacterModel,
  type Skill,
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

    const character: CharacterModel = {
      ...blank.character,
      name: backendData.name || '',
      species: backendData.species || '',
      appearance: backendData.appearance || '',
      class: backendData.class || '',
      career: backendData.career || '',
      status: backendData.status || '',
      movement: backendData.movement || 4,
      xp: {
        current: backendData.xpCurrent || 0,
        spent: backendData.xpSpent || 0,
      },
      fate: backendData.fate || 0,
      fortune: backendData.fortune || 0,
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
      advances3: Array.isArray(backendData.advances3) && backendData.advances3.length === 10 ? backendData.advances3 : Array(10).fill(false),
      advances4: Array.isArray(backendData.advances4) && backendData.advances4.length === 10 ? backendData.advances4 : Array(10).fill(false),
      talents: Array.isArray(backendData.talents)
        ? backendData.talents.map((t: any) => ({
            name: t.name || '',
            desc: t.description || t.desc || '',
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
          }))
        : [],
      armour: Array.isArray(backendData.armour)
        ? backendData.armour.map((a: any) => ({
            name: a.name || '',
            locations: a.locations || '',
            enc: Number(a.enc) || 0,
            ap: Number(a.ap) || 0,
            qualities: a.qualities || '',
          }))
        : [],
      trappings: Array.isArray(backendData.trappings)
        ? backendData.trappings.map((t: any) => ({
            name: t.name || '',
            category: t.category || '',
            enc: Number(t.enc) || 0,
            qty: Number(t.qty) || 1,
            desc: t.description || t.desc || '',
          }))
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
    }

    return {
      character,
      basicSkills: backendData.skills?.filter((s: any) => s.type === 'Basic').map((s: any) => ({
        name: s.name || '',
        characteristic: s.characteristic || 'WS',
        adv: Number(s.adv) || 0,
      })) || blank.basicSkills,
      advancedSkills: backendData.skills?.filter((s: any) => s.type === 'Advanced').map((s: any) => ({
        name: s.name || '',
        characteristic: s.characteristic || 'WS',
        adv: Number(s.adv) || 0,
      })) || [],
      languages: backendData.languages?.map((l: any) => ({
        name: l.name || '',
        characteristic: 'Int',
        adv: Number(l.adv) || 0,
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

    return {
      uuid,
      name: character.name,
      species: character.species,
      appearance: character.appearance,
      class: character.class,
      career: character.career,
      status: character.status,
      movement: Number(character.movement) || 4,
      advances2: character.advances2,
      advances3: character.advances3,
      advances4: character.advances4,
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
      fortune: Number(character.fortune) || 0,
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
      weapons: character.weapons.map((w) => ({
        name: w.name || '',
        group: w.group || '',
        enc: Number(w.enc) || 0,
        rangeReach: w.rangeReach || '',
        damage: w.damage || '',
        qualities: w.qualities || '',
      })),
      armour: character.armour.map((a) => ({
        name: a.name || '',
        locations: a.locations || '',
        enc: Number(a.enc) || 0,
        ap: Number(a.ap) || 0,
        qualities: a.qualities || '',
      })),
      trappings: character.trappings.map((t) => ({
        name: t.name || '',
        category: t.category || '',
        enc: Number(t.enc) || 0,
        description: t.desc || '',
      })),
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
    }
  }
}

export const characterApi = new CharacterApiService()
