import http, { USE_MOCK } from './http'
import { requestApi as mockRequestApi } from './services'
import { categoryApi } from './categories'
import type {
  Paginated,
  Priority,
  RequestFilters,
  RequestStatus,
  ServiceRequest,
  StatusLog,
  User,
} from '@/types'

interface Envelope<T> {
  success: boolean
  message: string
  data: T
}

interface BackendRequest {
  id: number
  title: string
  description: string
  category_id: number
  category_name?: string
  created_by: number
  created_by_name?: string
  location: string
  priority: Priority
  status: RequestStatus
  evidence_path?: string | null
  assigned_officer_id?: number | null
  assigned_officer_name?: string | null
  created_at: string
  updated_at: string
}

interface BackendLog {
  id: number
  request_id: number
  changed_by: number
  old_status?: string | null
  new_status: RequestStatus
  note?: string | null
  created_at: string
}

interface BackendList {
  items: BackendRequest[]
  pagination: { page: number; page_size: number; total: number; total_pages: number }
}

interface BackendDetail {
  request: BackendRequest
  status_logs: BackendLog[]
}

// The backend returns lean records (ids only). We resolve category names from
// the categories endpoint and the requester's name from the persisted session
// (the common case: a user viewing their own request). A users endpoint would
// let us resolve every actor's name.
let categoryNames: Map<string, string> | null = null

async function nameMap(): Promise<Map<string, string>> {
  if (!categoryNames) {
    const cats = await categoryApi.list()
    categoryNames = new Map(cats.map((c) => [c.id, c.name]))
  }
  return categoryNames
}

function sessionUser(): { id: string; fullName: string } | null {
  try {
    const raw = JSON.parse(localStorage.getItem('auth') || 'null')
    return raw?.user ?? null
  } catch {
    return null
  }
}

function personName(userId: number): string {
  const me = sessionUser()
  return me && me.id === String(userId) ? me.fullName : `User #${userId}`
}

function mapRequest(b: BackendRequest, names: Map<string, string>): ServiceRequest {
  return {
    id: String(b.id),
    title: b.title,
    description: b.description,
    categoryId: String(b.category_id),
    categoryName: b.category_name || names.get(String(b.category_id)) || 'Uncategorized',
    createdBy: String(b.created_by),
    createdByName: b.created_by_name || personName(b.created_by),
    location: b.location,
    priority: b.priority,
    status: b.status,
    // evidence_path is a server-side path, not yet served over HTTP.
    evidenceUrl: undefined,
    assignedOfficerId: b.assigned_officer_id ? String(b.assigned_officer_id) : undefined,
    assignedOfficerName: b.assigned_officer_name ?? undefined,
    createdAt: b.created_at,
    updatedAt: b.updated_at,
  }
}

function mapLog(l: BackendLog): StatusLog {
  return {
    id: String(l.id),
    requestId: String(l.request_id),
    changedBy: String(l.changed_by),
    changedByName: personName(l.changed_by),
    oldStatus: (l.old_status as RequestStatus) ?? null,
    newStatus: l.new_status,
    note: l.note ?? undefined,
    createdAt: l.created_at,
  }
}

export interface CreateRequestInput {
  title: string
  description: string
  categoryId: string
  location: string
  priority: Priority
  evidenceUrl?: string
  evidenceFile?: File
  createdBy: User
}

export interface UploadEvidenceInput {
  requestId: string
  file: File
}

/** Requests API — real backend with mock fallback. Same shape as the mock. */
export const requestApi = {
  async list(
    filters: RequestFilters,
    scope?: { userId?: string; officerId?: string },
  ): Promise<Paginated<ServiceRequest>> {
    // The backend scopes results by the caller's role automatically, so `scope`
    // is only used by the mock.
    if (USE_MOCK) return mockRequestApi.list(filters, scope)

    const params: Record<string, string | number> = {}
    if (filters.q) params.q = filters.q
    if (filters.status) params.status = filters.status
    if (filters.categoryId) params.category_id = Number(filters.categoryId)
    if (filters.priority) params.priority = filters.priority
    if (filters.page) params.page = filters.page
    if (filters.pageSize) params.page_size = filters.pageSize

    const { data } = await http.get<Envelope<BackendList>>('/requests', { params })
    const names = await nameMap()
    return {
      items: data.data.items.map((r) => mapRequest(r, names)),
      page: data.data.pagination.page,
      pageSize: data.data.pagination.page_size,
      total: data.data.pagination.total,
    }
  },

  async get(id: string): Promise<ServiceRequest> {
    if (USE_MOCK) return mockRequestApi.get(id)
    const { data } = await http.get<Envelope<BackendDetail>>(`/requests/${id}`)
    const names = await nameMap()
    const req = mapRequest(data.data.request, names)
    req.logs = data.data.status_logs.map(mapLog)
    return req
  },

  async create(input: CreateRequestInput): Promise<ServiceRequest> {
    if (USE_MOCK) return mockRequestApi.create(input)

    const form = new FormData()
    form.append('title', input.title)
    form.append('description', input.description)
    form.append('category_id', String(Number(input.categoryId)))
    form.append('location', input.location)
    form.append('priority', input.priority)
    if (input.evidenceFile) form.append('evidence', input.evidenceFile)

    const { data } = await http.post<Envelope<BackendRequest>>('/requests', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return mapRequest(data.data, await nameMap())
  },

  async updateStatus(
    id: string,
    newStatus: RequestStatus,
    actor: User,
    note?: string,
  ): Promise<ServiceRequest> {
    if (USE_MOCK) return mockRequestApi.updateStatus(id, newStatus, actor, note)
    const { data } = await http.put<Envelope<BackendRequest>>(`/requests/${id}/status`, {
      status: newStatus,
      note,
    })
    return mapRequest(data.data, await nameMap())
  },

  async assign(id: string, officerId: string, actor: User): Promise<ServiceRequest> {
    if (USE_MOCK) return mockRequestApi.assign(id, officerId, actor)
    const { data } = await http.post<Envelope<BackendRequest>>(`/requests/${id}/assign`, {
      officer_id: Number(officerId),
    })
    return mapRequest(data.data, await nameMap())
  },

  async uploadEvidence(input: UploadEvidenceInput): Promise<ServiceRequest> {
    if (USE_MOCK) {
      throw new Error('Evidence replacement is not available in mock mode')
    }

    const form = new FormData()
    form.append('evidence', input.file)

    const { data } = await http.post<Envelope<BackendRequest>>(
      `/requests/${input.requestId}/evidence`,
      form,
      { headers: { 'Content-Type': 'multipart/form-data' } },
    )

    return mapRequest(data.data, await nameMap())
  },

  async remove(id: string, actor: User): Promise<void> {
    if (USE_MOCK) return mockRequestApi.remove(id, actor)
    await http.delete(`/requests/${id}`)
  },
}
