<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { reportApi } from '@/api/reports'
import { requestApi } from '@/api/requests'
import { userApi } from '@/api/users'
import { exportRequestsCsv, exportRequestsPdf } from '@/lib/export'
import { ALL_STATUSES, STATUS_LABELS } from '@/lib/constants'
import type { RequestStatus, ServiceRequest, User } from '@/types'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import BaseSelect from '@/components/ui/BaseSelect.vue'
import BaseTextarea from '@/components/ui/BaseTextarea.vue'
import RequestBrowser from '@/components/requests/RequestBrowser.vue'

const auth = useAuthStore()

const browserRef = ref<{ reload: () => Promise<void> } | null>(null)
const officers = ref<User[]>([])
const loadingOfficers = ref(true)

const showAssign = ref(false)
const selectedRequest = ref<ServiceRequest | null>(null)
const selectedOfficer = ref('')
const assigning = ref(false)

const showStatus = ref(false)
const selectedStatus = ref<RequestStatus>('in_progress')
const statusNote = ref('')
const updatingStatus = ref(false)

const statusOptions = ALL_STATUSES.map((status) => ({ value: status, label: STATUS_LABELS[status] }))

onMounted(async () => {
  try {
    officers.value = await userApi.officers()
  } finally {
    loadingOfficers.value = false
  }
})

function openAssign(request: ServiceRequest) {
  selectedRequest.value = request
  selectedOfficer.value = request.assignedOfficerId ?? ''
  showAssign.value = true
}

function openStatus(request: ServiceRequest) {
  selectedRequest.value = request
  selectedStatus.value = request.status
  statusNote.value = ''
  showStatus.value = true
}

async function assign() {
  if (!selectedRequest.value || !selectedOfficer.value || !auth.user) return
  assigning.value = true
  try {
    await requestApi.assign(selectedRequest.value.id, selectedOfficer.value, auth.user)
    toast.success('Request assigned')
    showAssign.value = false
    await browserRef.value?.reload()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Assignment failed')
  } finally {
    assigning.value = false
  }
}

async function updateStatus() {
  if (!selectedRequest.value || !auth.user) return
  updatingStatus.value = true
  try {
    await requestApi.updateStatus(
      selectedRequest.value.id,
      selectedStatus.value,
      auth.user,
      statusNote.value || undefined,
    )
    toast.success('Status updated')
    showStatus.value = false
    await browserRef.value?.reload()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Status update failed')
  } finally {
    updatingStatus.value = false
  }
}

async function download(kind: 'csv' | 'pdf') {
  const requests = await reportApi.allRequests()
  if (!requests.length) {
    toast.info('No requests to export')
    return
  }
  if (kind === 'csv') exportRequestsCsv(requests)
  else exportRequestsPdf(requests)
  toast.success(`Exported ${requests.length} requests`)
}
</script>

<template>
  <div>
    <PageHeader title="Manage requests" subtitle="Search, filter, assign and export all service requests.">
      <template #actions>
        <BaseButton variant="secondary" size="sm" @click="download('csv')">Export CSV</BaseButton>
        <BaseButton variant="secondary" size="sm" @click="download('pdf')">Export PDF</BaseButton>
      </template>
    </PageHeader>

    <RequestBrowser ref="browserRef">
      <template #card-actions="{ request }">
        <div class="flex justify-end gap-2">
          <BaseButton variant="secondary" size="sm" @click="openStatus(request)">
            Update status
          </BaseButton>
          <BaseButton variant="secondary" size="sm" @click="openAssign(request)">
            {{ request.assignedOfficerName ? 'Reassign' : 'Assign' }}
          </BaseButton>
        </div>
      </template>
    </RequestBrowser>

    <BaseModal :open="showAssign" title="Assign officer" @close="showAssign = false">
      <div class="space-y-4">
        <p v-if="selectedRequest" class="text-sm text-muted">
          Assigning: <span class="text-fg">{{ selectedRequest.title }}</span>
        </p>
        <BaseSelect
          v-model="selectedOfficer"
          label="Maintenance officer"
          :placeholder="loadingOfficers ? 'Loading officers…' : 'Select an officer'"
          :options="officers.map((o) => ({ value: o.id, label: o.fullName }))"
        />
        <div class="flex justify-end gap-3">
          <BaseButton variant="secondary" @click="showAssign = false">Cancel</BaseButton>
          <BaseButton :loading="assigning" :disabled="!selectedOfficer" @click="assign">Assign</BaseButton>
        </div>
      </div>
    </BaseModal>

    <BaseModal :open="showStatus" title="Update request status" @close="showStatus = false">
      <div class="space-y-4">
        <p v-if="selectedRequest" class="text-sm text-muted">
          Updating: <span class="text-fg">{{ selectedRequest.title }}</span>
        </p>
        <BaseSelect v-model="selectedStatus" label="Status" :options="statusOptions" />
        <BaseTextarea
          v-model="statusNote"
          label="Note (optional)"
          :rows="2"
          placeholder="Add a progress note..."
        />
        <div class="flex justify-end gap-3">
          <BaseButton variant="secondary" @click="showStatus = false">Cancel</BaseButton>
          <BaseButton :loading="updatingStatus" @click="updateStatus">Save update</BaseButton>
        </div>
      </div>
    </BaseModal>
  </div>
</template>
