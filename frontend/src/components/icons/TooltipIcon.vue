<script setup lang="ts">
import { ref } from 'vue'

const showTooltip = ref(false)
</script>

<template>
  <span class="tooltip-container" @mouseenter="showTooltip = true" @mouseleave="showTooltip = false">
    <span class="tooltip-icon">
      <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="currentColor">
        <path
          d="M478-240q21 0 35.5-14.5T528-290q0-21-14.5-35.5T478-340q-21 0-35.5 14.5T428-290q0 21 14.5 35.5T478-240Zm-36-154h74q0-33 7.5-52t42.5-52q26-26 41-49.5t15-56.5q0-56-41-86t-97-30q-57 0-92.5 30T342-618l66 26q5-18 22.5-39t53.5-21q32 0 48 17.5t16 38.5q0 20-12 37.5T506-526q-44 39-54 59t-10 73Zm38 314q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Zm0-320Z" />
      </svg>
    </span>
    <div v-if="showTooltip" class="tooltip-text">
      <slot name="tooltip">
        A helpful tooltip with additional information.
      </slot>
    </div>
  </span>
</template>

<style scoped>
.tooltip-container {
  position: relative;
  display: inline-flex;
  /* Use inline-flex to align icon and text */
  align-items: center;
  cursor: help;
  margin-left: 0.15rem;
}

.tooltip-icon {
  display: flex;
  /* Use flex for centering SVG */
  align-items: center;
  justify-content: center;
  width: 1.2em;
  /* Adjust size as needed */
  height: 1.2em;
  border-radius: 50%;
  color: var(--color-text);
  font-size: 0.8em;
  /* Controls the size of the SVG relative to parent font-size */
}

.tooltip-icon svg {
  width: 100%;
  height: 100%;
  fill: var(--color-text);
}

.tooltip-text {
  visibility: hidden;
  width: 200px;
  background-color: var(--color-background-mute);
  color: var(--color-text);
  text-align: center;
  border-radius: 6px;
  padding: 0.5rem;
  position: absolute;
  z-index: 1;
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
  white-space: normal;
  /* Allow text to wrap */
}

.tooltip-container:hover .tooltip-text {
  visibility: visible;
  opacity: 1;
}
</style>
