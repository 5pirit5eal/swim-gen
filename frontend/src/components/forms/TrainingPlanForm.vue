<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSettingsStore } from '@/stores/settings'
import type { QueryRequest } from '@/types'
import { DIFFICULTY_OPTIONS, TRAINING_TYPE_OPTIONS } from '@/types'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'

// Store access
const trainingStore = useTrainingPlanStore()
const settingsStore = useSettingsStore()

// Form data
const requestText = ref('')
const showAdvancedSettings = ref(false)

// Computed
const isFormValid = computed(() => {
  const content = requestText.value.trim()
  return content.length > 0 && content.length <= 3000
})
const tooMuchText = computed(() => requestText.value.trim().length > 3000)
const canSubmit = computed(() => isFormValid.value && !trainingStore.isLoading)

// Actions
async function handleSubmit() {
  if (!canSubmit.value) return

  const request: QueryRequest = {
    content: requestText.value.trim(),
    method: settingsStore.preferredMethod,
    filter: settingsStore.filters,
  }

  const success = await trainingStore.generatePlan(request)
  if (success) {
    // Plan generated successfully - it's now in trainingStore.currentPlan
    console.log('Plan generated successfully!')
  }
}

function toggleAdvancedSettings() {
  showAdvancedSettings.value = !showAdvancedSettings.value
}
</script>

<template>
  <div class="training-plan-form">
    <form @submit.prevent="handleSubmit" class="form-container">
      <!-- Main text input -->
      <div class="form-group">
        <label for="request-text" class="form-label">
          Describe your training needs
          <BaseTooltip>
            <template #tooltip>
              Be specific about your goals, experience level, time constraints, and preferences.
            </template>
          </BaseTooltip>
        </label>
        <textarea id="request-text" v-model="requestText" class="form-textarea"
          placeholder="Example: I need a 45-minute freestyle endurance workout for an intermediate swimmer..." rows="4"
          :disabled="trainingStore.isLoading" />
      </div>

      <!-- Advanced settings toggle -->
      <button type="button" @click="toggleAdvancedSettings" class="toggle-settings-btn"
        :disabled="trainingStore.isLoading">
        {{ showAdvancedSettings ? 'Hide' : 'Show' }} Advanced Settings
      </button>

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
          </div>

           Pool Length
          <div class="setting-group">
            <label class="setting-label">Pool Length</label>
            <div class="radio-group">
              <label class="radio-option">
                <input
                  type="radio"
                  :value="25"
                  v-model.number="settingsStore.poolLength"
                  :disabled="trainingStore.isLoading"
                />
                25 meters
              </label>
              <label class="radio-option">
                <input
                  type="radio"
                  :value="50"
                  v-model.number="settingsStore.poolLength"
                  :disabled="trainingStore.isLoading"
                />
                50 meters
              </label>
            </div>
            <p class="setting-help">Specify your pool length for accurate distance calculations</p>
          </div> -->

          <!-- Swimming Strokes Filter -->
          <div class="setting-group">
            <label class="setting-label">
              Swimming Strokes
              <BaseTooltip>
                <template #tooltip>Select specific swimming strokes to focus on</template>
              </BaseTooltip>
            </label>
            <div class="checkbox-group">
              <label class="checkbox-option">
                <input type="checkbox" :checked="settingsStore.filters.freistil === true" @change="
                  settingsStore.updateStrokeFilter(
                    'freistil',
                    ($event.target as HTMLInputElement).checked ? true : undefined,
                  )
                  " :disabled="trainingStore.isLoading" />
                Freestyle
              </label>
              <label class="checkbox-option">
                <input type="checkbox" :checked="settingsStore.filters.brust === true" @change="
                  settingsStore.updateStrokeFilter(
                    'brust',
                    ($event.target as HTMLInputElement).checked ? true : undefined,
                  )
                  " :disabled="trainingStore.isLoading" />
                Breaststroke
              </label>
              <label class="checkbox-option">
                <input type="checkbox" :checked="settingsStore.filters.ruecken === true" @change="
                  settingsStore.updateStrokeFilter(
                    'ruecken',
                    ($event.target as HTMLInputElement).checked ? true : undefined,
                  )
                  " :disabled="trainingStore.isLoading" />
                Backstroke
              </label>
              <label class="checkbox-option">
                <input type="checkbox" :checked="settingsStore.filters.delfin === true" @change="
                  settingsStore.updateStrokeFilter(
                    'delfin',
                    ($event.target as HTMLInputElement).checked ? true : undefined,
                  )
                  " :disabled="trainingStore.isLoading" />
                Butterfly
              </label>
              <label class="checkbox-option">
                <input type="checkbox" :checked="settingsStore.filters.lagen === true" @change="
                  settingsStore.updateStrokeFilter(
                    'lagen',
                    ($event.target as HTMLInputElement).checked ? true : undefined,
                  )
                  " :disabled="trainingStore.isLoading" />
                Individual Medley
              </label>
            </div>
          </div>

          <div class="settings-grid">
            <!-- Difficulty Level -->
            <div class="setting-group">
              <label class="setting-label">
                Difficulty Level
                <BaseTooltip>
                  <template #tooltip>Filter plans by swimmer experience level</template>
                </BaseTooltip>
              </label>
              <select v-model="settingsStore.filters.schwierigkeitsgrad" :disabled="trainingStore.isLoading"
                class="select-input">
                <option :value="undefined">Any difficulty</option>
                <option v-for="option in DIFFICULTY_OPTIONS" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>

            <!-- Training Type -->
            <div class="setting-group">
              <label class="setting-label">
                Training Type
                <BaseTooltip>
                  <template #tooltip>Filter plans by training focus and goals</template>
                </BaseTooltip>
              </label>
              <select v-model="settingsStore.filters.trainingstyp" :disabled="trainingStore.isLoading"
                class="select-input">
                <option :value="undefined">Any training type</option>
                <option v-for="option in TRAINING_TYPE_OPTIONS" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
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
          <div class="setting-group">
            <button type="button" @click="settingsStore.clearFilters" :disabled="trainingStore.isLoading"
              class="clear-filters-btn">
              Clear All Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Submit button and status -->
      <div class="form-actions">
        <button type="submit" class="submit-btn" :disabled="!canSubmit" :class="{ loading: trainingStore.isLoading }">
          {{ trainingStore.isLoading ? 'Generating...' : 'Generate Training Plan' }}
        </button>

        <!-- Too much text error -->
        <div v-if="tooMuchText" class="form-hint text-warning">
          Your request is too long! Please limit it to 3000 characters.
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
  border-radius: 0.5rem;
  border: 1px solid var(--color-border);
  width: 100%;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .form-container {
    padding: 1.5rem;
  }
}

.form-group {
  margin-bottom: 1.5rem;
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
  border-radius: 0.25rem;
  font-family: inherit;
  font-size: 1rem;
  resize: vertical;
  min-height: 100px;
  background-color: var(--color-background);
  color: var(--color-text);
}

.form-textarea:focus {
  outline: none;
  border-color: var(--color-border-hover);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.form-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-hint {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.text-warning {
  color: #dc2626;
  font-weight: 600;
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
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) .submit-btn.loading {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  background: #fef2f2;
  color: #dc2626;
  padding: 0.75rem;
  border-radius: 0.25rem;
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
  border-radius: 0.25rem;
  cursor: pointer;
  margin-bottom: 1rem;
  color: var(--color-heading);
}

.toggle-settings-btn:hover {
  background: var(--color-background-soft);
}

.advanced-settings {
  background: var(--color-background);
  padding: 1.5rem;
  border-radius: 0.25rem;
  border: 1px solid var(--color-border);
  margin-bottom: 1.5rem;
}

.settings-grid {
  display: grid;
  gap: 1.5rem;
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  color: var(--color-text-light);
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

.setting-help {
  font-size: 0.8rem;
  color: var(--color-text-light);
  margin: 0;
  line-height: 1.4;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.select-input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 0.25rem;
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
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--color-text);
}

.clear-filters-btn:hover:not(:disabled) {
  background: var(--color-background);
  color: var(--color-text-light);
}

.clear-filters-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (min-width: 768px) {
  .settings-grid {
    grid-template-columns: 1fr 1fr;
  }

  .setting-group:nth-child(3),
  .setting-group:nth-child(4),
  .setting-group:last-child {
    grid-column: 1 / -1;
  }
}
</style>
