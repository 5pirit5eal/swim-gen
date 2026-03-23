<script setup lang="ts">
import ContentWithDrillLinks from '@/components/training/ContentWithDrillLinks.vue'
import type { Row, PlanStore } from '@/types'
import { EQUIPMENT_I18N_KEYS, EQUIPMENT_TYPES, MAX_NESTING_DEPTH } from '@/utils/rowHelpers'
import { computed } from 'vue'
import Multiselect from 'vue-multiselect'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  row: Row
  path: number[]
  depth: number
  isEditing: boolean
  store: PlanStore
  isFirst: boolean
  isLast: boolean
}>()

const { t } = useI18n()

type EquipmentOption = {
  value: string
  label: string
}

const equipmentOptions = computed<EquipmentOption[]>(() =>
  EQUIPMENT_TYPES.map((equipment) => ({
    value: equipment,
    label: getEquipmentLabel(equipment),
  })),
)

// ── Computed helpers ──────────────────────────────────────────────────────────

const hasSubRows = computed(() => props.row.SubRows && props.row.SubRows.length > 0)
const canAddSubRow = computed(() => props.depth < MAX_NESTING_DEPTH)
const hasEquipment = computed(() => props.row.Equipment && props.row.Equipment.length > 0)
const shouldShowEquipmentMetric = computed(() => props.isEditing || hasEquipment.value)

/** Parent rows (with SubRows) have a computed Distance — not directly editable */
const isDistanceEditable = computed(() => !hasSubRows.value)

// ── Store action handlers ─────────────────────────────────────────────────────

function handleAddRow() {
  props.store.addRow(props.path)
}

function handleRemoveRow() {
  props.store.removeRow(props.path)
}

function handleMoveRow(direction: 'up' | 'down') {
  props.store.moveRow(props.path, direction)
}

function handleAddSubRow() {
  props.store.addSubRow(props.path, props.depth)
}

function handleFieldBlur(event: Event, field: keyof Row) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  const rawValue = target.value

  if (field === 'Amount' || field === 'Distance') {
    const parsed = parseInt(rawValue, 10)
    const value = isNaN(parsed) ? 0 : parsed
    props.store.updatePlanRow(props.path, field, value)
  } else {
    props.store.updatePlanRow(props.path, field, rawValue)
  }
}

const selectedEquipmentOptions = computed<EquipmentOption[]>({
  get() {
    const selectedEquipment = props.row.Equipment ?? []
    return equipmentOptions.value.filter((option) => selectedEquipment.includes(option.value))
  },
  set(selectedOptions) {
    props.store.updatePlanRowEquipment(
      props.path,
      selectedOptions.map((option) => option.value),
    )
  },
})

function getEquipmentLabel(equipment: string): string {
  const translationKey = EQUIPMENT_I18N_KEYS[equipment as keyof typeof EQUIPMENT_I18N_KEYS]
  return translationKey ? t(`equipment.${translationKey}`) : equipment
}

// ── Sub-row path helper ───────────────────────────────────────────────────────

function subRowPath(subIndex: number): number[] {
  return [...props.path, subIndex]
}
</script>

<template>
  <div
    class="plan-row-card"
    :class="[`plan-row-card--depth-${depth}`, { 'plan-row-card--parent': hasSubRows }]"
    :data-testid="depth === 0 ? 'plan-card' : 'plan-card-nested'"
  >
    <!-- ── Content body ─────────────────────────────────────────────────── -->
    <div>
      <textarea
        v-if="isEditing"
        :value="row.Content"
        @blur="handleFieldBlur($event, 'Content')"
        @keyup.enter.prevent="handleFieldBlur($event, 'Content')"
        class="plan-row-card__textarea"
        rows="2"
        :aria-label="t('display.content')"
      ></textarea>
      <div v-else class="plan-row-card__content-view">
        <ContentWithDrillLinks :content="row.Content" />
      </div>
    </div>

    <!-- ── Card header: metrics + actions ──────────────────────────────── -->
    <div class="plan-row-card__data">
      <!-- Metrics row -->
      <div class="plan-row-card__metrics">
        <!-- Amount -->
        <div class="plan-row-card__metric">
          <span class="plan-row-card__metric-label">{{ t('display.amount') }}</span>
          <input
            v-if="isEditing"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            :value="row.Amount"
            @blur="handleFieldBlur($event, 'Amount')"
            @keyup.enter="handleFieldBlur($event, 'Amount')"
            class="plan-row-card__input plan-row-card__input--small"
            :aria-label="t('display.amount')"
          />
          <span v-else class="plan-row-card__metric-value">{{ row.Amount }}</span>
        </div>

        <!-- Multiplier placeholder to preserve flex layout without visual redundancy -->
        <div class="plan-row-card__metric">
          <span
            class="plan-row-card__metric-label plan-row-card__metric--placeholder"
            aria-hidden="true"
            >{{ t('display.multiplier') }}</span
          >
          <span class="plan-row-card__metric-value">{{ row.Multiplier }}</span>
        </div>

        <!-- Distance -->
        <div class="plan-row-card__metric">
          <span class="plan-row-card__metric-label">{{ t('display.distance') }}</span>
          <input
            v-if="isEditing && isDistanceEditable"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            :value="row.Distance"
            @blur="handleFieldBlur($event, 'Distance')"
            @keyup.enter="handleFieldBlur($event, 'Distance')"
            class="plan-row-card__input plan-row-card__input--small"
            :aria-label="t('display.distance')"
          />
          <span v-else class="plan-row-card__metric-value">{{ row.Distance }}</span>
        </div>

        <!-- Break -->
        <div class="plan-row-card__metric">
          <span class="plan-row-card__metric-label">{{ t('display.break') }}</span>
          <input
            v-if="isEditing"
            type="text"
            :value="row.Break"
            @blur="handleFieldBlur($event, 'Break')"
            @keyup.enter="handleFieldBlur($event, 'Break')"
            class="plan-row-card__input plan-row-card__input--small"
            :aria-label="t('display.break')"
          />
          <span v-else class="plan-row-card__metric-value">{{ row.Break }}</span>
        </div>

        <!-- Intensity -->
        <div class="plan-row-card__metric">
          <span class="plan-row-card__metric-label">{{ t('display.intensity') }}</span>
          <input
            v-if="isEditing"
            type="text"
            :value="row.Intensity"
            @blur="handleFieldBlur($event, 'Intensity')"
            @keyup.enter="handleFieldBlur($event, 'Intensity')"
            class="plan-row-card__input plan-row-card__input--small"
            :aria-label="t('display.intensity')"
          />
          <span v-else class="plan-row-card__metric-value plan-row-card__metric-value--intensity">{{
            row.Intensity
          }}</span>
        </div>

        <!-- Equipment badges — inline with metrics on wide screens, wraps below on narrow -->
        <div
          v-if="shouldShowEquipmentMetric"
          data-testid="equipment-metric"
          class="plan-row-card__metric plan-row-card__metric--equipment"
        >
          <span class="plan-row-card__metric-label">{{ t('display.equipment') }}</span>
          <Multiselect
            v-if="isEditing"
            v-model="selectedEquipmentOptions"
            :options="equipmentOptions"
            :multiple="true"
            :close-on-select="false"
            :clear-on-select="false"
            :preserve-search="true"
            :show-labels="false"
            :placeholder="t('display.equipment')"
            label="label"
            track-by="value"
            class="plan-row-card__multiselect"
            data-testid="equipment-multiselect"
            :aria-label="t('display.equipment')"
          />
          <span v-else class="plan-row-card__equipment-badges">
            <span v-for="eq in row.Equipment" :key="eq" class="plan-row-card__equipment-badge">{{
              getEquipmentLabel(eq)
            }}</span>
          </span>
        </div>

        <!-- Sum (always read-only — computed by store) -->
        <div class="plan-row-card__metric plan-row-card__metric--sum">
          <span class="plan-row-card__metric-label">{{ t('display.sum') }}</span>
          <span class="plan-row-card__metric-value plan-row-card__metric-value--sum">{{
            row.Sum
          }}</span>
        </div>
      </div>

      <!-- Edit-mode inline action controls (NOT hover-only) -->
      <div
        v-if="isEditing"
        class="plan-row-card__actions"
        role="toolbar"
        :aria-label="t('display.row_actions')"
      >
        <button
          class="plan-row-card__action-btn plan-row-card__action-btn--move-up"
          :disabled="isFirst"
          @click.stop="handleMoveRow('up')"
          :title="t('display.move_row_up')"
          :aria-label="t('display.move_row_up')"
        >
          <!-- Up arrow pseudo-element via CSS; inner text for a11y -->
          <span
            aria-hidden="true"
            class="plan-row-card__action-icon plan-row-card__action-icon--up"
          ></span>
        </button>
        <button
          class="plan-row-card__action-btn plan-row-card__action-btn--move-down"
          :disabled="isLast"
          @click.stop="handleMoveRow('down')"
          :title="t('display.move_row_down')"
          :aria-label="t('display.move_row_down')"
        >
          <span
            aria-hidden="true"
            class="plan-row-card__action-icon plan-row-card__action-icon--down"
          ></span>
        </button>
        <button
          class="plan-row-card__action-btn plan-row-card__action-btn--add"
          @click.stop="handleAddRow"
          :title="t('display.add_row')"
          :aria-label="t('display.add_row')"
        >
          <span
            aria-hidden="true"
            class="plan-row-card__action-icon plan-row-card__action-icon--add"
          ></span>
        </button>
        <button
          class="plan-row-card__action-btn plan-row-card__action-btn--remove"
          @click.stop="handleRemoveRow"
          :title="t('display.remove_row')"
          :aria-label="t('display.remove_row')"
        >
          <span
            aria-hidden="true"
            class="plan-row-card__action-icon plan-row-card__action-icon--remove"
          ></span>
        </button>
        <button
          v-if="canAddSubRow"
          class="plan-row-card__action-btn plan-row-card__action-btn--add-subrow"
          @click.stop="handleAddSubRow"
          :title="t('display.add_subrow')"
          :aria-label="t('display.add_subrow')"
        >
          <span
            aria-hidden="true"
            class="plan-row-card__action-icon plan-row-card__action-icon--subrow"
          ></span>
        </button>
      </div>
    </div>

    <!-- ── Nested SubRows ───────────────────────────────────────────────── -->
    <TransitionGroup tag="div" name="list" v-if="hasSubRows" class="plan-row-card__subrows">
      <PlanRowCard
        v-for="(subRow, subIndex) in row.SubRows"
        :key="subRow._id || subIndex"
        :row="subRow"
        :path="subRowPath(subIndex)"
        :depth="depth + 1"
        :is-editing="isEditing"
        :store="store"
        :is-first="subIndex === 0"
        :is-last="subIndex === (row.SubRows?.length ?? 1) - 1"
      />
    </TransitionGroup>
  </div>
</template>

<script lang="ts">
// Self-referencing for recursive component
export default {
  name: 'PlanRowCard',
}
</script>

<style scoped>
/* ── Card shell ─────────────────────────────────────────────────────────────── */

.plan-row-card {
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  position: relative;
}

/* Parent cards get a left accent */
.plan-row-card--parent {
  background: var(--color-background-mute);
}

/* Depth-specific tinting/indentation — gradient of prominence depth 0→4 */

/* Depth 0: top-level, most prominent — strong primary accent */
.plan-row-card--depth-0 {
  border-left: 4px solid var(--color-primary);
  background: var(--color-background);
  box-shadow: 0 1px 6px var(--color-shadow);
}

/* Depth 1: first-level nesting — solid but lighter */
.plan-row-card--depth-1 {
  border-left: 3px solid var(--color-primary);
  margin-left: 0.75rem;
  background: var(--color-background-soft);
  font-size: 0.93rem;
}

/* Depth 2: second-level nesting — muted border, more indent */
.plan-row-card--depth-2 {
  border-left: 2px solid var(--color-border-hover);
  margin-left: 0.5rem;
  background: var(--color-background-mute);
  font-size: 0.875rem;
}

/* Depth 3: very nested — subtle border */
.plan-row-card--depth-3 {
  border-left: 2px solid var(--color-border);
  margin-left: 0.35rem;
  background: var(--color-background-soft);
  font-size: 0.82rem;
}

/* Depth 4: maximum depth — minimal, near-invisible border */
.plan-row-card--depth-4 {
  border-left: 1px solid var(--color-border);
  margin-left: 0.25rem;
  background: var(--color-background-mute);
  font-size: 0.78rem;
  opacity: 0.88;
}

/* ── Data row ─────────────────────────────────────────────────────────────── */

.plan-row-card__data {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

/* ── Metrics ────────────────────────────────────────────────────────────────── */

.plan-row-card__metrics {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  flex-wrap: wrap;
  flex: 1;
}

.plan-row-card__metric {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.15rem;
  min-width: 2.5rem;
}

.plan-row-card__metric--equipment {
  min-width: 10rem;
}

.plan-row-card__metric--sum {
  margin-left: auto;
}

.plan-row-card__metric--placeholder {
  pointer-events: none;
  visibility: hidden;
}

.plan-row-card__metric-label {
  font-size: 0.6rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text);
  white-space: nowrap;
}

.plan-row-card__metric-value {
  font-weight: 600;
  color: var(--color-heading);
}

.plan-row-card__metric-value--intensity {
  color: color-mix(var(--color-primary) 50%, var(--color-text) 50%);
}

.plan-row-card__metric-value--sum {
  font-weight: 700;
}

/* ── Inputs ─────────────────────────────────────────────────────────────────── */

.plan-row-card__input {
  border: 1px solid var(--color-shadow);
  border-radius: 6px;
  background: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
  padding: 0.15rem 0.25rem;
}

.plan-row-card__input--small {
  width: 3.5rem;
  text-align: center;
}

.plan-row-card__multiselect {
  min-width: 14rem;
  width: 100%;
  max-width: 18rem;
  font-size: 0.9rem;
}

.plan-row-card__input:focus {
  outline: 1px solid var(--color-shadow);
  border: 1px solid var(--color-primary);
}

.plan-row-card__multiselect:deep(.multiselect__tags) {
  border: 1px solid var(--color-shadow);
  border-radius: 6px;
  background: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
  min-height: 2.25rem;
}

.plan-row-card__multiselect:deep(.multiselect__content-wrapper) {
  border-color: var(--color-shadow);
  background: var(--color-background);
}

.plan-row-card__multiselect:deep(.multiselect__input),
.plan-row-card__multiselect:deep(.multiselect__single) {
  background: var(--color-background);
  color: var(--color-text);
  font-size: 0.9rem;
  margin-bottom: 0;
}

.plan-row-card__multiselect:deep(.multiselect__placeholder) {
  color: var(--color-heading);
  opacity: 0.6;
  font-size: 0.75rem;
}

.plan-row-card__multiselect:deep(.multiselect__tag) {
  background: var(--color-primary);
  color: white;
  border-radius: 4px;
  font-size: 0.85rem;
  margin-bottom: 0.2rem;
}

.plan-row-card__multiselect:deep(.multiselect__tag-icon::after) {
  color: white;
}

.plan-row-card__multiselect:deep(.multiselect__tag-icon:hover) {
  background: var(--color-primary-hover, color-mix(in srgb, var(--color-primary) 80%, black));
}

.plan-row-card__multiselect:deep(.multiselect__option--highlight) {
  background: var(--color-primary);
  color: white;
}

.plan-row-card__multiselect:deep(.multiselect__option--selected) {
  background: var(--color-background-mute);
  color: var(--color-heading);
  font-weight: 600;
}

.plan-row-card__multiselect:deep(.multiselect__select::before) {
  border-color: var(--color-text) transparent transparent;
}

.plan-row-card__multiselect:deep(.multiselect--active .multiselect__tags) {
  outline: 1px solid var(--color-shadow);
  border-color: var(--color-primary);
}

/* ── Action controls (always visible in edit mode — NOT hover-only) ──────────── */

.plan-row-card__actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  flex-shrink: 0;
}

.plan-row-card__action-btn {
  width: 1.75rem;
  height: 1.75rem;
  border: none;
  background-color: var(--color-primary);
  cursor: pointer;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.15s;
  position: relative;
  flex-shrink: 0;
}

.plan-row-card__action-btn:hover:not(:disabled) {
  background-color: var(--color-primary-hover, color-mix(in srgb, var(--color-primary) 80%, black));
}

.plan-row-card__action-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* ── Action icons via pseudo-elements (same icon shapes as BaseTableAction) ─── */

/* Up arrow */
.plan-row-card__action-icon--up::before {
  content: '';
  display: block;
  width: 0;
  height: 0;
  border-left: 0.35rem solid transparent;
  border-right: 0.35rem solid transparent;
  border-bottom: 0.45rem solid white;
}

/* Down arrow */
.plan-row-card__action-icon--down::before {
  content: '';
  display: block;
  width: 0;
  height: 0;
  border-left: 0.35rem solid transparent;
  border-right: 0.35rem solid transparent;
  border-top: 0.45rem solid white;
}

/* Plus icon (add row) */
.plan-row-card__action-icon--add {
  position: relative;
  width: 0.75rem;
  height: 0.75rem;
}

.plan-row-card__action-icon--add::before,
.plan-row-card__action-icon--add::after {
  content: '';
  position: absolute;
  background: white;
}

.plan-row-card__action-icon--add::before {
  top: 50%;
  left: 0;
  width: 100%;
  height: 0.125rem;
  transform: translateY(-50%);
}

.plan-row-card__action-icon--add::after {
  left: 50%;
  top: 0;
  height: 100%;
  width: 0.125rem;
  transform: translateX(-50%);
}

/* Minus icon (remove row) */
.plan-row-card__action-icon--remove {
  position: relative;
  width: 0.75rem;
  height: 0.75rem;
}

.plan-row-card__action-icon--remove::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  width: 100%;
  height: 0.125rem;
  transform: translateY(-50%);
  background: white;
}

/* Subrow icon (L-bracket) */
.plan-row-card__action-icon--subrow {
  position: relative;
  width: 0.75rem;
  height: 0.75rem;
}

.plan-row-card__action-icon--subrow::before {
  content: '';
  position: absolute;
  top: 10%;
  left: 25%;
  width: 0.125rem;
  height: 55%;
  background: white;
}

.plan-row-card__action-icon--subrow::after {
  content: '';
  position: absolute;
  top: 65%;
  left: 25%;
  width: 55%;
  height: 0.125rem;
  background: white;
  transform: translateY(-50%);
}

/* ── Content area ───────────────────────────────────────────────────────────── */

.plan-row-card__content-view {
  color: var(--color-heading);
  font-size: 1rem;
  line-height: 1.5;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.plan-row-card__textarea {
  width: 100%;
  min-height: 3rem;
  padding: 0.3rem 0.4rem;
  border: 1px solid var(--color-shadow);
  border-radius: 8px;
  background: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
  resize: vertical;
}

.plan-row-card__textarea:focus {
  outline: 1px solid var(--color-shadow);
  border: 1px solid var(--color-primary);
}

/* ── Equipment badges ───────────────────────────────────────────────────────── */

.plan-row-card__equipment-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
  align-items: center;
  flex-shrink: 0;
}

.plan-row-card__equipment-badge {
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

/* ── Nested SubRows container ───────────────────────────────────────────────── */

.plan-row-card__subrows {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  padding-left: 0.75rem;
  margin-top: 0.25rem;
  position: relative;
}

/* List Transitions for nested cards */
.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.list-leave-active {
  position: absolute;
  left: 0;
}

/* ── Inline add-subrow button ───────────────────────────────────────────────── */

.plan-row-card__add-subrow-inline {
  display: block;
  width: 100%;
  padding: 0.3rem 0.5rem;
  border: none;
  border-top: 1px dashed var(--color-border);
  background: transparent;
  color: var(--color-primary);
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  text-align: left;
  border-radius: 0 0 6px 6px;
  transition: background-color 0.15s;
}

.plan-row-card__add-subrow-inline:hover {
  background-color: var(--color-background-soft);
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */

@media (max-width: 740px) {
  .plan-row-card {
    padding: 0.5rem 0.6rem;
    gap: 0.35rem;
  }

  .plan-row-card--depth-1 {
    margin-left: 0.5rem;
  }

  .plan-row-card--depth-2 {
    margin-left: 0.35rem;
  }

  .plan-row-card--depth-3 {
    margin-left: 0.25rem;
  }

  .plan-row-card--depth-4 {
    margin-left: 0.15rem;
  }

  .plan-row-card__metrics {
    gap: 0.4rem;
  }

  .plan-row-card__metric {
    min-width: 2rem;
  }

  .plan-row-card__metric-label {
    font-size: 0.5rem;
  }

  .plan-row-card__metric-value {
    font-size: 0.82rem;
  }

  .plan-row-card__input--small {
    width: 2.75rem;
    font-size: 0.8rem;
  }

  .plan-row-card__multiselect {
    min-width: 100%;
    max-width: none;
  }

  /* Ensure touch-friendly button targets (min 44px) */
  .plan-row-card__action-btn {
    width: 2.2rem;
    height: 2.2rem;
    min-width: 44px;
    min-height: 44px;
  }

  .plan-row-card__actions {
    gap: 0.15rem;
  }

  .plan-row-card__subrows {
    padding-left: 0.4rem;
  }

  .plan-row-card__equipment-badge {
    font-size: 0.58rem;
    padding: 0.1rem 0.3rem;
  }

  .plan-row-card__equipment-badges::before {
    font-size: 0.55rem;
  }

  /* On narrow screens, push equipment badges to their own line below metrics/actions */
  .plan-row-card__equipment-badges {
    order: 99;
    flex-basis: 100%;
    margin-top: 0.2rem;
  }
}

@media (max-width: 480px) {
  .plan-row-card__data {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
