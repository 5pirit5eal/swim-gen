<script setup lang="ts">
import IconCheck from '@/components/icons/IconCheck.vue'
import IconCopy from '@/components/icons/IconCopy.vue'
import IconShare from '@/components/icons/IconShare.vue'
import type { PlanStore, ShareUrlRequest } from '@/types'
import { isIOS } from '@/utils/platform'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { apiClient } from '@/api/client'
import { toast } from 'vue3-toastify'

const props = defineProps<{
  store: PlanStore
}>()

const { t } = useI18n()

const isLoading = ref(false)
const shareUrl = ref<string | null>(null)
const justCopied = ref(false)

// Creates a shareable URL for a plan
async function createShareUrl(request: ShareUrlRequest): Promise<string | null> {
  isLoading.value = true
  const result = await apiClient.createShareUrl(request)
  isLoading.value = false

  if (result.success && result.data) {
    return `${window.location.origin}/shared/${result.data.url_hash}`
  } else {
    toast.error(t('share.create_error'))
    return null
  }
}

// Helper to copy URL to clipboard and show success state
async function copyToClipboard(url: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(url)
    toast.success(t('share.copied'))

    justCopied.value = true
    setTimeout(() => {
      justCopied.value = false
      shareUrl.value = null
    }, 2000)
    return true
  } catch (err) {
    console.error('Failed to copy:', err)
    toast.error(t('share.copy_error'))
    return false
  }
}

async function handleShare() {
  if (!props.store.currentPlan || !props.store.currentPlan.plan_id) {
    return
  }

  // iOS two-step process: Step 2 - Copy existing URL to clipboard
  if (isIOS() && shareUrl.value) {
    await copyToClipboard(shareUrl.value)
    return
  }

  // Generate share URL
  try {
    await props.store.keepForever(props.store.currentPlan.plan_id)

    const url = await createShareUrl({
      plan_id: props.store.currentPlan.plan_id,
      method: 'link',
    })

    if (url) {
      if (isIOS()) {
        // iOS: Store URL for second click to copy
        shareUrl.value = url
      } else {
        // Non-iOS: Copy immediately in same user gesture
        await copyToClipboard(url)
      }
    }
  } catch (err) {
    console.error('Failed to create share URL:', err)
    toast.error(t('share.create_error'))
  }
}
</script>

<template>
  <div class="share-container">
    <button @click="handleShare" :disabled="isLoading" class="share-btn" :class="{ success: justCopied }">
      <span v-if="isLoading" class="loading-spinner"></span>
      <template v-else>
        <transition name="scale" mode="out-in">
          <IconCheck v-if="justCopied" class="icon" />
          <IconCopy v-else-if="shareUrl && isIOS()" class="icon" />
          <IconShare v-else class="icon" />
        </transition>
        <span v-if="justCopied">{{ t('share.copied') }}</span>
        <span v-else-if="shareUrl && isIOS()">{{ t('share.copy') }}</span>
        <span v-else>{{ t('share.share_plan') }}</span>
      </template>
    </button>
  </div>
</template>

<style scoped>
.share-container {
  max-width: 200px;
  position: relative;
}

.share-btn {
  width: fit-content;
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

.share-btn.success {
  background-color: var(--color-success);
}

.share-btn.success:hover {
  background-color: var(--color-success);
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
    padding: 0.75rem 0.5rem;
    overflow-wrap: break-word;
    font-size: 0.8rem;
  }
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

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Icon transition */
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
