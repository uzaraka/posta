import api from './client'
import type { PaginatedResponse, Bounce, Suppression, ApiResponse } from './types'

export const bouncesApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<Bounce>>('/users/me/bounces', { params: { page, size } })
  },
  create(data: { email_id: number; recipient: string; type: string; reason: string }) {
    return api.post<ApiResponse<Bounce>>('/users/me/bounces', data)
  },
}

export const suppressionsApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<Suppression>>('/users/me/suppressions', { params: { page, size } })
  },
  create(data: { email: string; reason: string }) {
    return api.post<ApiResponse<Suppression>>('/users/me/suppressions', data)
  },
  delete(email: string) {
    return api.delete('/users/me/suppressions', { data: { email } })
  },
}
