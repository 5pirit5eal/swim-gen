import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Sidebar from '../Sidebar.vue'
import { createTestingPinia } from '@pinia/testing'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import i18n from '@/plugins/i18n'

describe('Sidebar.vue', () => {
    beforeEach(() => {
        const pinia = createTestingPinia({
            createSpy: vi.fn,
        })

        const trainingPlanStore = useTrainingPlanStore(pinia)
        trainingPlanStore.generationHistory = [
            { plan_id: '1', title: 'Test Plan 1', description: 'Desc 1', table: [] },
        ]

        const sidebarStore = useSidebarStore(pinia)
        sidebarStore.isOpen = true
    })

    it('renders the sidebar when open', () => {
        const wrapper = mount(Sidebar, {
            global: {
                plugins: [i18n],
            },
        })

        expect(wrapper.find('.sidebar').classes()).toContain('is-open')
        expect(wrapper.find('h3').text()).toBe('History')
        expect(wrapper.find('.plan-list li').text()).toBe('Test Plan 1')
    })

    it('closes the sidebar when the close button is clicked', async () => {
        const wrapper = mount(Sidebar, {
            global: {
                plugins: [i18n],
            },
        })

        const sidebarStore = useSidebarStore()
        await wrapper.find('.close-btn').trigger('click')
        expect(sidebarStore.close).toHaveBeenCalledTimes(1)
    })

    it('loads a plan when a plan is clicked', async () => {
        const wrapper = mount(Sidebar, {
            global: {
                plugins: [i18n],
            },
        })

        const trainingPlanStore = useTrainingPlanStore()
        const sidebarStore = useSidebarStore()

        await wrapper.find('.plan-list li').trigger('click')
        expect(trainingPlanStore.loadPlanFromHistory).toHaveBeenCalledTimes(1)
        expect(sidebarStore.close).toHaveBeenCalledTimes(1)
    })
})
