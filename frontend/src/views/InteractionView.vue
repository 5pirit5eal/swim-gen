<script setup lang="ts">
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import SimplePlanDisplay from '@/components/training/SimplePlanDisplay.vue'
import IconSend from '@/components/icons/IconSend.vue'
import IconStar from '@/components/icons/IconStar.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useAuthStore } from '@/stores/auth'
import type { RAGResponse, Message } from '@/types'
import { storeToRefs } from 'pinia'
import { onActivated, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import FeedbackForm from '@/components/forms/FeedbackForm.vue'

import { useTutorial } from '@/tutorial/useTutorial'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const trainingStore = useTrainingPlanStore()
const authStore = useAuthStore()
const { startInteractionTutorial } = useTutorial()

const { currentPlan, isLoading, isFetchingConversation, error, conversation, historyMetadata } =
  storeToRefs(trainingStore)

const planMetadata = ref<
  | {
      plan_id: string
      created_at: string
      updated_at: string
      exported_at?: string
      feedback_rating?: number
    }
  | undefined
>()

// Track which messages have expanded plan snapshots
const expandedSnapshots = ref<Set<string>>(new Set())
const chatInput = ref('')
const displayedMessages = ref<Message[]>([])
const showFeedbackForm = ref(false)

// Layout & Tabs
const activeTab = ref<'plan' | 'chat'>('plan')

watch(
  () => conversation.value,
  async (newVal) => {
    if (newVal.length === 0) {
      displayedMessages.value = []
      return
    }

    // If we have a fresh load (displayedMessages is empty but newVal has many)
    if (displayedMessages.value.length === 0 && newVal.length > 0) {
      // Staggered load
      for (const msg of newVal) {
        displayedMessages.value.push(msg)
        await new Promise((resolve) => setTimeout(resolve, 150)) // 150ms delay for smooth transition
      }
    } else if (newVal.length > displayedMessages.value.length) {
      // New message(s) added
      const newMessages = newVal.slice(displayedMessages.value.length)
      for (const msg of newMessages) {
        displayedMessages.value.push(msg)
      }
    } else {
      // Reset or other change, just sync
      displayedMessages.value = [...newVal]
    }
  },
  { deep: true, immediate: true },
)

function toggleSnapshot(messageId: string) {
  if (expandedSnapshots.value.has(messageId)) {
    expandedSnapshots.value.delete(messageId)
  } else {
    expandedSnapshots.value.add(messageId)
  }
}

function isExpanded(messageId: string): boolean {
  return expandedSnapshots.value.has(messageId)
}

function handleSaveSnapshot(plan: RAGResponse) {
  trainingStore.saveSnapshot(plan)
  toast.success(t('interaction.snapshot_saved'))
}

async function handleSendMessage() {
  if (!chatInput.value.trim()) return

  const message = chatInput.value
  chatInput.value = ''
  await trainingStore.sendMessage(message)
}

async function handleFeedbackSubmit(payload: {
  rating: number
  was_swam: boolean
  difficulty_rating: number
  comment?: string
}) {
  if (!currentPlan.value || !currentPlan.value.plan_id) return

  const success = await trainingStore.submitFeedback({
    plan_id: currentPlan.value.plan_id,
    rating: payload.rating,
    was_swam: payload.was_swam,
    difficulty_rating: payload.difficulty_rating,
    comment: payload.comment,
  })

  if (success) {
    toast.success(t('feedback.submit_success'))
    showFeedbackForm.value = false
    // Update metadata to reflect new rating
    planMetadata.value = getMetadata()
  }
}

async function initializeView() {
  const planId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  if (!planId) {
    console.log('No plan ID provided in route.')
    router.push('/')
    return
  }

  // If the plan is already loaded in the store and matches the route, we don't need to do anything
  if (currentPlan.value?.plan_id === planId) {
    return
  }

  if (trainingStore.planHistory.length === 0) {
    if (!trainingStore.isFetchingHistory) {
      await trainingStore.fetchHistory()
    } else {
      while (trainingStore.isFetchingHistory) {
        await new Promise((resolve) => setTimeout(resolve, 100))
      }
    }
  }

  const planFromHistory = trainingStore.planHistory.find((p) => p.plan_id === planId)

  if (planFromHistory) {
    trainingStore.loadPlanFromHistory(planFromHistory)
  } else {
    console.log('Plan not found in history.')
    toast.error(t('interaction.not_found'))
    router.push('/')
    return
  }
}

function getMetadata() {
  const planId = route.params.id as string
  return historyMetadata.value.find((m) => m.plan_id === planId)
}

onMounted(async () => {
  await initializeView()
  planMetadata.value = getMetadata()
  window.scrollTo(0, 0)
  startInteractionTutorial()
})

onUnmounted(() => {
  trainingStore.clear()
})

onActivated(() => {
  startInteractionTutorial()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      // Reset state for new plan
      expandedSnapshots.value.clear()
      await initializeView()
      // Update metadata for the new plan
      planMetadata.value = getMetadata()
      startInteractionTutorial()
    }
  },
)
</script>

<template>
  <div class="interaction-view">
    <div v-if="currentPlan" class="layout-container">
      <!-- Tab Switcher -->
      <div class="tab-header">
        <div class="tab-switcher" id="tutorial-tab-switcher">
          <button
            class="tab-button"
            :class="{ active: activeTab === 'plan' }"
            @click="activeTab = 'plan'"
          >
            {{ t('interaction.plan_tab') }}
          </button>
          <button
            class="tab-button"
            :class="{ active: activeTab === 'chat' }"
            @click="activeTab = 'chat'"
          >
            {{ t('interaction.conversation_tab') }}
          </button>
        </div>
        <button class="rate-plan-button" @click="showFeedbackForm = true" id="tutorial-rate-btn">
          {{ t('interaction.rate_plan') }}
        </button>
      </div>

      <!-- Plan Tab -->
      <Transition name="fade">
        <div class="tab-content" v-show="activeTab === 'plan'">
          <!-- Current Plan Display -->
          <section>
            <div v-if="isLoading" class="loading-state">
              <div class="loading-spinner"></div>
              <p>{{ t('shared.loading') }}</p>
            </div>
            <TrainingPlanDisplay v-else :store="trainingStore" :show-share-button="true" />
          </section>

          <!-- Metadata Section -->
          <section v-if="planMetadata" class="metadata-section">
            <h3>{{ t('interaction.metadata') }}</h3>
            <div class="metadata-grid">
              <div class="metadata-item">
                <span class="label">{{ t('interaction.created_at') }}</span>
                <span class="value">{{ new Date(planMetadata.created_at).toLocaleString() }}</span>
              </div>
              <div class="metadata-item">
                <span class="label">{{ t('interaction.updated_at') }}</span>
                <span class="value">{{ new Date(planMetadata.updated_at).toLocaleString() }}</span>
              </div>
              <div class="metadata-item" v-if="planMetadata.exported_at">
                <span class="label">{{ t('interaction.exported_at') }}</span>
                <span class="value">{{ new Date(planMetadata.exported_at).toLocaleString() }}</span>
              </div>
              <div class="metadata-item" v-if="planMetadata.feedback_rating">
                <span class="label">{{ t('feedback.rating_label') }}</span>
                <div class="rating-value-container">
                  <span class="value">{{ planMetadata.feedback_rating }}</span>
                  <IconStar class="rating-icon" />
                </div>
              </div>
            </div>
          </section>
        </div>
      </Transition>

      <!-- Chat Tab -->
      <Transition name="fade">
        <div class="tab-content chat-container" v-show="activeTab === 'chat'">
          <!-- Chat Messages Area -->
          <div class="chat-messages">
            <div
              v-if="displayedMessages.length === 0 && !isFetchingConversation"
              class="empty-chat"
            >
              <p>{{ t('interaction.no_messages') }}</p>
            </div>
            <TransitionGroup name="message">
              <div
                v-for="message in displayedMessages"
                :key="message.id"
                :class="['message', `message-${message.role}`]"
              >
                <div class="message-header">
                  <span class="message-role">{{
                    message.role === 'user'
                      ? authStore.user?.user_metadata?.username || t('interaction.you')
                      : t('interaction.ai')
                  }}</span>
                  <span class="message-time">{{
                    new Date(message.created_at).toLocaleString()
                  }}</span>
                </div>

                <div class="message-content">
                  {{ message.content }}
                </div>

                <!-- Plan Snapshot (for AI messages) -->
                <div
                  v-if="message.plan_snapshot && message.role === 'ai'"
                  class="snapshot-container"
                >
                  <button @click="toggleSnapshot(message.id)" class="snapshot-toggle">
                    <span class="toggle-icon">{{ isExpanded(message.id) ? '▼' : '▶' }}</span>
                    {{
                      isExpanded(message.id)
                        ? t('interaction.hide_plan')
                        : t('interaction.show_plan')
                    }}
                  </button>

                  <div v-if="isExpanded(message.id)" class="snapshot-content">
                    <SimplePlanDisplay
                      :title="message.plan_snapshot.title"
                      :description="message.plan_snapshot.description"
                      :table="message.plan_snapshot.table"
                      :plan-id="message.plan_snapshot.plan_id"
                      @save="handleSaveSnapshot"
                    />
                  </div>
                </div>
              </div>
            </TransitionGroup>

            <!-- Three-dot loader while waiting for AI response -->
            <div v-if="isLoading" class="chat-loader">
              <div class="loader"></div>
            </div>
          </div>

          <!-- Chat Input Area -->
          <div class="chat-input-wrapper">
            <label class="input-label">{{ t('interaction.describe_changes') }}</label>
            <form @submit.prevent="handleSendMessage" class="chat-form">
              <input
                v-model="chatInput"
                type="text"
                :placeholder="t('interaction.chat_placeholder')"
                class="chat-input"
                :disabled="isLoading"
              />
              <button type="submit" class="send-button" :disabled="isLoading || !chatInput.trim()">
                <IconSend class="send-icon" />
              </button>
            </form>
          </div>
        </div>
      </Transition>
    </div>

    <div v-else class="error-state">
      <p>{{ error || t('interaction.not_found') }}</p>
    </div>

    <FeedbackForm
      :show="showFeedbackForm"
      :plan-title="currentPlan?.title || ''"
      @submit="handleFeedbackSubmit"
      @close="showFeedbackForm = false"
    />
  </div>
</template>

<style scoped>
.interaction-view {
  padding: 0.25rem 0 1rem 0;
  display: flex;
  flex-direction: column;
  max-width: 1080px;
  margin: 0 auto;
}

.layout-container {
  margin: 0 auto;
  padding: 0 1rem;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.tab-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* Chat Specific Styles */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 500px;
  background: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.tab-header {
  display: flex;
  margin: 1rem 0;
  padding: 0.5rem 0;
  position: relative;
  z-index: 10;
}

.tab-switcher {
  display: flex;
  gap: 1rem;
}

.tab-button {
  background: color-mix(in srgb, var(--color-background), transparent 40%);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--color-border);
  padding: 0.75rem 1.75rem;
  font-size: 1rem;
  font-weight: 700;
  color: var(--color-text);
  cursor: pointer;
  border-radius: 24px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px var(--color-shadow);
}

.tab-button:hover {
  background: color-mix(in srgb, var(--color-background), transparent 20%);
  transform: translateY(-1px);
}

.tab-button.active {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
  box-shadow: 0 4px 12px color-mix(in srgb, var(--color-primary), transparent 70%);
}

.rate-plan-button {
  background: color-mix(in srgb, var(--color-background), transparent 40%);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 2px solid var(--color-success);
  color: var(--color-success);
  padding: 0.75rem 1.75rem;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  border-radius: 24px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px var(--color-shadow);
  margin-left: auto;
}

.rate-plan-button:hover {
  background: color-mix(in srgb, var(--color-background), transparent 20%);
  transform: translateY(-1px);
  opacity: 0.9;
}

/* Existing Styles Refined */
.loading-state,
.error-state {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--color-text);
  font-style: italic;
  background: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin: 2rem auto;
  max-width: 1080px;
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

.metadata-section {
  padding: 1.5rem;
  background: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin: 1rem 0;
}

.metadata-section h3 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: var(--color-heading);
}

.metadata-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.metadata-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.metadata-item .label {
  font-size: 0.875rem;
  color: var(--color-text-soft);
}

.metadata-item .value {
  font-weight: 500;
  color: var(--color-text);
}

.rating-value-container {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--color-warning);
}

.rating-icon {
  width: 16px;
  height: 16px;
  fill: currentColor;
}

.feedback-btn {
  background-color: var(--color-background-soft);
  color: var(--color-success);
  border: 2px solid var(--color-success);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: opacity 0.2s;
  font-size: 0.875rem;
  align-self: flex-start;
}

.feedback-btn:hover {
  opacity: 0.9;
}

.message {
  background: var(--color-background-mute);
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid var(--color-border);
  max-width: 90%;
  position: relative;
  transition: all 0.3s ease;
}

.message-user {
  border-right: 3px solid var(--color-primary-hover);
  align-self: flex-end;
  margin-left: 10%;
  border-bottom-right-radius: 0;
  background: var(--color-background);
}

.message-ai {
  border-left: 3px solid var(--color-primary);
  align-self: flex-start;
  margin-right: 10%;
  border-bottom-left-radius: 0;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  gap: 1rem;
}

.message-role {
  font-weight: 600;
  color: var(--color-heading);
}

.message-time {
  color: var(--color-text-soft);
  font-size: 0.75rem;
}

.message-content {
  line-height: 1.6;
  color: var(--color-text);
  white-space: pre-wrap;
}

.snapshot-container {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.snapshot-toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--color-background);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.snapshot-toggle:hover {
  background: var(--color-primary);
  color: white;
}

.toggle-icon {
  font-size: 0.75rem;
}

.snapshot-content {
  margin-top: 1rem;
}

.chat-form {
  display: flex;
  gap: 1rem;
}

.chat-input-wrapper {
  padding: 1rem;
  background: var(--color-background-soft);
  border-top: 1px solid var(--color-border);
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
}

.input-label {
  display: block;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-heading);
  margin-bottom: 0.5rem;
}

.chat-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-background);
  color: var(--color-text);
  font-size: 1rem;
  transition: border-color 0.2s;
}

.chat-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-shadow);
}

.send-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.send-button:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: scale(1.05);
}

.send-icon {
  transform: rotate(45deg) translateX(-2px) translateY(1px);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Transitions */
.message-enter-active,
.message-leave-active {
  transition: all 0.5s ease;
}

.message-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.message-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.chat-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem 0 0;
}

/* Three-dot jumping loader */
.loader {
  width: 60px;
  aspect-ratio: 2;
  --_g: no-repeat radial-gradient(circle closest-side, var(--color-text) 90%, #0000);
  background:
    var(--_g) 0% 50%,
    var(--_g) 50% 50%,
    var(--_g) 100% 50%;
  background-size: calc(100% / 3) 50%;
  animation: l3 1s infinite linear;
}

@keyframes l3 {
  20% {
    background-position:
      0% 0%,
      50% 50%,
      100% 50%;
  }

  40% {
    background-position:
      0% 100%,
      50% 0%,
      100% 50%;
  }

  60% {
    background-position:
      0% 50%,
      50% 100%,
      100% 0%;
  }

  80% {
    background-position:
      0% 50%,
      50% 50%,
      100% 100%;
  }
}
</style>
