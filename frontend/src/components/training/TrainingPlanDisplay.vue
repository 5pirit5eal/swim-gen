<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useExportStore } from '@/stores/export'
import type { PlanToPDFRequest, Row } from '@/types'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { useI18n } from 'vue-i18n'

const trainingStore = useTrainingPlanStore()
const exportStore = useExportStore()
const { t } = useI18n()

// Ref to track editing state
const isEditing = ref(false)
const editingCell = ref<{ rowIndex: number; field: keyof Row } | null>(null)

// Computed for separating exercise rows from total row
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
const totalExercises = computed(() => {
  return exerciseRows.value.length
})

// Start editing a specific cell
function startEditing(rowIndex: number, field: keyof Row) {
  if (isEditing.value) {
    editingCell.value = { rowIndex, field }
  }
}

// Stop editing the current cell and save the changes
function stopEditing(event: Event, rowIndex: number, field: keyof Row) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  let newValue: string | number = target.value

  if (['Amount', 'Distance'].includes(field as string)) {
    // Convert numeric fields to numbers
    const numValue = parseFloat(newValue as string)
    newValue = isNaN(numValue) ? 0 : Math.max(0, numValue)
  }
  trainingStore.updatePlanRow(rowIndex, field, newValue)
  editingCell.value = null
}

async function handleExport() {
  if (!trainingStore.currentPlan) return

  const pdfUri = await exportStore.exportToPDF(trainingStore.currentPlan as PlanToPDFRequest)
  if (pdfUri) {
    // Trigger download
    window.open(pdfUri, '_blank')
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
                  <template #tooltip>
                    {{ t('display.amount_tooltip') }}
                  </template>
                </BaseTooltip>
              </th>
              <th></th>
              <th>
                {{ t('display.distance') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('display.distance_tooltip') }}
                  </template>
                </BaseTooltip>
              </th>
              <th>
                {{ t('display.break') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('display.break_tooltip') }}
                  </template>
                </BaseTooltip>
              </th>
              <th>
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
                  <template #tooltip>
                    {{ t('display.total_tooltip') }}
                  </template>
                </BaseTooltip>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in exerciseRows" :key="index" class="exercise-row">
              <!-- Amount Cell -->
              <td @click="startEditing(index, 'Amount')">
                <input type="number" min="0" v-if="isEditing" :value="row.Amount"
                  @blur="stopEditing($event, index, 'Amount')" @keyup.enter="stopEditing($event, index, 'Amount')"
                  class="editable-small" />
                <span v-else>{{ row.Amount }}</span>
              </td>
              <td>{{ row.Multiplier }}</td>
              <!-- Distance Cell -->
              <td @click="startEditing(index, 'Distance')">
                <input type="number" min="0" max="100000" step="25" v-if="isEditing" :value="row.Distance"
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
                  @blur="stopEditing($event, index, 'Intensity')" @keyup.enter="stopEditing($event, index, 'Intensity')"
                  class="editable-small" />
                <span v-else>{{ row.Intensity }}</span>
              </td>
              <td class="total-cell">{{ row.Sum }}</td>
            </tr>
            <!-- Total row from backend -->
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
    <button @click="isEditing = !isEditing" class="export-btn">
      {{ isEditing ? t('display.done_editing') : t('display.refine_plan') }}
    </button>
    <!-- Export Action -->
    <button @click="handleExport" class="export-btn" :disabled="exportStore.isExporting">
      {{ exportStore.isExporting ? t('display.exporting') : t('display.export_pdf') }}
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
  overflow: visible;
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
}

.exercise-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.exercise-table th,
.exercise-table td {
  border: 1px solid var(--color-border);
  padding: 0.75rem 0.5rem;
  text-align: center;
  color: var(--color-text-light);
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
  white-space: nowrap;
  /* Prevent text from wrapping */
}

.exercise-table td>span,
.exercise-table td>textarea {
  display: block;
}

.exercise-row:nth-child(even) {
  background-color: var(--color-background);
}

.exercise-row:nth-child(odd) {
  background-color: var(--color-background-soft);
}

.exercise-row:hover {
  background-color: var(--color-background-mute);
}

.content-cell {
  text-align: left;
  font-weight: 500;
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
  width: 70px;
  text-align: center;
  border: 1px solid var(--color-primary);
  border-radius: 0.25rem;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
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

.export-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.export-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: var(--color-text-light);
}
</style>
