import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { Role } from '@/types'

const AppLayout = () => import('@/components/layout/AppLayout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      // Student / Staff
      {
        path: 'app',
        name: 'student-home',
        component: () => import('@/views/student/StudentDashboard.vue'),
        meta: { roles: ['student_staff'] as Role[] },
      },
      {
        path: 'app/new',
        name: 'new-request',
        component: () => import('@/views/student/NewRequestView.vue'),
        meta: { roles: ['student_staff'] as Role[] },
      },
      {
        path: 'app/requests',
        name: 'my-requests',
        component: () => import('@/views/student/MyRequestsView.vue'),
        meta: { roles: ['student_staff'] as Role[] },
      },
      // Maintenance Officer
      {
        path: 'officer',
        name: 'officer-home',
        component: () => import('@/views/officer/OfficerDashboard.vue'),
        meta: { roles: ['maintenance_officer'] as Role[] },
      },
      {
        path: 'officer/assigned',
        name: 'officer-assigned',
        component: () => import('@/views/officer/AssignedView.vue'),
        meta: { roles: ['maintenance_officer'] as Role[] },
      },
      // Administrator
      {
        path: 'admin',
        name: 'admin-home',
        component: () => import('@/views/admin/AdminDashboard.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      {
        path: 'admin/requests',
        name: 'admin-requests',
        component: () => import('@/views/admin/AdminRequestsView.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      {
        path: 'admin/users',
        name: 'admin-users',
        component: () => import('@/views/admin/UsersView.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      {
        path: 'admin/reports',
        name: 'admin-reports',
        component: () => import('@/views/admin/ReportsView.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      {
        path: 'admin/audit',
        name: 'admin-audit',
        component: () => import('@/views/admin/AuditView.vue'),
        meta: { roles: ['admin'] as Role[] },
      },
      // Shared request detail (any authenticated role)
      {
        path: 'requests/:id',
        name: 'request-detail',
        component: () => import('@/views/RequestDetailView.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  // Root → role home (or login)
  if (to.path === '/') {
    return auth.isAuthenticated ? auth.homeRoute : '/login'
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && auth.isAuthenticated) {
    return auth.homeRoute
  }

  // Role gate
  const roles = to.meta.roles as Role[] | undefined
  if (roles && auth.role && !roles.includes(auth.role)) {
    return auth.homeRoute
  }

  return true
})

export default router
