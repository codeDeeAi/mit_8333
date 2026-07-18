<script setup lang="ts">
withDefaults(defineProps<{ padded?: boolean; hover?: boolean }>(), {
  padded: true,
  hover: false,
})
</script>

<template>
  <div
    class="rounded-xl border border-border bg-surface"
    :class="[padded && 'p-5', hover && 'transition-colors hover:border-border-strong']"
  >
    <slot />
  </div>
</template>
