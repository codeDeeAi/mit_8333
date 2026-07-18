<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { auditApi } from '@/api/audit'
import type { AuditLog } from '@/types'
import { formatDateTime } from '@/lib/format'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseCard from '@/components/ui/BaseCard.vue'

const loading = ref(true)
const logs = ref<AuditLog[]>([])

const actionStyles: Record<string, string> = {
  LOGIN: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
  REGISTER: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  CREATE_REQUEST: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  ASSIGN: 'bg-purple-500/10 text-purple-400 border-purple-500/20',
  STATUS_CHANGE: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  UPDATE_ROLE: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  DELETE_REQUEST: 'bg-red-500/10 text-red-400 border-red-500/20',
  DELETE_USER: 'bg-red-500/10 text-red-400 border-red-500/20',
}

onMounted(async () => {
  logs.value = await auditApi.list()
  loading.value = false
})
</script>

<template>
  <div>
    <PageHeader title="Audit log" subtitle="A tamper-evident trail of every important action." />

    <BaseCard :padded="false">
      <div v-if="loading" class="p-6 text-sm text-subtle">Loading activity…</div>
      <ul v-else class="divide-y divide-border/60">
        <li v-for="log in logs" :key="log.id" class="flex flex-wrap items-center gap-3 px-5 py-3.5 text-sm">
          <span
            class="rounded-full border px-2.5 py-0.5 text-xs font-medium"
            :class="actionStyles[log.action] ?? 'bg-surface-2 text-muted border-border-strong'"
          >
            {{ log.action.replace(/_/g, ' ') }}
          </span>
          <span class="text-fg">{{ log.userName }}</span>
          <span class="text-subtle">·</span>
          <span class="text-muted">{{ log.entity }} <span class="text-subtle">#{{ log.entityId }}</span></span>
          <span class="ml-auto text-xs text-subtle">{{ formatDateTime(log.createdAt) }}</span>
        </li>
      </ul>
    </BaseCard>
  </div>
</template>
