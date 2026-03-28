---
sidebar_position: 4
title: Localization
description: Multi-language template support
---

# Template Localization

Add localized versions of templates for different languages. Localizations are attached to specific template versions.

## Managing Languages

Before adding localizations, create the languages you need:

### Create a Language

```
POST /api/v1/users/me/languages
```

```json
{
  "code": "fr",
  "name": "French"
}
```

### List Languages

```
GET /api/v1/users/me/languages
```

## Adding Localizations

### Add a Localization to a Version

```
POST /api/v1/users/me/templates/{templateId}/versions/{versionId}/localizations
```

```json
{
  "language_id": "language-uuid",
  "html": "<h1>Bienvenue, {{name}} !</h1><p>Merci de nous avoir rejoints.</p>",
  "text": "Bienvenue, {{name}} ! Merci de nous avoir rejoints."
}
```

Response (`201`):

```json
{
  "success": true,
  "data": {
    "id": "localization-uuid",
    "language_id": "language-uuid",
    "html": "...",
    "text": "..."
  }
}
```

:::info
A `409 Conflict` is returned if the language already exists for this version.
:::

### List Localizations

```
GET /api/v1/users/me/templates/{templateId}/versions/{versionId}/localizations
```

### Update a Localization

```
PUT /api/v1/users/me/localizations/{localizationId}
```

```json
{
  "html": "<h1>Updated French content</h1>",
  "text": "Updated French content"
}
```

### Delete a Localization

```
DELETE /api/v1/users/me/localizations/{localizationId}
```

## Sending with Language

Specify the language when sending:

```bash
curl -X POST http://localhost:9000/api/v1/emails/send-template \
  -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "template": "welcome",
    "to": ["user@example.fr"],
    "language": "fr",
    "template_data": {"name": "Marie"}
  }'
```

### Language Resolution

1. Look for a localization matching the requested language on the active version
2. If not found, fall back to the template's default language
3. If no default is set, use the base template content

### Per-Recipient Language in Batch

```json
{
  "template": "newsletter",
  "recipients": [
    {"email": "bob@example.com", "language": "en", "template_data": {"name": "Bob"}},
    {"email": "marie@example.fr", "language": "fr", "template_data": {"name": "Marie"}},
    {"email": "hans@example.de", "language": "de", "template_data": {"name": "Hans"}}
  ]
}
```

## Preview a Localized Version

```
POST /api/v1/users/me/templates/{templateId}/versions/{versionId}/preview
```

```json
{
  "language_id": "language-uuid",
  "template_data": {"name": "Marie"}
}
```
