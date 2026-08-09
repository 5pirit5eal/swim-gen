<script setup lang="ts">
import { ref, watch } from 'vue'
import { useExportStore } from '@/stores/export'
import IconDownload from '@/components/icons/IconDownload.vue'
import IconCopy from '@/components/icons/IconCopy.vue'
import type { PlanToPDFRequest, PlanStore } from '@/types'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue3-toastify'
import { stripRowIds } from '@/utils/rowHelpers'
import TextExportModal from '@/components/training/TextExportModal.vue'
import { isIOS } from '@/utils/platform'

const props = defineProps<{
  store: PlanStore
}>()

const exportStore = useExportStore()
const { t } = useI18n()

const exportPhase = ref<'idle' | 'exporting' | 'done'>('idle')
const pdfUrl = ref<string | null>(null)
const exportHorizontal = ref(false)
const exportLargeFont = ref(false)
const omitTrainerNotes = ref(true)
const isExportMenuOpen = ref(false)
const isTextExportOpen = ref(false)

// Reset export if plan changes (deep watch) or options change
watch(
  () => props.store.currentPlan,
  () => {
    resetExportState()
  },
  { deep: true },
)

// Reset export if options change
watch([exportHorizontal, exportLargeFont, omitTrainerNotes], () => {
  resetExportState()
})

// Utility to reset export state (re-used)
function resetExportState() {
  pdfUrl.value = null
  exportPhase.value = 'idle'
}

function openPDF() {
  if (pdfUrl.value) {
    const w = window.open(pdfUrl.value, '_blank')
    if (!w) window.location.href = pdfUrl.value
  }
}

async function handlePDFExport() {
  isExportMenuOpen.value = false // Close menu on export
  // Phase 2: user clicks "Open PDF" (mostly for iOS or if auto-open fails)
  if (exportPhase.value === 'done' && pdfUrl.value) {
    openPDF()
    return
  }

  // Prevent double starts
  if (exportPhase.value === 'exporting') return
  if (!props.store.currentPlan) return

  // Phase 1: user clicks "Export PDF"
  exportPhase.value = 'exporting'
  try {
    // Strip _id from table rows (including nested SubRows) before sending to backend
    const tableWithoutIds = stripRowIds(props.store.currentPlan.table)

    const payload: PlanToPDFRequest = {
      plan_id: props.store.currentPlan.plan_id,
      title: props.store.currentPlan.title,
      description: omitTrainerNotes.value ? '' : props.store.currentPlan.description,
      table: tableWithoutIds,
      horizontal: exportHorizontal.value,
      large_font: exportLargeFont.value,
      language: navigator.language.split('-')[0] || 'en',
      frontend_base_url: window.location.origin,
    }
    pdfUrl.value = await exportStore.exportToPDF(payload)
    if (!pdfUrl.value) {
      exportPhase.value = 'idle'
      return
    }
    exportPhase.value = 'done'

    // Auto-open for non-iOS
    if (!isIOS()) {
      openPDF()
      exportPhase.value = 'idle' // On non-iOS, reset so user can re-export with different options; iOS keeps "done" so a second tap just re-opens the existing PDF
    }
  } catch (e) {
    console.error('PDF export failed', e)
    toast.error(t('errors.failed_to_export_plan'))
    exportPhase.value = 'idle'
  }
}

function toggleExportMenu() {
  isExportMenuOpen.value = !isExportMenuOpen.value
}

function openTextExport() {
  isExportMenuOpen.value = false
  if (props.store.currentPlan) isTextExportOpen.value = true
}
</script>

<template>
  <div class="export-actions">
    <button @click="toggleExportMenu" class="export-btn main-action" :disabled="exportPhase === 'exporting'">
      <IconDownload class="icon" />
      {{ t('display.export') }}
    </button>
    <div class="dropdown-container">
      <Transition name="dropdown-transform">
        <div v-if="isExportMenuOpen" class="dropdown-menu">
          <button class="dropdown-item dropdown-item-action" @click="handlePDFExport">
            <IconDownload class="menu-icon" />
            <span>
              {{
                exportPhase === 'exporting'
                  ? t('display.exporting')
                  : exportPhase === 'done'
                    ? t('display.open_pdf')
                    : t('display.export_pdf')
              }}
            </span>
          </button>
          <label class="pdf-option">
            <input type="checkbox" v-model="exportHorizontal" />
            {{ t('display.export_horizontal') }}
          </label>
          <label class="pdf-option">
            <input type="checkbox" v-model="exportLargeFont" />
            {{ t('display.export_large_font') }}
          </label>
          <label class="pdf-option">
            <input type="checkbox" v-model="omitTrainerNotes" />
            {{ t('display.export_omit_trainer_notes') }}
          </label>
          <button class="dropdown-item dropdown-item-action" @click="openTextExport">
            <IconCopy class="menu-icon" />
            <span>{{ t('text_export.export_text') }}</span>
          </button>
        </div>
      </Transition>
    </div>
    <TextExportModal v-if="store.currentPlan" :show="isTextExportOpen" :plan="store.currentPlan"
      @close="isTextExportOpen = false" />
  </div>
</template>

<style scoped>
.export-actions {
  display: flex;
  position: relative;
  width: fit-content;
  max-width: 200px;
}

.export-actions .main-action {
  position: relative;
  padding-right: 2.5rem;
}

.export-btn {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

@media (max-width: 740px) {
  .export-btn {
    min-width: 90px;
    padding: 0.25rem 0.5rem;
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

.dropdown-container {
  display: flex;
  position: static;
}

.main-action::after {
  content: '';
  position: absolute;
  top: 50%;
  right: 1rem;
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
  border-radius: 8px;
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

.dropdown-menu .pdf-option {
  padding-left: 2.25rem;
  position: relative;
}

.dropdown-menu .pdf-option::before {
  content: '';
  position: absolute;
  left: 1rem;
  top: 0;
  bottom: 0;
  border-left: 2px solid var(--color-border);
}

.dropdown-menu label:hover {
  background-color: var(--color-background-mute);
}

.dropdown-item {
  width: 100%;
  border: 0;
  background: transparent;
  color: var(--color-text);
  text-align: left;
  font: inherit;
}

.dropdown-item-action {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  cursor: pointer;
}

.dropdown-item-action:hover {
  background-color: var(--color-background-mute);
}

.menu-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
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
