<script setup lang="ts">
import { ref } from 'vue'
import TooltipIcon from '@/components/icons/TooltipIcon.vue'

const showTooltip = ref(false)
</script>

<template>
  <span
    class="tooltip-container"
    @mouseenter="showTooltip = true"
    @mouseleave="showTooltip = false"
  >
    <TooltipIcon />
    <div v-if="showTooltip" class="tooltip-text">
      <slot name="tooltip"> A helpful tooltip with additional information. </slot>
    </div>
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
  visibility: hidden;
  background-color: var(--color-background-mute);
  color: var(--color-text);
  text-align: start;
  text-transform: none;
  border-radius: 6px;
  padding: 0.5rem;
  position: absolute;
  z-index: 9999;
  bottom: 125%;
  /* Position the tooltip above the icon */
  left: 50%;
  transform: translateX(-50%);
  /* Center the tooltip precisely */
  opacity: 0;
  transition: opacity 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  font-size: 0.875rem;
  line-height: 1.4;
  text-wrap: wrap;
  overflow: auto;
}

@media (max-width: 740px) {
  .tooltip-text {
    font-size: 0.75rem;
    padding: 0.25rem;
    white-space: normal;
    width: 200px;
  }
}

.tooltip-container:hover .tooltip-text {
  visibility: visible;
  opacity: 1;
}
</style>
