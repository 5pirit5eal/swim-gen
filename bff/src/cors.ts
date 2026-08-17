import type { CorsOptions } from "cors";

/**
 * Normalizes an origin string:
 * - Ensures scheme exists (defaults to https:// if omitted)
 * - Strips trailing slashes
 * - Normalizes scheme + hostname to lowercase
 */
export function normalizeOrigin(origin: string): string {
  const trimmed = origin.trim();
  if (!trimmed || trimmed === "*") {
    return trimmed;
  }

  let withScheme = trimmed;
  if (!/^https?:\/\//i.test(withScheme)) {
    withScheme = `https://${withScheme}`;
  }

  try {
    const parsed = new URL(withScheme);
    // parsed.origin gives lowercase scheme + host + port (without path or trailing slash)
    return parsed.origin;
  } catch {
    return withScheme.replace(/\/+$/, "");
  }
}

/**
 * Parses and returns the list of allowed origins from an environment variable or parameter.
 * Supports comma-separated origin lists.
 */
export function getAllowedOrigins(envValue?: string): string[] {
  const raw = typeof envValue === "string" ? envValue : process.env.FRONTEND_URL;
  if (!raw || typeof raw !== "string" || raw.trim().length === 0) {
    return [];
  }

  return raw
    .split(",")
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
    .map(normalizeOrigin);
}

/**
 * Checks whether an incoming origin header matches the allowed origins list.
 */
export function isOriginAllowed(
  origin: string | undefined,
  allowedOrigins: string[] = getAllowedOrigins(),
): boolean {
  // Allow requests without Origin header (curl, server-to-server, health checks)
  if (!origin) {
    return true;
  }

  // Reject literal "null" origin strings (e.g. sandboxed iframes, data URLs) unless explicitly allowed
  if (origin === "null") {
    return allowedOrigins.includes("null");
  }

  if (allowedOrigins.includes("*")) {
    return true;
  }

  const normalized = normalizeOrigin(origin);
  return allowedOrigins.includes(normalized);
}

/**
 * Creates Express CORS middleware options enforcing strict origin matching and safe defaults.
 */
export function createCorsOptions(): CorsOptions {
  return {
    origin: (origin, callback) => {
      const allowed = getAllowedOrigins();
      if (isOriginAllowed(origin, allowed)) {
        // Passing true will echo the exact requested origin in Access-Control-Allow-Origin
        callback(null, true);
      } else {
        // Passing false does NOT emit Access-Control-Allow-Origin, rejecting the request securely
        callback(null, false);
      }
    },
    // Credentials (cookies) are disabled because the app authenticates exclusively via Bearer tokens
    credentials: false,
    optionsSuccessStatus: 200,
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allowedHeaders: ["Content-Type", "Authorization"],
  };
}
