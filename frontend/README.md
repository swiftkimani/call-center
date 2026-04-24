# Frontend

This is the Next.js 16 frontend for the Go call center backend.

## Stack

- Next.js 16 App Router
- React 19
- Tailwind CSS 4
- shadcn/ui-style local components
- TanStack Query provider for future live views

## Run

1. Copy `.env.local.example` to `.env.local`.
2. Set `CALLCENTER_API_BASE_URL` to the Go API URL.
3. Optionally set `CALLCENTER_API_TOKEN` to call protected endpoints like `/api/v1/me`.
4. Install dependencies with `npm install`.
5. Start the app with `npm run dev`.

## Structure

```text
frontend/
├── src/app                    # App Router entrypoints
├── src/components             # shared layout, providers, and UI primitives
├── src/features               # vertical feature slices
├── src/lib                    # env, HTTP client, shared utilities
└── components.json            # shadcn/ui registry config
```
