<script setup lang="ts">
import TrainingPlanForm from '@/components/forms/TrainingPlanForm.vue'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'
</script>

<template>
  <div class="home-view">
    <div class="container">
      <section class="hero">
        <h1>AI-Powered Swimming Training Plans</h1>
        <p class="hero-description">
          Generate personalized training plans for swimmers, triathletes, and coaches. Just describe
          what you're looking for, and our AI will create a custom plan for you.
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
  padding: 2rem 0;
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

@media (max-width: 768px) {
  .hero h1 {
    font-size: 2rem;
  }

  .hero-description {
    font-size: 1rem;
  }
}
</style>
