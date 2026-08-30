import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import DrillsView from '../DrillsView.vue'
import DrillList from '@/components/drills/DrillList.vue'

describe('DrillsView.vue', () => {
  it('mounts DrillList in full mode and scrolls to top', () => {
    const scrollToSpy = vi.spyOn(window, 'scrollTo').mockImplementation(() => {})

    const wrapper = mount(DrillsView, {
      global: {
        stubs: {
          DrillList: {
            name: 'DrillList',
            props: ['featuredMode'],
            template: '<div class="stubbed-drill-list" :data-featured="featuredMode"></div>',
          },
        },
      },
    })

    expect(wrapper.find('.drills-view').exists()).toBe(true)
    const drillList = wrapper.findComponent(DrillList)
    expect(drillList.exists()).toBe(true)
    expect(drillList.props('featuredMode')).toBe(false)
    expect(scrollToSpy).toHaveBeenCalledWith(0, 0)
  })
})
