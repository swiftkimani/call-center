# MCP in This Project

MCP stands for Model Context Protocol. It lets an AI tool talk to external systems through structured servers instead of ad hoc shell scripts or prompt copy-paste.

## Why you would use MCP here

In a call center platform, MCP is useful when you want AI-assisted workflows to safely access:

- PostgreSQL schema or read-only analytics queries
- CRM data
- ticketing systems
- internal documentation
- OpenAPI specs
- deployment metadata

Instead of hardcoding all of that into prompts, you expose it through MCP servers and let the client request structured resources or tools.

## Good MCP use cases for this repo

- Read the database schema and explain a table without giving the model raw production credentials.
- Surface queue health, campaign throughput, or call volume from a controlled internal tool.
- Let AI inspect architecture docs, runbooks, and API contracts as first-class resources.
- Connect the frontend or ops assistants to a CRM or helpdesk system through a stable tool boundary.

## Practical MCP model

Use three categories of servers:

1. `docs-mcp`
   Serves markdown docs, runbooks, ADRs, and API contracts from this repository.

2. `data-mcp`
   Exposes read-only Postgres or analytics queries such as campaign counts, queue occupancy, or disposition summaries.

3. `ops-mcp`
   Exposes safe operational actions like checking deployment status, queue lag, or object storage health.

## Example mental model

```text
AI Client
  -> docs-mcp      # architecture and runbooks
  -> data-mcp      # safe reads from Postgres and analytics
  -> ops-mcp       # status, logs, metrics, non-destructive operations
```

## How to start

1. Pick one narrow server first.
2. Make it read-only.
3. Expose a few high-value resources or tools.
4. Add authentication and audit logs before allowing write actions.

## Example MCP resource ideas

- `docs://architecture/scalable-structure`
- `docs://runbooks/telephony-webhooks`
- `db://tables/calls`
- `db://reports/daily-call-volume`
- `ops://deployments/api/status`

## Important rule

Do not use MCP as a shortcut around normal application architecture. MCP is for tool-assisted context and automation. Your frontend should still call your backend API. Your backend should still own business rules, auth, and persistence.
