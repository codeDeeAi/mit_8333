<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import type { Category, Paginated, RequestStatus, Priority, ServiceRequest } from '@/types'
import { requestApi } from '@/api/requests'
import { categoryApi } from '@/api/categories'
import { ALL_STATUSES, ALL_PRIORITIES, STATUS_LABELS, PRIORITY_LABELS } from '@/lib/constants'
import RequestCard from './RequestCard.vue'
import BaseInput from '@/components/ui/BaseInput.vue'
import BaseSelect from '@/components/ui/BaseSelect.vue'
import BasePagination from '@/components/ui/BasePagination.vue'
import EmptyState from '@/components/ui/EmptyState.vue'

const props = defineProps<{
  scope?: { userId?: string; officerId?: string }
  detailBase?: string
}>()

const loading = ref(true)
const categories = ref<Category[]>([])
const result = ref<Paginated<ServiceRequest>>({ items: [], page: 1, pageSize: 6, total: 0 })

const filters = reactive({
  q: '',
  status: '' as RequestStatus | '',
  categoryId: '',
  priority: '' as Priority | '',
  page: 1,
})

const statusOptions = ALL_STATUSES.map((s) => ({ value: s, label: STATUS_LABELS[s] }))
const priorityOptions = ALL_PRIORITIES.map((p) => ({ value: p, label: PRIORITY_LABELS[p] }))

async function load() {
  loading.value = true
  try {
    result.value = await requestApi.list({ ...filters }, props.scope)
  } finally {
    loading.value = false
  }
}

let debounce: ReturnType<typeof setTimeout>
watch(
  () => filters.q,
  () => {
    clearTimeout(debounce)
    debounce = setTimeout(() => {
      filters.page = 1
      load()
    }, 300)
  },
)
watch([() => filters.status, () => filters.categoryId, () => filters.priority], () => {
  filters.page = 1
  load()
})
watch(() => filters.page, load)

onMounted(async () => {
  categories.value = await categoryApi.list()
  await load()
})

defineExpose({ reload: load })
</script>

<template>
  <div class="space-y-5">
    <!-- Filters -->
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <BaseInput v-model="filters.q" class="sm:col-span-2 lg:col-span-1" placeholder="Search requests…" />
      <BaseSelect v-model="filters.status" :options="statusOptions" placeholder="All statuses" />
      <BaseSelect
        v-model="filters.categoryId"
        :options="categories.map((c) => ({ value: c.id, label: c.name }))"
        placeholder="All categories"
      />
      <BaseSelect v-model="filters.priority" :options="priorityOptions" placeholder="All priorities" />
    </div>

    <!-- List -->
    <div v-if="loading" class="grid gap-4 sm:grid-cols-2">
      <div v-for="i in 4" :key="i" class="h-40 animate-pulse rounded-xl border border-border bg-surface" />
    </div>

    <EmptyState
      v-else-if="!result.items.length"
      title="No requests found"
      subtitle="Try adjusting your search or filters."
    />

    <template v-else>
      <div class="grid gap-4 sm:grid-cols-2">
        <div v-for="req in result.items" :key="req.id" class="space-y-2">
          <RequestCard
            :request="req"
            :to="`${detailBase ?? '/requests'}/${req.id}`"
          />
          <slot name="card-actions" :request="req" />
        </div>
      </div>
      <BasePagination
        :page="result.page"
        :page-size="result.pageSize"
        :total="result.total"
        @update:page="filters.page = $event"
      />
    </template>
  </div>
</template>
