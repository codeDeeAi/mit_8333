<script setup lang="ts">
defineProps<{ label: string; value: number | string; accent?: string }>()
</script>

<template>
  <div class="rounded-xl border border-border bg-surface p-5">
    <p class="text-sm text-muted">{{ label }}</p>
    <p class="mt-2 text-3xl font-semibold tracking-tight" :class="accent ?? 'text-fg'">
      {{ value }}
    </p>
  </div>
</template>
