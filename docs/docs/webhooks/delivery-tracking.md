---
sidebar_position: 3
title: Delivery Tracking
description: Track webhook delivery history
---

# Webhook Delivery Tracking

Monitor the delivery status of your webhook notifications.

## View Delivery History

```
GET /api/v1/users/me/webhook-deliveries?page=1&size=20
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "url": "https://your-app.com/webhooks/posta",
      "status": 200,
      "retry_count": 0,
      "response": "{\"received\": true}",
      "timestamp": "2026-01-01T00:00:01Z"
    },
    {
      "url": "https://your-app.com/webhooks/posta",
      "status": 500,
      "retry_count": 2,
      "response": "Internal Server Error",
      "timestamp": "2026-01-01T00:01:00Z"
    }
  ]
}
```

## Delivery Details

| Field | Description |
|-------|-------------|
| `url` | Webhook endpoint URL |
| `status` | HTTP status code of the response |
| `retry_count` | Number of delivery attempts |
| `response` | Response body from your server |
| `timestamp` | When the delivery was attempted |

## Retry Policy

- **Max retries:** Configured via `POSTA_WEBHOOK_MAX_RETRIES` (default: 3)
- **Timeout:** Configured via `POSTA_WEBHOOK_TIMEOUT_SECS` (default: 10s)
- **Backoff:** Exponential backoff between retries
- **Success:** Any `2xx` status code is considered successful

## Data Retention

Webhook delivery history is retained according to the `webhook_retention_days` platform setting, configurable by administrators.
