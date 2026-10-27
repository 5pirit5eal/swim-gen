import { extractDrillIdFromUrl } from './markdownParser'

export type InlineSegment =
  | { type: 'text'; content: string }
  | { type: 'strong' | 'em'; children: InlineSegment[] }
  | { type: 'drill-link'; drillId: string; children: InlineSegment[] }
  | { type: 'link'; href: string; children: InlineSegment[] }

// Only inline links and paired emphasis are supported, not block Markdown or HTML.
// Match links first so underscores/asterisks in their URLs remain untouched.
const INLINE_PATTERN = new RegExp(
  [
    String.raw`(?<!!)\[(?<label>[^\]\n]+)\]\((?<url>[^\s()]+(?:\([^\s()]*\)[^\s()]*)*)\)`,
    String.raw`(?<!\*)(?<stars>\*{1,3})(?!\*)(?=\S)(?<starText>\S(?:[\s\S]*?\S)?)(?<!\*)\k<stars>(?!\*)`,
    String.raw`(?<![\p{L}\p{N}_])(?<underscores>_{1,3})(?!_)(?=\S)(?<underscoreText>\S(?:[\s\S]*?\S)?)(?<!_)\k<underscores>(?![\p{L}\p{N}_])`,
  ].join('|'),
  'gu',
)

function isSafeUrl(url: string): boolean {
  // Reject browser-normalized control characters and executable URL schemes,
  // including URLs that happen to contain a drill path.
  if (
    [...url].some((char) => char.charCodeAt(0) < 33 || char.charCodeAt(0) === 127 || char === '\\')
  ) {
    return false
  }
  const scheme = url.match(/^([a-z][a-z\d+.-]*):/i)?.[1]
  return !scheme || /^(https?)$/i.test(scheme)
}

export function parseInlineMarkdown(
  content: string,
  allowLinks = true,
  depth = 0,
): InlineSegment[] {
  if (!content) return []
  if (depth >= 16) return [{ type: 'text', content }]

  const segments: InlineSegment[] = []
  let lastIndex = 0

  // matchAll uses its own regex state, allowing recursive emphasis parsing.
  for (const match of content.matchAll(INLINE_PATTERN)) {
    const groups = match.groups!
    if (match.index > lastIndex) {
      segments.push({ type: 'text', content: content.slice(lastIndex, match.index) })
    }

    if (groups.label !== undefined && groups.url !== undefined) {
      const children = parseInlineMarkdown(groups.label, false, depth + 1)
      if (allowLinks && isSafeUrl(groups.url)) {
        const drillId = extractDrillIdFromUrl(groups.url)
        segments.push(
          drillId
            ? { type: 'drill-link', drillId, children }
            : { type: 'link', href: groups.url, children },
        )
      } else {
        segments.push(...children)
      }
    } else {
      const delimiter = groups.stars ?? groups.underscores!
      const children = parseInlineMarkdown(
        groups.starText ?? groups.underscoreText!,
        allowLinks,
        depth + 1,
      )
      segments.push({
        type: delimiter.length === 1 ? 'em' : 'strong',
        children: delimiter.length === 3 ? [{ type: 'em', children }] : children,
      })
    }

    lastIndex = match.index + match[0].length
  }

  if (lastIndex < content.length) {
    segments.push({ type: 'text', content: content.slice(lastIndex) })
  }
  return segments
}
