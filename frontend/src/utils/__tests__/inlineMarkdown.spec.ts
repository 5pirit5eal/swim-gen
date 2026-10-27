import { describe, expect, it } from 'vitest'
import { parseInlineMarkdown } from '../inlineMarkdown'

const text = (content: string) => ({ type: 'text', content })

describe('parseInlineMarkdown', () => {
  it('preserves plain text and line breaks', () => {
    expect(parseInlineMarkdown('')).toEqual([])
    expect(parseInlineMarkdown('First line\nSecond line')).toEqual([
      text('First line\nSecond line'),
    ])
  })

  it.each(['*', '_', '**', '__', '***', '___'])('parses paired %s emphasis', (delimiter) => {
    const children = [text('Swim')]
    expect(parseInlineMarkdown(`${delimiter}Swim${delimiter}`)).toEqual([
      {
        type: delimiter.length === 1 ? 'em' : 'strong',
        children: delimiter.length === 3 ? [{ type: 'em', children }] : children,
      },
    ])
  })

  it('parses nested emphasis', () => {
    expect(parseInlineMarkdown('**Swim _slowly_ now**')).toEqual([
      {
        type: 'strong',
        children: [text('Swim '), { type: 'em', children: [text('slowly')] }, text(' now')],
      },
    ])
    expect(parseInlineMarkdown('*Swim **slowly** now*')).toEqual([
      {
        type: 'em',
        children: [text('Swim '), { type: 'strong', children: [text('slowly')] }, text(' now')],
      },
    ])
  })

  it.each([
    '/drill/finger_drag',
    '/drills/finger_drag.webp',
    'drill/finger_drag.png',
    'drills/finger_drag',
    'https://example.com/drills/finger_drag?lang=de#video',
    'finger_drag.png',
  ])('uses the table drill recognition for %s', (url) => {
    expect(parseInlineMarkdown(`Try [Finger Drag](${url}) now`)).toEqual([
      text('Try '),
      { type: 'drill-link', drillId: 'finger_drag', children: [text('Finger Drag')] },
      text(' now'),
    ])
  })

  it('handles links inside emphasis and emphasis inside link labels', () => {
    expect(parseInlineMarkdown('**Try [_Finger Drag_](/drills/finger_drag)**')).toEqual([
      {
        type: 'strong',
        children: [
          text('Try '),
          {
            type: 'drill-link',
            drillId: 'finger_drag',
            children: [{ type: 'em', children: [text('Finger Drag')] }],
          },
        ],
      },
    ])
  })

  it('keeps safe non-drill links, including parentheses in URLs', () => {
    expect(parseInlineMarkdown('[Info](https://example.com/swim_(sport))')).toEqual([
      { type: 'link', href: 'https://example.com/swim_(sport)', children: [text('Info')] },
    ])
  })

  it.each([
    'javascript:alert',
    'javascript:/drills/unsafe',
    'data:text/html,unsafe',
    'vbscript:unsafe',
    'java\u0000script:unsafe',
    '\\evil.example.com',
  ])('does not create links for unsafe URL %s', (url) => {
    expect(parseInlineMarkdown(`[Label](${url})`)).toEqual([text('Label')])
  })

  it.each(['*unfinished', '__unfinished', '* spaced *', 'some_drill_name', 'some__drill__name'])(
    'leaves unpaired markers, whitespace and word underscores untouched: %s',
    (content) => {
      expect(parseInlineMarkdown(content)).toEqual([text(content)])
    },
  )

  it('does not interpret other Markdown or raw HTML', () => {
    const content =
      '# Heading\n- List item\n`code`\n~~strike~~\n![Image](https://example.com/image.png)\n<b>HTML</b>'
    expect(parseInlineMarkdown(content)).toEqual([text(content)])
  })
})
