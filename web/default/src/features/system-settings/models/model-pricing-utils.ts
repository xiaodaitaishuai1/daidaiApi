function toFiniteNumber(value: unknown): number | null {
  if (value === '' || value === null || value === undefined) return null

  const numberValue = Number(value)
  return Number.isFinite(numberValue) ? numberValue : null
}

/** Format calculated prices without exposing floating-point representation noise. */
export function formatPriceForDisplay(value: unknown): string {
  const numberValue = toFiniteNumber(value)
  if (numberValue === null) return ''

  return Number.parseFloat(numberValue.toPrecision(10)).toString()
}

/** Preserve the shortest decimal that round-trips to the calculated ratio. */
export function formatRatioForStorage(value: unknown): string {
  const numberValue = toFiniteNumber(value)
  return numberValue === null ? '' : numberValue.toString()
}
