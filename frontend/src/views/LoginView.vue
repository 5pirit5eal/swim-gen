<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { useSidebarStore } from '@/stores/sidebar'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue3-toastify'
import IconGoogle from '@/components/icons/IconGoogle.vue'

const { t } = useI18n()
const auth = useAuthStore()
const trainingPlanStore = useTrainingPlanStore()
const sharedPlanStore = useSharedPlanStore()
const sidebarStore = useSidebarStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const username = ref('')
const loading = ref(false)

const features = ['history', 'share', 'upload', 'personalize', 'interactive'] as const

onMounted(() => {
  if (auth.user) {
    router.replace('/')
  }
})

const isSignUp = computed(() => route.query.register === 'true')

const canSubmit = computed(() => {
  if (isSignUp.value) {
    return email.value && password.value && username.value
  } else {
    return email.value && password.value
  }
})

async function handleLogin() {
  loading.value = true
  try {
    await auth.signInWithPassword(email.value, password.value)
    sidebarStore.open()
    toast.success(t('login.loginSuccess'))
    router.push('/')
  } catch (error) {
    console.error('Login failed:', error) // Log the full error
    if (error instanceof Error) {
      if (error.message.includes('Invalid login credentials')) {
        toast.error(t('login.invalidLogin'))
      } else if (error.message.includes('Email not confirmed')) {
        toast.error(t('login.emailNotConfirmed'))
      } else {
        toast.error(t('login.unknownError'))
      }
    } else {
      toast.error(t('login.unknownError'))
    }
  } finally {
    loading.value = false
  }
}

async function handleResetPassword() {
  if (!email.value) return
  loading.value = true
  try {
    const redirectTo = `${window.location.origin}/profile/update-password`
    await auth.resetPassword(email.value, redirectTo)
    toast.success(t('profile.reset_password_success'))
  } catch (error) {
    console.error('Reset password failed:', error)
    toast.error((error as Error).message || t('profile.reset_password_error'))
  } finally {
    loading.value = false
  }
}

async function handleGoogleLogin() {
  loading.value = true
  try {
    await auth.signInWithOAuth()
  } catch (error) {
    console.error('Google Login failed:', error)
    toast.error(t('login.unknownError'))
    loading.value = false
  }
}

async function handleSignUp() {
  loading.value = true
  let response
  try {
    response = await auth.signUp(email.value, password.value, username.value)
  } catch (error) {
    if (error instanceof Error) {
      if (error.message.includes('Username already taken')) {
        toast.error(t('login.usernameTaken'))
        loading.value = false
        return
      } else if (error.message.includes('User already registered')) {
        console.log('User exists, attempting login...')
      } else {
        toast.error(t('login.unknownError'))
      }
    } else {
      toast.error(t('login.unknownError'))
    }
  }

  if (!response?.user?.identities?.length) {
    try {
      response = await auth.signInWithPassword(email.value, password.value)
      await trainingPlanStore.fetchHistory()
      await sharedPlanStore.fetchSharedHistory()
      sidebarStore.open()
      toast.success(t('login.userExistsLoginSuccess'))
      router.push({ path: '/', state: { redirectedFromLogin: true } })
    } catch {
      toast.error(t('login.userExistsLoginFailed'))
    } finally {
      loading.value = false
    }
  } else {
    toast.success(t('login.registrationSuccess'))
    // Clear form
    email.value = ''
    password.value = ''
    router.push('/login')
  }
  loading.value = false
}
</script>

<template>
  <div class="login-view">
    <div class="column" v-if="isSignUp">
      <div class="features-box">
        <h2>{{ t('login.features.title') }}</h2>
        <p class="features-subtitle">{{ t('login.features.subtitle') }}</p>

        <ul class="features-list">
          <li v-for="feature in features" :key="feature" class="feature-item">
            <span class="feature-short">{{ t(`login.features.items.${feature}.short`) }}</span>
            <span class="feature-long">{{ t(`login.features.items.${feature}.long`) }}</span>
          </li>
        </ul>
        <p class="contact-text">{{ t('login.features.contact') }}</p>
      </div>
    </div>
    <div class="column">
      <div class="login-box">
        <h1>{{ isSignUp ? t('login.signUp') : t('login.login') }}</h1>
        <form @submit.prevent="isSignUp ? handleSignUp() : handleLogin()">
          <div class="form-group" v-if="isSignUp">
            <label for="username">{{ t('login.username') }}*</label>
            <input
              id="username"
              type="text"
              :placeholder="t('login.username')"
              v-model="username"
              required
            />
          </div>
          <div class="form-group">
            <label for="email">{{ t('login.email') }}*</label>
            <input
              id="email"
              type="email"
              :placeholder="t('login.email')"
              v-model="email"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">{{ t('login.password') }}*</label>
            <input
              id="password"
              type="password"
              :placeholder="t('login.password')"
              v-model="password"
              required
            />
          </div>
          <div class="switch-form">
            <router-link v-if="isSignUp" to="/login">{{ t('login.haveAccount') }}</router-link>
            <div v-else class="login-links">
              <router-link to="/login?register=true">{{ t('login.needAccount') }}</router-link>
              <button
                type="button"
                class="text-btn"
                @click="handleResetPassword"
                :disabled="!email || loading"
              >
                {{ t('login.forgot_password') }}
              </button>
            </div>
          </div>
          <button type="submit" :disabled="!canSubmit || loading">
            {{ loading ? t('login.loading') : isSignUp ? t('login.signUp') : t('login.login') }}
          </button>

          <div class="divider">
            <span>{{ t('login.or') }}</span>
          </div>

          <button type="button" class="google-btn" @click="handleGoogleLogin" :disabled="loading">
            <IconGoogle class="google-icon" />
            {{ t('login.signInWithGoogle') }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  flex-direction: row;
  padding: 2rem 1rem;
  gap: 2rem;
  margin: 0 auto;
}

.column {
  display: flex;
  flex-direction: column;
}

.features-box {
  height: 100%;
  max-width: 700px;
  padding: 2rem;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.features-box h2 {
  color: var(--color-heading);
  margin-bottom: 0.75rem;
  font-size: 1.75rem;
}

.features-subtitle {
  color: var(--color-heading);
  font-weight: 500;
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.contact-text {
  margin-top: 1.5rem;
  color: var(--color-text);
  font-size: 0.9rem;
  font-style: italic;
  text-align: center;
}

.features-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature-item {
  display: flex;
  flex-direction: column;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--color-border);
}

.feature-item:last-child {
  border-bottom: none;
}

.feature-short {
  font-weight: 600;
  color: var(--color-heading);
  margin-bottom: 0.25rem;
}

.feature-long {
  color: var(--color-text);
  font-size: 0.9rem;
  line-height: 1.4;
}

.login-box {
  height: 100%;
  width: fit-content;
  padding: 2rem;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.login-box h1 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: var(--color-heading);
}

form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-weight: 500;
}

input {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

button {
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  background-color: var(--color-primary);
  color: white;
  cursor: pointer;
  transition: background-color 0.2s;
}

button:hover {
  background-color: var(--color-primary-hover);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  color: var(--color-text);
  margin: 0.5rem 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--color-border);
}

.divider span {
  padding: 0 10px;
}

.google-btn {
  background-color: white;
  color: #333;
  border: 1px solid var(--color-primary);
  box-shadow: 0 2px 4px var(--color-shadow);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.google-btn:hover {
  background-color: color-mix(in srgb, white, var(--color-shadow));
}

.google-icon {
  width: 22px;
  height: 22px;
}

.error-message {
  color: var(--color-error);
  text-align: center;
}

.success-message {
  color: var(--color-success);
  text-align: center;
}

.switch-form {
  text-align: center;
}

.switch-form a {
  color: var(--color-text);
  text-decoration: underline;
  cursor: pointer;
}

@media (max-width: 950px) {
  .login-view {
    flex-direction: column-reverse;
    align-items: center;
  }
}

.login-links {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: center;
}

.text-btn {
  background: none;
  border: none;
  color: var(--color-text);
  padding: 0;
  font-size: 1rem;
  text-decoration: underline;
  cursor: pointer;
  width: auto;
  border-radius: 0;
}

.text-btn:hover {
  color: var(--color-primary);
  background: none;
}

.text-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  text-decoration: none;
}
</style>
