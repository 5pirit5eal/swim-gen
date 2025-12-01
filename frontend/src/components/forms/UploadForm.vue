<script setup lang="ts">
import { useUploadFormStore } from '@/stores/uploadForm'
import { apiClient } from '@/api/client'
import BaseModal from '@/components/ui/BaseModal.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import IconUpload from '@/components/icons/IconUpload.vue'
import { useI18n } from 'vue-i18n'
import { onMounted, ref } from 'vue'
import { toast } from 'vue3-toastify'
import IconCross from '@/components/icons/IconCross.vue'

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
const imagePreviewUrl = ref<string | null>(null)
const isPreviewOpen = ref(false)
const allowSharing = ref(false)

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

    // Create preview URL
    const reader = new FileReader()
    reader.onload = (e) => {
        if (e.target?.result) {
            imagePreviewUrl.value = e.target.result as string
        }
    }
    reader.readAsDataURL(file)

    try {
        const language = navigator.language.split('-')[0] || 'en'
        const result = await apiClient.imageToPlan(file, language)
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

    const success = await uploadFormStore.uploadCurrentPlan(allowSharing.value)
    if (success) {
        toast.success(t('donation.success_toast'))
        uploadFormStore.clear()
        imagePreviewUrl.value = null
        allowSharing.value = false
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

                <div v-if="imagePreviewUrl" class="preview-container">
                    <img :src="imagePreviewUrl" class="preview-image" @click="isPreviewOpen = true" />
                    <button class="remove-preview" @click="imagePreviewUrl = null">
                        <IconCross />
                    </button>
                </div>

                <button class="upload-image-btn" @click="triggerFileInput" :disabled="isUploadingImage">
                    <IconUpload class="icon" />
                    {{ isUploadingImage ? t('donation.processing_image') : t('donation.upload_image_btn') }}
                </button>
                <p class="upload-hint">{{ t('donation.upload_image_hint') }}</p>
            </div>

            <div v-if="isPreviewOpen" class="preview-modal" @click="isPreviewOpen = false">
                <div class="preview-content">
                    <img :src="imagePreviewUrl || ''" />
                    <button class="close-preview" @click="isPreviewOpen = false">×</button>
                </div>
            </div>

            <div class="plan-editor">
                <TrainingPlanDisplay :store="uploadFormStore" :show-share-button="false" />
            </div>

            <div class="sharing-option">
                <label class="checkbox-label">
                    <input type="checkbox" v-model="allowSharing" />
                    <span>{{ t('donation.allow_sharing') }}</span>
                </label>
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

.preview-container {
    position: relative;
    margin-bottom: 1rem;
    max-width: 200px;
}

.preview-image {
    width: 100%;
    border-radius: 8px;
    cursor: zoom-in;
    border: 1px solid var(--color-border);
}

.remove-preview {
    position: absolute;
    top: -8px;
    right: -8px;
    background: var(--color-error);
    color: white;
    border: none;
    border-radius: 50%;
    width: 20px;
    height: 20px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
}

.preview-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    cursor: zoom-out;
}

.preview-content {
    position: relative;
    max-width: 90vw;
    max-height: 90vh;
}

.preview-content img {
    max-width: 100%;
    max-height: 90vh;
    border-radius: 8px;
}

.close-preview {
    position: absolute;
    top: -40px;
    right: 0;
    background: none;
    border: none;
    color: white;
    font-size: 2rem;
    cursor: pointer;
}

.sharing-option {
    margin: 0 2rem 1.5rem 2rem;
}

.checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    color: var(--color-text);
    user-select: none;
}

.checkbox-label input[type="checkbox"] {
    width: 1.2rem;
    height: 1.2rem;
    accent-color: var(--color-primary);
}
</style>
