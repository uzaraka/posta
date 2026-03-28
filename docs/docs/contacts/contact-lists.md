---
sidebar_position: 2
title: Contact Lists
description: Organize contacts into lists
---

# Contact Lists

Organize your contacts into reusable mailing lists.

## Create a List

```
POST /api/v1/users/me/contact-lists
```

```json
{
  "name": "Newsletter Subscribers",
  "description": "Users who opted in to the monthly newsletter"
}
```

## List All Contact Lists

```
GET /api/v1/users/me/contact-lists?page=1&size=20
```

Response includes `member_count` for each list.

## Update a List

```
PUT /api/v1/users/me/contact-lists/{id}
```

```json
{
  "name": "Weekly Newsletter",
  "description": "Updated description"
}
```

## Delete a List

```
DELETE /api/v1/users/me/contact-lists/{id}
```

## Managing Members

### Add a Member

```
POST /api/v1/users/me/contact-lists/{listId}/members
```

```json
{
  "email": "user@example.com"
}
```

Returns `409 Conflict` if the email is already a member.

### Remove a Member

```
DELETE /api/v1/users/me/contact-lists/{listId}/members
```

```json
{
  "email": "user@example.com"
}
```

### List Members

```
GET /api/v1/users/me/contact-lists/{listId}/members?page=1&size=20
```
