<script setup lang="ts">
import { computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useExportStore } from '@/stores/export'
import type { PlanToPDFRequest } from '@/types'

const trainingStore = useTrainingPlanStore()

const exportStore = useExportStore()

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
    <div v-if="trainingStore.isGenerating" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Generating your training plan...</p>
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
              <th>Amount</th>
              <th>Distance (m)</th>
              <th>Break</th>
              <th>Content</th>
              <th>Intensity</th>
              <th>Total (m)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in exerciseRows" :key="index" class="exercise-row">
              <td>{{ row.Amount }}{{ row.Multiplier }}</td>
              <td>{{ row.Distance }}</td>
              <td>{{ row.Break }}</td>
              <td class="content-cell">{{ row.Content }}</td>
              <td class="intensity-cell">{{ row.Intensity }}</td>
              <td class="total-cell">{{ row.Sum }}</td>
            </tr>
            <!-- Total row from backend -->
            <tr v-if="totalRow" class="total-row">
              <td colspan="5">
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
          <div class="summary-label">Meters Total</div>
        </div>
        <div class="summary-item">
          <div class="summary-value">{{ totalExercises }}</div>
          <div class="summary-label">Exercise Sets</div>
        </div>
      </div>

      <!-- Export Action -->
      <div class="export-section">
        <button @click="handleExport" class="export-btn" :disabled="exportStore.isExporting">
          {{ exportStore.isExporting ? 'Exporting...' : 'Export PDF' }}
        </button>
      </div>
    </div>

    <div v-else class="no-plan">
      <p>No training plan generated yet. Use the form above to create one!</p>
    </div>
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
  overflow: hidden;
}

.plan-header {
  background: var(--color-primary, #3b82f6);
  color: white;
  padding: 2rem;
  text-align: center;
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
  overflow-x: auto;
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
  color: white;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-size: 0.8rem;
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

.total-cell {
  font-weight: 600;
}

.total-row {
  background: var(--color-border) !important;
  color: white;
  font-weight: 700;
  font-size: 1rem;
}

.total-row td {
  border-color: var(--color-border);
}

.summary-section {
  display: flex;
  justify-content: space-around;
  padding: 1.5rem;
  background: var(--color-background-soft);
  gap: 3rem;
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
  background-color: var(--color-primary, #3b82f6);
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
  padding: 1.5rem;
  background: var(--color-background-soft);
  text-align: center;
}

.export-btn {
  background: var(--color-primary, #3b82f6);
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
  background: var(--color-primary-hover, #2563eb);
}

.export-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: var(--color-text-light);
}
</style>
