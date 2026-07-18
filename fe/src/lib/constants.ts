import type { Priority, RequestStatus, Role } from '@/types'

export const ROLE_LABELS: Record<Role, string> = {
  student_staff: 'Student / Staff',
  maintenance_officer: 'Maintenance Officer',
  admin: 'Administrator',
}

export const STATUS_LABELS: Record<RequestStatus, string> = {
  pending: 'Pending',
  assigned: 'Assigned',
  in_progress: 'In Progress',
  completed: 'Completed',
  rejected: 'Rejected',
}

/** Tailwind classes for status pills (dark theme). */
export const STATUS_STYLES: Record<RequestStatus, string> = {
  pending: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20',
  assigned: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
  in_progress: 'bg-purple-500/10 text-purple-400 border-purple-500/20',
  completed: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  rejected: 'bg-red-500/10 text-red-400 border-red-500/20',
}

export const PRIORITY_LABELS: Record<Priority, string> = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
}

export const PRIORITY_STYLES: Record<Priority, string> = {
  low: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
  medium: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  high: 'bg-red-500/10 text-red-400 border-red-500/20',
}

export const ALL_STATUSES: RequestStatus[] = [
  'pending',
  'assigned',
  'in_progress',
  'completed',
  'rejected',
]

export const ALL_PRIORITIES: Priority[] = ['low', 'medium', 'high']
