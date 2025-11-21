import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import SharedView from '../SharedView.vue'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import en from '@/locales/en.json'
import TrainingPlanDisplay from '@/components/training/TrainingPlanDisplay.vue'

const i18n = createI18n({
    legacy: false,
    locale: 'en',
    messages: {
        en,
    },
})

// Mock useRoute
vi.mock('vue-router', () => ({
    useRoute: vi.fn(() => ({
        params: {
            urlHash: 'test-hash',
        },
    })),
}))

describe('SharedView.vue', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('fetches the shared plan on mount', async () => {
        mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        const sharedPlanStore = useSharedPlanStore()

        expect(sharedPlanStore.fetchSharedPlanByHash).toHaveBeenCalledWith('test-hash')
    })

    it('displays loading state', async () => {
        const wrapper = mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            sharedPlan: {
                                isLoading: true,
                            },
                        },
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        expect(wrapper.find('.loading-state').exists()).toBe(true)
        expect(wrapper.find('.error-state').exists()).toBe(false)
        expect(wrapper.find('.no-plan').exists()).toBe(false)
    })

    it('displays error state', async () => {
        const wrapper = mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            sharedPlan: {
                                isLoading: false,
                                error: 'Test Error',
                            },
                        },
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        expect(wrapper.find('.error-state').text()).toBe('Test Error')
        expect(wrapper.find('.loading-state').exists()).toBe(false)
    })

    it('displays plan when loaded', async () => {
        const wrapper = mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            sharedPlan: {
                                isLoading: false,
                                error: null,
                                sharedPlan: {
                                    sharer_username: 'Test User',
                                    plan: { title: 'Test Plan' },
                                },
                            },
                        },
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        expect(wrapper.find('.shared-info').text()).toContain('Shared by: Test User')
        expect(wrapper.findComponent(TrainingPlanDisplay).exists()).toBe(true)
    })

    it('clears plan on unmount', async () => {
        const wrapper = mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        const sharedPlanStore = useSharedPlanStore()

        wrapper.unmount()

        expect(sharedPlanStore.clear).toHaveBeenCalled()
    })

    it('passes correct props to TrainingPlanDisplay', async () => {
        const wrapper = mount(SharedView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            sharedPlan: {
                                isLoading: false,
                                error: null,
                                sharedPlan: {
                                    sharer_username: 'Test User',
                                    plan: { title: 'Test Plan' },
                                },
                            },
                        },
                    }),
                    i18n,
                ],
                stubs: {
                    TrainingPlanDisplay: true
                }
            },
        })

        const displayComponent = wrapper.findComponent(TrainingPlanDisplay)
        expect(displayComponent.exists()).toBe(true)
        // Check if store prop is passed (it should be the sharedPlanStore instance)
        // Since we use createTestingPinia, useSharedPlanStore() returns the same store instance
        const sharedPlanStore = useSharedPlanStore()
        expect(displayComponent.props('store')).toBe(sharedPlanStore)
        expect(displayComponent.props('showShareButton')).toBe(false)
    })
})
