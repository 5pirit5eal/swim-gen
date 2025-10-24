<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useExportStore } from '@/stores/export'
import type { PlanToPDFRequest, Row } from '@/types'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import BaseTableAction from '@/components/ui/BaseTableAction.vue'
import { useI18n } from 'vue-i18n'

const trainingStore = useTrainingPlanStore()
const exportStore = useExportStore()
const { t } = useI18n()

// Ref to track editing state
const isEditing = ref(false)
const editingCell = ref<{ rowIndex: number; field: keyof Row } | null>(null)
const exportPhase = ref<'idle' | 'exporting' | 'done'>('idle')
const pdfUrl = ref<string | null>(null)

const exerciseRows = computed(() => {
  if (!trainingStore.currentPlan?.table) return []
  // All rows except the last one (which should be the total)
  return trainingStore.currentPlan.table.slice(0, -1)
})

const totalRow = computed(() => {
  if (!trainingStore.currentPlan?.table) return null
  // The last row should be the total
  const table = trainingStore.currentPlan.table
  return table.length > 0 ? table[table.length - 1] : null
})

// Total exercises count (excluding the total row)
const totalExercises = computed(() => exerciseRows.value.length)

// Reset export if plan title changes (new plan)
watch(
  () => trainingStore.currentPlan?.title,
  () => {
    resetExportState()
  },
)

// Start editing a specific cell
function startEditing(rowIndex: number, field: keyof Row) {
  if (isEditing.value) {
    editingCell.value = { rowIndex, field }
  }
}

// Utility to reset export state (re-used)
function resetExportState() {
  pdfUrl.value = null
  exportPhase.value = 'idle'
}

// Toggle editing and always clear any previously generated PDF URL
function toggleEditing() {
  isEditing.value = !isEditing.value
  resetExportState()
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
      const originalRow = trainingStore.currentPlan?.table[rowIndex]
      newValue = originalRow ? originalRow[field] : 0
    }
  }
  trainingStore.updatePlanRow(rowIndex, field, newValue)
  editingCell.value = null
}

// Add a new row after the specified index
function handleAddRow(index: number) {
  trainingStore.addRow(index)
}

function handleRemoveRow(index: number) {
  trainingStore.removeRow(index)
}


async function handleExport() {
  // Phase 2: user clicks "Open PDF"
  if (exportPhase.value === 'done' && pdfUrl.value) {
    const w = window.open(pdfUrl.value, '_blank')
    if (!w) window.location.href = pdfUrl.value
    return
  }

  // Prevent double starts
  if (exportPhase.value === 'exporting') return
  if (!trainingStore.currentPlan) return

  // Phase 1: user clicks "Export PDF"
  exportPhase.value = 'exporting'
  try {
    pdfUrl.value = await exportStore.exportToPDF(trainingStore.currentPlan as PlanToPDFRequest)
    if (!pdfUrl.value) {
      exportPhase.value = 'idle'
      return
    }
    exportPhase.value = 'done'
  } catch (e) {
    console.error('PDF export failed', e)
    exportPhase.value = 'idle'
  }
}
</script>

<template>
  <div class="training-plan-display">
    <div v-if="trainingStore.isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('display.generating_plan_message') }}</p>
    </div>
    <div v-else-if="trainingStore.hasPlan && trainingStore.currentPlan" class="plan-container">
      <!-- Header -->
      <header class="plan-header">
        <h2 class="plan-title">{{ trainingStore.currentPlan.title }}</h2>
        <div class="plan-description">
          {{ trainingStore.currentPlan.description }}
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
          <tbody>
            <template v-for="(row, index) in exerciseRows" :key="index">
              <tr class="exercise-row">
                <!-- Amount Cell -->
                <td @click="startEditing(index, 'Amount')" class="anchor-cell">
                  <BaseTableAction v-if="isEditing">
                    <template #actions>
                      <div class="action-buttons">
                        <button @click.stop="handleAddRow(index)" class="action-btn add-btn"
                          :title="t('display.add_row')">
                          <!-- CSS Icon -->
                        </button>
                        <button @click.stop="handleRemoveRow(index)" class="action-btn remove-btn"
                          :title="t('display.remove_row')">
                          <!-- CSS Icon -->
                        </button>
                      </div>
                    </template>
                  </BaseTableAction>
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
                    @keyup.enter="stopEditing($event, index, 'Content')" class="editable-area"></textarea>
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
                <strong>{{ totalRow.Content }}</strong>
              </td>
              <td>
                <strong>{{ totalRow.Sum }} m</strong>
              </td>
            </tr>
          </tbody>
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

  <div v-if="trainingStore.hasPlan && trainingStore.currentPlan && !trainingStore.isLoading" class="export-section">
    <!-- Edit Action -->
    <button @click="toggleEditing" class="export-btn">
      {{ isEditing ? t('display.done_editing') : t('display.refine_plan') }}
    </button>
    <!-- Export Action -->
    <button @click="handleExport" class="export-btn" :disabled="exportPhase === 'exporting'">
      <template v-if="exportPhase === 'exporting'">
        {{ t('display.exporting') }}
      </template>
      <template v-else-if="exportPhase === 'done'">
        {{ t('display.open_pdf') }}
      </template>
      <template v-else>
        {{ t('display.export_pdf') }}
      </template>
    </button>
  </div>
</template>

<style scoped>
.training-plan-display {
  margin: 2rem auto;
  background: var(--color-background-soft);
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
}

.plan-container {
  background: var(--color-background);
  border-radius: 0.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
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
  border-top-right-radius: 0.5rem;
  border-top-left-radius: 0.5rem;
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
  color: var(--color-text-light);
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
  color: var(--color-text);
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
  border: 1px solid var(--color-primary);
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
  border: 1px solid var(--color-primary);
  border-radius: 0.25rem;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
}

.anchor-cell {
  position: relative;
  border-left: none;
}

/* Show action container on row hover */
.exercise-row:hover .anchor-cell :deep(.action-container) {
  opacity: 1;
}

.action-buttons {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  /* Space between buttons */
}

.action-btn {
  background-color: var(--color-primary);
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  position: relative;
  /* Needed for pseudo-element positioning */
  transition: background-color 0.2s;
}

/* Common style for icon bars */
.action-btn::before,
.action-btn::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: white;
  width: 12px;
  /* Width of the icon bar */
  height: 2px;
  /* Thickness of the icon bar */
  border-radius: 1px;
}

/* Create the vertical bar for the plus icon */
.add-btn::after {
  transform: translate(-50%, -50%) rotate(90deg);
}

.action-btn:hover {
  background-color: var(--color-primary-hover);
}


.total-cell {
  font-weight: 600;
}

.total-row {
  background: var(--color-border) !important;
  color: var(--color-text) !important;
  font-weight: 700;
  font-size: 1rem;
}

.total-row td {
  border-color: var(--color-border);
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
  color: var(--color-text-light);
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 1px;
}

.loading-state,
.no-plan {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--color-text-light);
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

.export-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
  padding: 1.5rem;
  background: var(--color-background-soft);
  text-align: center;
}

.export-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 160px;
}

@media (max-width: 740px) {
  .export-btn {
    width: 100%;
    min-width: 10%;
    padding: 0.5rem 1rem;
  }

  button+button {
    margin-left: 2rem;
  }
}

.export-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.export-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: var(--color-text-light);
}
</style>
