<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import type { Drill } from '@/types'
import IconArrowRight from '@/components/icons/IconArrowRight.vue'

const props = defineProps<{
  drill: Drill
}>()

const router = useRouter()
const { t } = useI18n()

const imageUrl = computed(() => {
  if (!props.drill.img_name) return ''
  return `https://storage.googleapis.com/${import.meta.env.VITE_PUBLIC_BUCKET_NAME}/${props.drill.img_name}`
})

function getDifficultyLevel(difficulty: string): number {
  if (!difficulty) return 1
  const d = difficulty.toLowerCase()
  if (d === 'easy' || d === 'leicht') return 1
  if (d === 'medium' || d === 'mittel') return 2
  if (d === 'hard' || d === 'schwer') return 3
  return 1
}

function navigateToDrill() {
  router.push({ name: 'drill', params: { id: props.drill.img_name } }) // Using img_name as ID because that's what backend expects currently in search response but usually ID is what we want.
  // Wait, backend `SearchDrills` returns `models.Drill`. The frontend `Drill` interface has `img_name` but also `slug`?
  // `GetDrillHandler` takes `id` which corresponds to `img_name` in DB.
  // So `img_name` is effectively the ID.
}

</script>

<template>
  <article class="drill-card" @click="navigateToDrill">
    <div class="card-image-container">
      <img :src="imageUrl" :alt="drill.title" class="card-image" loading="lazy"
        @error="($event.target as HTMLImageElement).style.display = 'none'" />

      <!-- Top Left: Target or Style Badge -->
      <span v-if="drill.targets && drill.targets.length > 0" class="image-overlay-badge target">
        {{ drill.targets[0] }}
      </span>
      <span v-else-if="drill.styles && drill.styles.length > 0" class="image-overlay-badge style">
        {{ drill.styles[0] }}
      </span>
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ drill.title }}</h3>
      <p class="card-description">{{ drill.short_description }}</p>

      <div class="card-footer">
        <div class="difficulty-info">
          <!-- <IconClock class="meta-icon" /> -->
          <!-- We don't have duration in data yet, skipping -->
          <span class="difficulty-text">{{ drill.difficulty }}</span>
          <div class="difficulty-dots">
            <span v-for="i in 3" :key="i" class="difficulty-dot"
              :class="{ active: i <= getDifficultyLevel(drill.difficulty) }"></span>
          </div>
        </div>

        <button class="action-button" :aria-label="t('drill.view_details', 'View Details')">
          <IconArrowRight class="icon-arrow" />
        </button>
      </div>
    </div>
  </article>
</template>

<style scoped>
.drill-card {
  background: var(--color-background-soft);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--color-border);
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.drill-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px -10px var(--color-shadow);
  border-color: var(--color-primary);
}

.card-image-container {
  aspect-ratio: 16/9;
  background: var(--color-background-mute);
  position: relative;
  overflow: hidden;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.drill-card:hover .card-image {
  transform: scale(1.05);
}

.image-overlay-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  padding: 4px 8px;
  border-radius: 6px;
  letter-spacing: 0.05em;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  z-index: 2;
}

.image-overlay-badge.target {
  background-color: var(--color-primary);
  color: white;
}

.image-overlay-badge.style {
  background-color: var(--color-secondary, #64748b);
  /* Fallback color */
  color: white;
}

.card-content {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--color-heading);
  margin: 0 0 0.5rem 0;
  line-height: 1.3;
}

.card-description {
  font-size: 0.9rem;
  color: var(--color-text);
  margin: 0 0 1.5rem 0;
  line-height: 1.5;
  opacity: 0.85;
  flex: 1;
  /* Pushes footer down */
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.difficulty-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text);
  font-size: 0.85rem;
  font-weight: 500;
}

.difficulty-dots {
  display: flex;
  gap: 3px;
}

.difficulty-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: var(--color-border);
}

.difficulty-dot.active {
  background-color: var(--color-text);
  /* Use text color for dots as per image dark theme intuition, or stick to primary */
  background-color: color-mix(in srgb, var(--color-primary), var(--color-heading));
  box-shadow: 0 0 2px var(--color-primary);
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-background-mute);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  transition: all 0.2s;
  cursor: pointer;
}

.drill-card:hover .action-button {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.icon-arrow {
  width: 16px;
  height: 16px;
}
</style>
