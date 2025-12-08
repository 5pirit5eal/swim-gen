<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import type { EmailOtpType } from '@supabase/supabase-js'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)
const isSuccess = ref(false)

onMounted(async () => {
  const token_hash = route.query.token_hash as string
  const type = route.query.type as EmailOtpType

  if (token_hash && type) {
    try {
      await authStore.verifyOtp(token_hash, type)
    } catch {
      toast.error(t('errors.unknown_error'))
      router.push('/profile')
    }
  }
})

async function handleUpdatePassword() {
  error.value = ''
  success.value = ''

  if (password.value !== confirmPassword.value) {
    error.value = t('profile.passwords_do_not_match')
    return
  }

  loading.value = true
  try {
    await authStore.updatePassword(password.value)
    success.value = t('profile.password_updated')
    toast.success(t('profile.password_updated'))
    isSuccess.value = true
    setTimeout(() => {
      router.push('/profile')
    }, 2000)
  } catch (e: unknown) {
    if (e instanceof Error) {
      error.value = e.message
      toast.error(e.message)
    } else {
      error.value = t('errors.unknown_error')
      toast.error(t('errors.unknown_error'))
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="update-password-view">
    <div class="container">
      <div class="card">
        <h1>{{ t('profile.update_password_title') }}</h1>
        <form @submit.prevent="handleUpdatePassword">
          <div class="form-group">
            <label>{{ t('profile.new_password') }}</label>
            <input type="password" v-model="password" required minlength="6" />
          </div>
          <div class="form-group">
            <label>{{ t('profile.confirm_password') }}</label>
            <input type="password" v-model="confirmPassword" required minlength="6" />
          </div>
          <div v-if="error" class="error-message">{{ error }}</div>
          <div v-if="success" class="success-message">{{ success }}</div>
          <button type="submit" :disabled="loading" class="submit-btn">
            {{ loading ? t('common.loading') : t('profile.update_password_btn') }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.update-password-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 1rem;
}

.container {
  width: 100%;
  max-width: 400px;
}

.card {
  background: var(--color-background-soft);
  padding: 2rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  box-shadow: 0 2px 4px var(--color-shadow);
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--color-heading);
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  color: var(--color-heading);
  font-weight: 500;
}

input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-background);
  color: var(--color-text);
  font-size: 1rem;
}

input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 2px 4px var(--color-shadow);
}

.submit-btn {
  width: 100%;
  padding: 0.75rem;
  background: var(--color-primary);
  color: white;
  border: none;
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
  opacity: 0.7;
  cursor: not-allowed;
}

.error-message {
  color: var(--color-error);
  margin-bottom: 1rem;
  padding: 0.5rem;
  background: rgba(220, 38, 38, 0.1);
  border-radius: 4px;
  text-align: center;
}

.success-message {
  color: var(--color-success);
  /* Assuming variable, otherwise fallback to green */
  margin-bottom: 1rem;
  padding: 0.5rem;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 4px;
  text-align: center;
}
</style>
