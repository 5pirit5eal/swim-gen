<script setup lang="ts">
import { nextTick, onUnmounted, ref } from 'vue'
import IconTooltip from '@/components/icons/IconTooltip.vue'
import { calculateOverlayPosition, type OverlayPlacement } from '@/utils/overlayPosition'

const showTooltip = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const tooltipRef = ref<HTMLElement | null>(null)
const tooltipPosition = ref({ left: 0, top: 0 })
const placement = ref<OverlayPlacement>('top')

async function updatePosition() {
  await nextTick()
  if (!containerRef.value || !tooltipRef.value) return

  const anchor = containerRef.value.getBoundingClientRect()
  const tooltip = tooltipRef.value.getBoundingClientRect()
  const position = calculateOverlayPosition(
    anchor,
    tooltip.width,
    tooltip.height,
    window.innerWidth,
    window.innerHeight,
  )

  tooltipPosition.value = { left: position.left, top: position.top }
  placement.value = position.placement
}

function show() {
  showTooltip.value = true
  updatePosition()
  window.addEventListener('resize', updatePosition)
  window.addEventListener('scroll', updatePosition, true)
}

function hide() {
  showTooltip.value = false
  window.removeEventListener('resize', updatePosition)
  window.removeEventListener('scroll', updatePosition, true)
}

onUnmounted(hide)
</script>

<template>
  <span ref="containerRef" class="tooltip-container" @mouseenter="show" @mouseleave="hide">
    <IconTooltip />
    <Teleport to="body">
      <div
        v-if="showTooltip"
        ref="tooltipRef"
        class="tooltip-text"
        :class="`position-${placement}`"
        :style="{ left: `${tooltipPosition.left}px`, top: `${tooltipPosition.top}px` }"
      >
        <slot name="tooltip"> A helpful tooltip with additional information. </slot>
      </div>
    </Teleport>
  </span>
</template>

<style scoped>
.tooltip-container {
  all: unset;
  position: relative;
  display: inline-flex;
  /* Use inline-flex to align icon and text */
  cursor: help;
  margin-left: 0.1rem;
  word-break: normal;
}

.tooltip-icon svg {
  width: 100%;
  height: 100%;
  fill: var(--color-text);
}

.tooltip-text {
  background-color: var(--color-background-mute);
  color: var(--color-text);
  text-align: start;
  text-transform: none;
  border-radius: 6px;
  padding: 0.5rem;
  position: fixed;
  z-index: 9999;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  font-size: 0.875rem;
  line-height: 1.4;
  text-wrap: wrap;
  overflow: auto;
  min-width: 300px;
  max-width: calc(100vw - 24px);
  max-height: calc(100vh - 24px);
  box-sizing: border-box;
  pointer-events: none;
  animation: tooltip-in 0.2s ease-out;
}

@media (max-width: 740px) {
  .tooltip-text {
    font-size: 0.75rem;
    padding: 0.25rem;
    white-space: normal;
    width: 200px;
    min-width: unset;
  }
}

@keyframes tooltip-in {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
}

.tooltip-text.position-top {
  transform-origin: bottom center;
}
</style>
