<script setup lang="ts">
import { ref, computed, nextTick, onUnmounted } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSettingsStore } from '@/stores/settings'
import { useProfileStore } from '@/stores/profile'
import { apiClient, formatError } from '@/api/client'
import type { QueryRequest, PromptGenerationRequest, AudienceType } from '@/types'
import { DIFFICULTY_OPTIONS, TRAINING_TYPE_OPTIONS } from '@/types'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { calculateOverlayPosition, type OverlayPlacement } from '@/utils/overlayPosition'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

// Store access
const authStore = useAuthStore()
const profileStore = useProfileStore()
const trainingStore = useTrainingPlanStore()
const settingsStore = useSettingsStore()

// i18n
const { t, locale } = useI18n()

// Form data
const requestText = ref('')
const showAdvancedSettings = ref(false)

// Loading state for prompt generation
const generatingPrompt = ref(false)
const highlightPromptBtn = ref(false)
let highlightTransitionTimer: ReturnType<typeof setTimeout> | null = null
let highlightTimer: ReturnType<typeof setTimeout> | null = null

// Audience options configuration sorted by experience: Anfänger, Hobbyschwimmer, Triathlet, Leistungsschwimmer
const audienceOptions: { id: AudienceType; labelKey: string; hintKey: string }[] = [
  { id: 'beginner', labelKey: 'form.audience_beginner', hintKey: 'form.audience_hint_beginner' },
  { id: 'hobby', labelKey: 'form.audience_hobby', hintKey: 'form.audience_hint_hobby' },
  {
    id: 'triathlete',
    labelKey: 'form.audience_triathlete',
    hintKey: 'form.audience_hint_triathlete',
  },
  {
    id: 'competitive_swimmer',
    labelKey: 'form.audience_competitive_swimmer',
    hintKey: 'form.audience_hint_competitive_swimmer',
  },
]

// Hover tooltip state for audience buttons (longer hover)
const hoveredAudience = ref<AudienceType | null>(null)
const activeTooltipId = ref<AudienceType | null>(null)
const tooltipPosition = ref({ left: 0, top: 0 })
const tooltipPlacement = ref<OverlayPlacement>('top')
const audienceButtonRefs = ref<Record<string, HTMLElement | null>>({})
let hoverTimer: ReturnType<typeof setTimeout> | null = null

// Submit button tooltip state
const submitWrapperRef = ref<HTMLElement | null>(null)
const isSubmitHovered = ref(false)
const showSubmitTooltip = ref(false)
const submitTooltipPosition = ref({ left: 0, top: 0 })
const submitTooltipPlacement = ref<OverlayPlacement>('top')
let submitHoverTimer: ReturnType<typeof setTimeout> | null = null

function setAudienceBtnRef(id: AudienceType, el: unknown) {
  if (el instanceof HTMLElement) {
    audienceButtonRefs.value[id] = el
  } else {
    delete audienceButtonRefs.value[id]
  }
}

async function updateTooltipPosition() {
  await nextTick()
  if (!activeTooltipId.value) return
  const anchorEl = audienceButtonRefs.value[activeTooltipId.value]
  const tooltipEl = document.getElementById('audience-hover-tooltip')
  if (!anchorEl || !tooltipEl) return

  const anchor = anchorEl.getBoundingClientRect()
  const tooltip = tooltipEl.getBoundingClientRect()
  const position = calculateOverlayPosition(
    anchor,
    tooltip.width,
    tooltip.height,
    window.innerWidth,
    window.innerHeight,
  )

  tooltipPosition.value = { left: position.left, top: position.top }
  tooltipPlacement.value = position.placement
}

function handleAudienceMouseEnter(id: AudienceType) {
  hoveredAudience.value = id
  if (hoverTimer) {
    clearTimeout(hoverTimer)
  }
  // Show tooltip only on a longer hover (500ms)
  hoverTimer = setTimeout(async () => {
    if (hoveredAudience.value === id) {
      activeTooltipId.value = id
      await updateTooltipPosition()
      window.addEventListener('resize', updateTooltipPosition)
      window.addEventListener('scroll', updateTooltipPosition, true)
    }
  }, 500)
}

function handleAudienceMouseLeave(id: AudienceType) {
  if (hoveredAudience.value === id) {
    hoveredAudience.value = null
    if (hoverTimer) {
      clearTimeout(hoverTimer)
      hoverTimer = null
    }
    activeTooltipId.value = null
    window.removeEventListener('resize', updateTooltipPosition)
    window.removeEventListener('scroll', updateTooltipPosition, true)
  }
}

async function updateSubmitTooltipPosition() {
  await nextTick()
  if (!showSubmitTooltip.value || !submitWrapperRef.value) return
  const anchorEl = submitWrapperRef.value
  const tooltipEl = document.getElementById('submit-disabled-tooltip')
  if (!anchorEl || !tooltipEl) return

  const anchor = anchorEl.getBoundingClientRect()
  const tooltip = tooltipEl.getBoundingClientRect()
  const position = calculateOverlayPosition(
    anchor,
    tooltip.width,
    tooltip.height,
    window.innerWidth,
    window.innerHeight,
  )

  submitTooltipPosition.value = { left: position.left, top: position.top }
  submitTooltipPlacement.value = position.placement
}

function handleSubmitMouseEnter() {
  if (canSubmit.value) return
  isSubmitHovered.value = true
  if (submitHoverTimer) {
    clearTimeout(submitHoverTimer)
  }
  submitHoverTimer = setTimeout(async () => {
    if (isSubmitHovered.value && !canSubmit.value) {
      showSubmitTooltip.value = true
      await updateSubmitTooltipPosition()
      window.addEventListener('resize', updateSubmitTooltipPosition)
      window.addEventListener('scroll', updateSubmitTooltipPosition, true)
    }
  }, 200)
}

function handleSubmitMouseLeave() {
  isSubmitHovered.value = false
  if (submitHoverTimer) {
    clearTimeout(submitHoverTimer)
    submitHoverTimer = null
  }
  showSubmitTooltip.value = false
  window.removeEventListener('resize', updateSubmitTooltipPosition)
  window.removeEventListener('scroll', updateSubmitTooltipPosition, true)
}

onUnmounted(() => {
  if (hoverTimer) {
    clearTimeout(hoverTimer)
  }
  if (submitHoverTimer) {
    clearTimeout(submitHoverTimer)
  }
  if (highlightTransitionTimer) {
    clearTimeout(highlightTransitionTimer)
  }
  if (highlightTimer) {
    clearTimeout(highlightTimer)
  }
  window.removeEventListener('resize', updateTooltipPosition)
  window.removeEventListener('scroll', updateTooltipPosition, true)
  window.removeEventListener('resize', updateSubmitTooltipPosition)
  window.removeEventListener('scroll', updateSubmitTooltipPosition, true)
})

// Dynamic placeholder based on selected audience
const currentPlaceholder = computed(() => {
  switch (settingsStore.selectedAudience) {
    case 'beginner':
      return t('form.example_placeholder_beginner')
    case 'triathlete':
      return t('form.example_placeholder_triathlete')
    case 'competitive_swimmer':
      return t('form.example_placeholder_competitive_swimmer')
    case 'hobby':
      return t('form.example_placeholder_hobby')
    default:
      return t('form.example_placeholder')
  }
})

// Computed
const isFormValid = computed(() => {
  const content = requestText.value.trim()
  return content.length > 0 && content.length <= 3000
})
const tooMuchText = computed(() => requestText.value.trim().length > 3000)
const canSubmit = computed(
  () => isFormValid.value && !trainingStore.isLoading && !generatingPrompt.value,
)

function clearPromptHighlight() {
  if (highlightTransitionTimer) {
    clearTimeout(highlightTransitionTimer)
    highlightTransitionTimer = null
  }
  if (highlightTimer) {
    clearTimeout(highlightTimer)
    highlightTimer = null
  }
  highlightPromptBtn.value = false
}

// Actions
async function handleSubmit() {
  if (!canSubmit.value) return

  const request: QueryRequest = {
    content: requestText.value.trim(),
    method: settingsStore.preferredMethod,
    filter: settingsStore.filters,
    language: navigator.language,
    pool_length: settingsStore.poolLength,
    preferences: settingsStore.useProfilePreferences,
  }

  const success = await trainingStore.generatePlan(request)
  if (!success) {
    trainingStore.error = t('errors.failed_to_generate_plan')
  }
}

function toggleAdvancedSettings() {
  showAdvancedSettings.value = !showAdvancedSettings.value
}

function applyAudienceConfiguration(audience: AudienceType) {
  settingsStore.setAudience(audience)

  if (audience === 'beginner') {
    settingsStore.filters.schwierigkeitsgrad = 'Anfaenger'
    settingsStore.filters.trainingstyp = 'Techniktraining'
  } else if (audience === 'triathlete') {
    settingsStore.filters.schwierigkeitsgrad = 'Fortgeschritten'
    settingsStore.filters.trainingstyp = 'Grundlagenausdauer'
  } else if (audience === 'competitive_swimmer') {
    settingsStore.filters.schwierigkeitsgrad = 'Leistungsschwimmer'
    settingsStore.filters.trainingstyp = 'Wettkampfvorbereitung'
  } else if (audience === 'hobby') {
    settingsStore.filters.schwierigkeitsgrad = undefined
    settingsStore.filters.trainingstyp = undefined
  }
}

function handleAudienceClick(audience: AudienceType) {
  applyAudienceConfiguration(audience)

  // Trigger visual highlight on the prompt generation button instead of immediate API generation
  if (highlightTransitionTimer) {
    clearTimeout(highlightTransitionTimer)
    highlightTransitionTimer = null
  }
  if (highlightTimer) {
    clearTimeout(highlightTimer)
    highlightTimer = null
  }
  highlightPromptBtn.value = false
  // Short delay to reset transition if repeatedly clicked
  highlightTransitionTimer = setTimeout(() => {
    highlightPromptBtn.value = true
    highlightTimer = setTimeout(() => {
      highlightPromptBtn.value = false
    }, 3500)
  }, 50)
}

function mapCategoryToAudience(category?: string): AudienceType | undefined {
  if (!category) return undefined
  switch (category) {
    case 'Triathlet':
    case 'Triathlete':
    case 'triathlete':
      return 'triathlete'
    case 'Leistungsschwimmer':
    case 'Swimmer':
    case 'Competitive Swimmer':
    case 'competitive_swimmer':
    case 'Trainer':
    case 'Coach':
      return 'competitive_swimmer'
    case 'Hobbyschwimmer':
    case 'Hobby':
    case 'hobby':
      return 'hobby'
    case 'Anfaenger':
    case 'Beginner':
    case 'beginner':
      return 'beginner'
    default:
      return undefined
  }
}

async function handlePromptGeneration(audienceOverride?: AudienceType) {
  generatingPrompt.value = true

  let targetAudience: AudienceType | undefined = audienceOverride

  if (!targetAudience) {
    if (authStore.user) {
      const categories = profileStore.profile?.categories
      const userCategory =
        categories && categories.length > 0
          ? categories[Math.floor(Math.random() * categories.length)]
          : undefined
      targetAudience = mapCategoryToAudience(userCategory)
    } else if (settingsStore.selectedAudience) {
      targetAudience = settingsStore.selectedAudience
    }
  }

  const promptRequest: PromptGenerationRequest = {
    language: locale.value,
    audience: targetAudience || undefined,
  }

  const response = await apiClient.generatePrompt(promptRequest)
  if (response.success) {
    requestText.value = response.data?.prompt || ''
  } else {
    trainingStore.error = response.error
      ? formatError(response.error)
      : t('errors.failed_to_generate_prompt')
  }

  generatingPrompt.value = false
}
</script>

<template>
  <div class="training-plan-form">
    <form @submit.prevent="handleSubmit" class="form-container">
      <!-- Audience selector (only for non-logged-in visitors) -->
      <div v-if="!authStore.user" class="audience-section">
        <label class="audience-label" id="audience-selector-label">
          {{ t('form.audience_label') }}
        </label>
        <div class="audience-selector" role="radiogroup" aria-labelledby="audience-selector-label">
          <button
            v-for="option in audienceOptions"
            :key="option.id"
            :ref="(el) => setAudienceBtnRef(option.id, el)"
            type="button"
            class="audience-btn"
            :class="{ active: settingsStore.selectedAudience === option.id }"
            :aria-checked="settingsStore.selectedAudience === option.id"
            role="radio"
            :disabled="trainingStore.isLoading || generatingPrompt"
            @click="handleAudienceClick(option.id)"
            @mouseenter="handleAudienceMouseEnter(option.id)"
            @mouseleave="handleAudienceMouseLeave(option.id)"
            @focus="handleAudienceMouseEnter(option.id)"
            @blur="handleAudienceMouseLeave(option.id)"
          >
            {{ t(option.labelKey) }}
          </button>
        </div>

        <!-- Teleported hover tooltip for audience buttons (shown on longer hover) -->
        <Teleport to="body">
          <div
            v-if="activeTooltipId"
            id="audience-hover-tooltip"
            class="audience-tooltip"
            :class="`position-${tooltipPlacement}`"
            :style="{ left: `${tooltipPosition.left}px`, top: `${tooltipPosition.top}px` }"
            role="tooltip"
          >
            {{
              t(
                audienceOptions.find((opt) => opt.id === activeTooltipId)?.hintKey ||
                  'form.audience_hint_all',
              )
            }}
          </div>
        </Teleport>
      </div>

      <!-- Main text input -->
      <div class="form-group">
        <label for="request-text" class="form-label">
          {{ t('form.describe_training_needs') }}
          <BaseTooltip>
            <template #tooltip>
              {{ t('form.describe_training_needs_tooltip') }}
            </template>
          </BaseTooltip>
        </label>
        <textarea
          id="request-text"
          v-model="requestText"
          class="form-textarea"
          :placeholder="currentPlaceholder"
          rows="4"
          :disabled="trainingStore.isLoading"
          @input="clearPromptHighlight"
        />
      </div>

      <!-- Advanced settings toggle -->
      <div class="form-middle">
        <button
          type="button"
          @click="toggleAdvancedSettings"
          class="toggle-settings-btn"
          :disabled="trainingStore.isLoading"
        >
          {{
            showAdvancedSettings
              ? t('form.hide_advanced_settings')
              : t('form.show_advanced_settings')
          }}
        </button>

        <!-- Prompt generation button -->
        <button
          type="button"
          @click="handlePromptGeneration()"
          class="toggle-settings-btn"
          :class="{ 'btn-highlight-pulse': highlightPromptBtn }"
          :disabled="trainingStore.isLoading || generatingPrompt"
        >
          <div v-if="!generatingPrompt">{{ t('form.suggest_workout') }}</div>
          <div v-else>{{ t('form.generating') }}</div>
        </button>
      </div>

      <!-- Advanced settings panel -->
      <Transition name="settings-expand">
        <div v-show="showAdvancedSettings" class="advanced-settings">
          <div class="settings-grid">
            <!-- NOTE: this is for v2 Generation Method
          <div class="setting-group">
            <label class="setting-label">Generation Method</label>
            <div class="radio-group">
              <label class="radio-option">
                <input
                  type="radio"
                  value="generate"
                  v-model="settingsStore.preferredMethod"
                  :disabled="trainingStore.isLoading"
                />
                Generate new plan
              </label>
              <label class="radio-option">
                <input
                  type="radio"
                  value="choose"
                  v-model="settingsStore.preferredMethod"
                  :disabled="trainingStore.isLoading"
                />
                Choose existing plan
              </label>
            </div>
            <p class="setting-help">
              Generate creates a new plan, Choose selects from existing plans
            </p>
          </div> -->

            <!-- Pool Length -->
            <div class="setting-group">
              <label class="setting-label">
                {{ t('form.pool_length') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('form.pool_length_tooltip') }}
                  </template>
                </BaseTooltip>
              </label>
              <div class="radio-group">
                <label class="radio-option">
                  <input
                    type="radio"
                    :value="25"
                    v-model="settingsStore.poolLength"
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.pool_length_twenty_five_meters') }}
                </label>
                <label class="radio-option">
                  <input
                    type="radio"
                    :value="50"
                    v-model="settingsStore.poolLength"
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.pool_length_fifty_meters') }}
                </label>
                <label class="radio-option">
                  <input
                    type="radio"
                    :value="'Freiwasser'"
                    v-model="settingsStore.poolLength"
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.pool_length_open_water') }}
                </label>
              </div>
            </div>

            <!-- Swimming Strokes Filter -->
            <div class="setting-group">
              <label class="setting-label">
                {{ t('form.swimming_strokes') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('form.swimming_strokes_tooltip') }}
                  </template>
                </BaseTooltip>
              </label>
              <div class="checkbox-group">
                <label class="checkbox-option">
                  <input
                    type="checkbox"
                    :checked="settingsStore.filters.freistil === true"
                    @change="
                      settingsStore.updateStrokeFilter(
                        'freistil',
                        ($event.target as HTMLInputElement).checked ? true : undefined,
                      )
                    "
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.freestyle') }}
                </label>
                <label class="checkbox-option">
                  <input
                    type="checkbox"
                    :checked="settingsStore.filters.brust === true"
                    @change="
                      settingsStore.updateStrokeFilter(
                        'brust',
                        ($event.target as HTMLInputElement).checked ? true : undefined,
                      )
                    "
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.breaststroke') }}
                </label>
                <label class="checkbox-option">
                  <input
                    type="checkbox"
                    :checked="settingsStore.filters.ruecken === true"
                    @change="
                      settingsStore.updateStrokeFilter(
                        'ruecken',
                        ($event.target as HTMLInputElement).checked ? true : undefined,
                      )
                    "
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.backstroke') }}
                </label>
                <label class="checkbox-option">
                  <input
                    type="checkbox"
                    :checked="settingsStore.filters.delfin === true"
                    @change="
                      settingsStore.updateStrokeFilter(
                        'delfin',
                        ($event.target as HTMLInputElement).checked ? true : undefined,
                      )
                    "
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.butterfly') }}
                </label>
                <label class="checkbox-option">
                  <input
                    type="checkbox"
                    :checked="settingsStore.filters.lagen === true"
                    @change="
                      settingsStore.updateStrokeFilter(
                        'lagen',
                        ($event.target as HTMLInputElement).checked ? true : undefined,
                      )
                    "
                    :disabled="trainingStore.isLoading"
                  />
                  {{ t('form.individual_medley') }}
                </label>
              </div>
            </div>

            <!-- Difficulty Level -->
            <div class="setting-group">
              <label class="setting-label">
                {{ t('form.difficulty_level') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('form.difficulty_level_tooltip') }}
                  </template>
                </BaseTooltip>
              </label>
              <select
                v-model="settingsStore.filters.schwierigkeitsgrad"
                :disabled="trainingStore.isLoading"
                class="select-input"
              >
                <option :value="undefined">{{ t('form.any_difficulty') }}</option>
                <option
                  v-for="option in DIFFICULTY_OPTIONS"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ t(option.label) }}
                </option>
              </select>
            </div>

            <!-- Training Type -->
            <div class="setting-group">
              <label class="setting-label">
                {{ t('form.training_type') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('form.training_type_tooltip') }}
                  </template>
                </BaseTooltip>
              </label>
              <select
                v-model="settingsStore.filters.trainingstyp"
                :disabled="trainingStore.isLoading"
                class="select-input"
              >
                <option :value="undefined">{{ t('form.any_training_type') }}</option>
                <option
                  v-for="option in TRAINING_TYPE_OPTIONS"
                  :key="option.value"
                  :value="option.value"
                >
                  {{ t(option.label) }}
                </option>
              </select>
            </div>

            <!-- Profile Preferences -->
            <div v-if="authStore.user" class="setting-group">
              <label class="setting-label">
                <input
                  type="checkbox"
                  v-model="settingsStore.useProfilePreferences"
                  :disabled="trainingStore.isLoading"
                />
                {{ t('form.use_profile_preferences') }}
                <BaseTooltip>
                  <template #tooltip>
                    {{ t('form.use_profile_preferences_tooltip') }}
                  </template>
                </BaseTooltip>
              </label>
            </div>

            <!-- Data Donation -->
            <!-- <div class="setting-group">
            <label class="setting-label">Privacy Settings</label>
            <label class="checkbox-option">
              <input
                type="checkbox"
                v-model="settingsStore.dataDonationOptOut"
                :disabled="trainingStore.isLoading"
              />
              Opt out of data donation
            </label>
            <p class="setting-help">
              When enabled, your training requests won't be used to improve the system
            </p>
          </div> -->

            <!-- Clear Filters -->
          </div>
          <div class="setting-group">
            <button
              type="button"
              @click="settingsStore.clearFilters"
              :disabled="trainingStore.isLoading"
              class="clear-filters-btn"
            >
              {{ t('form.clear_all_filters') }}
            </button>
            <p class="ai-disclosure">{{ t('ai_disclosure') }}</p>
          </div>
        </div>
      </Transition>

      <!-- Submit button and status -->
      <div class="form-actions">
        <div
          ref="submitWrapperRef"
          class="submit-btn-wrapper"
          @mouseenter="handleSubmitMouseEnter"
          @mouseleave="handleSubmitMouseLeave"
        >
          <button
            type="submit"
            class="submit-btn"
            :disabled="!canSubmit"
            :class="{ loading: trainingStore.isLoading }"
          >
            {{
              trainingStore.isLoading ? t('form.generating_plan') : t('form.generate_training_plan')
            }}
          </button>
        </div>

        <Teleport to="body">
          <div
            v-if="showSubmitTooltip && !canSubmit"
            id="submit-disabled-tooltip"
            class="submit-tooltip"
            :class="`position-${submitTooltipPlacement}`"
            :style="{
              left: `${submitTooltipPosition.left}px`,
              top: `${submitTooltipPosition.top}px`,
            }"
            role="tooltip"
          >
            {{ t('form.generate_training_plan_disabled_tooltip') }}
          </div>
        </Teleport>

        <!-- Too much text error -->
        <div v-if="tooMuchText" class="form-hint text-warning">
          {{ t('form.request_too_long') }}
        </div>

        <!-- Error display -->
        <div v-if="trainingStore.error" class="error-message">
          {{ trainingStore.error }}
          <button type="button" @click="trainingStore.clearError" class="clear-error-btn">×</button>
        </div>
      </div>
    </form>
  </div>
</template>

<style scoped>
.form-container {
  background: var(--color-background-soft);
  padding: 2rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  box-shadow: 0 4px 12px var(--color-shadow);
  width: 100%;
  box-sizing: border-box;
}

@media (max-width: 740px) {
  .form-container {
    padding: 1.25rem;
  }

  .audience-btn {
    min-height: 44px;
    padding: 0.6rem 1rem;
    font-size: 0.95rem;
  }

  .toggle-settings-btn {
    min-height: 44px;
    padding: 0.6rem 1rem;
  }
}

.audience-section {
  margin-bottom: 1.25rem;
}

.audience-label {
  display: block;
  margin-bottom: 0.6rem;
  font-weight: 600;
  color: var(--color-heading);
  font-size: 0.95rem;
  text-align: left;
}

.audience-selector {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 0.6rem;
  width: 100%;
}

.audience-btn {
  background-color: var(--color-background);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  font-family: inherit;
  transition: all 0.2s;
  min-height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  flex: 1 1 calc(25% - 0.6rem);
  max-width: 220px;
  min-width: 140px;
}

.audience-btn:hover:not(:disabled):not(.active) {
  background-color: var(--color-background-mute);
  border-color: var(--color-border-hover);
  color: var(--color-primary);
}

.audience-btn:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px var(--color-shadow);
}

.audience-btn.active {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
  font-weight: 600;
}

.audience-btn.active:hover {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.audience-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.audience-tooltip {
  background-color: var(--color-background-mute);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
  text-align: center;
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  position: fixed;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  font-size: 0.82rem;
  font-weight: 500;
  line-height: 1.4;
  max-width: 260px;
  box-sizing: border-box;
  pointer-events: none;
  animation: tooltip-fade 0.2s ease-out;
}

@keyframes tooltip-fade {
  from {
    opacity: 0;
    transform: translateY(4px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-heading);
}

.form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 1rem;
  resize: vertical;
  min-height: 100px;
  background-color: var(--color-background);
  color: var(--color-text);
}

.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-shadow);
}

.form-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-textarea::placeholder {
  color: color-mix(in srgb, var(--color-text), transparent 40%);
}

.form-hint {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--color-heading);
}

.text-warning {
  color: var(--color-error);
  font-weight: 600;
}

.form-middle {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.submit-btn-wrapper {
  display: flex;
  width: 100%;
}

.submit-btn {
  width: 100%;
  background: var(--color-primary);
  color: white;
  border: 1px solid transparent;
  padding: 0.85rem 1.75rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 4px 12px var(--color-shadow);
  transition: background-color 0.2s ease;
}

.submit-btn:hover:not(:disabled),
.submit-btn.loading {
  background: var(--color-primary-hover);
  box-shadow: 0 6px 12px var(--color-shadow);
  transform: translateY(-1px);
}

.submit-btn:focus-visible {
  outline: none;
  box-shadow:
    0 0 0 3px var(--color-shadow),
    0 4px 12px var(--color-shadow);
}

.submit-btn:disabled {
  opacity: 0.75;
  cursor: not-allowed;
  transform: none;
}

.submit-tooltip {
  background-color: var(--color-background-mute);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
  text-align: center;
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  position: fixed;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  font-size: 0.85rem;
  font-weight: 500;
  line-height: 1.4;
  max-width: 320px;
  box-sizing: border-box;
  pointer-events: none;
  animation: tooltip-fade 0.2s ease-out;
}

.error-message {
  background: #fef2f2;
  color: var(--color-error);
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid #fecaca;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.clear-error-btn {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  font-size: 1.25rem;
  line-height: 1;
}

.toggle-settings-btn {
  background: var(--color-background);
  border: 1px solid var(--color-border);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 1rem;
  color: var(--color-heading);
  font-weight: 500;
  transition: all 0.25s ease;
}

.toggle-settings-btn:hover:not(:disabled) {
  background: var(--color-background-soft);
  border-color: var(--color-border-hover);
  color: var(--color-primary);
}

.toggle-settings-btn:focus-visible {
  outline: none;
  box-shadow: 0 0 0 2px var(--color-shadow);
}

.toggle-settings-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toggle-settings-btn.btn-highlight-pulse {
  border-color: var(--color-primary);
  background-color: var(--color-primary);
  color: #ffffff;
  font-weight: 600;
  box-shadow:
    0 0 0 4px var(--color-shadow),
    0 4px 12px var(--color-shadow);
  animation: prompt-pulse 1.2s ease-in-out infinite alternate;
}

@keyframes prompt-pulse {
  0% {
    transform: scale(1);
    box-shadow:
      0 0 0 3px var(--color-shadow),
      0 2px 6px var(--color-shadow);
  }

  100% {
    transform: scale(1.06);
    box-shadow:
      0 0 0 7px var(--color-shadow),
      0 6px 18px var(--color-shadow);
  }
}

@media (prefers-reduced-motion: reduce) {
  .toggle-settings-btn.btn-highlight-pulse {
    animation: none;
    transform: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 4px var(--color-shadow);
  }
}

.advanced-settings {
  background: var(--color-background);
  padding: 1.5rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  margin-bottom: 1.5rem;
}

.settings-grid {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: 1fr 1fr;
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  color: var(--color-text);
}

@media (max-width: 768px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }

  .setting-group:nth-child(3),
  .setting-group:nth-child(4),
  .setting-group:last-child {
    grid-column: 1 / -1;
  }
}

.setting-label {
  font-weight: 600;
  color: var(--color-heading);
  font-size: 0.9rem;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.radio-option,
.checkbox-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
}

.radio-option input,
.checkbox-option input {
  margin: 0;
}

.radio-option:hover,
.checkbox-option:hover {
  color: var(--color-heading);
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.select-input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.9rem;
  background: var(--color-background);
  color: var(--color-text);
  width: max-content;
}

.select-input:focus {
  outline: none;
  border-color: var(--color-border-hover);
}

.select-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.clear-filters-btn {
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--color-text);
  margin-top: 2rem;
}

.clear-filters-btn:hover:not(:disabled) {
  background: var(--color-background);
}

.clear-filters-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Settings expand transition */
.settings-expand-enter-active,
.settings-expand-leave-active {
  transition: all 0.3s ease;
  max-height: 500px;
  overflow: hidden;
}

.settings-expand-enter-from,
.settings-expand-leave-to {
  opacity: 0;
  max-height: 0;
  margin-top: 0;
  margin-bottom: 0;
  padding-top: 0;
  padding-bottom: 0;
}
</style>
