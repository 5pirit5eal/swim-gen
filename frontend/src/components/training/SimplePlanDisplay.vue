<script setup lang="ts">
import type { Row, RAGResponse } from '@/types'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import ContentWithDrillLinks from '@/components/training/ContentWithDrillLinks.vue'

const props = defineProps<{
  title: string
  description: string
  table: Row[]
  planId: string
}>()

const emit = defineEmits<{
  save: [plan: RAGResponse]
}>()

const { t } = useI18n()

const exerciseRows = computed(() => {
  if (!props.table) return []
  return props.table.slice(0, -1)
})

const totalRow = computed(() => {
  if (!props.table || props.table.length === 0) return null
  return props.table[props.table.length - 1]
})

const totalExercises = computed(() => exerciseRows.value.length)

function onSave() {
  emit('save', {
    plan_id: props.planId,
    title: props.title,
    description: props.description,
    table: props.table,
  })
}
</script>

<template>
  <div class="simple-plan-display">
    <!-- Header -->
    <header class="plan-header">
      <h3 class="plan-title">{{ title }}</h3>
      <div class="plan-description">{{ description }}</div>
    </header>

    <!-- Compact Exercise Table -->
    <div class="table-container">
      <table class="exercise-table">
        <thead>
          <tr>
            <th>{{ t('display.amount') }}</th>
            <th>x</th>
            <th>{{ t('display.distance') }}</th>
            <th>{{ t('display.break') }}</th>
            <th class="content-header">{{ t('display.content') }}</th>
            <th>{{ t('display.intensity') }}</th>
            <th>{{ t('display.total') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, index) in exerciseRows" :key="index" class="exercise-row">
            <td>{{ row.Amount }}</td>
            <td>{{ row.Multiplier }}</td>
            <td>{{ row.Distance }}</td>
            <td>{{ row.Break }}</td>
            <td class="content-cell">
              <ContentWithDrillLinks :content="row.Content" />
            </td>
            <td class="intensity-cell">{{ row.Intensity }}</td>
            <td class="total-cell">{{ row.Sum }}</td>
          </tr>
          <tr v-if="totalRow" class="total-row">
            <td colspan="6">
              <strong>{{ t('display.meters_total') }}</strong>
            </td>
            <td>
              <strong>{{ totalRow.Sum }} m</strong>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Summary and Actions -->
    <div class="footer-section">
      <div class="summary-compact">
        <span class="summary-item">{{ totalRow?.Sum || 0 }} m</span>
        <span class="separator">•</span>
        <span class="summary-item">{{ totalExercises }} {{ t('display.exercise_sets') }}</span>
      </div>
      <button @click="onSave" class="save-btn">
        {{ t('interaction.save_to_history') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.simple-plan-display {
  background: var(--color-background);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.plan-header {
  background: var(--color-background-soft);
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.plan-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-heading);
}

.plan-description {
  font-size: 0.9rem;
  line-height: 1.4;
  color: var(--color-text);
}

.table-container {
  padding: 0.75rem;
  overflow-x: auto;
}

.exercise-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.exercise-table th,
.exercise-table td {
  border: 1px solid var(--color-border);
  padding: 0.4rem 0.3rem;
  text-align: center;
  color: var(--color-text);
}

.exercise-table th {
  background: var(--color-border);
  color: var(--color-heading);
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.content-header {
  width: 30%;
}

.content-cell {
  text-align: left;
  font-size: 0.8rem;
}

.intensity-cell {
  font-weight: 600;
  color: var(--color-primary);
}

.total-cell {
  font-weight: 600;
}

.exercise-row:nth-child(even) {
  background-color: var(--color-background-soft);
}

.total-row {
  background: var(--color-border);
  font-weight: 700;
}

.total-row td {
  border-color: var(--color-border);
  color: var(--color-heading);
}

.footer-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: var(--color-background-soft);
  border-top: 1px solid var(--color-border);
}

.summary-compact {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: var(--color-text);
}

.summary-item {
  font-weight: 500;
}

.separator {
  color: var(--color-border);
}

.save-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.save-btn:hover {
  background: var(--color-primary-hover);
}
</style>
