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
    return `/images/drills/${preview.value.img_name}`
})

// Difficulty class for styling
const difficultyClass = computed(() => {
    if (!preview.value?.difficulty) return ''
    return preview.value.difficulty.toLowerCase()
})

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
            {{ text || drillId }}
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
                    </div>
                    <div class="card-content">
                        <h4 class="card-title">{{ preview.title }}</h4>
                        <p class="card-description">{{ preview.short_description }}</p>
                        <span class="card-difficulty" :class="difficultyClass">
                            {{ preview.difficulty }}
                        </span>
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
    text-underline-offset: 2px;
    cursor: pointer;
    transition: color 0.2s;
}

.drill-link:hover {
    color: var(--color-primary-hover);
    text-decoration-style: solid;
}

.drill-preview-card {
    position: absolute;
    z-index: 1000;
    top: calc(100% + 8px);
    left: 0;
    width: 300px;
    background: var(--color-background);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    box-shadow:
        0 4px 6px -1px rgba(0, 0, 0, 0.1),
        0 2px 4px -1px rgba(0, 0, 0, 0.06);
    overflow: hidden;
    pointer-events: none;
}

.drill-preview-card.position-top {
    top: auto;
    bottom: calc(100% + 8px);
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
    padding: 2rem;
    color: var(--color-text-mute);
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
    height: 140px;
    background: var(--color-background-mute);
    overflow: hidden;
}

.card-image {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.card-content {
    padding: 0.75rem 1rem;
}

.card-title {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--color-heading);
    margin: 0 0 0.5rem 0;
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.card-description {
    font-size: 0.85rem;
    color: var(--color-text);
    margin: 0 0 0.75rem 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.card-difficulty {
    display: inline-block;
    padding: 0.2rem 0.5rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 500;
    color: white;
}

.card-difficulty.easy,
.card-difficulty.leicht {
    background-color: var(--color-success, #22c55e);
}

.card-difficulty.medium,
.card-difficulty.mittel {
    background-color: var(--color-warning, #f59e0b);
}

.card-difficulty.hard,
.card-difficulty.schwer {
    background-color: var(--color-error, #ef4444);
}

/* Card transitions */
.card-enter-active {
    transition: all 0.2s ease-out;
}

.card-leave-active {
    transition: all 0.15s ease-in;
}

.card-enter-from {
    opacity: 0;
    transform: translateY(-4px) scale(0.98);
}

.card-leave-to {
    opacity: 0;
    transform: translateY(-4px) scale(0.98);
}

.position-top .card-enter-from,
.position-top .card-leave-to {
    transform: translateY(4px) scale(0.98);
}
</style>
