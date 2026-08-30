<script setup lang="ts">
import TrainingPlanForm from '@/components/forms/TrainingPlanForm.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import DrillList from '@/components/drills/DrillList.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import IconSave from '@/components/icons/IconSave.vue'
import IconEdit from '@/components/icons/IconEdit.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { nextTick, onActivated, onMounted, ref, watch } from 'vue'
import { useTutorial } from '@/tutorial/useTutorial'

const trainingPlanStore = useTrainingPlanStore()
const authStore = useAuthStore()
const { t } = useI18n()
const router = useRouter()
const planDisplayContainer = ref<HTMLDivElement | null>(null)
const { startHomeTutorial } = useTutorial()

// Restore anonymous plan from localStorage if it exists (e.g. after OAuth redirect)

async function checkAndLinkAnonymousPlan() {
  console.debug(
    'Checking for anonymous plan linking...',
    authStore.user,
    trainingPlanStore.initialQuery,
  )
  if (authStore.user && trainingPlanStore.currentPlan && trainingPlanStore.initialQuery) {
    await trainingPlanStore.linkAnonymousPlan()
  }
}

// Restore anonymous plan from localStorage if it exists (e.g. after OAuth redirect)
onMounted(async () => {
  console.debug('Checking for anonymous plan restoration...')
  const savedQuery = localStorage.getItem('anonymousQuery')
  const savedPlan = localStorage.getItem('anonymousPlan')

  if (savedPlan && savedQuery) {
    try {
      const plan = JSON.parse(savedPlan)
      trainingPlanStore.currentPlan = plan
      trainingPlanStore.initialQuery = savedQuery
      // Clean up
      localStorage.removeItem('anonymousPlan')
      localStorage.removeItem('anonymousQuery')
    } catch (e) {
      console.error('Failed to restore anonymous plan', e)
    }
  }

  // Try to link immediately after potential restoration
  await checkAndLinkAnonymousPlan()
})

watch(
  () => authStore.user,
  async () => {
    await checkAndLinkAnonymousPlan()
  },
  { immediate: true, deep: true },
)

function navigateToLogin(register = false) {
  if (register) {
    router.push({ name: 'login', query: { register: 'true' } })
  } else {
    router.push({ name: 'login' })
  }
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
  if (authStore.user) {
    startHomeTutorial()
  }
})

onActivated(() => {
  if (trainingPlanStore.currentPlan) {
    scrollToPlan()
  }
  if (authStore.user) {
    startHomeTutorial()
  }
})

watch(
  () => authStore.user,
  (user) => {
    if (user) {
      startHomeTutorial()
    }
  },
)
</script>

<template>
  <div class="home-view">
    <div class="container">
      <section class="hero">
        <h1>{{ t('app.hero_title') }}</h1>
        <p class="hero-description">
          {{ t('app.hero_description') }}
        </p>
        <div class="hero-badges" role="list">
          <span class="hero-badge" role="listitem">
            <IconCheck class="hero-badge-icon" aria-hidden="true" />
            {{ t('app.hero_badge_free') }}
          </span>
          <span class="hero-badge" role="listitem">
            <IconCheck class="hero-badge-icon" aria-hidden="true" />
            {{ t('app.hero_badge_no_signup') }}
          </span>
          <span class="hero-badge" role="listitem">
            <IconCheck class="hero-badge-icon" aria-hidden="true" />
            {{ t('app.hero_badge_audience') }}
          </span>
        </div>
      </section>

      <!-- Main content -->
      <section>
        <TrainingPlanForm class="training-plan-form" />
        <div ref="planDisplayContainer">
          <TrainingPlanDisplay :store="trainingPlanStore" :show-share-button="!!authStore.user" />
          <div v-if="trainingPlanStore.currentPlan" class="cta-banner">
            <div v-if="!authStore.user" class="cta-content">
              <div class="cta-text-group">
                <h3 class="cta-title">{{ t('home.banner.not_logged_in.title') }}</h3>
                <p class="cta-description">{{ t('home.banner.not_logged_in.text') }}</p>
              </div>
              <div class="cta-actions">
                <button @click="navigateToLogin(true)" class="cta-button btn-primary">
                  <IconSave class="cta-button-icon" aria-hidden="true" />
                  {{ t('home.banner.not_logged_in.button') }}
                </button>
                <button
                  type="button"
                  @click="navigateToLogin(false)"
                  class="cta-secondary-link text-btn"
                >
                  {{ t('home.banner.not_logged_in.secondary_login') }}
                </button>
              </div>
            </div>
            <div v-else class="cta-content">
              <div class="cta-text-group">
                <h3 class="cta-title">{{ t('home.banner.logged_in.title') }}</h3>
                <p class="cta-description">{{ t('home.banner.logged_in.text') }}</p>
              </div>
              <div class="cta-actions">
                <button @click="navigateToInteraction" class="cta-button btn-primary">
                  <IconEdit class="cta-button-icon" aria-hidden="true" />
                  {{ t('home.banner.logged_in.button') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Featured Drill List Section -->
        <DrillList :featured-mode="true" :limit="4" />
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
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
  border: 1px solid var(--color-border);
  box-shadow: 0 4px 12px var(--color-shadow);
  border-radius: 8px;
  padding: 1.75rem 1.25rem;
  margin: 2rem auto 1.5rem auto;
}

.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 0.75rem;
}

.hero-description {
  font-size: 1.25rem;
  color: var(--color-heading);
  font-weight: 500;
  max-width: 640px;
  margin: 0 auto;
  line-height: 1.6;
}

.hero-badges {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 0.6rem;
  margin-top: 1.25rem;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.82rem;
  font-weight: 600;
  color: #ffffff;
  background-color: var(--color-primary);
  border: 1px solid var(--color-primary);
  border-radius: 6px;
  padding: 0.35rem 0.75rem;
  box-shadow: 0 1px 3px var(--color-shadow);
  cursor: default;
  user-select: none;
  letter-spacing: 0.02em;
}

.hero-badge-icon {
  width: 0.85rem;
  height: 0.85rem;
  flex-shrink: 0;
}

@media (max-width: 740px) {
  .hero {
    padding: 1.25rem 0.75rem;
    margin: 1.25rem auto 1rem auto;
  }

  .training-plan-form {
    margin: 1.25rem auto;
  }

  .hero h1 {
    font-size: 1.85rem;
    margin-bottom: 0.5rem;
  }

  .hero-description {
    font-size: 0.95rem;
  }

  .hero-badges {
    gap: 0.4rem;
    margin-top: 0.85rem;
  }

  .hero-badge {
    font-size: 0.75rem;
    padding: 0.25rem 0.55rem;
  }
}

.cta-banner {
  margin-top: 2rem;
  padding: 1.75rem;
  background-color: var(--color-background-soft);
  border-radius: 8px;
  border: 1px solid var(--color-border);
  box-shadow: 0 2px 8px var(--color-shadow);
}

.cta-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  text-align: left;
}

.cta-text-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  max-width: 620px;
}

.cta-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--color-heading);
  margin: 0;
}

.cta-description {
  font-size: 0.95rem;
  color: var(--color-text);
  line-height: 1.5;
  margin: 0;
}

.cta-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.cta-button {
  padding: 0.75rem 1.5rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  white-space: nowrap;
}

.cta-button:hover {
  background-color: var(--color-primary-hover);
}

.cta-button-icon {
  width: 1.1rem;
  height: 1.1rem;
  flex-shrink: 0;
}

.cta-secondary-link {
  background: none;
  border: none;
  color: var(--color-text);
  font-size: 0.82rem;
  text-decoration: underline;
  cursor: pointer;
  padding: 0.2rem;
  transition: color 0.2s;
}

.cta-secondary-link:hover {
  color: var(--color-primary);
}

@media (max-width: 740px) {
  .cta-banner {
    padding: 1.25rem;
    margin-top: 1.5rem;
  }

  .cta-content {
    flex-direction: column;
    align-items: stretch;
    text-align: left;
    gap: 1rem;
  }

  .cta-actions {
    width: 100%;
    align-items: stretch;
  }

  .cta-button {
    width: 100%;
  }

  .cta-secondary-link {
    text-align: center;
    margin-top: 0.25rem;
  }
}
</style>
