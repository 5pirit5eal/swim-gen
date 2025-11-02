import { createRouter, createWebHistory } from 'vue-router'
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
  ],
})

router.beforeEach(async (to, from) => {
  const auth = useAuthStore()
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  if (requiresAuth || to.name === 'login') await auth.getUser()

  if (requiresAuth && !auth.user) {
    if (from.name !== 'login') {
      return {
        name: 'login',
        query: {
          redirectTo: to.fullPath,
        },
      }
    } else {
      return false
    }
  } else if (to.name === 'login' && auth.user) {
    return { name: 'home' }
  }
})

export default router
