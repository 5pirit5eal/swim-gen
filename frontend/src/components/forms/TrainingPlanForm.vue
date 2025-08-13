<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSettingsStore } from '@/stores/settings'
import type { QueryRequest } from '@/types'
import { DIFFICULTY_OPTIONS, TRAINING_TYPE_OPTIONS } from '@/types'

// Store access
const trainingStore = useTrainingPlanStore()
const settingsStore = useSettingsStore()

// Form data
const requestText = ref('')
const showAdvancedSettings = ref(false)

// Computed
const isFormValid = computed(() => requestText.value.trim().length > 0)
const canSubmit = computed(() => isFormValid.value && !trainingStore.isGenerating)

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
  <div class="training-plan-form content-container">
    <form @submit.prevent="handleSubmit" class="form-container">
      <!-- Main text input -->
      <div class="form-group">
        <label for="request-text" class="form-label"> Describe your training needs </label>
        <textarea
          id="request-text"
          v-model="requestText"
          class="form-textarea"
          placeholder="Example: I need a 45-minute freestyle endurance workout for an intermediate swimmer..."
          rows="4"
          :disabled="trainingStore.isGenerating"
        />
        <p class="form-hint">
          Be specific about your goals, experience level, time constraints, and preferences.
        </p>
      </div>

      <!-- Advanced settings toggle -->
      <button
        type="button"
        @click="toggleAdvancedSettings"
        class="toggle-settings-btn"
        :disabled="trainingStore.isGenerating"
      >
        {{ showAdvancedSettings ? 'Hide' : 'Show' }} Advanced Settings
      </button>

      <!-- Advanced settings panel -->
      <div v-if="showAdvancedSettings" class="advanced-settings">
        <h3>Advanced Settings</h3>

        <div class="settings-grid">
          <!-- Generation Method -->
          <div class="setting-group">
            <label class="setting-label">Generation Method</label>
            <div class="radio-group">
              <label class="radio-option">
                <input
                  type="radio"
                  value="generate"
                  v-model="settingsStore.preferredMethod"
                  :disabled="trainingStore.isGenerating"
                />
                Generate new plan
              </label>
              <label class="radio-option">
                <input
                  type="radio"
                  value="choose"
                  v-model="settingsStore.preferredMethod"
                  :disabled="trainingStore.isGenerating"
                />
                Choose existing plan
              </label>
            </div>
            <p class="setting-help">
              Generate creates a new plan, Choose selects from existing plans
            </p>
          </div>

          <!-- Pool Length -->
          <div class="setting-group">
            <label class="setting-label">Pool Length</label>
            <div class="radio-group">
              <label class="radio-option">
                <input
                  type="radio"
                  :value="25"
                  v-model.number="settingsStore.poolLength"
                  :disabled="trainingStore.isGenerating"
                />
                25 meters
              </label>
              <label class="radio-option">
                <input
                  type="radio"
                  :value="50"
                  v-model.number="settingsStore.poolLength"
                  :disabled="trainingStore.isGenerating"
                />
                50 meters
              </label>
            </div>
            <p class="setting-help">Specify your pool length for accurate distance calculations</p>
          </div>

          <!-- Swimming Strokes Filter -->
          <div class="setting-group">
            <label class="setting-label">Swimming Strokes</label>
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
                  :disabled="trainingStore.isGenerating"
                />
                Freestyle
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
                  :disabled="trainingStore.isGenerating"
                />
                Breaststroke
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
                  :disabled="trainingStore.isGenerating"
                />
                Backstroke
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
                  :disabled="trainingStore.isGenerating"
                />
                Butterfly
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
                  :disabled="trainingStore.isGenerating"
                />
                Individual Medley
              </label>
            </div>
            <p class="setting-help">Select specific swimming strokes to focus on</p>
          </div>

          <!-- Difficulty Level -->
          <div class="setting-group">
            <label class="setting-label">Difficulty Level</label>
            <select
              v-model="settingsStore.filters.schwierigkeitsgrad"
              :disabled="trainingStore.isGenerating"
              class="select-input"
            >
              <option :value="undefined">Any difficulty</option>
              <option
                v-for="option in DIFFICULTY_OPTIONS"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </select>
            <p class="setting-help">Filter plans by swimmer experience level</p>
          </div>

          <!-- Training Type -->
          <div class="setting-group">
            <label class="setting-label">Training Type</label>
            <select
              v-model="settingsStore.filters.trainingstyp"
              :disabled="trainingStore.isGenerating"
              class="select-input"
            >
              <option :value="undefined">Any training type</option>
              <option
                v-for="option in TRAINING_TYPE_OPTIONS"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </select>
            <p class="setting-help">Filter plans by training focus and goals</p>
          </div>

          <!-- Data Donation -->
          <!-- <div class="setting-group">
            <label class="setting-label">Privacy Settings</label>
            <label class="checkbox-option">
              <input
                type="checkbox"
                v-model="settingsStore.dataDonationOptOut"
                :disabled="trainingStore.isGenerating"
              />
              Opt out of data donation
            </label>
            <p class="setting-help">
              When enabled, your training requests won't be used to improve the system
            </p>
          </div> -->

          <!-- Clear Filters -->
          <div class="setting-group">
            <button
              type="button"
              @click="settingsStore.clearFilters"
              :disabled="trainingStore.isGenerating"
              class="clear-filters-btn"
            >
              Clear All Filters
            </button>
            <p class="setting-help">Reset all filter options to default values</p>
          </div>
        </div>
      </div>

      <!-- Submit button and status -->
      <div class="form-actions">
        <button
          type="submit"
          class="submit-btn"
          :disabled="!canSubmit"
          :class="{ loading: trainingStore.isGenerating }"
        >
          {{ trainingStore.isGenerating ? 'Generating...' : 'Generate Training Plan' }}
        </button>

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

.toggle-settings-btn {
  background: none;
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

.advanced-settings h3 {
  margin: 0 0 1rem 0;
  font-size: 1.125rem;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.submit-btn {
  background: var(--color-primary, #3b82f6);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover, #2563eb);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.submit-btn.loading {
  background: var(--color-text-light);
  border-color: var(--color-border);
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
  color: var(--color-text-light);
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
