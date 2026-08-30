<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useDrillsStore } from '@/stores/drills'
import { useI18n } from 'vue-i18n'
import DrillCard from '@/components/drills/DrillCard.vue'
import DrillFilter from '@/components/drills/DrillFilter.vue'
import IconArrowRight from '@/components/icons/IconArrowRight.vue'
import type { DrillSearchParams } from '@/types'

const props = withDefaults(
  defineProps<{
    featuredMode?: boolean
    limit?: number
  }>(),
  {
    featuredMode: false,
    limit: 4,
  },
)

const drillsStore = useDrillsStore()
const { searchResults, searchTotal, isLoading, searchParams, error } = storeToRefs(drillsStore)
const { t, locale } = useI18n()

// Initial fetch
onMounted(async () => {
  if (!props.featuredMode) {
    // Fetch filter options first in full mode
    await drillsStore.fetchFilterOptions(locale.value)
  }

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
const displayedDrills = computed(() => {
  if (props.featuredMode) {
    return searchResults.value.slice(0, props.limit)
  }
  return searchResults.value
})

const totalPages = computed(() => {
  return Math.ceil(searchTotal.value / searchParams.value.limit)
})
</script>

<template>
  <div
    class="drill-list-container"
    id="drill-list-section"
    :class="{ 'is-featured': featuredMode }"
  >
    <!-- Featured Mode Header -->
    <div v-if="featuredMode" class="featured-header">
      <div class="featured-title-group">
        <h2>{{ t('drill.featured_drills') }}</h2>
        <p class="featured-subtitle">{{ t('drill.featured_drills_subtitle') }}</p>
      </div>
      <router-link to="/drills" class="browse-all-btn btn-secondary">
        {{ t('drill.find_suitable_drills') }}
        <IconArrowRight class="icon-arrow-browse" aria-hidden="true" />
      </router-link>
    </div>

    <!-- Full Mode Header with Search and Filters -->
    <div v-else class="drill-list-header">
      <div class="header-row">
        <h2>
          {{ t('drill.explore_drills', 'Trainingsübungen') }}
          <span class="count-badge" v-if="searchTotal > 0">{{ searchTotal }}</span>
        </h2>
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
    <div v-else class="drill-grid" :class="{ 'featured-grid': featuredMode }">
      <DrillCard v-for="drill in displayedDrills" :key="drill.slug" :drill="drill" />
    </div>

    <!-- Featured Mode Bottom Action -->
    <div v-if="featuredMode && searchTotal > 0" class="featured-footer-action">
      <router-link to="/drills" class="browse-all-cta btn-primary">
        {{ t('drill.browse_all_drills', { count: searchTotal }) }}
        <IconArrowRight class="icon-arrow-browse" aria-hidden="true" />
      </router-link>
    </div>

    <!-- Full Mode Pagination -->
    <div v-if="!featuredMode && totalPages > 1" class="pagination">
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

.drill-list-container.is-featured {
  padding: 1.5rem 0 2rem 0;
}

.featured-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 1.5rem;
  background-color: var(--color-transparent);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid var(--color-border);
  box-shadow: 0 4px 16px var(--color-shadow);
  border-radius: 12px;
  padding: 1.5rem;
  margin: 2rem auto 1.5rem auto;
}

.featured-title-group h2 {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-heading);
  margin: 0 0 0.4rem 0;
}

.featured-subtitle {
  font-size: 0.95rem;
  color: var(--color-text);
  margin: 0;
  line-height: 1.5;
}

.browse-all-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  background-color: var(--color-background);
  color: var(--color-heading);
  border: 1px solid var(--color-border);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.95rem;
  text-decoration: none;
  white-space: nowrap;
  flex-shrink: 0;
  transition: all 0.2s;
}

.browse-all-btn:hover {
  background-color: var(--color-background-soft);
  border-color: var(--color-border-hover);
  color: var(--color-primary);
}

.browse-all-btn:focus-visible {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-shadow);
}

.browse-all-cta {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  background-color: var(--color-primary);
  color: white;
  border: 1px solid transparent;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  white-space: nowrap;
  box-shadow: 0 2px 6px var(--color-shadow);
  transition:
    background-color 0.2s,
    box-shadow 0.2s;
}

.browse-all-cta:hover {
  background-color: var(--color-primary-hover);
  box-shadow: 0 4px 12px var(--color-shadow);
}

.browse-all-cta:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--color-shadow);
}

.icon-arrow-browse {
  width: 16px;
  height: 16px;
}

.featured-footer-action {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.drill-list-header {
  background-color: var(--color-transparent);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid var(--color-border);
  box-shadow: 0 4px 16px var(--color-shadow);
  border-radius: 12px;
  padding: 1.5rem;
  margin: 1.5rem auto 2rem auto;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
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
  border-radius: 8px;
  vertical-align: middle;
}

.drill-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 2.5rem;
}

.drill-grid.featured-grid {
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 250px;
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

@media (max-width: 740px) {
  .featured-header {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
    padding: 1.25rem;
  }

  .browse-all-btn {
    justify-content: center;
    width: 100%;
  }

  .browse-all-cta {
    width: 100%;
  }

  .featured-title-group h2 {
    font-size: 1.4rem;
  }

  .drill-grid {
    grid-template-columns: 1fr;
    gap: 1.25rem;
  }
}
</style>
