import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { apiClient } from '../client'
import { ApiEndpoints } from '@/types'
import i18n from '@/plugins/i18n'

// Mock global fetch
const mockFetch = vi.fn()
global.fetch = mockFetch

describe('ApiClient', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    mockFetch.mockClear()
  })

  afterEach(() => {
    vi.runOnlyPendingTimers()
    vi.useRealTimers()
  })

  describe('checkHealth', () => {
    it('should return success true and data on successful health check', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
        headers: new Headers({ 'Content-Type': 'text/plain' }),
        text: () => Promise.resolve('API is healthy'),
      })

      const result = await apiClient.checkHealth()

      expect(mockFetch).toHaveBeenCalledWith(`/api/${ApiEndpoints.HEALTH}`, expect.any(Object))
      expect(result).toEqual({ success: true, data: 'API is healthy' })
    })

    it('should return success false and error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
      })

      const result = await apiClient.checkHealth()

      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.api_request_failed', { endpoint: ApiEndpoints.HEALTH }),
          status: 500,
          details: 'Internal Server Error',
        },
      })
    })

    it('should return success false and error on network failure', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network down'))

      const result = await apiClient.checkHealth()

      expect(result).toEqual({
        success: false,
        error: {
          message: 'Network down',
          status: 0,
          details: i18n.global.t('errors.connection_failed'),
        },
      })
    })

    it('should abort the request if it times out', async () => {
      const error = new Error('The user aborted a request.')
      error.name = 'AbortError'
      mockFetch.mockRejectedValueOnce(error)

      const promise = apiClient.checkHealth()
      vi.advanceTimersByTime(5000)

      await expect(promise).resolves.toEqual({
        success: false,
        error: {
          message: 'The user aborted a request.',
          status: 0,
          details: i18n.global.t('errors.timeout', { time: 5 }),
        },
      })
    })
  })

  describe('generatePrompt', () => {
    const mockRequest = { language: 'en' }
    const mockResponseData = {
      prompt: 'Generate a short swim training plan.',
    }

    it('should return success true and data on successful prompt generation', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        headers: new Headers({ 'Content-Type': 'application/json' }),
        json: () => Promise.resolve(mockResponseData),
      })

      const result = await apiClient.generatePrompt(mockRequest)

      expect(mockFetch).toHaveBeenCalledWith(
        `/api/${ApiEndpoints.PROMPT}`,
        expect.objectContaining({
          method: 'POST',
          headers: new Headers({ 'Content-Type': 'application/json' }),
          body: JSON.stringify(mockRequest),
        }),
      )
      expect(result).toEqual({ success: true, data: mockResponseData })
    })

    it('should return success false and error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        headers: new Headers(),
      })

      const result = await apiClient.generatePrompt(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.api_request_failed', { endpoint: ApiEndpoints.PROMPT }),
          status: 400,
          details: 'Bad Request',
        },
      })
    })

    it('should return success false and error on network failure', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const result = await apiClient.generatePrompt(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: 'Network error',
          status: 0,
          details: i18n.global.t('errors.connection_failed'),
        },
      })
    })

    it('should return timeout error if request takes longer than 10 seconds', async () => {
      const error = new Error('The user aborted a request.')
      error.name = 'AbortError'
      mockFetch.mockRejectedValueOnce(error)

      const promise = apiClient.generatePrompt(mockRequest)
      vi.advanceTimersByTime(10000)

      await expect(promise).resolves.toEqual({
        success: false,
        error: {
          message: 'The user aborted a request.',
          status: 0,
          details: i18n.global.t('errors.timeout', { time: 10 }),
        },
      })
    })
  })

  describe('query', () => {
    const mockRequest = {
      content: 'Generate a 3-day beginner swim plan', // Renamed 'query' to 'content'
      language: 'en',
    }
    const mockResponseData = {
      plan: {
        title: 'Beginner Swim Plan',
        description: 'A plan for beginners',
        table: [],
      },
    }

    it('should return success true and data on successful query', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        headers: new Headers({ 'Content-Type': 'application/json' }),
        json: () => Promise.resolve(mockResponseData),
      })

      const result = await apiClient.query(mockRequest)

      expect(mockFetch).toHaveBeenCalledWith(
        `/api/${ApiEndpoints.QUERY}`,
        expect.objectContaining({
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization: 'Bearer null',
          }),
          body: JSON.stringify(mockRequest),
        }),
      )
      expect(result).toEqual({ success: true, data: mockResponseData })
    })

    it('should return success false and error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
      })

      const result = await apiClient.query(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.api_request_failed', { endpoint: ApiEndpoints.QUERY }),
          status: 500,
          details: 'Internal Server Error',
        },
      })
    })

    it('should return success false and error on network failure', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const result = await apiClient.query(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: 'Network error',
          status: 0,
          details: i18n.global.t('errors.connection_failed'),
        },
      })
    })

    it('should return timeout error if request takes longer than 60 seconds', async () => {
      const error = new Error('The user aborted a request.')
      error.name = 'AbortError'
      mockFetch.mockRejectedValueOnce(error)

      const promise = apiClient.query(mockRequest)
      vi.advanceTimersByTime(60000)

      await expect(promise).resolves.toEqual({
        success: false,
        error: {
          message: 'The user aborted a request.',
          status: 0,
          details: i18n.global.t('errors.timeout', { time: 60 }),
        },
      })
    })
  })

  describe('exportPDF', () => {
    const mockRequest = {
      title: 'My Plan',
      description: 'Description',
      table: [], // Added missing property
    }
    const mockResponseData = {
      pdfUri: 'http://example.com/plan.pdf',
    }

    it('should return success true and data on successful PDF export', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        headers: new Headers({ 'Content-Type': 'application/json' }),
        json: () => Promise.resolve(mockResponseData),
      })

      const result = await apiClient.exportPDF(mockRequest)

      expect(mockFetch).toHaveBeenCalledWith(
        `/api/${ApiEndpoints.EXPORT_PDF}`,
        expect.objectContaining({
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization: 'Bearer null',
          }),
          body: JSON.stringify(mockRequest),
        }),
      )
      expect(result).toEqual({ success: true, data: mockResponseData })
    })

    it('should return success false and error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        headers: new Headers(),
      })

      const result = await apiClient.exportPDF(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.api_request_failed', {
            endpoint: ApiEndpoints.EXPORT_PDF,
          }),
          status: 400,
          details: 'Bad Request',
        },
      })
    })

    it('should return success false and error on network failure', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const result = await apiClient.exportPDF(mockRequest)

      expect(result).toEqual({
        success: false,
        error: {
          message: 'Network error',
          status: 0,
          details: i18n.global.t('errors.connection_failed'),
        },
      })
    })
  })
})
