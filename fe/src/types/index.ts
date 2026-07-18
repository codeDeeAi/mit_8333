/** Domain types for the University Maintenance & Service Request system. */

export type Role = 'student_staff' | 'maintenance_officer' | 'admin'

export type RequestStatus = 'pending' | 'assigned' | 'in_progress' | 'completed' | 'rejected'

export type Priority = 'low' | 'medium' | 'high'

export interface User {
  id: string
  fullName: string
  email: string
  role: Role
  phone?: string
  createdAt: string
}

export interface Category {
  id: string
  name: string
  description?: string
}

export interface StatusLog {
  id: string
  requestId: string
  changedBy: string
  changedByName: string
  oldStatus: RequestStatus | null
  newStatus: RequestStatus
  note?: string
  createdAt: string
}

export interface ServiceRequest {
  id: string
  title: string
  description: string
  categoryId: string
  categoryName: string
  createdBy: string
  createdByName: string
  location: string
  priority: Priority
  status: RequestStatus
  evidenceUrl?: string
  assignedOfficerId?: string
  assignedOfficerName?: string
  createdAt: string
  updatedAt: string
  logs?: StatusLog[]
}

export interface Notification {
  id: string
  userId: string
  message: string
  isRead: boolean
  createdAt: string
}

export interface AuditLog {
  id: string
  userId: string
  userName: string
  action: string
  entity: string
  entityId: string
  createdAt: string
}

export interface Paginated<T> {
  items: T[]
  page: number
  pageSize: number
  total: number
}

export interface RequestFilters {
  q?: string
  status?: RequestStatus | ''
  categoryId?: string
  priority?: Priority | ''
  page?: number
  pageSize?: number
}
