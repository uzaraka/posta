import api from './client'
import type { ApiResponse, PaginatedResponse, SmtpServer, SmtpServerInput } from './types'

export const smtpApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<SmtpServer>>('/users/me/smtp-servers', { params: { page, size } })
  },
  get(id: number) {
    return api.get<ApiResponse<SmtpServer>>(`/users/me/smtp-servers/${id}`)
  },
  create(data: SmtpServerInput) {
    return api.post<ApiResponse<SmtpServer>>('/users/me/smtp-servers', data)
  },
  update(id: number, data: Partial<SmtpServerInput>) {
    return api.put<ApiResponse<SmtpServer>>(`/users/me/smtp-servers/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/users/me/smtp-servers/${id}`)
  },
  test(id: number) {
    return api.post(`/users/me/smtp-servers/${id}/test`)
  },
}
