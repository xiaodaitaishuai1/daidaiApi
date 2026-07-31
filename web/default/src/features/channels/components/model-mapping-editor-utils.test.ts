import assert from 'node:assert/strict'
import { describe, test } from 'node:test'
import { shouldSyncExternalMappingValue } from './model-mapping-editor-utils'

describe('model mapping editor synchronization', () => {
  test('does not re-sync a value just emitted by the editor', () => {
    const emittedValue = '{\n  "gpt-4": "gpt-4.1"\n}'

    assert.equal(
      shouldSyncExternalMappingValue(emittedValue, emittedValue),
      false
    )
  })

  test('syncs a value supplied by a different channel or form reset', () => {
    assert.equal(
      shouldSyncExternalMappingValue('{"gpt-4":"gpt-4.1"}', '{}'),
      true
    )
  })
})
