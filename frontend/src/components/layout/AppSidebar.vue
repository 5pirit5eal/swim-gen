<script setup lang="ts">
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useI18n } from 'vue-i18n'
import type { RAGResponse } from '@/types'
import { useRouter } from 'vue-router'
import IconHourglass from '@/components/icons/IconHourglass.vue'
import IconHeart from '@/components/icons/IconHeart.vue'
import IconCross from '@/components/icons/IconCross.vue'

const trainingPlanStore = useTrainingPlanStore()
const sidebarStore = useSidebarStore()
const { t } = useI18n()
const router = useRouter()

function loadPlan(plan: RAGResponse) {
  trainingPlanStore.loadPlanFromHistory(plan)
  sidebarStore.close()
  router.push('/')
}
</script>

<template>
  <aside class="sidebar" :class="{ 'is-open': sidebarStore.isOpen }">
    <div class="sidebar-header">
      <button @click="sidebarStore.close" class="close-btn">
        <IconCross />
      </button>
      <h3>{{ t('sidebar.history') }}</h3>
    </div>
    <div class="sidebar-content">
      <section>
        <h3>{{ t('sidebar.generated') }}</h3>
        <p v-if="trainingPlanStore.planHistory.length === 0">
          {{ t('sidebar.generated_placeholder') }}
        </p>
        <ul v-else class="plan-list">
          <li v-for="plan in trainingPlanStore.planHistory" :key="plan.plan_id">
            <div class="plan-item-main">
              <div
                class="status-icon-container"
                @click.stop="trainingPlanStore.toggleKeepForever(plan.plan_id)"
              >
                <IconHeart v-if="plan.keep_forever" class="status-icon" />
                <IconHourglass v-else class="status-icon" />
              </div>
              <div class="plan-title" @click="loadPlan(plan)">
                <span>{{ plan.title }}</span>
              </div>
            </div>
          </li>
        </ul>
      </section>
      <section>
        <h3>{{ t('sidebar.donated') }}</h3>
        <p>{{ t('sidebar.donated_placeholder') }}</p>
      </section>
      <section>
        <h3>{{ t('sidebar.shared') }}</h3>
        <p>{{ t('sidebar.shared_placeholder') }}</p>
      </section>
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
  background-color: var(--color-transparent);
  backdrop-filter: blur(4px);
  border-right: 1px solid var(--color-border);
  transition: left 0.3s ease;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.sidebar.is-open {
  left: 0;
  border-top-right-radius: 8px;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1rem 1rem 1rem;
  border-bottom: 1px solid var(--color-border);
}

.sidebar-header h3 {
  font-size: 1.25rem;
  color: var(--color-heading);
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-heading);
}

.close-btn:hover {
  color: var(--color-error);
}

.sidebar-content {
  margin: 0.75rem;
  overflow-y: auto;
}

.sidebar-content section {
  margin-bottom: 1.5rem;
}

.sidebar-content section h3 {
  font-size: 1.125rem;
  margin-bottom: 0.5rem;
}

.sidebar-content section p {
  color: var(--color-text);
}

.sidebar-content h3 {
  text-align: left;
  font-size: 1rem;
  padding: 0 0 0.5rem 0;
  color: var(--color-heading);
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
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.plan-title {
  display: flex;
  align-items: center;
  color: var(--color-heading);
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.25rem;
}

.plan-title:hover {
  color: var(--color-text);
}

.plan-title span {
  font-weight: 500;
}

.status-icon {
  width: 1.5rem;
  height: 1.5rem;
  padding: 0.11rem;
  color: var(--color-primary);
}

.status-icon:hover {
  stroke: var(--color-primary-hover);
  stroke-width: 3px;
}

@media (max-width: 768px) {
  .sidebar {
    left: -100%;
    width: 100%;
    background-color: var(--color-background-soft);
    backdrop-filter: none;
  }

  .sidebar.is-open {
    left: 0;
  }
}
</style>
