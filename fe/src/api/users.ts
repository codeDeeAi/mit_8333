import http, { USE_MOCK } from './http'
import { userApi as mockUserApi } from './services'
import type { Role, User } from '@/types'

interface Envelope<T> {
  success: boolean
  message: string
  data: T
}

interface BackendUser {
  id: number
  full_name: string
  email: string
  role_id: number
  role: string
  phone?: string | null
  created_at: string
}

function mapUser(u: BackendUser): User {
  return {
    id: String(u.id),
    fullName: u.full_name,
    email: u.email,
    role: (u.role || 'student_staff') as Role,
    phone: u.phone ?? undefined,
    createdAt: u.created_at,
  }
}

/** Users API (admin) — real backend with mock fallback. Same shape as the mock. */
export const userApi = {
  async list(): Promise<User[]> {
    if (USE_MOCK) return mockUserApi.list()
    const { data } = await http.get<Envelope<BackendUser[]>>('/users')
    return data.data.map(mapUser)
  },

  async officers(): Promise<User[]> {
    if (USE_MOCK) return mockUserApi.officers()
    const { data } = await http.get<Envelope<BackendUser[]>>('/users/officers')
    return data.data.map(mapUser)
  },

  async updateRole(id: string, role: Role, actor: User): Promise<User> {
    if (USE_MOCK) return mockUserApi.updateRole(id, role, actor)
    const { data } = await http.put<Envelope<BackendUser>>(`/users/${id}/role`, { role })
    return mapUser(data.data)
  },

  async remove(id: string, actor: User): Promise<void> {
    if (USE_MOCK) return mockUserApi.remove(id, actor)
    await http.delete(`/users/${id}`)
  },
}
