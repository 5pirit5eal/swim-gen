<script setup lang="ts">
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
            <button class="modal-close" @click="$emit('close')">×</button>
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
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9998;
  display: flex;
  justify-content: center;
  align-items: center;
  backdrop-filter: blur(2px);
}

.modal-container {
  border-radius: 0.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
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
  border-radius: 0.5rem 0.5rem 0 0;
  background: var(--color-background-soft);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
  background: var(--color-background);
  color: var(--color-text);
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
  color: var(--color-text);
  padding: 0.25rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s;
}

.modal-close:hover {
  background-color: var(--color-border);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-close:focus {
  outline: 2px solid var(--color-text);
  outline-offset: 2px;
}

.modal-footer button:focus {
  outline: 2px solid var(--color-text);
  outline-offset: 2px;
}

/* Add responsive padding */
@media (max-width: 480px) {
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
