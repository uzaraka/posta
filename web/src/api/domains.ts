import api from './client'
import type { ApiResponse, PaginatedResponse, Domain } from './types'

export const domainsApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<Domain>>('/users/me/domains', { params: { page, size } })
  },
  get(id: number) {
    return api.get<ApiResponse<Domain>>(`/users/me/domains/${id}`)
  },
  create(domain: string) {
    return api.post<ApiResponse<Domain>>('/users/me/domains', { domain })
  },
  verify(id: number) {
    return api.post<ApiResponse<Domain>>(`/users/me/domains/${id}/verify`)
  },
  delete(id: number) {
    return api.delete(`/users/me/domains/${id}`)
  },
}
