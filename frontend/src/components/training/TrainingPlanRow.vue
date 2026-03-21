<script setup lang="ts">
import BaseTableAction from '@/components/ui/BaseTableAction.vue'
import ContentWithDrillLinks from '@/components/training/ContentWithDrillLinks.vue'
import type { Row, PlanStore } from '@/types'
import { MAX_NESTING_DEPTH } from '@/utils/rowHelpers'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  row: Row
  path: number[]
  depth: number
  isEditing: boolean
  store: PlanStore
  isFirst: boolean
  isLast: boolean
}>()

const emit = defineEmits<{
  (e: 'start-editing', path: number[], field: keyof Row): void
  (e: 'stop-editing', event: Event, path: number[], field: keyof Row): void
  (e: 'auto-resize', event: Event): void
}>()

const { t } = useI18n()

const hasSubRows = computed(() => props.row.SubRows && props.row.SubRows.length > 0)
const canAddSubRow = computed(() => props.depth < MAX_NESTING_DEPTH)
const hasEquipment = computed(
  () => props.row.Equipment && props.row.Equipment.length > 0,
)

// Sub-row display: for a parent row, Distance is computed (sum of child sums),
// so we show it read-only. Only leaf rows allow Distance editing.
const isDistanceEditable = computed(() => !hasSubRows.value)

function handleAddRow() {
  props.store.addRow(props.path)
}

function handleRemoveRow() {
  props.store.removeRow(props.path)
}

function handleMoveRow(direction: 'up' | 'down') {
  props.store.moveRow(props.path, direction)
}

function handleAddSubRow() {
  props.store.addSubRow(props.path, props.depth)
}

function subRowPath(subIndex: number): number[] {
  return [...props.path, subIndex]
}
</script>

<template>
  <!-- Main row -->
  <tr
    class="exercise-row"
    :class="{
      'parent-row': hasSubRows,
      [`depth-${depth}`]: depth > 0,
    }"
    data-testid="plan-card"
  >
    <!-- Amount Cell -->
    <td
      @click="emit('start-editing', path, 'Amount')"
      class="anchor-cell"
      :class="{ 'nested-amount': depth > 0 }"
      data-testid="plan-row-actions"
    >
      <BaseTableAction
        v-if="isEditing && depth === 0"
        :is-first="isFirst"
        :is-last="isLast"
        :can-add-sub-row="canAddSubRow"
        @add="handleAddRow"
        @remove="handleRemoveRow"
        @move-up="handleMoveRow('up')"
        @move-down="handleMoveRow('down')"
        @add-subrow="handleAddSubRow"
      />
      <input
        type="text"
        inputmode="numeric"
        pattern="[0-9]*"
        v-if="isEditing"
        :value="row.Amount"
        @blur="emit('stop-editing', $event, path, 'Amount')"
        @keyup.enter="emit('stop-editing', $event, path, 'Amount')"
        class="editable-small"
      />
      <span v-else>{{ row.Amount }}</span>
    </td>
    <td>{{ row.Multiplier }}</td>
    <!-- Distance Cell -->
    <td @click="isDistanceEditable && emit('start-editing', path, 'Distance')">
      <input
        type="text"
        inputmode="numeric"
        pattern="[0-9]*"
        v-if="isEditing && isDistanceEditable"
        :value="row.Distance"
        @blur="emit('stop-editing', $event, path, 'Distance')"
        @keyup.enter="emit('stop-editing', $event, path, 'Distance')"
        class="editable-small"
      />
      <span v-else>{{ row.Distance }}</span>
    </td>
    <!-- Break Cell -->
    <td @click="emit('start-editing', path, 'Break')">
      <input
        type="text"
        v-if="isEditing"
        :value="row.Break"
        @blur="emit('stop-editing', $event, path, 'Break')"
        @keyup.enter="emit('stop-editing', $event, path, 'Break')"
        class="editable-small"
      />
      <span v-else>{{ row.Break }}</span>
    </td>
    <!-- Content Cell -->
    <td class="content-cell" @click="emit('start-editing', path, 'Content')">
      <textarea
        v-if="isEditing"
        :value="row.Content"
        @blur="emit('stop-editing', $event, path, 'Content')"
        @keyup.enter="emit('stop-editing', $event, path, 'Content')"
        @input="emit('auto-resize', $event)"
        class="editable-area"
      ></textarea>
      <span v-else>
        <ContentWithDrillLinks :content="row.Content" />
        <span v-if="hasEquipment" class="equipment-badges" data-testid="plan-equipment">
          <span
            v-for="eq in row.Equipment"
            :key="eq"
            class="equipment-badge"
          >{{ eq }}</span>
        </span>
      </span>
    </td>
    <!-- Intensity Cell -->
    <td class="intensity-cell" @click="emit('start-editing', path, 'Intensity')">
      <input
        type="text"
        v-if="isEditing"
        :value="row.Intensity"
        @blur="emit('stop-editing', $event, path, 'Intensity')"
        @keyup.enter="emit('stop-editing', $event, path, 'Intensity')"
        class="editable-small"
      />
      <span v-else>{{ row.Intensity }}</span>
    </td>
    <td class="total-cell">{{ row.Sum }}</td>
  </tr>

  <!-- SubRows: rendered as a nested area below the parent row -->
  <tr v-if="hasSubRows" class="subrow-container-row" data-testid="plan-card-nested">
    <td colspan="7" class="subrow-container-cell">
      <div class="subrow-container" :class="`depth-indent-${depth}`">
        <table class="subrow-table">
          <template v-for="(subRow, subIndex) in row.SubRows" :key="subRow._id || subIndex">
            <TrainingPlanRow
              :row="subRow"
              :path="subRowPath(subIndex)"
              :depth="depth + 1"
              :is-editing="isEditing"
              :store="store"
              :is-first="subIndex === 0"
              :is-last="subIndex === (row.SubRows?.length ?? 1) - 1"
              @start-editing="(p, f) => emit('start-editing', p, f)"
              @stop-editing="(ev, p, f) => emit('stop-editing', ev, p, f)"
              @auto-resize="(ev) => emit('auto-resize', ev)"
            />
          </template>
        </table>
        <!-- Add subrow button in edit mode (inside nested area) -->
        <button
          v-if="isEditing && canAddSubRow"
          class="add-subrow-inline"
          @click="handleAddSubRow"
          :title="t('display.add_subrow')"
        >
          + {{ t('display.add_subrow') }}
        </button>
      </div>
    </td>
  </tr>
</template>

<script lang="ts">
// Self-referencing for recursive component
export default {
  name: 'TrainingPlanRow',
}
</script>

<style scoped>
/* Parent row: slightly different background to show it has children */
.parent-row {
  background-color: var(--color-background-mute);
}

/* Nested indentation for subrow containers */
.subrow-container-cell {
  padding: 0 !important;
  border: none !important;
}

.subrow-container {
  border-left: 2px solid var(--color-primary);
  margin-left: 1rem;
  padding-left: 0;
  background: var(--color-background);
}

.depth-indent-0 {
  margin-left: 1rem;
}

.depth-indent-1 {
  margin-left: 0.75rem;
}

.depth-indent-2 {
  margin-left: 0.5rem;
}

.depth-indent-3 {
  margin-left: 0.25rem;
}

.subrow-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.subrow-table td {
  border: 1px solid var(--color-border);
  padding: 0.5rem 0.4rem;
  text-align: center;
  color: var(--color-heading);
}

/* Alternate SubRow backgrounds */
.subrow-table .exercise-row:nth-child(odd) {
  background-color: var(--color-background);
}

.subrow-table .exercise-row:nth-child(even) {
  background-color: var(--color-background-soft);
}

.subrow-table .exercise-row:hover {
  background-color: var(--color-background-mute);
}

/* SubRow container row itself should not have visible styling */
.subrow-container-row {
  background: transparent !important;
}

.subrow-container-row:hover {
  background: transparent !important;
}

/* Equipment badges */
.equipment-badges {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  margin-left: 0.5rem;
}

.equipment-badge {
  display: inline-block;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  background: var(--color-primary);
  color: white;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

/* Content cell */
.content-cell {
  text-align: left;
  font-weight: 500;
}

.intensity-cell {
  font-weight: 600;
  color: var(--color-primary);
}

.total-cell {
  font-weight: 600;
}

/* Anchor cell for action container */
.anchor-cell {
  position: relative;
  border-left: none;
}

/* Show action container on row hover */
.exercise-row:hover .anchor-cell :deep(.action-container) {
  opacity: 1;
  transform: translateX(0);
}

/* Editable inputs */
.editable-area {
  width: 100%;
  padding: 0.25rem;
  border: 1px solid var(--color-shadow);
  border-radius: 8px;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
}

.editable-small {
  width: 70%;
  text-align: center;
  border: 1px solid var(--color-shadow);
  border-radius: 8px;
  background-color: var(--color-background);
  color: var(--color-text);
  font-family: inherit;
  font-size: inherit;
  box-sizing: border-box;
}

.editable-area:focus,
.editable-small:focus {
  outline: 2px solid var(--color-primary);
}

/* Inline add subrow button */
.add-subrow-inline {
  display: block;
  width: 100%;
  padding: 0.3rem 0.75rem;
  border: none;
  border-top: 1px dashed var(--color-border);
  background: transparent;
  color: var(--color-primary);
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.2s;
}

.add-subrow-inline:hover {
  background-color: var(--color-background-soft);
}

@media (max-width: 740px) {
  .subrow-container {
    margin-left: 0.5rem;
  }

  .subrow-table td {
    padding: 0.25rem 0.2rem;
    font-size: 0.75rem;
  }

  .equipment-badge {
    font-size: 0.55rem;
    padding: 0.05rem 0.25rem;
  }
}
</style>
