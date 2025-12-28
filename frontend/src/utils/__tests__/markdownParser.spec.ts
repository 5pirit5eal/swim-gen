import { describe, it, expect } from 'vitest'
import { extractDrillIdFromUrl, parseContentForDrillLinks, hasDrillLinks } from '../markdownParser'

describe('markdownParser', () => {
  describe('extractDrillIdFromUrl', () => {
    it('extracts drill ID from path-based URL with /drills/ (plural)', () => {
      expect(extractDrillIdFromUrl('/drills/123')).toBe('123')
      expect(extractDrillIdFromUrl('/drills/abc-def')).toBe('abc-def')
      expect(extractDrillIdFromUrl('/drills/drill_with_underscore')).toBe('drill_with_underscore')
    })

    it('extracts drill ID from path-based URL with /drill/ (singular)', () => {
      expect(extractDrillIdFromUrl('/drill/123')).toBe('123')
      expect(extractDrillIdFromUrl('/drill/abc-def')).toBe('abc-def')
      expect(extractDrillIdFromUrl('/drill/drill_with_underscore')).toBe('drill_with_underscore')
    })

    it('extracts drill ID from path without leading slash (drill/)', () => {
      expect(extractDrillIdFromUrl('drill/123')).toBe('123')
      expect(extractDrillIdFromUrl('drill/abc-def')).toBe('abc-def')
    })

    it('extracts drill ID from path without leading slash (drills/)', () => {
      expect(extractDrillIdFromUrl('drills/123')).toBe('123')
      expect(extractDrillIdFromUrl('drills/abc-def')).toBe('abc-def')
    })

    it('extracts drill ID from full URL', () => {
      expect(extractDrillIdFromUrl('https://example.com/drills/123')).toBe('123')
      expect(extractDrillIdFromUrl('http://localhost:3000/drills/test-drill')).toBe('test-drill')
    })

    it('handles URL with query parameters', () => {
      expect(extractDrillIdFromUrl('/drills/123?lang=en')).toBe('123')
    })

    it('handles URL with hash fragment', () => {
      expect(extractDrillIdFromUrl('/drills/123#section')).toBe('123')
    })

    it('returns the URL as ID if no path separators present', () => {
      expect(extractDrillIdFromUrl('simple-id')).toBe('simple-id')
      expect(extractDrillIdFromUrl('123')).toBe('123')
    })

    it('returns .png filename as drill ID', () => {
      expect(extractDrillIdFromUrl('some_drill.png')).toBe('some_drill.png')
      expect(extractDrillIdFromUrl('drill-image.png')).toBe('drill-image.png')
    })

    it('returns null for non-drill URLs', () => {
      expect(extractDrillIdFromUrl('/other/path')).toBeNull()
      expect(extractDrillIdFromUrl('https://example.com/about')).toBeNull()
    })

    it('returns null for URLs with protocol but no drill path', () => {
      expect(extractDrillIdFromUrl('https://example.com')).toBeNull()
    })
  })

  describe('parseContentForDrillLinks', () => {
    it('returns empty array for empty content', () => {
      expect(parseContentForDrillLinks('')).toEqual([])
    })

    it('returns text segment for content without links', () => {
      const result = parseContentForDrillLinks('Plain text content')
      expect(result).toEqual([{ type: 'text', content: 'Plain text content' }])
    })

    it('parses single drill link', () => {
      const content = '[Drill Name](/drills/123)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([{ type: 'drill-link', drillId: '123', text: 'Drill Name' }])
    })

    it('parses drill link with surrounding text', () => {
      const content = 'Start with [Drill Name](/drills/123) and continue'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'text', content: 'Start with ' },
        { type: 'drill-link', drillId: '123', text: 'Drill Name' },
        { type: 'text', content: ' and continue' },
      ])
    })

    it('parses multiple drill links', () => {
      const content = 'Do [First Drill](/drills/first) then [Second Drill](/drills/second)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'text', content: 'Do ' },
        { type: 'drill-link', drillId: 'first', text: 'First Drill' },
        { type: 'text', content: ' then ' },
        { type: 'drill-link', drillId: 'second', text: 'Second Drill' },
      ])
    })

    it('handles non-drill links by returning just the link text', () => {
      const content = 'Check [this link](https://example.com) out'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'text', content: 'Check ' },
        { type: 'text', content: 'this link' },
        { type: 'text', content: ' out' },
      ])
    })

    it('handles mixed drill and non-drill links', () => {
      const content = 'Try [Drill](/drills/abc) or visit [website](https://example.com)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'text', content: 'Try ' },
        { type: 'drill-link', drillId: 'abc', text: 'Drill' },
        { type: 'text', content: ' or visit ' },
        { type: 'text', content: 'website' },
      ])
    })

    it('handles drill links without leading slash in URL', () => {
      const content = '[Drill Name](drill/123)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([{ type: 'drill-link', drillId: '123', text: 'Drill Name' }])
    })

    it('handles drill links with .png file as ID', () => {
      const content = '[Drill Image](some_drill.png)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'drill-link', drillId: 'some_drill.png', text: 'Drill Image' },
      ])
    })

    it('handles consecutive drill links', () => {
      const content = '[A](/drills/a)[B](/drills/b)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'drill-link', drillId: 'a', text: 'A' },
        { type: 'drill-link', drillId: 'b', text: 'B' },
      ])
    })

    it('handles drill link at start', () => {
      const content = '[Drill](/drills/123) at the start'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'drill-link', drillId: '123', text: 'Drill' },
        { type: 'text', content: ' at the start' },
      ])
    })

    it('handles drill link at end', () => {
      const content = 'Finish with [Drill](/drills/123)'
      const result = parseContentForDrillLinks(content)
      expect(result).toEqual([
        { type: 'text', content: 'Finish with ' },
        { type: 'drill-link', drillId: '123', text: 'Drill' },
      ])
    })
  })

  describe('hasDrillLinks', () => {
    it('returns false for empty content', () => {
      expect(hasDrillLinks('')).toBe(false)
    })

    it('returns false for content without links', () => {
      expect(hasDrillLinks('Plain text content')).toBe(false)
    })

    it('returns false for non-drill links', () => {
      expect(hasDrillLinks('[link](https://example.com)')).toBe(false)
    })

    it('returns true for drill links', () => {
      expect(hasDrillLinks('[Drill](/drills/123)')).toBe(true)
    })

    it('returns true for mixed content with drill links', () => {
      expect(hasDrillLinks('Text with [Drill](/drills/123) in it')).toBe(true)
    })
  })
})
