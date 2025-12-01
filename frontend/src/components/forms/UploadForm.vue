<script setup lang="ts">
import { useUploadFormStore } from '@/stores/uploadForm'
import { apiClient } from '@/api/client'
import BaseModal from '@/components/ui/BaseModal.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import IconUpload from '@/components/icons/IconUpload.vue'
import { useI18n } from 'vue-i18n'
import { onMounted, ref } from 'vue'
import { toast } from 'vue3-toastify'

defineProps<{
    show: boolean
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const { t } = useI18n()
const uploadFormStore = useUploadFormStore()
const isUploadingImage = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
    uploadFormStore.initNewPlan()
})

function triggerFileInput() {
    fileInput.value?.click()
}

async function handleImageUpload(event: Event) {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return

    isUploadingImage.value = true
    try {
        const result = await apiClient.imageToPlan(file)
        if (result.success && result.data) {
            // Update the store with the extracted plan
            uploadFormStore.currentPlan = {
                ...uploadFormStore.currentPlan,
                ...result.data,
                // Ensure we keep the ID if it exists, or let the store handle it
            }
            console.log(uploadFormStore.currentPlan)
            toast.success(t('donation.image_upload_success'))
        } else {
            toast.error(t('donation.image_upload_error'))
        }
    } catch (e) {
        console.error(e)
        toast.error(t('donation.image_upload_error'))
    } finally {
        isUploadingImage.value = false
        // Reset file input
        if (fileInput.value) fileInput.value.value = ''
    }
}

async function submit() {
    if (!uploadFormStore.currentPlan?.title?.trim()) {
        toast.error(t('donation.title_required'))
        return
    }

    const success = await uploadFormStore.uploadCurrentPlan()
    if (success) {
        toast.success(t('donation.success_toast'))
        uploadFormStore.clear()
        emit('close')
    } else {
        toast.error(t('donation.error_toast', { error: uploadFormStore.error || '' }))
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

            <div class="image-upload-section">
                <input type="file" ref="fileInput" accept="image/*" class="hidden-input" @change="handleImageUpload" />
                <button class="upload-image-btn" @click="triggerFileInput" :disabled="isUploadingImage">
                    <IconUpload class="icon" />
                    {{ isUploadingImage ? t('donation.processing_image') : t('donation.upload_image_btn') }}
                </button>
                <p class="upload-hint">{{ t('donation.upload_image_hint') }}</p>
            </div>

            <div class="plan-editor">
                <TrainingPlanDisplay :store="uploadFormStore" :show-share-button="false" />
            </div>
        </template>

        <template #footer>
            <button class="submit-btn" @click="submit" :disabled="uploadFormStore.isLoading">
                {{ uploadFormStore.isLoading ? t('donation.loading') : t('donation.submit') }}
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

.image-upload-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
    padding: 1rem;
    border: 1px dashed var(--color-border);
    border-radius: 8px;
    background-color: var(--color-background-soft);
}

.hidden-input {
    display: none;
}

.upload-image-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: var(--color-background);
    border: 1px solid var(--color-border);
    padding: 0.5rem 1rem;
    border-radius: 8px;
    cursor: pointer;
    color: var(--color-heading);
    font-weight: 500;
    transition: all 0.2s;
}

.upload-image-btn:hover:not(:disabled) {
    border-color: var(--color-primary);
    color: var(--color-primary);
}

.upload-image-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.upload-hint {
    font-size: 0.875rem;
    color: var(--color-text-mute);
    margin: 0;
}

.icon {
    width: 18px;
    height: 18px;
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
