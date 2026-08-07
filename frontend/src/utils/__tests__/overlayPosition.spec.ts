import { describe, expect, it } from 'vitest'
import { calculateOverlayPosition } from '../overlayPosition'

function rect(left: number, top: number, width = 20, height = 20): DOMRect {
  return {
    left,
    top,
    width,
    height,
    right: left + width,
    bottom: top + height,
  } as DOMRect
}

describe('calculateOverlayPosition', () => {
  it('centers below an anchor when there is enough room', () => {
    expect(calculateOverlayPosition(rect(390, 100), 300, 200, 800, 600)).toEqual({
      left: 250,
      top: 132,
      placement: 'bottom',
    })
  })

  it('places the overlay above an anchor near the bottom', () => {
    expect(calculateOverlayPosition(rect(390, 550), 300, 200, 800, 600)).toEqual({
      left: 250,
      top: 338,
      placement: 'top',
    })
  })

  it('shifts overlays inward at both viewport edges', () => {
    expect(calculateOverlayPosition(rect(0, 100), 300, 100, 800, 600).left).toBe(12)
    expect(calculateOverlayPosition(rect(790, 100), 300, 100, 800, 600).left).toBe(488)
  })
})
