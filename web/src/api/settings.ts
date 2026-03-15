import api from './client'
import type { ApiResponse, AdminSetting, AdminSettingInput, UserSettings } from './types'

export const settingsApi = {
  // Admin platform settings
  getAdminSettings() {
    return api.get<ApiResponse<AdminSetting[]>>('/admin/settings')
  },
  updateAdminSettings(settings: AdminSettingInput[]) {
    return api.put<ApiResponse<AdminSetting[]>>('/admin/settings', { settings })
  },

  // User settings
  getUserSettings() {
    return api.get<ApiResponse<UserSettings>>('/users/me/settings')
  },
  updateUserSettings(data: Partial<Omit<UserSettings, 'id' | 'user_id' | 'created_at' | 'updated_at'>>) {
    return api.put<ApiResponse<UserSettings>>('/users/me/settings', data)
  },
}
