---
sidebar_position: 2
title: Email Analytics
description: Detailed email analytics and trends
---

# Email Analytics

Posta provides detailed analytics on email delivery performance with date range filtering.

## Basic Analytics

```
GET /api/v1/users/me/analytics?start_date=2026-01-01&end_date=2026-01-31&page=1&size=30
```

Response:

```json
{
  "success": true,
  "data": [
    {"date": "2026-01-01", "sent": 150, "failed": 3, "bounced": 1},
    {"date": "2026-01-02", "sent": 200, "failed": 5, "bounced": 2}
  ],
  "pageable": {
    "current_page": 1,
    "size": 30,
    "total_pages": 1,
    "total_elements": 31
  }
}
```

## Advanced Dashboard Analytics

```
GET /api/v1/users/me/analytics/dashboard?start_date=2026-01-01&end_date=2026-01-31
```

Provides:

- **Delivery rate trend** — Daily delivery success rate over time
- **Bounce rate graph** — Daily bounce rate over time
- **Latency percentiles** — p50, p95, p99 email delivery latency

## Admin Platform Analytics

Administrators can view platform-wide analytics:

```
GET /api/v1/admin/analytics?start_date=2026-01-01&end_date=2026-01-31
```

```
GET /api/v1/admin/analytics/dashboard?start_date=2026-01-01&end_date=2026-01-31
```

These endpoints aggregate data across all users.
