<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useDrillsStore } from '@/stores/drills'
import { useI18n } from 'vue-i18n'
import DrillCard from '@/components/drills/DrillCard.vue'
import DrillFilter from '@/components/drills/DrillFilter.vue'
import IconArrowRight from '@/components/icons/IconArrowRight.vue'
import type { DrillSearchParams } from '@/types'

const drillsStore = useDrillsStore()
const { searchResults, searchTotal, isLoading, searchParams, error } = storeToRefs(drillsStore)
const { t, locale } = useI18n()

// Initial fetch
onMounted(async () => {
  // Fetch filter options first
  await drillsStore.fetchFilterOptions(locale.value)

  // If we still don't have results (or we just applied filters), fetch
  if (searchResults.value.length === 0) {
    drillsStore.searchDrills({ page: 1 })
  }
})

function handleFilterUpdate(newFilters: Partial<DrillSearchParams>) {
  drillsStore.searchDrills({ ...newFilters, page: 1 })
}

function handlePageChange(newPage: number) {
  if (newPage < 1 || newPage > totalPages.value) return
  drillsStore.searchDrills({ page: newPage })
  // Scroll to top of list
  const el = document.getElementById('drill-list-section')
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}

// Computeds
const totalPages = computed(() => {
  return Math.ceil(searchTotal.value / searchParams.value.limit)
})
</script>

<template>
  <div class="drill-list-container" id="drill-list-section">
    <div class="drill-list-header">
      <div class="header-row">
        <h2>
          {{ t('drill.explore_drills', 'Trainingsübungen') }}
          <span class="count-badge" v-if="searchTotal > 0">{{ searchTotal }}</span>
        </h2>
        <!-- Sorting could go here -->
        <!-- <div class="sort-options">Sortieren: Relevanz</div> -->
      </div>

      <DrillFilter :initial-filters="searchParams" @update:filters="handleFilterUpdate" />
    </div>
    <!-- Loading State -->
    <div v-if="isLoading && searchResults.length === 0" class="loading-state">
      <div class="loading-spinner"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="drillsStore.searchDrills({ page: 1 })" class="retry-btn">
        {{ t('common.retry', 'Retry') }}
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="!isLoading && searchResults.length === 0" class="empty-state">
      <p>{{ t('drill.no_results', 'No drills found matching your criteria.') }}</p>
    </div>

    <!-- Grid -->
    <div v-else class="drill-grid">
      <DrillCard v-for="drill in searchResults" :key="drill.slug" :drill="drill" />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        class="page-btn"
        :disabled="searchParams.page <= 1"
        @click="handlePageChange(searchParams.page - 1)"
        :aria-label="t('common.previous', 'Previous')"
      >
        <IconArrowRight class="icon-arrow-prev" />
      </button>

      <span class="page-info"> {{ searchParams.page }} / {{ totalPages }} </span>

      <button
        class="page-btn"
        :disabled="searchParams.page >= totalPages"
        @click="handlePageChange(searchParams.page + 1)"
        :aria-label="t('common.next', 'Next')"
      >
        <IconArrowRight class="icon-arrow-next" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.drill-list-container {
  padding: 2rem 0;
}

.drill-list-header {
  background-color: var(--color-transparent);
  backdrop-filter: blur(3px);
  border-radius: 8px;
  padding: 1rem 1rem 1rem 1rem;
  margin: 2rem auto;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header-row h2 {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-heading);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin: 0;
}

.count-badge {
  background: var(--color-primary);
  color: white;
  font-size: 0.9rem;
  padding: 2px 8px;
  border-radius: 12px;
  vertical-align: middle;
}

.drill-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  text-align: center;
  color: var(--color-text);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.retry-btn {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
  margin-top: 2rem;
}

.page-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 1px solid var(--color-border);
  background: var(--color-background);
  color: var(--color-heading);
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--color-background-soft);
}

.icon-arrow-prev,
.icon-arrow-next {
  width: 18px;
  height: 18px;
}

.icon-arrow-prev {
  transform: rotate(180deg);
}

.page-info {
  font-weight: 600;
  color: var(--color-text);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
