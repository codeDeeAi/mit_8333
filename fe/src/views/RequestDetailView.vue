<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { requestApi } from '@/api/requests'
import { userApi } from '@/api/users'
import type { RequestStatus, ServiceRequest, User } from '@/types'
import { ALL_STATUSES, STATUS_LABELS } from '@/lib/constants'
import { formatDateTime } from '@/lib/format'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseCard from '@/components/ui/BaseCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseSelect from '@/components/ui/BaseSelect.vue'
import BaseTextarea from '@/components/ui/BaseTextarea.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import PriorityBadge from '@/components/ui/PriorityBadge.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const id = route.params.id as string
const request = ref<ServiceRequest | null>(null)
const loading = ref(true)
const officers = ref<User[]>([])

// Status update form
const newStatus = ref<RequestStatus>('in_progress')
const statusNote = ref('')
const updating = ref(false)

// Assign modal
const showAssign = ref(false)
const selectedOfficer = ref('')
const assigning = ref(false)

// Delete modal
const showDelete = ref(false)

const canUpdateStatus = computed(
  () => auth.isAdmin || (auth.isOfficer && request.value?.assignedOfficerId === auth.user?.id),
)
const canAssign = computed(() => auth.isAdmin)
const statusOptions = ALL_STATUSES.map((s) => ({ value: s, label: STATUS_LABELS[s] }))

async function load() {
  loading.value = true
  try {
    request.value = await requestApi.get(id)
    newStatus.value = request.value.status
  } catch {
    toast.error('Request not found')
    router.back()
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
  if (canAssign.value) officers.value = await userApi.officers()
})

async function submitStatus() {
  if (!request.value) return
  updating.value = true
  try {
    await requestApi.updateStatus(id, newStatus.value, auth.user!, statusNote.value || undefined)
    statusNote.value = ''
    toast.success('Status updated')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Update failed')
  } finally {
    updating.value = false
  }
}

async function assign() {
  if (!selectedOfficer.value) return
  assigning.value = true
  try {
    await requestApi.assign(id, selectedOfficer.value, auth.user!)
    toast.success('Request assigned')
    showAssign.value = false
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Assignment failed')
  } finally {
    assigning.value = false
  }
}

async function remove() {
  try {
    await requestApi.remove(id, auth.user!)
    toast.success('Request deleted')
    router.push('/admin/requests')
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Delete failed')
  }
}
</script>

<template>
  <div v-if="loading" class="space-y-4">
    <div class="h-8 w-48 animate-pulse rounded bg-surface" />
    <div class="h-64 animate-pulse rounded-xl border border-border bg-surface" />
  </div>

  <div v-else-if="request">
    <button class="mb-4 flex items-center gap-1.5 text-sm text-subtle hover:text-fg" @click="router.back()">
      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
      </svg>
      Back
    </button>

    <PageHeader :title="request.title">
      <template #actions>
        <template v-if="auth.isAdmin">
          <BaseButton v-if="canAssign" variant="secondary" @click="showAssign = true">
            {{ request.assignedOfficerName ? 'Reassign' : 'Assign' }}
          </BaseButton>
          <BaseButton variant="danger" @click="showDelete = true">Delete</BaseButton>
        </template>
      </template>
    </PageHeader>

    <div class="grid gap-6 lg:grid-cols-3">
      <!-- Main -->
      <div class="space-y-6 lg:col-span-2">
        <BaseCard>
          <div class="mb-4 flex flex-wrap items-center gap-2">
            <StatusBadge :status="request.status" />
            <PriorityBadge :priority="request.priority" />
            <span class="rounded-full border border-border-strong px-2.5 py-0.5 text-xs text-muted">
              {{ request.categoryName }}
            </span>
          </div>
          <p class="whitespace-pre-line text-sm leading-relaxed text-muted">{{ request.description }}</p>

          <div v-if="request.evidenceUrl" class="mt-5">
            <p class="mb-2 text-xs font-medium text-subtle">Evidence</p>
            <img :src="request.evidenceUrl" alt="Evidence" class="max-h-72 rounded-lg border border-border" />
          </div>
        </BaseCard>

        <!-- Status update (officer/admin) -->
        <BaseCard v-if="canUpdateStatus">
          <h3 class="mb-4 text-sm font-semibold text-fg">Update status</h3>
          <div class="space-y-4">
            <BaseSelect v-model="newStatus" label="New status" :options="statusOptions" />
            <BaseTextarea v-model="statusNote" label="Note (optional)" :rows="2" placeholder="Add a progress note…" />
            <div class="flex justify-end">
              <BaseButton :loading="updating" @click="submitStatus">Save update</BaseButton>
            </div>
          </div>
        </BaseCard>

        <!-- Timeline -->
        <BaseCard>
          <h3 class="mb-4 text-sm font-semibold text-fg">Activity timeline</h3>
          <ol class="relative space-y-5 border-l border-border pl-5">
            <li v-for="log in request.logs" :key="log.id" class="relative">
              <span class="absolute -left-[27px] top-1 h-3 w-3 rounded-full border-2 border-bg bg-accent" />
              <div class="flex flex-wrap items-center gap-2 text-sm">
                <span class="text-fg">{{ log.changedByName }}</span>
                <span class="text-subtle">
                  {{ log.oldStatus ? `moved to` : 'created request' }}
                </span>
                <StatusBadge v-if="log.oldStatus" :status="log.newStatus" />
              </div>
              <p v-if="log.note" class="mt-1 text-sm text-muted">{{ log.note }}</p>
              <p class="mt-1 text-xs text-subtle">{{ formatDateTime(log.createdAt) }}</p>
            </li>
          </ol>
        </BaseCard>
      </div>

      <!-- Sidebar meta -->
      <div class="space-y-4">
        <BaseCard>
          <dl class="space-y-4 text-sm">
            <div>
              <dt class="text-xs text-subtle">Requested by</dt>
              <dd class="mt-0.5 text-fg">{{ request.createdByName }}</dd>
            </div>
            <div>
              <dt class="text-xs text-subtle">Location</dt>
              <dd class="mt-0.5 text-fg">{{ request.location }}</dd>
            </div>
            <div>
              <dt class="text-xs text-subtle">Assigned officer</dt>
              <dd class="mt-0.5 text-fg">{{ request.assignedOfficerName ?? 'Unassigned' }}</dd>
            </div>
            <div>
              <dt class="text-xs text-subtle">Created</dt>
              <dd class="mt-0.5 text-fg">{{ formatDateTime(request.createdAt) }}</dd>
            </div>
            <div>
              <dt class="text-xs text-subtle">Last updated</dt>
              <dd class="mt-0.5 text-fg">{{ formatDateTime(request.updatedAt) }}</dd>
            </div>
          </dl>
        </BaseCard>
      </div>
    </div>

    <!-- Assign modal -->
    <BaseModal v-if="canAssign" :open="showAssign" title="Assign officer" @close="showAssign = false">
      <div class="space-y-4">
        <BaseSelect
          v-model="selectedOfficer"
          label="Maintenance officer"
          placeholder="Select an officer"
          :options="officers.map((o) => ({ value: o.id, label: o.fullName }))"
        />
        <div class="flex justify-end gap-3">
          <BaseButton variant="secondary" @click="showAssign = false">Cancel</BaseButton>
          <BaseButton :loading="assigning" :disabled="!selectedOfficer" @click="assign">Assign</BaseButton>
        </div>
      </div>
    </BaseModal>

    <!-- Delete modal -->
    <BaseModal :open="showDelete" title="Delete request" @close="showDelete = false">
      <p class="text-sm text-muted">This action cannot be undone. Are you sure?</p>
      <div class="mt-5 flex justify-end gap-3">
        <BaseButton variant="secondary" @click="showDelete = false">Cancel</BaseButton>
        <BaseButton variant="danger" @click="remove">Delete</BaseButton>
      </div>
    </BaseModal>
  </div>
</template>
