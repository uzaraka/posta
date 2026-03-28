---
sidebar_position: 5
title: Scheduled Jobs
description: View and manage scheduled background jobs
---

# Scheduled Jobs

Posta runs scheduled background jobs for maintenance tasks.

## List Jobs

```
GET /api/v1/admin/jobs
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "name": "retention_cleanup",
      "schedule": "0 2 * * *",
      "last_execution": "2026-01-15T02:00:00Z",
      "next_execution": "2026-01-16T02:00:00Z",
      "status": "success"
    },
    {
      "name": "daily_report",
      "schedule": "0 8 * * *",
      "last_execution": "2026-01-15T08:00:00Z",
      "next_execution": "2026-01-16T08:00:00Z",
      "status": "success"
    }
  ]
}
```

## Built-in Jobs

| Job | Schedule | Description |
|-----|----------|-------------|
| `retention_cleanup` | Daily at 2:00 AM | Deletes expired email logs, webhook deliveries, and audit entries based on retention settings |
| `daily_report` | Daily at 8:00 AM | Sends daily email summary reports to users who opted in |
