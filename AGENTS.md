# swim-gen - AI Swimming Training Plan Generator

**Tech Stack**: Vue 3 + Express BFF + Go Backend + Supabase + Cloud Run  
**RAG System**: Langchain + pgvector + Gemini embeddings  
**Observability**: OpenTelemetry → Google Cloud Trace

## Project Structure

```
swim-gen/
├── frontend/          # Vue 3 + TypeScript + Vite
│   ├── src/
│   │   ├── stores/    # Pinia state (10 stores) - see ./frontend/src/stores/AGENTS.md
│   │   ├── views/     # Route components (12 views) - see ./frontend/src/views/AGENTS.md
│   │   ├── components/
│   │   ├── api/       # ApiClient singleton (axios wrapper)
│   │   └── types/     # TypeScript models mirroring backend
│   └── __tests__/     # Vitest + jsdom, colocated with src/
├── bff/               # Node.js Express proxy (rate limiting, auth forwarding)
│   └── src/
│       ├── main.ts
│       └── instrumentation.ts  # OTEL setup (MUST load first)
├── backend/           # Go + Chi + pgvector RAG
│   ├── main.go
│   └── internal/
│       ├── server/    # HTTP handlers (27 funcs) - see ./backend/internal/server/AGENTS.md
│       ├── rag/       # RAG core (53 funcs) - see ./backend/internal/rag/AGENTS.md
│       ├── genai/     # Vertex AI client
│       ├── models/    # Shared structs
│       └── telemetry/ # OTEL setup
└── deployments/       # Terraform/OpenTofu - see ./deployments/AGENTS.md
    ├── dev/
    │   ├── 0-infra/   # IAM, secrets, WIF
    │   └── 1-services/# Cloud Run
    └── prod/          # Same structure
```

## Code Map (High-Level Flow)

### User Journey
1. **Auth**: Frontend → Supabase (JWT) → BFF forwards token → Backend validates with Supabase middleware
2. **Plan Generation**: 
   - User submits query (TrainingPlanView) → store action → API call
   - BFF rate-limits → Backend receives → RAG.Query() → Gemini generates plan → pgvector retrieves drills
   - Response flows back → store updates → UI re-renders
3. **Chat**: Similar flow, but uses RAG.ChatWithContext() with conversation history from memory store
4. **Export**: Frontend → BFF → Backend generates PDF (chromedp) → Cloud Storage → signed URL

### Data Flow
```
Frontend Store (Pinia)
    ↓ (axios via ApiClient)
BFF (Express + rate-limit-flexible)
    ↓ (adds service token)
Backend Chi Router (otelhttp middleware)
    ↓ (Supabase auth middleware validates JWT)
Handler (server/)
    ↓
RAG Core (rag/)
    ↓ (pgx + pgxscan)
PostgreSQL + pgvector
    ↓ (langchain embeddings)
Vertex AI (Gemini)
```

### Tracing Flow
```
Frontend (no tracing)
    ↓ (HTTP request)
BFF (OTEL auto-instrumentation)
    ↓ (W3C trace context in headers: traceparent, tracestate)
Backend (otelhttp.NewHandler wraps Chi router)
    ↓ (trace context extracted → spans linked)
Google Cloud Trace (OTLP gRPC exporter)
```

## Conventions & Patterns

### Frontend
- **Composition API**: `<script setup lang="ts">`, no Options API
- **Path Alias**: `@/` maps to `src/`
- **Tests**: Vitest with jsdom, colocated in `__tests__/` directories
- **Stores**: Pinia with composition style (`defineStore('name', () => { ... })`)
  - **CRITICAL DUPLICATION**: `ensureRowIds()` utility duplicated across `trainingPlan.ts`, `sharedPlan.ts`, `drills.ts` (see stores/AGENTS.md)
- **API Client**: Singleton pattern in `api/client.ts`, DO NOT create new axios instances
- **Type Safety**: Models in `types/` mirror backend Go structs exactly

### BFF (Node.js)
- **Instrumentation Order**: `instrumentation.ts` MUST be imported before anything else (see `--require` in package.json)
- **Rate Limiting**: MemoryStore (not Redis), 100 req/15min per IP
- **Auth**: Dual tokens: user JWT (forwarded) + service token (added by BFF)
- **OTEL**: Auto-instrumentation for Express, HTTP, Axios via `@opentelemetry/auto-instrumentations-node`

### Backend (Go)
- **Router**: Chi with `otelhttp.NewHandler()` wrapping root handler
- **Middleware Order**: CORS → httplog → OTEL → Supabase auth → handler
- **Context**: Always pass `ctx context.Context` from request, DO NOT create background contexts
- **Error Handling**: Return errors up the stack, log at handler level only
- **Telemetry**: `telemetry.InitTelemetry()` MUST be called in `main()` before router setup
- **RAG**: All database operations through `RAGDB` struct methods (see rag/AGENTS.md)
- **Testing**: testify/assert, table-driven tests preferred

### Infrastructure (Terraform)
- **State**: Local state (COMMITTED - KNOWN ISSUE), should use GCS backend
- **Structure**: `0-infra/` (foundational: IAM, secrets) → `1-services/` (Cloud Run, dependent on 0-infra outputs)
- **WIF**: Workload Identity Federation for GitHub Actions (no service account keys)
- **Dual Env**: `dev/` (100% trace sampling) vs `prod/` (10% sampling)

## Anti-Patterns (Project-Specific)

### Frontend
- ❌ **DO NOT** suppress TypeScript errors with `as any`, `@ts-ignore`, `@ts-expect-error`
- ❌ **DO NOT** create new axios instances, use `ApiClient.getInstance()`
- ❌ **DO NOT** mutate store state directly, use store actions
- ❌ **DO NOT** duplicate `ensureRowIds()` further - refactor to shared utility instead

### BFF
- ❌ **DO NOT** skip instrumentation preload (breaks OTEL)
- ❌ **DO NOT** add authentication logic here (BFF is a proxy, not an auth service)

### Backend
- ❌ **DO NOT** use `context.Background()` in handlers, always use `req.Context()`
- ❌ **DO NOT** suppress errors with empty `if err != nil { }` blocks
- ❌ **DO NOT** commit test databases or credentials (use `.env` with `.gitignore`)
- ❌ **DO NOT** bypass Supabase middleware for protected routes

### Infrastructure
- ❌ **DO NOT** hardcode secrets in `.tf` files, use Secret Manager
- ❌ **DO NOT** commit `.terraform/` directories or `terraform.tfstate` (KNOWN ISSUE)
- ❌ **DO NOT** skip `0-infra` when deploying `1-services` (dependencies will break)

## Build & Run

### Local Development
```bash
# Frontend
cd frontend && npm install && npm run dev  # localhost:5173

# BFF (requires instrumentation preload)
cd bff && npm install && npm run dev       # localhost:3001

# Backend (requires .env with DB credentials)
cd backend && go run main.go               # localhost:8080

# Docker Compose (full stack)
docker-compose up                          # frontend:5173, bff:3001, backend:8080
```

### Testing
```bash
# Frontend
cd frontend && npm test                    # Vitest

# Backend
cd backend && go test ./...                # Go test + testify
```

### Deployment
```bash
# Triggered by GitHub comment: `.deploy dev` or `.deploy prod`
# CI/CD: .github/workflows/deploy-dev.yml, deploy-prod.yml
# Process:
#   1. Build Docker images (bff, backend)
#   2. Push to Artifact Registry
#   3. Apply Terraform: cd deployments/{env}/0-infra && tofu apply
#   4. Apply Terraform: cd deployments/{env}/1-services && tofu apply
#   5. Update Cloud Run with new images
```

### OTEL Verification
```bash
# Check traces in Google Cloud Console
# https://console.cloud.google.com/traces/list?project={PROJECT_ID}

# Verify trace context propagation:
# 1. Generate a plan in frontend
# 2. Find trace in Cloud Trace
# 3. Verify spans: BFF (Express) → Backend (Chi handler) → RAG methods
# 4. Check span attributes: http.method, http.route, http.status_code
```

## Gotchas & Known Issues

### 1. Build Artifacts Committed (CRITICAL)
- **Issue**: `frontend/dist/`, `bff/dist/`, `bff/node_modules/` tracked in git
- **Why**: Likely accident, not intentional deployment strategy
- **Fix**: Add to `.gitignore`, `git rm -r --cached`, commit cleanup

### 2. Terraform State Committed (CRITICAL)
- **Issue**: `.terraform/` directories and `terraform.tfstate` tracked in git
- **Why**: State should be remote (GCS backend), not local
- **Fix**: Configure GCS backend in `backend.tf`, migrate state, `.gitignore` cleanup

### 3. Instrumentation Load Order (BFF)
- **Issue**: If `instrumentation.ts` loads after app code, OTEL won't instrument
- **Solution**: `package.json` scripts use `--require ./dist/instrumentation.js` OR import at top of `main.ts`
- **Verify**: Check `npm run dev` and `npm run start` scripts

### 4. Trace Context Loss
- **Issue**: If BFF doesn't propagate `traceparent` header, backend traces are disconnected
- **Solution**: OTEL auto-instrumentation handles this IF loaded correctly
- **Debug**: Check Network tab for `traceparent` header in BFF → Backend requests

### 5. Large Files (>500 lines)
- **Files**: 12 files >500 lines (mostly in frontend stores, backend RAG)
- **Impact**: Complexity hotspots, harder to test/maintain
- **Refactor Priority**: `backend/internal/rag/scraper.go` (466 lines), `frontend/src/stores/trainingPlan.ts` (450+ lines)

### 6. IAM Permissions (Tracing)
- **Issue**: New Cloud Trace roles added to Terraform but NOT deployed yet
- **Symptom**: Traces fail to export with 403 Forbidden
- **Fix**: Deploy `deployments/{dev,prod}/0-infra` to apply IAM changes

### 7. Duplicated Utilities
- **Issue**: `ensureRowIds()` duplicated in `trainingPlan.ts`, `sharedPlan.ts`, `drills.ts`
- **Impact**: Bug fixes must be applied 3x, inconsistency risk
- **Refactor**: Extract to `frontend/src/utils/table.ts`, import in stores

### 8. Table Manipulation Logic
- **Issue**: Complex table state management spread across stores (add/delete rows, reorder, etc.)
- **Pattern**: Each store reimplements similar logic for different table types (weeks, days, sets)
- **Consider**: Shared composable or store mixin for common table operations

### 9. Comment-Triggered Deploys
- **Trigger**: `.deploy dev` or `.deploy prod` in PR comments
- **Caveat**: Only works on PRs with write access, not forks
- **Workflow**: `.github/workflows/deploy-{dev,prod}.yml`

### 10. Supabase Auth Token Validation
- **Issue**: Backend validates JWTs against Supabase, requires network call
- **Performance**: ~50-100ms per request, consider caching if traffic scales
- **Alternative**: JWT signature verification with public key (faster, but more complex)

## Tips for AI Assistants

### When Modifying Stores
1. Read existing store first to understand pattern (composition API + Pinia)
2. Check if `ensureRowIds()` is used → consider shared utility refactor
3. Add tests in `__tests__/` directory (Vitest + jsdom)
4. Update TypeScript types in `frontend/src/types/`

### When Adding Backend Handlers
1. Add handler function in `backend/internal/server/{domain}.go`
2. Register route in `main.go` Chi router with `otelhttp.NewHandler()` wrapper
3. Use `req.Context()` for all downstream calls
4. Add tests in `{domain}_test.go` with testify/assert
5. Update API client in `frontend/src/api/client.ts`

### When Changing Infrastructure
1. ALWAYS modify `dev/` and `prod/` in parallel (keep parity)
2. Apply `0-infra` before `1-services` (dependencies)
3. Test in `dev` first, then promote to `prod`
4. Document env-specific differences (e.g., trace sampling rates)

### When Debugging Traces
1. Check BFF instrumentation loaded: look for OTEL logs in `npm run dev` output
2. Verify `traceparent` header in BFF → Backend requests (Network tab)
3. Check Cloud Trace UI for disconnected spans (indicates context loss)
4. Verify IAM roles: `roles/cloudtrace.agent` on service accounts

### When Refactoring
1. Check for duplicates FIRST (grep for function names, patterns)
2. Run tests BEFORE and AFTER refactoring (ensure green → green)
3. Use LSP diagnostics to catch TypeScript/Go errors early
4. Update AGENTS.md if patterns change significantly
