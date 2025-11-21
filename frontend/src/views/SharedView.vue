<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { storeToRefs } from 'pinia'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'

const route = useRoute()
const sharedPlanStore = useSharedPlanStore()

const { sharedPlan, isLoading, error } = storeToRefs(sharedPlanStore)

onMounted(async () => {
  const urlHash = route.params.urlHash
  if (typeof urlHash === 'string') {
    await sharedPlanStore.fetchSharedPlanByHash(urlHash)
  }
})

onUnmounted(() => {
  sharedPlanStore.clear()
})
</script>

<template>
  <div class="shared-view">
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading shared plan...</p>
    </div>
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
    </div>
    <div v-else-if="sharedPlan">
      <div class="shared-info">
        <p>Shared by: <strong>{{ sharedPlan.sharer_username }}</strong></p>
      </div>
      <TrainingPlanDisplay :store="sharedPlanStore" :show-share-button="false" />
    </div>
    <div v-else class="no-plan">
      <p>Plan not found.</p>
    </div>
  </div>
</template>

<style scoped>
.shared-view {
  padding: 2rem;
  max-width: 1080px;
  margin: 0 auto;
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
}

.error-state {
  color: red;
}

.shared-info {
  text-align: center;
  margin-bottom: 1rem;
  color: var(--color-text);
  font-size: 1.1rem;
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
</style>
