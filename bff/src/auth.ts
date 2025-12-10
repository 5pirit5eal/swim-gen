import { GoogleAuth } from "google-auth-library";
import { logger } from "./logger";

// Initialize Google Auth client once at module level for better performance
const auth = new GoogleAuth();

/**
 * Build headers for backend request with dual authentication:
 * 1. Pass through user's Authorization header from frontend (Supabase token)
 * 2. Add Google Identity token for Cloud Run service-to-service authentication
 *
 * @param userAuthHeader - The Authorization header from the frontend request (optional)
 * @returns Headers object with both user auth and service-to-service auth
 */
export async function getAuthHeaders(userAuthHeader?: string): Promise<Record<string, string>> {
  const backendUrl = process.env.BACKEND_URL;
  const headers: Record<string, string> = {};

  // Pass through user's Authorization header if present
  if (userAuthHeader) {
    headers["Authorization"] = userAuthHeader;
  }

  // For local development, skip Google Identity token
  if (!backendUrl || process.env.NODE_ENV === "development") {
    logger.debug("Skipping Google Identity auth for local development.");
    return headers;
  }

  // Add Google Identity token for Cloud Run service-to-service authentication
  // The backend (if running on Cloud Run) will verify this token
  logger.debug(`Fetching Google Identity token for backend: ${backendUrl}`);
  try {
    const client = await auth.getIdTokenClient(backendUrl);
    const googleHeaders = await client.getRequestHeaders();
    const googleAuthHeader = googleHeaders.get("Authorization");
    if (!googleAuthHeader) {
      throw new Error("Authorization header not found in response from Google Auth.");
    }
    // Add Google Identity token as X-Serverless-Authorization
    // This keeps it separate from the user's Supabase token
    headers["X-Serverless-Authorization"] = googleAuthHeader;
  } catch (error) {
    logger.error("Failed to get Google Identity token:", error);
    throw new Error("Failed to authenticate with backend service.");
  }

  return headers;
}
