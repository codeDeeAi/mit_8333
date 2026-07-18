<script setup lang="ts">
defineProps<{
  label?: string
  modelValue: string
  error?: string
  options: { value: string; label: string }[]
  placeholder?: string
}>()
defineEmits<{ (e: 'update:modelValue', v: string): void }>()
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 block text-sm font-medium text-fg">{{ label }}</span>
    <div class="relative">
      <select
        :value="modelValue"
        class="w-full h-10 appearance-none rounded-lg bg-surface-2 border px-3 pr-9 text-sm text-fg transition-colors focus:outline-none focus:border-accent"
        :class="error ? 'border-danger/60' : 'border-border-strong'"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
      >
        <option v-if="placeholder" value="">{{ placeholder }}</option>
        <option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
      <svg
        class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-subtle"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fill-rule="evenodd"
          d="M5.23 7.21a.75.75 0 011.06.02L10 11.17l3.71-3.94a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
    <span v-if="error" class="mt-1 block text-xs text-red-400">{{ error }}</span>
  </label>
</template>
