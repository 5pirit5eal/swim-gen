<script setup lang="ts">
import { computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'

const trainingStore = useTrainingPlanStore()

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
</script>

<template>
  <div class="training-plan-display">
    <div v-if="trainingStore.hasPlan && trainingStore.currentPlan" class="plan-container">
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
    </div>

    <div v-else-if="trainingStore.isGenerating" class="loading-state">
      <p>Generating your training plan...</p>
    </div>

    <div v-else class="no-plan">
      <p>No training plan generated yet. Use the form above to create one!</p>
    </div>
  </div>
</template>

<style scoped>
.training-plan-display {
  max-width: 800px;
  margin: 2rem auto;
}

.plan-container {
  background: var(--color-background);
  border-radius: 0.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.plan-header {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
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
}

.exercise-table th {
  background: linear-gradient(135deg, #374151, #1f2937);
  color: white;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-size: 0.8rem;
}

.exercise-row:nth-child(even) {
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
  background: linear-gradient(135deg, #1f2937, #374151) !important;
  color: white;
  font-weight: 700;
  font-size: 1rem;
}

.total-row td {
  border-color: #374151;
}

.summary-section {
  display: flex;
  justify-content: space-around;
  padding: 1.5rem;
  background: var(--color-background-soft);
  gap: 1rem;
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

/* Mobile responsiveness */
@media (max-width: 768px) {
  .plan-header {
    padding: 1.5rem 1rem;
  }

  .plan-title {
    font-size: 1.25rem;
  }

  .table-container {
    padding: 1rem;
  }

  .exercise-table {
    font-size: 0.8rem;
  }

  .exercise-table th,
  .exercise-table td {
    padding: 0.5rem 0.25rem;
  }

  .summary-section {
    flex-direction: column;
    padding: 1rem;
  }
}
</style>
