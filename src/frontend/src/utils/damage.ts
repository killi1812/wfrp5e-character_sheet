/**
 * Parse damage string (e.g. "+SB+4" or "SB+2") and compute total damage if SB is present.
 */
export function formatDamage(damage: string, sb: number): string {
  if (!damage) return ''
  const sbMatch = damage.match(/\+?\s*SB\s*\+?\s*(\d+)/i)
  if (sbMatch) {
    const bonus = Number(sbMatch[1]) || 0
    return `${damage} (${sb + bonus})`
  }
  return damage
}
