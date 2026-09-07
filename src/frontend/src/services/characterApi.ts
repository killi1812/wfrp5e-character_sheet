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
   * Fetch a character sheet by UUID, or returns cached / mock data
   */
  async getCharacter(uuid?: string): Promise<CharacterSheetPayload> {
    if (uuid && uuid !== 'demo' && uuid !== 'mock') {
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
        console.warn('API fetch failed, falling back to local/mock store:', err)
      }

      // Check localStorage cache
      const cached = localStorage.getItem(`sheet_${uuid}`)
      if (cached) {
        try {
          return JSON.parse(cached)
        } catch {
          // ignore parse error
        }
      }
    }

    // Default to mock data when no UUID or in offline/demo mode
    return this.getMockData()
  }

  /**
   * Save character sheet to backend API (or localStorage fallback)
   */
  async saveCharacter(payload: CharacterSheetPayload, uuid?: string): Promise<{ success: boolean; uuid: string }> {
    const targetUuid = uuid || (window.crypto?.randomUUID ? window.crypto.randomUUID() : `sheet_${Date.now()}`)

    try {
      const response = await fetch(`${this.baseUrl}/${uuid ? uuid : ''}`, {
        method: uuid ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...this.getAuthHeader(),
        },
        body: JSON.stringify(this.adaptToBackend(payload, targetUuid)),
      })

      if (response.ok) {
        const result = await response.json()
        return { success: true, uuid: result.uuid || targetUuid }
      }
    } catch (err) {
      console.warn('Backend save failed, saving to local cache:', err)
    }

    // Fallback: Cache in localStorage
    localStorage.setItem(`sheet_${targetUuid}`, JSON.stringify(payload))
    // Simulate brief API latency
    await new Promise((resolve) => setTimeout(resolve, 300))

    return { success: true, uuid: targetUuid }
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
