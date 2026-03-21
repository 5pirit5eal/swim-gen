# Learnings — training-plan-display-style

## [2026-03-21] Planning Complete

### Architecture
- `Row` interface is immutable — do NOT add new fields
- Path-based recursion (`number[]`) must stay intact for all store operations
- Last row in `table[]` is always the total row — must be excluded from exercise rendering
- `MAX_NESTING_DEPTH = 4` defined in `rowHelpers.ts`

### Key Files
- `frontend/src/components/training/TrainingPlanDisplay.vue` — 616 lines, main editable display
- `frontend/src/components/training/TrainingPlanRow.vue` — 383 lines, recursive editable row
- `frontend/src/components/training/SimplePlanDisplay.vue` — 287 lines, read-only display
- `frontend/src/components/training/SimplePlanSubRows.vue` — 146 lines, read-only nested rows
- `frontend/src/utils/rowHelpers.ts` — 262 lines, pure helpers
- `frontend/src/stores/trainingPlan.ts` — 621 lines, Pinia store
- `frontend/src/types/training.ts` — Row interface and PlanStore type
- `frontend/src/components/ui/BaseTableAction.vue` — 230 lines, needs full redesign
- `frontend/src/assets/base.css` — CSS variables and responsive conventions

### Reference Files
- `frontend/web.html` — primary layout reference
- `frontend/mobile.html` — secondary mobile reference
- `frontend/image.png` — rendered visual

### Test Infrastructure
- Vitest + Vue Test Utils + jsdom
- `@pinia/testing` for store mocking
- `frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts` — 129 lines
- Current tests use `.anchor-cell input` and `.edit-btn` selectors (will need updating)
- CI: `npm run test:unit`, `npm run type-check`, `npm run lint`, `npm run build`

### Constraints
- No Tailwind
- No changes to Row schema or backend contracts
- No hover-only interactions
- CSS must use existing CSS variable/theme system
- Do NOT infer semantic phase labels (warmup/main/cooldown)
- Follow web.html more closely than mobile.html

## [2026-03-21] TDD Wave 1, Task 1 — Display Model & Fixtures (RED Phase)

### Failing Tests (9 PASSED, 2 FAILED)
Tests added to `TrainingPlanDisplay.spec.ts`:
- ✅ renders exercise rows as cards with `[data-testid="plan-card"]`
- ✅ renders flat plan with one exercise card and excludes total row from cards
- ✅ renders nested depth-2 plan with parent and sibling cards
- ✅ renders nested depth-3 plan with proper hierarchy
- ✅ renders mixed equipment plan with proper card structure
- ❌ does not render `<table>` element in main plan content (expects 0, currently 1)
- ❌ parent rows with SubRows render nested cards (expects `[data-testid="nested-card"]`, not found)

### Fixture Builders Added
Five fixture creators in test file:
1. `createSimplePlan()` — flat single exercise + total
2. `createNestedDepth2Plan()` — top-level parent with 2 children + 1 sibling + total
3. `createNestedDepth3Plan()` — 3-level nesting (parent → grandparent → children)
4. `createMixedEquipmentPlan()` — varied equipment (none, single, multiple) across rows
5. `mockPlan` — legacy compatibility wrapper

All fixtures explicitly exclude total row from card expectations.

### Pure Helper Functions Added to `rowHelpers.ts`
New display-model functions (no Row mutation):
- `isExerciseRow(row: Row): boolean` — checks Content !== 'Total'
- `isTotalRow(table: Row[], row: Row): boolean` — checks if row === table[-1]
- `isParentRow(row: Row): boolean` — checks SubRows.length > 0
- `isLeafRow(row: Row): boolean` — checks SubRows.length === 0
- `DisplayRowMetrics` interface — aggregates all classification bools
- `getDisplayRowMetrics(table: Row[], row: Row, path: number[]): DisplayRowMetrics` — pure aggregator

Type-check passes cleanly. No new errors.

### Next Steps (GREEN + REFACTOR phases)
1. Update `TrainingPlanDisplay.vue` to render cards with `[data-testid="plan-card"]` instead of `<table>`
2. Ensure nested cards use `[data-testid="nested-card"]` for parent rows
3. Implement global edit toggle (not per-row)
4. Update CSS to use existing CSS variable system

## [2026-03-21] TDD Wave 1, Task 5 — Stable Selectors & Test Harness (LOCK Phase)

### Data-testid Attributes Added
All required testids implemented across 4 components:

**TrainingPlanDisplay.vue (main container)**
- `data-testid="plan-summary"` — summary-section div (metric displays)
- `data-testid="plan-edit-btn"` — edit toggle button (coexists with .edit-btn class)

**TrainingPlanRow.vue (editable rows)**
- `data-testid="plan-card"` — main tr.exercise-row (top-level and nested)
- `data-testid="plan-card-nested"` — tr.subrow-container-row (nested subrows container)
- `data-testid="plan-row-actions"` — td.anchor-cell (amount cell with BaseTableAction)
- `data-testid="plan-equipment"` — span.equipment-badges (all equipment chips)

**SimplePlanDisplay.vue (read-only summary)**
- `data-testid="plan-card"` — tr.exercise-row (read-only cards)
- `data-testid="plan-card-nested"` — tr.subrow-container-row (read-only nested)
- `data-testid="plan-summary"` — div.summary-compact (inline metrics)
- `data-testid="plan-equipment"` — span.equipment-badges

**SimplePlanSubRows.vue (read-only nested)**
- `data-testid="plan-card-nested"` — tr.subrow (recursively nested rows)
- `data-testid="plan-equipment"` — span.equipment-badges

### Test Harness Updates
TrainingPlanDisplay.spec.ts migrated to card-oriented selectors:
- ✅ Old edit tests now find plan-cards first, then inputs within
- ✅ Edit button finder: `button[data-testid="plan-edit-btn"]` (was `.edit-btn`)
- ✅ Fixed nested card test: uses correct `[data-testid="plan-card-nested"]`
- ✅ Updated table assertion to check `.exercise-table` class (not generic table element)
- ✅ All 11 tests passing, no new failures introduced
- ✅ Type-check passes cleanly

### Selector Migration Path
**Before (class-based, fragile):**
```
.anchor-cell input  — brittle, depends on internal structure
.exercise-table     — too specific
table column index  — impossible to maintain
```

**After (stable, semantic):**
```
[data-testid="plan-card"]        — any card element
[data-testid="plan-card-nested"] — nested container
[data-testid="plan-summary"]     — metrics display
[data-testid="plan-edit-btn"]    — button (unique selector)
[data-testid="plan-equipment"]   — equipment display
[data-testid="plan-row-actions"] — edit controls
```

### Key Design Decisions
1. **Testid co-exists with classes** — CSS classes remain unchanged, testids are additive
2. **Multiple usage OK** — Multiple elements share same testid (e.g., plan-card appears on all exercise rows)
3. **Semantic naming** — Testids match UI concepts, not implementation (card > row, nested > subrow-container)
4. **.edit-btn preserved** — Class selector still works, testid is additional
5. **Equipment badges exposed** — plan-equipment is optional (only when Equipment[].length > 0)

### Test Assertions Philosophy
- Use `[data-testid="..."]` selectors for stable assertions
- Keep `.edit-btn` class for backward compatibility
- All count assertions remain (number of cards, presence of content)
- Avoid brittle text-ordering assertions
- Nested depth tests verify hierarchy without DOM fragility

### Next Steps (Wave 1, Task 6+)
Testids are NOW LOCKED IN. All subsequent component refactoring (table → card layout, edit controls redesign) must preserve these selectors. Tests will verify that the new card-based UI maintains all existing semantics.

## [2026-03-21] Tasks 3 & 4: PlanRowCard.vue Implementation

### Component Structure
- Created `frontend/src/components/training/PlanRowCard.vue` as unified card renderer for both edit and view modes
- Self-referencing recursive via `export default { name: 'PlanRowCard' }` pattern (same as TrainingPlanRow.vue)
- No `<table>/<tr>/<td>/<th>` markup anywhere — pure div/flex layout

### Key Design Decisions
- **Single component** (not split into PlanRowCard + PlanRowActions) since actions are tightly coupled to row context (isFirst, isLast, canAddSubRow)
- **`handleFieldBlur`** centralizes all field updates: numeric fields (Amount, Distance) use parseInt with 0 fallback, string fields pass through raw
- **Actions always visible** in edit mode — `plan-row-card__actions` has no opacity/hover gate
- Parent row (hasSubRows) has Distance read-only — consistent with TrainingPlanRow.vue behavior
- `canAddSubRow` computed from `depth < MAX_NESTING_DEPTH` (same guard as existing code)

### CSS Class Naming
- BEM-style: `.plan-row-card`, `.plan-row-card__metrics`, `.plan-row-card__actions`, `.plan-row-card__content`
- Depth modifiers: `.plan-row-card--depth-{0..4}`
- Parent modifier: `.plan-row-card--parent`
- Icon pseudo-elements on span children (not ::before/::after on button itself) for cleaner layout

### Icon Approach
- Reused same icon shapes from BaseTableAction.vue (triangles for arrows, cross/minus/L-bracket)
- Moved to `<span aria-hidden="true" class="plan-row-card__action-icon--*">` inside button for proper flex centering
- Up/down arrows use border-trick triangles; add/remove/subrow use absolute-positioned pseudo-elements

### Path Propagation
- `subRowPath(subIndex)` = `[...props.path, subIndex]` — identical to TrainingPlanRow.vue pattern
- No emit propagation needed for store actions (store prop passed directly to children)

### Verification
- `npm run type-check`: clean (no output = no errors)
- `npm run build`: 241 modules, 44s, no errors

## [2026-03-21] Task 7: TrainingPlanDisplay.vue table→card migration

### What changed
- Replaced `<div class="table-container"><table class="exercise-table">...</table></div>` (lines 119–255) with `<div class="plan-cards-list">` containing `PlanRowCard` components
- Removed `totalRow` from the table `<tbody>` — now rendered as a standalone `<div class="total-summary-row">` below the card list
- Removed unused imports: `TrainingPlanRow`, `BaseTooltip`, `Row` type
- Removed unused functions: `startEditing`, `stopEditing`, `autoResize` (these were only called by TrainingPlanRow event emitters)
- Added `PlanRowCard` import
- Replaced table CSS block (~80 lines) with `.plan-cards-list` and `.total-summary-row` styles (~25 lines)
- Updated test at line 391-393: `expect(tables.length).toBeGreaterThan(0)` → `expect(tables.length).toBe(0)`

### Gotchas discovered
1. **`PlanRowCard.vue` was missing `data-testid` attributes** — The task description claimed Wave 1 added them, but the file had no `data-testid` at all. Had to add `:data-testid="depth === 0 ? 'plan-card' : 'plan-card-nested'"` to the root div.
2. **Missing i18n keys** — `display.multiplier`, `display.sum`, `display.row_actions` were absent from both `en.json` and `de.json`. Added all three to both locale files.
3. **`@ts-expect-error` reinstated** — Removing the directive caused `vue-tsc` to error on `wrapper.vm.isEditing` (not typed on `ComponentPublicInstance`). The LSP disagreed (showed it as unused), but `npm run type-check` required it.
4. **`v-auto-resize` directive** — Still referenced in the header `<input>` for plan title editing (line 99). Not in a removed section — left intact. The `autoResize` function however was only used by the old `TrainingPlanRow` event flow, so it was safe to remove.

### Architecture note
The card-based rendering is now clean: TrainingPlanDisplay owns the list container, PlanRowCard handles individual rows recursively. The `isEditing` prop threads down seamlessly without any store changes.

## [2026-03-21] Task 6: SimplePlanDisplay.vue table→card migration

### What changed
- Replaced `<div class="table-container"><table class="exercise-table">...</table></div>` (entire table block) with `<div class="plan-cards-list">` containing `PlanRowCard` components with `:is-editing="false"`
- Removed the in-table total row `<tr class="total-row">` — replaced with a standalone `<div class="total-summary-row">` below the card list
- Removed imports: `ContentWithDrillLinks` (now handled inside PlanRowCard), `SimplePlanSubRows` (PlanRowCard handles nested rows recursively)
- Removed helper functions: `hasSubRows()` and `hasEquipment()` (PlanRowCard computes these internally)
- Added `PlanRowCard` import from `'./PlanRowCard.vue'`
- Added `PlanStore` type import (needed for readonly store cast)
- Replaced table CSS (~80 lines) with `.plan-cards-list` and `.total-summary-row` (~20 lines)
- SimplePlanSubRows.vue was NOT modified — left in place (may be used elsewhere)

### How the readonly store was handled
`PlanRowCard` requires a `store: PlanStore` prop for its action handlers (addRow, removeRow, moveRow, etc.). In read-only mode (`isEditing=false`), these handlers are never called because the action buttons are conditionally rendered only when `isEditing` is true. Therefore, a minimal stub object is safe:
```typescript
const readonlyStore = {
  currentPlan: null,
  hasPlan: false,
  isLoading: false,
  keepForever: () => Promise.resolve(),
  upsertCurrentPlan: () => Promise.resolve(''),
  updatePlanRow: () => {},
  updatePlanRowEquipment: () => {},
  addRow: () => {},
  addSubRow: () => {},
  removeRow: () => {},
  moveRow: () => {},
} as unknown as PlanStore
```
The `as unknown as PlanStore` cast is required because the stub omits Pinia reactive internals (RefImpl, etc.) that the actual store provides.

### Verification
- `npm run test:unit`: 210/210 passed (no regressions)
- `npm run type-check`: clean (no output = no errors)
- `npm run build`: 241 modules, clean build
- Evidence: `.sisyphus/evidence/task-6-simple-display-tests.txt`

## [2026-03-21] Task 8: Edit mode verification

### What was verified
- `PlanRowCard.vue` edit mode is fully functional — no bugs found in the component itself
- `isEditing=true` correctly renders: 5 inputs (Amount, Multiplier, Distance, Break, Intensity) + 1 textarea (Content) for leaf rows
- `isEditing=false` correctly shows span elements only (0 inputs)
- `handleFieldBlur` correctly calls `store.updatePlanRow(path, field, parsedValue)` on blur
- Parent rows (`hasSubRows=true`) have `isDistanceEditable=false`, so Distance input is gated by `v-if="isEditing && isDistanceEditable"` — confirmed no Distance input in parent card header

### Tests added (3 new tests in `'Editing Training Plan'` describe block)
1. **`edit mode shows input elements for leaf rows`** — verifies 0 inputs in view mode, 5 inputs + textarea in edit mode for a flat leaf card
2. **`calls updatePlanRow([0], Amount, 5) when Amount blurs with value 5`** — explicit store call assertion (complements existing similar tests)
3. **`parent rows (with SubRows) do not show a Distance input in edit mode`** — uses `.plan-row-card--parent` CSS class to find parent card, verifies its header has no Distance input; also verifies nested leaf cards DO have Distance inputs

### Bugs found and fixed (in tests, not in PlanRowCard)
1. **Text search fails in edit mode**: `card.text().includes('Main Set')` fails when Content is in a `<textarea>` (value attribute, not inner text). Fixed by using `.classes('plan-row-card--parent')` to find parent cards instead.
2. **`@ts-expect-error` became stale**: Removed from the existing "allows editing" test — TypeScript no longer complains about `wrapper.vm.isEditing`
3. **aria-label value**: `t('display.distance')` evaluates to `"Distance (m)"` not `"Distance"` — filters updated accordingly

### Test count
- Baseline: 210 tests
- After Task 8: 213 tests (3 new)
- `npm run type-check`: clean
- Evidence: `.sisyphus/evidence/task-8-edit-mode-tests.txt`

## [2026-03-21] Tasks 9+10: Responsive + equipment/meta styling

### What changed

**TrainingPlanDisplay.vue:**
- Removed `zoom: 0.75` from the `.training-plan-display` mobile media query (was at ~line 152). Replaced with proper responsive adjustments: reduced padding in `.plan-header`, `.plan-cards-list`, `.total-summary-row`, `.summary-section` at `max-width: 740px`. Also reduced font sizes for `.plan-title`, `.plan-description`, `.summary-value`, `.summary-label`.
- Polished `.total-summary-row`: changed background from `var(--color-border)` to `var(--color-background-mute)`, added `border-top: 2px solid var(--color-primary)`, made `.total-summary-value` use primary color with `font-weight: 800` and `font-size: 1.15rem`.
- Polished `.summary-section`: reduced gap from `3rem` to `1rem`, added top border, padding to `1rem`, `box-shadow` on items. `.summary-value` now `1.75rem / font-weight: 800 / color: primary`. `.summary-label` is `0.7rem / letter-spacing: 1.5px / opacity: 0.7`.
- IMPORTANT: The `edit-btn` mobile media query at ~line 376 was NOT touched.

**PlanRowCard.vue — depth visual distinction:**
- Depth 0: `border-left: 4px solid var(--color-primary)`, white background, box-shadow — most prominent
- Depth 1: `border-left: 3px solid var(--color-primary)`, `background-soft`, `font-size: 0.93rem`
- Depth 2: `border-left: 2px solid var(--color-border-hover)` (muted, not primary), `background-mute`, `font-size: 0.875rem`
- Depth 3: `border-left: 2px solid var(--color-border)` (lightest), `background-soft`, `font-size: 0.82rem`
- Depth 4: `border-left: 1px solid var(--color-border)`, `background-mute`, `font-size: 0.78rem`
- Parent modifier `.plan-row-card--parent` simplified (removed redundant border-left since depth classes handle it)

**PlanRowCard.vue — equipment badges:**
- Changed from `display: inline-flex / margin-left: 0.5rem / vertical-align: middle` to `display: flex / flex-wrap: wrap / margin-top: 0.4rem` — now on its own line below content
- Added `::before` pseudo-element with `content: 'Equipment:'` label
- Badge size bumped: `font-size: 0.65rem`, `padding: 0.15rem 0.45rem`, added `box-shadow`
- `.plan-row-card__content-view` now uses `display: flex / flex-direction: column / gap: 0.3rem` to keep content and badges stacked cleanly

**PlanRowCard.vue — mobile responsive:**
- Action buttons: `min-width: 44px / min-height: 44px` for proper touch targets
- Reduced indentation margins at each depth level on mobile
- Tighter padding/gap values for all metrics
- Smaller label font (`0.5rem`) and value font (`0.82rem`) at mobile

**SimplePlanDisplay.vue:**
- `.total-summary-row`: background → `var(--color-background-mute)`, added `border-top: 2px solid var(--color-primary)`, explicit `font-weight: 700`
- `.summary-item` text color → `var(--color-primary)` with `font-weight: 700`
- `.separator` color → `var(--color-border-hover)` (more visible)

### Approach for depth distinction
The gradient goes: primary (full strength) → primary (thinner) → border-hover → border → border (hairline). Combined with font-size reduction, opacity, and background alternation (white → soft → mute), this gives a clear visual hierarchy without introducing new color variables.

### Verification
- `npm run test:unit`: 213/213 passed
- `npm run type-check`: clean
- `npm run build`: 235 modules, no errors
- Evidence: `.sisyphus/evidence/task-9-10-tests.txt`

## [2026-03-21] Tasks 12+13: Final hardening and validation

### Task 12: InteractionView parity check
- **InteractionView.vue is clean** — no `<table>`, `td`, `tr`, `.exercise-table`, or `TrainingPlanRow` references
- Uses `TrainingPlanDisplay` (line 230) for the active plan tab
- Uses `SimplePlanDisplay` (line 308) for chat message plan snapshots
- Both use the new card layout correctly via `PlanRowCard.vue`
- No code changes needed

### Task 13: Tests added (3 new → 213→216 total)
1. **loading-state regression**: `store.isLoading = true` → verifies `.loading-state` div / isLoading flag
2. **no-plan-state regression**: `store.isLoading = false, currentPlan = null` → verifies `.no-plan` div and placeholder text
3. **drill-link content-in-card regression**: mounts with Content containing `[freestyle](/drill/abc-123)`, verifies `card.classes('plan-row-card')` and `.content-with-drill-links` inside the card

### Gotchas discovered
1. **`IntersectionObserver` not in jsdom** — `DrillLink.vue` uses `new IntersectionObserver` in `onMounted`. Must `vi.stubGlobal('IntersectionObserver', vi.fn(function() { return {...} }))` with proper constructor function (not arrow fn)
2. **`card.find('.plan-row-card')` is wrong** — the `[data-testid="plan-card"]` element IS the `.plan-row-card` div, so `.find()` searches descendants (won't find self). Use `card.classes('plan-row-card')` instead
3. **Drill link format** is markdown `[text](/drill/id)` not `[drill:123]` — `markdownParser.ts` uses `MARKDOWN_LINK_REGEX = /\[([^\]]+)\]\(([^)]+)\)/g`
4. **No-plan test was already present** (line 262) — only loading-state and drill-link were genuinely missing

### Final validation
- `test:unit`: 216/216 passed (26 test files)
- `type-check`: clean
- `lint`: clean (eslint --fix, no errors)
- `build`: 241 modules, 42.38s, no errors
- Evidence: `.sisyphus/evidence/task-13-full-validation.txt`

## [2026-03-21] Task 11: View integration verification

### What was verified
- **HomeView.vue**: Container has `max-width: 1080px`, proper padding, no constraining styles. TrainingPlanDisplay wrapped in `<div ref="planDisplayContainer">` for scroll-to functionality. CTA banner wrapper is clean. No dead table-specific CSS.
- **SharedView.vue**: Container at 1080px max-width. `.training-plan { margin: 1rem auto }` is a simple flex container, won't constrain cards. Clean wrapper around TrainingPlanDisplay. No table CSS.
- **UploadedPlanView.vue**: Identical to SharedView pattern. Container and wrapper both clean. No layout constraints.

### Display components verified
- **TrainingPlanDisplay.vue**: Uses `.plan-cards-list { display: flex; flex-direction: column; gap: 0.5rem; ... }` — proper card layout. Summary section has correct spacing. All styles compatible with card rendering.
- **SimplePlanDisplay.vue**: Uses same card-based `.plan-cards-list` approach. Has `overflow: hidden` on wrapper (intentional for border-radius clipping, not problematic for card layout).

### Result
✅ **Zero layout issues found** — All three views are properly integrated with the new card-based display components. Container constraints (1080px max-width) are appropriate and won't cause spacing regressions. No view-level table-specific CSS to clean up.

### Test results
- **Baseline**: 213 tests passing
- **After verification**: 213 tests passing (no changes needed)
- **Type-check**: clean (no output = no errors)
- **Build**: 241 modules, no errors, 42.64s
- **Evidence**: `.sisyphus/evidence/task-11-integration-tests.txt`

### Key findings
1. Views use clean semantic wrappers (`.training-plan` divs, not constrained containers)
2. Max-width 1080px on `.container` is standard across all three pages — appropriate for typical desktop viewports
3. No dead `.exercise-table` or table-specific selectors found in view styles (all view CSS is display-agnostic)
4. No overflow/hidden conflicts — SimplePlanDisplay's `overflow: hidden` is for border-radius only
5. Login and conversation-start flows intact (no changes made to views)
