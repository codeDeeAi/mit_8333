<script setup lang="ts">
defineProps<{
  label?: string
  modelValue: string
  type?: string
  placeholder?: string
  error?: string
  autocomplete?: string
}>()
defineEmits<{ (e: 'update:modelValue', v: string): void }>()
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 block text-sm font-medium text-fg">{{ label }}</span>
    <input
      :type="type ?? 'text'"
      :value="modelValue"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      class="w-full h-10 rounded-lg bg-surface-2 border px-3 text-sm text-fg placeholder:text-subtle transition-colors focus:outline-none focus:border-accent"
      :class="error ? 'border-danger/60' : 'border-border-strong'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="mt-1 block text-xs text-red-400">{{ error }}</span>
  </label>
</template>
