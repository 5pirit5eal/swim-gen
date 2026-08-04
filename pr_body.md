## Summary
Automated dependency update performed on 2026-08-01.

## Packages Updated
### Backend (Go)
| Package | Previous | New |
|---------|----------|-----|
| No updates |

### BFF (Node.js)
| Package | Previous | New | Type |
|---------|----------|-----|------|
| @opentelemetry/sdk-trace-node | ^2.8.0 | ^2.10.0 | dependency |
| @opentelemetry/semantic-conventions | ^1.41.1 | ^1.43.0 | dependency |
| axios | ^1.18.1 | ^1.19.0 | dependency |
| express-rate-limit | ^8.5.2 | ^8.6.1 | dependency |
| google-auth-library | ^10.9.0 | ^10.9.1 | dependency |
| @types/node | ^24.13.2 | ^24.13.3 | devDependency |
| @vitest/eslint-plugin | ^1.6.20 | ^1.6.24 | devDependency |
| eslint | ^9.39.4 | ^9.39.5 | devDependency |
| prettier | ^3.9.4 | ^3.9.6 | devDependency |
| typescript-eslint | ^8.62.1 | ^8.65.0 | devDependency |

### Frontend (Node.js)
| Package | Previous | New | Type |
|---------|----------|-----|------|
| @supabase/supabase-js | ^2.110.0 | ^2.111.0 | dependency |
| driver.js | ^1.6.0 | ^1.8.0 | dependency |
| vue | ^3.5.25 | ^3.5.40 | dependency |
| vue-i18n | ^11.4.6 | ^11.4.8 | dependency |
| @types/node | ^24.13.2 | ^24.13.3 | devDependency |
| @vitejs/plugin-vue | ^6.0.7 | ^6.0.8 | devDependency |
| @vitest/eslint-plugin | ^1.6.20 | ^1.6.24 | devDependency |
| eslint | ^9.39.4 | ^9.39.5 | devDependency |
| svgo | ^4.0.1 | ^4.0.2 | devDependency |
| vite-plugin-vue-devtools | ^8.1.5 | ^8.2.1 | devDependency |
| vue-tsc | ^3.3.6 | ^3.3.9 | devDependency |

## Breaking Changes
None

## Validation Status
- Backend (go vet + go test): SKIPPED - no changes
- BFF (type-check + lint + test): PASS
- Frontend (type-check + lint + test): PASS

## Notes
All tests and static analysis passed after updates.
