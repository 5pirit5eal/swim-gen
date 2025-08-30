<script setup lang="ts">
import { useI18n } from 'vue-i18n'
// Header component for the swim training plan generator
// Currently minimal for V1 (single page), expandable for V2

const { t } = useI18n()
</script>

<template>
  <header class="app-header">
    <div class="header-container">
      <div class="logo">
        <h1>{{ t('app.name') }}</h1>
        <!-- <span class="subtitle">Training Plan Generator</span> -->
      </div>

      <!-- Navigation for future use (V2) -->
      <nav class="navigation" v-if="false">
        <router-link to="/" class="nav-link">{{ t('header.home') }}</router-link>
        <router-link to="/about" class="nav-link">{{ t('header.about') }}</router-link>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  background-color: var(--color-background-soft, #f8f9fa);
  border-bottom: 1px solid var(--color-border, #e2e8f0);
  padding: 1rem 0;
}

.header-container {
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  container-type: inline-size;
}

.logo h1 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--color-heading, #2c3e50);
}

.subtitle {
  font-size: 0.875rem;
  color: var(--color-text-light, #64748b);
}

.navigation {
  display: flex;
  gap: 1rem;
}

.nav-link {
  text-decoration: none;
  color: var(--color-text, #374151);
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: var(--color-background-mute, #f1f5f9);
}

.nav-link.router-link-active {
  background-color: var(--color-background-mute, #f1f5f9);
  color: var(--color-text, #1f2937);
}
</style>
