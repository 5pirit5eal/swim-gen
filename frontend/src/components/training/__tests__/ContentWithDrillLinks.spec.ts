import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ContentWithDrillLinks from '../ContentWithDrillLinks.vue'

// Mock the DrillLink component to simplify testing
vi.mock('@/components/drills/DrillLink.vue', () => ({
  default: {
    name: 'DrillLink',
    props: ['drillId', 'text'],
    template: '<a class="drill-link" :data-drill-id="drillId">{{ text }}</a>',
  },
}))

describe('ContentWithDrillLinks', () => {
  it('renders plain text content without links', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: 'Simple text without any links',
      },
    })

    expect(wrapper.text()).toBe('Simple text without any links')
    expect(wrapper.findAll('.drill-link')).toHaveLength(0)
  })

  it('renders empty content', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: '',
      },
    })

    expect(wrapper.text()).toBe('')
  })

  it('renders a single drill link', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: '[Kraul Drill](/drills/kraul-drill-123)',
      },
    })

    const drillLink = wrapper.find('.drill-link')
    expect(drillLink.exists()).toBe(true)
    expect(drillLink.attributes('data-drill-id')).toBe('kraul-drill-123')
    expect(drillLink.text()).toBe('Kraul Drill')
  })

  it('renders drill link with surrounding text', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: 'Start with [Kraul Drill](/drills/abc) and continue swimming',
      },
    })

    expect(wrapper.text()).toContain('Start with')
    expect(wrapper.text()).toContain('and continue swimming')

    const drillLink = wrapper.find('.drill-link')
    expect(drillLink.exists()).toBe(true)
    expect(drillLink.attributes('data-drill-id')).toBe('abc')
    expect(drillLink.text()).toBe('Kraul Drill')
  })

  it('renders multiple drill links', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: '[First Drill](/drills/first) then [Second Drill](/drills/second)',
      },
    })

    const drillLinks = wrapper.findAll('.drill-link')
    expect(drillLinks).toHaveLength(2)

    const firstLink = drillLinks[0]
    const secondLink = drillLinks[1]
    expect(firstLink).toBeDefined()
    expect(secondLink).toBeDefined()
    expect(firstLink!.attributes('data-drill-id')).toBe('first')
    expect(firstLink!.text()).toBe('First Drill')
    expect(secondLink!.attributes('data-drill-id')).toBe('second')
    expect(secondLink!.text()).toBe('Second Drill')
  })

  it('does not render non-drill links as DrillLink components', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: 'Check [this link](https://example.com) out',
      },
    })

    expect(wrapper.findAll('.drill-link')).toHaveLength(0)
    // Non-drill markdown links should be rendered as plain text
    expect(wrapper.text()).toContain('[this link](https://example.com)')
  })

  it('handles mixed drill and non-drill links', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: '[Drill](/drills/abc) and [website](https://example.com)',
      },
    })

    const drillLinks = wrapper.findAll('.drill-link')
    expect(drillLinks).toHaveLength(1)

    const firstLink = drillLinks[0]
    expect(firstLink).toBeDefined()
    expect(firstLink!.attributes('data-drill-id')).toBe('abc')
    expect(wrapper.text()).toContain('[website](https://example.com)')
  })

  it('updates when content prop changes', async () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: 'Initial content',
      },
    })

    expect(wrapper.text()).toBe('Initial content')
    expect(wrapper.findAll('.drill-link')).toHaveLength(0)

    await wrapper.setProps({
      content: '[New Drill](/drills/new-drill)',
    })

    const drillLink = wrapper.find('.drill-link')
    expect(drillLink.exists()).toBe(true)
    expect(drillLink.attributes('data-drill-id')).toBe('new-drill')
  })

  it('handles content with special characters', () => {
    const wrapper = mount(ContentWithDrillLinks, {
      props: {
        content: '4x100m [Finger & Thumb Drill](/drills/finger-thumb) @ GA1',
      },
    })

    expect(wrapper.text()).toContain('4x100m')
    expect(wrapper.text()).toContain('@ GA1')

    const drillLink = wrapper.find('.drill-link')
    expect(drillLink.exists()).toBe(true)
    expect(drillLink.text()).toBe('Finger & Thumb Drill')
  })
})
