<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSettingsStore } from '@/stores/settings'
import { apiClient, formatError } from '@/api/client'
import type { QueryRequest, PromptGenerationRequest } from '@/types'
import { DIFFICULTY_OPTIONS, TRAINING_TYPE_OPTIONS } from '@/types'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

// Store access
const authStore = useAuthStore()
const trainingStore = useTrainingPlanStore()
const settingsStore = useSettingsStore()

// i18n
const { t } = useI18n()

// Form data
const requestText = ref('')
const showAdvancedSettings = ref(false)

// Loading state for prompt generation
const generatingPrompt = ref(false)

// Computed
const isFormValid = computed(() => {
  const content = requestText.value.trim()
  return content.length > 0 && content.length <= 3000
})
const tooMuchText = computed(() => requestText.value.trim().length > 3000)
const canSubmit = computed(
  () => isFormValid.value && !trainingStore.isLoading && !generatingPrompt.value,
)

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

async function handlePromptGeneration() {
  generatingPrompt.value = true
  const promptRequest: PromptGenerationRequest = {
    language: navigator.language, // Use current locale
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
          :placeholder="t('form.example_placeholder')"
          rows="4"
          :disabled="trainingStore.isLoading"
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
          @click="handlePromptGeneration"
          class="toggle-settings-btn"
          :disabled="trainingStore.isLoading || generatingPrompt"
        >
          <div v-if="!generatingPrompt">{{ t('form.i_feel_lucky') }}</div>
          <div v-else>{{ t('form.generating') }}</div>
        </button>
      </div>

      <!-- Advanced settings panel -->
      <div v-if="showAdvancedSettings" class="advanced-settings">
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
          <div class="setting-group">
            <label class="setting-label">
              <input
                type="checkbox"
                v-model="settingsStore.useProfilePreferences"
                :disabled="trainingStore.isLoading || !authStore.user"
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
        </div>
      </div>

      <!-- Submit button and status -->
      <div class="form-actions">
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
  width: 100%;
  box-sizing: border-box;
}

@media (max-width: 740px) {
  .form-container {
    padding: 1.5rem;
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
  box-shadow: 0 0 0 2px var(--color-shadow);
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
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.submit-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled),
.submit-btn.loading {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
}

.toggle-settings-btn:hover {
  background: var(--color-background-soft);
}

.toggle-settings-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
</style>
