<script setup lang="ts">
import type { Row, RAGResponse, PlanStore } from '@/types'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import PlanRowCard from '@/components/training/PlanRowCard.vue'

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

// Minimal read-only store-compatible object — PlanRowCard requires a store prop,
// but in read-only mode (isEditing=false) it never calls any mutating methods.
const readonlyStore = {
  currentPlan: null,
  hasPlan: false,
  isLoading: false,
  keepForever: () => Promise.resolve(),
  upsertCurrentPlan: () => Promise.resolve(''),
  updatePlanRow: () => {},
  updatePlanRowEquipment: () => {},
  addRow: () => {},
  addSubRow: () => {},
  removeRow: () => {},
  moveRow: () => {},
} as unknown as PlanStore

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
    </header>

    <!-- Exercise Cards -->
    <div class="plan-cards-list">
      <PlanRowCard
        v-for="(row, index) in exerciseRows"
        :key="row._id || index"
        :row="row"
        :path="[index]"
        :depth="0"
        :is-editing="false"
        :store="readonlyStore"
        :is-first="index === 0"
        :is-last="index === exerciseRows.length - 1"
      />
    </div>

    <!-- Summary and Actions -->
    <div class="footer-section">
      <div class="summary-compact" data-testid="plan-summary">
        <span class="summary-item">{{ totalRow?.Sum || 0 }} m</span>
        <span class="separator">&bull;</span>
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
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--color-heading);
}

.plan-cards-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.75rem;
  background: var(--color-background-soft);
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
  font-weight: 700;
  color: var(--color-primary);
}

.separator {
  color: var(--color-border-hover);
  font-size: 0.75rem;
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
