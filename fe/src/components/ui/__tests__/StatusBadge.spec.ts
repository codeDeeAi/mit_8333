import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from '../StatusBadge.vue'

describe('StatusBadge', () => {
  it('renders the human label for a status', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'in_progress' } })
    expect(wrapper.text()).toContain('In Progress')
  })

  it('applies status-specific styling', () => {
    const wrapper = mount(StatusBadge, { props: { status: 'completed' } })
    expect(wrapper.classes().join(' ')).toContain('emerald')
  })
})
