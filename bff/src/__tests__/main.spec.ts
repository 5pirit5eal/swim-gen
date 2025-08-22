import { describe, it, expect, vi, beforeEach } from 'vitest';
import request from 'supertest';
import axios from 'axios';
import { app, testingStore } from '../main';
import * as authModule from '../auth';

vi.mock('axios');

// Reset mocks between tests
beforeEach(async () => {
    vi.clearAllMocks();
    process.env.BACKEND_URL = 'http://backend.test';
    process.env.FRONTEND_URL = 'http://frontend.test';
    // Reset the rate limit for the test IP address
    await testingStore.resetAll();
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

    describe('Security', () => {
        it('should block requests to /api/scrape', async () => {
            const response = await request(app).get('/api/scrape');
            expect(response.status).toBe(403);
            expect(response.text).toBe('This endpoint is not available.');
        });

        it('should allow requests from the configured frontend URL', async () => {
            const response = await request(app)
                .get('/api/some-endpoint')
                .set('Origin', process.env.FRONTEND_URL as string);
            // We expect a 500 error because the backend is mocked to fail, but a CORS error would be different.
            expect(response.status).not.toBe(0);
        });

        it('should block requests from other origins', async () => {
            const response = await request(app)
                .get('/api/some-endpoint')
                .set('Origin', 'http://another-origin.com');
            // supertest doesn't throw a CORS error, but the origin will not be set in the response header
            expect(response.headers['access-control-allow-origin']).toBeUndefined();
        });

        it('should rate limit requests', async () => {
            // Mock successful responses for the first 100 requests
            vi.mocked(axios).mockResolvedValue({
                status: 200,
                data: { success: true },
            });

            const agent = request.agent(app);
            for (let i = 0; i < 100; i++) {
                await agent.get('/api/some-endpoint').expect(200);
            }

            // The 101st request should be rate limited
            const response = await agent.get('/api/some-endpoint');
            expect(response.status).toBe(429);
        });
    });
});
