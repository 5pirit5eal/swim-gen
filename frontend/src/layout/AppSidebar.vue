<script setup lang="ts">
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { useI18n } from 'vue-i18n'
import type { HistoryMetadata, RAGResponse, SharedHistoryItem } from '@/types'
import { useRouter } from 'vue-router'
import IconHourglass from '@/components/icons/IconHourglass.vue'
import IconHeart from '@/components/icons/IconHeart.vue'
import IconCross from '@/components/icons/IconCross.vue'
import IconPlus from '@/components/icons/IconPlus.vue'
import IconUpload from '@/components/icons/IconUpload.vue'
import UploadForm from '@/components/forms/UploadForm.vue'
import { useUploadStore } from '@/stores/uploads'
import { ref } from 'vue'

const trainingPlanStore = useTrainingPlanStore()
const sharedPlanStore = useSharedPlanStore()
const donationStore = useUploadStore()
const sidebarStore = useSidebarStore()
const { t } = useI18n()
const router = useRouter()

async function loadPlan(plan: RAGResponse & HistoryMetadata) {
  // Load plan and fetch conversation before navigation
  await trainingPlanStore.loadPlanFromHistory(plan)
  if (window.innerWidth <= 768) sidebarStore.close()
  router.push(`/plan/${plan.plan_id}`)
}

async function loadSharedPlan(plan: SharedHistoryItem) {
  await sharedPlanStore.loadPlanFromHistory(plan)
  if (window.innerWidth <= 768) sidebarStore.close()
  router.push('/shared/')
}

function createNewPlan() {
  trainingPlanStore.clear()
  if (window.innerWidth <= 768) sidebarStore.close()
  router.push('/')
}

const showDonationForm = ref(false)

async function loadUploadedPlan(plan_id: string) {
  await donationStore.loadPlanFromHistory(plan_id)
  if (window.innerWidth <= 768) sidebarStore.close()
  router.push(`/uploaded/${plan_id}`)
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
      <div class="action-buttons">
        <button
          @click="createNewPlan"
          class="create-new-btn secondary"
          :title="t('sidebar.create_new')"
        >
          <IconPlus class="icon-small" />
          <span>{{ t('sidebar.create_new') }}</span>
        </button>
        <button
          @click="showDonationForm = true"
          class="create-new-btn secondary"
          :title="t('sidebar.upload_plan')"
        >
          <IconUpload class="icon-small" />
          <span>{{ t('sidebar.upload_plan') }}</span>
        </button>
      </div>
      <section>
        <div class="section-header">
          <h3>{{ t('sidebar.generated') }}</h3>
          <div v-if="trainingPlanStore.isFetchingHistory" class="loading-spinner"></div>
        </div>
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
        <div class="section-header">
          <h3>{{ t('sidebar.shared') }}</h3>
          <div v-if="sharedPlanStore.isFetchingHistory" class="loading-spinner"></div>
        </div>
        <p v-if="sharedPlanStore.sharedHistory.length === 0">
          {{ t('sidebar.shared_placeholder') }}
        </p>
        <ul v-else class="plan-list">
          <li v-for="item in sharedPlanStore.sharedHistory" :key="item.plan_id">
            <div class="plan-item-main">
              <div class="plan-title" @click="loadSharedPlan(item)">
                <span>{{ item.plan.title }}</span>
              </div>
            </div>
          </li>
        </ul>
      </section>
      <section>
        <div class="section-header">
          <h3>{{ t('sidebar.uploaded') }}</h3>
          <div v-if="donationStore.isFetchingUploads" class="loading-spinner"></div>
        </div>
        <p v-if="donationStore.uploadedPlans.length === 0">
          {{ t('sidebar.uploaded_placeholder') }}
        </p>
        <ul v-else class="plan-list">
          <li v-for="plan in donationStore.uploadedPlans" :key="plan.plan_id">
            <div class="plan-item-main">
              <div class="plan-title" @click="loadUploadedPlan(plan.plan_id)">
                <span>{{ plan.title }}</span>
              </div>
            </div>
          </li>
        </ul>
      </section>
    </div>
    <UploadForm :show="showDonationForm" @close="showDonationForm = false" />
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
  color: var(--color-primary-hover);
}

.sidebar-content {
  margin: 0.75rem;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: var(--color-text) var(--color-shadow);
}

.sidebar-content section {
  margin-bottom: 1.5rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  width: 100%;
  justify-content: space-evenly;
}

.create-new-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  gap: 0.5rem;
  font-weight: 500;
  font-size: 0.9rem;
}

.create-new-btn:hover {
  background-color: var(--color-primary-hover);
}

.create-new-btn.secondary {
  background-color: var(--color-background-soft);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.create-new-btn.secondary:hover {
  background-color: var(--color-background-mute);
}

.icon-small {
  width: 22px;
  height: 22px;
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
  border-radius: 8px;
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

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.section-header h3 {
  margin-bottom: 0 !important;
  padding-bottom: 0 !important;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--color-border);
  border-top: 2px solid var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
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
