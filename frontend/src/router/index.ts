import { createRouter, createWebHistory } from 'vue-router'
import { watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile/update-password',
      name: 'update-password',
      component: () => import('@/views/UpdatePasswordView.vue'),
    },
    {
      path: '/shared/:urlHash',
      name: 'shared',
      component: () => import('@/views/SharedView.vue'),
    },
    {
      path: '/shared',
      name: 'shared_empty',
      component: () => import('@/views/SharedView.vue'),
    },
    {
      path: '/plan/:id',
      name: 'plan',
      component: () => import('@/views/InteractionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/uploaded/:planId',
      name: 'uploaded',
      component: () => import('@/views/UploadedPlanView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to, from) => {
  const auth = useAuthStore()
  if (!auth.hasInitialized) {
    await new Promise<void>((resolve) => {
      const unwatch = watch(
        () => auth.hasInitialized,
        (isInitialized) => {
          if (isInitialized) {
            unwatch()
            resolve()
          }
        },
      )
    })
  }
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)

  if (requiresAuth && !auth.user) {
    if (from.name !== 'login') {
      console.log('User is not logged in, redirecting to home.')
      return {
        name: 'login',
        query: {
          redirectTo: to.fullPath,
        },
      }
    } else {
      console.log('User is not logged in, not redirecting.')
      return false
    }
  } else if (to.name === 'login' && auth.user) {
    console.log('User is already logged in, redirecting to home.')
    return { name: 'home' }
  }
})

export default router
