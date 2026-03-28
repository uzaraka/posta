---
sidebar_position: 2
title: Domain Verification
description: Verify sending domains with SPF, DKIM, and DMARC
---

# Domain Verification

Verify your sending domains to improve deliverability and prevent spoofing. Posta checks SPF, DKIM, and DMARC DNS records.

## Register a Domain

```
POST /api/v1/users/me/domains
```

```json
{
  "domain": "yourdomain.com"
}
```

Response:

```json
{
  "success": true,
  "data": {
    "id": "domain-uuid",
    "domain": "yourdomain.com",
    "ownership_verified": false,
    "verification_token": "posta-verify=abc123def456",
    "spf_record": "v=spf1 include:yourdomain.com ~all",
    "dkim_record": "...",
    "dmarc_record": "v=DMARC1; p=quarantine; rua=mailto:dmarc@yourdomain.com"
  }
}
```

## Add DNS Records

Add the following DNS records to your domain:

### 1. Ownership Verification (TXT)

| Type | Host | Value |
|------|------|-------|
| TXT | `@` or `yourdomain.com` | `posta-verify=abc123def456` |

### 2. SPF Record (TXT)

| Type | Host | Value |
|------|------|-------|
| TXT | `@` | `v=spf1 include:yourdomain.com ~all` |

### 3. DKIM Record (TXT)

| Type | Host | Value |
|------|------|-------|
| TXT | `posta._domainkey` | *(provided in response)* |

### 4. DMARC Record (TXT)

| Type | Host | Value |
|------|------|-------|
| TXT | `_dmarc` | `v=DMARC1; p=quarantine; rua=mailto:dmarc@yourdomain.com` |

## Verify DNS Records

After adding DNS records, trigger verification:

```
POST /api/v1/users/me/domains/{id}/verify
```

Response:

```json
{
  "success": true,
  "data": {
    "ownership_verified": true,
    "spf_verified": true,
    "dkim_verified": true,
    "dmarc_verified": false
  }
}
```

:::tip
DNS propagation can take up to 48 hours. If verification fails, wait and try again.
:::

## Domain Enforcement

When `require_verified_domain` is enabled in user settings, Posta will reject emails from unverified domains. This adds an extra layer of security to prevent unauthorized sending.

## List Domains

```
GET /api/v1/users/me/domains?page=1&size=20
```

## Delete a Domain

```
DELETE /api/v1/users/me/domains/{id}
```
