<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/api/users'
import type { Role, User } from '@/types'
import { ROLE_LABELS } from '@/lib/constants'
import { initials } from '@/lib/format'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseCard from '@/components/ui/BaseCard.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
const loading = ref(true)
const users = ref<User[]>([])

const roleOptions: Role[] = ['student_staff', 'maintenance_officer', 'admin']
const toDelete = ref<User | null>(null)

async function load() {
  users.value = await userApi.list()
  loading.value = false
}
onMounted(load)

async function changeRole(user: User, role: Role) {
  if (user.role === role) return
  try {
    await userApi.updateRole(user.id, role, auth.user!)
    user.role = role
    toast.success(`${user.fullName} is now ${ROLE_LABELS[role]}`)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Update failed')
  }
}

async function remove() {
  if (!toDelete.value) return
  try {
    await userApi.remove(toDelete.value.id, auth.user!)
    users.value = users.value.filter((u) => u.id !== toDelete.value!.id)
    toast.success('User removed')
    toDelete.value = null
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Delete failed')
  }
}
</script>

<template>
  <div>
    <PageHeader title="Users" subtitle="Manage accounts and assign roles." />

    <BaseCard :padded="false">
      <div v-if="loading" class="p-6 text-sm text-subtle">Loading users…</div>
      <table v-else class="w-full text-sm">
        <thead>
          <tr class="border-b border-border text-left text-xs text-subtle">
            <th class="px-5 py-3 font-medium">User</th>
            <th class="px-5 py-3 font-medium">Role</th>
            <th class="hidden px-5 py-3 font-medium sm:table-cell">Phone</th>
            <th class="px-5 py-3" />
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="border-b border-border/60 last:border-0">
            <td class="px-5 py-3">
              <div class="flex items-center gap-3">
                <div class="flex h-8 w-8 items-center justify-center rounded-full bg-surface-2 text-xs font-semibold">
                  {{ initials(u.fullName) }}
                </div>
                <div>
                  <p class="font-medium text-fg">{{ u.fullName }}</p>
                  <p class="text-xs text-subtle">{{ u.email }}</p>
                </div>
              </div>
            </td>
            <td class="px-5 py-3">
              <select
                :value="u.role"
                class="h-8 rounded-lg border border-border-strong bg-surface-2 px-2 text-xs text-fg focus:border-accent focus:outline-none"
                :disabled="u.id === auth.user?.id"
                @change="changeRole(u, ($event.target as HTMLSelectElement).value as Role)"
              >
                <option v-for="r in roleOptions" :key="r" :value="r">{{ ROLE_LABELS[r] }}</option>
              </select>
            </td>
            <td class="hidden px-5 py-3 text-muted sm:table-cell">{{ u.phone ?? '—' }}</td>
            <td class="px-5 py-3 text-right">
              <button
                v-if="u.id !== auth.user?.id"
                class="text-xs text-red-400 hover:underline"
                @click="toDelete = u"
              >
                Remove
              </button>
              <span v-else class="text-xs text-subtle">You</span>
            </td>
          </tr>
        </tbody>
      </table>
    </BaseCard>

    <BaseModal :open="!!toDelete" title="Remove user" @close="toDelete = null">
      <p class="text-sm text-muted">
        Remove <span class="text-fg">{{ toDelete?.fullName }}</span
        >? This cannot be undone.
      </p>
      <div class="mt-5 flex justify-end gap-3">
        <BaseButton variant="secondary" @click="toDelete = null">Cancel</BaseButton>
        <BaseButton variant="danger" @click="remove">Remove</BaseButton>
      </div>
    </BaseModal>
  </div>
</template>
