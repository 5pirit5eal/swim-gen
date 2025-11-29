<script setup lang="ts">
import TrainingPlanForm from '@/components/forms/TrainingPlanForm.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { nextTick, onMounted, ref, watch } from 'vue'

const trainingPlanStore = useTrainingPlanStore()
const authStore = useAuthStore()
const { t } = useI18n()
const router = useRouter()
const planDisplayContainer = ref<HTMLDivElement | null>(null)

function navigateToLogin() {
  router.push({ name: 'login' })
}

function navigateToInteraction() {
  if (trainingPlanStore.currentPlan?.plan_id) {
    router.push({ name: 'plan', params: { id: trainingPlanStore.currentPlan.plan_id } })
  }
}

function scrollToPlan() {
  if (planDisplayContainer.value) {
    nextTick(() => {
      planDisplayContainer.value?.scrollIntoView?.({ behavior: 'smooth', block: 'nearest' })
    })
  }
}


watch(
  () => trainingPlanStore.currentPlan,
  (newPlan) => {
    if (newPlan) {
      scrollToPlan()
    }
  },
)

onMounted(async () => {
  if (trainingPlanStore.currentPlan) {
    scrollToPlan()
  }
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
        <TrainingPlanForm class="training-plan-form" />
        <div ref="planDisplayContainer">
          <TrainingPlanDisplay :store="trainingPlanStore" :show-share-button="!!authStore.user" />
          <div v-if="trainingPlanStore.currentPlan" class="cta-banner">
            <div v-if="!authStore.user" class="cta-content">
              <p>{{ t('home.banner.not_logged_in.text') }}</p>
              <button @click="navigateToLogin" class="cta-button">
                {{ t('home.banner.not_logged_in.button') }}
              </button>
            </div>
            <div v-else class="cta-content">
              <p>{{ t('home.banner.logged_in.text') }}</p>
              <button @click="navigateToInteraction" class="cta-button">
                {{ t('home.banner.logged_in.button') }}
              </button>
            </div>
          </div>
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

.training-plan-form {
  margin: 2rem auto;
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

.cta-banner {
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  text-align: center;
}

.cta-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  text-align: left;
}

.cta-content p {
  font-size: 1.1rem;
  color: var(--color-heading);
  max-width: 600px;
  margin: 0;
}

.cta-button {
  padding: 0.75rem 1.5rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.cta-button:hover {
  background-color: var(--color-primary-hover);
}
</style>
