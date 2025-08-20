import { describe, it, expect, vi, Mock } from 'vitest';
import request from 'supertest';
import { app } from '../main'; // Assuming app is exported from main.ts
import axios from 'axios';
import { GoogleAuth } from 'google-auth-library';

// Mock the external dependencies
vi.mock('axios');
vi.mock('google-auth-library');

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
            // Mock Google Auth
            const mockGetIdTokenClient = vi.fn().mockResolvedValue({
                getRequestHeaders: vi.fn().mockResolvedValue({
                    Authorization: 'Bearer mocked-token',
                }),
            });
            (GoogleAuth as Mock).mockImplementation(() => ({
                getIdTokenClient: mockGetIdTokenClient,
            }));

            // Mock Axios
            (axios as any).mockResolvedValue({
                status: 200,
                data: { success: true, message: 'Backend response' },
            });

            const requestBody = { query: 'test' };
            const response = await request(app)
                .post('/api/query')
                .send(requestBody);

            expect(response.status).toBe(200);
            expect(response.body).toEqual({ success: true, message: 'Backend response' });

            // Verify Axios was called correctly
            expect(axios).toHaveBeenCalledWith({
                method: 'POST',
                url: `${process.env.BACKEND_URL}/query`,
                data: requestBody,
                headers: {
                    'Authorization': 'Bearer mocked-token',
                    'Content-Type': 'application/json',
                },
            });
        });

        it('should handle errors from the backend', async () => {
            // Mock Google Auth (can be simpler for this test)
            const mockGetIdTokenClient = vi.fn().mockResolvedValue({
                getRequestHeaders: vi.fn().mockResolvedValue({
                    Authorization: 'Bearer mocked-token',
                }),
            });
            (GoogleAuth as Mock).mockImplementation(() => ({
                getIdTokenClient: mockGetIdTokenClient,
            }));

            // Mock Axios to simulate a backend error
            (axios as any).mockRejectedValue({
                response: {
                    status: 500,
                    data: { message: 'Internal Server Error' },
                },
            });

            const response = await request(app)
                .post('/api/some-endpoint')
                .send({});

            expect(response.status).toBe(500);
            expect(response.body).toEqual({ message: 'Internal Server Error' });
        });
    });
});
