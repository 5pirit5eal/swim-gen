<script setup lang="ts">
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useI18n } from 'vue-i18n'
import type { RAGResponse } from '@/types'
import { useRouter } from 'vue-router'
import IconHourglass from '@/components/icons/IconHourglass.vue'
import IconHeart from '@/components/icons/IconHeart.vue'

const trainingPlanStore = useTrainingPlanStore()
const sidebarStore = useSidebarStore()
const { t } = useI18n()
const router = useRouter()

function loadPlan(plan: RAGResponse) {
  trainingPlanStore.loadPlanFromHistory(plan)
  sidebarStore.close()
  router.push('/')
}

function formatTimestamp(timestamp: string | undefined) {
  if (!timestamp) {
    return ''
  }
  const date = new Date(timestamp)
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>

<template>
  <aside class="sidebar" :class="{ 'is-open': sidebarStore.isOpen }">
    <div class="sidebar-header">
      <button @click="sidebarStore.close" class="close-btn">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
      <h3>{{ t('sidebar.history') }}</h3>
    </div>
    <div class="sidebar-content">
      <ul class="plan-list">
        <li v-for="plan in trainingPlanStore.planHistory" :key="plan.plan_id">
          <div class="plan-item-main">
            <div class="status-icon-container" @click.stop="trainingPlanStore.toggleKeepForever(plan.plan_id)">
              <IconHeart v-if="plan.keep_forever" class="status-icon" />
              <IconHourglass v-else class="status-icon" />
            </div>
            <div class="plan-title" @click="loadPlan(plan)">
              <span>{{ plan.title }}</span>
            </div>
          </div>
          <div class="plan-timestamps">
            <span>{{ t('labels.created_at') }}: {{ formatTimestamp(plan.created_at) }}</span>
            <span>{{ t('labels.updated_at') }}: {{ formatTimestamp(plan.updated_at) }}</span>
          </div>
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
  margin: 0.75rem;
  overflow-y: auto;
}

.plan-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.plan-list li {
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text);
  display: flex;
  flex-direction: column;
  padding: 0.5rem;
}

.plan-item-main {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.status-icon-container {
  padding: 0.25rem;
  border: 1px solid transparent;
  border-radius: 0.25rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-icon-container:hover {
  border-color: var(--color-border-hover);
  background-color: var(--color-background-mute);
}

.plan-title {
  display: flex;
  align-items: center;
  font-weight: 600;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.25rem;
}

.plan-title:hover {
  background-color: var(--color-background-mute);
}

.status-icon {
  width: 1.5rem;
  height: 1.5rem;
  padding: 0.15rem;
  color: var(--color-primary);
}

.plan-timestamps {
  font-size: 0.6rem;
  color: var(--color-text-mute);
  margin-top: 0.25rem;
  display: flex;
  flex-direction: column;
  padding: 0 0 0.5rem 2.75rem;
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
