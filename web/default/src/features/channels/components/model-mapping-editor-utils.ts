export type ModelMappingRow = {
  id: string
  from: string
  to: string
}

export function shouldSyncExternalMappingValue(
  value: string,
  lastEmittedValue: string | null
): boolean {
  return value !== lastEmittedValue
}

export function parseModelMappingRows(
  value: string,
  createId: () => string
): ModelMappingRow[] | null {
  if (!value.trim()) return []

  try {
    const parsed = JSON.parse(value)
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return null
    }

    return Object.entries(parsed).map(([from, to]) => ({
      id: createId(),
      from,
      to: String(to),
    }))
  } catch {
    return null
  }
}

export function serializeModelMappingRows(rows: ModelMappingRow[]): string {
  const mapping: Record<string, string> = {}

  rows.forEach((row) => {
    if (row.from.trim()) {
      mapping[row.from.trim()] = row.to.trim()
    }
  })

  return Object.keys(mapping).length > 0
    ? JSON.stringify(mapping, null, 2)
    : ''
}
