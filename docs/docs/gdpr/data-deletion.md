---
sidebar_position: 2
title: Data Deletion
description: Delete contact and email data for GDPR compliance
---

# Data Deletion

Posta provides endpoints for GDPR-compliant data deletion.

## Delete Contact Data

Remove all data associated with a specific email address, or all contacts:

```
POST /api/v1/users/me/gdpr/delete-contacts
```

### Delete a Specific Contact

```json
{
  "email": "user@example.com"
}
```

### Delete All Contacts

```json
{
  "email": null
}
```

Response:

```json
{
  "success": true,
  "data": {
    "deleted_count": 1
  }
}
```

This removes the contact from:
- The contacts list
- All contact lists
- The suppression list

## Delete Email Logs

Remove email logs older than a specified number of days:

```
POST /api/v1/users/me/gdpr/delete-email-logs
```

### Delete Logs Older Than 30 Days

```json
{
  "older_than_days": 30
}
```

### Delete All Email Logs

```json
{
  "older_than_days": 0
}
```

Response:

```json
{
  "success": true,
  "data": {
    "deleted_count": 1500
  }
}
```

## Automatic Retention

Administrators can configure automatic data retention via [Platform Settings](/docs/admin/platform-settings):

- `email_retention_days` — Auto-delete email logs after N days
- `webhook_retention_days` — Auto-delete webhook delivery history after N days
- `audit_log_retention_days` — Auto-delete audit log entries after N days

The retention cleanup job runs daily at 2:00 AM.
