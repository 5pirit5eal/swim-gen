<script setup lang="ts">
import DrillList from '@/components/drills/DrillList.vue'
import { onMounted } from 'vue'

onMounted(() => {
  window.scrollTo(0, 0)
})
</script>

<template>
  <div class="drills-view">
    <div class="container">
      <DrillList :featured-mode="false" />
    </div>
  </div>
</template>

<style scoped>
.drills-view {
  padding: 1.5rem 0 3rem 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}
</style>
