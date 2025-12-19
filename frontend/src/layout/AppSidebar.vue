<script setup lang="ts">
import { apiClient } from '@/api/client'
import UploadForm from '@/components/forms/UploadForm.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import IconCopy from '@/components/icons/IconCopy.vue'
import IconCross from '@/components/icons/IconCross.vue'
import IconDots from '@/components/icons/IconDots.vue'
import IconHeart from '@/components/icons/IconHeart.vue'
import IconHourglass from '@/components/icons/IconHourglass.vue'
import IconPlus from '@/components/icons/IconPlus.vue'
import IconSearch from '@/components/icons/IconSearch.vue'
import IconShare from '@/components/icons/IconShare.vue'
import IconUpload from '@/components/icons/IconUpload.vue'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useUploadStore } from '@/stores/uploads'
import type { HistoryMetadata, RAGResponse, SharedHistoryItem } from '@/types'
import { isIOS } from '@/utils/platform'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useTutorial } from '@/tutorial/useTutorial'
import { toast } from 'vue3-toastify'

// Search debounce utility
function debounce<T extends (...args: Parameters<T>) => void>(fn: T, delay: number): T {
  let timeoutId: ReturnType<typeof setTimeout>
  return ((...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }) as T
}

const trainingPlanStore = useTrainingPlanStore()
const sharedPlanStore = useSharedPlanStore()
const donationStore = useUploadStore()
const sidebarStore = useSidebarStore()
const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const { startSidebarTutorial } = useTutorial()

// Track which menu is open
const openMenuPlanId = ref<string | null>(null)
const editingPlanId = ref<string | null>(null)
const editingTitle = ref('')

// Share functionality state (two-step process for iOS compatibility)
const shareUrl = ref<string | null>(null)
const sharingPlanId = ref<string | null>(null)
const sharingSuccess = ref<string | null>(null)

// Compute currently viewed plan ID from route
const currentPlanId = computed(() => {
  if (route.name === 'plan' && route.params.id) {
    return route.params.id as string
  }
  // For shared plans, check if we have a loaded shared plan
  if (
    (route.name === 'shared' || route.name === 'shared_empty') &&
    sharedPlanStore.sharedPlan?.plan?.plan_id
  ) {
    return sharedPlanStore.sharedPlan.plan.plan_id
  }
  if (route.name === 'uploaded' && route.params.planId) {
    return route.params.planId as string
  }
  return null
})

// Search functionality
const searchQuery = ref('')
const debouncedSearch = debounce((query: string) => {
  trainingPlanStore.searchPlans(query)
}, 300)

watch(searchQuery, (query) => {
  debouncedSearch(query)
})

watch(
  () => sidebarStore.isOpen,
  (isOpen) => {
    if (isOpen) {
      setTimeout(() => {
        startSidebarTutorial()
      }, 300)
    }
  },
)

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

function toggleMenu(planId: string) {
  openMenuPlanId.value = openMenuPlanId.value === planId ? null : planId
}

function closeMenu() {
  openMenuPlanId.value = null
}

async function deletePlan(planId: string) {
  if (!confirm(t('sidebar.confirm_delete_plan'))) {
    return
  }
  const result = await apiClient.deletePlan(planId)
  if (result.success) {
    await trainingPlanStore.fetchHistory()
    // If we're currently viewing this plan, go home
    if (currentPlanId.value === planId) {
      router.push('/')
    }
  } else {
    alert('Failed to delete plan: ' + (result.error?.message || 'Unknown error'))
  }
  closeMenu()
}

function startEditTitle(plan: RAGResponse & HistoryMetadata) {
  editingPlanId.value = plan.plan_id
  editingTitle.value = plan.title
  closeMenu()
}

async function saveTitle(planId: string) {
  if (!editingTitle.value.trim()) {
    editingPlanId.value = null
    return
  }

  const plan = trainingPlanStore.planHistory.find((p) => p.plan_id === planId)
  if (!plan) return

  const result = await apiClient.upsertPlan({
    plan_id: planId,
    title: editingTitle.value,
    description: plan.description,
    table: plan.table,
  })

  if (result.success) {
    await trainingPlanStore.fetchHistory()
  }
  editingPlanId.value = null
}

function cancelEdit() {
  editingPlanId.value = null
  editingTitle.value = ''
}

function startEditTitleUpload(plan: { plan_id: string; title: string }) {
  editingPlanId.value = plan.plan_id
  editingTitle.value = plan.title
  closeMenu()
}

async function saveUploadedTitle(planId: string) {
  if (!editingTitle.value.trim()) {
    editingPlanId.value = null
    return
  }

  const plan = donationStore.uploadedPlans.find((p) => p.plan_id === planId)
  if (!plan) return

  const result = await apiClient.upsertPlan({
    plan_id: planId,
    title: editingTitle.value,
    description: plan.description,
    table: plan.table,
  })

  if (result.success) {
    await donationStore.fetchUploadedPlans()
  }
  editingPlanId.value = null
}

// Helper to copy URL to clipboard and show success state
async function copyToClipboard(url: string, planId: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(url)
    toast.success(t('share.copied'))

    sharingSuccess.value = planId
    setTimeout(() => {
      sharingSuccess.value = null
      shareUrl.value = null
      sharingPlanId.value = null
    }, 2000)
    return true
  } catch (err) {
    console.error('Failed to copy:', err)
    toast.error(t('share.copy_error'))
    return false
  }
}

async function sharePlan(plan: RAGResponse & HistoryMetadata) {
  // iOS two-step process: Step 2 - Copy existing URL to clipboard
  if (isIOS() && shareUrl.value && sharingPlanId.value === plan.plan_id) {
    await copyToClipboard(shareUrl.value, plan.plan_id)
    return
  }

  // Generate share URL
  try {
    if (!plan.keep_forever) await trainingPlanStore.toggleKeepForever(plan.plan_id)
    const result = await apiClient.createShareUrl({ plan_id: plan.plan_id, method: 'link' })

    if (result.success && result.data) {
      const url = `${window.location.origin}/shared/${result.data.url_hash}`
      if (isIOS()) {
        // iOS: Store URL for second click to copy
        shareUrl.value = url
        sharingPlanId.value = plan.plan_id
      } else {
        // Non-iOS: Copy immediately in same user gesture
        await copyToClipboard(url, plan.plan_id)
      }
    } else {
      toast.error(t('share.create_error'))
    }
  } catch (err) {
    console.error('Failed to create share URL:', err)
    toast.error(t('share.create_error'))
  }
}

async function deleteSharedPlan(planId: string) {
  if (!confirm(t('sidebar.confirm_delete_shared'))) {
    return
  }
  const result = await apiClient.deletePlan(planId)
  if (result.success) {
    await sharedPlanStore.fetchSharedHistory()
    // If we're currently viewing this plan, go home
    if (currentPlanId.value === planId) {
      router.push('/')
    }
  } else {
    alert('Failed to delete plan: ' + (result.error?.message || 'Unknown error'))
  }
  closeMenu()
}

async function deleteUploadedPlan(planId: string) {
  if (!confirm(t('sidebar.confirm_delete_uploaded'))) {
    return
  }
  const result = await apiClient.deletePlan(planId)
  if (result.success) {
    await donationStore.fetchUploadedPlans()
    // If we're currently viewing this plan, go home
    if (currentPlanId.value === planId) {
      router.push('/')
    }
  } else {
    alert('Failed to delete plan: ' + (result.error?.message || 'Unknown error'))
  }
  closeMenu()
}

async function shareUploadedPlan(plan: { plan_id: string; title: string }) {
  // iOS two-step process: Step 2 - Copy existing URL to clipboard
  if (isIOS() && shareUrl.value && sharingPlanId.value === plan.plan_id) {
    await copyToClipboard(shareUrl.value, plan.plan_id)
    return
  }

  // Generate share URL
  try {
    const result = await apiClient.createShareUrl({ plan_id: plan.plan_id, method: 'link' })

    if (result.success && result.data) {
      const url = `${window.location.origin}/shared/${result.data.url_hash}`
      if (isIOS()) {
        // iOS: Store URL for second click to copy
        shareUrl.value = url
        sharingPlanId.value = plan.plan_id
      } else {
        // Non-iOS: Copy immediately in same user gesture
        await copyToClipboard(url, plan.plan_id)
      }
    } else {
      toast.error(t('share.create_error'))
    }
  } catch (err) {
    console.error('Failed to create share URL:', err)
    toast.error(t('share.create_error'))
  }
}

// Close menu when clicking outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.menu-container')) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

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
          class="create-new-btn"
          :title="t('sidebar.create_new')"
          id="tutorial-new-plan-btn"
        >
          <IconPlus class="icon-small" />
          <span>{{ t('sidebar.create_new') }}</span>
        </button>
        <button
          @click="showDonationForm = true"
          class="create-new-btn"
          :title="t('sidebar.upload_plan')"
          id="tutorial-upload-btn"
        >
          <IconUpload class="icon-small" />
          <span>{{ t('sidebar.upload_plan') }}</span>
        </button>
      </div>
      <section>
        <div class="section-header" id="tutorial-history-generated">
          <h3>{{ t('sidebar.generated') }}</h3>
          <div v-if="trainingPlanStore.isFetchingHistory" class="loading-spinner"></div>
        </div>
        <!-- Search input -->
        <div class="search-container">
          <IconSearch class="search-icon" />
          <input
            v-model="searchQuery"
            type="search"
            :placeholder="t('sidebar.search_placeholder')"
            class="search-input"
          />
          <div v-if="trainingPlanStore.isSearching" class="loading-spinner small" />
        </div>
        <p v-if="trainingPlanStore.planHistory.length === 0 && !searchQuery">
          {{ t('sidebar.generated_placeholder') }}
        </p>
        <p
          v-else-if="trainingPlanStore.planHistory.length === 0 && searchQuery"
          class="search-info"
        >
          {{ t('sidebar.search_no_results') }}
        </p>
        <ul v-else class="plan-list">
          <li
            v-for="plan in trainingPlanStore.planHistory"
            :key="plan.plan_id"
            :class="{ 'active-plan': currentPlanId === plan.plan_id }"
          >
            <div class="plan-item-main">
              <div
                class="status-icon-container"
                :title="
                  plan.keep_forever
                    ? t('sidebar.tooltip_permanent')
                    : t('sidebar.tooltip_temporary')
                "
                @click.stop="trainingPlanStore.toggleKeepForever(plan.plan_id)"
              >
                <IconHeart v-if="plan.keep_forever" class="status-icon" />
                <IconHourglass v-else class="status-icon" />
              </div>
              <div v-if="editingPlanId === plan.plan_id" class="plan-title-edit">
                <input
                  ref="titleInputRef"
                  v-model="editingTitle"
                  type="text"
                  class="title-input"
                  @keyup.enter="saveTitle(plan.plan_id)"
                  @keyup.escape="cancelEdit"
                  @blur="saveTitle(plan.plan_id)"
                />
              </div>
              <div v-else class="plan-title" @click="loadPlan(plan)">
                <span>{{ plan.title }}</span>
              </div>
              <div class="menu-container">
                <button class="menu-button" @click.stop="toggleMenu(plan.plan_id)">
                  <IconDots class="dots-icon" />
                </button>
                <transition name="dropdown">
                  <div v-if="openMenuPlanId === plan.plan_id" class="dropdown-menu">
                    <button class="menu-item" @click="startEditTitle(plan)">
                      {{ t('sidebar.menu_edit_title') }}
                    </button>
                    <button class="menu-item" @click="sharePlan(plan)">
                      <transition name="scale" mode="out-in">
                        <IconCheck v-if="sharingSuccess === plan.plan_id" class="menu-icon" />
                        <IconCopy
                          v-else-if="isIOS() && shareUrl && sharingPlanId === plan.plan_id"
                          class="menu-icon"
                        />
                        <IconShare v-else class="menu-icon" />
                      </transition>
                      <span v-if="sharingSuccess === plan.plan_id">{{ t('share.copied') }}</span>
                      <span v-else-if="isIOS() && shareUrl && sharingPlanId === plan.plan_id">{{
                        t('share.copy')
                      }}</span>
                      <span v-else>{{ t('sidebar.menu_share') }}</span>
                    </button>
                    <button class="menu-item delete" @click="deletePlan(plan.plan_id)">
                      {{ t('sidebar.menu_delete') }}
                    </button>
                  </div>
                </transition>
              </div>
            </div>
          </li>
        </ul>
        <!-- Search results info -->
        <p
          v-if="
            searchQuery &&
            trainingPlanStore.planHistory.length > 0 &&
            trainingPlanStore.searchHitLimit
          "
          class="search-info search-limit-warning"
        >
          {{ t('sidebar.search_limit_hit', { count: trainingPlanStore.planHistory.length }) }}
        </p>
        <p v-else-if="searchQuery && trainingPlanStore.planHistory.length > 0" class="search-info">
          {{ t('sidebar.search_results_info', { count: trainingPlanStore.planHistory.length }) }}
        </p>
        <!-- Load more button for generated plans -->
        <button
          v-if="trainingPlanStore.historyHasMore && !searchQuery"
          @click="trainingPlanStore.fetchMoreHistory()"
          :disabled="trainingPlanStore.isLoadingMore"
          class="load-more-btn"
        >
          <span v-if="trainingPlanStore.isLoadingMore">{{ t('common.loading') }}</span>
          <span v-else>{{ t('sidebar.load_more') }}</span>
        </button>
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
          <li
            v-for="item in sharedPlanStore.sharedHistory"
            :key="item.plan_id"
            :class="{ 'active-plan': currentPlanId === item.plan_id }"
          >
            <div class="plan-item-main">
              <div class="plan-title" @click="loadSharedPlan(item)">
                <span>{{ item.plan.title }}</span>
              </div>
              <div class="menu-container">
                <button class="menu-button" @click.stop="toggleMenu(item.plan_id)">
                  <IconDots class="dots-icon" />
                </button>
                <transition name="dropdown">
                  <div v-if="openMenuPlanId === item.plan_id" class="dropdown-menu">
                    <button class="menu-item delete" @click="deleteSharedPlan(item.plan_id)">
                      {{ t('sidebar.menu_delete') }}
                    </button>
                  </div>
                </transition>
              </div>
            </div>
          </li>
        </ul>
        <!-- Load more button for shared plans -->
        <button
          v-if="sharedPlanStore.historyHasMore"
          @click="sharedPlanStore.fetchMoreSharedHistory()"
          :disabled="sharedPlanStore.isLoadingMore"
          class="load-more-btn"
        >
          <span v-if="sharedPlanStore.isLoadingMore">{{ t('common.loading') }}</span>
          <span v-else>{{ t('sidebar.load_more') }}</span>
        </button>
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
          <li
            v-for="plan in donationStore.uploadedPlans"
            :key="plan.plan_id"
            :class="{ 'active-plan': currentPlanId === plan.plan_id }"
          >
            <div class="plan-item-main">
              <div v-if="editingPlanId === plan.plan_id" class="plan-title-edit">
                <input
                  ref="titleInputRef"
                  v-model="editingTitle"
                  type="text"
                  class="title-input"
                  @keyup.enter="saveUploadedTitle(plan.plan_id)"
                  @keyup.escape="cancelEdit"
                  @blur="saveUploadedTitle(plan.plan_id)"
                />
              </div>
              <div v-else class="plan-title" @click="loadUploadedPlan(plan.plan_id)">
                <span>{{ plan.title }}</span>
              </div>
              <div class="menu-container">
                <button class="menu-button" @click.stop="toggleMenu(plan.plan_id)">
                  <IconDots class="dots-icon" />
                </button>
                <transition name="dropdown">
                  <div v-if="openMenuPlanId === plan.plan_id" class="dropdown-menu">
                    <button class="menu-item" @click="startEditTitleUpload(plan)">
                      {{ t('sidebar.menu_edit_title') }}
                    </button>
                    <button class="menu-item" @click="shareUploadedPlan(plan)">
                      <transition name="scale" mode="out-in">
                        <IconCheck v-if="sharingSuccess === plan.plan_id" class="menu-icon" />
                        <IconCopy
                          v-else-if="isIOS() && shareUrl && sharingPlanId === plan.plan_id"
                          class="menu-icon"
                        />
                        <IconShare v-else class="menu-icon" />
                      </transition>
                      <span v-if="sharingSuccess === plan.plan_id">{{ t('share.copied') }}</span>
                      <span v-else-if="isIOS() && shareUrl && sharingPlanId === plan.plan_id">{{
                        t('share.copy')
                      }}</span>
                      <span v-else>{{ t('sidebar.menu_share') }}</span>
                    </button>
                    <button class="menu-item delete" @click="deleteUploadedPlan(plan.plan_id)">
                      {{ t('sidebar.menu_delete') }}
                    </button>
                  </div>
                </transition>
              </div>
            </div>
          </li>
        </ul>
        <!-- Load more button for uploaded plans -->
        <button
          v-if="donationStore.historyHasMore"
          @click="donationStore.fetchMoreUploadedPlans()"
          :disabled="donationStore.isLoadingMore"
          class="load-more-btn"
        >
          <span v-if="donationStore.isLoadingMore">{{ t('common.loading') }}</span>
          <span v-else>{{ t('sidebar.load_more') }}</span>
        </button>
      </section>
    </div>
    <UploadForm :show="showDonationForm" @close="showDonationForm = false" />
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: -400px;
  width: 400px;
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
  font-size: 1.5rem;
  font-weight: 600;
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
  font-size: 1rem;
  font-weight: 600;
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
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.sidebar-content section p {
  color: var(--color-text);
}

.sidebar-content h3 {
  text-align: left;
  font-size: 1.25rem;
  font-weight: 600;
  padding: 0 0 0.5rem 0;
  color: var(--color-heading);
}

.plan-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.plan-list li {
  border: 1px solid transparent;
  border-bottom: 1px solid var(--color-border);
  border-bottom-color: var(--color-border);
  color: var(--color-text);
  display: flex;
  flex-direction: column;
  padding: 0.5rem;
  font-size: 1.125rem;
  font-weight: 500;
}

/* Active/highlighted plan */
.plan-list li.active-plan {
  background-color: var(--color-shadow);
  border: 1px solid var(--color-primary-hover);
  border-radius: 8px;
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
  color: var(--color-text);
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
}

.plan-title:hover {
  color: var(--color-primary-hover);
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

/* Three-dot menu and dropdown */
.menu-container {
  position: relative;
  margin-left: auto;
}

.menu-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.menu-button:hover .dots-icon {
  color: var(--color-primary);
}

.dots-icon {
  width: 24px;
  height: 24px;
  color: var(--color-text);
  transition: color 0.2s;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 0.25rem;
  margin-right: 0.25rem;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  box-shadow: 0 4px 6px var(--color-shadow);
  z-index: 1000;
  min-width: 150px;
  overflow: hidden;
}

.menu-item {
  width: 100%;
  padding: 0.75rem 1rem;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  color: var(--color-text);
  font-size: 0.9rem;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.menu-item:hover {
  background-color: var(--color-background-mute);
}

.menu-item.delete {
  color: var(--color-error);
}

.menu-item.delete:hover {
  background-color: rgba(231, 76, 60, 0.1);
}

.menu-icon {
  width: 1rem;
  height: 1rem;
}

/* Dropdown transition */
.dropdown-enter-active {
  animation: dropdown-in 0.2s ease-out;
}

.dropdown-leave-active {
  animation: dropdown-out 0.15s ease-in;
}

@keyframes dropdown-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes dropdown-out {
  from {
    opacity: 1;
    transform: translateY(0);
  }

  to {
    opacity: 0;
    transform: translateY(-5px);
  }
}

/* Scale transition for icon changes */
.scale-enter-active,
.scale-leave-active {
  transition: all 0.2s ease;
}

.scale-enter-from,
.scale-leave-to {
  transform: scale(0.5);
  opacity: 0;
}

/* Title editing */
.plan-title-edit {
  flex: 1;
  padding: 0.25rem;
}

.title-input {
  width: 100%;
  padding: 0.5rem;
  background-color: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  color: var(--color-heading);
  font-size: 0.9rem;
  font-weight: 500;
}

.title-input:focus {
  outline: none;
  border-color: var(--color-primary);
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

@media (max-width: 1124px) {
  .sidebar {
    left: -300px;
    width: 300px;
    background-color: var(--color-background-soft);
    backdrop-filter: none;
  }

  .plan-list li {
    border: 1px solid transparent;
    border-bottom: 1px solid var(--color-border);
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    padding: 0.5rem;
    font-size: 0.85rem;
    font-weight: 500;
  }
}

@media (max-width: 740px) {
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

/* Search styles */
.search-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  margin: 1rem 1rem 0.5rem 1rem;
  background-color: var(--color-background-mute);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.search-icon {
  width: 18px;
  height: 18px;
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 0.95rem;
  color: var(--color-text);
  outline: none;
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-input::-webkit-search-cancel-button {
  -webkit-appearance: none;
  appearance: none;
  height: 16px;
  width: 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23888' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cline x1='18' y1='6' x2='6' y2='18'%3E%3C/line%3E%3Cline x1='6' y1='6' x2='18' y2='18'%3E%3C/line%3E%3C/svg%3E");
  cursor: pointer;
}

.loading-spinner.small {
  width: 14px;
  height: 14px;
  border-width: 2px;
  flex-shrink: 0;
}

/* Load more button */
.load-more-btn {
  display: block;
  width: calc(100% - 1rem);
  margin: 0.75rem 0.5rem;
  padding: 0.5rem 1rem;
  background-color: var(--color-background-mute);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  color: var(--color-text);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.load-more-btn:hover:not(:disabled) {
  background-color: var(--color-background-soft);
  border-color: var(--color-primary);
}

.load-more-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.search-info {
  font-size: 0.85rem;
  color: var(--color-text-muted);
  padding: 0.5rem 1rem;
  margin: 0;
  text-align: center;
  font-style: italic;
}
</style>
