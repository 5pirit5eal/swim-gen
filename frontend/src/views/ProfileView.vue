<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useProfileStore } from '@/stores/profile'
import { useI18n } from 'vue-i18n'
import BaseTooltip from '@/components/ui/BaseTooltip.vue'

const { t } = useI18n()
const profileStore = useProfileStore()
const isEditMode = ref(false)

const experienceOptions = ['Beginner', 'Intermediate', 'Advanced']
const strokeOptions = ['Freestyle', 'Breaststroke', 'Backstroke', 'Butterfly', 'Individual Medley']
const categoryOptions = ['Triathlete', 'Swimmer', 'Coach', 'Hobby']

const editableProfile = ref({
  experience: '',
  preferred_strokes: [] as string[],
  categories: [] as string[]
})

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
        categories: newProfile.categories || []
      }
    }
  },
  { immediate: true }
)

function saveProfile() {
  profileStore.updateProfile(editableProfile.value)
  isEditMode.value = false
}

function toggleEditMode() {
  isEditMode.value = !isEditMode.value
}
</script>

<template>
  <div class="profile-view">
    <div class="container">
      <section class="hero">
        <h1>{{ t('profile.title') }}</h1>
        <p class="hero-description">{{ t('profile.description') }}</p>
      </section>

      <section class="profile-content">
        <div class="profile-card">
          <h3>{{ t('profile.your_information') }}</h3>
          <p>{{ t('profile.info_description', { username: profileStore.profile?.username }) }}</p>
          <div v-if="!isEditMode" class="display-mode">
            <div class="info-grid">
              <div class="info-group">
                <label>
                  {{ t('profile.experience') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.experience_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <p>{{ editableProfile.experience }}</p>
              </div>
              <div class="info-group">
                <label>
                  {{ t('profile.preferred_strokes') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.preferred_strokes_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <p>{{ editableProfile.preferred_strokes.join(', ') }}</p>
              </div>
              <div class="info-group">
                <label>
                  {{ t('profile.categories') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.categories_explanation') }}</template>
                  </BaseTooltip>
                </label>
                <p>{{ editableProfile.categories.join(', ') }}</p>
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
                  <label class="form-label">{{ t('profile.experience') }}</label>
                  <div class="radio-group">
                    <label v-for="option in experienceOptions" :key="option" class="radio-option">
                      <input type="radio" :value="option" v-model="editableProfile.experience"
                        :disabled="profileStore.loading" />
                      {{ option }}
                    </label>
                  </div>
                </div>
                <div class="form-group">
                  <label class="form-label">{{ t('profile.preferred_strokes') }}</label>
                  <div class="checkbox-group">
                    <label v-for="option in strokeOptions" :key="option" class="checkbox-option">
                      <input type="checkbox" :value="option" v-model="editableProfile.preferred_strokes"
                        :disabled="profileStore.loading" />
                      {{ option }}
                    </label>
                  </div>
                </div>
              </div>
              <div class="form-column">
                <div class="form-group">
                  <label class="form-label">{{ t('profile.categories') }}</label>
                  <div class="checkbox-group">
                    <label v-for="option in categoryOptions" :key="option" class="checkbox-option">
                      <input type="checkbox" :value="option" v-model="editableProfile.categories"
                        :disabled="profileStore.loading" />
                      {{ option }}
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
            <div class="statistics-grid">
              <div class="stat-item">
                <h3>{{ t('profile.generated_plans') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.generated_plans_tooltip') }}</template>
                  </BaseTooltip>
                </h3>
                <p>0</p>
              </div>
              <div class="stat-item">
                <h3>{{ t('profile.exported_plans') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.exported_plans_tooltip') }}</template>
                  </BaseTooltip>
                </h3>
                <p>0</p>
              </div>
              <div class="stat-item">
                <h3>{{ t('profile.monthly_quota') }}
                  <BaseTooltip>
                    <template #tooltip>{{ t('profile.monthly_quota_tooltip') }}</template>
                  </BaseTooltip>
                </h3>
                <p>10 / 100</p>
                <div class="progress-bar">
                  <div class="progress" style="width: 10%"></div>
                </div>
              </div>
            </div>
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
  border-radius: 0.5rem;
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
  align-items: column;
}

.info-group {
  margin: 0.25rem 0.5rem 1.5rem;
}

.info-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-heading);
}

.info-group p {
  color: var(--color-text);
  margin-top: 0.5rem;
}

.edit-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  align-self: self-end;
}

.edit-btn:hover {
  background: var(--color-primary-hover);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
}

.form-column {
  padding: 1rem 0rem;
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

.radio-group,
.checkbox-group {
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
  font-size: 1rem;
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

.submit-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.statistics-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  text-align: center;
}

.stat-item h3 {
  font-size: 1rem;
  font-weight: 1rem;
  color: var(--color-heading);
}

.stat-item p {
  font-size: 1.5rem;
  font-weight: 0.5rem;
  color: var(--color-text);
}

.progress-bar {
  background-color: var(--color-border);
  border-radius: 0.25rem;
  height: 0.5rem;
  margin-top: 0.5rem;
}

.progress {
  background-color: var(--color-primary);
  height: 100%;
  border-radius: 0.25rem;
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
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.delete-btn:hover {
  background: var(--color-error-soft);
}
</style>
