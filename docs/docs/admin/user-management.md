---
sidebar_position: 1
title: User Management
description: Manage users from the admin panel
---

# User Management

Administrators can create, update, and delete user accounts.

## Create a User

```
POST /api/v1/admin/users
```

```json
{
  "name": "New User",
  "email": "user@example.com",
  "password": "secure-password",
  "role": "user"
}
```

| Role | Description |
|------|-------------|
| `user` | Standard user — can send emails, manage own resources |
| `admin` | Administrator — full platform access |

Returns `409 Conflict` if the email already exists.

## List Users

```
GET /api/v1/admin/users?page=1&size=20
```

## Update a User

```
PUT /api/v1/admin/users/{userId}
```

```json
{
  "name": "Updated Name",
  "email": "new-email@example.com",
  "role": "admin"
}
```

## User Metrics

Get sending statistics for a specific user:

```
GET /api/v1/admin/users/{userId}/metrics
```

```json
{
  "success": true,
  "data": {
    "total_emails": 5000,
    "sent": 4850,
    "failed": 100,
    "bounce_rate": 1.0,
    "api_keys_count": 3
  }
}
```

## Delete a User

```
DELETE /api/v1/admin/users/{userId}
```

## Disable 2FA for a User

If a user loses access to their authenticator:

```
DELETE /api/v1/admin/users/{userId}/2fa
```

## Manage API Keys

View and revoke API keys across all users:

```
GET /api/v1/admin/api-keys?page=1&size=20
DELETE /api/v1/admin/api-keys/{keyId}
```

## View All Emails

```
GET /api/v1/admin/emails?page=1&size=20
```

Lists all emails across all users with pagination.
