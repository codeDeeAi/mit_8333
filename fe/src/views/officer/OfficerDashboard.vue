<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { requestApi } from '@/api/requests'
import type { ServiceRequest } from '@/types'
import PageHeader from '@/components/ui/PageHeader.vue'
import StatCard from '@/components/ui/StatCard.vue'
import RequestCard from '@/components/requests/RequestCard.vue'
import EmptyState from '@/components/ui/EmptyState.vue'

const auth = useAuthStore()
const loading = ref(true)
const items = ref<ServiceRequest[]>([])

const stats = computed(() => ({
  assigned: items.value.filter((r) => r.status === 'assigned').length,
  inProgress: items.value.filter((r) => r.status === 'in_progress').length,
  completed: items.value.filter((r) => r.status === 'completed').length,
}))
const active = computed(() =>
  items.value.filter((r) => !['completed', 'rejected'].includes(r.status)).slice(0, 4),
)

onMounted(async () => {
  const res = await requestApi.list({ pageSize: 100 }, { officerId: auth.user!.id })
  items.value = res.items
  loading.value = false
})
</script>

<template>
  <div>
    <PageHeader
      :title="`Hello, ${auth.user?.fullName.split(' ')[0]}`"
      subtitle="Your assigned jobs and their current progress."
    />

    <div class="mb-8 grid gap-4 sm:grid-cols-3">
      <StatCard label="Assigned" :value="stats.assigned" accent="text-blue-400" />
      <StatCard label="In progress" :value="stats.inProgress" accent="text-purple-400" />
      <StatCard label="Completed" :value="stats.completed" accent="text-emerald-400" />
    </div>

    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-fg">Active jobs</h2>
      <RouterLink to="/officer/assigned" class="text-sm text-accent hover:underline">View all</RouterLink>
    </div>

    <div v-if="loading" class="grid gap-4 sm:grid-cols-2">
      <div v-for="i in 2" :key="i" class="h-40 animate-pulse rounded-xl border border-border bg-surface" />
    </div>
    <EmptyState v-else-if="!active.length" title="No active jobs" subtitle="You're all caught up." />
    <div v-else class="grid gap-4 sm:grid-cols-2">
      <RequestCard v-for="req in active" :key="req.id" :request="req" :to="`/requests/${req.id}`" />
    </div>
  </div>
</template>
