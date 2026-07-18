<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { ServiceRequest } from '@/types'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'
import { timeAgo } from '@/lib/format'

defineProps<{ request: ServiceRequest; to: string }>()
</script>

<template>
  <RouterLink
    :to="to"
    class="block rounded-xl border border-border bg-surface p-4 transition-colors hover:border-border-strong"
  >
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <h3 class="truncate text-sm font-semibold text-fg">{{ request.title }}</h3>
        <p class="mt-1 flex items-center gap-1.5 text-xs text-subtle">
          <span>{{ request.categoryName }}</span>
          <span aria-hidden="true">·</span>
          <span class="truncate">{{ request.location }}</span>
        </p>
      </div>
      <PriorityBadge :priority="request.priority" />
    </div>

    <p class="mt-3 line-clamp-2 text-sm text-muted">{{ request.description }}</p>

    <div class="mt-4 flex items-center justify-between">
      <StatusBadge :status="request.status" />
      <span class="text-xs text-subtle">{{ timeAgo(request.createdAt) }}</span>
    </div>
    <p v-if="request.assignedOfficerName" class="mt-2 text-xs text-subtle">
      Officer: <span class="text-muted">{{ request.assignedOfficerName }}</span>
    </p>
  </RouterLink>
</template>
