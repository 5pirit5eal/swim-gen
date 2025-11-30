<script setup lang="ts">
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import IconSend from '@/components/icons/IconSend.vue'
import { useDonationStore } from '@/stores/uploads'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { storeToRefs } from 'pinia'
import { onMounted, onUnmounted, watch, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const donationStore = useDonationStore()

const { currentPlan, isLoading, error } = storeToRefs(donationStore)

const chatInput = ref('')

async function initializeView() {
    const planId = route.params.planId
    if (typeof planId === 'string') {
        if (await donationStore.fetchUploadedPlan(planId)) return
        if (currentPlan.value === null) {
            noPlanFound()
        }
    } else {
        noPlanFound()
    }
}

onMounted(async () => {
    await initializeView()
})

onUnmounted(() => {
    donationStore.clear()
})

watch(
    () => route.params.planId,
    async (newId) => {
        if (newId) {
            await initializeView()
        }
    },
)

function noPlanFound() {
    toast.error(t('donation.no_plan_toast', { error: error.value || '' }))
    router.push('/')
}

// Starts a conversation by adding the plan to the history
async function handleStartConversation() {
    if (!chatInput.value.trim() || !currentPlan.value) return

    const message = chatInput.value
    chatInput.value = ''

    const trainingPlanStore = useTrainingPlanStore()

    if (currentPlan.value.plan_id) {
        // Load the plan into the training plan store
        await trainingPlanStore.loadPlanFromHistory(currentPlan.value)

        // Send the message
        trainingPlanStore.sendMessage(message).catch((err) => {
            console.error('Failed to send initial message:', err)
            toast.error(t('errors.send_message_failed'))
        })

        // Navigate to the interaction view
        router.push({ name: 'plan', params: { id: currentPlan.value.plan_id } })
    } else {
        toast.error(t('errors.generic'))
    }
}
</script>

<template>
    <div class="donated-view">
        <Transition name="fade">
            <div v-if="isLoading" class="loading-state">
                <div class="loading-spinner"></div>
                <p>{{ t('common.loading') }}</p>
            </div>
            <div v-else-if="currentPlan">
                <div class="container">
                    <section class="hero">
                        <h1>{{ currentPlan.title }}</h1>
                        <p class="hero-description">{{ currentPlan.description }}</p>
                    </section>

                    <!-- Main content -->
                    <section class="training-plan">
                        <TrainingPlanDisplay :store="donationStore" :show-share-button="true" />
                    </section>

                    <!-- Chat Transition Area -->
                    <section class="chat-transition">
                        <label class="input-label">{{ t('shared.start_conversation') }}</label>
                        <form @submit.prevent="handleStartConversation" class="chat-form">
                            <input v-model="chatInput" type="text" :placeholder="t('interaction.chat_placeholder')"
                                class="chat-input" :disabled="isLoading" />
                            <button type="submit" class="send-button" :disabled="isLoading || !chatInput.trim()">
                                <IconSend class="send-icon" />
                            </button>
                        </form>
                    </section>
                </div>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
.donated-view {
    padding: 0.25rem 0 2rem 0;
}

.container {
    max-width: 1080px;
    margin: 0 auto;
    padding: 0 1rem;
}

.hero {
    text-align: center;
    background-color: var(--color-transparent);
    backdrop-filter: blur(2px);
    border-radius: 8px;
    padding: 1rem;
    margin: 2rem auto;
}

.hero h1 {
    font-size: 2.5rem;
    font-weight: 700;
    color: var(--color-heading);
    margin-bottom: 1rem;
}

.hero-description {
    font-size: 1.25rem;
    color: var(--color-heading);
    font-weight: 500;
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
}

@media (max-width: 740px) {
    .hero h1 {
        font-size: 2rem;
    }

    .hero-description {
        font-size: 1rem;
    }
}

.loading-state,
.no-plan,
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

.error-state {
    color: var(--color-error);
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

.training-plan {
    margin: 1rem auto;
}

.chat-transition {
    background: var(--color-background-soft);
    border-radius: 8px;
    border: 1px solid var(--color-border);
    padding: 1rem;
}

.input-label {
    display: block;
    font-size: 1rem;
    font-weight: 500;
    color: var(--color-heading);
    margin-bottom: 0.5rem;
}

.chat-form {
    display: flex;
    gap: 1rem;
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
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
