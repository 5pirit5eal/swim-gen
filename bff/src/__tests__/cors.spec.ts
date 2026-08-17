import { describe, it, expect, beforeEach, afterEach } from "vitest";
import request from "supertest";
import { app } from "../main";
import { normalizeOrigin, getAllowedOrigins, isOriginAllowed } from "../cors";

describe("CORS Configuration", () => {
  const originalFrontendUrl = process.env.FRONTEND_URL;

  afterEach(() => {
    process.env.FRONTEND_URL = originalFrontendUrl;
  });

  describe("normalizeOrigin", () => {
    it("should prepend https:// if scheme is missing", () => {
      expect(normalizeOrigin("swim-gen.com")).toBe("https://swim-gen.com");
      expect(normalizeOrigin("app.swim-gen.com")).toBe("https://app.swim-gen.com");
    });

    it("should preserve existing http:// or https:// schemes", () => {
      expect(normalizeOrigin("http://localhost:5173")).toBe("http://localhost:5173");
      expect(normalizeOrigin("https://swim-gen.com")).toBe("https://swim-gen.com");
    });

    it("should strip trailing slashes", () => {
      expect(normalizeOrigin("https://swim-gen.com/")).toBe("https://swim-gen.com");
      expect(normalizeOrigin("http://localhost:5173///")).toBe("http://localhost:5173");
      expect(normalizeOrigin("swim-gen.com/")).toBe("https://swim-gen.com");
    });

    it("should normalize hostname to lowercase", () => {
      expect(normalizeOrigin("https://SWIM-GEN.COM")).toBe("https://swim-gen.com");
      expect(normalizeOrigin("HTTP://LOCALHOST:5173")).toBe("http://localhost:5173");
    });

    it("should preserve ports", () => {
      expect(normalizeOrigin("http://localhost:8080")).toBe("http://localhost:8080");
      expect(normalizeOrigin("https://example.com:8443")).toBe("https://example.com:8443");
    });

    it("should handle wildcard *", () => {
      expect(normalizeOrigin("*")).toBe("*");
    });
  });

  describe("getAllowedOrigins", () => {
    it("should return empty array when FRONTEND_URL is unset or empty", () => {
      delete process.env.FRONTEND_URL;
      expect(getAllowedOrigins()).toEqual([]);
      expect(getAllowedOrigins("")).toEqual([]);
      expect(getAllowedOrigins("   ")).toEqual([]);
    });

    it("should parse single origin with or without scheme", () => {
      expect(getAllowedOrigins("swim-gen.com")).toEqual(["https://swim-gen.com"]);
      expect(getAllowedOrigins("https://swim-gen.com")).toEqual(["https://swim-gen.com"]);
    });

    it("should parse comma-separated list of origins with trimming", () => {
      const input = " https://swim-gen.com , http://localhost:5173 , http://frontend:8080/ ";
      expect(getAllowedOrigins(input)).toEqual([
        "https://swim-gen.com",
        "http://localhost:5173",
        "http://frontend:8080",
      ]);
    });
  });

  describe("isOriginAllowed", () => {
    const allowed = ["https://swim-gen.com", "http://localhost:5173"];

    it("should allow requests without origin header (server-to-server, curl)", () => {
      expect(isOriginAllowed(undefined, allowed)).toBe(true);
      expect(isOriginAllowed("", allowed)).toBe(true);
    });

    it("should allow exact matched origins", () => {
      expect(isOriginAllowed("https://swim-gen.com", allowed)).toBe(true);
      expect(isOriginAllowed("http://localhost:5173", allowed)).toBe(true);
    });

    it("should allow matching origins with trailing slashes", () => {
      expect(isOriginAllowed("https://swim-gen.com/", allowed)).toBe(true);
    });

    it("should allow matching origins regardless of case", () => {
      expect(isOriginAllowed("HTTPS://SWIM-GEN.COM", allowed)).toBe(true);
    });

    it("should reject untrusted domains", () => {
      expect(isOriginAllowed("https://malicious.com", allowed)).toBe(false);
      expect(isOriginAllowed("http://evil-swim-gen.com", allowed)).toBe(false);
    });

    it("should reject subdomain hijacking attempts", () => {
      expect(isOriginAllowed("https://swim-gen.com.attacker.com", allowed)).toBe(false);
      expect(isOriginAllowed("https://not-swim-gen.com", allowed)).toBe(false);
    });

    it("should reject scheme mismatch", () => {
      expect(isOriginAllowed("http://swim-gen.com", allowed)).toBe(false);
      expect(isOriginAllowed("https://localhost:5173", allowed)).toBe(false);
    });

    it("should reject port mismatch", () => {
      expect(isOriginAllowed("https://swim-gen.com:8080", allowed)).toBe(false);
      expect(isOriginAllowed("http://localhost:3000", allowed)).toBe(false);
    });

    it("should reject 'null' origins", () => {
      expect(isOriginAllowed("null", allowed)).toBe(false);
    });
  });

  describe("HTTP CORS Integration", () => {
    beforeEach(() => {
      process.env.FRONTEND_URL = "https://swim-gen.com";
    });

    it("should handle OPTIONS preflight for allowed production origin", async () => {
      const res = await request(app)
        .options("/health")
        .set("Origin", "https://swim-gen.com")
        .set("Access-Control-Request-Method", "POST")
        .set("Access-Control-Request-Headers", "Content-Type,Authorization");

      expect(res.status).toBe(200);
      expect(res.headers["access-control-allow-origin"]).toBe("https://swim-gen.com");
      expect(res.headers["access-control-allow-methods"]).toBe("GET,POST,PUT,DELETE,OPTIONS");
      expect(res.headers["access-control-allow-headers"]).toBe("Content-Type,Authorization");
      expect(res.headers["access-control-allow-credentials"]).toBeUndefined();
    });

    it("should emit Access-Control-Allow-Origin on actual GET requests for allowed origin", async () => {
      const res = await request(app).get("/health").set("Origin", "https://swim-gen.com");

      expect(res.status).toBe(200);
      expect(res.headers["access-control-allow-origin"]).toBe("https://swim-gen.com");
      expect(res.headers["access-control-allow-credentials"]).toBeUndefined();
    });

    it("should NOT emit Access-Control-Allow-Origin for untrusted origins on preflight", async () => {
      const res = await request(app)
        .options("/health")
        .set("Origin", "https://attacker.com")
        .set("Access-Control-Request-Method", "POST");

      expect(res.headers["access-control-allow-origin"]).toBeUndefined();
      expect(res.headers["access-control-allow-credentials"]).toBeUndefined();
    });

    it("should NOT emit Access-Control-Allow-Origin for untrusted origins on actual requests", async () => {
      const res = await request(app).get("/health").set("Origin", "https://attacker.com");

      expect(res.status).toBe(200);
      expect(res.headers["access-control-allow-origin"]).toBeUndefined();
    });

    it("should NOT emit Access-Control-Allow-Origin for 'null' origin", async () => {
      const res = await request(app).get("/health").set("Origin", "null");

      expect(res.status).toBe(200);
      expect(res.headers["access-control-allow-origin"]).toBeUndefined();
    });

    it("should NOT emit Access-Control-Allow-Origin for subdomain spoofing", async () => {
      const res = await request(app).get("/health").set("Origin", "https://swim-gen.com.evil.com");

      expect(res.status).toBe(200);
      expect(res.headers["access-control-allow-origin"]).toBeUndefined();
    });

    it("should allow requests with no Origin header (e.g. server-to-server)", async () => {
      const res = await request(app).get("/health");

      expect(res.status).toBe(200);
      expect(res.text).toBe("OK");
      expect(res.headers["access-control-allow-origin"]).toBeUndefined();
    });

    it("should support multiple origins when configured as comma-separated list", async () => {
      process.env.FRONTEND_URL = "https://swim-gen.com, http://localhost:5173";

      const resProd = await request(app).get("/health").set("Origin", "https://swim-gen.com");
      expect(resProd.headers["access-control-allow-origin"]).toBe("https://swim-gen.com");

      const resLocal = await request(app).get("/health").set("Origin", "http://localhost:5173");
      expect(resLocal.headers["access-control-allow-origin"]).toBe("http://localhost:5173");

      const resOther = await request(app).get("/health").set("Origin", "http://localhost:3000");
      expect(resOther.headers["access-control-allow-origin"]).toBeUndefined();
    });

    it("should normalize FRONTEND_URL if scheme was omitted in env variable", async () => {
      process.env.FRONTEND_URL = "swim-gen.com";

      const res = await request(app).get("/health").set("Origin", "https://swim-gen.com");

      expect(res.headers["access-control-allow-origin"]).toBe("https://swim-gen.com");
    });
  });
});
