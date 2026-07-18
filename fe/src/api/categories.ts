import http, { USE_MOCK } from './http'
import { categoryApi as mockCategoryApi } from './services'
import type { Category } from '@/types'

interface Envelope<T> {
  success: boolean
  message: string
  data: T
}

interface BackendCategory {
  id: number
  name: string
  description?: string | null
}

/** Categories API — real backend with mock fallback. Same shape as the mock. */
export const categoryApi = {
  async list(): Promise<Category[]> {
    if (USE_MOCK) return mockCategoryApi.list()
    const { data } = await http.get<Envelope<BackendCategory[]>>('/categories')
    return data.data.map((c) => ({
      id: String(c.id),
      name: c.name,
      description: c.description ?? undefined,
    }))
  },
}
