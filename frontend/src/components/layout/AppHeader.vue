<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue3-toastify'
// Header component for the swim training plan generator

const { t } = useI18n()
const auth = useAuthStore()

async function handleLogout() {
  try {
    await auth.signOut()
    toast.success(t('login.logoutSuccess'))
  } catch (error) {
    console.error('Logout failed:', error)
    toast.error(t('login.unknownError'))
  }
}
</script>

<template>
  <header class="app-header">
    <div class="header-container">
      <router-link to="/" class="logo-link">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M19.5 11.195C20.9497 11.195 22.125 10.0198 22.125 8.57001C22.125 7.12026 20.9497 5.94501 19.5 5.94501C18.0503 5.94501 16.875 7.12026 16.875 8.57001C16.875 10.0198 18.0503 11.195 19.5 11.195Z"
              stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            <path
              d="M7.37891 14L11.2099 11.6L8.39991 8.39401C8.09795 8.04897 7.8791 7.6393 7.76016 7.19648C7.64123 6.75366 7.62538 6.28948 7.71384 5.83958C7.80229 5.38968 7.99269 4.96604 8.27042 4.60121C8.54814 4.23637 8.9058 3.94005 9.31591 3.73501L13.5789 1.60001C13.7553 1.51007 13.9477 1.45594 14.1451 1.44072C14.3424 1.4255 14.5409 1.4495 14.7289 1.51135C14.917 1.57319 15.091 1.67164 15.2408 1.80104C15.3906 1.93043 15.5133 2.08821 15.6019 2.26526C15.6905 2.44231 15.7431 2.63515 15.7568 2.83264C15.7705 3.03013 15.7449 3.22838 15.6816 3.41595C15.6183 3.60352 15.5185 3.77671 15.3879 3.92552C15.2574 4.07434 15.0986 4.19583 14.9209 4.28301L11.5209 5.98301C11.4189 6.03513 11.3302 6.10988 11.2616 6.20154C11.1929 6.29321 11.1461 6.39936 11.1248 6.51188C11.1035 6.6244 11.1082 6.74031 11.1385 6.85074C11.1689 6.96117 11.2241 7.06319 11.2999 7.14901L16.6498 13.263"
              stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            <path
              d="M23.25 17C21.75 17 21 18.5 19.5 18.5C18 18.5 17.25 17 15.75 17C14.25 17 13.5 18.5 12 18.5C10.5 18.5 9.74995 17 8.24995 17C6.74995 17 5.99995 18.5 4.49995 18.5C3.00911 18.5 2.25914 17.0183 0.777344 17.0002"
              stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            <path
              d="M0.75 21C2.25 21 3 22.5 4.5 22.5C6 22.5 6.75 21 8.25 21C9.75 21 10.5 22.5 12 22.5C13.5 22.5 14.25 21 15.75 21C17.25 21 18 22.5 19.5 22.5C20.9921 22.5 21.7421 21.0157 23.2265 21.0001"
              stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          <h1>{{ t('app.name') }}</h1>
          <!-- <span class="subtitle">Training Plan Generator</span> -->
        </div>
      </router-link>

      <!-- Navigation for future use (V2) -->
      <nav class="navigation" v-if="false">
        <router-link to="/" class="nav-link">{{ t('header.home') }}</router-link>
        <router-link to="/about" class="nav-link">{{ t('header.about') }}</router-link>
      </nav>

      <div class="auth-actions">
        <div v-if="!auth.user" class="login-actions">
          <router-link to="/login" custom v-slot="{ navigate }">
            <button @click="navigate" class="login-btn btn-secondary">
              {{ t('login.login') }}
            </button>
          </router-link>
          <router-link to="/login?register=true" custom v-slot="{ navigate }">
            <button @click="navigate" class="signup-btn btn-primary">
              {{ t('login.signUp') }}
            </button>
          </router-link>
        </div>
        <div v-else class="login-actions">
          <router-link to="/profile" custom v-slot="{ navigate }">
            <button @click="navigate" class="login-btn btn-secondary">
              {{ t('header.profile') }}
            </button>
          </router-link>
          <button @click="handleLogout" class="logout-btn">{{ t('login.logout') }}</button>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  background: transparent;
  padding: 1rem 0 0 0;
}

.header-container {
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  container-type: inline-size;
}

.logo-link {
  text-decoration: none;
}

.logo {
  display: flex;
  align-items: center;
  gap: 1rem;
  background-color: var(--color-transparent);
  backdrop-filter: blur(3px);
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
}

.logo-icon {
  width: 2rem;
  height: 2rem;
  color: var(--color-heading);
  background: transparent;
}

.logo h1 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--color-heading);
}

.subtitle {
  font-size: 0.875rem;
  color: var(--color-text-light);
}

.navigation {
  display: flex;
  gap: 1rem;
}

.nav-link {
  text-decoration: none;
  color: var(--color-text);
  font-weight: 500;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: var(--color-background-mute);
}

.nav-link.router-link-active {
  background-color: var(--color-background-mute);
  color: var(--color-text);
}

.auth-actions .login-actions {
  display: flex;
  gap: 0.5rem;
}

.auth-actions button {
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  background-color: var(--color-transparent);
  backdrop-filter: blur(2px);
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  margin-left: 0.5rem;
}

.auth-actions .btn-primary {
  border: none;
  background-color: var(--color-primary);
  color: white;
}

.auth-actions .btn-primary:hover {
  background-color: var(--color-primary-hover);
}

.auth-actions .btn-secondary {
  background-color: transparent;
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

.auth-actions .btn-secondary:hover {
  border: 1px solid var(--color-primary-hover);
  color: var(--color-primary-hover);
}

.auth-actions .logout-btn {
  background-color: transparent;
  color: var(--color-heading);
  border: 1px solid var(--color-heading);
  margin-left: unset;
}

.auth-actions .logout-btn:hover {
  color: var(--color-text);
  border: 1px solid var(--color-text);
}
</style>
