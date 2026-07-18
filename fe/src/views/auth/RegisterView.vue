<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { getRegistrationData, type RoleOption } from '@/api/auth'
import { ROLE_LABELS } from '@/lib/constants'
import type { Role } from '@/types'
import { registerSchema, validate } from '@/lib/validation'
import BaseInput from '@/components/ui/BaseInput.vue'
import BaseSelect from '@/components/ui/BaseSelect.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
const router = useRouter()

const roles = ref<RoleOption[]>([])
const rolesLoading = ref(true)
const SELF_REGISTRATION_ROLES: Role[] = ['student_staff']

const roleOptions = computed(() =>
  roles.value
    .filter((role) => SELF_REGISTRATION_ROLES.includes(role.name as Role))
    .map((role) => ({ value: role.id, label: ROLE_LABELS[role.name as Role] ?? role.name })),
)

onMounted(async () => {
  try {
    roles.value = await getRegistrationData()
    const defaultRole = roles.value.find((role) => role.name === 'student_staff')
    if (defaultRole) form.roleId = defaultRole.id
  } catch {
    // Non-fatal: the user can still submit; the backend defaults the role.
  } finally {
    rolesLoading.value = false
  }
})

const form = reactive({
  fullName: '',
  email: '',
  phone: '',
  roleId: '',
  password: '',
  confirmPassword: '',
})
const errors = ref<Record<string, string>>({})

async function submit() {
  errors.value = (await validate(registerSchema, form)) ?? {}
  if (Object.keys(errors.value).length) return
  try {
    await auth.register({
      fullName: form.fullName,
      email: form.email,
      phone: form.phone || undefined,
      roleId: form.roleId || undefined,
      password: form.password,
    })
    toast.success('Account created — welcome!')
    router.push(auth.homeRoute)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Registration failed')
  }
}
</script>

<template>
  <div class="bg-grid-glow flex min-h-screen items-center justify-center px-4 py-10">
    <div class="w-full max-w-sm animate-fade-in-up">
      <div class="mb-8 text-center">
        <div class="mx-auto mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-white text-lg font-bold text-black">
          M
        </div>
        <h1 class="text-2xl font-semibold tracking-tight">Create your account</h1>
        <p class="mt-1 text-sm text-subtle">Submit and track maintenance requests</p>
      </div>

      <form class="space-y-4 rounded-2xl border border-border bg-surface/80 p-6 backdrop-blur" @submit.prevent="submit">
        <BaseInput v-model="form.fullName" label="Full name" placeholder="Jane Doe" :error="errors.fullName" />
        <BaseInput v-model="form.email" label="Email" type="email" placeholder="you@miva.edu" :error="errors.email" />
        <BaseInput v-model="form.phone" label="Phone (optional)" placeholder="+234…" :error="errors.phone" />
        <BaseSelect
          v-model="form.roleId"
          label="Account type"
          :options="roleOptions"
          :placeholder="rolesLoading ? 'Loading account types…' : 'Select an account type'"
          :error="errors.roleId"
        />
        <BaseInput v-model="form.password" label="Password" type="password" placeholder="••••••••" :error="errors.password" />
        <BaseInput v-model="form.confirmPassword" label="Confirm password" type="password" placeholder="••••••••" :error="errors.confirmPassword" />
        <BaseButton type="submit" block :loading="auth.loading">Create account</BaseButton>
      </form>

      <p class="mt-4 text-center text-sm text-subtle">
        Already have an account?
        <RouterLink to="/login" class="text-fg hover:underline">Sign in</RouterLink>
      </p>
    </div>
  </div>
</template>
