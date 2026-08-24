import { describe, it, expect, vi, beforeEach } from "vitest";
import request from "supertest";
import axios from "axios";
import { app, testingStore, MAX_JSON_BYTES, MAX_MULTIPART_BYTES } from "../main";
import * as authModule from "../auth";

vi.mock("axios");

// Reset mocks between tests
beforeEach(async () => {
  vi.restoreAllMocks();
  vi.clearAllMocks();
  process.env.NODE_ENV = "development";
  process.env.BACKEND_URL = "http://backend.test";
  process.env.FRONTEND_URL = "http://frontend.test";
  vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({});
  // Reset the rate limit for the test IP address
  await testingStore.resetAll();
});

describe("BFF Server", () => {
  describe("GET /health", () => {
    it("should return 200 OK", async () => {
      const response = await request(app).get("/health");
      expect(response.status).toBe(200);
      expect(response.text).toBe("OK");
      expect(response.headers["x-powered-by"]).toBeUndefined();
      expect(response.headers["referrer-policy"]).toBe("strict-origin-when-cross-origin");
      expect(response.headers["permissions-policy"]).toContain("camera=()");
    });
  });

  describe("API Proxy", () => {
    it("should proxy POST requests to the backend with auth headers", async () => {
      process.env.NODE_ENV = "test";
      // Mock getAuthHeaders to provide both user auth and Google Identity token
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({
        Authorization: "Bearer user-supabase-token",
        "X-Serverless-Authorization": "Bearer google-identity-token",
      });

      // Mock Axios
      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true, message: "Backend response" },
      });

      const requestBody = { query: "test" };
      const response = await request(app)
        .post("/api/query")
        .set("Authorization", "Bearer user-supabase-token")
        .send(requestBody);

      expect(response.status).toBe(200);
      expect(response.body).toEqual({ success: true, message: "Backend response" });
      expect(response.headers["x-powered-by"]).toBeUndefined();
      expect(response.headers["cache-control"]).toBe("no-store");

      expect(authModule.getAuthHeaders).toHaveBeenCalledWith("Bearer user-supabase-token");
      expect(axios).toHaveBeenCalledWith({
        method: "POST",
        url: `${process.env.BACKEND_URL}/query`,
        data: requestBody,
        headers: {
          Authorization: "Bearer user-supabase-token",
          "X-Serverless-Authorization": "Bearer google-identity-token",
          "Content-Type": "application/json",
        },
        maxBodyLength: MAX_JSON_BYTES,
        maxContentLength: MAX_JSON_BYTES,
      });
    });

    it("should proxy without auth headers in development mode", async () => {
      process.env.NODE_ENV = "development";
      // getAuthHeaders returns empty in development; ensure spy returns empty
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({});

      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { ok: true },
      });

      const response = await request(app).post("/api/ping").send({ a: 1 });
      expect(response.status).toBe(200);
      expect(response.body).toEqual({ ok: true });

      expect(authModule.getAuthHeaders).toHaveBeenCalledWith(undefined);
      expect(axios).toHaveBeenCalledWith({
        method: "POST",
        url: `${process.env.BACKEND_URL}/ping`,
        data: { a: 1 },
        headers: {
          "Content-Type": "application/json",
        },
        maxBodyLength: MAX_JSON_BYTES,
        maxContentLength: MAX_JSON_BYTES,
      });
    });

    it("should handle errors from the backend", async () => {
      process.env.NODE_ENV = "test";
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({
        Authorization: "Bearer user-token",
        "X-Serverless-Authorization": "Bearer google-token",
      });

      vi.mocked(axios).mockRejectedValue({
        response: {
          status: 500,
          data: { message: "Internal Server Error" },
        },
      });

      const response = await request(app)
        .post("/api/some-endpoint")
        .set("Authorization", "Bearer user-token")
        .send({});

      expect(response.status).toBe(500);
      expect(response.body).toEqual({ message: "Internal Server Error" });
    });

    it("should forward text/plain response and error from backend properly", async () => {
      process.env.NODE_ENV = "test";
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({});

      vi.mocked(axios).mockResolvedValueOnce({
        status: 200,
        headers: { "content-type": "text/plain; charset=utf-8" },
        data: "plain text payload",
      });

      const successResponse = await request(app).get("/api/text-endpoint");
      expect(successResponse.status).toBe(200);
      expect(successResponse.text).toBe("plain text payload");
      expect(successResponse.headers["content-type"]).toContain("text/plain");

      vi.mocked(axios).mockRejectedValueOnce({
        response: {
          status: 400,
          headers: { "content-type": "text/plain; charset=utf-8" },
          data: "bad request plain text",
        },
      });

      const errorResponse = await request(app).get("/api/text-error");
      expect(errorResponse.status).toBe(400);
      expect(errorResponse.text).toBe("bad request plain text");
      expect(errorResponse.headers["content-type"]).toContain("text/plain");
    });

    it("should handle anonymous requests without user Authorization header", async () => {
      process.env.NODE_ENV = "test";
      // Mock getAuthHeaders to return only Google Identity token for anonymous requests
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({
        "X-Serverless-Authorization": "Bearer google-identity-token",
      });

      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true },
      });

      const response = await request(app).post("/api/public-endpoint").send({ data: "test" });

      expect(response.status).toBe(200);
      expect(authModule.getAuthHeaders).toHaveBeenCalledWith(undefined);
      expect(axios).toHaveBeenCalledWith({
        method: "POST",
        url: `${process.env.BACKEND_URL}/public-endpoint`,
        data: { data: "test" },
        headers: {
          "X-Serverless-Authorization": "Bearer google-identity-token",
          "Content-Type": "application/json",
        },
        maxBodyLength: MAX_JSON_BYTES,
        maxContentLength: MAX_JSON_BYTES,
      });
    });

    it("should proxy multipart requests with MAX_MULTIPART_BYTES limits", async () => {
      process.env.NODE_ENV = "test";
      vi.spyOn(authModule, "getAuthHeaders").mockResolvedValue({
        Authorization: "Bearer user-token",
        "X-Serverless-Authorization": "Bearer google-token",
      });

      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true },
      });

      const response = await request(app)
        .post("/api/file-to-plan")
        .set("Authorization", "Bearer user-token")
        .attach("file", Buffer.from("dummy-content"), "plan.png");

      expect(response.status).toBe(200);
      expect(axios).toHaveBeenCalledWith(
        expect.objectContaining({
          method: "POST",
          url: `${process.env.BACKEND_URL}/file-to-plan`,
          maxBodyLength: MAX_MULTIPART_BYTES,
          maxContentLength: MAX_MULTIPART_BYTES,
        }),
      );
    });
  });

  describe("Security", () => {
    it("should allow requests from the configured frontend URL", async () => {
      const response = await request(app)
        .get("/api/some-endpoint")
        .set("Origin", process.env.FRONTEND_URL as string);
      expect(response.headers["access-control-allow-origin"]).toBe(process.env.FRONTEND_URL);
      expect(response.headers["access-control-allow-credentials"]).toBeUndefined();
    });

    it("should rate limit requests", async () => {
      // Mock successful responses for the first 100 requests
      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true },
      });

      const clientIp = "192.0.2.99";
      for (let i = 0; i < 100; i++) {
        await request(app).get("/api/some-endpoint").set("X-Forwarded-For", clientIp).expect(200);
      }

      // The 101st request should be rate limited
      const response = await request(app)
        .get("/api/some-endpoint")
        .set("X-Forwarded-For", clientIp);
      expect(response.status).toBe(429);
    }, 15000);

    it("should resolve client IP across multi-hop X-Forwarded-For headers", async () => {
      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true },
      });

      // 100 requests with multi-hop X-Forwarded-For from IP 198.51.100.1
      for (let i = 0; i < 10; i++) {
        await Promise.all(
          Array.from({ length: 10 }, () =>
            request(app)
              .get("/api/test-hops")
              .set("X-Forwarded-For", "198.51.100.1, 10.0.0.1, 10.0.0.2")
              .expect(200),
          ),
        );
      }

      // 101st request from 198.51.100.1 is rate limited
      const blocked = await request(app)
        .get("/api/test-hops")
        .set("X-Forwarded-For", "198.51.100.1, 10.0.0.1, 10.0.0.2");
      expect(blocked.status).toBe(429);

      // A different client IP 198.51.100.2 through the same proxy chain is NOT rate limited
      const allowed = await request(app)
        .get("/api/test-hops")
        .set("X-Forwarded-For", "198.51.100.2, 10.0.0.1, 10.0.0.2");
      expect(allowed.status).toBe(200);
    });

    it("should isolate rate limits per authenticated user even on shared IP", async () => {
      vi.mocked(axios).mockResolvedValue({
        status: 200,
        data: { success: true },
      });

      const makeJwt = (userId: string) => {
        const payload = Buffer.from(JSON.stringify({ sub: userId })).toString("base64url");
        return `Bearer header.${payload}.sig`;
      };

      const user1Token = makeJwt("user-111");
      const user2Token = makeJwt("user-222");

      // Exhaust limit for user 1
      for (let i = 0; i < 10; i++) {
        await Promise.all(
          Array.from({ length: 10 }, () =>
            request(app)
              .get("/api/user-test")
              .set("Authorization", user1Token)
              .set("X-Forwarded-For", "203.0.113.50")
              .expect(200),
          ),
        );
      }

      // 101st request from user 1 should be blocked
      const user1Blocked = await request(app)
        .get("/api/user-test")
        .set("Authorization", user1Token)
        .set("X-Forwarded-For", "203.0.113.50");
      expect(user1Blocked.status).toBe(429);

      // User 2 from the same IP should still be allowed
      const user2Allowed = await request(app)
        .get("/api/user-test")
        .set("Authorization", user2Token)
        .set("X-Forwarded-For", "203.0.113.50");
      expect(user2Allowed.status).toBe(200);
    });
  });
});
