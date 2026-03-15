import api from './client'
import type { ApiResponse, DashboardStats } from './types'

export const dashboardApi = {
  getStats() {
    return api.get<ApiResponse<DashboardStats>>('/users/me/dashboard/stats')
  },
}
