<script setup lang="ts">
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { storeToRefs } from 'pinia'
import { onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const sharedPlanStore = useSharedPlanStore()

const { sharedPlan, isLoading, error } = storeToRefs(sharedPlanStore)

onMounted(async () => {
  const urlHash = route.params.urlHash
  if (typeof urlHash === 'string') {
    if (await sharedPlanStore.fetchSharedPlanByHash(urlHash)) return
    if (sharedPlan.value === null) {
      noPlanFound()
    }
  } else if (typeof urlHash === 'undefined' && sharedPlan.value === null) {
    noPlanFound()
  }
})

onUnmounted(() => {
  sharedPlanStore.clear()
})

function noPlanFound() {
  toast.error(t('shared.no_plan_toast', { error: error.value || '' }))
  router.push('/')
}
</script>

<template>
  <div class="shared-view">
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>{{ t('shared.loading') }}</p>
    </div>
    <div v-else-if="sharedPlan">
      <div class="container">
        <section class="hero">
          <h1>{{ t('shared.hero_title') }}</h1>
          <p class="hero-description">
            {{ t('shared.hero_description', { username: sharedPlan.sharer_username }) }}
          </p>
        </section>

        <!-- Main content -->
        <section>
          <TrainingPlanDisplay :store="sharedPlanStore" :show-share-button="false" />
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.shared-view {
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
</style>
