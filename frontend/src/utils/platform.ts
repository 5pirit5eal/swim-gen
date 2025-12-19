/**
 * Detects if the current device is running iOS (iPhone, iPad, iPod)
 * iOS has restrictions on clipboard access that require user gestures,
 * so we need to use a two-step sharing process on these devices.
 *
 * Note: iPadOS 13+ reports as Mac in user agent to request desktop sites,
 * so we also check for touch support with Mac platform to detect iPads.
 */
export function isIOS(): boolean {
  const win = window as unknown as { MSStream?: unknown }
  // Standard iOS detection (iPhone, iPod, older iPads)
  if (/iPad|iPhone|iPod/.test(navigator.userAgent) && !win.MSStream) {
    return true
  }
  // iPadOS 13+ detection: reports as Mac but has touch support
  if (/Mac/.test(navigator.userAgent) && navigator.maxTouchPoints > 2) {
    return true
  }
  return false
}
