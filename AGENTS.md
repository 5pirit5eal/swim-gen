# Style & Frontend Reference Guide

This document defines the UI/UX conventions, design system, and styling rules for the `swim-gen` frontend (`frontend/src`). All agent tasks modifying or creating frontend components must follow these established patterns.

---

## 1. Design System & Theme Architecture

The application uses a **custom CSS variable-based design system** built on top of native CSS and Vue 3 Scoped Styles (no Tailwind or external UI component frameworks).

### Dynamic Theming & Light/Dark Mode

Theme variables are defined in `frontend/src/assets/base.css` and automatically switch based on `@media (prefers-color-scheme: dark)`.

| Token / Variable | Light Mode | Dark Mode | Purpose |
| :--- | :--- | :--- | :--- |
| `--color-background` | `#ffffff` | `#0b171f` | Main page background / solid surfaces |
| `--color-background-soft` | `#efefef` | `#132330` | Default card & container background |
| `--color-background-mute` | `#e9e9e9` | `#152531` | Nested sections, subtle input backgrounds |
| `--color-transparent` | `rgba(237, 237, 237, 0.25)` | `rgba(15, 34, 50, 0.75)` | Frosted glass containers over background |
| `--color-heading` | `#2c3e50` (Indigo) | `#ffffff` | Primary headings and prominent labels |
| `--color-text` | `rgba(37, 36, 36, 0.9)` | `rgba(235, 245, 250, 0.85)` | Primary body text and general labels |
| `--color-primary` | `#1a6d93` | `#254d6d` | Primary CTA buttons, active states, badges |
| `--color-primary-hover` | `#145775` | `#1a3c5b` | Hover/active states for primary actions |
| `--color-border` | `rgba(60, 60, 60, 0.12)` | `rgba(61, 73, 87, 0.48)` | Standard element & card borders |
| `--color-border-hover` | `rgba(60, 60, 60, 0.29)` | `rgba(59, 70, 84, 0.65)` | Hover state for borders |
| `--color-shadow` | `rgba(47, 124, 248, 0.1)` | `rgba(56, 191, 248, 0.12)` | Drop shadows and focus glow |
| `--color-success` | `#07bc0c` | `#07bc0c` | Success indicators, save confirmation |
| `--color-error` | `#e74c3c` / `#b91c1c` | `#e74c3c` / `#b91c1c` | Error messages, destructive buttons |

### Background Water Layer

`AppLayout.vue` provides a fixed swimming pool backdrop (`light_mode.webp` in light mode, `dark_mode.webp` in dark mode). Transparent frosted cards (`--color-transparent` with `backdrop-filter: blur(...)`) sit above this layer.

---

## 2. Typography & Core Layout

- **Font Family**:

  ```css
  font-family: Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  ```

- **Base Body Size**: `15px` with line height `1.6`.
- **Headings**:
  - `h1` (Hero / Views): `2.5rem` (`2rem` on mobile), `font-weight: 700`, `color: var(--color-heading)`.
  - `h2` (Section / Dialog): `1.75rem`, `font-weight: 600`, `color: var(--color-heading)`.
  - `h3` (Card Titles / Subsections): `1.25rem` – `1.5rem`, `font-weight: 600`.
- **Page Container**:

  ```css
  .container {
    max-width: 1080px;
    margin: 0 auto;
    padding: 0 1rem;
  }
  ```

---

## 3. Element Styling Patterns

### 1. Buttons

Always provide disabled, hover, and focus states. Keep touch targets $\ge 44\text{px}$ on mobile.

#### Primary Action Button (`.btn-primary` / `.submit-btn` / `.cta-button`)

```css
.btn-primary,
.submit-btn {
  background-color: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.btn-primary:hover:not(:disabled),
.submit-btn:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-primary:disabled,
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
```

#### Secondary Button (`.btn-secondary` / `.toggle-settings-btn`)

```css
.btn-secondary {
  background-color: var(--color-background);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-secondary:hover:not(:disabled) {
  background-color: var(--color-background-soft);
  border-color: var(--color-border-hover);
  color: var(--color-primary);
}
```

#### Destructive Button (`.delete-btn` / `.confirm-delete-btn`)

```css
.delete-btn {
  background-color: var(--color-error);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.delete-btn:hover:not(:disabled) {
  background-color: var(--color-error-soft);
}
```

#### Icon-Only Action Button (`.icon-btn` / `.action-button`)

```css
.icon-btn {
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  padding: 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: var(--color-background-mute);
  color: var(--color-primary);
}
```

#### Text / Link Button (`.text-btn`)

```css
.text-btn {
  background: none;
  border: none;
  color: var(--color-text);
  padding: 0;
  font-size: 1rem;
  text-decoration: underline;
  cursor: pointer;
}

.text-btn:hover {
  color: var(--color-primary);
}
```

---

### 2. Form Inputs & Controls

#### Text Inputs & Textareas (`input[type='text']`, `textarea`)

```css
.form-input,
.form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 1rem;
  background-color: var(--color-background);
  color: var(--color-text);
  box-sizing: border-box;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-shadow);
}

.form-input:disabled,
.form-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: color-mix(in srgb, var(--color-text), transparent 40%);
}
```

#### Select Dropdowns

```css
.select-input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.9rem;
  background: var(--color-background);
  color: var(--color-text);
}

.select-input:focus {
  outline: none;
  border-color: var(--color-primary);
}
```

#### Radio and Checkbox Options

```css
.radio-option,
.checkbox-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--color-text);
}

.radio-option:hover,
.checkbox-option:hover {
  color: var(--color-heading);
}
```

#### Error and Validation States

- Form error banner: Light red background (`#fef2f2`), border `#fecaca`, text `var(--color-error)`, `border-radius: 8px`, `padding: 0.75rem`.
- Field error text: `.field-error` with `color: var(--color-error); font-size: 0.8rem; margin-top: 0.3rem;`.
- Invalid input border: `border-color: var(--color-error);`.

---

### 3. Cards & Content Panels

#### Standard Content Card (`.profile-card`, `.form-container`, `.cta-banner`)

```css
.card {
  background: var(--color-background-soft);
  padding: 2rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  box-sizing: border-box;
}

@media (max-width: 740px) {
  .card {
    padding: 1.25rem;
  }
}
```

#### Frosted Hero Panel (`.hero`)

```css
.hero {
  text-align: center;
  background-color: var(--color-transparent);
  backdrop-filter: blur(2px);
  border-radius: 8px;
  padding: 1rem;
  margin: 2rem auto;
}
```

#### Interactive Media Card (`.drill-card`)

```css
.drill-card {
  background: var(--color-background-soft);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--color-border);
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.drill-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px -10px var(--color-shadow);
  border-color: var(--color-primary);
}
```

---

### 4. Badges, Tags & Chips

#### Accent Badges & Chips

```css
.badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
  background: var(--color-primary);
  color: white;
  letter-spacing: 0.5px;
  white-space: nowrap;
  box-shadow: 0 1px 3px var(--color-shadow);
}
```

#### Overlay Badges (Images)

```css
.image-overlay-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  padding: 4px 8px;
  border-radius: 6px;
  letter-spacing: 0.05em;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  z-index: 2;
  background-color: var(--color-primary);
  color: white;
}
```

---

### 5. Modals, Dropdowns & Tooltips

- **Modals (`BaseModal.vue`)**: Teleported to `body`. Fixed backdrop `rgba(...)` with `backdrop-filter: blur(3px)`. Container width `90%`, `max-width: 1000px`, `border-radius: 8px`, `box-shadow: 0 2px 10px var(--color-shadow)`. Header, Body, and Footer sections divided by `1px solid var(--color-border)`.
- **Dropdown Menus (`.dropdown-menu`)**: `position: absolute`, `background-color: var(--color-background-soft)`, `border: 1px solid var(--color-border)`, `border-radius: 8px`, `box-shadow: 0 4px 6px var(--color-shadow)`, `z-index: 1000`.
- **Tooltips (`BaseTooltip.vue`)**: Fixed overlay computed dynamically via `calculateOverlayPosition`, `background-color: var(--color-background-mute)`, `border-radius: 6px`, `box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15)`, `z-index: 9999`.

---

### 6. Spinners & Loading States

```css
.loading-spinner {
  width: 1rem;
  height: 1rem;
  border: 2px solid rgb(255 255 255 / 45%);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
```

---

## 4. Responsive Breakpoints & Viewport Rules

- **Desktop** (> 1124px): Multi-column layouts, 400px sidebars, standard font sizes.
- **Laptop / Tablet** (741px – 1124px): `zoom: 0.85` on layout, 300px sidebar, 2-column grids collapse gracefully.
- **Mobile** ($\le$ 740px / 390px / 320px):
  - `.app-layout`: `zoom: 0.75` and `background-attachment: scroll`.
  - Full-width sidebar (`width: 100%`).
  - Single-column stacked form grids (`grid-template-columns: 1fr`).
  - Minimum touch target size: 44px $\times$ 44px for buttons.
  - Horizontal overflow is prohibited (`overflow-x: hidden` / wrap chips and metrics).

---

## 5. Accessibility & Internationalization (i18n) Rules

1. **Color Contrast**: All text must meet WCAG AA contrast against its background in both light and dark modes. Avoid dark text directly over the pool photo without a backdrop/card container.
2. **Focus Visibility**: Keep explicit focus outlines or box-shadow halos (`0 0 0 2px var(--color-shadow)`) on interactive elements.
3. **i18n Requirements**:
   - **Never** hardcode user-facing strings in `.vue` files.
   - Always add translation keys to both `frontend/src/locales/en.json` and `frontend/src/locales/de.json`.
   - German strings are typically 20–30% longer than English. Ensure button layouts and grid columns use flex-wrap or auto-sizing to prevent truncation or overflow.
4. **Reduced Motion**: Respect `prefers-reduced-motion` for transitions and smooth scrolling.
