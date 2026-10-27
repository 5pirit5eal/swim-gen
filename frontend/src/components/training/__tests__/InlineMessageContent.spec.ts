import { describe, expect, it } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import InlineMessageContent from '../InlineMessageContent.vue'

function mountContent(content: string) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/drill/:id', name: 'drill', component: { template: '<div />' } },
    ],
  })
  return mount(InlineMessageContent, { props: { content }, global: { plugins: [router] } })
}

describe('InlineMessageContent', () => {
  it('renders drill links as real inline anchors without hover cards', async () => {
    const wrapper = mountContent('Try [Finger Drag](/drills/finger_drag.webp) now')
    const link = wrapper.get('a')

    expect(link.text()).toBe('Finger Drag')
    expect(link.attributes('href')).toBe('/drill/finger_drag')
    expect(wrapper.text()).toBe('Try Finger Drag now')
    await link.trigger('mouseenter')
    expect(wrapper.find('.drill-preview-card').exists()).toBe(false)
    expect(wrapper.find('.drill-link-wrapper').exists()).toBe(false)

    await link.trigger('click')
    await flushPromises()
    expect(wrapper.vm.$router.currentRoute.value.fullPath).toBe('/drill/finger_drag')
  })

  it('renders emphasis and nested links as semantic elements', () => {
    const wrapper = mountContent('**Fett** __auch fett__ *kursiv* _auch kursiv_ **[Drill](abc)**')

    expect(wrapper.findAll('strong').map((node) => node.text())).toEqual([
      'Fett',
      'auch fett',
      'Drill',
    ])
    expect(wrapper.findAll('em').map((node) => node.text())).toEqual(['kursiv', 'auch kursiv'])
    expect(wrapper.get('strong a').attributes('href')).toBe('/drill/abc')
  })

  it('renders safe external links', () => {
    const wrapper = mountContent('[Info](https://example.com)')
    expect(wrapper.get('a').attributes('href')).toBe('https://example.com')
  })

  it('never renders model HTML or executable links', () => {
    const content =
      '<img src=x onerror=alert(1)> **<script>alert(1)</script>** [Bad](javascript:/drills/bad)'
    const wrapper = mountContent(content)

    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.find('script').exists()).toBe(false)
    expect(wrapper.find('a').exists()).toBe(false)
    expect(wrapper.get('strong').text()).toBe('<script>alert(1)</script>')
    expect(wrapper.text()).toContain('<img src=x onerror=alert(1)>')
  })

  it('updates when message content changes', async () => {
    const wrapper = mountContent('Original\nmessage')
    expect(wrapper.text()).toBe('Original\nmessage')

    await wrapper.setProps({ content: '**Updated** [Drill](updated.webp)' })
    expect(wrapper.get('strong').text()).toBe('Updated')
    expect(wrapper.get('a').attributes('href')).toBe('/drill/updated')
  })
})
