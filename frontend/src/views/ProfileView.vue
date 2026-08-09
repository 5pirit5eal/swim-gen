<script setup lang="ts">
import BaseTooltip from '@/components/ui/BaseTooltip.vue'
import IconEdit from '@/components/icons/IconEdit.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import IconCross from '@/components/icons/IconCross.vue'
import { toast } from 'vue3-toastify'
import BaseModal from '@/components/ui/BaseModal.vue'
import { useProfileStore } from '@/stores/profile'
import { useAuthStore } from '@/stores/auth'
import { apiClient } from '@/api/client'
import { DIFFICULTY_OPTIONS } from '@/types'
import { onMounted, ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRouter } from 'vue-router'
import {
  calculateCssPace,
  calculateCssZones,
  formatPaceRange,
  formatSwimTime,
} from '@/utils/criticalSwimSpeed'

const { t } = useI18n()
const profileStore = useProfileStore()
const authStore = useAuthStore()
const router = useRouter()

// Feature flag
const showMonthlyQuota = ref(false)

// Delete account state
const showDeleteModal = ref(false)
const deleteConfirmationText = ref('')
const deletingAccount = ref(false)
const deleteError = ref('')

const canDelete = computed(() => deleteConfirmationText.value === 'DELETE')
const isEmailUser = computed(() => authStore.user?.app_metadata?.provider === 'email')

const strokeOptions = ['Freestyle', 'Breaststroke', 'Backstroke', 'Butterfly', 'Individual Medley']
const categoryOptions = ['Triathlete', 'Swimmer', 'Coach', 'Hobby']

const editableProfile = ref({
  experience: '',
  preferred_strokes: [] as string[],
  categories: [] as string[],
  preferred_language: '',
  css_200m_time: '',
  css_400m_time: '',
})
const username = ref('')
const isUsernameEditMode = ref(false)
const usernameEditValue = ref('')
const savedProfileSnapshot = ref('')
const saveConfirmationVisible = ref(false)

function profileSnapshot() {
  return JSON.stringify({
    ...editableProfile.value,
    username: username.value,
  })
}

const hasUnsavedChanges = computed(
  () => Boolean(savedProfileSnapshot.value) && profileSnapshot() !== savedProfileSnapshot.value,
)

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
        css_200m_time: formatSwimTime(newProfile.css_200m_seconds),
        css_400m_time: formatSwimTime(newProfile.css_400m_seconds),
      }
      username.value = newProfile.username || ''
      usernameEditValue.value = newProfile.username || ''
      savedProfileSnapshot.value = profileSnapshot()
    }
  },
  { immediate: true },
)

watch(
  () => navigator.language.split('-')[0] || 'en',
  (lang) => {
    editableProfile.value.preferred_language = lang
    profileStore.updateProfile({ preferred_language: lang })
    if (profileStore.error) {
      toast.error(profileStore.error)
    }
  },
  { immediate: true },
)

async function saveProfile() {
  const { css_200m_time, css_400m_time, ...profileFields } = editableProfile.value
  await profileStore.updateProfile({
    ...profileFields,
    css_200m_seconds: parseProfileTime(css_200m_time),
    css_400m_seconds: parseProfileTime(css_400m_time),
  })
  if (profileStore.error) {
    toast.error(profileStore.error)
    return
  }

  savedProfileSnapshot.value = profileSnapshot()
  saveConfirmationVisible.value = true
  toast.success(t('profile.profile_saved'))
  window.setTimeout(() => {
    saveConfirmationVisible.value = false
  }, 3000)
}

onBeforeRouteLeave(() => {
  if (hasUnsavedChanges.value) {
    toast.warning(t('profile.unsaved_changes'))
  }
})

function parseProfileTime(value: string): number | null {
  const match = value.trim().match(/^(\d+):([0-5]\d)$/)
  return match ? Number(match[1]) * 60 + Number(match[2]) : null
}

const cssPaceSeconds = computed(() =>
  calculateCssPace(editableProfile.value.css_200m_time, editableProfile.value.css_400m_time),
)
const cssZones = computed(() => calculateCssZones(cssPaceSeconds.value))
const cssTimesAreInvalid = computed(() => {
  const time200 = editableProfile.value.css_200m_time.trim()
  const time400 = editableProfile.value.css_400m_time.trim()
  const validFormat = (value: string) => !value || /^\d+:[0-5]\d$/.test(value)

  return (
    !validFormat(time200) ||
    !validFormat(time400) ||
    Boolean(time200 && time400 && !cssPaceSeconds.value)
  )
})

function isInvalidCssTime(value: string): boolean {
  return Boolean(value.trim() && !/^\d+:[0-5]\d$/.test(value.trim()))
}

function openDeleteModal() {
  deleteConfirmationText.value = ''
  deleteError.value = ''
  showDeleteModal.value = true
}

function closeDeleteModal() {
  showDeleteModal.value = false
  deleteConfirmationText.value = ''
  deleteError.value = ''
}

async function confirmDeleteAccount() {
  if (!canDelete.value) return

  deletingAccount.value = true
  deleteError.value = ''

  try {
    const result = await apiClient.deleteUser()
    if (result.success) {
      // Sign out the user and redirect to home
      await authStore.signOut()
      router.push('/')
    } else {
      deleteError.value = result.error?.message || t('profile.delete_error')
    }
  } catch {
    deleteError.value = t('profile.delete_error')
  } finally {
    deletingAccount.value = false
  }
}

async function saveUsername() {
  if (!usernameEditValue.value.trim()) return

  await profileStore.updateProfile({ username: usernameEditValue.value })
  if (profileStore.error) {
    toast.error(profileStore.error)
  } else {
    username.value = usernameEditValue.value
    isUsernameEditMode.value = false
    savedProfileSnapshot.value = profileSnapshot()
  }
}

function cancelUsernameEdit() {
  usernameEditValue.value = username.value
  isUsernameEditMode.value = false
}

const resetCooldown = ref(false)

async function handleResetPassword() {
  if (authStore.user?.email && !resetCooldown.value) {
    const redirectTo = `${window.location.origin}/profile/update-password`
    try {
      resetCooldown.value = true
      await authStore.resetPassword(authStore.user.email, redirectTo)
      toast.success(t('profile.reset_password_success'))

      // Disable button for 10 seconds
      setTimeout(() => {
        resetCooldown.value = false
      }, 10000)
    } catch (error) {
      toast.error((error as Error).message || t('profile.reset_password_error'))
      resetCooldown.value = false
    }
  }
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
          <div class="card-header">
            <h3>{{ t('profile.your_information') }}</h3>
            <button
              @click="saveProfile"
              class="submit-btn"
              :disabled="profileStore.loading || cssTimesAreInvalid"
            >
              <IconCheck v-if="saveConfirmationVisible && !profileStore.loading" />
              {{
                profileStore.loading
                  ? t('profile.saving')
                  : saveConfirmationVisible
                    ? t('profile.saved')
                    : t('profile.save')
              }}
            </button>
          </div>
          <p>{{ t('profile.info_description') }}</p>
          <div class="edit-mode">
            <div class="form-grid">
              <div class="form-column experience-column">
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
              </div>
              <div class="form-column strokes-column">
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
              <div class="form-column categories-column">
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
                  <th v-if="!showMonthlyQuota">
                    {{ t('profile.monthly_generated') }}
                    <BaseTooltip>
                      <template #tooltip>{{ t('profile.monthly_generated_tooltip') }}</template>
                    </BaseTooltip>
                  </th>
                  <th v-else>
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
                  <td v-if="!showMonthlyQuota">
                    {{ profileStore.profile?.monthly_generations ?? 0 }}
                  </td>
                  <td v-else>
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
            <p>{{ t('profile.delete_account_warning') }}</p>
            <button class="delete-btn" @click="openDeleteModal">
              {{ t('profile.delete_account_button') }}
            </button>
          </div>
        </div>
        <div class="css-profile-section">
          <h4>{{ t('profile.css_title') }}</h4>
          <p class="setting-help">{{ t('profile.css_explanation') }}</p>
          <div class="css-input-grid">
            <div class="form-group">
              <label class="form-label" for="css-200m">{{ t('profile.css_200m') }}</label>
              <input
                id="css-200m"
                v-model="editableProfile.css_200m_time"
                class="select-input"
                type="text"
                inputmode="numeric"
                placeholder="3:20"
                :aria-invalid="isInvalidCssTime(editableProfile.css_200m_time)"
                :class="{ 'input-invalid': isInvalidCssTime(editableProfile.css_200m_time) }"
                :disabled="profileStore.loading"
              />
              <small v-if="isInvalidCssTime(editableProfile.css_200m_time)" class="field-error">
                {{ t('profile.css_format_hint') }}
              </small>
            </div>
            <div class="form-group">
              <label class="form-label" for="css-400m">{{ t('profile.css_400m') }}</label>
              <input
                id="css-400m"
                v-model="editableProfile.css_400m_time"
                class="select-input"
                type="text"
                inputmode="numeric"
                placeholder="7:00"
                :aria-invalid="isInvalidCssTime(editableProfile.css_400m_time)"
                :class="{ 'input-invalid': isInvalidCssTime(editableProfile.css_400m_time) }"
                :disabled="profileStore.loading"
              />
              <small v-if="isInvalidCssTime(editableProfile.css_400m_time)" class="field-error">
                {{ t('profile.css_format_hint') }}
              </small>
            </div>
            <div v-if="cssPaceSeconds" class="css-pace-summary">
              <span class="css-result-label">{{ t('profile.css_pace') }}</span>
              <strong class="css-result-value">{{ formatSwimTime(cssPaceSeconds) }} / 100 m</strong>
            </div>
          </div>
          <p v-if="cssTimesAreInvalid" class="css-error">{{ t('profile.css_invalid') }}</p>
          <div v-if="cssPaceSeconds" class="css-result">
            <div class="css-zone-bar" aria-label="CSS training zones">
              <div
                v-for="zone in cssZones"
                :key="zone.key"
                class="css-zone-segment"
                :class="`css-zone-${zone.key}`"
              >
                <strong>{{ t(`profile.css_${zone.key}`) }}</strong>
                <span>{{ t(`profile.css_${zone.key}_focus`) }}</span>
                <b>{{ formatPaceRange(zone) }}</b>
              </div>
            </div>
          </div>
          <p v-else class="setting-help">{{ t('profile.css_missing') }}</p>
        </div>
        <section class="credentials-section">
          <div class="profile-card">
            <h3>{{ t('profile.user_credentials') }}</h3>
            <div class="credentials-grid">
              <div class="info-group">
                <label>{{ t('profile.username') }}</label>
                <div v-if="!isUsernameEditMode" class="value-display">
                  <p>{{ username }}</p>
                  <button @click="isUsernameEditMode = true" class="icon-btn">
                    <IconEdit />
                  </button>
                </div>
                <div v-else class="edit-display">
                  <input
                    v-model="usernameEditValue"
                    type="text"
                    class="select-input"
                    @keyup.enter="saveUsername"
                  />
                  <div class="action-buttons">
                    <button @click="saveUsername" class="icon-btn success">
                      <IconCheck />
                    </button>
                    <button @click="cancelUsernameEdit" class="icon-btn">
                      <IconCross />
                    </button>
                  </div>
                </div>
              </div>
              <div class="info-group">
                <label>{{ t('profile.email') }}</label>
                <p>{{ authStore.user?.email }}</p>
              </div>
              <div class="info-group">
                <label>{{ t('profile.password') }}</label>
                <div class="value-display">
                  <p>{{ t('profile.password_placeholder') }}</p>
                  <button
                    v-if="isEmailUser"
                    @click="handleResetPassword"
                    class="icon-btn"
                    :disabled="resetCooldown"
                  >
                    <IconEdit />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>
      </section>

      <!-- Delete Account Confirmation Modal -->
      <BaseModal :show="showDeleteModal" @close="closeDeleteModal">
        <template #header>
          <h2>{{ t('profile.delete_account_title') }}</h2>
        </template>
        <template #body>
          <div class="delete-modal-content">
            <p class="delete-warning">{{ t('profile.delete_account_confirmation') }}</p>
            <p class="delete-instruction">{{ t('profile.delete_account_instruction') }}</p>
            <input
              v-model="deleteConfirmationText"
              type="text"
              class="delete-confirmation-input"
              :placeholder="t('profile.delete_account_placeholder_input')"
              :disabled="deletingAccount"
            />
            <p v-if="deleteError" class="delete-error">{{ deleteError }}</p>
          </div>
        </template>
        <template #footer>
          <button class="cancel-btn" @click="closeDeleteModal" :disabled="deletingAccount">
            {{ t('profile.cancel') }}
          </button>
          <button
            class="confirm-delete-btn"
            @click="confirmDeleteAccount"
            :disabled="!canDelete || deletingAccount"
          >
            <span v-if="deletingAccount" class="spinner"></span>
            {{ deletingAccount ? t('profile.deleting') : t('profile.delete_account_confirm') }}
          </button>
        </template>
      </BaseModal>
    </div>
  </div>
</template>

<style scoped>
.profile-view {
  padding: 0.25rem 0 2rem 0;
}

.credentials-section {
  margin-bottom: 2rem;
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
  font-weight: 500;
  color: var(--color-heading);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

.profile-content {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(18rem, 0.5fr);
  gap: 2rem;
}

@media (max-width: 1250px) {
  .profile-content {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .profile-card,
  .css-profile-section,
  .statistics-and-delete,
  .credentials-section {
    grid-column: auto;
  }
}

.profile-card {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  position: relative;
  grid-column: 1 / -1;
  grid-row: 1;
}

.statistics-and-delete {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  grid-column: 1 / -1;
  grid-row: 4;
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
  margin-bottom: 0.5rem;
  color: var(--color-heading);
  font-size: 1.5rem;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.card-header h3 {
  margin: 0;
}

.profile-card p {
  margin: 0.5rem 0;
  color: var(--color-text);
  font-size: 1rem;
}

.info-grid {
  display: flex;
  margin: 1rem auto;
  justify-content: space-between;
  gap: 0.5rem;
}

@media (max-width: 1250px) {
  .info-grid {
    margin-bottom: 3.5rem;
  }
}

@media (max-width: 460px) {
  .profile-card {
    padding: 1rem;
  }

  .info-grid {
    flex-direction: column;
    margin-bottom: 3rem;
  }
}

.info-group label {
  gap: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-heading);
}

.info-group p {
  margin: 0.5rem 0;
}

.info-group ul {
  padding-left: 1rem;
  margin: 0.5rem 0;
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
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 2rem;
  margin-bottom: 1.5rem;
}

@media (max-width: 460px) {
  .form-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
}

.form-column {
  min-width: 0;
}

.form-group {
  margin: 0.5rem 0;
}

.css-profile-section {
  grid-column: 1 / -1;
  grid-row: 2;
  margin: 0;
  padding: 2rem;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding-top: 1rem;
}

.credentials-section {
  grid-column: 1 / -1;
  grid-row: 3;
  margin: 0;
}

.css-profile-section h4 {
  margin: 0 0 0.5rem;
  color: var(--color-heading);
  font-size: 1.5rem;
  font-weight: 700;
}

.css-input-grid {
  display: grid;
  gap: 0.75rem;
  grid-template-columns: repeat(2, minmax(0, 1fr)) minmax(15rem, 1fr);
  align-items: end;
}

.css-pace-summary {
  min-height: 2.25rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  color: var(--color-heading);
}

.css-zone-bar {
  display: flex;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  min-height: 10rem;
}

.css-zone-segment {
  display: flex;
  flex: 1 1 0;
  min-width: 0;
  flex-direction: column;
  justify-content: space-between;
  gap: 0.6rem;
  padding: 0.9rem 0.75rem;
  color: #17202a;
}

.css-zone-segment strong {
  font-size: 0.95rem;
  font-weight: 800;
}

.css-zone-segment span {
  font-size: 0.78rem;
  line-height: 1.3;
}

.css-zone-segment b {
  font-size: 0.85rem;
  font-weight: 800;
  white-space: nowrap;
}

.css-zone-z1 {
  background: #d1d5db;
  color: #17202a;
}
.css-zone-z2 {
  background: #3b82f6;
  color: white;
}
.css-zone-z3 {
  background: #16a34a;
  color: white;
}
.css-zone-z4 {
  background: #ea580c;
  color: white;
}
.css-zone-z5 {
  background: #dc2626;
  color: white;
}

.css-zone-segment + .css-zone-segment {
  border-left: 1px solid rgb(255 255 255 / 45%);
}

.css-result {
  display: grid;
  gap: 0.75rem;
  margin-top: 1rem;
}

.css-result-label {
  color: var(--color-heading);
  font-size: 1rem;
  font-weight: 600;
}

.css-result-value {
  color: var(--color-heading);
  font-size: 1.75rem;
  font-weight: 800;
}

.setting-help,
.css-error {
  color: var(--color-text);
  font-size: 0.9rem;
  line-height: 1.5;
}

.css-error {
  color: var(--color-danger, #b42318);
}

.form-label {
  display: block;
  font-size: 1rem;
  font-weight: 600;
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

.css-input-grid .select-input {
  background: #fff;
}

.select-input.input-invalid {
  border-color: var(--color-danger, #b42318);
}

.field-error {
  display: block;
  margin-top: 0.3rem;
  color: var(--color-danger, #b42318);
  font-size: 0.8rem;
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

@media (max-width: 460px) {
  .css-input-grid {
    grid-template-columns: 1fr;
  }

  .css-zone-bar {
    flex-direction: column;
  }

  .css-zone-segment {
    min-height: 7rem;
  }

  .css-zone-segment + .css-zone-segment {
    border-top: 1px solid rgb(255 255 255 / 45%);
    border-left: 0;
  }
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

/* Delete Modal Styles */
.delete-modal-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.delete-warning {
  color: var(--color-error);
  font-weight: 600;
  font-size: 1rem;
}

.delete-instruction {
  color: var(--color-text);
  font-size: 0.9rem;
}

.delete-confirmation-input {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 1rem;
  background: var(--color-background-soft);
  color: var(--color-text);
  width: 100%;
}

.delete-confirmation-input:focus {
  outline: none;
  border-color: var(--color-error);
}

.delete-confirmation-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.delete-error {
  color: var(--color-error);
  font-size: 0.9rem;
}

.cancel-btn {
  background: var(--color-background-soft);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.cancel-btn:hover:not(:disabled) {
  background: var(--color-background-mute);
}

.cancel-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.confirm-delete-btn {
  background: var(--color-error);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.confirm-delete-btn:hover:not(:disabled) {
  background: var(--color-error-soft);
}

.confirm-delete-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner {
  width: 1rem;
  height: 1rem;
  border: 2px solid transparent;
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.credentials-grid {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  background: var(--color-background-soft);
  gap: 1.5rem;
}

@media (max-width: 460px) {
  .credentials-grid {
    flex-direction: column;
  }
}

.value-display {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.icon-btn {
  background: var(--color-background-soft);
  border: 1px solid var(--color-background-soft);
  color: var(--color-text);
  padding: 0.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 0.5rem;
}

.icon-btn svg {
  width: 16px;
  height: 16px;
}

.icon-btn:hover {
  background: var(--color-background-mute);
  border-color: var(--color-border-hover);
}

.edit-display {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.edit-display input {
  max-width: 200px;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.save-btn.small,
.cancel-btn.small {
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
}

.save-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}

.save-btn:hover {
  background: var(--color-primary-hover);
}
</style>
