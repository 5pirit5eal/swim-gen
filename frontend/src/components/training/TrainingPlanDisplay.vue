<script setup lang="ts">
import ExportPlanButton from '@/components/buttons/ButtonExportPlan.vue'
import SharePlanButton from '@/components/buttons/ButtonSharePlan.vue'
import IconEdit from '@/components/icons/IconEdit.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import BaseTableAction from '@/components/ui/BaseTableAction.vue'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import type { Row, PlanStore, RAGResponse } from '@/types'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    store: PlanStore
    showShareButton?: boolean
    planOverride?: RAGResponse | null
  }>(),
  {
    showShareButton: false,
    planOverride: undefined,
  },
)

const { t } = useI18n()

// Ref to track editing state
const isEditing = ref(false)
const editingCell = ref<{ rowIndex: number; field: keyof Row } | null>(null)

const exerciseRows = computed(() => {
  const plan = props.planOverride || props.store.currentPlan
  if (!plan?.table) return []
  // All rows except the last one (which should be the total)
  return plan.table.slice(0, -1)
})

const totalRow = computed(() => {
  const plan = props.planOverride || props.store.currentPlan
  if (!plan?.table) return null
  // The last row should be the total
  const table = plan.table
  return table.length > 0 ? table[table.length - 1] : null
})

// Total exercises count (excluding the total row)
const totalExercises = computed(() => exerciseRows.value.length)

// Start editing a specific cell
function startEditing(rowIndex: number, field: keyof Row) {
  if (isEditing.value) {
    editingCell.value = { rowIndex, field }
  }
}

// Toggle editing
async function toggleEditing() {
  isEditing.value = !isEditing.value
  if (!isEditing.value) {
    // Upsert the current plan when done editing
    await props.store.upsertCurrentPlan()
  }
}

// Stop editing the current cell and save the changes
function stopEditing(event: Event, rowIndex: number, field: keyof Row) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  let newValue: string | number = target.value

  if (['Amount', 'Distance'].includes(field as string)) {
    // Convert numeric fields to numbers, ensuring it's a valid number
    const numValue = parseFloat(newValue as string)
    if (!isNaN(numValue) && /^\d*\.?\d*$/.test(newValue as string)) {
      newValue = Math.max(0, numValue)
      newValue = Math.round(newValue as number)
    } else {
      // Revert to the original value if input is invalid
      const originalRow = props.store.currentPlan?.table[rowIndex]
      const val = originalRow ? originalRow[field] : 0
      newValue = val !== undefined ? val : 0
    }
  }
  props.store.updatePlanRow(rowIndex, field, newValue)
  editingCell.value = null
}

// Add a new row after the specified index
function handleAddRow(index: number) {
  props.store.addRow(index)
}

function handleRemoveRow(index: number) {
  props.store.removeRow(index)
}

function handleMoveRow(index: number, direction: 'up' | 'down') {
  props.store.moveRow(index, direction)
}

// Auto-resize directive for textarea
const vAutoResize = {
  mounted: (el: HTMLTextAreaElement) => {
    el.style.height = 'auto'
    el.style.height = el.scrollHeight + 'px'
    el.style.overflowY = 'hidden'
  }
}

function autoResize(event: Event) {
  const target = event.target as HTMLTextAreaElement
  target.style.height = 'auto'
  target.style.height = target.scrollHeight + 'px'
}
</script>

<template>
  <div class="training-plan-display">
    <div v-if="store.isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('display.generating_plan_message') }}</p>
    </div>
    <div v-else-if="store.hasPlan && (store.currentPlan || planOverride)" class="plan-container">
      <!-- Header -->
      <header class="plan-header">
        <div v-if="isEditing" class="edit-header">
          <input v-model="store.currentPlan!.title" class="edit-title" v-auto-resize
            :placeholder="t('display.plan_title')" />
          <textarea v-model="store.currentPlan!.description" v-auto-resize class="edit-description"
            :placeholder="t('display.plan_description')" rows="3"></textarea>
        </div>
        <div v-else>
          <h2 class="plan-title">{{ (planOverride || store.currentPlan)?.title }}</h2>
          <div class="plan-description">
            {{ (planOverride || store.currentPlan)?.description }}
          </div>
        </div>
      </header>

      <!-- Exercise Table -->
      <div class="table-container">
        <table class="exercise-table">
          <thead>
            <tr>
              <th>
                {{ t('display.amount') }}
                <BaseTooltip>
                  <template #tooltip>{{ t('display.amount_tooltip') }}</template>
                </BaseTooltip>
              </th>
              <th class="multiplier"></th>
              <th>
                {{ t('display.distance') }}
                <BaseTooltip>
                  <template #tooltip>{{ t('display.distance_tooltip') }}</template>
                </BaseTooltip>
              </th>
              <th>
                {{ t('display.break') }}
                <BaseTooltip>
                  <template #tooltip>{{ t('display.break_tooltip') }}</template>
                </BaseTooltip>
              </th>
              <th class="content-header">
                {{ t('display.content') }}
                <BaseTooltip>
                  <template #tooltip>
                    <p>{{ t('display.content_tooltip.title') }}</p>
                    <ul>
                      <li>
                        <strong>{{ t('display.content_tooltip.freestyle') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.content_tooltip.backstroke') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.content_tooltip.breaststroke') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.content_tooltip.leg_work') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.content_tooltip.butterfly') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.content_tooltip.individual_medley') }}</strong>
                      </li>
                    </ul>
                  </template>
                </BaseTooltip>
              </th>
              <th>
                {{ t('display.intensity') }}
                <BaseTooltip>
                  <template #tooltip>
                    <p>{{ t('display.intensity_tooltip.title') }}</p>
                    <ul>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.ga') }}</strong>
                        <ul>
                          <li>
                            <strong>{{ t('display.intensity_tooltip.ga1') }}</strong>
                          </li>
                          <li>
                            <strong>{{ t('display.intensity_tooltip.ga1_2') }}</strong>
                          </li>
                          <li>
                            <strong>{{ t('display.intensity_tooltip.ga2') }}</strong>
                          </li>
                        </ul>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.lza') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.hf') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.lt') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.sa') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.ta') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.tue') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.ts') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.sprint') }}</strong>
                      </li>
                      <li>
                        <strong>{{ t('display.intensity_tooltip.recovery') }}</strong>
                      </li>
                    </ul>
                  </template>
                </BaseTooltip>
              </th>
              <th>
                {{ t('display.total') }}
                <BaseTooltip>
                  <template #tooltip>{{ t('display.total_tooltip') }}</template>
                </BaseTooltip>
              </th>
            </tr>
          </thead>
          <TransitionGroup tag="tbody" name="list">
            <template v-for="(row, index) in exerciseRows" :key="row._id || index">
              <tr class="exercise-row">
                <!-- Amount Cell -->
                <td @click="startEditing(index, 'Amount')" class="anchor-cell">
                  <BaseTableAction v-if="isEditing" :is-first="index === 0" :is-last="index === exerciseRows.length - 1"
                    @add="handleAddRow(index)" @remove="handleRemoveRow(index)" @move-up="handleMoveRow(index, 'up')"
                    @move-down="handleMoveRow(index, 'down')" />
                  <input type="text" inputmode="numeric" pattern="[0-9]*" v-if="isEditing" :value="row.Amount"
                    @blur="stopEditing($event, index, 'Amount')" @keyup.enter="stopEditing($event, index, 'Amount')"
                    class="editable-small" />
                  <span v-else>{{ row.Amount }}</span>
                </td>
                <td>{{ row.Multiplier }}</td>
                <!-- Distance Cell -->
                <td @click="startEditing(index, 'Distance')">
                  <input type="text" inputmode="numeric" pattern="[0-9]*" v-if="isEditing" :value="row.Distance"
                    @blur="stopEditing($event, index, 'Distance')" @keyup.enter="stopEditing($event, index, 'Distance')"
                    class="editable-small" />
                  <span v-else>{{ row.Distance }}</span>
                </td>
                <!-- Intensity Cell -->
                <td @click="startEditing(index, 'Break')">
                  <input type="text" v-if="isEditing" :value="row.Break" @blur="stopEditing($event, index, 'Break')"
                    @keyup.enter="stopEditing($event, index, 'Break')" class="editable-small" />
                  <span v-else>{{ row.Break }}</span>
                </td>
                <!-- Content Cell -->
                <td class="content-cell" @click="startEditing(index, 'Content')">
                  <textarea v-if="isEditing" :value="row.Content" @blur="stopEditing($event, index, 'Content')"
                    @keyup.enter="stopEditing($event, index, 'Content')" @input="autoResize" v-auto-resize
                    class="editable-area"></textarea>
                  <span v-else>{{ row.Content }}</span>
                </td>
                <!-- Intensity Cell -->
                <td class="intensity-cell" @click="startEditing(index, 'Intensity')">
                  <input type="text" v-if="isEditing" :value="row.Intensity"
                    @blur="stopEditing($event, index, 'Intensity')"
                    @keyup.enter="stopEditing($event, index, 'Intensity')" class="editable-small" />
                  <span v-else>{{ row.Intensity }}</span>
                </td>
                <td class="total-cell">{{ row.Sum }}</td>
              </tr>
            </template>
            <!-- Total row -->
            <tr v-if="totalRow" class="total-row">
              <td colspan="6">
                <strong>{{ t('display.meters_total') }}</strong>
              </td>
              <td>
                <strong>{{ totalRow.Sum }} m</strong>
              </td>
            </tr>
          </TransitionGroup>
        </table>
      </div>

      <!-- Summary Statistics -->
      <div class="summary-section">
        <div class="summary-item">
          <div class="summary-value">{{ totalRow?.Sum || 0 }}</div>
          <div class="summary-label">{{ t('display.meters_total') }}</div>
        </div>
        <div class="summary-item">
          <div class="summary-value">{{ totalExercises }}</div>
          <div class="summary-label">{{ t('display.exercise_sets') }}</div>
        </div>
      </div>
    </div>

    <div v-else class="no-plan">
      <p>{{ t('display.no_plan_placeholder') }}</p>
    </div>
  </div>

  <div v-if="store.hasPlan && store.currentPlan && !store.isLoading" class="button-section">
    <!-- Edit Action -->
    <button @click="toggleEditing" class="edit-btn">
      <IconCheck v-if="isEditing" class="icon" />
      <IconEdit v-else class="icon" />
      {{ isEditing ? t('display.done_editing') : t('display.refine_plan') }}
    </button>
    <div class="gap"></div>
    <SharePlanButton v-if="showShareButton" :store="store" />
    <div class="gap"></div>
    <ExportPlanButton :store="store" />
  </div>
</template>

<style scoped>
.training-plan-display {
  background: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.plan-container {
  background: var(--color-background);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
}

@media (max-width: 740px) {
  .training-plan-display {
    zoom: 0.75;
  }
}

.plan-header {
  background: var(--color-primary);
  color: white;
  padding: 2rem;
  text-align: center;
  border-top-right-radius: 8px;
  border-top-left-radius: 8px;
  outline: 1px solid var(--color-primary);
  border: 2px solid var(--color-primary);
}

.plan-title {
  margin: 0 0 1rem 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.plan-description {
  font-size: 1rem;
  line-height: 1.6;
  opacity: 0.95;
}

.edit-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.edit-title {
  font-size: 1.5rem;
  font-weight: 700;
  padding: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  text-align: center;
}

.edit-title::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.edit-description {
  font-size: 1rem;
  line-height: 1.6;
  padding: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-family: inherit;
  resize: vertical;
}

.edit-description::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.table-container {
  padding: 1.5rem;
  background: var(--color-background-soft);
  width: inherit;
  /* Set table to take full width of its container */
  table-layout: fixed;
}

.exercise-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
  table-layout: fixed;
}

.exercise-table th,
.exercise-table td {
  border: 1px solid var(--color-border);
  padding: 0.75rem 0.5rem;
  text-align: center;
  color: var(--color-text);
  width: auto;
}

.exercise-table th.multiplier,
.exercise-table td.multiplier {
  width: 5%;
}

.exercise-table th.content-header {
  width: 30%;
}

.exercise-table th {
  background: var(--color-border);
  color: var(--color-heading);
  /* Use flexbox for alignment */
  align-items: center;
  /* Horizontally center items */
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-size: 0.8rem;
  /* white-space: nowrap; */
  word-break: break-all;
  /* Prevent text from wrapping */
  padding: 0.5rem 0.25rem;
}

@media (max-width: 740px) {
  .exercise-table th {
    font-size: 0.5rem;
    padding: 0.5rem 0.25rem;
    white-space: normal;
    word-break: break-all;
  }

  .exercise-table td {
    padding: 0.5rem 0.25rem;
    white-space: normal;
    padding: 0.25rem 0.2rem;
    font-size: 0.75rem;
  }
}

.exercise-table td>span,
.exercise-table td>textarea {
  display: block;
}

/* Apply alternating backgrounds to data cells */
.exercise-row:nth-child(even) {
  background-color: var(--color-background);
}

.exercise-row:nth-child(odd) {
  background-color: var(--color-background-soft);
}

/* Apply hover effect to data cells */
.exercise-row:hover {
  background-color: var(--color-background-mute);
}

.exercise-row:hover {
  --action-bg-color: var(--color-background-mute);
}

.content-cell {
  text-align: left;
  font-weight: 500;
  width: 300px;
}

.intensity-cell {
  font-weight: 600;
  color: var(--color-primary);
}

.editable-area {
  width: 100%;
  padding: 0.25rem;
  border: 1px solid var(--color-shadow);
  border-radius: 0.25rem;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
  /* Include padding and border in the element's total width and height */
}

.editable-small {
  width: 70%;
  text-align: center;
  border: 1px solid var(--color-shadow);
  border-radius: 0.25rem;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
}

.editable-area:focus,
.editable-small:focus {
  outline: 2px solid var(--color-primary);
}

.anchor-cell {
  position: relative;
  border-left: none;
}

/* Show action container on row hover */
.exercise-row:hover .anchor-cell :deep(.action-container) {
  opacity: 1;
  transform: translateX(0);
}

.total-cell {
  font-weight: 600;
}

.total-row {
  background: var(--color-border);
  font-weight: 700;
  font-size: 1rem;
}

.total-row td {
  border-color: var(--color-border);
  color: var(--color-heading);
}

.summary-section {
  display: flex;
  justify-content: space-around;
  padding-bottom: 1rem;
  padding-top: 0rem;
  padding-left: 1rem;
  padding-right: 1rem;
  background: var(--color-background-soft);
  gap: 3rem;
  border-bottom-right-radius: 0.5rem;
  border-bottom-left-radius: 0.5rem;
}

.summary-item {
  background: var(--color-background);
  padding: 1rem;
  border-radius: 0.5rem;
  text-align: center;
  flex: 1;
  border: 1px solid var(--color-border);
}

.summary-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 0.25rem;
}

.summary-label {
  color: var(--color-heading);
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 1px;
}

.loading-state,
.no-plan {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--color-text);
  font-style: italic;
}

.loading-spinner {
  width: 120px;
  height: 40px;
  background-color: var(--color-background-soft);
  position: relative;
  border-radius: 50px;
  box-shadow: inset 0 0 0 2px var(--color-border);
  margin: 0 auto 1rem auto;
}

.loading-spinner:after {
  border-radius: 50px;
  content: '';
  position: absolute;
  background-color: var(--color-primary);
  left: 2px;
  top: 2px;
  bottom: 2px;
  right: 80px;
  animation: slide 2s linear infinite;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

@keyframes slide {
  0% {
    right: 80px;
    left: 2px;
  }

  5% {
    left: 2px;
  }

  50% {
    right: 2px;
    left: 80px;
  }

  55% {
    right: 2px;
  }

  100% {
    right: 80px;
    left: 2px;
  }
}

.button-section {
  display: flex;
  justify-content: space-between;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
  padding: 1.5rem;
  background: var(--color-background-soft);
  text-align: center;
  margin-top: 1rem;
}

.edit-btn {
  flex: 1;
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

@media (max-width: 740px) {
  .edit-btn {
    /* width: 100%;
    min-width: 10%; */
    padding: 0.5rem 1rem;
    overflow-wrap: break-word;
    font-size: 0.8rem;
  }
}

.edit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.icon {
  width: 24px;
  height: 24px;
}

.gap {
  flex: 2;
  display: flex;
}

/* List Transitions */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* Ensure the leaving item is taken out of flow so others can move */
/* Note: position: absolute on table rows can be tricky, but often needed for smooth 'move' during 'leave' */
/* If this breaks table layout during animation, remove the absolute positioning */
/* .list-leave-active { position: absolute; } */
</style>
