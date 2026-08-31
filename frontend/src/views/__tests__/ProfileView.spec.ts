import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import { createMemoryHistory, createRouter } from 'vue-router'
import ProfileView from '../ProfileView.vue'
import { useProfileStore } from '@/stores/profile'
import { describe, it, expect, vi } from 'vitest'
import en from '@/locales/en.json'

const { toastSuccess, toastWarning, toastError } = vi.hoisted(() => ({
  toastSuccess: vi.fn(),
  toastWarning: vi.fn(),
  toastError: vi.fn(),
}))

vi.mock('vue3-toastify', () => ({
  toast: {
    success: toastSuccess,
    warning: toastWarning,
    error: toastError,
  },
}))

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en,
  },
})

describe('ProfileView.vue', () => {
  it('renders editable profile fields and saves them', async () => {
    const wrapper = mount(ProfileView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
          }),
          i18n,
        ],
      },
    })

    const profileStore = useProfileStore()
    profileStore.profile = {
      user_id: '123',
      updated_at: new Date().toISOString(),
      username: 'testuser',
      experience: 'Beginner',
      preferred_language: 'en',
      preferred_strokes: ['Freestyle'],
      categories: ['Swimmer'],
      overall_generations: 10,
      monthly_generations: 5,
      exports: 2,
      css_200m_seconds: null,
      css_400m_seconds: null,
    }

    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Profile')
    expect(wrapper.find('.submit-btn').exists()).toBe(true)
    expect(wrapper.find('#css-200m').exists()).toBe(true)

    const select = wrapper.find('select')
    expect(select.exists()).toBe(true)
    const options = select.findAll('option')
    expect(options.length).toBe(4) // 3 options + 1 empty
    const checkboxInputs = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxInputs.length).toBe(10)

    await wrapper.find('#css-200m').setValue('3:20')
    await wrapper.find('#css-400m').setValue('7:00')
    await wrapper.find('.submit-btn').trigger('click')
    await wrapper.vm.$nextTick()

    expect(profileStore.updateProfile).toHaveBeenCalled()
    expect(wrapper.find('.submit-btn').text()).toContain('Saved')
    expect(toastSuccess).toHaveBeenCalledWith('Profile saved successfully.')
  })

  it('shows the unsaved changes warning when leaving after editing', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/profile', component: ProfileView },
        { path: '/', component: { template: '<div />' } },
      ],
    })
    await router.push('/profile')
    await router.isReady()

    const wrapper = mount(
      { template: '<router-view />' },
      {
        global: {
          plugins: [createTestingPinia({ createSpy: vi.fn }), i18n, router],
        },
      },
    )
    const profileStore = useProfileStore()
    profileStore.profile = {
      user_id: '123',
      updated_at: new Date().toISOString(),
      username: 'testuser',
      experience: 'Beginner',
      preferred_language: 'en',
      preferred_strokes: [],
      categories: [],
      overall_generations: 0,
      monthly_generations: 0,
      exports: 0,
      css_200m_seconds: null,
      css_400m_seconds: null,
    }

    await wrapper.vm.$nextTick()
    await wrapper.find('select').setValue('Advanced')
    await router.push('/')

    expect(toastWarning).toHaveBeenCalledWith('Your profile changes were not saved.')
  })
})
