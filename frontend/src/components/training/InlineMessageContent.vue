<script lang="ts">
import { computed, defineComponent, h, type VNodeChild } from 'vue'
import { RouterLink } from 'vue-router'
import { parseInlineMarkdown, type InlineSegment } from '@/utils/inlineMarkdown'

function renderSegments(segments: InlineSegment[]): VNodeChild[] {
  return segments.map((segment) => {
    if (segment.type === 'text') return segment.content

    if (segment.type === 'drill-link') {
      return h(RouterLink, { to: { name: 'drill', params: { id: segment.drillId } } }, () =>
        renderSegments(segment.children),
      )
    }

    if (segment.type === 'link') {
      return h('a', { href: segment.href }, renderSegments(segment.children))
    }

    return h(segment.type, renderSegments(segment.children))
  })
}

export default defineComponent({
  props: {
    content: { type: String, required: true },
  },
  setup(props) {
    const segments = computed(() => parseInlineMarkdown(props.content))
    // Text children are escaped by Vue; model output is never inserted as HTML.
    return () => h('span', { class: 'inline-message-content' }, renderSegments(segments.value))
  },
})
</script>

<style scoped>
.inline-message-content {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.inline-message-content :deep(strong),
.inline-message-content :deep(strong *) {
  font-weight: 700;
}

.inline-message-content :deep(em) {
  font-style: italic;
}

.inline-message-content :deep(a) {
  color: var(--color-heading);
  text-decoration: underline;
  text-underline-offset: 0.15em;
  border-radius: 2px;
}

.inline-message-content :deep(a:hover) {
  text-decoration-thickness: 2px;
}

.inline-message-content :deep(a:focus-visible) {
  outline: 2px solid var(--color-heading);
  outline-offset: 2px;
}
</style>
