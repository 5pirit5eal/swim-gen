<script setup lang="ts">
const VITE_IMPRESSUM_NAME = import.meta.env.VITE_IMPRESSUM_NAME;
const VITE_IMPRESSUM_ADDRESS = import.meta.env.VITE_IMPRESSUM_ADDRESS;
const VITE_IMPRESSUM_CITY = import.meta.env.VITE_IMPRESSUM_CITY;
const VITE_IMPRESSUM_PHONE = import.meta.env.VITE_IMPRESSUM_PHONE;
const VITE_IMPRESSUM_EMAIL = import.meta.env.VITE_IMPRESSUM_EMAIL;
</script>

<template>
  <div class="impressum-container">
    <section>
      <h3>{{ $t('impressum.title') }}</h3>
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
</template>

<style scoped>
.impressum-container {
  padding: 1rem;
  line-height: 1.6;
  color: var(--color-text);
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
