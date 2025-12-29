<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import IconSearch from '@/components/icons/IconSearch.vue'
import type { DrillSearchParams } from '@/types'
// import { DIFFICULTY_OPTIONS } from '@/types' // Assuming this exists or I'll define it locally if not exported properly

const props = defineProps<{
  // Initial values
  initialFilters: Partial<DrillSearchParams>
}>()

const emit = defineEmits<{
  (e: 'update:filters', filters: Partial<DrillSearchParams>): void
  (e: 'search'): void
}>()

const { t } = useI18n()

// Local state for filters
const searchQuery = ref(props.initialFilters.q || '')
const selectedDifficulty = ref(props.initialFilters.difficulty || '')
// const selectedStyles = ref<string[]>(props.initialFilters.styles || [])
// const selectedTargets = ref<string[]>(props.initialFilters.target_groups || [])

// Difficulty Options (Hardcoded for now based on known values or I could import if available)
const difficultyOptions = [
  { value: '', label: t('drill.all_difficulties', 'All Difficulties') },
  { value: 'easy', label: t('drill.difficulty.easy', 'Easy') },
  { value: 'medium', label: t('drill.difficulty.medium', 'Medium') },
  { value: 'hard', label: t('drill.difficulty.hard', 'Hard') },
]

// Debounce search
let debounceTimeout: ReturnType<typeof setTimeout> | null = null

function handleSearchInput() {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    emitUpdate()
  }, 500)
}

function handleFilterChange() {
  emitUpdate()
}

function emitUpdate() {
  emit('update:filters', {
    q: searchQuery.value,
    difficulty: selectedDifficulty.value || undefined,
    page: 1, // Reset to page 1 on filter change
  })
}

// Watch prop changes to update local state if needed (e.g. reset from parent)
watch(
  () => props.initialFilters,
  (newFilters) => {
    if (newFilters.q !== undefined && newFilters.q !== searchQuery.value) {
      searchQuery.value = newFilters.q
    }
    if (newFilters.difficulty !== undefined && newFilters.difficulty !== selectedDifficulty.value) {
      selectedDifficulty.value = newFilters.difficulty
    }
  },
  { deep: true }
)
</script>

<template>
  <div class="drill-filter">
    <!-- Search Bar -->
    <div class="search-bar">
      <IconSearch class="search-icon" />
      <input v-model="searchQuery" type="text" :placeholder="t('drill.search_placeholder', 'Search drills...')"
        class="search-input" @input="handleSearchInput" />
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

      <!-- Additional filters for Style/Target could go here -->
    </div>
  </div>
</template>

<style scoped>
.drill-filter {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  /* margin-bottom: 2rem; */
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
  transition: border-color 0.2s, box-shadow 0.2s;
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
  background-image: url("data:image/svg+xml;charset=US-ASCII,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%22292.4%22%20height%3D%22292.4%22%3E%3Cpath%20fill%3D%22%23999%22%20d%3D%22M287%2069.4a17.6%2017.6%200%200%200-13-5.4H18.4c-5%200-9.3%201.8-12.9%205.4A17.6%2017.6%200%200%200%200%2082.2c0%205%201.8%209.3%205.4%2012.9l128%20127.9c3.6%203.6%207.8%205.4%2012.8%205.4s9.2-1.8%2012.8-5.4L287%2095c3.5-3.5%205.4-7.8%205.4-12.8%200-5-1.9-9.2-5.5-12.8z%22%2F%3E%3C%2Fsvg%3E");
  background-repeat: no-repeat;
  background-position: right 0.7rem top 50%;
  background-size: 0.65rem auto;
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
