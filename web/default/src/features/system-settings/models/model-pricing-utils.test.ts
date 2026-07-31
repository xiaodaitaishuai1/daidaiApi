import assert from 'node:assert/strict'
import { describe, test } from 'node:test'
import {
  formatPriceForDisplay,
  formatRatioForStorage,
} from './model-pricing-utils'

describe('model pricing number utilities', () => {
  test('normalizes a persisted floating-point tail for display', () => {
    assert.equal(formatPriceForDisplay(3 * 0.008333333333), '0.025')
  })

  test('preserves full round-trip precision when deriving a ratio', () => {
    assert.equal(formatRatioForStorage(0.025 / 3), '0.008333333333333333')
  })
})
