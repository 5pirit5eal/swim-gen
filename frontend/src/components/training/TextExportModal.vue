<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseModal from '@/components/ui/BaseModal.vue'
import IconCopy from '@/components/icons/IconCopy.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import type { RAGResponse, Row } from '@/types'
import { EQUIPMENT_I18N_KEYS } from '@/utils/rowHelpers'
import { createBulletExport, createMarkdownExport, type TextExportLabels } from '@/utils/textExport'
import { toast } from 'vue3-toastify'

const props = defineProps<{ show: boolean; plan: RAGResponse }>()
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n()
const activeTab = ref<'markdown' | 'bullets'>('bullets')
const copied = ref(false)

function equipment(row: Row): string {
  return (row.Equipment ?? [])
    .map((item) => {
      const key = EQUIPMENT_I18N_KEYS[item as keyof typeof EQUIPMENT_I18N_KEYS]
      return key ? t(`equipment.${key}`) : item
    })
    .join(', ')
}

const labels = computed<TextExportLabels>(() => ({
  title: t('text_export.title'),
  description: t('display.coach_notes'),
  set: t('text_export.set'),
  exercise: t('display.content'),
  distance: t('display.distance'),
  break: t('display.break'),
  intensity: t('display.intensity'),
  equipment: t('display.equipment'),
  total: t('display.sum'),
  meters: 'm',
}))

const markdown = computed(() => createMarkdownExport(props.plan, labels.value, equipment))
const bullets = computed(() => createBulletExport(props.plan, labels.value, equipment))
const text = computed(() => (activeTab.value === 'markdown' ? markdown.value : bullets.value))

async function copyText() {
  try {
    await navigator.clipboard.writeText(text.value)
    copied.value = true
    toast.success(t('text_export.copied'))
    window.setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to copy text export:', error)
    toast.error(t('text_export.copy_error'))
  }
}

function close() {
  copied.value = false
  emit('close')
}
</script>

<template>
  <BaseModal :show="show" @close="close">
    <template #header>
      <div class="text-export-header">
        <div class="text-export-tabs" role="tablist">
          <button
            class="text-export-tab"
            :class="{ active: activeTab === 'bullets' }"
            role="tab"
            :aria-selected="activeTab === 'bullets'"
            @click="activeTab = 'bullets'"
          >
            {{ t('text_export.bullets') }}
          </button>
          <button
            class="text-export-tab"
            :class="{ active: activeTab === 'markdown' }"
            role="tab"
            :aria-selected="activeTab === 'markdown'"
            @click="activeTab = 'markdown'"
          >
            {{ t('text_export.markdown') }}
          </button>
        </div>
      </div>
    </template>
    <template #body>
      <div class="text-export-title-row">
        <h2 class="text-export-plan-title">{{ plan.title }}</h2>
        <button class="copy-text-button" :class="{ success: copied }" @click="copyText">
          <IconCheck v-if="copied" class="icon" />
          <IconCopy v-else class="icon" />
          {{ copied ? t('text_export.copied') : t('text_export.copy') }}
        </button>
      </div>
      <pre class="text-export-content">{{ text }}</pre>
    </template>
  </BaseModal>
</template>

<style scoped>
.text-export-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
  min-width: 0;
}

.text-export-tabs {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.text-export-tab {
  border: 1px solid var(--color-border);
  border-radius: 999px;
  padding: 0.6rem 1rem;
  background: var(--color-background);
  color: var(--color-text);
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 700;
  line-height: 1.2;
  white-space: nowrap;
  transition:
    background-color 0.2s,
    border-color 0.2s,
    color 0.2s,
    box-shadow 0.2s;
}

.text-export-tab:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-background-soft);
}

.text-export-tab.active {
  border-color: var(--color-primary);
  background: var(--color-primary);
  color: white;
  box-shadow: 0 2px 5px var(--color-shadow);
}

.text-export-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.copy-text-button {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  border: 0;
  border-radius: 6px;
  padding: 0.5rem 0.75rem;
  background: var(--color-primary);
  color: white;
  cursor: pointer;
  font-weight: 600;
}

.copy-text-button.success {
  background: var(--color-success);
}

.icon {
  width: 1.1rem;
  height: 1.1rem;
}

.text-export-plan-title {
  margin: 0;
  color: var(--color-heading);
  font-size: 1.1rem;
}

.text-export-content {
  margin: 0.75rem 0 0;
  padding: 1rem;
  height: 28rem;
  max-height: calc(90vh - 220px);
  min-height: 0;
  overflow: auto;
  white-space: pre;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: var(--color-background-soft);
  color: var(--color-text);
  font:
    0.9rem/1.5 ui-monospace,
    SFMono-Regular,
    Menlo,
    monospace;
}

@media (max-width: 500px) {
  .text-export-header {
    align-items: flex-start;
    flex-direction: column;
    gap: 0.5rem;
  }

  .text-export-tabs {
    width: 100%;
  }

  .text-export-tab {
    flex: 1;
    padding-inline: 0.4rem;
  }

  .text-export-title-row {
    align-items: flex-start;
    flex-direction: column;
  }

  .copy-text-button {
    align-self: flex-end;
  }
}
</style>
