<script setup lang="ts">
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import SimplePlanDisplay from '@/components/training/SimplePlanDisplay.vue'
import IconSend from '@/components/icons/IconSend.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useAuthStore } from '@/stores/auth'
import type { RAGResponse, Message } from '@/types'
import { storeToRefs } from 'pinia'
import { onMounted, onUnmounted, ref, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const trainingStore = useTrainingPlanStore()
const authStore = useAuthStore()

const { currentPlan, isLoading, error, conversation, historyMetadata } = storeToRefs(trainingStore)

// Track which messages have expanded plan snapshots
const expandedSnapshots = ref<Set<string>>(new Set())
const chatInput = ref('')
const chatInputSection = ref<HTMLElement | null>(null)
const displayedMessages = ref<Message[]>([])

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
        await new Promise(resolve => setTimeout(resolve, 150)) // 150ms delay for smooth transition
      }
      await scrollToChatInput()
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
  { deep: true }
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

async function scrollToChatInput() {
  await nextTick()
  if (chatInputSection.value) {
    chatInputSection.value.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

async function initializeView() {
  const planId = route.params.id as string
  if (!planId) {
    console.log('No plan ID provided in route.')
    router.push('/')
    return
  }

  // If the plan is already loaded in the store and matches the route, we don't need to do anything
  if (currentPlan.value?.plan_id === planId && conversation.value.length > 0) {
    return
  }

  if (trainingStore.planHistory.length === 0) {
    await trainingStore.fetchHistory()
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

  await trainingStore.fetchConversation(planId)
  await scrollToChatInput()
}

onMounted(async () => {
  await initializeView()
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
    }
  }
)



const planMetadata = ref<{ plan_id: string; created_at: string; updated_at: string } | undefined>()

function getMetadata() {
  const planId = route.params.id as string
  return historyMetadata.value.find((m) => m.plan_id === planId)
}

onMounted(() => {
  planMetadata.value = getMetadata()
})
onUnmounted(() => {
  trainingStore.clear()
})
</script>

<template>
  <div class="interaction-view">
    <div v-if="currentPlan" class="container">
      <!-- Chat Messages -->
      <section v-if="displayedMessages.length !== 0" class="chat-section">
        <div class="messages">
          <TransitionGroup name="message">
            <div v-for="message in displayedMessages" :key="message.id" :class="['message', `message-${message.role}`]">
              <div class="message-header">
                <span class="message-role">{{
                  message.role === 'user' ? (authStore.user?.user_metadata?.username || t('interaction.you')) :
                    t('interaction.ai')
                }}</span>
                <span class="message-time">{{
                  new Date(message.created_at).toLocaleString()
                }}</span>
              </div>

              <div class="message-content">
                {{ message.content }}
              </div>

              <!-- Plan Snapshot (for AI messages) -->
              <div v-if="message.plan_snapshot && message.role === 'ai'" class="snapshot-container">
                <button @click="toggleSnapshot(message.id)" class="snapshot-toggle">
                  <span class="toggle-icon">{{ isExpanded(message.id) ? '▼' : '▶' }}</span>
                  {{ isExpanded(message.id) ? t('interaction.hide_plan') : t('interaction.show_plan') }}
                </button>

                <div v-if="isExpanded(message.id)" class="snapshot-content">
                  <SimplePlanDisplay :title="message.plan_snapshot.title"
                    :description="message.plan_snapshot.description" :table="message.plan_snapshot.table"
                    :plan-id="message.plan_snapshot.plan_id" @save="handleSaveSnapshot" />
                </div>
              </div>
            </div>
          </TransitionGroup>
        </div>
      </section>

      <!-- Current Plan Display -->
      <section class="current-plan-section">
        <div v-if="isLoading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>{{ t('shared.loading') }}</p>
        </div>
        <TrainingPlanDisplay v-else :store="trainingStore" :show-share-button="true" />
      </section>

      <!-- Chat Input -->
      <section class="chat-input-section" ref="chatInputSection">
        <form @submit.prevent="handleSendMessage" class="chat-form">
          <input v-model="chatInput" type="text"
            :placeholder="t('interaction.chat_placeholder', 'Nachricht eingeben...')" class="chat-input"
            :disabled="isLoading" />
          <button type="submit" class="send-button" :disabled="isLoading || !chatInput.trim()">
            <IconSend class="send-icon" />
          </button>
        </form>
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
        </div>
      </section>
    </div>

    <div v-else-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('shared.loading') }}</p>
    </div>

    <div v-else class="error-state">
      <p>{{ error || t('interaction.not_found') }}</p>
    </div>
  </div>
</template>

<style scoped>
.interaction-view {
  padding: 0.25rem 0 1rem 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

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
  max-width: 600px;
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

.chat-section,
.current-plan-section,
.metadata-section,
.chat-input-section {
  margin-bottom: 2rem;
}

.chat-section {
  margin-top: 2rem;
}

.metadata-section h3 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  color: var(--color-heading);
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message {
  background: var(--color-background-soft);
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid var(--color-border);
  max-width: 80%;
  position: relative;
  transition: all 0.3s ease;
}

.message-user {
  border-left: 3px solid var(--color-primary);
  align-self: flex-end;
  margin-left: 20%;
  border-bottom-right-radius: 0;
}

.message-ai {
  border-left: 3px solid var(--color-secondary, #6366f1);
  align-self: flex-start;
  margin-right: 20%;
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
  color: var(--color-text);
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
  border-radius: 4px;
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

.metadata-section {
  padding: 1.5rem;
  background: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
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

/* Chat Input Styles */
.chat-form {
  display: flex;
  gap: 1rem;
  background: var(--color-background-soft);
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.chat-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
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

.message-enter-from,
.message-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
