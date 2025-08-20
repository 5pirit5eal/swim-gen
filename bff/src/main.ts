import express from 'express';
import dotenv from 'dotenv';
import axios from 'axios';
import * as authModule from './auth';

dotenv.config();

const app = express();
const port = process.env.PORT || 8080;

// Middleware to handle JSON bodies
app.use(express.json());

// getAuthHeaders is imported from './auth'

// Generic proxy handler for all API requests
async function proxyRequest(req: express.Request, res: express.Response) {
    const { method, originalUrl, body } = req;

    // The frontend will call paths like /api/query. We strip /api before forwarding.
    const backendPath = originalUrl.replace(/^\/api/, '');
    const targetUrl = `${process.env.BACKEND_URL}${backendPath}`;

    console.log(`Proxying request: ${method} ${originalUrl} -> ${targetUrl}`)

    try {
        const authHeaders = await authModule.getAuthHeaders()

        const response = await axios({
            method,
            url: targetUrl,
            data: body,
            headers: {
                ...authHeaders,
                'Content-Type': 'application/json',
            },
        });

        res.status(response.status).json(response.data);
    } catch (error) {
        console.error(`Error proxying request to ${targetUrl}:`, error);
        if (error && typeof error === 'object' && 'response' in error && error.response) {
            const axiosError = error as { response: { status: number; data: unknown } };
            res.status(axiosError.response.status).json(axiosError.response.data);
        } else {
            res.status(500).json({ message: 'Error proxying request to backend' });
        }
    }
}

app.get('/health', (req, res) => {
    res.status(200).send('OK')
});

// All API routes from the frontend are prefixed with /api
app.use('/api', proxyRequest);

// Avoid binding a port during tests
if (process.env.NODE_ENV !== 'test') {
    app.listen(port, () => {
        console.log(`BFF server listening on port ${port}`)
    });
}

// Exported for testing/mocking in unit tests
export { app }
