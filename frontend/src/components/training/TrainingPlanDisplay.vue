<script setup lang="ts">
import ExportPlanButton from '@/components/buttons/ButtonExportPlan.vue'
import SharePlanButton from '@/components/buttons/ButtonSharePlan.vue'
import IconEdit from '@/components/icons/IconEdit.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import PlanRowCard from '@/components/training/PlanRowCard.vue'
import type { PlanStore } from '@/types'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = withDefaults(
  defineProps<{
    store: PlanStore
    showShareButton?: boolean
  }>(),
  {
    showShareButton: false,
  },
)

const { t } = useI18n()

// Ref to track editing state
const isEditing = ref(false)

const exerciseRows = computed(() => {
  const plan = props.store.currentPlan
  if (!plan?.table) return []
  // All rows except the last one (which should be the total)
  return plan.table.slice(0, -1)
})

const totalRow = computed(() => {
  const plan = props.store.currentPlan
  if (!plan?.table) return null
  // The last row should be the total
  const table = plan.table
  return table.length > 0 ? table[table.length - 1] : null
})

// Total exercises count (excluding the total row)
const totalExercises = computed(() => exerciseRows.value.length)

// Toggle editing
async function toggleEditing() {
  isEditing.value = !isEditing.value
  if (!isEditing.value) {
    // Upsert the current plan when done editing
    await props.store.upsertCurrentPlan()
  }
}
</script>

<template>
  <div class="training-plan-display" id="tutorial-plan-display">
    <div v-if="store.isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('display.generating_plan_message') }}</p>
    </div>
    <div v-else-if="store.hasPlan && store.currentPlan" class="plan-container">
      <!-- Header -->
      <header class="plan-header">
        <div v-if="isEditing" class="edit-header">
          <input
            v-model="store.currentPlan!.title"
            class="edit-title"
            v-auto-resize
            :placeholder="t('display.plan_title')"
          />
          <textarea
            v-model="store.currentPlan!.description"
            v-auto-resize
            class="edit-description"
            :placeholder="t('display.plan_description')"
            rows="3"
          ></textarea>
        </div>
        <div v-else>
          <h2 class="plan-title">{{ store.currentPlan?.title }}</h2>
          <div class="plan-description">
            {{ store.currentPlan?.description }}
          </div>
        </div>
      </header>

      <!-- Exercise Cards -->
      <div class="plan-cards-list" data-testid="plan-cards-list">
        <PlanRowCard
          v-for="(row, index) in exerciseRows"
          :key="row._id || index"
          :row="row"
          :path="[index]"
          :depth="0"
          :is-editing="isEditing"
          :store="store"
          :is-first="index === 0"
          :is-last="index === exerciseRows.length - 1"
        />
      </div>

      <!-- Total row summary -->
      <div v-if="totalRow" class="total-summary-row">
        <span class="total-summary-label">{{ t('display.meters_total') }}</span>
        <span class="total-summary-value">{{ totalRow.Sum }} m</span>
      </div>

      <!-- Summary Statistics -->
      <div class="summary-section" data-testid="plan-summary">
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
    <button @click="toggleEditing" class="edit-btn" data-testid="plan-edit-btn">
      <IconCheck v-if="isEditing" class="icon" />
      <IconEdit v-else class="icon" />
      {{ isEditing ? t('display.done_editing') : t('display.refine_plan') }}
    </button>
    <SharePlanButton v-if="showShareButton" :store="store" id="tutorial-share-btn" />
    <ExportPlanButton :store="store" />
  </div>
</template>

<style scoped>
.training-plan-display {
  background: var(--color-background-soft);
  border-radius: 8px;
  border-top-right-radius: 11px;
  border-top-left-radius: 11px;
  border: 1px solid var(--color-border);
}

.plan-container {
  background: var(--color-background);
  border-radius: 8px;
  box-shadow: 0 2px 10px var(--color-shadow);
}

@media (max-width: 740px) {
  .training-plan-display {
    border-radius: 6px;
    border-top-right-radius: 8px;
    border-top-left-radius: 8px;
  }

  .plan-header {
    padding: 1.25rem 1rem;
  }

  .plan-title {
    font-size: 1.15rem;
    margin-bottom: 0.5rem;
  }

  .plan-description {
    font-size: 0.875rem;
  }

  .plan-cards-list {
    padding: 0.75rem;
    gap: 0.4rem;
  }

  .total-summary-row {
    padding: 0.6rem 0.75rem;
    font-size: 0.875rem;
  }

  .summary-section {
    padding: 0 0.75rem 0.75rem 0.75rem;
    gap: 0.5rem;
  }

  .summary-item {
    padding: 0.75rem 0.5rem;
  }

  .summary-value {
    font-size: 1.15rem;
  }

  .summary-label {
    font-size: 0.65rem;
    letter-spacing: 0.5px;
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
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-background-soft);
  color: var(--color-text);
  text-align: center;
}

.edit-title::placeholder {
  color: var(--color-text);
}

.edit-description {
  font-size: 1rem;
  line-height: 1.6;
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-background-soft);
  color: var(--color-text);
  font-family: inherit;
  resize: vertical;
}

.edit-description::placeholder {
  color: var(--color-text);
}

.edit-title:focus,
.edit-description:focus {
  outline: 2px solid var(--color-text);
  border: 1px solid var(--color-primary);
}

.plan-cards-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1.25rem;
  background: var(--color-background-soft);
}

.total-summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1.25rem;
  background: var(--color-background-mute);
  border-top: 2px solid var(--color-primary);
  font-weight: 700;
  font-size: 1rem;
  color: var(--color-heading);
}

.total-summary-label {
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-size: 0.85rem;
  opacity: 0.75;
}

.total-summary-value {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--color-primary);
}

.summary-section {
  display: flex;
  justify-content: space-around;
  padding: 1rem;
  background: var(--color-background-soft);
  gap: 1rem;
  border-bottom-right-radius: 8px;
  border-bottom-left-radius: 8px;
  border-top: 1px solid var(--color-border);
}

.summary-item {
  background: var(--color-background);
  padding: 1rem;
  border-radius: 8px;
  text-align: center;
  flex: 1;
  border: 1px solid var(--color-border);
  box-shadow: 0 1px 4px var(--color-shadow);
}

.summary-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--color-primary);
  margin-bottom: 0.25rem;
  line-height: 1;
}

.summary-label {
  color: var(--color-heading);
  text-transform: uppercase;
  font-size: 0.7rem;
  letter-spacing: 1.5px;
  opacity: 0.7;
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
  border-radius: 8px;
  border: 1px solid var(--color-border);
  padding: 1.5rem;
  background: var(--color-background-soft);
  text-align: center;
  margin-top: 1rem;
  gap: 1rem;
  width: 100%;
}

.edit-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  width: fit-content;
  max-width: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

@media (max-width: 740px) {
  .edit-btn {
    padding: 0.75rem 0.5rem;
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
</style>
