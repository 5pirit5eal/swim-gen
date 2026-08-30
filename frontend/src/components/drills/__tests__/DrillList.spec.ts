import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { createI18n } from 'vue-i18n'
import DrillList from '../DrillList.vue'
import en from '@/locales/en.json'
import type { Drill } from '@/types'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: { en },
})

const mockDrills: Drill[] = Array.from({ length: 6 }, (_, i) => ({
  slug: `drill-${i + 1}`,
  targets: ['Technique'],
  short_description: `Short description for drill ${i + 1}`,
  img_name: `drill-${i + 1}.webp`,
  img_description: `Drill ${i + 1} illustration`,
  title: `Drill ${i + 1} Title`,
  description: [`Full description for drill ${i + 1}`],
  video_url: [],
  styles: ['Freestyle'],
  difficulty: 'Easy',
  target_groups: ['Beginner'],
}))

describe('DrillList.vue', () => {
  it('renders featured mode with limited drills and browse-all links', async () => {
    const wrapper = mount(DrillList, {
      props: {
        featuredMode: true,
        limit: 4,
      },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              drills: {
                searchResults: mockDrills,
                searchTotal: 6,
                isLoading: false,
                error: null,
                searchParams: { page: 1, limit: 12 },
              },
            },
          }),
          i18n,
        ],
        stubs: {
          DrillCard: {
            props: ['drill'],
            template: '<div class="stubbed-drill-card">{{ drill.title }}</div>',
          },
          DrillFilter: true,
          RouterLink: {
            props: ['to'],
            template: '<a :href="to"><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.find('.featured-header').exists()).toBe(true)
    expect(wrapper.find('.drill-list-header').exists()).toBe(false)
    expect(wrapper.findComponent({ name: 'DrillFilter' }).exists()).toBe(false)

    // Only 4 items rendered
    const cards = wrapper.findAll('.stubbed-drill-card')
    expect(cards.length).toBe(4)

    // Contains browse links to /drills
    const links = wrapper.findAll('a')
    expect(links.some((l) => l.attributes('href') === '/drills')).toBe(true)
  })

  it('renders full mode with filters, full items, and pagination', async () => {
    const wrapper = mount(DrillList, {
      props: {
        featuredMode: false,
      },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              drills: {
                searchResults: mockDrills,
                searchTotal: 24,
                isLoading: false,
                error: null,
                searchParams: { page: 1, limit: 12 },
              },
            },
          }),
          i18n,
        ],
        stubs: {
          DrillCard: {
            props: ['drill'],
            template: '<div class="stubbed-drill-card">{{ drill.title }}</div>',
          },
          DrillFilter: {
            template: '<div class="stubbed-drill-filter">Filter</div>',
          },
          RouterLink: true,
        },
      },
    })

    expect(wrapper.find('.featured-header').exists()).toBe(false)
    expect(wrapper.find('.drill-list-header').exists()).toBe(true)
    expect(wrapper.find('.stubbed-drill-filter').exists()).toBe(true)

    const cards = wrapper.findAll('.stubbed-drill-card')
    expect(cards.length).toBe(6)

    // Pagination visible
    expect(wrapper.find('.pagination').exists()).toBe(true)
    expect(wrapper.find('.page-info').text()).toBe('1 / 2')
  })
})
