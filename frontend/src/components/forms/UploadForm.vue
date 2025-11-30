<script setup lang="ts">
import { useDonationStore } from '@/stores/uploads'
import BaseModal from '@/components/ui/BaseModal.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import { useI18n } from 'vue-i18n'
import { onMounted } from 'vue'
import { toast } from 'vue3-toastify'

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const { t } = useI18n()
const donationStore = useDonationStore()

onMounted(() => {
    donationStore.initNewPlan()
})

async function submit() {
    const success = await donationStore.uploadCurrentPlan()
    if (success) {
        toast.success(t('donation.success_toast'))
        emit('close')
    } else {
        toast.error(t('donation.error_toast', { error: donationStore.error || '' }))
    }
}

function close() {
    emit('close')
}
</script>

<template>
    <BaseModal :show="show" @close="close">
        <template #header>
            <h2>{{ t('donation.title') }}</h2>
        </template>

        <template #body>
            <p class="intro-text">{{ t('donation.intro') }}</p>
            <div class="plan-editor">
                <TrainingPlanDisplay :store="donationStore" :show-share-button="false" />
            </div>
        </template>

        <template #footer>
            <button class="submit-btn" @click="submit" :disabled="donationStore.isLoading">
                {{ donationStore.isLoading ? t('donation.loading') : t('donation.submit') }}
            </button>
        </template>
    </BaseModal>
</template>

<style scoped>
h2 {
    margin: 0 auto;
    color: var(--color-heading);
}

.intro-text {
    margin-bottom: 1.5rem;
    color: var(--color-text);
}

.plan-editor {
    margin: 1rem 2rem 1.5rem 2rem;
}

.submit-btn {
    background-color: var(--color-primary);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
    background-color: var(--color-primary-hover);
}

.submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>
