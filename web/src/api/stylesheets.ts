import api from './client'
import type { ApiResponse, PaginatedResponse, StyleSheet, StyleSheetInput } from './types'

export const stylesheetsApi = {
  list(page = 0, size = 100) {
    return api.get<PaginatedResponse<StyleSheet>>('/users/me/stylesheets', { params: { page, size } })
  },
  create(data: StyleSheetInput) {
    return api.post<ApiResponse<StyleSheet>>('/users/me/stylesheets', data)
  },
  update(id: number, data: StyleSheetInput) {
    return api.put<ApiResponse<StyleSheet>>(`/users/me/stylesheets/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/users/me/stylesheets/${id}`)
  },
}
