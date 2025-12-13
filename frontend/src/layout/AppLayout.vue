<script setup lang="ts">
import AppHeader from './AppHeader.vue'
import AppFooter from './AppFooter.vue'
import Sidebar from './AppSidebar.vue'
import { useSidebarStore } from '@/stores/sidebar'

const sidebarStore = useSidebarStore()
</script>

<template>
  <div class="app-layout">
    <Sidebar />
    <div class="content-wrapper" :class="{ 'sidebar-open': sidebarStore.isOpen }">
      <AppHeader />
      <main class="main-content">
        <router-view />
      </main>
      <AppFooter />
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  position: relative;
  /* Background image setup */
  background-attachment: fixed;
  background-position: center center;
  background-repeat: no-repeat;
  background-size: cover;
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

.content-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  flex-grow: 1;
  transition: margin-left 0.3s ease;
}

.content-wrapper.sidebar-open {
  margin-left: 400px;
}

.main-content {
  container-type: inherit;
  position: relative;
  z-index: 1;
}

@media (max-width: 1124px) {
  .app-layout {
    zoom: 0.85;
  }

  .content-wrapper.sidebar-open {
    margin-left: 300px;
  }
}

/* Responsive background adjustments */
@media (max-width: 740px) {
  .app-layout {
    background-attachment: scroll;
    /* Better performance on mobile */
    background-size: cover;
    zoom: 0.75;
  }

  .content-wrapper.sidebar-open {
    margin-left: 0;
  }
}
</style>
