import { describe, expect, it } from 'vitest'
import {
  calculateCssPace,
  calculateCssZones,
  formatPaceRange,
  parseSwimTime,
} from '../criticalSwimSpeed'

describe('critical swim speed', () => {
  it('calculates CSS pace from 200m and 400m times', () => {
    expect(calculateCssPace('3:20', '7:00')).toBe(110)
  })

  it('rejects invalid test times', () => {
    expect(parseSwimTime('3:60')).toBeNull()
    expect(calculateCssPace('7:00', '3:20')).toBeNull()
  })

  it('calculates slower displayed paces from CSS speed percentages', () => {
    const zones = calculateCssZones(110)

    expect(zones).toHaveLength(5)
    expect(zones[2]!.name).toBe('Threshold')
    expect(zones[3]!.name).toBe('Anaerobic / VO2 Max')
    expect(zones[4]!.name).toBe('Sprint')
    expect(formatPaceRange(zones[0]!)).toBe('2:06-2:23 / 100 m')
    expect(formatPaceRange(zones[4]!)).toBe('< 1:39 / 100 m')
  })
})
