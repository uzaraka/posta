---
sidebar_position: 1
title: Dashboard Stats
description: Dashboard statistics and metrics
---

# Dashboard Statistics

Get an overview of your email sending activity.

## User Dashboard Stats

```
GET /api/v1/users/me/dashboard/stats
```

Response:

```json
{
  "success": true,
  "data": {
    "total_emails": 15420,
    "sent": 14890,
    "failed": 230,
    "bounce_rate": 1.5,
    "delivery_rate": 96.6
  }
}
```

## Email Listing

View your sent emails with pagination:

```
GET /api/v1/users/me/emails?page=1&size=20
```

## Email Details

Get full details for a specific email:

```
GET /api/v1/users/me/emails/{emailId}
```

:::note
Email content may be redacted depending on privacy settings configured by the administrator.
:::
