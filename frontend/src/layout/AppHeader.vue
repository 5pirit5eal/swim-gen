<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useSidebarStore } from '@/stores/sidebar'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { toast } from 'vue3-toastify'
import { useRouter } from 'vue-router'
import IconMenu from '@/components/icons/IconMenu.vue'
import IconLogo from '@/components/icons/IconLogo.vue'
// Header component for the swim training plan generator

const { t } = useI18n()
const auth = useAuthStore()
const router = useRouter()
const sidebarStore = useSidebarStore()
const trainingStore = useTrainingPlanStore()

async function handleLogout() {
  try {
    await auth.signOut()
    toast.success(t('login.logoutSuccess'))
    trainingStore.clear()
    router.push('/')
    sidebarStore.close()
  } catch (error) {
    console.error('Logout failed:', error)
    toast.error(t('login.unknownError'))
  }
}
</script>

<template>
  <header class="app-header">
    <div class="header-container">
      <div class="logo" :class="{ 'sidebar-open': sidebarStore.isOpen }">
        <button v-if="auth.user" @click="sidebarStore.toggle" class="sidebar-toggle-btn" id="tutorial-sidebar-toggle">
          <IconMenu />
        </button>
        <router-link to="/" class="logo-link">
          <IconLogo />
          <h1>{{ t('app.name') }}</h1>
        </router-link>
      </div>

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
            <button @click="navigate" class="login-btn btn-secondary" id="tutorial-profile-btn">
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

.sidebar-toggle-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-heading);
  width: 2rem;
  margin-right: 1rem;
  opacity: 1;
  flex-shrink: 0;
  transition:
    opacity 0.3s ease,
    width 0.3s ease,
    margin-right 0.3s ease;
  overflow: hidden;
}

.sidebar-toggle-btn:hover {
  color: var(--color-primary);
}

.logo.sidebar-open .sidebar-toggle-btn {
  opacity: 0;
  width: 0;
  margin-right: 0;
  pointer-events: none;
}

.logo {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
  background-color: var(--color-transparent);
  backdrop-filter: blur(3px);
  padding: 0.25rem 0.75rem;
  border-radius: 8px;
  transition:
    padding 0.3s ease,
    gap 0.3s ease;
}

.logo-link {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 1rem;
  text-decoration: none;
}

.logo-icon {
  width: 38px;
  height: 38px;
  color: var(--color-heading);
  background: transparent;
  stroke-width: 6;
}

.logo h1 {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 500;
  color: var(--color-heading);
  text-wrap: nowrap;
}

.logo-link:hover .logo-icon,
.logo-link:hover h1 {
  color: var(--color-primary);
}

.auth-actions .login-actions {
  display: flex;
  gap: 0.5rem;
}

.auth-actions button {
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  font-size: 1rem;
  background-color: var(--color-transparent);
  backdrop-filter: blur(2px);
  padding: 0.5rem 1rem;
  border-radius: 8px;
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
  color: var(--color-heading);
  border: 1px solid var(--color-text);
}

.auth-actions .btn-secondary:hover {
  border: 1px solid var(--color-primary-hover);
  color: var(--color-primary);
}

.auth-actions .logout-btn {
  color: var(--color-heading);
  border: 1px solid var(--color-text);
}

.auth-actions .logout-btn:hover {
  color: var(--color-error);
  border: 1px solid var(--color-error);
}

@media (max-width: 740px) {
  .logo-icon {
    width: 28px;
    height: 28px;
  }

  .logo h1 {
    font-size: 1rem;
  }
}
</style>
