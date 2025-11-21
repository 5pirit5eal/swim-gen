<script setup lang="ts">
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { storeToRefs } from 'pinia'
import type { PlanStore } from '@/types'

const props = defineProps<{
  store: PlanStore
}>()

const sharedPlanStore = useSharedPlanStore()

const { isLoading, error, shareUrl } = storeToRefs(sharedPlanStore)

async function handleShare() {
  if (!props.store.currentPlan || !props.store.currentPlan.plan_id) {
    return
  }
  await sharedPlanStore.createShareUrl({ plan_id: props.store.currentPlan.plan_id, method: 'link' })
}

function copyUrl() {
  if (shareUrl.value) {
    navigator.clipboard.writeText(shareUrl.value)
  }
}
</script>

<template>
  <div>
    <button @click="handleShare" :disabled="isLoading" class="share-btn">
      {{ isLoading ? 'Sharing...' : 'Share Plan' }}
    </button>
    <div v-if="shareUrl">
      <p>Share this URL:</p>
      <input type="text" :value="shareUrl" readonly />
      <button @click="copyUrl">Copy</button>
    </div>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<style scoped>
.share-btn {
  flex: 1;
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 160px;
}

.error {
  color: red;
}
</style>
