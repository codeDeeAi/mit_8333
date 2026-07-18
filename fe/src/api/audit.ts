import http, { USE_MOCK } from './http'
import { auditApi as mockAuditApi } from './services'
import type { AuditLog } from '@/types'

interface Envelope<T> {
  success: boolean
  message: string
  data: T
}

interface BackendAuditLog {
  id: number
  user_id?: number | null
  user_name?: string | null
  action: string
  entity: string
  entity_id: string
  created_at: string
}

function mapAuditLog(log: BackendAuditLog): AuditLog {
  return {
    id: String(log.id),
    userId: log.user_id ? String(log.user_id) : '0',
    userName: log.user_name ?? 'System',
    action: log.action,
    entity: log.entity,
    entityId: log.entity_id,
    createdAt: log.created_at,
  }
}

export const auditApi = {
  async list(): Promise<AuditLog[]> {
    if (USE_MOCK) return mockAuditApi.list()

    try {
      const { data } = await http.get<Envelope<BackendAuditLog[]>>('/audit-logs')
      return data.data.map(mapAuditLog)
    } catch (error: any) {
      if (error?.response?.status === 404) {
        return []
      }
      throw error
    }
  },
}
