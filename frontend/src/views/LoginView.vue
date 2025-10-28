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

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''
  try {
    await auth.signInWithPassword(email.value, password.value)
    router.push('/')
  } catch (error) {
    if (error instanceof Error) {
      errorMsg.value = error.message
    }
  } finally {
    loading.value = false
  }
}

async function handleSignUp() {
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await auth.signUp(email.value, password.value, username.value)
    successMsg.value = t('login.checkEmail')
  } catch (error) {
    if (error instanceof Error) {
      errorMsg.value = error.message
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-view">
    <div class="login-box">
      <h1>{{ isSignUp ? t('login.signUp') : t('login.login') }}</h1>
      <form @submit.prevent="isSignUp ? handleSignUp() : handleLogin()">
        <div class="form-group" v-if="isSignUp">
          <label for="username">{{ t('login.username') }}</label>
          <input id="username" type="text" v-model="username" required />
        </div>
        <div class="form-group">
          <label for="email">{{ t('login.email') }}</label>
          <input id="email" type="email" v-model="email" required />
        </div>
        <div class="form-group">
          <label for="password">{{ t('login.password') }}</label>
          <input id="password" type="password" v-model="password" required />
        </div>
        <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
        <p v-if="successMsg" class="success-message">{{ successMsg }}</p>
        <button type="submit" :disabled="loading">
          {{ loading ? t('login.loading') : isSignUp ? t('login.signUp') : t('login.login') }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.login-box {
  width: 300px;
  padding: 2rem;
  background-color: var(--color-background-soft);
  border-radius: 0.5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  margin-bottom: 1.5rem;
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
  background-color: var(--color-primary-dark);
}

button:disabled {
  background-color: var(--color-border);
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
</style>
