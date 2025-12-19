/**
 * Detects if the current device is running iOS (iPhone, iPad, iPod)
 * iOS has restrictions on clipboard access that require user gestures,
 * so we need to use a two-step sharing process on these devices.
 */
export function isIOS(): boolean {
    return /iPad|iPhone|iPod/.test(navigator.userAgent) && !(window as unknown as { MSStream?: unknown }).MSStream
}
