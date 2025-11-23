<script setup lang="ts">
import { ref, watch } from 'vue'
import { useExportStore } from '@/stores/export'
import IconDownload from '@/components/icons/IconDownload.vue'
import type { PlanToPDFRequest, PlanStore } from '@/types'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  store: PlanStore
}>()

const exportStore = useExportStore()
const { t } = useI18n()

const exportPhase = ref<'idle' | 'exporting' | 'done'>('idle')
const pdfUrl = ref<string | null>(null)
const exportHorizontal = ref(false)
const exportLargeFont = ref(false)
const isExportMenuOpen = ref(false)

// Reset export if plan changes (deep watch) or options change
watch(
  () => props.store.currentPlan,
  () => {
    resetExportState()
  },
  { deep: true },
)

// Reset export if options change
watch([exportHorizontal, exportLargeFont], () => {
  resetExportState()
})

// Utility to reset export state (re-used)
function resetExportState() {
  pdfUrl.value = null
  exportPhase.value = 'idle'
}

async function handleExport() {
  isExportMenuOpen.value = false // Close menu on export
  // Phase 2: user clicks "Open PDF"
  if (exportPhase.value === 'done' && pdfUrl.value) {
    const w = window.open(pdfUrl.value, '_blank')
    if (!w) window.location.href = pdfUrl.value
    return
  }

  // Prevent double starts
  if (exportPhase.value === 'exporting') return
  if (!props.store.currentPlan) return

  // Phase 1: user clicks "Export PDF"
  exportPhase.value = 'exporting'
  try {
    const payload: PlanToPDFRequest = {
      ...props.store.currentPlan,
      horizontal: exportHorizontal.value,
      large_font: exportLargeFont.value,
      language: navigator.language.split('-')[0] || 'en',
    }
    pdfUrl.value = await exportStore.exportToPDF(payload)
    if (!pdfUrl.value) {
      exportPhase.value = 'idle'
      return
    }
    exportPhase.value = 'done'
  } catch (e) {
    console.error('PDF export failed', e)
    exportPhase.value = 'idle'
  }
}
</script>

<template>
  <div class="export-actions">
    <button
      @click="handleExport"
      class="export-btn main-action"
      :disabled="exportPhase === 'exporting'"
    >
      <template v-if="exportPhase === 'exporting'">
        {{ t('display.exporting') }}
      </template>
      <template v-else-if="exportPhase === 'done'">
        <IconDownload class="icon" />
        {{ t('display.open_pdf') }}
      </template>
      <template v-else>
        <IconDownload class="icon" />
        {{ t('display.export_pdf') }}
      </template>
    </button>
    <div class="dropdown-container">
      <button
        @click="isExportMenuOpen = !isExportMenuOpen"
        class="export-btn dropdown-toggle"
        :disabled="exportPhase === 'exporting'"
      ></button>
      <Transition name="dropdown-transform">
        <div v-if="isExportMenuOpen" class="dropdown-menu">
          <label>
            <input type="checkbox" v-model="exportHorizontal" />
            {{ t('display.export_horizontal') }}
          </label>
          <label>
            <input type="checkbox" v-model="exportLargeFont" />
            {{ t('display.export_large_font') }}
          </label>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.export-btn {
  flex: 1;
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 160px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

@media (max-width: 740px) {
  .export-btn {
    width: 100%;
    min-width: 10%;
    padding: 0.5rem 1rem;
    overflow-wrap: break-word;
    font-size: 0.8rem;
  }
}

.export-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.export-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.icon {
  width: 24px;
  height: 24px;
}

.export-actions {
  display: flex;
  flex: 1;
  position: relative;
}

.export-actions .main-action {
  flex: 3;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  min-width: 140px;
}

.dropdown-container {
  flex: 1;
  display: flex;
  position: static;
  min-width: 20px;
}

.export-actions .dropdown-toggle {
  flex: 1;
  position: relative;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  border-left: 1px solid var(--color-primary-hover);
  padding: 0.75rem 1rem;
  min-width: 0;
}

.dropdown-toggle::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-left: 0.375rem solid transparent;
  border-right: 0.375rem solid transparent;
  border-top: 0.5rem solid white;
  transform: translate(-50%, -50%);
  transition: border-color 0.2s;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 0.25rem;
  padding: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 10;
  margin-top: 0.5rem;
}

.dropdown-menu label {
  display: block;
  padding: 0.5rem;
  cursor: pointer;
  color: var(--color-text);
  text-align: left;
}

.dropdown-menu label:hover {
  background-color: var(--color-background-mute);
}

.dropdown-menu input {
  margin-right: 0.5rem;
}

/* Dropdown Transition using transform */
.dropdown-transform-enter-active,
.dropdown-transform-leave-active {
  transition:
    opacity 0.2s ease-in-out,
    transform 0.2s ease-in-out;
  transform-origin: top;
}

.dropdown-transform-enter-from,
.dropdown-transform-leave-to {
  opacity: 0;
  transform: scaleY(0.9) translateY(-0.5rem);
}
</style>
