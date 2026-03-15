import api from './client'
import type { ApiResponse, PaginatedResponse, Webhook, WebhookDelivery, WebhookInput } from './types'

export const webhooksApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<Webhook>>('/users/me/webhooks', { params: { page, size } })
  },
  create(data: WebhookInput) {
    return api.post<ApiResponse<Webhook>>('/users/me/webhooks', data)
  },
  delete(id: number) {
    return api.delete(`/users/me/webhooks/${id}`)
  },
}

export const webhookDeliveriesApi = {
  list(page = 0, size = 20) {
    return api.get<PaginatedResponse<WebhookDelivery>>('/users/me/webhook-deliveries', { params: { page, size } })
  },
}
