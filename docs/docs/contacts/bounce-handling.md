---
sidebar_position: 3
title: Bounce Handling
description: Track and manage email bounces
---

# Bounce Handling

Posta tracks email bounces and automatically suppresses addresses that hard bounce.

## Bounce Types

| Type | Description | Action |
|------|-------------|--------|
| `hard` | Permanent failure (e.g., mailbox doesn't exist) | Auto-suppressed |
| `soft` | Temporary failure (e.g., mailbox full) | Tracked, not suppressed |
| `complaint` | Recipient reported as spam | Auto-suppressed |

## Record a Bounce

```
POST /api/v1/users/me/bounces
```

```json
{
  "email": "bounced@example.com",
  "bounce_type": "hard",
  "message": "550 5.1.1 The email account does not exist"
}
```

## List Bounces

```
GET /api/v1/users/me/bounces?page=1&size=20
```

## Automatic Suppression

When a **hard bounce** or **complaint** is recorded:

1. The email address is automatically added to the [suppression list](/docs/contacts/suppression-list)
2. Future sends to that address are blocked with status `suppressed`
3. Batch sends skip suppressed addresses automatically

**Soft bounces** are tracked but do not trigger automatic suppression. They may resolve on their own (e.g., when a full mailbox is cleared).
