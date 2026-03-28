---
sidebar_position: 1
title: Data Export & Import
description: Export and import user data
---

# Data Export & Import

Posta supports full data export and import for GDPR compliance and environment migration.

## Export All User Data

```
GET /api/v1/users/me/data/export
```

Returns all user data as JSON:

```json
{
  "success": true,
  "data": {
    "templates": [...],
    "stylesheets": [...],
    "languages": [...],
    "contacts": [...],
    "contact_lists": [...],
    "webhooks": [...],
    "suppressions": [...],
    "settings": {...}
  }
}
```

## Import User Data

```
POST /api/v1/users/me/data/import
```

Send the exported JSON as the request body. Existing items are skipped (not overwritten).

```json
{
  "success": true,
  "data": {
    "imported_count": 42
  }
}
```

## Use Cases

- **GDPR data portability** — Users can export all their data
- **Environment migration** — Move settings between staging and production
- **Backup** — Periodic data exports for disaster recovery
