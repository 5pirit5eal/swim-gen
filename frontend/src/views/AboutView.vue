<script setup lang="ts">
// About page for Swim RAG application
</script>

<template>
  <div class="about-view">
    <div class="container">
      <h1>About Swim RAG</h1>

      <section class="intro">
        <p>
          Swim RAG is an AI-powered training plan generator designed specifically for swimmers,
          triathletes, and swimming coaches. Using advanced retrieval-augmented generation
          technology, we create personalized training plans based on your specific needs and goals.
        </p>
      </section>

      <section class="features">
        <h2>Features</h2>
        <div class="feature-grid">
          <div class="feature-card">
            <h3>🤖 AI-Powered Generation</h3>
            <p>Advanced AI creates personalized training plans based on your input</p>
          </div>
          <div class="feature-card">
            <h3>📝 Natural Language Input</h3>
            <p>Describe your training needs in plain language - no complex forms</p>
          </div>
          <div class="feature-card">
            <h3>📄 PDF Export</h3>
            <p>Export your training plans as professional PDFs for easy sharing</p>
          </div>
          <div class="feature-card">
            <h3>🔒 Privacy Focused</h3>
            <p>Anonymous usage with optional data donation for improvement</p>
          </div>
        </div>
      </section>

      <section class="technology">
        <h2>Technology</h2>
        <p>
          Built with modern web technologies including Vue.js, TypeScript, and a robust Go backend.
          Our RAG system ensures training plans are both scientifically sound and personalized to
          your specific requirements.
        </p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.about-view {
  padding: 2rem 0;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 1rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 2rem;
  text-align: center;
}

h2 {
  font-size: 1.5rem;
  color: var(--color-heading);
  margin: 2rem 0 1rem 0;
}

h3 {
  font-size: 1.125rem;
  color: var(--color-heading);
  margin-bottom: 0.5rem;
}

section {
  margin-bottom: 2rem;
}

.intro p {
  font-size: 1.125rem;
  line-height: 1.7;
  color: var(--color-text);
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-top: 1rem;
}

.feature-card {
  background: var(--color-background-soft);
  padding: 1.5rem;
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
}

.feature-card p {
  color: var(--color-text-light);
  line-height: 1.6;
}

.technology p {
  font-size: 1rem;
  line-height: 1.7;
  color: var(--color-text);
}

@media (max-width: 768px) {
  h1 {
    font-size: 2rem;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .feature-card {
    padding: 1rem;
  }
}
</style>
