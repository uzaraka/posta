# Posta Dashboard

The admin dashboard for [Posta](https://github.com/jonasfroeller/posta), built with Vue 3, TypeScript, and Vite.

## Setup

```bash
npm install
npm run dev
```

The dev server runs on port 3000 and proxies `/api/v1` requests to the backend at `localhost:9000`.

## Scripts

| Command | Description |
|---|---|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run preview` | Preview production build |

## Tech Stack

- **Vue 3** with `<script setup>` SFCs
- **TypeScript**
- **Vite** for bundling
- **Pinia** for state management
- **Vue Router** with auth guards
- **Axios** for API requests

## Project Structure

```
src/
  api/           API client modules (one per feature)
  components/    Reusable components (ConfirmDialog, Pagination)
  composables/   Composition utilities
  layouts/       Dashboard layout with sidebar navigation
  router/        Route definitions and auth guards
  stores/        Pinia stores (auth, notification, theme)
  views/         Page components organized by feature
```

## Features

- JWT authentication with 2FA support
- Email management and delivery tracking
- Template editor with versioning and localization
- SMTP server configuration
- Domain verification (SPF, DKIM, DMARC)
- Contact and list management
- API key management
- Webhook configuration and delivery logs
- Analytics dashboard
- Admin panel (users, jobs, metrics, events)
- Light/dark/system theme support
