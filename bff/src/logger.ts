/**
 * Simple logger utility that respects the LOG_LEVEL environment variable.
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

export const logger = {
  debug(message: string, ...args: unknown[]): void {
    if (shouldLog("DEBUG")) {
      console.debug(formatMessage("DEBUG", message, ...args));
    }
  },

  info(message: string, ...args: unknown[]): void {
    if (shouldLog("INFO")) {
      console.info(formatMessage("INFO", message, ...args));
    }
  },

  warn(message: string, ...args: unknown[]): void {
    if (shouldLog("WARN")) {
      console.warn(formatMessage("WARN", message, ...args));
    }
  },

  error(message: string, ...args: unknown[]): void {
    if (shouldLog("ERROR")) {
      console.error(formatMessage("ERROR", message, ...args));
    }
  },

  /**
   * Returns the currently configured log level
   */
  getLevel(): LogLevel {
    return getConfiguredLogLevel();
  },
};
