import { computed, type Ref } from 'vue'
import type { CharacterModel } from '../constants/placeholders'

export interface EncumbranceBreakdown {
  weapons: number
  armour: number
  trappings: number
  coins: number
  totalCoins: number
}

/**
 * Encapsulates WFRP encumbrance rules, breakdown, and over-encumbrance penalties.
 */
export function useEncumbrance(
  character: Ref<CharacterModel>,
  getCharBonus: (code: string) => number,
  getCharCurrent: (code: string) => number
) {
  const computedMaxEnc = computed(() => getCharBonus('S') + getCharBonus('T'))

  const encBreakdown = computed<EncumbranceBreakdown>(() => {
    let weapons = 0
    character.value.weapons.forEach((w) => {
      const raw = Number(w.enc) || 0
      weapons += w.worn ? Math.max(0, raw - 1) : raw
    })

    let armour = 0
    character.value.armour.forEach((a) => {
      const raw = Number(a.enc) || 0
      armour += a.worn ? Math.max(0, raw - 1) : raw
    })

    let trappings = 0
    character.value.trappings.forEach((t) => {
      const raw = Number(t.enc) || 0
      const qty = Number(t.qty) || 1
      const perItem = t.worn ? Math.max(0, raw - 1) : raw
      trappings += perItem * qty
    })

    const totalCoins =
      (Number(character.value.wealth.gc) || 0) +
      (Number(character.value.wealth.ss) || 0) +
      (Number(character.value.wealth.bp) || 0)
    const coins = Math.floor(totalCoins / 200)

    return {
      weapons,
      armour,
      trappings,
      coins,
      totalCoins,
    }
  })

  const computedTotalEnc = computed(() => {
    const b = encBreakdown.value
    return b.weapons + b.armour + b.trappings + b.coins
  })

  // Over-encumbrance penalties (visual only, not persisted)
  const movementPenalty = computed(() => {
    const tot = computedTotalEnc.value
    const max = computedMaxEnc.value
    if (tot <= max) return 0
    const baseMove = Number(character.value.movement) || 0
    if (tot > 3 * max) return baseMove
    if (tot > 2 * max) {
      const penalized = Math.max(2, baseMove - 2)
      return Math.max(0, baseMove - penalized)
    }
    const penalized = Math.max(3, baseMove - 1)
    return Math.max(0, baseMove - penalized)
  })

  const effectiveMovement = computed(() => {
    const base = Number(character.value.movement) || 0
    return Math.max(0, base - movementPenalty.value)
  })

  const computedWalk = computed(() => effectiveMovement.value * 2)
  const computedRun = computed(() => effectiveMovement.value * 4)

  const agilityPenalty = computed(() => {
    const tot = computedTotalEnc.value
    const max = computedMaxEnc.value
    if (tot <= max) return 0
    const baseAg = getCharCurrent('Ag')
    if (tot > 3 * max) return baseAg
    if (tot > 2 * max) {
      const penalized = Math.max(10, baseAg - 20)
      return Math.max(0, baseAg - penalized)
    }
    return Math.min(10, baseAg)
  })

  const travelFatigue = computed(() => {
    const tot = computedTotalEnc.value
    const max = computedMaxEnc.value
    if (tot <= max) return 0
    if (tot > 3 * max) return 3
    if (tot > 2 * max) return 2
    return 1
  })

  return {
    computedMaxEnc,
    encBreakdown,
    computedTotalEnc,
    movementPenalty,
    effectiveMovement,
    computedWalk,
    computedRun,
    agilityPenalty,
    travelFatigue,
  }
}
