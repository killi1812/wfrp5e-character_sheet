/**
 * Validates whether an array is a boolean array of expected length, returning fallback if invalid.
 */
export function ensureBoolArray(arr: unknown, length: number): boolean[] {
  return Array.isArray(arr) && arr.length === length ? arr : Array(length).fill(false)
}
