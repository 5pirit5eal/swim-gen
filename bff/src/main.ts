import express from 'express';
import dotenv from 'dotenv';
import axios from 'axios';
import { GoogleAuth } from 'google-auth-library';

dotenv.config();

const app = express();
const port = process.env.PORT || 8080;
const backendUrl = process.env.BACKEND_URL;

// Middleware to handle JSON bodies
app.use(express.json());

// Google Auth setup
const auth = new GoogleAuth();

async function getAuthHeaders() {
    // For local development against a local backend, we might not have a token.
    // The NODE_ENV check allows skipping auth.
    if (!backendUrl || process.env.NODE_ENV === 'development') {
        console.log('Skipping auth for local development.');
        return {};
    }

    console.log(`Fetching token for backend: ${backendUrl}`);
    try {
        const client = await auth.getIdTokenClient(backendUrl);
        const headers = await client.getRequestHeaders();
        const authorizationHeader = headers.get('Authorization');
        if (!authorizationHeader) {
            throw new Error('Authorization header not found in response from Google Auth.');
        }
        return { Authorization: authorizationHeader };
    } catch (error) {
        console.error('Failed to get auth token:', error);
        throw new Error('Failed to authenticate with backend service.');
    }
}

// Generic proxy handler for all API requests
async function proxyRequest(req: express.Request, res: express.Response) {
    const { method, originalUrl, body } = req;

    // The frontend will call paths like /api/query. We strip /api before forwarding.
    const backendPath = originalUrl.replace(/^\/api/, '');
    const targetUrl = `${backendUrl}${backendPath}`;

    console.log(`Proxying request: ${method} ${originalUrl} -> ${targetUrl}`);

    try {
        const authHeaders = await getAuthHeaders();

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
    } catch (error: any) {
        console.error(`Error proxying request to ${targetUrl}:`, error.message);
        if (error.response) {
            res.status(error.response.status).json(error.response.data);
        } else {
            res.status(500).json({ message: 'Error proxying request to backend' });
        }
    }
}

app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

// All API routes from the frontend are prefixed with /api
app.use('/api', proxyRequest);

app.listen(port, () => {
    console.log(`BFF server listening on port ${port}`);
});
