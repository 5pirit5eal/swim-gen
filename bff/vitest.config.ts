import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    pool: "forks",
    globals: true,
    include: ["src/**/__tests__/**/*.spec.ts"],
    exclude: ["dist/**", "node_modules/**"],
  },
});


