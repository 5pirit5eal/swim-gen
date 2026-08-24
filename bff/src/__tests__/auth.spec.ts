import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { getAuthHeaders } from "../auth";
import { logger } from "../logger";

const { mockGetIdTokenClient, mockGetRequestHeaders } = vi.hoisted(() => {
  return {
    mockGetIdTokenClient: vi.fn(),
    mockGetRequestHeaders: vi.fn(),
  };
});

vi.mock("google-auth-library", () => {
  return {
    GoogleAuth: class {
      getIdTokenClient = mockGetIdTokenClient;
    },
  };
});

describe("Auth Module", () => {
  const originalEnv = { ...process.env };

  beforeEach(() => {
    vi.clearAllMocks();
    process.env = {
      ...originalEnv,
      NODE_ENV: "production",
      BACKEND_URL: "https://backend.example.com",
    };
  });

  afterEach(() => {
    process.env = { ...originalEnv };
  });

  it("should pass through user authorization header", async () => {
    process.env.NODE_ENV = "development";
    const headers = await getAuthHeaders("Bearer test-token");
    expect(headers).toEqual({ Authorization: "Bearer test-token" });
  });

  it("should attach Google Identity token when available", async () => {
    const headersMap = new Map<string, string>();
    headersMap.set("Authorization", "Bearer google-id-token");
    mockGetIdTokenClient.mockResolvedValueOnce({
      getRequestHeaders: mockGetRequestHeaders.mockResolvedValueOnce(headersMap),
    });

    const headers = await getAuthHeaders("Bearer user-token");
    expect(headers).toEqual({
      Authorization: "Bearer user-token",
      "X-Serverless-Authorization": "Bearer google-id-token",
    });
  });

  it("should log sanitized error message and throw when Google auth fails", async () => {
    mockGetIdTokenClient.mockRejectedValueOnce(new Error("Network timeout connecting to metadata server"));
    const loggerSpy = vi.spyOn(logger, "error").mockImplementation(() => {});

    await expect(getAuthHeaders()).rejects.toThrow("Failed to authenticate with backend service.");
    expect(loggerSpy).toHaveBeenCalledWith(
      "Failed to get Google Identity token: Network timeout connecting to metadata server",
    );
  });
});
