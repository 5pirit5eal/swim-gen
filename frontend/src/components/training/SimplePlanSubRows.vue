<script setup lang="ts">
import type { Row } from '@/types'
import ContentWithDrillLinks from '@/components/training/ContentWithDrillLinks.vue'

defineProps<{
  subRows: Row[]
  depth: number
}>()

function hasSubRows(row: Row): boolean {
  return !!row.SubRows && row.SubRows.length > 0
}

function hasEquipment(row: Row): boolean {
  return !!row.Equipment && row.Equipment.length > 0
}
</script>

<template>
  <div class="subrow-container" :class="`depth-indent-${depth}`">
    <table class="subrow-table">
      <tbody>
        <template v-for="(subRow, index) in subRows" :key="index">
          <tr class="subrow" data-testid="plan-card-nested">
            <td>{{ subRow.Amount }}</td>
            <td>{{ subRow.Multiplier }}</td>
            <td>{{ subRow.Distance }}</td>
            <td>{{ subRow.Break }}</td>
             <td class="content-cell">
               <ContentWithDrillLinks :content="subRow.Content" />
               <span v-if="hasEquipment(subRow)" class="equipment-badges" data-testid="plan-equipment">
                 <span
                   v-for="eq in subRow.Equipment"
                   :key="eq"
                   class="equipment-badge"
                 >{{ eq }}</span>
               </span>
             </td>
            <td class="intensity-cell">{{ subRow.Intensity }}</td>
            <td class="total-cell">{{ subRow.Sum }}</td>
          </tr>
          <!-- Recursively render deeper SubRows -->
          <tr v-if="hasSubRows(subRow)" class="nested-subrow-container-row">
            <td colspan="7" class="nested-subrow-container-cell">
              <SimplePlanSubRows :sub-rows="subRow.SubRows!" :depth="depth + 1" />
            </td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
export default {
  name: 'SimplePlanSubRows',
}
</script>

<style scoped>
.subrow-container {
  border-left: 2px solid var(--color-primary);
  margin-left: 0.75rem;
  background: var(--color-background);
}

.depth-indent-0 {
  margin-left: 0.75rem;
}

.depth-indent-1 {
  margin-left: 0.5rem;
}

.depth-indent-2 {
  margin-left: 0.35rem;
}

.depth-indent-3 {
  margin-left: 0.25rem;
}

.subrow-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8rem;
}

.subrow-table td {
  border: 1px solid var(--color-border);
  padding: 0.3rem 0.25rem;
  text-align: center;
  color: var(--color-text);
}

.subrow:nth-child(odd) {
  background-color: var(--color-background);
}

.subrow:nth-child(even) {
  background-color: var(--color-background-soft);
}

.content-cell {
  text-align: left;
  font-size: 0.75rem;
}

.intensity-cell {
  font-weight: 600;
  color: var(--color-primary);
}

.total-cell {
  font-weight: 600;
}

.nested-subrow-container-row {
  background: transparent;
}

.nested-subrow-container-cell {
  padding: 0 !important;
  border: none !important;
}

.equipment-badges {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 0.15rem;
  margin-left: 0.3rem;
}

.equipment-badge {
  display: inline-block;
  font-size: 0.55rem;
  font-weight: 600;
  text-transform: uppercase;
  padding: 0.05rem 0.25rem;
  border-radius: 3px;
  background: var(--color-primary);
  color: white;
  letter-spacing: 0.3px;
  white-space: nowrap;
}
</style>
