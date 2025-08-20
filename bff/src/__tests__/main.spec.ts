import { describe, it, expect, vi, beforeEach } from 'vitest';
import request from 'supertest';
import axios from 'axios';
import { app } from '../main';
import * as authModule from '../auth';

vi.mock('axios');

// Reset mocks between tests
beforeEach(() => {
    vi.clearAllMocks();
    process.env.BACKEND_URL = 'http://backend.test';
});

describe('BFF Server', () => {
    describe('GET /health', () => {
        it('should return 200 OK', async () => {
            const response = await request(app).get('/health');
            expect(response.status).toBe(200);
            expect(response.text).toBe('OK');
        });
    });

    describe('API Proxy', () => {
        it('should proxy POST requests to the backend with auth headers', async () => {
            process.env.NODE_ENV = 'test'
            // Mock getAuthHeaders to provide Authorization header
            vi.spyOn(authModule, 'getAuthHeaders').mockResolvedValue({
                Authorization: 'Bearer mocked-token',
            });

            // Mock Axios
            vi.mocked(axios).mockResolvedValue({
                status: 200,
                data: { success: true, message: 'Backend response' },
            });

            const requestBody = { query: 'test' };
            const response = await request(app).post('/api/query').send(requestBody);

            expect(response.status).toBe(200);
            expect(response.body).toEqual({ success: true, message: 'Backend response' });

            expect(axios).toHaveBeenCalledWith({
                method: 'POST',
                url: `${process.env.BACKEND_URL}/query`,
                data: requestBody,
                headers: {
                    Authorization: 'Bearer mocked-token',
                    'Content-Type': 'application/json',
                },
            });
        });

        it('should proxy without auth headers in development mode', async () => {
            process.env.NODE_ENV = 'development'
            // getAuthHeaders returns empty in development; ensure spy returns empty
            vi.spyOn(authModule, 'getAuthHeaders').mockResolvedValue({});

            vi.mocked(axios).mockResolvedValue({
                status: 200,
                data: { ok: true },
            });

            const response = await request(app).post('/api/ping').send({ a: 1 })
            expect(response.status).toBe(200)
            expect(response.body).toEqual({ ok: true })

            expect(axios).toHaveBeenCalledWith({
                method: 'POST',
                url: `${process.env.BACKEND_URL}/ping`,
                data: { a: 1 },
                headers: {
                    'Content-Type': 'application/json',
                },
            })
        })

        it('should handle errors from the backend', async () => {
            process.env.NODE_ENV = 'test'
            vi.spyOn(authModule, 'getAuthHeaders').mockResolvedValue({
                Authorization: 'Bearer mocked-token',
            });

            vi.mocked(axios).mockRejectedValue({
                response: {
                    status: 500,
                    data: { message: 'Internal Server Error' },
                },
            });

            const response = await request(app).post('/api/some-endpoint').send({});

            expect(response.status).toBe(500);
            expect(response.body).toEqual({ message: 'Internal Server Error' });
        });
    });
});
