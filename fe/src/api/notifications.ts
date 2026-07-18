import http, { USE_MOCK } from './http'
import { notificationApi as mockNotificationApi } from './services'
import type { Notification } from '@/types'

interface Envelope<T> {
  success: boolean
  message: string
  data: T
}

interface BackendNotification {
  id: number
  user_id: number
  message: string
  is_read: boolean
  created_at: string
}

function mapNotification(n: BackendNotification): Notification {
  return {
    id: String(n.id),
    userId: String(n.user_id),
    message: n.message,
    isRead: n.is_read,
    createdAt: n.created_at,
  }
}

/**
 * Notifications API — real backend with mock fallback. The backend scopes to the
 * authenticated user, so `userId` is only used by the mock.
 */
export const notificationApi = {
  async list(userId: string): Promise<Notification[]> {
    if (USE_MOCK) return mockNotificationApi.list(userId)
    const { data } = await http.get<Envelope<BackendNotification[]>>('/notifications')
    return data.data.map(mapNotification)
  },

  async markRead(id: string): Promise<void> {
    if (USE_MOCK) return mockNotificationApi.markRead(id)
    await http.put(`/notifications/${id}/read`)
  },

  async markAllRead(userId: string): Promise<void> {
    if (USE_MOCK) return mockNotificationApi.markAllRead(userId)
    await http.put('/notifications/read-all')
  },
}
