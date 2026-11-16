import { describe, expect, it, vi, afterEach } from 'vitest'
import { shuffleArray } from '../shuffle'

describe('shuffleArray', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns a permutation of the input without mutating it', () => {
    const input = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    const snapshot = [...input]

    const result = shuffleArray(input)

    expect(input).toEqual(snapshot)
    expect(result).toHaveLength(input.length)
    expect([...result].sort((a, b) => a - b)).toEqual(snapshot)
  })

  it('handles empty and single-element arrays', () => {
    expect(shuffleArray([])).toEqual([])
    expect(shuffleArray(['only'])).toEqual(['only'])
  })

  it('shuffles deterministically when Math.random is mocked', () => {
    vi.spyOn(Math, 'random').mockReturnValue(0)

    expect(shuffleArray(['a', 'b', 'c', 'd'])).toEqual(['b', 'c', 'd', 'a'])
  })
})
