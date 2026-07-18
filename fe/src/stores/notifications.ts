import { defineStore } from 'pinia'
import { notificationApi } from '@/api/notifications'
import type { Notification } from '@/types'

export const useNotificationStore = defineStore('notifications', {
  state: (): { items: Notification[] } => ({ items: [] }),

  getters: {
    unread: (s) => s.items.filter((n) => !n.isRead),
    unreadCount(): number {
      return this.unread.length
    },
  },

  actions: {
    async fetch(userId: string) {
      this.items = await notificationApi.list(userId)
    },
    async markRead(id: string) {
      await notificationApi.markRead(id)
      const n = this.items.find((x) => x.id === id)
      if (n) n.isRead = true
    },
    async markAllRead(userId: string) {
      await notificationApi.markAllRead(userId)
      this.items.forEach((n) => (n.isRead = true))
    },
  },
})
