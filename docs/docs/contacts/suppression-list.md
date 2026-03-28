---
sidebar_position: 4
title: Suppression List
description: Manage suppressed email addresses
---

# Suppression List

The suppression list prevents emails from being sent to specific addresses. Addresses are added automatically on hard bounces/complaints, or manually.

## Add to Suppression List

```
POST /api/v1/users/me/suppressions
```

```json
{
  "email": "unsubscribed@example.com",
  "reason": "User requested removal"
}
```

Returns `409 Conflict` if already suppressed.

## List Suppressed Addresses

```
GET /api/v1/users/me/suppressions?page=1&size=20
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "email": "unsubscribed@example.com",
      "reason": "User requested removal",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

## Remove from Suppression List

```
DELETE /api/v1/users/me/suppressions
```

```json
{
  "email": "resubscribed@example.com"
}
```

## How Suppression Works

When sending an email:

1. Posta checks the recipient against the suppression list
2. If suppressed, the email is marked as `suppressed` and not delivered
3. In batch sends, suppressed recipients are skipped and reported in the response

```json
{
  "results": [
    {"email": "active@example.com", "id": "uuid", "status": "queued"},
    {"email": "suppressed@example.com", "status": "suppressed"}
  ]
}
```
