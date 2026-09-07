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
    const targetUuid = uuid || (window.crypto?.randomUUID ? window.crypto.randomUUID() : `sheet_${Date.now()}`)
    const response = await fetch(`${this.baseUrl}/${uuid ? uuid : ''}`, {
      method: uuid ? 'PUT' : 'POST',
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
   * List available character sheets
   */
  async listCharacters(): Promise<CharacterSheetSummary[]> {
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
      console.warn('Backend list failed, returning mock summary:', err)
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
   * Adapter: map backend schema to frontend model
   */
  private adaptFromBackend(backendData: any): CharacterSheetPayload {
    const blank = this.getBlankData()
    if (!backendData) return blank

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
      sin: backendData.sin || 0,
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
      advances2: Array(10).fill(false),
      advances3: Array(10).fill(false),
      advances4: Array(10).fill(false),
      talents: backendData.talents || [],
      weapons: backendData.weapons || [],
      armour: backendData.armour || [],
      trappings: backendData.trappings || [],
      spells: backendData.spellsAndPrayers || [],
      mutations: backendData.mutations || [],
    }

    return {
      character,
      basicSkills: backendData.skills?.filter((s: any) => s.type === 'Basic') || blank.basicSkills,
      advancedSkills: backendData.skills?.filter((s: any) => s.type === 'Advanced') || [],
      languages: backendData.languages?.map((l: any) => ({ name: l.name, characteristic: 'Int', adv: l.adv || 0 })) || [],
    }
  }

  /**
   * Adapter: map frontend model to backend schema
   */
  private adaptToBackend(payload: CharacterSheetPayload, uuid: string): any {
    const { character, basicSkills, advancedSkills, languages } = payload

    return {
      uuid,
      name: character.name,
      species: character.species,
      appearance: character.appearance,
      class: character.class,
      career: character.career,
      status: character.status,
      movement: character.movement,
      xpCurrent: character.xp.current,
      xpSpent: character.xp.spent,
      xpTotal: (character.xp.current || 0) + (character.xp.spent || 0),
      fate: character.fate,
      fortune: character.fortune,
      personalAmbition: character.personalAmbition,
      partyAmbition: character.partyAmbition,
      notes: character.notes,
      wounds: {
        current: character.wounds.current,
        hardy: character.wounds.hardy,
      },
      wealth: character.wealth,
      armourPoints: {
        head: character.armourPoints.head,
        primaryArm: character.armourPoints.leftArm,
        secondaryArm: character.armourPoints.rightArm,
        body: character.armourPoints.body,
        primaryLeg: character.armourPoints.leftLeg,
        secondaryLeg: character.armourPoints.rightLeg,
        shield: character.armourPoints.shield,
      },
      corruptionPoints: character.corruption.current,
      mutations: character.mutations,
      talents: character.talents,
      weapons: character.weapons,
      armour: character.armour,
      trappings: character.trappings,
      spellsAndPrayers: character.spells,
      skills: [
        ...basicSkills.map((s) => ({ ...s, type: 'Basic' })),
        ...advancedSkills.map((s) => ({ ...s, type: 'Advanced' })),
      ],
      languages: languages.map((l) => ({ name: l.name, int: 0, adv: l.adv, skill: 0 })),
    }
  }
}

export const characterApi = new CharacterApiService()
