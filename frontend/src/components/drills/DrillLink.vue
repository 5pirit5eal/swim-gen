<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDrillsStore } from '@/stores/drills'
import type { DrillPreview } from '@/types'

const props = defineProps<{
    drillId: string
    text?: string
}>()

const { locale } = useI18n()
const router = useRouter()
const drillsStore = useDrillsStore()

const isHovering = ref(false)
const isLoading = ref(false)
const preview = ref<DrillPreview | null>(null)
const hoverTimeout = ref<ReturnType<typeof setTimeout> | null>(null)
const cardPosition = ref<{ top: boolean; left: boolean }>({ top: false, left: false })
const linkRef = ref<HTMLElement | null>(null)

// Image URL
const imageUrl = computed(() => {
    if (!preview.value?.img_name) return ''
    return `https://storage.googleapis.com/${import.meta.env.VITE_PUBLIC_BUCKET_NAME}/${preview.value.img_name}`
})

function getDifficultyLevel(difficulty: string): number {
    const d = difficulty.toLowerCase()
    if (d === 'easy' || d === 'leicht') return 1
    if (d === 'medium' || d === 'mittel') return 2
    if (d === 'hard' || d === 'schwer') return 3
    return 1
}

async function handleMouseEnter(event: MouseEvent) {
    // Delay before showing the card to avoid flickering
    hoverTimeout.value = setTimeout(async () => {
        isHovering.value = true
        calculatePosition(event)

        if (!preview.value) {
            isLoading.value = true
            const result = await drillsStore.fetchDrillPreview(props.drillId, locale.value)
            // Only update state if still hovering (prevents race condition)
            if (isHovering.value) {
                if (result) {
                    preview.value = result
                }
                isLoading.value = false
            }
        }
    }, 200)
}

function handleMouseLeave() {
    if (hoverTimeout.value) {
        clearTimeout(hoverTimeout.value)
        hoverTimeout.value = null
    }
    isHovering.value = false
}

function calculatePosition(event: MouseEvent) {
    const element = event.target as HTMLElement
    const rect = element.getBoundingClientRect()
    const windowHeight = window.innerHeight
    const windowWidth = window.innerWidth

    // Check if card should appear above or below
    cardPosition.value.top = rect.bottom + 300 > windowHeight

    // Check if card should appear to the left or right
    cardPosition.value.left = rect.left + 320 > windowWidth
}

function navigateToDrill() {
    router.push({ name: 'drill', params: { id: props.drillId } })
}

onUnmounted(() => {
    if (hoverTimeout.value) {
        clearTimeout(hoverTimeout.value)
    }
})
</script>

<template>
    <span class="drill-link-wrapper" ref="linkRef" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
        <a class="drill-link" @click.prevent="navigateToDrill" href="#">
            {{ text }}
        </a>

        <Transition name="card">
            <div v-if="isHovering" class="drill-preview-card"
                :class="{ 'position-top': cardPosition.top, 'position-left': cardPosition.left }">
                <div v-if="isLoading" class="card-loading">
                    <div class="loading-spinner-small"></div>
                </div>
                <template v-else-if="preview">
                    <div class="card-image-container">
                        <img :src="imageUrl" :alt="preview.title" class="card-image"
                            @error="($event.target as HTMLImageElement).style.display = 'none'" />

                        <!-- Top Left: Target -->
                        <span v-if="preview.target" class="image-overlay-badge">{{ preview.target }}</span>

                        <!-- Bottom Right: Difficulty -->
                        <div class="image-overlay-difficulty">
                            <span class="difficulty-text">{{ preview.difficulty }}</span>
                            <div class="difficulty-dots">
                                <span v-for="i in 3" :key="i" class="difficulty-dot"
                                    :class="{ active: i <= getDifficultyLevel(preview.difficulty) }"></span>
                            </div>
                        </div>
                    </div>

                    <div class="card-content">
                        <h4 class="card-title">{{ preview.title }}</h4>
                        <p class="card-description">{{ preview.short_description }}</p>
                    </div>
                </template>
                <div v-else class="card-error">
                    <p>Unable to load preview</p>
                </div>
            </div>
        </Transition>
    </span>
</template>

<style scoped>
.drill-link-wrapper {
    position: relative;
    display: inline;
}

.drill-link {
    color: var(--color-primary);
    text-decoration: underline;
    text-decoration-style: dotted;
    text-underline-offset: 3px;
    cursor: pointer;
    transition: all 0.2s;
    font-weight: 500;
}

.drill-link:hover {
    color: var(--color-primary-hover);
    text-decoration-style: solid;
    background: rgba(59, 130, 246, 0.1);
    border-radius: 4px;
}

.drill-preview-card {
    position: absolute;
    z-index: 1000;
    top: calc(100% + 12px);
    left: 0;
    width: 320px;
    background: var(--color-background);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(0, 0, 0, 0.05);
    overflow: hidden;
    pointer-events: none;
}

.drill-preview-card.position-top {
    top: auto;
    bottom: calc(100% + 12px);
}

.drill-preview-card.position-left {
    left: auto;
    right: 0;
}

.card-loading,
.card-error {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2.5rem;
    color: var(--color-text);
}

.loading-spinner-small {
    width: 24px;
    height: 24px;
    border: 2px solid var(--color-border);
    border-top-color: var(--color-primary);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}

.card-image-container {
    width: 100%;
    aspect-ratio: 16/9;
    background: var(--color-background-mute);
    position: relative;
    overflow: hidden;
}

.card-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
}

.drill-link-wrapper:hover .card-image {
    transform: scale(1.05);
}

/* Image Overlays */
.image-overlay-badge {
    position: absolute;
    top: 12px;
    left: 12px;
    background-color: var(--color-primary);
    color: white;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    padding: 4px 8px;
    border-radius: 4px;
    letter-spacing: 0.05em;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.image-overlay-difficulty {
    position: absolute;
    bottom: 12px;
    right: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--color-transparent);
    backdrop-filter: blur(2px);
    padding: 4px 10px;
    border: 1px solid var(--color-primary-hover);
    border-radius: 20px;
    color: var(--color-heading);
}

.difficulty-text {
    font-size: 0.8rem;
    font-weight: 600;
}

.difficulty-dots {
    display: flex;
    gap: 3px;
}

.difficulty-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: var(--color-background);
}

.difficulty-dot.active {
    background-color: var(--color-primary);
    box-shadow: 0 0 4px var(--color-primary);
}

/* Card Content */
.card-content {
    padding: 1rem 1.25rem;
    background: var(--color-background);
}

.card-title {
    font-size: 1.15rem;
    font-weight: 700;
    color: var(--color-heading);
    margin: 0 0 0.5rem 0;
    line-height: 1.3;
}

.card-description {
    font-size: 0.9rem;
    color: var(--color-text);
    margin: 0;
    line-height: 1.5;
    opacity: 0.8;
}

/* Transitions */
.card-enter-active {
    transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.card-leave-active {
    transition: all 0.15s ease-in;
}

.card-enter-from,
.card-leave-to {
    opacity: 0;
    transform: translateY(8px) scale(0.96);
}

.position-top .card-enter-from,
.position-top .card-leave-to {
    transform: translateY(-8px) scale(0.96);
}
</style>
