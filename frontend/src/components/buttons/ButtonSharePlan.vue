<script setup lang="ts">
import { useSharedPlanStore } from '@/stores/sharedPlan';
import type { PlanStore } from '@/types';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  store: PlanStore
}>()

const { t } = useI18n()

const sharedPlanStore = useSharedPlanStore()

const { isLoading, error, shareUrl } = storeToRefs(sharedPlanStore)

async function handleShare() {
  if (!props.store.currentPlan || !props.store.currentPlan.plan_id) {
    return
  }
  await props.store.toggleKeepForever(props.store.currentPlan.plan_id)
  await sharedPlanStore.createShareUrl({ plan_id: props.store.currentPlan.plan_id, method: 'link' })
}

async function copyUrl() {
  if (shareUrl.value) {
    await navigator.clipboard.writeText(shareUrl.value)
  }
}
</script>

<template>
  <div>
    <button @click="handleShare" :disabled="isLoading" class="share-btn">
      {{ isLoading ? t('share.fetching_url') : t('share.share_plan') }}
    </button>
    <div v-if="shareUrl">
      <p>{{ t('share.show_url') }}</p>
      <input type="text" :value="shareUrl" readonly />
      <button @click="copyUrl">{{ t('share.copy') }}</button>
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
  color: var(--color-error);
}
</style>
