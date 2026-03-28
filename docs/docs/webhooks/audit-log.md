---
sidebar_position: 4
title: Audit Log
description: Track user and system events
---

# Audit Log

Posta maintains an audit trail of user and system events.

## View Audit Log

```
GET /api/v1/users/me/audit-log?page=1&size=20&category=email
```

### Query Parameters

| Parameter | Description |
|-----------|-------------|
| `page` | Page number |
| `size` | Items per page |
| `category` | Filter by category: `user`, `email`, `system`, `audit` |

Response:

```json
{
  "success": true,
  "data": [
    {
      "type": "email.sent",
      "actor_name": "admin@example.com",
      "message": "Email sent to recipient@example.com",
      "timestamp": "2026-01-01T00:00:01Z"
    },
    {
      "type": "user.login",
      "actor_name": "admin@example.com",
      "message": "User logged in from 203.0.113.42",
      "timestamp": "2026-01-01T00:00:00Z"
    }
  ]
}
```

## Event Categories

| Category | Events |
|----------|--------|
| `user` | Login, logout, profile updates, password changes, 2FA changes |
| `email` | Email sent, failed, retried, bounced |
| `system` | Server started, configuration changes |
| `audit` | API key created/revoked, template changes, SMTP changes |

## Admin: Platform-Wide Events

Administrators can view all events across all users:

```
GET /api/v1/admin/events?page=1&size=20&category=email
```

### Real-Time Event Stream (SSE)

Subscribe to real-time events via Server-Sent Events:

```
GET /api/v1/admin/events/stream?token=<jwt-token>
```

```javascript
const eventSource = new EventSource(
  'http://localhost:9000/api/v1/admin/events/stream?token=your-jwt'
);

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
};
```
