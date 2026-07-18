<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
</script>

<template>
  <div class="bg-grid-glow flex min-h-screen flex-col items-center justify-center px-4 text-center">
    <p class="text-6xl font-semibold tracking-tight text-fg">404</p>
    <p class="mt-3 text-muted">This page could not be found.</p>
    <RouterLink :to="auth.isAuthenticated ? auth.homeRoute : '/login'" class="mt-6">
      <BaseButton>Go home</BaseButton>
    </RouterLink>
  </div>
</template>
