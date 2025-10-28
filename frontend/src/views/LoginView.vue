<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const username = ref('')
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const isSignUp = computed(() => route.query.register === 'true')

const canSubmit = computed(() => {
  if (isSignUp.value) {
    return email.value !== '' && password.value !== '' && username.value !== ''
  } else {
    return email.value !== '' && password.value !== ''
  }
})

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''
  try {
    await auth.signInWithPassword(email.value, password.value)
    router.push('/')
  } catch (error) {
    if (error instanceof Error) {
      if (error.message.includes('invalid login')) {
        errorMsg.value = t('login.invalidLogin')
      } else {
        errorMsg.value = error.message
      }
    }
  } finally {
    loading.value = false
  }
}

async function handleSignUp() {
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''
  let response = null
  try {
    response = await auth.signUp(email.value, password.value, username.value)
    successMsg.value = t('login.checkEmail')
    console.log('Sign-up response:', JSON.stringify(response, null, 2))
  } catch (error) {
    if (error instanceof Error) {
      errorMsg.value = t('login.registrationFailed')
    }
  }
  loading.value = false

  if (response?.user?.identities?.length && response?.user?.identities?.length > 0) {
    loading.value = true
    try {
      await auth.signInWithPassword(email.value, password.value)
    } catch (error) {
      if (error instanceof Error) {
        errorMsg.value = t('login.registrationFailed')
      }
    } finally {
      successMsg.value = t('login.autoLogin')
      loading.value = false
    }

  }
}
</script>

<template>
  <div class="login-view">
    <div class="login-box">
      <h1>{{ isSignUp ? t('login.signUp') : t('login.login') }}</h1>
      <form @submit.prevent="isSignUp ? handleSignUp() : handleLogin()">
        <div class="form-group" v-if="isSignUp">
          <label for="username">{{ t('login.username') }}*</label>
          <input id="username" type="text" v-model="username" required />
        </div>
        <div class="form-group">
          <label for="email">{{ t('login.email') }}*</label>
          <input id="email" type="email" v-model="email" required />
        </div>
        <div class="form-group">
          <label for="password">{{ t('login.password') }}*</label>
          <input id="password" type="password" v-model="password" required />
        </div>
        <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
        <p v-if="successMsg" class="success-message">{{ successMsg }}</p>
        <div class="switch-form">
          <router-link v-if="isSignUp" to="/login">{{ t('login.haveAccount') }}</router-link>
          <router-link v-else to="/login?register=true">{{ t('login.needAccount') }}</router-link>
        </div>
        <button type="submit" :disabled="canSubmit === false || loading">
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
