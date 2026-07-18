<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { reportApi, type SummaryStats } from '@/api/reports'
import { requestApi } from '@/api/requests'
import type { ServiceRequest } from '@/types'
import { STATUS_LABELS, STATUS_STYLES } from '@/lib/constants'
import PageHeader from '@/components/ui/PageHeader.vue'
import StatCard from '@/components/ui/StatCard.vue'
import BaseCard from '@/components/ui/BaseCard.vue'
import RequestCard from '@/components/requests/RequestCard.vue'

const loading = ref(true)
const summary = ref<SummaryStats | null>(null)
const recent = ref<ServiceRequest[]>([])

const maxCategory = ref(1)

onMounted(async () => {
  summary.value = await reportApi.summary()
  maxCategory.value = Math.max(1, ...summary.value.byCategory.map((c) => c.count))
  const res = await requestApi.list({ pageSize: 4 })
  recent.value = res.items
  loading.value = false
})
</script>

<template>
  <div>
    <PageHeader title="Overview" subtitle="System-wide maintenance activity at a glance." />

    <div v-if="summary" class="mb-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <StatCard label="Total requests" :value="summary.total" />
      <StatCard label="Pending" :value="summary.byStatus.pending" accent="text-yellow-400" />
      <StatCard label="In progress" :value="summary.byStatus.in_progress" accent="text-purple-400" />
      <StatCard label="Completed" :value="summary.byStatus.completed" accent="text-emerald-400" />
    </div>

    <div class="grid gap-6 lg:grid-cols-2">
      <!-- Status breakdown -->
      <BaseCard>
        <h3 class="mb-4 text-sm font-semibold text-fg">By status</h3>
        <div v-if="summary" class="space-y-3">
          <div v-for="(count, status) in summary.byStatus" :key="status" class="flex items-center gap-3">
            <span class="w-24 shrink-0 text-xs text-muted">{{ STATUS_LABELS[status] }}</span>
            <div class="h-2 flex-1 overflow-hidden rounded-full bg-surface-2">
              <div
                class="h-full rounded-full border"
                :class="STATUS_STYLES[status]"
                :style="{ width: `${summary.total ? (count / summary.total) * 100 : 0}%` }"
              />
            </div>
            <span class="w-6 text-right text-xs text-fg">{{ count }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Category breakdown -->
      <BaseCard>
        <h3 class="mb-4 text-sm font-semibold text-fg">By category</h3>
        <div v-if="summary" class="space-y-3">
          <div v-for="cat in summary.byCategory" :key="cat.name" class="flex items-center gap-3">
            <span class="w-32 shrink-0 truncate text-xs text-muted">{{ cat.name }}</span>
            <div class="h-2 flex-1 overflow-hidden rounded-full bg-surface-2">
              <div class="h-full rounded-full bg-accent" :style="{ width: `${(cat.count / maxCategory) * 100}%` }" />
            </div>
            <span class="w-6 text-right text-xs text-fg">{{ cat.count }}</span>
          </div>
        </div>
      </BaseCard>
    </div>

    <div class="mt-8 mb-4 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-fg">Latest requests</h2>
      <RouterLink to="/admin/requests" class="text-sm text-accent hover:underline">Manage all</RouterLink>
    </div>
    <div class="grid gap-4 sm:grid-cols-2">
      <RequestCard v-for="req in recent" :key="req.id" :request="req" :to="`/requests/${req.id}`" />
    </div>
  </div>
</template>
