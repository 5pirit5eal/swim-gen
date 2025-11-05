import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import ProfileView from '../ProfileView.vue'
import { useProfileStore } from '@/stores/profile'
import { describe, it, expect, vi } from 'vitest'
import en from '@/locales/en.json'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en,
  },
})

describe('ProfileView.vue', () => {
  it('renders the profile in display mode and toggles to edit mode', async () => {
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
    }

    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Profile')
    expect(wrapper.find('.edit-btn').exists()).toBe(true)
    expect(wrapper.find('.submit-btn').exists()).toBe(false)

    await wrapper.find('.edit-btn').trigger('click')

    expect(wrapper.find('.edit-btn').exists()).toBe(false)
    expect(wrapper.find('.submit-btn').exists()).toBe(true)

    const select = wrapper.find('select')
    expect(select.exists()).toBe(true)
    const options = select.findAll('option')
    expect(options.length).toBe(4) // 3 options + 1 empty
    const checkboxInputs = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxInputs.length).toBe(9)

    await wrapper.find('.submit-btn').trigger('click')

    expect(profileStore.updateProfile).toHaveBeenCalled()
  })
})
