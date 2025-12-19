import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { isIOS } from '../platform'

describe('platform utils', () => {
  let originalUserAgent: string
  let originalPlatform: string
  let originalMaxTouchPoints: number

  beforeEach(() => {
    // Save original values
    originalUserAgent = navigator.userAgent
    originalPlatform = navigator.platform
    originalMaxTouchPoints = navigator.maxTouchPoints
  })

  afterEach(() => {
    // Restore original values
    Object.defineProperty(navigator, 'userAgent', {
      value: originalUserAgent,
      configurable: true,
    })
    Object.defineProperty(navigator, 'platform', {
      value: originalPlatform,
      configurable: true,
    })
    Object.defineProperty(navigator, 'maxTouchPoints', {
      value: originalMaxTouchPoints,
      configurable: true,
    })
  })

  describe('isIOS()', () => {
    it('detects iPhone', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X)',
        configurable: true,
      })
      expect(isIOS()).toBe(true)
    })

    it('detects iPod', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value: 'Mozilla/5.0 (iPod touch; CPU iPhone OS 12_0 like Mac OS X)',
        configurable: true,
      })
      expect(isIOS()).toBe(true)
    })

    it('detects older iPad with iPad in user agent', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value: 'Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X)',
        configurable: true,
      })
      expect(isIOS()).toBe(true)
    })

    it('detects iPadOS 13+ (MacIntel with touch support)', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value:
          'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)',
        configurable: true,
      })
      Object.defineProperty(navigator, 'platform', {
        value: 'MacIntel',
        configurable: true,
      })
      Object.defineProperty(navigator, 'maxTouchPoints', {
        value: 5,
        configurable: true,
      })
      expect(isIOS()).toBe(true)
    })

    it('does not detect Mac desktop (MacIntel without touch)', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value:
          'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko)',
        configurable: true,
      })
      Object.defineProperty(navigator, 'platform', {
        value: 'MacIntel',
        configurable: true,
      })
      Object.defineProperty(navigator, 'maxTouchPoints', {
        value: 0,
        configurable: true,
      })
      expect(isIOS()).toBe(false)
    })

    it('does not detect Windows', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value:
          'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        configurable: true,
      })
      Object.defineProperty(navigator, 'platform', {
        value: 'Win32',
        configurable: true,
      })
      expect(isIOS()).toBe(false)
    })

    it('does not detect Android', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value:
          'Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36',
        configurable: true,
      })
      Object.defineProperty(navigator, 'platform', {
        value: 'Linux armv8l',
        configurable: true,
      })
      expect(isIOS()).toBe(false)
    })

    it('does not detect Linux desktop', () => {
      Object.defineProperty(navigator, 'userAgent', {
        value:
          'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        configurable: true,
      })
      Object.defineProperty(navigator, 'platform', {
        value: 'Linux x86_64',
        configurable: true,
      })
      expect(isIOS()).toBe(false)
    })
  })
})
