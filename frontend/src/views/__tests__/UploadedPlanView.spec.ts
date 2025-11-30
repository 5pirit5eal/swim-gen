import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import UploadedPlanView from '../UploadedPlanView.vue'
import { useDonationStore } from '@/stores/uploads'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useRoute, useRouter } from 'vue-router'

// Mock vue-i18n
vi.mock('vue-i18n', () => ({
    useI18n: () => ({
        t: (key: string) => key,
    }),
}))

const pushMock = vi.fn()
vi.mock('vue-router', () => ({
    useRoute: vi.fn(),
    useRouter: vi.fn(() => ({
        push: pushMock,
    })),
}))

vi.mock('@/components/training/TrainingPlanDisplay.vue', () => ({
    default: { template: '<div class="plan-display-stub"></div>' },
}))

describe('UploadedPlanView.vue', () => {
    beforeEach(() => {
        vi.clearAllMocks()
            ; (useRoute as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
                params: { planId: 'uploaded-123' },
            })
        window.scrollTo = vi.fn()
    })

    it('fetches uploaded plan on mount', () => {
        mount(UploadedPlanView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                    }),
                ],
                stubs: {
                    IconSend: true,
                },
            },
        })

        const store = useDonationStore()
        expect(store.fetchUploadedPlan).toHaveBeenCalledWith('uploaded-123')
    })

    it('redirects if no plan found', async () => {
        const store = useDonationStore()
        store.fetchUploadedPlan = vi.fn().mockResolvedValue(false) // Failed to fetch
        store.currentPlan = null

        mount(UploadedPlanView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        stubActions: false,
                    }),
                ],
                stubs: {
                    IconSend: true,
                },
            },
        })

        // Wait for async onMounted
        await flushPromises()

        // Verify router push was called
        expect(pushMock).toHaveBeenCalledWith('/')
    })

    it('starts conversation', async () => {
        const wrapper = mount(UploadedPlanView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            donation: {
                                currentPlan: { plan_id: 'uploaded-123', title: 'Uploaded Plan' },
                            },
                        },
                    }),
                ],
                stubs: {
                    IconSend: true,
                },
            },
        })

        const trainingStore = useTrainingPlanStore()
        trainingStore.sendMessage = vi.fn().mockResolvedValue(true)
        const router = useRouter()

        const input = wrapper.find('.chat-input')
        await input.setValue('Hello')
        await wrapper.find('.chat-form').trigger('submit')

        expect(trainingStore.loadPlanFromHistory).toHaveBeenCalled()
        expect(trainingStore.sendMessage).toHaveBeenCalledWith('Hello')
        expect(pushMock).toHaveBeenCalledWith({ name: 'plan', params: { id: 'uploaded-123' } })
    })
})
