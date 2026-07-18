<script setup lang="ts">
import type { RequestStatus } from '@/types'
import { STATUS_LABELS, STATUS_STYLES } from '@/lib/constants'

defineProps<{ status: RequestStatus }>()
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium"
    :class="STATUS_STYLES[status]"
  >
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ STATUS_LABELS[status] }}
  </span>
</template>
