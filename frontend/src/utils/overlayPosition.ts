export type OverlayPlacement = 'top' | 'bottom'

export interface OverlayPosition {
  left: number
  top: number
  placement: OverlayPlacement
}

export function calculateOverlayPosition(
  anchor: DOMRect,
  overlayWidth: number,
  overlayHeight: number,
  viewportWidth: number,
  viewportHeight: number,
  gap = 12,
  margin = 12,
): OverlayPosition {
  const spaceAbove = anchor.top - gap - margin
  const spaceBelow = viewportHeight - anchor.bottom - gap - margin
  const placement: OverlayPlacement =
    overlayHeight <= spaceBelow || spaceBelow >= spaceAbove ? 'bottom' : 'top'
  const unclampedTop =
    placement === 'bottom' ? anchor.bottom + gap : anchor.top - overlayHeight - gap
  const maxLeft = Math.max(margin, viewportWidth - overlayWidth - margin)
  const maxTop = Math.max(margin, viewportHeight - overlayHeight - margin)

  return {
    left: Math.min(Math.max(anchor.left + anchor.width / 2 - overlayWidth / 2, margin), maxLeft),
    top: Math.min(Math.max(unclampedTop, margin), maxTop),
    placement,
  }
}
