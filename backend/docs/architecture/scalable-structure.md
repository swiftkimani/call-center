# Scalable Frontend and Backend Structure

This repository already has a strong backend foundation. The right move is not a full rewrite. The right move is to preserve the current domain split and evolve it into clearer deployment and ownership boundaries.

## Recommended top-level shape

```text
.
├── cmd/
│   ├── api/                   # stateless HTTP API process
│   └── worker/                # background job process
├── internal/
│   ├── platform/              # horizontal concerns shared by all modules
│   │   ├── auth/
│   │   ├── broker/
│   │   ├── config/
│   │   ├── db/
│   │   ├── observability/
│   │   ├── redis/
│   │   └── storage/
│   ├── modules/               # vertical business slices
│   │   ├── agents/
│   │   ├── calls/
│   │   ├── campaigns/
│   │   ├── customers/
│   │   ├── queues/
│   │   ├── reports/
│   │   └── telephony/
│   └── transport/             # HTTP, WS, webhooks, DTO mapping
├── frontend/                  # independently deployable Next.js 16 app
├── sql/                       # SQL-first contracts
├── migrations/                # schema evolution
└── docs/                      # architecture, runbooks, MCP notes
```

## What this means for scaling

### Vertical scalability

Vertical scaling means each business capability can grow in codebase size and team ownership without turning into a monolith of shared files.

- Keep each domain in its own module: `calls`, `campaigns`, `customers`, `agents`, `queues`.
- Put transport-only code in `internal/transport/http` and `internal/transport/ws`, not in domain services.
- Keep DB queries close to the module they serve, even if `sqlc` generation stays centralized.
- Treat the frontend the same way: feature folders own screens, server loaders, hooks, and local components.

### Horizontal scalability

Horizontal scaling means runtime instances can multiply without coordination bugs.

- The API should stay stateless except for cache and broker interactions.
- Realtime session state belongs in Redis, not in a single API instance.
- Slow work belongs in workers and queues, not in request handlers.
- Recording processing, campaign dialing, notification fan-out, and analytics rollups should each be able to move to their own worker queue.
- Frontend deployments should not depend on backend process memory; they only depend on the API contract.

## Frontend structure

```text
frontend/src/
├── app/                       # routes, layouts, metadata
├── components/
│   ├── layout/                # app chrome
│   ├── providers/             # query provider, theme provider
│   └── ui/                    # shadcn-style shared primitives
├── features/
│   ├── dashboard/
│   │   ├── components/
│   │   ├── server/
│   │   └── types.ts
│   ├── agents/
│   ├── calls/
│   ├── campaigns/
│   └── supervisor/
└── lib/
    ├── env.ts
    ├── http/
    └── utils.ts
```

Rules:

- Shared UI primitives go in `components/ui`.
- Route-specific composition stays in `features/<feature>/components`.
- Server fetching logic stays in `features/<feature>/server`.
- Browser-side caches and websockets should stay feature-local until they become truly shared.

## Backend evolution path

The current codebase already contains most of the right primitives. The next cleanups should be incremental:

1. Move cross-cutting libraries under `internal/platform`.
2. Move domain packages under `internal/modules`.
3. Keep `cmd/api` and `cmd/worker` thin composition roots only.
4. Add contract definitions under `api/` or `pkg/contracts/` when frontend and backend teams need generated SDKs.
5. Split workers by responsibility once queue traffic justifies it.

## Suggested deployment units

- `frontend`: Next.js app, separately deployable.
- `api`: Go HTTP service, horizontally scalable.
- `worker-campaigns`: outbound dialing and campaign orchestration.
- `worker-recordings`: recording ingestion and enrichment.
- `worker-notifications`: SMS, email, WhatsApp, CRM sync.
- `postgres`, `redis`, `rabbitmq`, `minio`: infrastructure tier.

This gives you a path from one team and one VPS to multiple teams and multiple instances without changing the mental model.
