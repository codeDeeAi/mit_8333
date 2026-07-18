<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ page: number; pageSize: number; total: number }>()
const emit = defineEmits<{ (e: 'update:page', v: number): void }>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1))
const to = computed(() => Math.min(props.page * props.pageSize, props.total))

function go(p: number) {
  if (p >= 1 && p <= totalPages.value) emit('update:page', p)
}
</script>

<template>
  <div class="flex items-center justify-between gap-4 text-sm">
    <p class="text-subtle">
      Showing <span class="text-fg">{{ from }}–{{ to }}</span> of
      <span class="text-fg">{{ total }}</span>
    </p>
    <div class="flex items-center gap-2">
      <button
        class="h-8 rounded-lg border border-border-strong px-3 text-xs text-muted transition-colors hover:bg-surface-2 disabled:opacity-40 disabled:pointer-events-none"
        :disabled="page <= 1"
        @click="go(page - 1)"
      >
        Previous
      </button>
      <span class="text-xs text-subtle">Page {{ page }} / {{ totalPages }}</span>
      <button
        class="h-8 rounded-lg border border-border-strong px-3 text-xs text-muted transition-colors hover:bg-surface-2 disabled:opacity-40 disabled:pointer-events-none"
        :disabled="page >= totalPages"
        @click="go(page + 1)"
      >
        Next
      </button>
    </div>
  </div>
</template>
