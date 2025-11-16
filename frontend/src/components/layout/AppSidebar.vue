<script setup lang="ts">
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useI18n } from 'vue-i18n'
import type { RAGResponse } from '@/types'

const trainingPlanStore = useTrainingPlanStore()
const sidebarStore = useSidebarStore()
const { t } = useI18n()

function loadPlan(plan: RAGResponse) {
  trainingPlanStore.loadPlanFromHistory(plan)
  sidebarStore.close()
}
</script>

<template>
  <aside class="sidebar" :class="{ 'is-open': sidebarStore.isOpen }">
    <div class="sidebar-header">
      <h3>{{ t('sidebar.history') }}</h3>
      <button @click="sidebarStore.close" class="close-btn">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
    <div class="sidebar-content">
      <ul class="plan-list">
        <li
          v-for="plan in trainingPlanStore.generationHistory"
          :key="plan.plan_id"
          @click="loadPlan(plan)"
        >
          {{ plan.title }}
        </li>
      </ul>
      <h3>{{ t('sidebar.donated') }}</h3>
      <p>{{ t('sidebar.donated_placeholder') }}</p>
      <h3>{{ t('sidebar.shared') }}</h3>
      <p>{{ t('sidebar.shared_placeholder') }}</p>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: -300px;
  width: 300px;
  height: 100%;
  background-color: var(--color-background-soft);
  border-right: 1px solid var(--color-border);
  transition: left 0.3s ease;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.sidebar.is-open {
  left: 0;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--color-heading);
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text);
}

.sidebar-content {
  padding: 1rem;
  overflow-y: auto;
}

.plan-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.plan-list li {
  padding: 0.75rem;
  cursor: pointer;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text);
}

.plan-list li:hover {
  background-color: var(--color-background-mute);
}

@media (max-width: 768px) {
  .sidebar {
    left: -100%;
    width: 100%;
  }

  .sidebar.is-open {
    left: 0;
  }
}
</style>
