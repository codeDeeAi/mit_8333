<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { useAuthStore } from '@/stores/auth'
import { ROLE_LABELS } from '@/lib/constants'
import { initials } from '@/lib/format'
import NotificationBell from './NotificationBell.vue'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(false)

interface NavItem {
  label: string
  to: string
  icon: string
}

const navItems = computed<NavItem[]>(() => {
  if (auth.isAdmin) {
    return [
      { label: 'Overview', to: '/admin', icon: 'grid' },
      { label: 'Requests', to: '/admin/requests', icon: 'clipboard' },
      { label: 'Users', to: '/admin/users', icon: 'users' },
      { label: 'Reports', to: '/admin/reports', icon: 'chart' },
      { label: 'Audit Log', to: '/admin/audit', icon: 'shield' },
    ]
  }
  if (auth.isOfficer) {
    return [
      { label: 'Overview', to: '/officer', icon: 'grid' },
      { label: 'Assigned to Me', to: '/officer/assigned', icon: 'clipboard' },
    ]
  }
  return [
    { label: 'Overview', to: '/app', icon: 'grid' },
    { label: 'New Request', to: '/app/new', icon: 'plus' },
    { label: 'My Requests', to: '/app/requests', icon: 'clipboard' },
  ]
})

const icons: Record<string, string> = {
  grid: 'M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zm10 0a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-4zm10 0a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z',
  clipboard: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
  users: 'M15 19.128a9.4 9.4 0 002.625.372 9.3 9.3 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.3 12.3 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z',
  chart: 'M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z',
  shield: 'M9 12.75L11.25 15 15 9.75M21 12c0 1.268-.63 2.39-1.593 3.068a3.745 3.745 0 01-1.043 3.296 3.745 3.745 0 01-3.296 1.043A3.745 3.745 0 0112 21c-1.268 0-2.39-.63-3.068-1.593a3.746 3.746 0 01-3.296-1.043 3.745 3.745 0 01-1.043-3.296A3.745 3.745 0 013 12c0-1.268.63-2.39 1.593-3.068a3.745 3.745 0 011.043-3.296 3.746 3.746 0 013.296-1.043A3.746 3.746 0 0112 3c1.268 0 2.39.63 3.068 1.593a3.746 3.746 0 013.296 1.043 3.746 3.746 0 011.043 3.296A3.745 3.745 0 0121 12z',
  plus: 'M12 4.5v15m7.5-7.5h-15',
}

function isActive(to: string): boolean {
  return route.path === to
}

async function logout() {
  await auth.logout()
  toast.info('Signed out')
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-bg text-fg">
    <!-- Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-40 w-64 transform border-r border-border bg-surface transition-transform lg:translate-x-0"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="flex h-16 items-center gap-2 border-b border-border px-6">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-white text-black font-bold">
          M
        </div>
        <div class="leading-tight">
          <p class="text-sm font-semibold">MIVA FixIt</p>
          <p class="text-[11px] text-subtle">Maintenance Desk</p>
        </div>
      </div>

      <nav class="space-y-1 p-3">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors"
          :class="
            isActive(item.to)
              ? 'bg-surface-2 text-fg'
              : 'text-muted hover:bg-surface-2 hover:text-fg'
          "
          @click="sidebarOpen = false"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" :d="icons[item.icon]" />
          </svg>
          {{ item.label }}
        </RouterLink>
      </nav>
    </aside>

    <!-- Sidebar backdrop (mobile) -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-30 bg-black/60 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Main -->
    <div class="lg:pl-64">
      <header
        class="sticky top-0 z-20 flex h-16 items-center justify-between gap-4 border-b border-border bg-bg/80 px-4 backdrop-blur sm:px-6"
      >
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg border border-border-strong text-muted lg:hidden"
          @click="sidebarOpen = true"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5M3.75 17.25h16.5" />
          </svg>
        </button>

        <div class="hidden text-sm text-subtle sm:block">
          {{ auth.user ? ROLE_LABELS[auth.user.role] : '' }} workspace
        </div>

        <div class="ml-auto flex items-center gap-3">
          <NotificationBell />
          <div class="flex items-center gap-3 border-l border-border pl-3">
            <div class="hidden text-right sm:block">
              <p class="text-sm font-medium leading-none">{{ auth.user?.fullName }}</p>
              <p class="mt-0.5 text-xs text-subtle">{{ auth.user?.email }}</p>
            </div>
            <div class="flex h-9 w-9 items-center justify-center rounded-full bg-surface-2 text-xs font-semibold text-fg">
              {{ auth.user ? initials(auth.user.fullName) : '?' }}
            </div>
            <button
              class="flex h-9 w-9 items-center justify-center rounded-lg border border-border-strong text-muted transition-colors hover:bg-surface-2 hover:text-fg"
              title="Sign out"
              @click="logout"
            >
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
              </svg>
            </button>
          </div>
        </div>
      </header>

      <main class="mx-auto max-w-6xl px-4 py-8 sm:px-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
