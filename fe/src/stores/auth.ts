import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'
import type { Role, User } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (s): boolean => !!s.token && !!s.user,
    role: (s): Role | null => s.user?.role ?? null,
    isAdmin: (s): boolean => s.user?.role === 'admin',
    isOfficer: (s): boolean => s.user?.role === 'maintenance_officer',
    isStudent: (s): boolean => s.user?.role === 'student_staff',
    /** Landing route for the current user's role. */
    homeRoute(): string {
      switch (this.user?.role) {
        case 'admin':
          return '/admin'
        case 'maintenance_officer':
          return '/officer'
        default:
          return '/app'
      }
    },
  },

  actions: {
    async login(email: string, password: string) {
      this.loading = true
      try {
        const { user, token } = await authApi.login(email, password)
        this.user = user
        this.token = token
      } finally {
        this.loading = false
      }
    },

    async register(input: {
      fullName: string
      email: string
      password: string
      phone?: string
      roleId?: string
    }) {
      this.loading = true
      try {
        const { user, token } = await authApi.register(input)
        this.user = user
        this.token = token
      } finally {
        this.loading = false
      }
    },

    async fetchMe() {
      this.user = await authApi.me()
    },

    async logout() {
      await authApi.logout()
      this.user = null
      this.token = null
    },
  },

  // Persist the whole auth slice under the `auth` key (read by the axios interceptor).
  persist: {
    key: 'auth',
    pick: ['user', 'token'],
  },
})
