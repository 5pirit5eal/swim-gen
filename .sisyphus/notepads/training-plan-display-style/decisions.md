# Decisions — training-plan-display-style

## [2026-03-21] Initial Decisions

- Use TDD: write failing tests first, then implement
- Global edit toggle (not per-row) is the default entry point for edit mode
- CSS: use existing CSS variable/theme system only
- `BaseTableAction` must be fully redesigned for card layout (not just restyled)
- Follow `web.html` more closely than `mobile.html`
- Do NOT infer semantic phase labels from row data
- Shared recursive renderer for both editable and read-only displays
