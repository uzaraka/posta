# Workspaces

Workspaces provide multi-tenant isolation within Posta. They work like GitHub Organizations — users have a **personal space** by default, and can optionally create **workspaces** to share resources with team members.

## Concepts

### Personal Space

Every user has a personal space where their resources (templates, SMTP servers, domains, contacts, API keys, etc.) live by default. No workspace is required — the platform works exactly as before for single users.

### Workspaces

A workspace is an isolated environment where team members can collaborate. Resources created within a workspace are only visible to workspace members. Each workspace has:

- A unique **name** and **slug**
- An **owner** (the creator)
- **Members** with assigned roles
- Isolated resources (templates, SMTP servers, domains, API keys, contacts, emails, webhooks, etc.)

### Roles

| Role | View | Create/Edit | Manage Members | Delete Workspace |
|------|------|-------------|----------------|------------------|
| **Owner** | Yes | Yes | Yes | Yes |
| **Admin** | Yes | Yes | Yes | No |
| **Editor** | Yes | Yes | No | No |
| **Viewer** | Yes | No | No | No |

## Creating a Workspace

Navigate to **Workspaces** in the sidebar and click **Create Workspace**. Provide a name and slug (URL-friendly identifier).

## Switching Context

Use the **workspace switcher** in the sidebar to switch between your personal space and any workspace you belong to. When you switch:

- All resource views (templates, emails, contacts, etc.) show data from the selected context
- New resources you create are scoped to the selected context
- Analytics and dashboard stats reflect the selected context

## Inviting Members

Workspace owners and admins can invite users:

1. Go to **Workspaces** > select your workspace > **Manage**
2. Open the **Invitations** tab
3. Click **Invite** and enter the user's email and role

Invited users will see pending invitations on the **Workspaces** page and can accept or decline directly from the dashboard.

## Data Transfer

Transfer your personal resources into a workspace:

1. Go to the workspace's **Manage** page
2. Open the **Data Transfer** tab
3. Select which resource types to transfer (templates, contacts, SMTP servers, etc.)
4. Click **Transfer**

Transferred resources will no longer appear in your personal space — they become workspace-scoped.

## API Usage

### Workspace Context Header

To interact with workspace-scoped resources via the API, include the `X-Workspace-ID` header:

```bash
curl -X GET http://localhost:9000/api/v1/users/me/templates \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Workspace-ID: 1"
```

If the header is omitted, the request operates in the user's personal space.

### Workspace-Scoped API Keys

API keys created within a workspace context are automatically scoped to that workspace. Requests using a workspace-scoped API key will operate within the workspace without needing the `X-Workspace-ID` header.

### Workspace Management Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/workspaces` | Create workspace |
| `GET` | `/api/v1/workspaces` | List user's workspaces |
| `GET` | `/api/v1/workspaces/current` | Get current workspace |
| `PUT` | `/api/v1/workspaces/current` | Update workspace |
| `DELETE` | `/api/v1/workspaces/current` | Delete workspace |
| `GET` | `/api/v1/workspaces/current/members` | List members |
| `PUT` | `/api/v1/workspaces/current/members/{id}` | Update member role |
| `DELETE` | `/api/v1/workspaces/current/members/{id}` | Remove member |
| `POST` | `/api/v1/workspaces/current/invitations` | Invite member |
| `GET` | `/api/v1/workspaces/current/invitations` | List invitations |
| `DELETE` | `/api/v1/workspaces/current/invitations/{id}` | Cancel invitation |
| `POST` | `/api/v1/workspaces/current/transfer` | Transfer personal data |
| `GET` | `/api/v1/invitations` | List my invitations |
| `POST` | `/api/v1/invitations/{id}/accept` | Accept invitation |
| `POST` | `/api/v1/invitations/{id}/decline` | Decline invitation |

All workspace management endpoints require the `X-Workspace-ID` header (except create, list, and invitation actions).
