<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseModal from '@/components/ui/BaseModal.vue'

defineProps<{
    show: boolean
    planTitle: string
}>()

const emit = defineEmits<{
    (e: 'close'): void
    (e: 'submit', payload: { rating: number; was_swam: boolean; difficulty_rating: number; comment?: string }): void
}>()

const { t } = useI18n()

const rating = ref(0)
const wasSwam = ref(false)
const difficultyRating = ref(5)
const comment = ref('')
const hoveredStar = ref(0)

const isValid = computed(() => rating.value > 0)

function setRating(star: number) {
    rating.value = star
}

function setHover(star: number) {
    hoveredStar.value = star
}

function clearHover() {
    hoveredStar.value = 0
}

function submit() {
    if (!isValid.value) return
    emit('submit', {
        rating: rating.value,
        was_swam: wasSwam.value,
        difficulty_rating: difficultyRating.value,
        comment: comment.value
    })
    resetForm()
}

function close() {
    resetForm()
    emit('close')
}

function resetForm() {
    rating.value = 0
    wasSwam.value = false
    difficultyRating.value = 5
    comment.value = ''
    hoveredStar.value = 0
}
</script>

<template>
    <BaseModal :show="show" @close="close">
        <template #header>
            <h2>{{ t('feedback.title') }}</h2>
        </template>

        <template #body>
            <div class="intro-text">
                <p class="subtitle">{{ t('feedback.subtitle', { plan: planTitle }) }}</p>
                <p class="thank-you">{{ t('feedback.thank_you') }}</p>
            </div>

            <div class="form-group">
                <label>{{ t('feedback.rating_label') }} <span class="required">*</span></label>
                <div class="stars">
                    <span v-for="star in 5" :key="star" class="star"
                        :class="{ active: star <= (hoveredStar || rating) }" @click="setRating(star)"
                        @mouseenter="setHover(star)" @mouseleave="clearHover">
                        ★
                    </span>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group checkbox-group">
                    <input type="checkbox" id="swam-checkbox" v-model="wasSwam">
                    <label for="swam-checkbox">{{ t('feedback.swam_label') }}</label>
                </div>

                <div class="form-group difficulty-group">
                    <label>{{ t('feedback.difficulty_label') }}: {{ difficultyRating }} <span
                            class="required">*</span></label>
                    <div class="slider-container">
                        <span class="slider-label">{{ t('feedback.easy') }}</span>
                        <input type="range" min="1" max="10" v-model.number="difficultyRating" class="slider">
                        <span class="slider-label">{{ t('feedback.hard') }}</span>
                    </div>
                </div>
            </div>

            <div class="form-group">
                <label for="comment">{{ t('feedback.comment_label') }} ({{ t('common.optional') }})</label>
                <textarea id="comment" v-model="comment" rows="3"
                    :placeholder="t('feedback.comment_placeholder')"></textarea>
            </div>

            <p class="disclaimer">{{ t('feedback.disclaimer') }}</p>
        </template>

        <template #footer>
            <button class="submit-btn" @click="submit" :disabled="!isValid">
                {{ t('common.submit') }}
            </button>
        </template>
    </BaseModal>
</template>

<style scoped>
h2 {
    margin: 0;
    color: var(--color-heading);
}

.intro-text {
    text-align: center;
    margin-bottom: 1.5rem;
}

.subtitle {
    color: var(--color-text);
    margin-bottom: 0.5rem;
    font-size: 0.95rem;
    font-weight: 500;
}

.thank-you {
    color: var(--color-text-soft);
    font-size: 0.85rem;
    margin: 0;
}

.required {
    color: var(--color-error);
    font-weight: bold;
}

.form-group {
    margin-bottom: 1.5rem;
}

.form-row {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 2rem;
    align-items: start;
    margin-bottom: 1.5rem;
}

.checkbox-group {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0;
}

.checkbox-group label {
    margin-bottom: 0;
    cursor: pointer;
    white-space: nowrap;
}

.checkbox-group input[type="checkbox"] {
    cursor: pointer;
}

.difficulty-group {
    margin-bottom: 0;
    flex: 1;
}

label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--color-text);
}

.stars {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    font-size: 2rem;
    cursor: pointer;
}

.star {
    color: var(--color-border);
    transition: color 0.2s;
}

.star.active {
    color: #fbbf24;
}

.slider-container {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.slider {
    flex: 1;
    accent-color: var(--color-primary);
}

.slider-label {
    font-size: 0.8rem;
    color: var(--color-text-soft);
}

textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background-color: var(--color-background);
    color: var(--color-text);
    resize: vertical;
}

.submit-btn {
    background-color: var(--color-primary);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
    background-color: var(--color-primary-hover);
}

.submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.disclaimer {
    text-align: center;
    font-size: 0.8rem;
    color: var(--color-text-soft);
    margin-top: 1rem;
}
</style>
