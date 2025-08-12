<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSettingsStore } from '@/stores/settings'
import type { QueryRequest } from '@/types'

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
        <!-- Settings controls will go here -->
        <p class="settings-placeholder">⚙️ Settings controls coming soon...</p>
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

.settings-placeholder {
  color: var(--color-text-light);
  font-style: italic;
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
</style>
