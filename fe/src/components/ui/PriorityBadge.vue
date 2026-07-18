<script setup lang="ts">
import type { Priority } from '@/types'
import { PRIORITY_LABELS, PRIORITY_STYLES } from '@/lib/constants'

defineProps<{ priority: Priority }>()
</script>

<template>
  <span
    class="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium"
    :class="PRIORITY_STYLES[priority]"
  >
    {{ PRIORITY_LABELS[priority] }}
  </span>
</template>
