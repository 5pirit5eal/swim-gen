<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import IconSearch from '@/components/icons/IconSearch.vue'
import type { DrillSearchParams } from '@/types'
import { useDrillsStore } from '@/stores/drills'

const props = defineProps<{
  initialFilters: Partial<DrillSearchParams>
}>()

const emit = defineEmits<{
  (e: 'update:filters', filters: Partial<DrillSearchParams>): void
}>()

const { t } = useI18n()
const drillsStore = useDrillsStore()

// Local state for filters
const searchQuery = ref(props.initialFilters.q || '')
const selectedDifficulty = ref(props.initialFilters.difficulty || '')
const selectedStyles = ref<string[]>(props.initialFilters.styles || [])
const selectedTargetGroups = ref<string[]>(props.initialFilters.target_groups || [])
// const selectedTargets = ref<string[]>([]) // Not currently in searchParams but available in options

// Options from store
const difficultyOptions = computed(() => {
  const options = [{ value: '', label: t('drill.all_difficulties', 'All Difficulties') }]

  if (drillsStore.filterOptions?.difficulties) {
    drillsStore.filterOptions.difficulties.forEach((diff: string) => {
      options.push({ value: diff, label: diff }) // Backend returns localized strings
    })
  } else {
    // Fallback static options if needed, but backend should provide them
    options.push(
      { value: 'Easy', label: t('drill.difficulty.easy', 'Easy') },
      { value: 'Medium', label: t('drill.difficulty.medium', 'Medium') },
      { value: 'Hard', label: t('drill.difficulty.hard', 'Hard') },
    )
  }
  return options
})

const styleOptions = computed(() => drillsStore.filterOptions?.styles || [])
const targetGroupOptions = computed(() => drillsStore.filterOptions?.target_groups || [])

// Debounce search
let debounceTimeout: ReturnType<typeof setTimeout> | null = null

function handleSearchInput() {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    emitUpdate()
  }, 500)
}

function handleFilterChange() {
  // Styles and TargetGroups are arrays but we use single select for simplicity OR multi-select logic
  // For this UI, let's assume single select dropdowns that add to the array or just single value for now?
  // User request: "Add filters for target, target_group and style... allow querying... list should be cached"
  // Let's implement as multi-select enabled or just simple selects.
  // Given standard HTML selects, we'll do single selection that sets the array to [value].
  emitUpdate()
}

function emitUpdate() {
  emit('update:filters', {
    q: searchQuery.value,
    difficulty: selectedDifficulty.value || undefined,
    styles: selectedStyles.value.length > 0 ? selectedStyles.value : undefined,
    target_groups: selectedTargetGroups.value.length > 0 ? selectedTargetGroups.value : undefined,
    page: 1,
  })
}

// Watch prop changes
watch(
  () => props.initialFilters,
  (newFilters) => {
    if (newFilters.q !== undefined && newFilters.q !== searchQuery.value) {
      searchQuery.value = newFilters.q
    }
    if (newFilters.difficulty !== undefined && newFilters.difficulty !== selectedDifficulty.value) {
      selectedDifficulty.value = newFilters.difficulty
    }
    if (newFilters.styles) {
      selectedStyles.value = newFilters.styles
    }
    if (newFilters.target_groups) {
      selectedTargetGroups.value = newFilters.target_groups
    }
  },
  { deep: true },
)

// Helpers to handle array selection via single select (simple UI)
const currentStyle = computed({
  get: () => selectedStyles.value[0] || '',
  set: (val) => {
    selectedStyles.value = val ? [val] : []
  },
})

const currentTargetGroup = computed({
  get: () => selectedTargetGroups.value[0] || '',
  set: (val) => {
    selectedTargetGroups.value = val ? [val] : []
  },
})
</script>

<template>
  <div class="drill-filter">
    <!-- Search Bar -->
    <div class="search-bar">
      <IconSearch class="search-icon" />
      <input
        v-model="searchQuery"
        type="text"
        :placeholder="t('drill.search_placeholder', 'Search drills...')"
        class="search-input"
        @input="handleSearchInput"
      />
    </div>

    <!-- Filters -->
    <div class="filters-row">
      <!-- Difficulty Filter -->
      <div class="filter-group">
        <select v-model="selectedDifficulty" @change="handleFilterChange" class="filter-select">
          <option v-for="opt in difficultyOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <!-- Style Filter -->
      <div class="filter-group">
        <select v-model="currentStyle" @change="handleFilterChange" class="filter-select">
          <option value="">{{ t('drill.styles', 'Styles') }}</option>
          <option v-for="style in styleOptions" :key="style" :value="style">
            {{ style }}
          </option>
        </select>
      </div>

      <!-- Target Group Filter -->
      <div class="filter-group">
        <select v-model="currentTargetGroup" @change="handleFilterChange" class="filter-select">
          <option value="">{{ t('drill.target_groups', 'Target Groups') }}</option>
          <option v-for="tg in targetGroupOptions" :key="tg" :value="tg">
            {{ tg }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>

<style scoped>
.drill-filter {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.search-bar {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 1rem;
  width: 20px;
  height: 20px;
  color: var(--color-text);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3rem;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  color: var(--color-heading);
  font-size: 1rem;
  transition:
    border-color 0.2s,
    box-shadow 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-soft, rgba(0, 150, 255, 0.1));
}

.filters-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
}

.filter-select {
  padding: 0.5rem 2rem 0.5rem 1rem;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  color: var(--color-text);
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  appearance: none;
  background-image: url('data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23999%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E');
  background-repeat: no-repeat;
  background-position: right 0.7rem top 50%;
  background-size: 0.65rem auto;
  min-width: 150px;
}

.filter-select:focus {
  outline: none;
  border-color: var(--color-primary);
}

.filter-select:hover {
  background-color: var(--color-background-mute);
}

@media (min-width: 768px) {
  .drill-filter {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }

  .search-bar {
    max-width: 400px;
  }
}
</style>
