---
sidebar_position: 2
title: Creating Templates
description: Create and manage email templates
---

# Creating Templates

## Create a Template

```
POST /api/v1/users/me/templates
```

**Authentication:** JWT token

```bash
curl -X POST http://localhost:9000/api/v1/users/me/templates \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "welcome",
    "description": "Welcome email for new users",
    "html": "<h1>Welcome, {{name}}!</h1><p>Thanks for joining.</p>",
    "text": "Welcome, {{name}}! Thanks for joining.",
    "sample_data": {
      "name": "Alice"
    }
  }'
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "welcome",
    "description": "Welcome email for new users",
    "active_version_id": "version-uuid",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

:::info
Template names must be unique. A `409 Conflict` is returned if the name already exists.
:::

## List Templates

```
GET /api/v1/users/me/templates?page=1&size=20
```

## Update a Template

```
PUT /api/v1/users/me/templates/{id}
```

```json
{
  "name": "welcome-v2",
  "description": "Updated welcome email",
  "default_language": "en"
}
```

## Delete a Template

```
DELETE /api/v1/users/me/templates/{id}
```

Returns `204 No Content`.

## Using Templates for Sending

Once created, reference templates by name in the send API:

```bash
curl -X POST http://localhost:9000/api/v1/emails/send-template \
  -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "template": "welcome",
    "to": ["user@example.com"],
    "template_data": {
      "name": "Bob"
    }
  }'
```
