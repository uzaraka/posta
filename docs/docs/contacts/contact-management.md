---
sidebar_position: 1
title: Contact Management
description: Track email recipients
---

# Contact Management

Posta automatically tracks recipients when you send emails. You can search and view contact details from the dashboard.

## List Contacts

```
GET /api/v1/users/me/contacts?page=1&size=20&search=alice
```

| Parameter | Description |
|-----------|-------------|
| `page` | Page number |
| `size` | Items per page |
| `search` | Search by email or name |

Response:

```json
{
  "success": true,
  "data": [
    {
      "email": "alice@example.com",
      "name": "Alice",
      "sent_count": 42,
      "failed_count": 1,
      "suppressed": false,
      "last_sent_at": "2026-01-15T10:00:00Z"
    }
  ]
}
```

## Contact Details

```
GET /api/v1/users/me/contacts/{contactId}
```

Returns the contact with suppression status and sending statistics.

## Automatic Tracking

Contacts are created automatically when you send to a new email address. Posta tracks:

- **Sent count** — Total emails successfully sent
- **Failed count** — Total delivery failures
- **Suppression status** — Whether the contact is suppressed
- **Last sent** — Timestamp of the most recent email
