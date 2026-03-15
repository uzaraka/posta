import api from './client'
import type { PaginatedResponse, Event } from './types'

export const auditApi = {
  list(page = 0, size = 20, category?: string) {
    const params: Record<string, any> = { page, size }
    if (category) params.category = category
    return api.get<PaginatedResponse<Event>>('/users/me/audit-log', { params })
  },
}
