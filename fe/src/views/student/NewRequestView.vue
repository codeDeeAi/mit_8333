<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { requestApi } from '@/api/requests'
import { categoryApi } from '@/api/categories'
import { requestSchema, validate } from '@/lib/validation'
import type { Category, Priority } from '@/types'
import { ALL_PRIORITIES, PRIORITY_LABELS } from '@/lib/constants'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseCard from '@/components/ui/BaseCard.vue'
import BaseInput from '@/components/ui/BaseInput.vue'
import BaseTextarea from '@/components/ui/BaseTextarea.vue'
import BaseSelect from '@/components/ui/BaseSelect.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const auth = useAuthStore()
const router = useRouter()

const categories = ref<Category[]>([])
const submitting = ref(false)
const evidenceName = ref('')
const evidenceUrl = ref<string | undefined>()
const evidenceFile = ref<File | undefined>()

const form = reactive({
  title: '',
  categoryId: '',
  location: '',
  priority: 'medium' as Priority,
  description: '',
})
const errors = ref<Record<string, string>>({})

const priorityOptions = ALL_PRIORITIES.map((p) => ({ value: p, label: PRIORITY_LABELS[p] }))

onMounted(async () => {
  categories.value = await categoryApi.list()
})

function onFile(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    toast.error('File must be under 5 MB')
    return
  }
  evidenceName.value = file.name
  evidenceFile.value = file
  // Object URL drives the local preview; the real API uploads the File itself.
  evidenceUrl.value = URL.createObjectURL(file)
}

async function submit() {
  errors.value = (await validate(requestSchema, form)) ?? {}
  if (Object.keys(errors.value).length) return
  submitting.value = true
  try {
    const req = await requestApi.create({
      ...form,
      evidenceUrl: evidenceUrl.value,
      evidenceFile: evidenceFile.value,
      createdBy: auth.user!,
    })
    toast.success('Request submitted successfully')
    router.push(`/requests/${req.id}`)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : 'Could not submit request')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl">
    <PageHeader title="New service request" subtitle="Describe the fault so it can be assigned and resolved quickly." />

    <BaseCard>
      <form class="space-y-5" @submit.prevent="submit">
        <BaseInput v-model="form.title" label="Title" placeholder="e.g. Broken socket in Room 12" :error="errors.title" />

        <div class="grid gap-5 sm:grid-cols-2">
          <BaseSelect
            v-model="form.categoryId"
            label="Category"
            placeholder="Select a category"
            :options="categories.map((c) => ({ value: c.id, label: c.name }))"
            :error="errors.categoryId"
          />
          <BaseSelect v-model="form.priority" label="Priority" :options="priorityOptions" :error="errors.priority" />
        </div>

        <BaseInput v-model="form.location" label="Location" placeholder="Building, room or hostel block" :error="errors.location" />

        <BaseTextarea
          v-model="form.description"
          label="Description"
          :rows="5"
          placeholder="Describe the problem, when it started, and any safety concerns…"
          :error="errors.description"
        />

        <!-- Evidence upload -->
        <div>
          <span class="mb-1.5 block text-sm font-medium text-fg">Evidence (optional)</span>
          <label class="flex cursor-pointer items-center gap-3 rounded-lg border border-dashed border-border-strong bg-surface-2 px-4 py-3 text-sm text-muted transition-colors hover:border-accent">
            <svg class="h-5 w-5 text-subtle" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
            </svg>
            <span>{{ evidenceName || 'Upload a photo of the fault (max 5 MB)' }}</span>
            <input type="file" accept="image/*" class="hidden" @change="onFile" />
          </label>
          <div v-if="evidenceUrl" class="mt-3">
            <img :src="evidenceUrl" alt="Evidence preview" class="max-h-40 rounded-lg border border-border" />
          </div>
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <BaseButton variant="secondary" type="button" @click="router.back()">Cancel</BaseButton>
          <BaseButton type="submit" :loading="submitting">Submit request</BaseButton>
        </div>
      </form>
    </BaseCard>
  </div>
</template>
