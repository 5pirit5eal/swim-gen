<script setup lang="ts">
import IconCross from '@/components/icons/IconCross.vue'
defineProps<{ show: boolean }>()
defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="modal-mask" @click="$emit('close')">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <slot name="header">
              <h1>Modal Title</h1>
            </slot>
            <button class="modal-close" @click="$emit('close')">
              <IconCross />
            </button>
          </div>
          <div class="modal-body">
            <slot name="body">
              <p>Modal content goes here</p>
            </slot>
          </div>
          <div class="modal-footer" v-if="$slots.footer">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: var(color-transparent);
  z-index: 9998;
  display: flex;
  justify-content: center;
  align-items: center;
  backdrop-filter: blur(3px);
}

.modal-container {
  border-radius: 8px;
  box-shadow: 0 2px 10px var(--color-shadow);
  border: 1px solid var(--color-border);
  width: 90%;
  max-width: 1000px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  background: var(--color-background-soft);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--color-border);
  border-radius: 8px 8px 0 0;
  background: var(--color-background-soft);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
  background: var(--color-background);
  color: var(--color-text);
  border-radius: 0 0 8px 8px;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--color-heading);
  transition: background-color 0.2s;
}

.modal-close:hover {
  color: var(--color-primary);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

/* Add responsive padding */
@media (max-width: 740px) {
  .modal-container {
    width: 95%;
    margin: 1rem;
  }

  .modal-header,
  .modal-body {
    padding: 1rem;
  }

  .modal-footer {
    padding: 0.75rem 1rem;
  }
}
</style>
