<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'
import { timeAgo } from '@/lib/format'

const auth = useAuthStore()
const notifications = useNotificationStore()
const open = ref(false)

onMounted(() => {
  if (auth.user) notifications.fetch(auth.user.id)
})

async function markAll() {
  if (auth.user) await notifications.markAllRead(auth.user.id)
}
</script>

<template>
  <div class="relative">
    <button
      class="relative flex h-9 w-9 items-center justify-center rounded-lg border border-border-strong text-muted transition-colors hover:bg-surface-2 hover:text-fg"
      aria-label="Notifications"
      @click="open = !open"
    >
      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.8 23.8 0 005.454-1.31A8.97 8.97 0 0118 9.75V9A6 6 0 006 9v.75a8.97 8.97 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.3 24.3 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
      </svg>
      <span
        v-if="notifications.unreadCount"
        class="absolute -right-1 -top-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-accent px-1 text-[10px] font-semibold text-white"
      >
        {{ notifications.unreadCount }}
      </span>
    </button>

    <Transition name="fade">
      <div
        v-if="open"
        class="absolute right-0 z-40 mt-2 w-80 rounded-xl border border-border-strong bg-surface shadow-2xl"
      >
        <div class="flex items-center justify-between border-b border-border px-4 py-3">
          <span class="text-sm font-semibold text-fg">Notifications</span>
          <button class="text-xs text-accent hover:underline" @click="markAll">Mark all read</button>
        </div>
        <div class="max-h-80 overflow-y-auto">
          <p v-if="!notifications.items.length" class="px-4 py-6 text-center text-sm text-subtle">
            You're all caught up.
          </p>
          <button
            v-for="n in notifications.items"
            :key="n.id"
            class="flex w-full gap-3 border-b border-border/60 px-4 py-3 text-left transition-colors hover:bg-surface-2"
            @click="notifications.markRead(n.id)"
          >
            <span
              class="mt-1.5 h-2 w-2 shrink-0 rounded-full"
              :class="n.isRead ? 'bg-transparent' : 'bg-accent'"
            />
            <span class="min-w-0">
              <span class="block text-sm text-fg">{{ n.message }}</span>
              <span class="mt-0.5 block text-xs text-subtle">{{ timeAgo(n.createdAt) }}</span>
            </span>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
