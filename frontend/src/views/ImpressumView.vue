<script setup lang="ts">
const VITE_IMPRESSUM_NAME = import.meta.env.VITE_IMPRESSUM_NAME
const VITE_IMPRESSUM_ADDRESS = import.meta.env.VITE_IMPRESSUM_ADDRESS
const VITE_IMPRESSUM_CITY = import.meta.env.VITE_IMPRESSUM_CITY
const VITE_IMPRESSUM_PHONE = import.meta.env.VITE_IMPRESSUM_PHONE
const VITE_IMPRESSUM_EMAIL = import.meta.env.VITE_IMPRESSUM_EMAIL
</script>

<template>
  <div class="legal-page">
    <div class="container">
      <div class="impressum-container">
        <h1>{{ $t('impressum.title') }}</h1>
        <section>
          <p>
            {{ VITE_IMPRESSUM_NAME }}<br />
            {{ VITE_IMPRESSUM_ADDRESS }}<br />
            {{ VITE_IMPRESSUM_CITY }}
          </p>
        </section>
        <section>
          <h3>{{ $t('impressum.contact_title') }}</h3>
          <p>
            {{ $t('impressum.phone') }}: {{ VITE_IMPRESSUM_PHONE }}<br />
            {{ $t('impressum.email') }}: {{ VITE_IMPRESSUM_EMAIL }}
          </p>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.legal-page {
  padding: 0.25rem 0 2rem 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

.impressum-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--color-heading);
}

h3 {
  color: var(--color-heading);
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  font-weight: 600;
}

p {
  margin-bottom: 1rem;
}

small {
  font-size: 0.8rem;
  color: var(--color-text-light);
}
</style>
