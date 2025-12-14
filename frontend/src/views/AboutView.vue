<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="legal-page">
    <div class="container">
      <div class="content-box">
        <h1>{{ t('about.title') }}</h1>
        <section class="intro">
          <p class="greeting">
            <strong>{{ t('about.greeting') }}</strong>
          </p>
        </section>

        <section>
          <p>{{ t('about.intro_text') }}</p>
          <p>{{ t('about.origin_story') }}</p>
          <p>{{ t('about.mission') }}</p>
          <p>{{ t('about.personal_note') }}</p>

          <div class="signature">
            <p>{{ t('about.cheers') }}</p>
            <p class="name">Ruben</p>
          </div>
        </section>

        <section class="support">
          <p>
            {{ t('about.support_text') }}
            <a href="https://buymeacoffee.com/spiritseal" target="_blank" rel="noopener">
              {{ t('about.support_link') }}
            </a>
          </p>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.legal-page {
  padding: 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

.content-box {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.content-box h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--color-heading);
}

.content-box section {
  margin-bottom: 1.5rem;
}

.content-box p {
  margin-bottom: 1rem;
}

.intro {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--color-background);
  padding: 1rem;
  border-radius: 8px;
}

.greeting {
  font-size: 1.5rem;
  color: var(--color-heading);
  margin-bottom: 0;
  text-align: center;
}

.signature {
  margin-top: 2rem;
  margin-bottom: 2rem;
}

.signature .name {
  font-weight: 600;
  color: var(--color-text);
  font-family: 'Handlee', cursive, sans-serif;
  /* Fallback if font not loaded, though not strictly required by brief */
}

.support {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid var(--color-border);
  text-align: center;
  font-style: italic;
}

.support a {
  color: var(--color-primary);
  text-decoration: underline;
  font-weight: 500;
}

.support a:hover {
  color: var(--color-primary-hover);
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .content-box {
    max-width: 100%;
    padding: 1.5rem;
  }
}
</style>
