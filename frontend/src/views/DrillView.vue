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

// Extract YouTube video ID from URL
function getYouTubeVideoId(url: string): string | null {
  if (!url) return null
  const match = url.match(
    /(?:youtube\.com\/(?:[^/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?/\s]{11})/,
  )
  return match?.[1] ?? null
}

const videoIds = computed(() => {
  if (!currentDrill.value?.video_url) return []
  return currentDrill.value.video_url
    .map((url) => getYouTubeVideoId(url))
    .filter((id): id is string => id !== null)
})

const hasVideos = computed(() => videoIds.value.length > 0)

// Image URL - assumes images are served from a static path
const imageUrl = computed(() => {
  if (!currentDrill.value?.img_name) {
    console.debug('No image name for current drill')
    return ''
  }
  console.debug('Current drill image:',
    `https://storage.googleapis.com/${import.meta.env.VITE_PUBLIC_BUCKET_NAME}/${currentDrill.value.img_name}`
  )
  return `https://storage.googleapis.com/${import.meta.env.VITE_PUBLIC_BUCKET_NAME}/${currentDrill.value.img_name}`

})

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
    <!-- Loading State -->
    <Transition name="fade">
      <div v-if="isLoading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>{{ t('drill.loading') }}</p>
      </div>
    </Transition>

    <!-- Error State -->
    <Transition name="fade">
      <div v-if="error && !isLoading" class="error-state">
        <p>{{ error }}</p>
      </div>
    </Transition>

    <!-- Drill Content -->
    <Transition name="fade">
      <div v-if="currentDrill && !isLoading" class="container">
        <!-- Main Card -->
        <article class="drill-card">
          <!-- Header Section -->
          <header class="drill-header">
            <div class="header-content">
              <div class="drill-image-container">
                <img :src="imageUrl" :alt="currentDrill.img_description" class="drill-image"
                  @error="($event.target as HTMLImageElement).style.display = 'none'" />
              </div>
              <div class="drill-info">
                <div class="title-row">
                  <h1 class="drill-title">{{ currentDrill.title }}</h1>
                  <!-- Difficulty Badge -->
                  <span class="difficulty-badge" :class="currentDrill.difficulty.toLowerCase()">
                    {{ currentDrill.difficulty }}
                  </span>
                </div>

                <p class="drill-short-description">{{ currentDrill.short_description }}</p>

                <!-- Tags Row -->
                <div class="drill-tags">
                  <span v-for="style in currentDrill.styles" :key="style" class="meta-tag style">
                    {{ style }}
                  </span>
                  <span v-for="group in currentDrill.target_groups" :key="group" class="meta-tag group">
                    {{ group }}
                  </span>
                </div>
              </div>
            </div>
          </header>

          <div class="drill-body">
            <!-- Description Section -->
            <section class="content-section">
              <h2>{{ t('drill.description') }}</h2>
              <div class="description-text">
                <p v-for="(paragraph, index) in currentDrill.description" :key="index">
                  {{ paragraph }}
                </p>
              </div>
            </section>

            <!-- Targets Section -->
            <section class="content-section">
              <h2>{{ t('drill.targets') }}</h2>
              <div class="targets-list">
                <span v-for="target in currentDrill.targets" :key="target" class="target-chip">
                  <span class="check-icon">✓</span>
                  {{ target }}
                </span>
              </div>
            </section>

            <!-- Video Section -->
            <section v-if="hasVideos" class="content-section video-section">
              <h2>{{ t('drill.video') }}</h2>
              <div class="video-grid">
                <div v-for="videoId in videoIds" :key="videoId" class="video-wrapper">
                  <iframe :src="`https://www.youtube.com/embed/${videoId}`" allow="
                      accelerometer;
                      autoplay;
                      clipboard-write;
                      encrypted-media;
                      gyroscope;
                      picture-in-picture;
                    " allowfullscreen class="video-iframe"></iframe>
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
  padding: 2rem 0;
  min-height: 80vh;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* Card Container */
.drill-card {
  background: var(--color-background);
  border-radius: 8px;
  box-shadow: 0 4px 20px var(--color-shadow);
  overflow: hidden;
  border: 1px solid var(--color-border);
}

/* Header Styling */
.drill-header {
  background: var(--color-background-soft);
  padding: 2.5rem;
  border-bottom: 1px solid var(--color-border);
}

.header-content {
  display: flex;
  gap: 2.5rem;
  align-items: flex-start;
}

.drill-image-container {
  flex-shrink: 0;
  width: 320px;
  aspect-ratio: 4/3;
  border-radius: 12px;
  overflow: hidden;
  background: var(--color-background-mute);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.drill-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.drill-image:hover {
  transform: scale(1.02);
}

.drill-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding-top: 0.5rem;
}

.title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.drill-title {
  font-size: 2.25rem;
  font-weight: 800;
  color: var(--color-heading);
  line-height: 1.1;
  margin: 0;
  letter-spacing: -0.5px;
}

.drill-short-description {
  font-size: 1.15rem;
  color: var(--color-text);
  line-height: 1.6;
  margin: 0 0 1.5rem 0;
  opacity: 0.9;
  max-width: 65ch;
}

/* Badges & Tags */
.difficulty-badge {
  padding: 0.4rem 1rem;
  border-radius: 2rem;
  font-size: 0.85rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: white;
  white-space: nowrap;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.difficulty-badge.easy,
.difficulty-badge.leicht {
  background-color: var(--color-success);
}

.difficulty-badge.medium,
.difficulty-badge.mittel {
  background-color: var(--color-warning);
}

.difficulty-badge.hard,
.difficulty-badge.schwer {
  background-color: var(--color-error);
}

.drill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.meta-tag {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.85rem;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.2s;
}

.meta-tag.style {
  background-color: var(--color-primary);
  color: white;
}

.meta-tag.group {
  background-color: var(--color-background);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
}

/* Body Content */
.drill-body {
  padding: 2.5rem;
}

.content-section {
  margin-bottom: 3rem;
}

.content-section:last-child {
  margin-bottom: 0;
}

.content-section h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin: 0 0 1.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.content-section h2::before {
  content: '';
  display: block;
  width: 4px;
  height: 24px;
  background: var(--color-primary);
  border-radius: 2px;
}

.description-text p {
  font-size: 1.05rem;
  line-height: 1.8;
  color: var(--color-text);
  margin-bottom: 1.5rem;
  max-width: 75ch;
}

.description-text p:last-child {
  margin-bottom: 0;
}

/* Targets List */
.targets-list {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.target-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--color-background-soft);
  padding: 0.6rem 1rem;
  border-radius: 8px;
  font-size: 0.95rem;
  color: var(--color-heading);
  border: 1px solid var(--color-border);
}

.check-icon {
  color: var(--color-primary);
  font-weight: bold;
}

/* Video Section */
.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
}

.video-wrapper {
  position: relative;
  padding-bottom: 56.25%;
  /* 16:9 Aspect Ratio */
  height: 0;
  overflow: hidden;
  border-radius: 12px;
  background: var(--color-background-mute);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--color-border);
}

.video-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

/* Loading & Error States */
.loading-state,
.error-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--color-text);
  background: var(--color-background-soft);
  border-radius: 12px;
  border: 1px dashed var(--color-border);
  margin: 2rem auto;
  max-width: 600px;
}

.error-state {
  color: var(--color-error);
  border-color: var(--color-error-soft);
  background: rgba(231, 76, 60, 0.05);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  margin: 0 auto 1.5rem auto;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Responsive Adjustments */
@media (max-width: 850px) {
  .header-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .drill-image-container {
    width: 100%;
    max-width: 400px;
  }

  .title-row {
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  .drill-title {
    font-size: 1.75rem;
  }

  .drill-tags {
    justify-content: center;
  }

  .content-section h2::before {
    display: none;
  }

  .content-section h2 {
    justify-content: center;
    border-bottom: 2px solid var(--color-border);
    padding-bottom: 0.5rem;
  }

  .video-grid {
    grid-template-columns: 1fr;
  }
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
