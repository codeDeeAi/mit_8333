<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { USE_MOCK } from '@/api/http'
import { loginSchema, validate } from '@/lib/validation'
import BaseInput from '@/components/ui/BaseInput.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const form = reactive({ email: '', password: '' })
const errors = ref<Record<string, string>>({})

async function submit() {
  errors.value = (await validate(loginSchema, form)) ?? {}
  if (Object.keys(errors.value).length) return
  try {
    await auth.login(form.email, form.password)
    toast.success(`Welcome back, ${auth.user?.fullName.split(' ')[0]}`)
    const redirect = (route.query.redirect as string) || auth.homeRoute
    router.push(redirect)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Login failed')
  }
}

function fill(email: string) {
  form.email = email
  form.password = 'password123'
}
</script>

<template>
  <div class="bg-grid-glow flex min-h-screen items-center justify-center px-4">
    <div class="w-full max-w-sm animate-fade-in-up">
      <div class="mb-8 text-center">
        <div class="mx-auto mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-white text-lg font-bold text-black">
          M
        </div>
        <h1 class="text-2xl font-semibold tracking-tight">Sign in to MIVA FixIt</h1>
        <p class="mt-1 text-sm text-subtle">University maintenance & service requests</p>
      </div>

      <form class="space-y-4 rounded-2xl border border-border bg-surface/80 p-6 backdrop-blur" @submit.prevent="submit">
        <BaseInput
          v-model="form.email"
          label="Email"
          type="email"
          placeholder="you@miva.edu"
          autocomplete="email"
          :error="errors.email"
        />
        <BaseInput
          v-model="form.password"
          label="Password"
          type="password"
          placeholder="••••••••"
          autocomplete="current-password"
          :error="errors.password"
        />
        <BaseButton type="submit" block :loading="auth.loading">Sign in</BaseButton>
      </form>

      <p class="mt-4 text-center text-sm text-subtle">
        No account?
        <RouterLink to="/register" class="text-fg hover:underline">Create one</RouterLink>
      </p>

      <div v-if="USE_MOCK" class="mt-6 rounded-xl border border-border bg-surface/50 p-3 text-xs text-subtle">
        <p class="mb-2 font-medium text-muted">Demo accounts (password: password123)</p>
        <div class="grid grid-cols-3 gap-2">
          <button class="rounded-lg border border-border-strong py-1.5 hover:bg-surface-2" @click="fill('student@miva.edu')">Student</button>
          <button class="rounded-lg border border-border-strong py-1.5 hover:bg-surface-2" @click="fill('officer@miva.edu')">Officer</button>
          <button class="rounded-lg border border-border-strong py-1.5 hover:bg-surface-2" @click="fill('admin@miva.edu')">Admin</button>
        </div>
      </div>
    </div>
  </div>
</template>
