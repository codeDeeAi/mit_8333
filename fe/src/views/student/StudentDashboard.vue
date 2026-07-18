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
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
const loading = ref(true)
const items = ref<ServiceRequest[]>([])

const stats = computed(() => ({
  total: items.value.length,
  active: items.value.filter((r) => !['completed', 'rejected'].includes(r.status)).length,
  completed: items.value.filter((r) => r.status === 'completed').length,
}))
const recent = computed(() => items.value.slice(0, 4))

onMounted(async () => {
  const res = await requestApi.list({ pageSize: 100 }, { userId: auth.user!.id })
  items.value = res.items
  loading.value = false
})
</script>

<template>
  <div>
    <PageHeader
      :title="`Welcome, ${auth.user?.fullName.split(' ')[0]}`"
      subtitle="Track your maintenance requests and submit new ones."
    >
      <template #actions>
        <RouterLink to="/app/new">
          <BaseButton>+ New request</BaseButton>
        </RouterLink>
      </template>
    </PageHeader>

    <div class="mb-8 grid gap-4 sm:grid-cols-3">
      <StatCard label="Total requests" :value="stats.total" />
      <StatCard label="Active" :value="stats.active" accent="text-blue-400" />
      <StatCard label="Completed" :value="stats.completed" accent="text-emerald-400" />
    </div>

    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-fg">Recent requests</h2>
      <RouterLink to="/app/requests" class="text-sm text-accent hover:underline">View all</RouterLink>
    </div>

    <div v-if="loading" class="grid gap-4 sm:grid-cols-2">
      <div v-for="i in 2" :key="i" class="h-40 animate-pulse rounded-xl border border-border bg-surface" />
    </div>
    <EmptyState
      v-else-if="!recent.length"
      title="No requests yet"
      subtitle="Submit your first maintenance request to get started."
    >
      <RouterLink to="/app/new" class="mt-4">
        <BaseButton>Create a request</BaseButton>
      </RouterLink>
    </EmptyState>
    <div v-else class="grid gap-4 sm:grid-cols-2">
      <RequestCard v-for="req in recent" :key="req.id" :request="req" :to="`/requests/${req.id}`" />
    </div>
  </div>
</template>
