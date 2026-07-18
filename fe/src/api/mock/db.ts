import type {
  AuditLog,
  Category,
  Notification,
  Priority,
  RequestStatus,
  ServiceRequest,
  StatusLog,
  User,
} from '@/types'

/**
 * In-memory mock database, persisted to localStorage so the demo keeps state
 * across reloads. Mirrors the shape of the Go/Postgres backend one-to-one.
 */

const LS_KEY = 'mock_db_v1'

interface MockUser extends User {
  password: string
}

interface DbShape {
  users: MockUser[]
  categories: Category[]
  requests: ServiceRequest[]
  logs: StatusLog[]
  notifications: Notification[]
  audit: AuditLog[]
}

export const uid = () => Math.random().toString(36).slice(2, 10)
const now = () => new Date().toISOString()
const daysAgo = (n: number) => new Date(Date.now() - n * 86400000).toISOString()

function seed(): DbShape {
  const admin: MockUser = {
    id: 'u-admin',
    fullName: 'Amara Okafor',
    email: 'admin@miva.edu',
    role: 'admin',
    phone: '+234 800 000 0001',
    password: 'password123',
    createdAt: daysAgo(120),
  }
  const officer: MockUser = {
    id: 'u-officer',
    fullName: 'Tunde Balogun',
    email: 'officer@miva.edu',
    role: 'maintenance_officer',
    phone: '+234 800 000 0002',
    password: 'password123',
    createdAt: daysAgo(90),
  }
  const officer2: MockUser = {
    id: 'u-officer2',
    fullName: 'Grace Eze',
    email: 'grace@miva.edu',
    role: 'maintenance_officer',
    phone: '+234 800 000 0003',
    password: 'password123',
    createdAt: daysAgo(80),
  }
  const student: MockUser = {
    id: 'u-student',
    fullName: 'Chidi Nwosu',
    email: 'student@miva.edu',
    role: 'student_staff',
    phone: '+234 800 000 0004',
    password: 'password123',
    createdAt: daysAgo(60),
  }

  const categories: Category[] = [
    { id: 'c-elec', name: 'Electricity', description: 'Power, lighting and electrical faults' },
    { id: 'c-furn', name: 'Furniture', description: 'Damaged desks, chairs and fittings' },
    { id: 'c-plumb', name: 'Plumbing', description: 'Leaking pipes, taps and drainage' },
    { id: 'c-net', name: 'Internet', description: 'Wi-Fi and network connectivity' },
    { id: 'c-equip', name: 'Classroom Equipment', description: 'Projectors, boards and AV gear' },
    { id: 'c-hostel', name: 'Hostel Maintenance', description: 'General hostel repairs' },
  ]

  const mkReq = (
    id: string,
    title: string,
    description: string,
    categoryId: string,
    location: string,
    priority: Priority,
    status: RequestStatus,
    created: string,
    officerRef?: MockUser,
  ): ServiceRequest => {
    const cat = categories.find((c) => c.id === categoryId)!
    return {
      id,
      title,
      description,
      categoryId,
      categoryName: cat.name,
      createdBy: student.id,
      createdByName: student.fullName,
      location,
      priority,
      status,
      assignedOfficerId: officerRef?.id,
      assignedOfficerName: officerRef?.fullName,
      createdAt: created,
      updatedAt: created,
    }
  }

  const requests: ServiceRequest[] = [
    mkReq('r-1', 'Flickering lights in Lecture Hall B', 'The ceiling lights flicker constantly and two are completely dead.', 'c-elec', 'Faculty of Science, Hall B', 'high', 'in_progress', daysAgo(3), officer),
    mkReq('r-2', 'Broken chair in Library', 'A reading chair on the second floor has a cracked leg and is unsafe.', 'c-furn', 'Main Library, 2nd Floor', 'low', 'pending', daysAgo(1)),
    mkReq('r-3', 'Leaking pipe in Hostel C bathroom', 'Water leaks continuously under the sink, flooding the floor.', 'c-plumb', 'Hostel Block C, Room 12', 'high', 'assigned', daysAgo(2), officer2),
    mkReq('r-4', 'No internet in Computer Lab 1', 'Wi-Fi has been down for two days across the entire lab.', 'c-net', 'ICT Building, Lab 1', 'medium', 'completed', daysAgo(8), officer),
    mkReq('r-5', 'Projector not turning on', 'The projector in Room 204 will not power on before lectures.', 'c-equip', 'Engineering Block, Room 204', 'medium', 'pending', daysAgo(1)),
    mkReq('r-6', 'Faulty socket in hostel room', 'Wall socket sparks when anything is plugged in.', 'c-elec', 'Hostel Block A, Room 30', 'high', 'completed', daysAgo(15), officer),
  ]

  const logs: StatusLog[] = [
    { id: uid(), requestId: 'r-1', changedBy: admin.id, changedByName: admin.fullName, oldStatus: 'pending', newStatus: 'assigned', note: 'Assigned to Tunde.', createdAt: daysAgo(3) },
    { id: uid(), requestId: 'r-1', changedBy: officer.id, changedByName: officer.fullName, oldStatus: 'assigned', newStatus: 'in_progress', note: 'On site, sourcing replacement tubes.', createdAt: daysAgo(2) },
    { id: uid(), requestId: 'r-4', changedBy: officer.id, changedByName: officer.fullName, oldStatus: 'in_progress', newStatus: 'completed', note: 'Replaced router, connectivity restored.', createdAt: daysAgo(6) },
  ]

  const notifications: Notification[] = [
    { id: uid(), userId: officer.id, message: 'You have been assigned: "Flickering lights in Lecture Hall B"', isRead: false, createdAt: daysAgo(3) },
    { id: uid(), userId: student.id, message: 'Your request "No internet in Computer Lab 1" was completed.', isRead: false, createdAt: daysAgo(6) },
  ]

  const audit: AuditLog[] = [
    { id: uid(), userId: student.id, userName: student.fullName, action: 'CREATE_REQUEST', entity: 'service_request', entityId: 'r-1', createdAt: daysAgo(3) },
    { id: uid(), userId: admin.id, userName: admin.fullName, action: 'ASSIGN', entity: 'service_request', entityId: 'r-1', createdAt: daysAgo(3) },
    { id: uid(), userId: officer.id, userName: officer.fullName, action: 'STATUS_CHANGE', entity: 'service_request', entityId: 'r-4', createdAt: daysAgo(6) },
  ]

  return { users: [admin, officer, officer2, student], categories, requests, logs, notifications, audit }
}

let db: DbShape

export function loadDb(): DbShape {
  if (db) return db
  const raw = localStorage.getItem(LS_KEY)
  if (raw) {
    try {
      db = JSON.parse(raw)
      return db
    } catch {
      /* fall through to reseed */
    }
  }
  db = seed()
  persist()
  return db
}

export function persist() {
  localStorage.setItem(LS_KEY, JSON.stringify(db))
}

export function resetDb() {
  db = seed()
  persist()
}

export { now, daysAgo }
export type { MockUser, DbShape }
