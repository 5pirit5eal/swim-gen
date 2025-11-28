<script setup lang="ts">
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import { useProfileStore } from '@/stores/profile'
import { DIFFICULTY_OPTIONS } from '@/types'
import { onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const profileStore = useProfileStore()
const isEditMode = ref(false)

const strokeOptions = ['Freestyle', 'Breaststroke', 'Backstroke', 'Butterfly', 'Individual Medley']
const categoryOptions = ['Triathlete', 'Swimmer', 'Coach', 'Hobby']

const editableProfile = ref({
  experience: '',
  preferred_strokes: [] as string[],
  categories: [] as string[],
  preferred_language: '',
})
const username = ref('')

onMounted(() => {
  profileStore.fetchProfile()
})

watch(
  () => profileStore.profile,
  (newProfile) => {
    if (newProfile) {
      editableProfile.value = {
        experience: newProfile.experience || '',
        preferred_strokes: newProfile.preferred_strokes || [],
        categories: newProfile.categories || [],
        preferred_language: newProfile.preferred_language || '',
      }
      username.value = newProfile.username || ''
    }
  },
  { immediate: true },
)

watch(
  () => navigator.language.split('-')[0] || 'en',
  (lang) => {
    editableProfile.value.preferred_language = lang
    profileStore.updateProfile({ preferred_language: lang })
  },
  { immediate: true },
)

function saveProfile() {
  profileStore.updateProfile(editableProfile.value)
  isEditMode.value = false
}

function toggleEditMode() {
  isEditMode.value = !isEditMode.value
}

function getExperienceLabel(value: string) {
  const option = DIFFICULTY_OPTIONS.find((opt) => opt.value === value)
  return option ? t(option.label) : ''
}
</script>

<template>
  <div class="profile-view">
    <div class="container">
      <section class="hero">
        <h1>{{ t('profile.title') }}</h1>
        <p class="hero-description">{{ t('profile.description', { user: username }) }}</p>
      </section>

      <section class="profile-content">
        <div class="profile-card">
          <h3>{{ t('profile.your_information') }}</h3>
          <p>{{ t('profile.info_description') }}</p>
          <div v-if="!isEditMode" class="display-mode">
            <div class="info-grid">
              <div class="info-group">
                <label>
                  {{ t('profile.experience') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.experience_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <p v-if="editableProfile.experience">
                  {{ getExperienceLabel(editableProfile.experience) }}
                </p>
                <p v-else>{{ t('profile.no_selection_placeholder') }}</p>
              </div>
              <div class="info-group">
                <label>
                  {{ t('profile.preferred_strokes') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.preferred_strokes_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <ul v-if="editableProfile.preferred_strokes.length > 0">
                  <li v-for="stroke in editableProfile.preferred_strokes" :key="stroke">
                    {{ t(`profile.${stroke.toLowerCase().replace(' ', '_')}`) }}
                  </li>
                </ul>
                <p v-else>{{ t('profile.no_selection_placeholder') }}</p>
              </div>
              <div class="info-group">
                <label>
                  {{ t('profile.categories') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.categories_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <ul v-if="editableProfile.categories.length > 0">
                  <li v-for="category in editableProfile.categories" :key="category">
                    {{ t(`profile.category_${category.toLowerCase()}`) }}
                  </li>
                </ul>
                <p v-else>{{ t('profile.no_selection_placeholder') }}</p>
              </div>
            </div>
            <button @click="toggleEditMode" class="edit-btn">
              {{ t('profile.edit_profile') }}
            </button>
          </div>
          <div v-else class="edit-mode">
            <div class="form-grid">
              <div class="form-column">
                <div class="form-group">
                  <label class="form-label"
                    >{{ t('profile.experience') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.experience_explanation') }}</template>
                    </BaseTooltip>
                  </label>
                  <div class="select-group">
                    <select
                      class="select-input"
                      v-model="editableProfile.experience"
                      :disabled="profileStore.loading"
                    >
                      <option value="">{{ t('form.any_difficulty') }}</option>
                      <option
                        v-for="option in DIFFICULTY_OPTIONS"
                        :key="option.value"
                        :value="option.value"
                      >
                        {{ t(option.label) }}
                      </option>
                    </select>
                  </div>
                </div>
                <div class="form-group">
                  <label class="form-label"
                    >{{ t('profile.preferred_strokes') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.preferred_strokes_explanation') }}</template>
                    </BaseTooltip>
                  </label>
                  <div class="checkbox-group">
                    <label v-for="option in strokeOptions" :key="option" class="checkbox-option">
                      <input
                        type="checkbox"
                        :value="option"
                        v-model="editableProfile.preferred_strokes"
                        :disabled="profileStore.loading"
                      />
                      {{ t(`profile.${option.toLowerCase().replace(' ', '_')}`) }}
                    </label>
                  </div>
                </div>
              </div>
              <div class="form-column">
                <div class="form-group">
                  <label class="form-label"
                    >{{ t('profile.categories') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.categories_explanation') }}</template>
                    </BaseTooltip>
                  </label>
                  <div class="checkbox-group">
                    <label v-for="option in categoryOptions" :key="option" class="checkbox-option">
                      <input
                        type="checkbox"
                        :value="option"
                        v-model="editableProfile.categories"
                        :disabled="profileStore.loading"
                      />
                      {{ t(`profile.category_${option.toLowerCase()}`) }}
                    </label>
                  </div>
                </div>
              </div>
            </div>
            <button @click="saveProfile" class="submit-btn" :disabled="profileStore.loading">
              {{ profileStore.loading ? t('profile.saving') : t('profile.save') }}
            </button>
          </div>
        </div>

        <div class="statistics-and-delete">
          <div class="statistics-card">
            <table class="statistics-table">
              <thead>
                <tr>
                  <th>
                    {{ t('profile.generated_plans') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.generated_plans_tooltip') }}</template>
                    </BaseTooltip>
                  </th>
                  <th>
                    {{ t('profile.exported_plans') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.exported_plans_tooltip') }}</template>
                    </BaseTooltip>
                  </th>
                  <th>
                    {{ t('profile.monthly_quota') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.monthly_quota_tooltip') }}</template>
                    </BaseTooltip>
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>{{ profileStore.profile?.overall_generations ?? 0 }}</td>
                  <td>{{ profileStore.profile?.exports ?? 0 }}</td>
                  <td>
                    <p>{{ profileStore.profile?.monthly_generations ?? 0 }} / 100</p>
                    <div class="progress-bar">
                      <div
                        class="progress"
                        :style="{ width: `${profileStore.profile?.monthly_generations ?? 0}%` }"
                      ></div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="delete-card">
            <p>{{ t('profile.delete_account_placeholder') }}</p>
            <button class="delete-btn">{{ t('profile.delete_account_button') }}</button>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.profile-view {
  padding: 0.25rem 0 2rem 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

.hero {
  text-align: center;
  margin-bottom: 2rem;
  background-color: var(--color-transparent);
  backdrop-filter: blur(2px);
  border-radius: 8px;
  padding: 1rem;
  margin: 2rem auto;
}

.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-heading);
  margin-bottom: 1rem;
}

.hero-description {
  font-size: 1.25rem;
  color: var(--color-text);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

.profile-content {
  display: flex;
  gap: 2rem;
}

@media (max-width: 1186px) {
  .profile-content {
    flex-direction: column;
  }
}

.profile-card {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  position: relative;
}

.statistics-and-delete {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.profile-card,
.statistics-card,
.delete-card {
  background: var(--color-background-soft);
  padding: 2rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
}

.profile-card h3 {
  margin-bottom: 1rem;
  color: var(--color-heading);
  font-size: 1.5rem;
}

.profile-card p {
  margin-bottom: 1rem;
  color: var(--color-text);
  font-size: 1rem;
}

.info-grid {
  display: flex;
  margin-bottom: 2rem;
  justify-content: space-between;
}

@media (max-width: 460px) {
  .info-grid {
    flex-direction: column;
  }
}

.info-group {
  margin: 0.25rem 0.5rem 1.5rem;
}

.info-group label {
  gap: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-heading);
}

.info-group p {
  margin-top: 0.5rem;
}

.info-group ul {
  padding-left: 1rem;
  margin-top: 0.5rem;
}

.edit-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  position: absolute;
  bottom: 2rem;
  left: 2rem;
}

.edit-btn:hover {
  background: var(--color-primary-hover);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  margin-bottom: 1.5rem;
}

.form-column {
  padding: 1rem 0rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-heading);
}

.radio-group,
.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  color: var(--color-text);
}

.select-input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.9rem;
  background: var(--color-background-soft);
  color: var(--color-text);
}

.select-input:focus {
  outline: none;
  border-color: var(--color-border-hover);
}

.radio-option,
.checkbox-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 1rem;
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
  position: absolute;
  bottom: 2rem;
  left: 2rem;
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.statistics-table {
  display: table;
  text-align: center;
  justify-content: space-between;
  width: 100%;
}

.statistics-table th {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-heading);
  vertical-align: top;
  padding: 0 calc(max(5%, 0.5rem));
  width: 25%;
}

.statistics-table td {
  font-size: 1.5rem;
  font-weight: 300;
  color: var(--color-text);
  padding: 0 calc(max(10%, 0.5rem));
  text-wrap: nowrap;
}

.statistics-table p {
  width: 200%;
  position: relative;
  transform: translateX(-20%);
}

.progress-bar {
  background-color: var(--color-border);
  border-radius: 8px;
  height: 0.5rem;
  margin-top: 0.5rem;
  width: 200%;
  transform: translateX(-20%);
}

.progress {
  background-color: var(--color-primary);
  height: 100%;
  border-radius: 8px;
}

.delete-card {
  border-color: var(--color-error);
  background-color: var(--color-background-soft);
}

.delete-card p {
  color: var(--color-error);
  margin-bottom: 1rem;
}

.delete-btn {
  background: var(--color-error);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.delete-btn:hover {
  background: var(--color-error-soft);
}
</style>
