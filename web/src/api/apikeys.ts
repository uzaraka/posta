import api from './client'
import type { ApiResponse, PaginatedResponse, ApiKey, ApiKeyCreateResponse } from './types'

export const apiKeysApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<ApiKey>>('/users/me/api-keys', { params: { page, size } })
  },
  create(name: string, allowedIPs?: string[], expiresInDays?: number) {
    const body: Record<string, any> = { name }
    if (allowedIPs && allowedIPs.length > 0) body.allowed_ips = allowedIPs
    if (expiresInDays !== undefined) body.expires_in_days = expiresInDays
    return api.post<ApiResponse<ApiKeyCreateResponse>>('/users/me/api-keys', body)
  },
  revoke(id: number) {
    return api.put(`/users/me/api-keys/${id}/revoke`)
  },
  delete(id: number) {
    return api.delete(`/users/me/api-keys/${id}`)
  },
}
