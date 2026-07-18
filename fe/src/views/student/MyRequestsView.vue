<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import PageHeader from '@/components/ui/PageHeader.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import RequestBrowser from '@/components/requests/RequestBrowser.vue'

const auth = useAuthStore()
</script>

<template>
  <div>
    <PageHeader title="My requests" subtitle="Search and track everything you've submitted.">
      <template #actions>
        <RouterLink to="/app/new">
          <BaseButton>+ New request</BaseButton>
        </RouterLink>
      </template>
    </PageHeader>

    <RequestBrowser :scope="{ userId: auth.user!.id }" />
  </div>
</template>
