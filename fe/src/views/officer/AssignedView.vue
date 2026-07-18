<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import PageHeader from '@/components/ui/PageHeader.vue'
import RequestBrowser from '@/components/requests/RequestBrowser.vue'

const auth = useAuthStore()
</script>

<template>
  <div>
    <PageHeader title="Assigned to me" subtitle="Jobs assigned to you by the administrator." />
    <RequestBrowser :scope="{ officerId: auth.user!.id }" />
  </div>
</template>
