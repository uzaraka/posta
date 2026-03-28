---
sidebar_position: 2
title: Installation
description: Deploy Posta with Docker Compose or from source
---

# Installation

Posta can be deployed using Docker Compose (recommended) or built from source.

## Docker Compose (Recommended)

Create a `compose.yml` file:

```yaml
services:
  posta:
    image: goposta/posta:latest
    ports:
      - "9000:9000"
    environment:
      POSTA_DB_HOST: postgres
      POSTA_DB_USER: posta
      POSTA_DB_PASSWORD: posta
      POSTA_DB_NAME: posta
      POSTA_DB_PORT: 5432
      POSTA_REDIS_ADDR: redis:6379
      POSTA_JWT_SECRET: change-me-in-production
      POSTA_ADMIN_EMAIL: admin@example.com
      POSTA_ADMIN_PASSWORD: admin1234
      POSTA_EMBEDDED_WORKER: "true"
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: posta
      POSTGRES_PASSWORD: posta
      POSTGRES_DB: posta
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redisdata:/data

volumes:
  pgdata:
  redisdata:
```

Start the services:

```bash
docker compose up -d
```

Posta will be available at `http://localhost:9000`.

## Build from Source

### Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Redis 7+
- Node.js 18+ (for building the dashboard)

### Steps

```bash
# Clone the repository
git clone https://github.com/goposta/posta.git
cd posta

# Build the binary
make build

# Run the server
./bin/posta server
```

## Worker Deployment

Posta processes emails asynchronously using background workers. You have two options:

### Embedded Worker (Simple)

Set `POSTA_EMBEDDED_WORKER=true` to run the worker within the API server process. This is the simplest setup for small deployments.

### Standalone Worker (Scalable)

For larger deployments, run the worker as a separate process:

```bash
# API server
./bin/posta server

# Worker (separate process/container)
./bin/posta worker
```

You can run multiple worker instances for horizontal scaling. All workers share the same Redis queue.

## Health Checks

Once running, verify the deployment:

```bash
# Liveness probe
curl http://localhost:9000/api/v1/healthz

# Readiness probe (checks DB + Redis)
curl http://localhost:9000/api/v1/readyz
```
