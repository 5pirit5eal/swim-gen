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
        /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/,
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
    if (!currentDrill.value?.img_name) return ''
    // Images are in data/images/ folder
    return `/images/drills/${currentDrill.value.img_name}`
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
                <!-- Header Section -->
                <section class="drill-header">
                    <div class="drill-image-container">
                        <img :src="imageUrl" :alt="currentDrill.img_description" class="drill-image"
                            @error="($event.target as HTMLImageElement).style.display = 'none'" />
                    </div>
                    <div class="drill-title-section">
                        <h1 class="drill-title">{{ currentDrill.title }}</h1>
                        <p class="drill-short-description">{{ currentDrill.short_description }}</p>

                        <!-- Tags -->
                        <div class="drill-tags">
                            <span class="tag difficulty" :class="currentDrill.difficulty.toLowerCase()">
                                {{ currentDrill.difficulty }}
                            </span>
                            <span v-for="style in currentDrill.styles" :key="style" class="tag style">
                                {{ style }}
                            </span>
                        </div>
                    </div>
                </section>

                <!-- Description Section -->
                <section class="drill-description">
                    <h2>{{ t('drill.description') }}</h2>
                    <div class="description-content">
                        <p v-for="(paragraph, index) in currentDrill.description" :key="index">
                            {{ paragraph }}
                        </p>
                    </div>
                </section>

                <!-- Targets Section -->
                <section class="drill-targets">
                    <h2>{{ t('drill.targets') }}</h2>
                    <div class="target-list">
                        <span v-for="target in currentDrill.targets" :key="target" class="tag target">
                            {{ target }}
                        </span>
                    </div>
                </section>

                <!-- Target Groups Section -->
                <section class="drill-target-groups">
                    <h2>{{ t('drill.target_groups') }}</h2>
                    <div class="target-group-list">
                        <span v-for="group in currentDrill.target_groups" :key="group" class="tag target-group">
                            {{ group }}
                        </span>
                    </div>
                </section>

                <!-- Video Section -->
                <section v-if="hasVideos" class="drill-videos">
                    <h2>{{ t('drill.video') }}</h2>
                    <div class="video-container">
                        <div v-for="videoId in videoIds" :key="videoId" class="video-wrapper">
                            <iframe :src="`https://www.youtube.com/embed/${videoId}`" frameborder="0"
                                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                                allowfullscreen class="video-iframe"></iframe>
                        </div>
                    </div>
                </section>
            </div>
        </Transition>
    </div>
</template>

<style scoped>
.drill-view {
    padding: 0.25rem 0 2rem 0;
}

.container {
    max-width: 900px;
    margin: 0 auto;
    padding: 0 1rem;
}

/* Loading State */
.loading-state,
.error-state {
    text-align: center;
    padding: 3rem 2rem;
    color: var(--color-text);
    font-style: italic;
    background: var(--color-background-soft);
    border-radius: 8px;
    border: 1px solid var(--color-border);
    margin: 2rem auto;
    max-width: 900px;
}

.error-state {
    color: var(--color-error);
}

.loading-spinner {
    width: 120px;
    height: 40px;
    background-color: var(--color-background-soft);
    position: relative;
    border-radius: 50px;
    box-shadow: inset 0 0 0 2px var(--color-border);
    margin: 0 auto 1rem auto;
}

.loading-spinner:after {
    border-radius: 50px;
    content: '';
    position: absolute;
    background-color: var(--color-primary);
    left: 2px;
    top: 2px;
    bottom: 2px;
    right: 80px;
    animation: slide 2s linear infinite;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

@keyframes slide {
    0% {
        right: 80px;
        left: 2px;
    }

    5% {
        left: 2px;
    }

    50% {
        right: 2px;
        left: 80px;
    }

    55% {
        right: 2px;
    }

    100% {
        right: 80px;
        left: 2px;
    }
}

/* Header Section */
.drill-header {
    display: flex;
    gap: 2rem;
    margin: 2rem 0;
    background: var(--color-background-soft);
    border-radius: 12px;
    border: 1px solid var(--color-border);
    padding: 1.5rem;
}

.drill-image-container {
    flex-shrink: 0;
    width: 280px;
    height: 280px;
    border-radius: 8px;
    overflow: hidden;
    background: var(--color-background-mute);
    display: flex;
    align-items: center;
    justify-content: center;
}

.drill-image {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.drill-title-section {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.drill-title {
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-heading);
    margin: 0 0 0.75rem 0;
    line-height: 1.2;
}

.drill-short-description {
    font-size: 1.1rem;
    color: var(--color-text);
    margin: 0 0 1rem 0;
    line-height: 1.5;
}

.drill-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.tag {
    display: inline-block;
    padding: 0.35rem 0.75rem;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 500;
}

.tag.difficulty {
    color: white;
}

.tag.difficulty.easy {
    background-color: var(--color-success, #22c55e);
}

.tag.difficulty.medium {
    background-color: var(--color-warning, #f59e0b);
}

.tag.difficulty.hard {
    background-color: var(--color-error, #ef4444);
}

/* German translations for difficulty */
.tag.difficulty.leicht {
    background-color: var(--color-success, #22c55e);
}

.tag.difficulty.mittel {
    background-color: var(--color-warning, #f59e0b);
}

.tag.difficulty.schwer {
    background-color: var(--color-error, #ef4444);
}

.tag.style {
    background-color: var(--color-primary);
    color: white;
}

.tag.target {
    background-color: var(--color-background-mute);
    color: var(--color-text);
    border: 1px solid var(--color-border);
}

.tag.target-group {
    background-color: var(--color-background);
    color: var(--color-text);
    border: 1px solid var(--color-primary);
}

/* Sections */
section {
    margin: 2rem 0;
}

section h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--color-heading);
    margin: 0 0 1rem 0;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid var(--color-border);
}

/* Description */
.description-content {
    background: var(--color-background-soft);
    border-radius: 8px;
    padding: 1.25rem;
    border: 1px solid var(--color-border);
}

.description-content p {
    margin: 0 0 1rem 0;
    line-height: 1.7;
    color: var(--color-text);
}

.description-content p:last-child {
    margin-bottom: 0;
}

/* Targets & Target Groups */
.target-list,
.target-group-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

/* Videos */
.video-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.video-wrapper {
    position: relative;
    padding-bottom: 56.25%;
    height: 0;
    overflow: hidden;
    border-radius: 8px;
    background: var(--color-background-mute);
}

.video-iframe {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
}

/* Responsive */
@media (max-width: 740px) {
    .drill-header {
        flex-direction: column;
        align-items: center;
        text-align: center;
    }

    .drill-image-container {
        width: 200px;
        height: 200px;
    }

    .drill-title {
        font-size: 1.5rem;
    }

    .drill-tags {
        justify-content: center;
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
