import { requestApi } from './requests'
import { USE_MOCK } from './http'
import { reportApi as mockReportApi } from './services'
import type { RequestStatus, ServiceRequest } from '@/types'

export interface SummaryStats {
  total: number
  byStatus: Record<RequestStatus, number>
  byCategory: { name: string; count: number }[]
}

async function collectAllRequests(pageSize = 100): Promise<ServiceRequest[]> {
  const all: ServiceRequest[] = []
  let page = 1

  while (true) {
    const res = await requestApi.list({ page, pageSize })
    all.push(...res.items)

    if (all.length >= res.total || res.items.length === 0) {
      break
    }

    page += 1
  }

  return all
}

function summarize(requests: ServiceRequest[]): SummaryStats {
  const byStatus: Record<RequestStatus, number> = {
    pending: 0,
    assigned: 0,
    in_progress: 0,
    completed: 0,
    rejected: 0,
  }

  const categoryCount = new Map<string, number>()
  for (const request of requests) {
    byStatus[request.status] += 1
    categoryCount.set(request.categoryName, (categoryCount.get(request.categoryName) ?? 0) + 1)
  }

  return {
    total: requests.length,
    byStatus,
    byCategory: Array.from(categoryCount.entries()).map(([name, count]) => ({ name, count })),
  }
}

export const reportApi = {
  async summary(): Promise<SummaryStats> {
    if (USE_MOCK) return mockReportApi.summary()
    const requests = await collectAllRequests()
    return summarize(requests)
  },

  async allRequests(): Promise<ServiceRequest[]> {
    if (USE_MOCK) return mockReportApi.allRequests()
    return collectAllRequests()
  },
}
