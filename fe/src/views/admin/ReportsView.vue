<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { toast } from 'vue3-toastify'
import { reportApi, type SummaryStats } from '@/api/reports'
import { exportRequestsCsv, exportRequestsPdf } from '@/lib/export'
import { STATUS_LABELS } from '@/lib/constants'
import PageHeader from '@/components/ui/PageHeader.vue'
import StatCard from '@/components/ui/StatCard.vue'
import BaseCard from '@/components/ui/BaseCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const summary = ref<SummaryStats | null>(null)

onMounted(async () => {
  summary.value = await reportApi.summary()
})

async function download(kind: 'csv' | 'pdf') {
  const requests = await reportApi.allRequests()
  if (!requests.length) return toast.info('No requests to export')
  kind === 'csv' ? exportRequestsCsv(requests) : exportRequestsPdf(requests)
  toast.success(`Exported ${requests.length} requests`)
}
</script>

<template>
  <div>
    <PageHeader title="Reports" subtitle="Summarise activity and export records.">
      <template #actions>
        <BaseButton variant="secondary" size="sm" @click="download('csv')">Download CSV</BaseButton>
        <BaseButton size="sm" @click="download('pdf')">Download PDF</BaseButton>
      </template>
    </PageHeader>

    <div v-if="summary" class="space-y-8">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
        <StatCard label="Total" :value="summary.total" />
        <StatCard
          v-for="(count, status) in summary.byStatus"
          :key="status"
          :label="STATUS_LABELS[status]"
          :value="count"
        />
      </div>

      <BaseCard>
        <h3 class="mb-4 text-sm font-semibold text-fg">Requests by category</h3>
        <div class="grid gap-3 sm:grid-cols-2">
          <div
            v-for="cat in summary.byCategory"
            :key="cat.name"
            class="flex items-center justify-between rounded-lg border border-border bg-surface-2 px-4 py-3"
          >
            <span class="text-sm text-muted">{{ cat.name }}</span>
            <span class="text-sm font-semibold text-fg">{{ cat.count }}</span>
          </div>
        </div>
      </BaseCard>
    </div>
  </div>
</template>
