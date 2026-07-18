<script setup lang="ts">
defineProps<{
  label?: string
  modelValue: string
  placeholder?: string
  error?: string
  rows?: number
}>()
defineEmits<{ (e: 'update:modelValue', v: string): void }>()
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 block text-sm font-medium text-fg">{{ label }}</span>
    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      :rows="rows ?? 4"
      class="w-full rounded-lg bg-surface-2 border px-3 py-2 text-sm text-fg placeholder:text-subtle transition-colors focus:outline-none focus:border-accent resize-y"
      :class="error ? 'border-danger/60' : 'border-border-strong'"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />
    <span v-if="error" class="mt-1 block text-xs text-red-400">{{ error }}</span>
  </label>
</template>
