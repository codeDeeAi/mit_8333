<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center rounded-xl border border-dashed border-border-strong py-16 text-center">
    <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-surface-2 text-subtle">
      <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6M5 8h14M5 8a2 2 0 012-2h10a2 2 0 012 2M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8" />
      </svg>
    </div>
    <p class="text-sm font-medium text-fg">{{ title }}</p>
    <p v-if="subtitle" class="mt-1 text-sm text-subtle">{{ subtitle }}</p>
    <slot />
  </div>
</template>
