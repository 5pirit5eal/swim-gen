/**
 * Represents a segment of parsed content - either plain text or a drill link
 */
export type ContentSegment =
  | { type: 'text'; content: string }
  | { type: 'drill-link'; drillId: string; text: string }

/**
 * Regex pattern to match markdown links: [text](url)
 * Captures: [1] = link text, [2] = url
 */
const MARKDOWN_LINK_REGEX = /\[([^\]]+)\]\(([^)]+)\)/g

/**
 * Extracts the drill ID from a drill URL
 * Expected URL format: /drills/{id} or full URL like https://example.com/drills/{id}
 */
export function extractDrillIdFromUrl(url: string): string | null {
  // Try to match /drills/{id} or /drill/{id} pattern (both singular and plural)
  const drillPathMatch = url.match(/\/drills?\/([^/?#]+)/)
  if (drillPathMatch && drillPathMatch[1]) {
    console.debug('[markdownParser] Found drill ID in path:', drillPathMatch[1])
    return drillPathMatch[1]
  }

  // If URL is just the ID itself (e.g., from a simplified link)
  if (!url.includes('/') && !url.includes(':')) {
    console.debug('[markdownParser] URL is just an ID:', url)
    return url
  }

  console.debug('[markdownParser] No drill ID found')
  return null
}

/**
 * Parses content string and extracts markdown links, converting them to typed segments.
 * Drill links are identified by URLs containing '/drills/' path.
 *
 * @param content - The content string potentially containing markdown links
 * @returns Array of ContentSegment objects representing text and drill links
 */
export function parseContentForDrillLinks(content: string): ContentSegment[] {
  if (!content) {
    console.debug('[markdownParser] Empty content, returning empty array')
    return []
  }

  const segments: ContentSegment[] = []
  let lastIndex = 0

  // Reset regex state
  MARKDOWN_LINK_REGEX.lastIndex = 0

  let match: RegExpExecArray | null
  while ((match = MARKDOWN_LINK_REGEX.exec(content)) !== null) {
    const fullMatch = match[0]
    const linkText = match[1]
    const url = match[2]
    const matchIndex = match.index

    // Skip if we don't have required parts
    if (!linkText || !url) {
      console.debug('[markdownParser] Skipping - missing linkText or url')
      continue
    }

    // Add text before this match if any
    if (matchIndex > lastIndex) {
      const textBefore = content.slice(lastIndex, matchIndex)
      if (textBefore) {
        console.debug(
          '[markdownParser] Adding text segment before match:',
          JSON.stringify(textBefore),
        )
        segments.push({ type: 'text', content: textBefore })
      }
    }

    // Check if this is a drill link
    const drillId = extractDrillIdFromUrl(url)
    console.debug('[markdownParser] Extracted drillId:', drillId, 'from url:', url)
    if (drillId) {
      console.debug('[markdownParser] Adding drill-link segment:', { drillId, text: linkText })
      segments.push({
        type: 'drill-link',
        drillId,
        text: linkText,
      })
    } else {
      // Not a drill link, keep as plain text with the original markdown
      segments.push({ type: 'text', content: fullMatch })
    }

    lastIndex = matchIndex + fullMatch.length
  }

  // Add remaining text after last match
  if (lastIndex < content.length) {
    const remaining = content.slice(lastIndex)
    segments.push({ type: 'text', content: remaining })
  }

  // If no segments were created, return the original content as text
  if (segments.length === 0 && content) {
    console.debug('[markdownParser] No segments created, returning original content as text')
    return [{ type: 'text', content }]
  }

  console.debug('[markdownParser] Final segments:', JSON.stringify(segments))
  return segments
}

/**
 * Checks if content contains any drill links
 */
export function hasDrillLinks(content: string): boolean {
  if (!content) return false
  MARKDOWN_LINK_REGEX.lastIndex = 0
  let match: RegExpExecArray | null
  while ((match = MARKDOWN_LINK_REGEX.exec(content)) !== null) {
    const url = match[2]
    if (url && extractDrillIdFromUrl(url)) {
      return true
    }
  }
  return false
}
