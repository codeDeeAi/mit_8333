import type {
  AuditLog,
  Category,
  Notification,
  Paginated,
  Priority,
  RequestFilters,
  RequestStatus,
  Role,
  ServiceRequest,
  User,
} from '@/types'
import { loadDb, persist, uid, now, type MockUser } from './mock/db'

/**
 * Service layer. Every function returns a Promise, matching a real HTTP client,
 * so swapping the mock for the Go API later is a drop-in change per function.
 */

const delay = (ms = 350) => new Promise((r) => setTimeout(r, ms))
const strip = ({ password, ...u }: MockUser): User => u

function pushAudit(userId: string, action: string, entity: string, entityId: string) {
  const db = loadDb()
  const user = db.users.find((u) => u.id === userId)
  db.audit.unshift({
    id: uid(),
    userId,
    userName: user?.fullName ?? 'System',
    action,
    entity,
    entityId,
    createdAt: now(),
  })
}

function notify(userId: string, message: string) {
  const db = loadDb()
  db.notifications.unshift({ id: uid(), userId, message, isRead: false, createdAt: now() })
}

/* ------------------------------- Auth ---------------------------------- */

export interface AuthResult {
  user: User
  token: string
}

export const authApi = {
  async login(email: string, password: string): Promise<AuthResult> {
    await delay()
    const db = loadDb()
    const found = db.users.find((u) => u.email.toLowerCase() === email.toLowerCase())
    if (!found || found.password !== password) {
      throw new Error('Invalid email or password')
    }
    pushAudit(found.id, 'LOGIN', 'user', found.id)
    persist()
    return { user: strip(found), token: `mock.jwt.${found.id}` }
  },

  async register(input: {
    fullName: string
    email: string
    password: string
    phone?: string
  }): Promise<AuthResult> {
    await delay()
    const db = loadDb()
    if (db.users.some((u) => u.email.toLowerCase() === input.email.toLowerCase())) {
      throw new Error('An account with this email already exists')
    }
    const user: MockUser = {
      id: uid(),
      fullName: input.fullName,
      email: input.email,
      role: 'student_staff',
      phone: input.phone,
      password: input.password,
      createdAt: now(),
    }
    db.users.push(user)
    pushAudit(user.id, 'REGISTER', 'user', user.id)
    persist()
    return { user: strip(user), token: `mock.jwt.${user.id}` }
  },
}

/* ----------------------------- Categories ------------------------------ */

export const categoryApi = {
  async list(): Promise<Category[]> {
    await delay(150)
    return [...loadDb().categories]
  },
}

/* --------------------------- Service Requests -------------------------- */

export const requestApi = {
  async list(
    filters: RequestFilters,
    scope?: { userId?: string; officerId?: string },
  ): Promise<Paginated<ServiceRequest>> {
    await delay()
    const db = loadDb()
    let items = [...db.requests].sort((a, b) => b.createdAt.localeCompare(a.createdAt))

    if (scope?.userId) items = items.filter((r) => r.createdBy === scope.userId)
    if (scope?.officerId) items = items.filter((r) => r.assignedOfficerId === scope.officerId)

    if (filters.q) {
      const q = filters.q.toLowerCase()
      items = items.filter(
        (r) =>
          r.title.toLowerCase().includes(q) ||
          r.description.toLowerCase().includes(q) ||
          r.location.toLowerCase().includes(q),
      )
    }
    if (filters.status) items = items.filter((r) => r.status === filters.status)
    if (filters.categoryId) items = items.filter((r) => r.categoryId === filters.categoryId)
    if (filters.priority) items = items.filter((r) => r.priority === filters.priority)

    const total = items.length
    const page = filters.page ?? 1
    const pageSize = filters.pageSize ?? 6
    const start = (page - 1) * pageSize
    return { items: items.slice(start, start + pageSize), page, pageSize, total }
  },

  async get(id: string): Promise<ServiceRequest> {
    await delay(200)
    const db = loadDb()
    const req = db.requests.find((r) => r.id === id)
    if (!req) throw new Error('Request not found')
    const logs = db.logs
      .filter((l) => l.requestId === id)
      .sort((a, b) => b.createdAt.localeCompare(a.createdAt))
    return { ...req, logs }
  },

  async create(input: {
    title: string
    description: string
    categoryId: string
    location: string
    priority: Priority
    evidenceUrl?: string
    createdBy: User
  }): Promise<ServiceRequest> {
    await delay()
    const db = loadDb()
    const cat = db.categories.find((c) => c.id === input.categoryId)
    if (!cat) throw new Error('Invalid category')
    const req: ServiceRequest = {
      id: uid(),
      title: input.title,
      description: input.description,
      categoryId: input.categoryId,
      categoryName: cat.name,
      createdBy: input.createdBy.id,
      createdByName: input.createdBy.fullName,
      location: input.location,
      priority: input.priority,
      status: 'pending',
      evidenceUrl: input.evidenceUrl,
      createdAt: now(),
      updatedAt: now(),
    }
    db.requests.unshift(req)
    db.logs.unshift({
      id: uid(),
      requestId: req.id,
      changedBy: input.createdBy.id,
      changedByName: input.createdBy.fullName,
      oldStatus: null,
      newStatus: 'pending',
      note: 'Request submitted.',
      createdAt: now(),
    })
    pushAudit(input.createdBy.id, 'CREATE_REQUEST', 'service_request', req.id)
    // Notify all admins.
    db.users.filter((u) => u.role === 'admin').forEach((a) => notify(a.id, `New request: "${req.title}"`))
    persist()
    return req
  },

  async updateStatus(
    id: string,
    newStatus: RequestStatus,
    actor: User,
    note?: string,
  ): Promise<ServiceRequest> {
    await delay()
    const db = loadDb()
    const req = db.requests.find((r) => r.id === id)
    if (!req) throw new Error('Request not found')
    const old = req.status
    req.status = newStatus
    req.updatedAt = now()
    db.logs.unshift({
      id: uid(),
      requestId: id,
      changedBy: actor.id,
      changedByName: actor.fullName,
      oldStatus: old,
      newStatus,
      note,
      createdAt: now(),
    })
    pushAudit(actor.id, 'STATUS_CHANGE', 'service_request', id)
    notify(req.createdBy, `Your request "${req.title}" is now ${newStatus.replace('_', ' ')}.`)
    persist()
    return req
  },

  async assign(id: string, officerId: string, actor: User): Promise<ServiceRequest> {
    await delay()
    const db = loadDb()
    const req = db.requests.find((r) => r.id === id)
    const officer = db.users.find((u) => u.id === officerId && u.role === 'maintenance_officer')
    if (!req) throw new Error('Request not found')
    if (!officer) throw new Error('Officer not found')
    req.assignedOfficerId = officer.id
    req.assignedOfficerName = officer.fullName
    const old = req.status
    if (req.status === 'pending') req.status = 'assigned'
    req.updatedAt = now()
    db.logs.unshift({
      id: uid(),
      requestId: id,
      changedBy: actor.id,
      changedByName: actor.fullName,
      oldStatus: old,
      newStatus: req.status,
      note: `Assigned to ${officer.fullName}.`,
      createdAt: now(),
    })
    pushAudit(actor.id, 'ASSIGN', 'service_request', id)
    notify(officer.id, `You have been assigned: "${req.title}"`)
    notify(req.createdBy, `Your request "${req.title}" was assigned to ${officer.fullName}.`)
    persist()
    return req
  },

  async remove(id: string, actor: User): Promise<void> {
    await delay()
    const db = loadDb()
    db.requests = db.requests.filter((r) => r.id !== id)
    pushAudit(actor.id, 'DELETE_REQUEST', 'service_request', id)
    persist()
  },
}

/* ------------------------------- Users --------------------------------- */

export const userApi = {
  async list(): Promise<User[]> {
    await delay()
    return loadDb().users.map(strip)
  },

  async officers(): Promise<User[]> {
    await delay(150)
    return loadDb()
      .users.filter((u) => u.role === 'maintenance_officer')
      .map(strip)
  },

  async updateRole(id: string, role: Role, actor: User): Promise<User> {
    await delay()
    const db = loadDb()
    const user = db.users.find((u) => u.id === id)
    if (!user) throw new Error('User not found')
    user.role = role
    pushAudit(actor.id, 'UPDATE_ROLE', 'user', id)
    persist()
    return strip(user)
  },

  async remove(id: string, actor: User): Promise<void> {
    await delay()
    const db = loadDb()
    db.users = db.users.filter((u) => u.id !== id)
    pushAudit(actor.id, 'DELETE_USER', 'user', id)
    persist()
  },
}

/* --------------------------- Notifications ----------------------------- */

export const notificationApi = {
  async list(userId: string): Promise<Notification[]> {
    await delay(150)
    return loadDb()
      .notifications.filter((n) => n.userId === userId)
      .sort((a, b) => b.createdAt.localeCompare(a.createdAt))
  },

  async markRead(id: string): Promise<void> {
    await delay(100)
    const db = loadDb()
    const n = db.notifications.find((x) => x.id === id)
    if (n) n.isRead = true
    persist()
  },

  async markAllRead(userId: string): Promise<void> {
    await delay(100)
    const db = loadDb()
    db.notifications.filter((n) => n.userId === userId).forEach((n) => (n.isRead = true))
    persist()
  },
}

/* ----------------------------- Reports --------------------------------- */

export interface SummaryStats {
  total: number
  byStatus: Record<RequestStatus, number>
  byCategory: { name: string; count: number }[]
}

export const reportApi = {
  async summary(): Promise<SummaryStats> {
    await delay()
    const db = loadDb()
    const byStatus = {
      pending: 0,
      assigned: 0,
      in_progress: 0,
      completed: 0,
      rejected: 0,
    } as Record<RequestStatus, number>
    const catMap = new Map<string, number>()
    for (const r of db.requests) {
      byStatus[r.status]++
      catMap.set(r.categoryName, (catMap.get(r.categoryName) ?? 0) + 1)
    }
    return {
      total: db.requests.length,
      byStatus,
      byCategory: [...catMap.entries()].map(([name, count]) => ({ name, count })),
    }
  },

  async allRequests(): Promise<ServiceRequest[]> {
    await delay(150)
    return [...loadDb().requests]
  },
}

/* ------------------------------ Audit ---------------------------------- */

export const auditApi = {
  async list(): Promise<AuditLog[]> {
    await delay()
    return [...loadDb().audit].sort((a, b) => b.createdAt.localeCompare(a.createdAt))
  },
}
