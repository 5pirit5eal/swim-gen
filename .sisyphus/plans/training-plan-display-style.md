# Training Plan Display Style Migration Plan

## TL;DR

> **Quick Summary**: Replace the current table-based training-plan rendering with a card/list layout inspired by `frontend/web.html` and `frontend/mobile.html`, while preserving the existing Swim Gen visual identity, Vue architecture, and row-editing data model.
>
> **Deliverables**:
> - Shared card-based training-plan presentation for editable and read-only flows
> - Separate edit mode compatible with the new card layout
> - TDD coverage for display, nesting, and edit-mode regressions
> - Responsive equipment/summary/meta sections aligned with the reference hierarchy
>
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 3 implementation waves + final verification
> **Critical Path**: Task 1 → Task 3 → Task 7 → Task 12 → Task 13

---

## Context

### Original Request
Change the training plan display to use the new display style shown in `frontend/web.html`, `frontend/mobile.html`, and `frontend/image.png`, while maintaining the app’s original visual style. The current tabular representation may be replaced.

### Interview Summary
**Key Discussions**:
- Apply the redesign to both editable and read-only/shared training-plan views.
- Follow `web.html` more closely than `mobile.html` for hierarchy and nesting.
- Nested sets can reach depth 4, so child items must be visually simplified and clearly contained within their parent card.
- Editable screens should keep a separate edit mode rather than always-on inline table editing.
- Use TDD for the refactor.

**Research Findings**:
- Editable plan display currently lives in `frontend/src/components/training/TrainingPlanDisplay.vue:86-285` and recursive row rendering in `frontend/src/components/training/TrainingPlanRow.vue:58-191`.
- Read-only plan display currently lives in `frontend/src/components/training/SimplePlanDisplay.vue:51-124` and `frontend/src/components/training/SimplePlanSubRows.vue:19-52`.
- The existing display is deeply coupled to `<table>` markup, `BaseTableAction`, and recursive row paths.
- The reference HTML uses accent-bar section cards, badge labels, a title/description hierarchy, compact right-aligned metric clusters, nested contained child items, and meta sections for equipment/notes.
- Frontend test infrastructure already exists through Vitest, Vue Test Utils, and CI (`frontend/vitest.config.ts:1-17`, `.github/workflows/frontend-validate.yaml:66-92`).

### Metis Review
**Identified Gaps** (addressed):
- Do not introduce Tailwind into the Vue app; translate the reference layout into the existing CSS-variable system.
- Do not invent new workout-phase semantics from missing data; default accent treatment to existing app-brand styling unless a real data signal already exists.
- Treat `BaseTableAction` as a full interaction redesign, not a styling tweak.
- Keep the current global edit toggle as the default edit-mode entry point unless implementation evidence forces a narrower toggle model.

---

## Work Objectives

### Core Objective
Deliver a reusable card/list presentation system for training plans that preserves existing Swim Gen branding, supports nested row hierarchies up to depth 4, and works consistently across editable and read-only screens without relying on table markup.

### Concrete Deliverables
- Card-based plan header and summary layout for training plans
- Top-level set cards and nested child-item rendering replacing table rows/cells
- Separate edit mode that still supports row add/remove/move/add-subrow workflows
- Read-only/shared plan parity with the same visual hierarchy
- Regression tests covering display, nesting, and edit-mode transitions

### Definition of Done
- [ ] `frontend/src/components/training/TrainingPlanDisplay.vue` no longer depends on `.exercise-table` / `<table>` rendering for the primary plan view.
- [ ] `frontend/src/components/training/SimplePlanDisplay.vue` renders the same card hierarchy as the editable display’s read-only mode.
- [ ] Nested rows remain clearly contained inside parent cards through depth 4.
- [ ] `npm run test:unit`, `npm run type-check`, `npm run lint`, and `npm run build` all pass in `frontend/`.

### Must Have
- Preserve the existing data model (`Row`, `PlanStore`, `rowHelpers`) and drill-link rendering.
- Preserve editable workflows through a separate edit mode.
- Preserve responsive behavior for both narrow and wide screens.
- Add stable selectors/test IDs required for automated QA of the new layout.

### Must NOT Have (Guardrails)
- Do **not** add Tailwind or copy the full standalone shell/navigation from `web.html` / `mobile.html`.
- Do **not** change backend contracts or the `Row` schema just to classify card colors.
- Do **not** leave table-only interaction patterns (especially `BaseTableAction`) awkwardly embedded in the new card layout.
- Do **not** break `ContentWithDrillLinks`, equipment badges, or row path-based store operations.
- Do **not** hide critical training data behind hover-only interactions.

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** — all verification must be agent-executed.

### Test Decision
- **Infrastructure exists**: YES
- **Automated tests**: TDD
- **Framework**: Vitest + Vue Test Utils + jsdom
- **If TDD**: Each implementation task follows RED → GREEN → REFACTOR within its scope.

### QA Policy
Every task must include agent-executed QA scenarios with saved evidence.

- **Frontend/UI**: Playwright skill for layout, responsive, and interaction verification
- **Component behavior**: `npm run test:unit -- <spec>` for focused verification
- **Integration smoke**: `npm run type-check`, `npm run lint`, `npm run build`
- **Evidence root**: `.sisyphus/evidence/`

---

## Execution Strategy

### Parallel Execution Waves

```text
Wave 1 (Start immediately — foundations + shared presentation primitives)
├── Task 1: Display-model helpers + fixture expansion [quick]
├── Task 2: Shared card shell, header, and summary primitives [visual-engineering]
├── Task 3: Shared nested item renderer contract [quick]
├── Task 4: Edit-mode action pattern redesign [unspecified-high]
└── Task 5: Selector/test harness foundation [quick]

Wave 2 (After Wave 1 — main feature migrations)
├── Task 6: Migrate SimplePlanDisplay to card layout [visual-engineering]
├── Task 7: Migrate TrainingPlanDisplay read/view mode to card layout [visual-engineering]
├── Task 8: Migrate TrainingPlanRow edit mode off table interactions [unspecified-high]
├── Task 9: Equipment/meta/summary section redesign [visual-engineering]
└── Task 10: Responsive layout + deep nesting styling pass [visual-engineering]

Wave 3 (After Wave 2 — screen integration + resilience)
├── Task 11: Home/shared/uploaded screen integration pass [quick]
├── Task 12: Interaction/snapshot parity pass [unspecified-high]
└── Task 13: Edge-case, accessibility, and regression hardening [deep]

Wave FINAL (After ALL tasks — independent review, 4 parallel)
├── Task F1: Plan compliance audit (oracle)
├── Task F2: Code quality review (unspecified-high)
├── Task F3: Real manual QA / Playwright execution (unspecified-high)
└── Task F4: Scope fidelity check (deep)
```

### Dependency Matrix

- **1**: — → 6, 7, 8, 13
- **2**: — → 6, 7, 9, 10
- **3**: — → 6, 7, 8, 10, 13
- **4**: — → 8, 12, 13
- **5**: — → 6, 7, 8, 11, 12, 13
- **6**: 1, 2, 3, 5 → 11, 12, 13
- **7**: 1, 2, 3, 5 → 11, 12, 13
- **8**: 1, 3, 4, 5 → 12, 13
- **9**: 2 → 11, 12, 13
- **10**: 2, 3 → 11, 12, 13
- **11**: 6, 7, 9, 10 → F1-F4
- **12**: 6, 7, 8, 9, 10 → F1-F4
- **13**: 1, 3, 4, 5, 6, 7, 8, 9, 10 → F1-F4

### Agent Dispatch Summary

- **Wave 1**: **5 agents** — T1 `quick`, T2 `visual-engineering`, T3 `quick`, T4 `unspecified-high`, T5 `quick`
- **Wave 2**: **5 agents** — T6 `visual-engineering`, T7 `visual-engineering`, T8 `unspecified-high`, T9 `visual-engineering`, T10 `visual-engineering`
- **Wave 3**: **3 agents** — T11 `quick`, T12 `unspecified-high`, T13 `deep`
- **FINAL**: **4 agents** — F1 `oracle`, F2 `unspecified-high`, F3 `unspecified-high` + `playwright`, F4 `deep`

---

## TODOs

- [ ] 1. Establish display-model helpers and richer nested fixtures

  **What to do**:
  - Add/update tests and fixture builders to cover top-level rows, nested rows through depth 4, total rows, and mixed equipment/content cases.
  - Introduce any pure helper(s) needed to derive UI-friendly card metrics from existing `Row` data without changing the backend contract.
  - Ensure helpers distinguish top-level exercise rows from the terminal total row and keep parent/child relationships path-based.

  **Must NOT do**:
  - Do not add new persisted fields to `Row`.
  - Do not infer semantic phase labels that the data does not actually provide.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: focused helper/test setup in a few files.
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: layout styling is not the primary concern here.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3, 4, 5)
  - **Blocks**: 6, 7, 8, 13
  - **Blocked By**: None

  **References**:
  - `frontend/src/types/training.ts:16-28` - Defines the immutable `Row` shape that helpers must preserve.
  - `frontend/src/utils/rowHelpers.ts:17-128` - Existing source-of-truth for row normalization, nesting depth, and sum recalculation.
  - `frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts:9-129` - Current component test entry point to expand with failing card-layout expectations.
  - `frontend/src/stores/__tests__/trainingPlan.spec.ts:73-106` - Existing plan fixture pattern to reuse when creating richer row trees.

  **Acceptance Criteria**:
  - [ ] Failing tests exist for nested card-oriented display expectations before UI refactor starts.
  - [ ] Helper tests prove parent rows, leaf rows, and total rows are classified correctly.
  - [ ] No type changes are required in `frontend/src/types/training.ts`.

  **QA Scenarios**:
  ```text
  Scenario: Helper classifies nested rows and totals correctly
    Tool: Bash
    Preconditions: frontend dependencies installed
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Confirm the added failing/passing expectations mention nested rows, total row handling, and card metrics
      3. Save the Vitest output
    Expected Result: Focused spec executes and reports assertions for nested display-model behavior
    Failure Indicators: Missing nested assertions, type errors, or row classification failures
    Evidence: .sisyphus/evidence/task-1-display-model-tests.txt

  Scenario: Helper rejects schema drift
    Tool: Bash
    Preconditions: helper changes applied
    Steps:
      1. Run `npm run type-check`
      2. Verify no new required fields were introduced on `Row`
    Expected Result: Type-check passes with existing `Row` contract intact
    Failure Indicators: Type errors requiring schema changes
    Evidence: .sisyphus/evidence/task-1-type-check.txt
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-1-display-model-tests.txt`
  - [ ] `.sisyphus/evidence/task-1-type-check.txt`

  **Commit**: NO

- [ ] 2. Build shared card-shell styling primitives from existing theme tokens

  **What to do**:
  - Create the shared structural/styling foundation for card containers, accent rails, badges, metric clusters, nested containers, and meta sections using existing CSS variables.
  - Translate the `web.html` / `mobile.html` spacing and hierarchy into Vue component styles without importing Tailwind utilities.
  - Define how top-level cards, child cards, and summary/meta blocks visually differ while still feeling native to Swim Gen.

  **Must NOT do**:
  - Do not copy the standalone app chrome from the HTML references.
  - Do not hardcode colors that bypass the current theme variable system unless there is no matching token and the value is isolated/documented.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: design translation and responsive visual hierarchy are central.
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: useful for preserving identity while translating a static reference to the app’s design system.
  - **Skills Evaluated but Omitted**:
    - `playwright`: verification, not implementation.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3, 4, 5)
  - **Blocks**: 6, 7, 9, 10
  - **Blocked By**: None

  **References**:
  - `frontend/web.html:192-407` - Canonical web hierarchy for header, set cards, nested rows, and footer meta sections.
  - `frontend/mobile.html:110-325` - Mobile compression pattern for cards, metrics, nested children, and meta chips.
  - `frontend/src/components/training/TrainingPlanDisplay.vue:287-616` - Existing theme usage and summary/button styling to preserve visually.
  - `frontend/src/components/training/SimplePlanDisplay.vue:126-287` - Current read-only styling baseline.
  - `frontend/src/assets/base.css` - Existing app color variables and responsive conventions.

  **Acceptance Criteria**:
  - [ ] Shared class names or shared component styles exist for card shells, nested child blocks, and metric clusters.
  - [ ] Styling uses app variables/tokens rather than Tailwind.
  - [ ] Desktop and mobile layout rules are explicitly defined.

  **QA Scenarios**:
  ```text
  Scenario: Card styling renders with preserved visual identity on desktop
    Tool: Playwright
    Preconditions: frontend dev server running with a populated plan view
    Steps:
      1. Open the training plan page at 1440x1200
      2. Capture the plan header, first top-level card, and nested-card region
      3. Verify accent rail, badge, title/description hierarchy, and right-aligned metrics are visible
    Expected Result: New card shell matches reference hierarchy without replacing the app shell/navigation
    Failure Indicators: Table still visible, Tailwind-like shell copied, or metrics collapse incorrectly
    Evidence: .sisyphus/evidence/task-2-desktop-card-shell.png

  Scenario: Card styling collapses cleanly on mobile width
    Tool: Playwright
    Preconditions: same page with mobile viewport
    Steps:
      1. Open the same page at 390x844
      2. Capture top-level and nested cards
      3. Verify content stacks vertically and remains readable without horizontal scrolling
    Expected Result: Mobile card presentation is compact and readable
    Failure Indicators: overflow, clipped metrics, or unreadable nested content
    Evidence: .sisyphus/evidence/task-2-mobile-card-shell.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-2-desktop-card-shell.png`
  - [ ] `.sisyphus/evidence/task-2-mobile-card-shell.png`

  **Commit**: NO

- [ ] 3. Define and implement the shared nested item renderer contract

  **What to do**:
  - Refactor recursive rendering so both editable and read-only displays can share the same nested card/list structure rather than separate table-based recursion.
  - Ensure child items are visually contained within parents and simplify their presentation compared with top-level cards.
  - Preserve path-based recursion so downstream edit actions still target the correct row.

  **Must NOT do**:
  - Do not duplicate separate recursive rendering logic for editable vs read-only if a shared structure can cover both.
  - Do not flatten the hierarchy in a way that obscures parent-child relationships.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: structural component contract refactor with bounded file scope.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: helpful but secondary to recursion/component architecture.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 4, 5)
  - **Blocks**: 6, 7, 8, 10, 13
  - **Blocked By**: None

  **References**:
  - `frontend/src/components/training/TrainingPlanRow.vue:58-191` - Current recursive editable row rendering and path propagation.
  - `frontend/src/components/training/SimplePlanSubRows.vue:19-52` - Current read-only recursive structure to unify.
  - `frontend/src/utils/rowHelpers.ts:72-98` - Path semantics that must remain intact.
  - `frontend/web.html:241-309` - Preferred nested-child composition to emulate more closely.

  **Acceptance Criteria**:
  - [ ] A shared recursive rendering approach exists for nested items.
  - [ ] Depth 4 data renders without table markup.
  - [ ] Parent-child containment remains visually obvious in both modes.

  **QA Scenarios**:
  ```text
  Scenario: Shared renderer shows four nesting levels without table markup
    Tool: Bash
    Preconditions: renderer tests added
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Confirm rendered HTML snapshots/assertions do not include `<table>` for the primary plan content
      3. Save test output
    Expected Result: Spec proves recursive rendering works through depth 4 in card/list structure
    Failure Indicators: `<table>` still present or deeper levels disappear
    Evidence: .sisyphus/evidence/task-3-recursion-tests.txt

  Scenario: Nested containment is visible in browser
    Tool: Playwright
    Preconditions: page with a nested sample plan
    Steps:
      1. Open the plan page on desktop
      2. Capture a parent card containing at least two nested child items
      3. Verify child blocks are rendered inside the parent boundary, not as sibling top-level cards
    Expected Result: Parent boundary encloses child items clearly
    Failure Indicators: Children appear detached or hierarchy is visually ambiguous
    Evidence: .sisyphus/evidence/task-3-nested-containment.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-3-recursion-tests.txt`
  - [ ] `.sisyphus/evidence/task-3-nested-containment.png`

  **Commit**: NO

- [ ] 4. Redesign edit-mode action controls for card layout

  **What to do**:
  - Replace the table-anchored `BaseTableAction` interaction pattern with a card-compatible action affordance for add/remove/move/add-subrow controls.
  - Keep the existing global edit toggle in `TrainingPlanDisplay.vue`, but make row controls discoverable and usable within cards.
  - Preserve all existing store method wiring (`addRow`, `removeRow`, `moveRow`, `addSubRow`, `updatePlanRow`, `updatePlanRowEquipment`).

  **Must NOT do**:
  - Do not rely on hover-only sidecar controls positioned outside the card.
  - Do not remove row-management capabilities in edit mode.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: this is the riskiest behavioral redesign in the refactor.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: useful visually, but the primary challenge is preserving editing behavior safely.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 3, 5)
  - **Blocks**: 8, 12, 13
  - **Blocked By**: None

  **References**:
  - `frontend/src/components/ui/BaseTableAction.vue:21-230` - Current control set that must be re-homed for cards.
  - `frontend/src/components/training/TrainingPlanRow.vue:67-187` - Existing edit wiring and nested add-subrow behavior.
  - `frontend/src/stores/trainingPlan.ts:575-620` - Store API surface that edit-mode controls must keep using.

  **Acceptance Criteria**:
  - [ ] Edit mode exposes add/remove/move/add-subrow actions in the new layout.
  - [ ] Controls are accessible on desktop and mobile without hover dependency.
  - [ ] Store action wiring remains unchanged from the caller’s perspective.

  **QA Scenarios**:
  ```text
  Scenario: Edit mode exposes row controls inside card layout
    Tool: Playwright
    Preconditions: editable plan page loaded with at least two top-level rows
    Steps:
      1. Click `.edit-btn`
      2. Locate controls inside the first rendered card
      3. Trigger add, move-down, and add-subrow once each
      4. Capture the updated UI state
    Expected Result: Controls are visible and act within the card layout without off-card hover panels
    Failure Indicators: No visible controls, controls only appear on hover, or actions fail silently
    Evidence: .sisyphus/evidence/task-4-edit-controls.png

  Scenario: Edit mode remains usable on mobile width
    Tool: Playwright
    Preconditions: same editable view at mobile viewport
    Steps:
      1. Open page at 390x844
      2. Enter edit mode and attempt to add a subrow from the first parent card
      3. Verify no controls are clipped offscreen
    Expected Result: Card controls remain tappable and visible on mobile
    Failure Indicators: overflow, clipped buttons, or inaccessible controls
    Evidence: .sisyphus/evidence/task-4-mobile-edit-controls.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-4-edit-controls.png`
  - [ ] `.sisyphus/evidence/task-4-mobile-edit-controls.png`

  **Commit**: NO

- [ ] 5. Install stable selectors and test harness for new card UI

  **What to do**:
  - Add stable selectors / `data-testid` / semantic hooks needed for automated tests and Playwright scenarios on top-level cards, nested items, edit-mode controls, summary blocks, and equipment chips.
  - Update test helpers so component tests no longer depend on table cell order or table-specific selectors.
  - Keep selector naming aligned across editable and read-only plan displays where feasible.

  **Must NOT do**:
  - Do not make tests depend on fragile text ordering alone.
  - Do not add selectors that expose implementation-only details with no test value.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: focused testability scaffolding.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `playwright`: consumers of selectors, not necessary to add them.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 3, 4)
  - **Blocks**: 6, 7, 8, 11, 12, 13
  - **Blocked By**: None

  **References**:
  - `frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts:37-129` - Existing tests currently tied to `.anchor-cell input` and `.edit-btn`.
  - `frontend/src/components/training/TrainingPlanDisplay.vue:275-283` - Existing stable `.edit-btn` entry point to keep or deliberately replace.

  **Acceptance Criteria**:
  - [ ] Tests can target cards/items/controls without referencing table cells.
  - [ ] Shared and editable displays expose at least the critical selectors needed by the QA scenarios in this plan.

  **QA Scenarios**:
  ```text
  Scenario: Component tests use stable card selectors instead of table selectors
    Tool: Bash
    Preconditions: tests updated
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Confirm no assertions depend on `.exercise-table`, `td`, or table column indexes
    Expected Result: Focused tests pass using card-oriented selectors
    Failure Indicators: Tests still break on table-only selectors
    Evidence: .sisyphus/evidence/task-5-selector-tests.txt

  Scenario: Browser DOM exposes selectors for automation
    Tool: Playwright
    Preconditions: dev server running with sample plan
    Steps:
      1. Open the plan page
      2. Query selectors for top-level card, nested child item, summary block, and edit toggle
      3. Save the DOM query results/log
    Expected Result: All core selectors resolve exactly one or more intended elements
    Failure Indicators: Missing or duplicate ambiguous selectors
    Evidence: .sisyphus/evidence/task-5-selector-dom.txt
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-5-selector-tests.txt`
  - [ ] `.sisyphus/evidence/task-5-selector-dom.txt`

  **Commit**: YES
  - Message: `test(training): lock card-layout expectations before refactor`
  - Files: `frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts`, helper/fixture files
  - Pre-commit: `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`

- [ ] 6. Migrate `SimplePlanDisplay` to the shared card layout

  **What to do**:
  - Replace the read-only top-level table structure in `SimplePlanDisplay.vue` with the new card/list presentation using the shared primitives and selectors from Wave 1.
  - Keep save-to-history behavior intact.
  - Ensure total summary, title/description, and row content remain visible in the new hierarchy.

  **Must NOT do**:
  - Do not leave a parallel table-based rendering path in place for the normal read-only view.
  - Do not regress the save action in the footer.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: component is mostly read-only layout work.
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: needed to preserve hierarchy and visual consistency.
  - **Skills Evaluated but Omitted**:
    - `playwright`: verification only.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 7, 8, 9, 10)
  - **Blocks**: 11, 12, 13
  - **Blocked By**: 1, 2, 3, 5

  **References**:
  - `frontend/src/components/training/SimplePlanDisplay.vue:51-124` - Existing read-only render flow to replace.
  - `frontend/src/components/training/SimplePlanSubRows.vue:19-52` - Current nested read-only rendering to consolidate.
  - `frontend/web.html:214-407` - Desired top-level + nested + meta section pattern.

  **Acceptance Criteria**:
  - [ ] `SimplePlanDisplay.vue` uses card/list markup instead of a table.
  - [ ] Save button still emits the same payload.
  - [ ] Read-only cards expose the stable selectors defined in Task 5.

  **QA Scenarios**:
  ```text
  Scenario: Read-only simple display renders cards and still saves
    Tool: Bash
    Preconditions: component tests updated for simple display path
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Confirm assertions cover read-only card rendering and save action availability
    Expected Result: Tests pass for read-only card layout behavior
    Failure Indicators: Table assertions remain or save action is lost
    Evidence: .sisyphus/evidence/task-6-simple-display-tests.txt

  Scenario: Read-only display visually matches card hierarchy
    Tool: Playwright
    Preconditions: a shared/uploaded/read-only plan route available
    Steps:
      1. Open a read-only plan page
      2. Capture the plan header, first top-level card, and nested child items
      3. Verify save action remains present if the flow expects it
    Expected Result: Read-only screen uses the new hierarchy consistently
    Failure Indicators: Old table layout still appears or save affordance is missing incorrectly
    Evidence: .sisyphus/evidence/task-6-simple-display.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-6-simple-display-tests.txt`
  - [ ] `.sisyphus/evidence/task-6-simple-display.png`

  **Commit**: NO

- [ ] 7. Migrate `TrainingPlanDisplay` view mode to the shared card layout

  **What to do**:
  - Replace the top-level table structure in `TrainingPlanDisplay.vue` with the shared card/list layout for non-editing mode.
  - Preserve plan title, description, loading state, summary metrics, and button section.
  - Make the display feel closer to `web.html` while remaining consistent with Swim Gen page chrome and current CTA placement.

  **Must NOT do**:
  - Do not remove existing loading/no-plan states.
  - Do not move the entire page into a copied dashboard shell.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: main user-facing presentation shift.
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: useful for adapting the reference hierarchy into the current app.
  - **Skills Evaluated but Omitted**:
    - `playwright`: verification only.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 6, 8, 9, 10)
  - **Blocks**: 11, 12, 13
  - **Blocked By**: 1, 2, 3, 5

  **References**:
  - `frontend/src/components/training/TrainingPlanDisplay.vue:86-285` - Main view-mode rendering to replace.
  - `frontend/src/components/training/TrainingPlanDisplay.vue:287-616` - Existing header/summary/button styling to carry forward.
  - `frontend/src/views/HomeView.vue:128-147` - Main embedding context for the plan display.
  - `frontend/web.html:192-367` - Primary source for top-level card composition and metrics alignment.

  **Acceptance Criteria**:
  - [ ] View mode renders top-level cards and nested child blocks instead of a table.
  - [ ] Summary and CTA button section still render below the plan.
  - [ ] No-plan and loading states still render correctly.

  **QA Scenarios**:
  ```text
  Scenario: Main plan page renders card layout in view mode
    Tool: Playwright
    Preconditions: generated plan available on home route
    Steps:
      1. Open the home page with a populated plan
      2. Verify the plan area contains card selectors rather than table selectors
      3. Capture the full plan display
    Expected Result: Main display renders in card mode with summary and action buttons preserved
    Failure Indicators: Table still present, summary missing, or layout overflow
    Evidence: .sisyphus/evidence/task-7-main-display.png

  Scenario: Component tests lock in view-mode hierarchy
    Tool: Bash
    Preconditions: updated training plan display tests
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Verify assertions check title, description, card items, and absence of the no-plan placeholder when data exists
    Expected Result: View-mode tests pass against card layout expectations
    Failure Indicators: stale table assumptions or missing hierarchy assertions
    Evidence: .sisyphus/evidence/task-7-main-display-tests.txt
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-7-main-display.png`
  - [ ] `.sisyphus/evidence/task-7-main-display-tests.txt`

  **Commit**: NO

- [ ] 8. Rebuild editable row rendering for separate card-based edit mode

  **What to do**:
  - Replace table-cell editing in `TrainingPlanRow.vue` with card-compatible edit-mode inputs/controls.
  - Keep `toggleEditing`, `stopEditing`, `autoResize`, and row-path updates functioning with the new layout.
  - Ensure leaf vs parent distance-editability rules still apply.
  - Preserve content drill rendering in view mode and textareas/inputs in edit mode.

  **Must NOT do**:
  - Do not break `updatePlanRow`, `updatePlanRowEquipment`, or nested path propagation.
  - Do not make parent-row distance editable when it is derived from children.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: intertwined UI and behavior refactor with recursion.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: helpful, but preserving interaction correctness is primary.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 6, 7, 9, 10)
  - **Blocks**: 12, 13
  - **Blocked By**: 1, 3, 4, 5

  **References**:
  - `frontend/src/components/training/TrainingPlanRow.vue:27-55` - Current computed rules and store action handlers.
  - `frontend/src/components/training/TrainingPlanRow.vue:67-187` - Existing editable fields and event wiring.
  - `frontend/src/utils/rowHelpers.ts:100-226` - Sum recalculation and row mutation semantics that edit mode must preserve.

  **Acceptance Criteria**:
  - [ ] Edit mode uses card-based inputs/controls rather than table cells.
  - [ ] Valid numeric edits still call `store.updatePlanRow(path, field, value)` correctly.
  - [ ] Invalid numeric edits still normalize to `0` as they do now.
  - [ ] Parent rows keep computed distance behavior.

  **QA Scenarios**:
  ```text
  Scenario: Card edit mode updates numeric field correctly
    Tool: Bash
    Preconditions: updated Vitest specs for edit mode
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Verify a test toggles edit mode, changes the first amount input to `5`, and expects `updatePlanRow([0], 'Amount', 5)`
    Expected Result: Edit-mode card input updates the store exactly as before
    Failure Indicators: wrong path/field/value or missing input selectors
    Evidence: .sisyphus/evidence/task-8-edit-mode-tests.txt

  Scenario: Card edit mode handles invalid numeric input gracefully
    Tool: Playwright
    Preconditions: editable plan loaded
    Steps:
      1. Enter edit mode
      2. Input `abc` into the first amount field and blur
      3. Capture resulting UI and any visible normalized value/state
    Expected Result: Invalid number is rejected/normalized without crashing the UI
    Failure Indicators: uncaught error, broken card, or stale invalid value handling
    Evidence: .sisyphus/evidence/task-8-invalid-edit.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-8-edit-mode-tests.txt`
  - [ ] `.sisyphus/evidence/task-8-invalid-edit.png`

  **Commit**: YES
  - Message: `refactor(training): replace plan tables with shared card layout`
  - Files: `frontend/src/components/training/TrainingPlanDisplay.vue`, `TrainingPlanRow.vue`, `SimplePlanDisplay.vue`, `SimplePlanSubRows.vue`, shared helpers/styles/tests
  - Pre-commit: `npm run test:unit && npm run type-check`

- [ ] 9. Redesign equipment, totals, and footer/meta sections in card style

  **What to do**:
  - Convert summary/footer/meta presentation to match the reference hierarchy: compact metric blocks near the header, equipment chips/cards, and note-like content blocks where applicable.
  - Ensure the existing total meters and exercise-set counts remain visible and meaningful.
  - Aggregate equipment consistently across the plan where the current UX expects it.

  **Must NOT do**:
  - Do not hide totals solely inside nested cards.
  - Do not lose equipment visibility when row content moves away from table cells.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: summary/meta sections are presentational and hierarchy-driven.
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: useful for translating image/web footer patterns into native components.
  - **Skills Evaluated but Omitted**:
    - `playwright`: verification only.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 6, 7, 8, 10)
  - **Blocks**: 11, 12, 13
  - **Blocked By**: 2

  **References**:
  - `frontend/src/components/training/TrainingPlanDisplay.vue:257-284` - Existing summary and button section structure.
  - `frontend/src/components/training/SimplePlanDisplay.vue:112-123` - Existing compact footer structure.
  - `frontend/web.html:368-407` - Equipment and notes section treatment.
  - `frontend/mobile.html:290-325` - Compact mobile equipment/note pattern.

  **Acceptance Criteria**:
  - [ ] Totals and exercise-set counts remain visible in the redesigned layout.
  - [ ] Equipment is rendered in chip/card form without depending on table content cells.
  - [ ] Header/summary/meta hierarchy is consistent across editable and read-only screens.

  **QA Scenarios**:
  ```text
  Scenario: Equipment and totals appear in redesigned meta sections
    Tool: Playwright
    Preconditions: plan with equipment-bearing rows
    Steps:
      1. Open the plan page
      2. Verify total meters block, exercise-count block, and equipment chip section are visible
      3. Capture screenshot of the lower meta region
    Expected Result: Summary and equipment remain visible in card-style sections
    Failure Indicators: totals lost, equipment hidden, or footer collapsed
    Evidence: .sisyphus/evidence/task-9-meta-sections.png

  Scenario: Tests cover footer/meta rendering
    Tool: Bash
    Preconditions: specs updated
    Steps:
      1. Run `npm run test:unit -- src/components/training/__tests__/TrainingPlanDisplay.spec.ts`
      2. Confirm assertions check totals/exercise counts/equipment output
    Expected Result: Meta-section behavior is covered in component tests
    Failure Indicators: no assertions for redesigned meta blocks
    Evidence: .sisyphus/evidence/task-9-meta-tests.txt
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-9-meta-sections.png`
  - [ ] `.sisyphus/evidence/task-9-meta-tests.txt`

  **Commit**: NO

- [ ] 10. Complete responsive and deep-nesting styling pass

  **What to do**:
  - Replace the old table zoom/mobile strategy with responsive card behavior that works at narrow widths without layout collapse.
  - Tune spacing, typography, metric wrapping, and nested containment for depth 1-4.
  - Verify that long content and drill-link text wrap sensibly inside cards.

  **Must NOT do**:
  - Do not depend on `zoom: 0.75` as the primary mobile strategy for the new layout.
  - Do not allow nested depth to become visually indistinguishable or unreadable.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: responsive refinement and visual tuning.
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: useful for balancing compactness and readability.
  - **Skills Evaluated but Omitted**:
    - `playwright`: verification only.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 6, 7, 8, 9)
  - **Blocks**: 11, 12, 13
  - **Blocked By**: 2, 3

  **References**:
  - `frontend/src/components/training/TrainingPlanDisplay.vue:302-306, 425-439, 600-606` - Existing mobile behavior to retire/replace.
  - `frontend/mobile.html:128-325` - Compact mobile card behavior.
  - `frontend/src/components/training/ContentWithDrillLinks.vue:16-33` - Inline content rendering that must wrap correctly in the new layout.

  **Acceptance Criteria**:
  - [ ] New layout is usable at mobile widths without global zoom hacks.
  - [ ] Long content, drill links, and equipment badges wrap without overlap.
  - [ ] Nested levels remain distinguishable through depth 4 on desktop and mobile.

  **QA Scenarios**:
  ```text
  Scenario: Mobile viewport stays readable without table zoom fallback
    Tool: Playwright
    Preconditions: responsive card layout implemented
    Steps:
      1. Open the plan page at 390x844
      2. Verify no horizontal scroll is required for the main plan content
      3. Capture screenshot of a nested card region with long content
    Expected Result: Mobile presentation is readable and contained
    Failure Indicators: horizontal overflow, clipped text, or reliance on shrunken zoomed layout
    Evidence: .sisyphus/evidence/task-10-mobile-responsive.png

  Scenario: Desktop deep nesting remains visually distinct
    Tool: Playwright
    Preconditions: sample plan with depth-4 nesting
    Steps:
      1. Open the plan page at 1440x1200
      2. Capture a depth-4 nested region
      3. Verify each level remains distinguishable through container/border/spacing differences
    Expected Result: Hierarchy is visible through depth 4
    Failure Indicators: nested levels blur together or child blocks escape parent containers
    Evidence: .sisyphus/evidence/task-10-depth4-desktop.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-10-mobile-responsive.png`
  - [ ] `.sisyphus/evidence/task-10-depth4-desktop.png`

  **Commit**: NO

- [ ] 11. Integrate redesigned display across Home, Shared, and Uploaded plan screens

  **What to do**:
  - Verify the redesigned display fits correctly inside `HomeView.vue`, `SharedView.vue`, and `UploadedPlanView.vue` containers without spacing regressions.
  - Adjust surrounding layout spacing only where necessary to accommodate the new card presentation.
  - Ensure shared/uploaded plan screens still route users into conversation flows correctly after viewing the redesigned plan.

  **Must NOT do**:
  - Do not redesign the page hero/CTA structure beyond what's needed to fit the new display.
  - Do not break login or conversation-start flows on shared/uploaded pages.

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: bounded integration polish across a few views.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: core layout system already handled earlier.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 12, 13)
  - **Blocks**: F1-F4
  - **Blocked By**: 6, 7, 9, 10

  **References**:
  - `frontend/src/views/HomeView.vue:128-151` - Main page container around `TrainingPlanDisplay`.
  - `frontend/src/views/SharedView.vue:112-145` - Shared plan presentation and chat transition area.
  - `frontend/src/views/UploadedPlanView.vue:109-129` - Uploaded-plan presentation and chat transition area.

  **Acceptance Criteria**:
  - [ ] Home, shared, and uploaded plan screens all display the new layout correctly.
  - [ ] Their surrounding CTA/chat sections remain functional and visually separated.

  **QA Scenarios**:
  ```text
  Scenario: Home route contains redesigned plan without surrounding layout regressions
    Tool: Playwright
    Preconditions: generated plan available on home route
    Steps:
      1. Open home route with a plan
      2. Capture hero, plan display, and CTA banner together
      3. Verify spacing between sections remains consistent
    Expected Result: New plan layout fits naturally into the page container
    Failure Indicators: overlapping sections, collapsed margins, or broken CTA placement
    Evidence: .sisyphus/evidence/task-11-home-integration.png

  Scenario: Shared/uploaded routes still allow conversation handoff
    Tool: Playwright
    Preconditions: shared/uploaded route fixtures available
    Steps:
      1. Open a shared or uploaded plan route
      2. Verify plan renders in new layout
      3. Interact with the conversation starter input/button area
    Expected Result: Viewing the new layout does not break downstream conversation entry
    Failure Indicators: missing input/button, layout collision, or unresponsive controls
    Evidence: .sisyphus/evidence/task-11-shared-uploaded-integration.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-11-home-integration.png`
  - [ ] `.sisyphus/evidence/task-11-shared-uploaded-integration.png`

  **Commit**: NO

- [ ] 12. Restore interaction-view and snapshot parity with the redesigned plan display

  **What to do**:
  - Verify the redesigned training plan works correctly inside `InteractionView.vue`, including the plan tab and any snapshot-rendered plan content.
  - Ensure edit mode, read-only snapshots, and tab-switching still behave correctly with the new shared display system.
  - Confirm no assumptions remain about table markup in interaction-related tests or logic.

  **Must NOT do**:
  - Do not break snapshot save/view flows.
  - Do not regress plan/chat tab behavior.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: cross-feature integration with editable and read-only variants.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: visual system already established; this is integration-heavy.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 11, 13)
  - **Blocks**: F1-F4
  - **Blocked By**: 6, 7, 8, 9, 10

  **References**:
  - `frontend/src/views/InteractionView.vue:195-220` - Tab shell around the plan display.
  - `frontend/src/views/InteractionView.vue` (remainder) - Snapshot behavior and plan/chat switching context.
  - `frontend/src/components/training/SimplePlanDisplay.vue:41-48` - Save snapshot payload behavior to preserve.

  **Acceptance Criteria**:
  - [ ] Interaction view plan tab renders the new layout correctly.
  - [ ] Snapshot/read-only plan renderings match the shared card system.
  - [ ] Edit mode and tab switching remain functional.

  **QA Scenarios**:
  ```text
  Scenario: Interaction view plan tab shows redesigned plan and still switches tabs
    Tool: Playwright
    Preconditions: interaction route with existing plan available
    Steps:
      1. Open interaction route
      2. Verify the plan tab shows the card layout
      3. Switch to chat tab and back to plan tab
      4. Capture the restored plan tab state
    Expected Result: Tab switching preserves the redesigned plan display
    Failure Indicators: blank tab, broken hydration, or stale table markup
    Evidence: .sisyphus/evidence/task-12-interaction-tabs.png

  Scenario: Snapshot read-only rendering matches shared layout
    Tool: Playwright
    Preconditions: interaction view with at least one plan snapshot message
    Steps:
      1. Expand a snapshot in the conversation
      2. Verify the snapshot uses the card layout rather than table markup
      3. Capture screenshot of the snapshot panel
    Expected Result: Snapshot rendering is layout-consistent with read-only plan views
    Failure Indicators: snapshot still renders in old table mode or breaks sizing
    Evidence: .sisyphus/evidence/task-12-snapshot-layout.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-12-interaction-tabs.png`
  - [ ] `.sisyphus/evidence/task-12-snapshot-layout.png`

  **Commit**: NO

- [ ] 13. Harden edge cases, accessibility, and regression coverage

  **What to do**:
  - Add/finish test coverage for loading state, no-plan placeholder, long content, drill-link content, parent rows with computed distance, equipment-heavy rows, and max-depth nesting.
  - Validate accessible semantics for edit toggle, card controls, and summary/meta content.
  - Run the full frontend validation stack and fix regressions.

  **Must NOT do**:
  - Do not leave old table CSS/classes creating dead styling conflicts.
  - Do not ship without explicit checks for loading/no-plan states.

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: broad regression hardening across many edge cases.
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `playwright`: used in QA, but this task is primarily coverage and correctness.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 11, 12)
  - **Blocks**: F1-F4
  - **Blocked By**: 1, 3, 4, 5, 6, 7, 8, 9, 10

  **References**:
  - `frontend/src/components/training/TrainingPlanDisplay.vue:87-92, 270-283` - Loading/no-plan/button-state handling to preserve.
  - `frontend/src/components/training/ContentWithDrillLinks.vue:16-33` - Drill-link rendering that must survive card wrapping.
  - `frontend/vitest.config.ts:1-17` - Test runner configuration.
  - `.github/workflows/frontend-validate.yaml:82-92` - CI validation commands that must stay green.
  - `frontend/package.json:9-20` - Canonical frontend validation commands.

  **Acceptance Criteria**:
  - [ ] Component tests cover loading, no-plan, edit mode, nested depth, and long-content regressions.
  - [ ] Accessibility-sensitive controls remain keyboard-usable and visibly labeled.
  - [ ] Full frontend validation passes.

  **QA Scenarios**:
  ```text
  Scenario: Full frontend validation passes after refactor
    Tool: Bash
    Preconditions: all refactor work complete
    Steps:
      1. Run `npm run test:unit`
      2. Run `npm run type-check`
      3. Run `npm run lint`
      4. Run `npm run build`
      5. Save combined outputs
    Expected Result: All commands exit successfully
    Failure Indicators: test, type, lint, or build failures
    Evidence: .sisyphus/evidence/task-13-full-validation.txt

  Scenario: Keyboard and content edge cases remain usable
    Tool: Playwright
    Preconditions: plan page available with long text and nested rows
    Steps:
      1. Navigate via keyboard to `.edit-btn` and first card action control
      2. Verify visible focus styles exist
      3. Capture a long-content nested card region and a focused control state
    Expected Result: Controls are focusable and long content remains readable
    Failure Indicators: focus trap, invisible focus, clipped text, or inaccessible controls
    Evidence: .sisyphus/evidence/task-13-a11y-edge-cases.png
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-13-full-validation.txt`
  - [ ] `.sisyphus/evidence/task-13-a11y-edge-cases.png`

  **Commit**: YES
  - Message: `polish(training): restore edit flows and responsive regressions`
  - Files: tests, view integrations, responsive/a11y cleanup
  - Pre-commit: `npm run test:unit && npm run type-check && npm run lint && npm run build`

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Read the plan end-to-end. Verify the card layout replaced table-first rendering where required, edit mode still exists, and all mandated selectors/evidence files exist.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `npm run type-check`, `npm run lint`, `npm run test:unit`, and `npm run build` in `frontend/`. Review changed files for dead table selectors, broken recursion, unused imports, and styling duplication.
  Output: `Build [PASS/FAIL] | Lint [PASS/FAIL] | Tests [PASS/FAIL] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high` + `playwright`
  Execute all QA scenarios in this plan across desktop and mobile widths, save screenshots/logs to `.sisyphus/evidence/final-qa/`, and verify both read-only and edit-mode flows.
  Output: `Scenarios [N/N pass] | Responsive [PASS/FAIL] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  Compare final diff against the requested scope: display redesign only, no unrelated shell redesign, no backend/schema changes, no Tailwind adoption.
  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | VERDICT`

---

## Commit Strategy

- **1**: `test(training): lock card-layout expectations before refactor`
- **2**: `refactor(training): replace plan tables with shared card layout`
- **3**: `polish(training): restore edit flows and responsive regressions`

---

## Success Criteria

### Verification Commands
```bash
npm run test:unit
npm run type-check
npm run lint
npm run build
```

### Final Checklist
- [ ] Editable and read-only plan views share the same card hierarchy
- [ ] Nested sets remain readable and visibly parent-contained through depth 4
- [ ] Edit mode remains available and functional from the top-level toggle
- [ ] Equipment, totals, and summary information remain visible without table columns
- [ ] No Tailwind or unrelated shell redesign was introduced
