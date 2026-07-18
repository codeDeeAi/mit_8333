<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    size?: 'sm' | 'md'
    type?: 'button' | 'submit'
    loading?: boolean
    disabled?: boolean
    block?: boolean
  }>(),
  { variant: 'primary', size: 'md', type: 'button', loading: false, disabled: false, block: false },
)

const classes = computed(() => {
  const base =
    'inline-flex items-center justify-center gap-2 font-medium rounded-lg border transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-accent/50 disabled:opacity-50 disabled:pointer-events-none'
  const sizes = { sm: 'text-xs px-3 h-8', md: 'text-sm px-4 h-10' }
  const variants = {
    primary: 'bg-white text-black border-white hover:bg-zinc-200',
    secondary: 'bg-surface-2 text-fg border-border-strong hover:bg-elevated',
    ghost: 'bg-transparent text-muted border-transparent hover:bg-surface-2 hover:text-fg',
    danger: 'bg-danger/10 text-red-400 border-danger/30 hover:bg-danger/20',
  }
  return [base, sizes[props.size], variants[props.variant], props.block && 'w-full']
})
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <svg
      v-if="loading"
      class="animate-spin h-4 w-4"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
    </svg>
    <slot />
  </button>
</template>
