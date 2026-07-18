import http, { USE_MOCK } from './http'
import { authApi as mockAuth } from './services'
import type { Role, User } from '@/types'

/**
 * Real auth API backed by the Go service. Falls back to the in-memory mock when
 * VITE_USE_MOCK is enabled, so the UI works with or without a running backend.
 */

export interface AuthResult {
  user: User
  token: string
}

export interface RoleOption {
  id: string
  name: Role | string
  description?: string
}

interface BackendUser {
  id: number
  full_name: string
  email: string
  role_id: number
  role: string
  phone?: string | null
  created_at: string
  updated_at: string
}

interface Envelope<T> {
  success: boolean
  message: string
  data: T
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

export async function login(email: string, password: string): Promise<AuthResult> {
  if (USE_MOCK) return mockAuth.login(email, password)
  const { data } = await http.post<Envelope<{ token: string; user: BackendUser }>>('/auth/login', {
    email,
    password,
  })
  return { user: mapUser(data.data.user), token: data.data.token }
}

export async function register(input: {
  fullName: string
  email: string
  password: string
  phone?: string
  roleId?: string
}): Promise<AuthResult> {
  if (USE_MOCK) return mockAuth.register(input)
  const { data } = await http.post<Envelope<{ token: string; user: BackendUser }>>(
    '/auth/register',
    {
      full_name: input.fullName,
      email: input.email,
      password: input.password,
      phone: input.phone,
      role_id: input.roleId ? Number(input.roleId) : undefined,
    },
  )
  return { user: mapUser(data.data.user), token: data.data.token }
}

export async function logout(): Promise<void> {
  if (USE_MOCK) return
  try {
    await http.post('/auth/logout')
  } catch {
    // Logging out is best-effort; the client clears its session regardless.
  }
}

export async function me(): Promise<User> {
  if (USE_MOCK) throw new Error('me() is not available in mock mode')
  const { data } = await http.get<Envelope<BackendUser>>('/auth/me')
  return mapUser(data.data)
}

export async function getRegistrationData(): Promise<RoleOption[]> {
  if (USE_MOCK) {
    return [
      { id: '1', name: 'student_staff', description: 'Student or staff user' },
      { id: '2', name: 'maintenance_officer', description: 'Maintenance officer' },
      { id: '3', name: 'admin', description: 'System administrator' },
    ]
  }
  const { data } = await http.get<Envelope<{ roles: { id: number; name: string; description?: string }[] }>>(
    '/auth/registration-data',
  )
  return data.data.roles.map((r) => ({ id: String(r.id), name: r.name, description: r.description }))
}
