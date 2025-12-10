import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { logger } from "../logger";

describe("Logger", () => {
  const originalLogLevel = process.env.LOG_LEVEL;

  beforeEach(() => {
    vi.spyOn(console, "debug").mockImplementation(() => {});
    vi.spyOn(console, "info").mockImplementation(() => {});
    vi.spyOn(console, "warn").mockImplementation(() => {});
    vi.spyOn(console, "error").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
    process.env.LOG_LEVEL = originalLogLevel;
  });

  describe("getLevel", () => {
    it("should return INFO by default", () => {
      delete process.env.LOG_LEVEL;
      expect(logger.getLevel()).toBe("INFO");
    });

    it("should return the configured log level", () => {
      process.env.LOG_LEVEL = "DEBUG";
      expect(logger.getLevel()).toBe("DEBUG");
    });

    it("should be case insensitive", () => {
      process.env.LOG_LEVEL = "debug";
      expect(logger.getLevel()).toBe("DEBUG");
    });

    it("should fall back to INFO for invalid values", () => {
      process.env.LOG_LEVEL = "INVALID";
      expect(logger.getLevel()).toBe("INFO");
    });
  });

  describe("log level filtering", () => {
    it("should log all levels when LOG_LEVEL is DEBUG", () => {
      process.env.LOG_LEVEL = "DEBUG";

      logger.debug("debug message");
      logger.info("info message");
      logger.warn("warn message");
      logger.error("error message");

      expect(console.debug).toHaveBeenCalled();
      expect(console.info).toHaveBeenCalled();
      expect(console.warn).toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });

    it("should not log DEBUG when LOG_LEVEL is INFO", () => {
      process.env.LOG_LEVEL = "INFO";

      logger.debug("debug message");
      logger.info("info message");
      logger.warn("warn message");
      logger.error("error message");

      expect(console.debug).not.toHaveBeenCalled();
      expect(console.info).toHaveBeenCalled();
      expect(console.warn).toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });

    it("should only log WARN and ERROR when LOG_LEVEL is WARN", () => {
      process.env.LOG_LEVEL = "WARN";

      logger.debug("debug message");
      logger.info("info message");
      logger.warn("warn message");
      logger.error("error message");

      expect(console.debug).not.toHaveBeenCalled();
      expect(console.info).not.toHaveBeenCalled();
      expect(console.warn).toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });

    it("should only log ERROR when LOG_LEVEL is ERROR", () => {
      process.env.LOG_LEVEL = "ERROR";

      logger.debug("debug message");
      logger.info("info message");
      logger.warn("warn message");
      logger.error("error message");

      expect(console.debug).not.toHaveBeenCalled();
      expect(console.info).not.toHaveBeenCalled();
      expect(console.warn).not.toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });
  });

  describe("message formatting", () => {
    it("should include timestamp, level, and message", () => {
      process.env.LOG_LEVEL = "DEBUG";

      logger.info("test message");

      expect(console.info).toHaveBeenCalledWith(
        expect.stringMatching(
          /^\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z\] \[INFO\] test message$/,
        ),
      );
    });

    it("should include additional arguments as JSON", () => {
      process.env.LOG_LEVEL = "DEBUG";

      logger.info("test message", { key: "value" });

      expect(console.info).toHaveBeenCalledWith(
        expect.stringMatching(/\[INFO\] test message \[{"key":"value"}\]$/),
      );
    });
  });
});
