<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue3-toastify'

const { t } = useI18n()
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const username = ref('')
const loading = ref(false)

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
      toast.success(t('login.userExistsLoginSuccess'))
      router.push('/')
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
          <input id="email" type="email" :placeholder="t('login.email')" v-model="email" required />
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
          <router-link v-else to="/login?register=true">{{ t('login.needAccount') }}</router-link>
        </div>
        <button type="submit" :disabled="!canSubmit || loading">
          {{ loading ? t('login.loading') : isSignUp ? t('login.signUp') : t('login.login') }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 0.25rem 0 2rem 0;
}

.login-box {
  max-width: 1080px;
  margin: 2rem auto;
  padding: 2rem;
  background-color: var(--color-background-soft);
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
  border-radius: 0.375rem;
}

button {
  padding: 0.75rem;
  border: none;
  border-radius: 0.375rem;
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
</style>
