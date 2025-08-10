<script setup lang="ts">
// Main page for the Swim Training Plan Generator
// This will contain the form and plan display components
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

      <!-- Main content will go here -->
      <section class="main-content">
        <div class="placeholder">
          <p>🏊‍♀️ Training plan form and display components will be added here</p>
          <p>📝 Form for training plan requests</p>
          <p>⚙️ Advanced settings panel</p>
          <p>📊 Generated plan display</p>
          <p>📄 PDF export functionality</p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  padding: 2rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.hero {
  text-align: center;
  margin-bottom: 3rem;
}

.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 1rem;
}

.hero-description {
  font-size: 1.125rem;
  color: var(--color-text-light);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

.main-content {
  background: var(--color-background-soft);
  border-radius: 0.5rem;
  padding: 2rem;
  border: 1px solid var(--color-border);
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

  .main-content {
    padding: 1.5rem;
  }
}
</style>
