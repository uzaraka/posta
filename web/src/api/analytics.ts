import api from './client'
import type { ApiResponse, AnalyticsResponse } from './types'

export const analyticsApi = {
  user(from?: string, to?: string, status?: string) {
    const params: Record<string, string> = {}
    if (from) params.from = from
    if (to) params.to = to
    if (status) params.status = status
    return api.get<ApiResponse<AnalyticsResponse>>('/users/me/analytics', { params })
  },
  admin(from?: string, to?: string, status?: string) {
    const params: Record<string, string> = {}
    if (from) params.from = from
    if (to) params.to = to
    if (status) params.status = status
    return api.get<ApiResponse<AnalyticsResponse>>('/admin/analytics', { params })
  },
}
