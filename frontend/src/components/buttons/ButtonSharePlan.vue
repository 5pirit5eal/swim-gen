<script setup lang="ts">
import IconCheck from '@/components/icons/IconCheck.vue'
import IconCopy from '@/components/icons/IconCopy.vue'
import IconShare from '@/components/icons/IconShare.vue'
import type { PlanStore, ShareUrlRequest } from '@/types'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { apiClient, formatError } from '@/api/client'

const props = defineProps<{
  store: PlanStore
}>()

const { t } = useI18n()

const isLoading = ref(false)
const error = ref<string | null>(null)
const shareUrl = ref<string | null>(null)
const copied = ref(false)

// Reset share URL when the plan changes
watch(
  () => props.store.currentPlan,
  () => {
    shareUrl.value = null
  },
)

// Creates a shareable URL for a plan
async function createShareUrl(request: ShareUrlRequest): Promise<string | null> {
  isLoading.value = true
  error.value = null
  const result = await apiClient.createShareUrl(request)
  isLoading.value = false

  if (result.success && result.data) {
    shareUrl.value = `${window.location.origin}/shared/${result.data.url_hash}`
    return shareUrl.value
  } else {
    error.value = result.error ? formatError(result.error) : t('errors.share_plan_failed')
    return null
  }
}

async function handleShare() {
  if (!props.store.currentPlan || !props.store.currentPlan.plan_id) {
    return
  }
  await props.store.keepForever(props.store.currentPlan.plan_id)
  await createShareUrl({ plan_id: props.store.currentPlan.plan_id, method: 'link' })
}

async function copyUrl() {
  if (shareUrl.value) {
    try {
      await navigator.clipboard.writeText(shareUrl.value)
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }
}
</script>

<template>
  <div class="share-container">
    <transition name="fade" mode="out-in">
      <!-- Initial Share Button -->
      <button
        v-if="!shareUrl"
        @click="handleShare"
        :disabled="isLoading"
        class="share-btn"
        key="share-btn"
      >
        <span v-if="isLoading" class="loading-spinner"></span>
        <template v-else>
          <IconShare class="icon" />
          {{ t('share.share_plan') }}
        </template>
      </button>

      <!-- Copy Link Button (Success State) -->
      <button
        v-else
        @click="copyUrl"
        class="share-btn copy-link-btn"
        :class="{ copied: copied }"
        key="copy-btn"
      >
        <transition name="scale" mode="out-in">
          <IconCheck v-if="copied" class="icon" />
          <IconCopy v-else class="icon" />
        </transition>
        <span>{{ copied ? t('share.copied') : t('share.copy') }}</span>
      </button>
    </transition>

    <p v-if="error" class="error-message">{{ error }}</p>
  </div>
</template>

<style scoped>
.share-container {
  max-width: 200px;
  position: relative;
}

.share-btn {
  width: 100%;
  height: 100%;
  padding: 0.75rem 1rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.share-btn:hover {
  background-color: var(--color-primary-hover);
}

.share-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 740px) {
  .share-container {
    min-width: 90px;
  }

  .share-btn {
    padding: 0.25rem 0.5rem;
    overflow-wrap: break-word;
    font-size: 0.8rem;
  }
}

.copy-link-btn {
  background-color: var(--color-primary);
}

.copy-link-btn span {
  font-weight: 600;
}

.copy-link-btn.copied {
  background-color: var(--color-success);
  cursor: default;
}

.icon {
  width: 24px;
  height: 24px;
}

.loading-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid var(--color-text);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

.error-message {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 0.5rem;
  color: var(--color-error);
  font-size: 0.8rem;
  text-align: center;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.scale-enter-active,
.scale-leave-active {
  transition: all 0.2s ease;
}

.scale-enter-from,
.scale-leave-to {
  transform: scale(0.5);
  opacity: 0;
}
</style>
