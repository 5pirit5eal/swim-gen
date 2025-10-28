<script setup lang="ts">
import TrainingPlanForm from '@/components/forms/TrainingPlanForm.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
import { useI18n } from 'vue-i18n'
// Header component for the swim training plan generator
// Currently minimal for V1 (single page), expandable for V2

const { t } = useI18n()
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
        <TrainingPlanDisplay />
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
  margin-bottom: 2rem;
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
  color: var(--color-text-light);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

.placeholder {
  text-align: center;
  color: var(--color-text-light);
}

.placeholder p {
  margin: 0.5rem 0;
  font-size: 1rem;
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
