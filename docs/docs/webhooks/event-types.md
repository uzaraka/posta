---
sidebar_position: 2
title: Event Types
description: Available webhook event types
---

# Event Types

## Email Events

Subscribe to these events when creating webhooks:

| Event | Description |
|-------|-------------|
| `email.sent` | Email was successfully delivered to the SMTP server |
| `email.failed` | Email delivery failed after all retries |
| `email.bounced` | Email bounced (hard or soft bounce) |

## Event Payload Structure

All events share a common structure:

```json
{
  "event": "email.sent",
  "timestamp": "2026-01-01T00:00:01Z",
  "data": {
    "email_id": "uuid",
    "to": "recipient@example.com",
    "from": "sender@example.com",
    "subject": "Hello",
    "status": "sent",
    "error_message": null
  }
}
```

### `email.sent`

Fired when an email is successfully delivered:

```json
{
  "event": "email.sent",
  "data": {
    "email_id": "uuid",
    "to": "recipient@example.com",
    "status": "sent",
    "sent_at": "2026-01-01T00:00:01Z"
  }
}
```

### `email.failed`

Fired when delivery fails after all retries:

```json
{
  "event": "email.failed",
  "data": {
    "email_id": "uuid",
    "to": "recipient@example.com",
    "status": "failed",
    "error_message": "Connection refused",
    "retry_count": 3
  }
}
```

### `email.bounced`

Fired when an email bounces:

```json
{
  "event": "email.bounced",
  "data": {
    "email_id": "uuid",
    "to": "recipient@example.com",
    "bounce_type": "hard",
    "message": "Mailbox not found"
  }
}
```
