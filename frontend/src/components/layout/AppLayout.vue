<script setup lang="ts">
// import AppHeader from './AppHeader.vue' // Necessary for v2
import AppFooter from './AppFooter.vue'
</script>

<template>
  <div class="app-layout">
    <main class="main-content">
      <router-view />
    </main>
    <AppFooter />
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  /* Background image setup */
  background-attachment: fixed;
  background-position: center center;
  background-repeat: no-repeat;
  background-size: cover;
  transition: background-image 0.3s ease;
  /* Ensure the layout sits above the body background */
  background-color: transparent;
}

/* Light mode background */
@media (prefers-color-scheme: light) {
  .app-layout {
    background-image: url('@/assets/light_mode.png');
  }
}

/* Dark mode backgrounds */
@media (prefers-color-scheme: dark) {
  .app-layout {
    background-image: url('@/assets/dark_mode.png');
  }
}

.main-content {
  flex: 1;
  container-type: inherit;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

@media (min-width: 1024px) {
  .main-content {
    padding: 2rem;
  }
}

/* Responsive background adjustments */
@media (max-width: 768px) {
  .app-layout {
    background-attachment: scroll;
    /* Better performance on mobile */
    background-size: cover;
    zoom: 0.75;
  }
}
</style>
