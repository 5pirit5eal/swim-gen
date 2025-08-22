import { GoogleAuth } from 'google-auth-library'

// Google Auth setup
const auth = new GoogleAuth()

export async function getAuthHeaders(): Promise<Record<string, string>> {
    const backendUrl = process.env.BACKEND_URL
    // For local development against a local backend, we might not have a token.
    // The NODE_ENV check allows skipping auth.
    if (!backendUrl || process.env.NODE_ENV === 'development') {
        console.log('Skipping auth for local development.')
        return {} as Record<string, string>
    }

    console.log(`Fetching token for backend: ${backendUrl}`)
    try {
        const client = await auth.getIdTokenClient(backendUrl)
        const headers = await client.getRequestHeaders()
        const authorizationHeader = headers.get('Authorization')
        if (!authorizationHeader) {
            throw new Error('Authorization header not found in response from Google Auth.')
        }
        return { Authorization: authorizationHeader } as Record<string, string>
    } catch (error) {
        console.error('Failed to get auth token:', error)
        throw new Error('Failed to authenticate with backend service.')
    }
}
