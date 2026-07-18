import axios from 'axios'

/**
 * Configured axios instance for the Go backend.
 *
 * The UI currently runs against an in-memory mock (see `@/api/mock`) so it is
 * fully demoable before the backend exists. When the Go API is ready, flip
 * `VITE_USE_MOCK=false` and the service modules will call these real endpoints.
 */
const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

// Attach JWT from persisted auth store on every request.
http.interceptors.request.use((config) => {
  const raw = localStorage.getItem('auth')
  if (raw) {
    try {
      const token = JSON.parse(raw)?.token
      if (token) config.headers.Authorization = `Bearer ${token}`
    } catch {
      /* ignore malformed storage */
    }
  }
  return config
})

// On 401, clear session and bounce to login.
http.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error?.response?.status === 401) {
      localStorage.removeItem('auth')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export default http

export const USE_MOCK = import.meta.env.VITE_USE_MOCK !== 'false'
