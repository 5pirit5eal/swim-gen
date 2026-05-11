/**
 * Logger utility that outputs structured JSON on Cloud Run (detected via K_SERVICE env var)
 * and human-readable text locally.
 *
 * Supported log levels (in order of severity):
 * - DEBUG: Detailed debugging information
 * - INFO: General informational messages (default)
 * - WARN: Warning messages
 * - ERROR: Error messages only
 */

export type LogLevel = "DEBUG" | "INFO" | "WARN" | "ERROR";

const LOG_LEVELS: Record<LogLevel, number> = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
};

// Map log levels to Cloud Logging severity names
const SEVERITY_MAP: Record<LogLevel, string> = {
  DEBUG: "DEBUG",
  INFO: "INFO",
  WARN: "WARNING",
  ERROR: "ERROR",
};

const isCloudRun = !!process.env.K_SERVICE;

function getConfiguredLogLevel(): LogLevel {
  const envLevel = process.env.LOG_LEVEL?.toUpperCase();
  if (envLevel && envLevel in LOG_LEVELS) {
    return envLevel as LogLevel;
  }
  return "INFO"; // Default to INFO
}

function shouldLog(messageLevel: LogLevel): boolean {
  const configuredLevel = getConfiguredLogLevel();
  return LOG_LEVELS[messageLevel] >= LOG_LEVELS[configuredLevel];
}

function formatMessage(level: LogLevel, message: string, ...args: unknown[]): string {
  const timestamp = new Date().toISOString();
  const formattedArgs = args.length > 0 ? ` ${JSON.stringify(args)}` : "";
  return `[${timestamp}] [${level}] ${message}${formattedArgs}`;
}

function writeStructured(level: LogLevel, message: string, ...args: unknown[]): void {
  if (!shouldLog(level)) return;

  if (isCloudRun) {
    const entry: Record<string, unknown> = {
      severity: SEVERITY_MAP[level],
      message,
      timestamp: new Date().toISOString(),
    };
    if (args.length > 0) {
      entry.data = args;
    }
    const writer = level === "ERROR" ? process.stderr : process.stdout;
    writer.write(JSON.stringify(entry) + "\n");
  } else {
    const writer =
      level === "ERROR"
        ? console.error
        : level === "WARN"
          ? console.warn
          : level === "DEBUG"
            ? console.debug
            : console.info;
    writer(formatMessage(level, message, ...args));
  }
}

export const logger = {
  debug(message: string, ...args: unknown[]): void {
    writeStructured("DEBUG", message, ...args);
  },

  info(message: string, ...args: unknown[]): void {
    writeStructured("INFO", message, ...args);
  },

  warn(message: string, ...args: unknown[]): void {
    writeStructured("WARN", message, ...args);
  },

  error(message: string, ...args: unknown[]): void {
    writeStructured("ERROR", message, ...args);
  },

  /**
   * Write a pre-built structured log entry directly to stdout as JSON.
   * Used by the request logging middleware for full control over fields.
   * Falls back to console.log in non-Cloud Run environments.
   */
  structured(entry: Record<string, unknown>): void {
    if (isCloudRun) {
      process.stdout.write(JSON.stringify(entry) + "\n");
    } else {
      const { severity, message, ...rest } = entry;
      const level = (severity as string) || "INFO";
      const fields = Object.entries(rest)
        .filter(([, v]) => v !== undefined)
        .map(([k, v]) => `${k}=${JSON.stringify(v)}`)
        .join(" ");
      console.log(`[${level}] ${message || ""} ${fields}`);
    }
  },

  /**
   * Returns the currently configured log level
   */
  getLevel(): LogLevel {
    return getConfiguredLogLevel();
  },
};
