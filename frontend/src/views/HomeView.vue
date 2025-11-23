<script setup lang="ts">
import TrainingPlanForm from '@/components/forms/TrainingPlanForm.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

const trainingStore = useTrainingPlanStore()
const authStore = useAuthStore()
const { t } = useI18n()

const planDisplayContainer = ref<HTMLDivElement | null>(null)

function scrollToPlan() {
  if (planDisplayContainer.value) {
    nextTick(() => {
      planDisplayContainer.value?.scrollIntoView?.({ behavior: 'smooth', block: 'nearest' })
    })
  }
}

watch(
  () => trainingStore.currentPlan,
  (newPlan) => {
    if (newPlan) {
      scrollToPlan()
    }
  },
)

onMounted(() => {
  if (trainingStore.currentPlan) {
    scrollToPlan()
  }
})

onUnmounted(() => {
  trainingStore.clear()
})
</script>

<template>
  <div class="home-view">
    <div class="container">
      <section class="hero">
        <h1>{{ t('app.hero_title') }}</h1>
        <p class="hero-description">
          {{ t('app.hero_description') }}
        </p>
      </section>

      <!-- Main content -->
      <section>
        <TrainingPlanForm />
        <div ref="planDisplayContainer">
          <TrainingPlanDisplay :store="trainingStore" :show-share-button="!!authStore.user" />
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.home-view {
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
</style>
