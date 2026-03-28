---
sidebar_position: 4
title: Scheduled Email
description: Schedule emails for future delivery
---

# Scheduled Email

Send emails at a specific future time using the `send_at` field.

## Usage

Include the `send_at` field with an ISO 8601 timestamp in the `SendEmail` request:

```bash
curl -X POST http://localhost:9000/api/v1/emails/send \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "sender@example.com",
    "to": ["recipient@example.com"],
    "subject": "Scheduled Reminder",
    "html": "<p>This is your scheduled reminder.</p>",
    "send_at": "2026-03-20T09:00:00Z"
  }'
```

## Response

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "scheduled"
  }
}
```

The email status will be `scheduled` until the specified time, then it transitions to `queued` and follows the normal delivery pipeline.

## Notes

- The `send_at` time must be in the future (UTC)
- The `send_at` field is only available on the `/emails/send` endpoint
- Scheduled emails can be tracked via the [email status](/docs/email-sending/email-status) endpoint
