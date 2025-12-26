<script setup lang="ts">
import { useDrillsStore } from '@/stores/drills'
import { storeToRefs } from 'pinia'
import { onMounted, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const drillsStore = useDrillsStore()

const { currentDrill, isLoading, error } = storeToRefs(drillsStore)

const imageUrl = computed(() => {
  if (currentDrill.value?.img_name) {
    const bucketName = import.meta.env.VITE_PUBLIC_BUCKET_NAME || 'swim-gen-public'
    return `https://storage.googleapis.com/${bucketName}/${currentDrill.value.img_name}`
  }
  return '/img/placeholder.webp' // Default placeholder image
})

const hasValidVideo = computed(() => {
  if (!currentDrill.value?.video_url) return false
  return currentDrill.value.video_url.some((url) => url && url.trim().length > 0)
})

function getYouTubeId(url: string): string | null {
  const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/
  const match = url.match(regExp)
  return match && match[2].length === 11 ? match[2] : null
}

function getDifficultyLevel(difficulty: string): number {
  const d = difficulty.toLowerCase()
  if (d === 'easy' || d === 'leicht') return 1
  if (d === 'medium' || d === 'mittel') return 2
  if (d === 'hard' || d === 'schwer') return 3
  return 1
}

async function initializeView() {
  const drillId = route.params.id
  if (typeof drillId === 'string') {
    const drill = await drillsStore.fetchDrill(drillId, locale.value)
    if (!drill) {
      noDrillFound()
    }
  } else {
    noDrillFound()
  }
}

function noDrillFound() {
  toast.error(t('drill.not_found'))
  router.push('/')
}

onMounted(async () => {
  window.scrollTo(0, 0)
  await initializeView()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await initializeView()
    }
  },
)

// Refetch when locale changes
watch(locale, async () => {
  const drillId = route.params.id
  if (typeof drillId === 'string') {
    await drillsStore.fetchDrill(drillId, locale.value)
  }
})
</script>

<template>
  <div class="drill-view">
    <Transition name="fade">
      <!-- Loading State -->
      <div v-if="isLoading" class="loading-state">
        <div class="loading-spinner"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <p>{{ error }}</p>
      </div>

      <!-- Main Content -->
      <div v-else-if="currentDrill" class="container">
        <article class="drill-card">
          <header class="drill-header">
            <div class="header-content">
              <div class="image-container">
                <img :src="imageUrl" :alt="currentDrill.img_description" class="drill-image" />
              </div>
              <div class="header-details">
                <div class="title-row">
                  <h1 class="drill-title">{{ currentDrill.title }}</h1>
                  <div class="difficulty-indicator">
                    <span v-for="i in 3" :key="i" class="difficulty-dot"
                      :class="{ active: i <= getDifficultyLevel(currentDrill.difficulty) }"></span>
                  </div>
                </div>

                <div class="meta-row">
                  <div class="meta-group">
                    <span class="meta-label">{{ t('drill.styles') }}:</span>
                    <span v-for="style in currentDrill.styles" :key="style" class="meta-value style">{{ style }}</span>
                  </div>
                  <div class="meta-group">
                    <span class="meta-label">{{ t('drill.targets') }}:</span>
                    <span v-for="target in currentDrill.targets" :key="target" class="meta-value target">{{ target
                    }}</span>
                  </div>
                </div>

                <p class="drill-short-description">{{ currentDrill.short_description }}</p>
              </div>
            </div>
          </header>

          <div class="drill-body">
            <section class="description-section">
              <div class="section-header">
                <h3>{{ t('drill.description') }}</h3>
              </div>
              <div class="description-text">
                <p v-for="(paragraph, index) in currentDrill.description" :key="index">
                  {{ paragraph }}
                </p>
              </div>
            </section>

            <section v-if="currentDrill.video_url?.length && hasValidVideo" class="video-section">
              <div class="section-header">
                <h3>{{ t('drill.videos') }}</h3>
              </div>
              <div class="video-grid">
                <div v-for="(url, index) in currentDrill.video_url" :key="index" class="video-wrapper">
                  <iframe v-if="getYouTubeId(url)" :src="`https://www.youtube.com/embed/${getYouTubeId(url)}`"
                    title="YouTube video player" frameborder="0"
                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                    allowfullscreen class="video-iframe"></iframe>
                </div>
              </div>
            </section>

            <section class="meta-footer">
              <div class="meta-group">
                <span class="meta-label">{{ t('drill.target_groups') }}:</span>
                <div class="tags-container">
                  <span v-for="group in currentDrill.target_groups" :key="group" class="meta-tag group">
                    {{ group }}
                  </span>
                </div>
              </div>
            </section>
          </div>
        </article>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.drill-view {
  min-height: 100vh;
  padding: 2rem 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* Loading & Error States */
.loading-state,
.error-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-state {
  color: var(--color-error);
  font-size: 1.1rem;
}

/* Card Container */
.drill-card {
  background: var(--color-background);
  border-radius: 8px;
  box-shadow: 0 4px 20px var(--color-shadow);
  overflow: hidden;
  border: 1px solid var(--color-border);
}

/* Header Section */
.drill-header {
  background: var(--color-background-soft);
  padding: 2.5rem;
  border-bottom: 1px solid var(--color-border);
}

.header-content {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 2.5rem;
  align-items: start;
}

.image-container {
  aspect-ratio: 16/10;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background: var(--color-background-mute);
}

.drill-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.drill-image:hover {
  transform: scale(1.02);
}

.header-details {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.drill-title {
  font-size: 2rem;
  font-weight: 800;
  color: var(--color-heading);
  line-height: 1.1;
  margin: 0;
  letter-spacing: -0.01em;
}

/* Difficulty Dots */
.difficulty-indicator {
  display: flex;
  gap: 4px;
  margin-top: 8px;
  padding: 6px 10px;
  background: var(--color-background-mute);
  border-radius: 20px;
}

.difficulty-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: var(--color-border);
  transition: background-color 0.3s ease;
}

.difficulty-dot.active {
  background-color: var(--color-primary);
  box-shadow: 0 0 4px var(--color-primary-hover);
}

.drill-short-description {
  font-size: 1.1rem;
  line-height: 1.5;
  color: var(--color-text);
  margin-top: 0.5rem;
  max-width: 60ch;
}

/* Metadata Styling */
.meta-row {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.meta-group {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.75rem;
}

.meta-label {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-text);
  opacity: 0.8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  min-width: 80px;
}

.meta-value {
  font-size: 0.95rem;
  font-weight: 500;
}

.meta-value.style {
  color: var(--color-primary);
}

.meta-value.target {
  color: var(--color-heading);
  background: var(--color-background-mute);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.9rem;
}

/* Body Content */
.drill-body {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.section-header {
  margin-bottom: 1.5rem;
  border-left: 4px solid var(--color-primary);
  padding-left: 1rem;
}

.section-header h3 {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--color-heading);
  margin: 0;
}

.description-text {
  font-size: 1rem;
  line-height: 1.6;
  color: var(--color-text);
  max-width: 75ch;
}

.description-text p {
  margin-bottom: 1rem;
}

.description-text p:last-child {
  margin-bottom: 0;
}

/* Video Section */
.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.video-wrapper {
  position: relative;
  padding-bottom: 56.25%;
  /* 16:9 Aspect Ratio */
  height: 0;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background: var(--color-background-soft);
}

.video-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

/* Meta Footer for Target Groups */
.meta-footer {
  padding-top: 2rem;
  border-top: 1px solid var(--color-border);
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.meta-tag.group {
  font-size: 0.85rem;
  padding: 0.35rem 0.85rem;
  border-radius: 6px;
  background: var(--color-background-soft);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  font-weight: 500;
}

/* Responsive */
@media (max-width: 900px) {
  .header-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .image-container {
    width: 100%;
  }

  .drill-body {
    padding: 2rem;
    gap: 2.5rem;
  }

  .video-grid {
    grid-template-columns: 1fr;
  }
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
