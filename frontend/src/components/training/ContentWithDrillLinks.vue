<script setup lang="ts">
import { computed, watch } from 'vue'
import DrillLink from '@/components/drills/DrillLink.vue'
import { parseContentForDrillLinks, type ContentSegment } from '@/utils/markdownParser'

const props = defineProps<{
  content: string
}>()

const segments = computed<ContentSegment[]>(() => {
  console.debug(
    '[ContentWithDrillLinks] Computing segments for content:',
    JSON.stringify(props.content),
  )
  const result = parseContentForDrillLinks(props.content)
  console.debug('[ContentWithDrillLinks] Computed segments:', JSON.stringify(result))
  return result
})

// Debug watcher to see when content changes
watch(
  () => props.content,
  (newContent) => {
    console.debug('[ContentWithDrillLinks] Content prop changed to:', JSON.stringify(newContent))
  },
  { immediate: true },
)
</script>

<template>
  <span class="content-with-drill-links">
    <template v-for="(segment, index) in segments" :key="index">
      <span v-if="segment.type === 'text'">{{ segment.content }}</span>
      <DrillLink
        v-else-if="segment.type === 'drill-link'"
        :drill-id="segment.drillId"
        :text="segment.text"
      />
    </template>
  </span>
</template>

<style scoped>
.content-with-drill-links {
  display: inline;
}
</style>
