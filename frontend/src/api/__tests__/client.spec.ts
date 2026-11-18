import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { apiClient } from '../client'
import { ApiEndpoints } from '@/types'
import i18n from '@/plugins/i18n'
import { toast } from 'vue3-toastify'

vi.mock('vue3-toastify', () => ({
  toast: {
    info: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
  },
}))

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
          code: 'network_error',
          message: i18n.global.t('errors.connection_failed'),
          status: 0,
          details: 'Network down',
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
          code: 'timeout',
          message: i18n.global.t('errors.timeout', { time: 5 }),
          status: 0,
          details: 'The user aborted a request.',
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
          code: 'network_error',
          message: i18n.global.t('errors.connection_failed'),
          status: 0,
          details: 'Network error',
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
          code: 'timeout',
          message: i18n.global.t('errors.timeout', { time: 10 }),
          status: 0,
          details: 'The user aborted a request.',
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
          code: 'network_error',
          message: i18n.global.t('errors.connection_failed'),
          status: 0,
          details: 'Network error',
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
          code: 'timeout',
          message: i18n.global.t('errors.timeout', { time: 60 }),
          status: 0,
          details: 'The user aborted a request.',
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
          code: 'network_error',
          message: i18n.global.t('errors.connection_failed'),
          status: 0,
          details: 'Network error',
        },
      })
    })
  })

  describe('rate limit (429) handling', () => {
    const mockQueryRequest = { content: 'Generate a swim plan', language: 'en' }
    const mockQueryResponse = { plan_id: '123', title: 'Plan', description: '', table: [] }
    const tooManyRequests = () => ({
      ok: false,
      status: 429,
      statusText: 'Too Many Requests',
      headers: new Headers(),
    })

    it('should retry a 429 response with backoff and succeed', async () => {
      mockFetch
        .mockResolvedValueOnce(tooManyRequests())
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'Content-Type': 'application/json' }),
          json: () => Promise.resolve(mockQueryResponse),
        })

      const promise = apiClient.query(mockQueryRequest)
      await vi.advanceTimersByTimeAsync(5000) // first backoff: 5 seconds
      const result = await promise

      expect(mockFetch).toHaveBeenCalledTimes(2)
      expect(result).toEqual({ success: true, data: mockQueryResponse })
      expect(toast.info).toHaveBeenCalledWith(
        i18n.global.t('errors.rate_limited_retrying', { attempt: 1, max: 3, seconds: 5 }),
      )
    })

    it('should retry up to 3 times with 5s then 15s backoff and then fail with a rate limit error', async () => {
      mockFetch.mockResolvedValue(tooManyRequests())

      const promise = apiClient.query(mockQueryRequest)
      await vi.advanceTimersByTimeAsync(5000) // first backoff: 5 seconds
      await vi.advanceTimersByTimeAsync(15000) // second backoff: 15 seconds
      await vi.advanceTimersByTimeAsync(15000) // third backoff: 15 seconds
      const result = await promise

      expect(mockFetch).toHaveBeenCalledTimes(4) // initial request + 3 retries
      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.rate_limited'),
          status: 429,
          details: 'Too Many Requests',
        },
      })
      expect(toast.info).toHaveBeenCalledTimes(3)
      expect(toast.info).toHaveBeenNthCalledWith(
        2,
        i18n.global.t('errors.rate_limited_retrying', { attempt: 2, max: 3, seconds: 15 }),
      )
      expect(toast.info).toHaveBeenNthCalledWith(
        3,
        i18n.global.t('errors.rate_limited_retrying', { attempt: 3, max: 3, seconds: 15 }),
      )
    })

    it('should not retry non-429 error responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        headers: new Headers(),
      })

      const result = await apiClient.query(mockQueryRequest)

      expect(mockFetch).toHaveBeenCalledTimes(1)
      expect(result.success).toBe(false)
      expect(result.error?.status).toBe(500)
    })

    it('should return the rate limit error message on 429 without retry for endpoints without retries', async () => {
      mockFetch.mockResolvedValueOnce(tooManyRequests())

      const result = await apiClient.checkHealth()

      expect(mockFetch).toHaveBeenCalledTimes(1)
      expect(result).toEqual({
        success: false,
        error: {
          message: i18n.global.t('errors.rate_limited'),
          status: 429,
          details: 'Too Many Requests',
        },
      })
    })
  })
})
