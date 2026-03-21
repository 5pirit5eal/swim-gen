# Issues — training-plan-display-style

## [2026-03-21] Pre-existing LSP Errors (NOT caused by our work)

### TrainingPlanDisplay.spec.ts
- ERROR [5:33] Cannot find module '../TrainingPlanDisplay.vue' or its corresponding type declarations
- ERROR [87:7] Unused '@ts-expect-error' directive
- NOTE: These are PRE-EXISTING errors. The component likely does not type-export; this is expected in Vue + Vite setups and may not affect runtime.

### router/index.ts
- Multiple "Cannot find module" errors for all view files
- NOTE: These are PRE-EXISTING LSP errors typical in Vite/Vue projects where the TS language server doesn't resolve `.vue` paths the same way Vite does. They do NOT indicate actual broken imports.
- These errors exist BEFORE any of our changes.
