# CALL CENTER PLATFORM
## Implementation Blueprint & Delivery Plan

## Repository Additions

- Frontend app: [`frontend/`](./frontend)
- Frontend env example: [`frontend/.env.local.example`](./frontend/.env.local.example)
- Scale-oriented structure guide: [`docs/architecture/scalable-structure.md`](./docs/architecture/scalable-structure.md)
- MCP guide: [`docs/mcp.md`](./docs/mcp.md)

A complete engineering and project-management reference for building, launching, and operating a cloud-native call center platform in 28 days.

**Stack:** Go · Next.js · PostgreSQL

**Supporting services:** Redis · RabbitMQ · Africa's Talking · MinIO · Prometheus · Grafana

Document version 1.0 · April 2026

**Audience:** engineering, product, QA, operations, and delivery partners

---

## How to read this document

This blueprint is written to be usable by every role that touches the project. Non-engineers can read Sections 1, 2, 10, 11, and 16 to understand the vision, architecture at a glance, project plan, and risks. Engineers should read front to back; infrastructure and DevOps teams will spend most of their time in Sections 3, 7, 9, and 13. Designers and frontend engineers should focus on Sections 3, 5, and 8.

Every technical decision is paired with the reasoning behind it. You will not find a bare list of tools — you will find why each tool was chosen and what it is responsible for.

### Document conventions

- Blue callout boxes highlight foundational concepts worth memorizing.
- Yellow callout boxes flag risks, common pitfalls, and production gotchas.
- Green callout boxes capture tips and accelerators for a first-time builder.
- Monospaced blocks are shell commands, API payloads, or code. They can be pasted directly into a terminal or IDE.
- Diagrams are the source of truth when text and pictures disagree. If a diagram is updated, the prose near it should be updated in the same commit.

### The guiding principle

Ship a boring, reliable MVP in four weeks. Use mature tools, avoid exotic dependencies, and prefer managed services wherever the cost is justified. Every week of engineering saved is a week earlier to paying customers and real-world feedback.

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Product Vision and Objectives](#2-product-vision-and-objectives)
3. [System Architecture](#3-system-architecture)
4. [Technology Stack — Every Choice Explained](#4-technology-stack--every-choice-explained)
5. [Data Layer — Schema, State, and Events](#5-data-layer--schema-state-and-events)
6. [API Surface](#6-api-surface)
7. [Data Flow and Workflow Diagrams](#7-data-flow-and-workflow-diagrams)
8. [Frontend Design and UX](#8-frontend-design-and-ux)
9. [Deployment Architecture](#9-deployment-architecture)
10. [Security and Compliance](#10-security-and-compliance)
11. [Observability and Runbooks](#11-observability-and-runbooks)
12. [Project Management and Delivery Approach](#12-project-management-and-delivery-approach)
13. [Team Structure, Roles, and Communication](#13-team-structure-roles-and-communication)
14. [Testing Strategy](#14-testing-strategy)
15. [Go-Live Checklist](#15-go-live-checklist)
16. [Risk Register](#16-risk-register)
17. [Post-Launch Roadmap](#17-post-launch-roadmap)
18. [Appendices](#18-appendices)

---

## 1. Executive Summary

We are building a modern, cloud-native call center platform that lets a business accept inbound customer calls, run outbound campaigns, manage agents and supervisors in real time, and generate the reporting required to run a commercial operation. The system is designed to go from zero to a live pilot with a small agent team in four weeks, and to scale to hundreds of concurrent agents without architectural changes.

### 1.1 What this product does

- Receive customer calls to a business phone number and route them to the most appropriate available agent in under three seconds.
- Place outbound calls through a click-to-dial interface or automated campaign, respecting do-not-call rules and working-hour constraints.
- Give agents a single browser-based workspace that rings, shows the customer's history, captures notes, and records every call.
- Give supervisors a live view of every queue and agent, with the ability to whisper, barge, coach, and pull reports on demand.
- Integrate naturally with WhatsApp, email, and SMS so customer conversations are not trapped in a single channel.

### 1.2 Why this stack

Go, Next.js, and PostgreSQL are the strongest pragmatic choices for a product with this risk and timeline profile.

- **Go** is purpose-built for highly concurrent network services. A single process handles thousands of WebSocket connections with predictable memory usage and sub-millisecond scheduling. Call routing is exactly the workload Go was designed for.
- **Next.js** gives agents a rich, interactive browser app without forcing us to stitch together a build toolchain. Server-side rendering accelerates first paint; React hooks make live UIs tractable.
- **PostgreSQL** is battle-tested for both transactional workloads (every call is a financial event) and analytical queries (supervisors want historical reports). Its JSONB columns give us NoSQL-style flexibility when we need it.
- **Everything else** in this blueprint — Redis, RabbitMQ, Africa's Talking, MinIO, Prometheus — is chosen because it is the mature, well-documented default for its role. No single dependency is exotic, which means any engineer can be productive in days.

### 1.3 Headline commitments

| Commitment | Target |
|---|---|
| Time to first production call | 28 days from day zero |
| Concurrent agents supported by MVP | 50 agents on a single VPS |
| Architectural ceiling without rewrite | 500+ agents, 1,000 concurrent calls |
| Call routing latency (webhook → ring) | under 500 ms (p95) |
| System availability target | 99.5% in month one, 99.9% after month three |
| Recording retention (default) | 90 days, encrypted, off-site backup |
| Onboarding time for a new agent | under 20 minutes |

### 1.4 What is explicitly out of scope for v1

Scope discipline is the single largest determinant of whether we hit the four-week deadline. The following are deliberately deferred:

- **Self-hosted SIP/media server** (Asterisk/FreeSWITCH). We use Africa's Talking or Twilio as the telephony provider. We revisit self-hosting when per-minute costs justify the operational overhead, typically at over 20,000 minutes per day.
- **AI-driven features** (call summarization, sentiment analysis, automatic disposition). The architecture leaves clear seams for these, but they are not in the MVP.
- **Mobile native apps.** The web app is mobile-responsive. Native apps are a post-launch investment.
- **Chatbot/IVR visual designer.** V1 ships with a handful of hard-coded IVR flows. A drag-and-drop builder is a v2 feature.
- **Multi-tenancy at the database level.** V1 serves a single business. If we later sell this as SaaS, we add a tenant column and row-level security.

---

## 2. Product Vision and Objectives

### 2.1 The one-sentence vision

> Every conversation between a business and its customer — inbound or outbound, voice or chat — happens in one workspace, is recorded, is measurable, and can be improved by the next sprint.

### 2.2 Business objectives (what "success" looks like)

1. Reduce customer wait time to answer by at least 40% compared to the baseline the business measures today.
2. Lift agent utilization (time on productive calls ÷ time logged in) to 70% or higher.
3. Provide a recording and transcript for 100% of completed calls, retrievable in under 10 seconds.
4. Let supervisors run a complete shift — live monitoring, coaching, and end-of-shift reporting — without leaving the platform.
5. Ship a commercial-grade pilot in under 30 days so real usage can inform v2 priorities.

### 2.3 Personas and their jobs-to-be-done

| Persona | Core job when using the system |
|---|---|
| Customer (caller) | Reach a human who can help as quickly as possible; do not repeat information. |
| Agent | Answer calls, capture accurate notes, move on to the next call with minimal clicks. |
| Supervisor / Team Lead | See who is doing what right now; coach people who are struggling; defend the team's KPIs. |
| Campaign Manager | Load a list of customers, schedule a campaign, monitor progress, and export the outcome. |
| Operations / IT | Add agents, reset passwords, rotate phone numbers, pull audit logs. |
| Business owner | See one weekly dashboard that tells them whether the call center is paying for itself. |

### 2.4 Non-functional requirements (the invisible work)

- **Reliability:** a dropped call is worse than a slow call. Every flow has a defined failure mode and a fallback.
- **Observability:** if something breaks at 2 a.m., the on-call engineer must see what, where, and when within 60 seconds of waking up.
- **Security:** call recordings contain personally identifiable information (PII) and sometimes payment data. Encryption at rest and in transit is mandatory, not optional.
- **Cost awareness:** telephony minutes dominate the monthly bill. The system logs per-call cost from day one so we can optimize.
- **Operability:** a second engineer should be able to run a deploy, a database restore, and an incident investigation using only the runbook in this document.

---

## 3. System Architecture

This section answers a simple question: if you draw a single picture that shows everything, what does it look like, and what does each piece do?

### 3.1 High-level architecture

*Figure 1. High-level architecture. Solid arrows are HTTP and WebSocket traffic; dashed orange arrows are the audio media path over WebRTC.*

Three kinds of traffic flow through the system and it helps to keep them mentally separate:

- **Signaling traffic** is the conversation between the telephony provider and our Go backend — who is calling whom, which agent should ring, when the call ends. This is plain HTTP (webhooks in, TwiML-style XML back out).
- **Media traffic** is the actual audio. It never touches our servers. It flows directly between the telephony provider and the agent's browser over WebRTC. This is why the platform scales so well: we only carry the small signaling messages, not the heavy audio streams.
- **UI traffic** is the WebSocket and REST traffic between the Next.js app and the Go backend — the things that make the agent's screen light up and the supervisor dashboard update in real time.

### 3.2 Architectural style: modular monolith, not microservices

We ship a single Go binary. It exposes REST endpoints, webhook endpoints, and WebSocket endpoints, and is internally organized into clean packages (routing, agents, queues, CRM, recordings, reporting). This is a deliberate choice.

- Microservices would cost us speed. Splitting into five services on day one means five repos, five deploy pipelines, five Dockerfiles, and a service-to-service communication protocol to design. We have four weeks.
- A modular monolith scales further than most people think. A single Go process on a modest VPS comfortably handles thousands of concurrent WebSocket connections and millions of HTTP requests per day.
- The package boundaries are the future service boundaries. When we do split — likely for the dialer engine or the recording pipeline — the internal interfaces become network APIs and nothing else has to change.

*Figure 2. Internal organization of the Go binary. The transport layer is the only part that knows about HTTP. Domain services hold the business logic. The infrastructure layer talks to databases and external providers.*

### 3.3 The three responsibilities of the Go backend

#### 3.3.1 Routing and state

The backend is the brain. It decides which agent gets which call, tracks who is busy, who is available, who is in wrap-up, and how long customers have been waiting. It holds this "now" state in Redis because reading from Redis takes microseconds and this state changes dozens of times per second under load.

#### 3.3.2 Real-time push to the UI

When something changes — a call arrives, an agent goes offline, a queue grows past its SLA — every connected supervisor dashboard needs to know within a second. We do this with WebSockets. Go's goroutines make maintaining tens of thousands of open WebSocket connections straightforward: each connection is a lightweight goroutine, not a heavyweight thread.

#### 3.3.3 Persistence and reporting

Everything that matters historically — every call, every disposition, every campaign outcome — goes to PostgreSQL. We do this asynchronously: the routing code writes a lightweight event to Redis or RabbitMQ first, and a background worker persists it. This keeps the hot path fast.

> **Why this split matters:** Redis is the answer to "what is happening right now?" and PostgreSQL is the answer to "what has ever happened?". Confusing these two — trying to route calls from a Postgres query, or trying to run monthly reports out of Redis — is the number one mistake first-time builders make. Keep them separate.

### 3.4 Why we do not build our own telephony

Africa's Talking (for Kenyan and pan-African traffic) or Twilio (for global traffic) bridges the public phone network to our software. They hand us a webhook every time a call arrives and expect XML instructions back. We make the intelligence; they carry the minutes.

Running our own Asterisk or FreeSWITCH server is a two-engineer-month project on its own. For a four-week MVP it is the wrong trade. The moment we are consistently paying more than roughly USD 2,500 per month in telephony fees and have an engineer who enjoys SIP, we revisit this decision.

---

## 4. Technology Stack — Every Choice Explained

Every row in the table below is a decision. The "Why" column is the decision record. If a team member ever wants to change one of these choices, they should write a new decision record with at least as much depth as what is here.

### 4.1 The full stack at a glance

| Layer | Tool | Version | Why this one |
|---|---|---|---|
| Backend language | Go | 1.22+ | Lightweight goroutines, strong std lib, single static binary, predictable performance. |
| HTTP router | chi | v5 | Idiomatic, no magic, plays well with net/http, easy middleware composition. |
| DB access | pgx + sqlc | pgx v5 | Fast Postgres driver + generated, type-safe queries from plain SQL. |
| DB migrations | golang-migrate | latest | Simple SQL-first migrations; same tool works locally and in CI. |
| Realtime | gorilla/websocket | v1.5 | Industry standard WebSocket lib; hub pattern is well documented. |
| Message broker | RabbitMQ | 3.13 | Durable queues, native retry, great Go client. Kafka is overkill for MVP. |
| In-memory state | Redis | 7 | Microsecond reads/writes, perfect for "who is online right now". |
| Frontend framework | Next.js | 15 (App Router) | SSR, React Server Components, mature ecosystem, great DX. |
| UI components | shadcn/ui + Tailwind | latest | Copy-paste components we control; no runtime dependency surprises. |
| Frontend state | Zustand + TanStack Query | latest | Zustand for client state; TanStack Query for server state. |
| Database | PostgreSQL | 16 | Mature, fast, supports JSONB and full-text search; one tool covers many needs. |
| Object storage | MinIO (self) or AWS S3 | latest | S3-compatible; store recordings with lifecycle policies. |
| Telephony (primary) | Africa's Talking | API v1 | Kenyan and pan-African coverage, transparent pricing, strong voice API. |
| Telephony (fallback/intl) | Twilio | API 2010 | Global coverage; use for international DIDs and as failover. |
| Messaging | WhatsApp Cloud API | v18+ | Direct Meta integration; consistent with existing OmniPOS integration path. |
| Email | Resend or SES | latest | Resend for transactional simplicity; SES for volume and cost. |
| Reverse proxy | NGINX | 1.27 | TLS termination, rate limiting, WebSocket upgrade support. |
| Container runtime | Docker + Compose | v2 | Predictable, reproducible, production-adequate for single-VPS stage. |
| CI/CD | GitHub Actions | latest | Free for the scale we need; good secrets management. |
| Observability — metrics | Prometheus + Grafana | latest | Ubiquitous, low-cost, rich alerting; Go exporter is built-in. |
| Observability — logs | Loki + Promtail | latest | Integrates with Grafana; cheap log storage. |
| Observability — traces | OpenTelemetry + Jaeger | latest | Distributed tracing for multi-hop call flows. |
| Error tracking | Sentry (self or cloud) | latest | Grouped exceptions, release tracking, minimal setup. |
| Secret management | Doppler or SOPS | latest | Never commit secrets; Doppler is easiest, SOPS if self-hosting. |
| Automation glue | n8n | 1.x | Already running on Contabo; reuse for non-critical workflows. |

### 4.2 How the pieces fit together

The best way to understand the stack is to follow a single call through it. A customer dials the business line. The telephony provider (Africa's Talking) receives the call and fires an HTTP POST to our Go backend's webhook endpoint. The backend queries Redis to find an available agent, reads the customer's profile from Postgres, pushes a WebSocket event to that agent's Next.js tab so the screen rings, and replies to Africa's Talking with XML telling it to bridge the call to the agent's browser over WebRTC. The audio never hits our servers. When the call ends, the provider pings us again with the duration and recording URL; we enqueue a RabbitMQ message so a background worker can download the recording to MinIO and write the final call record to Postgres.

That is the entire platform in one paragraph. Every component earns its place by handling one step of that sequence.

### 4.3 Provider choices specifically for Kenya

Because this platform is being built from Nairobi and the first users will likely be Kenyan, some defaults differ from what a San Francisco team would choose.

- **Africa's Talking over Twilio** as the primary telephony provider. It gives better per-minute rates for +254 numbers, simpler KYC, and faster number provisioning. Twilio stays in the architecture as a failover and for future international expansion.
- **M-Pesa integration readiness.** Although billing is out of scope for v1, the database schema reserves space for payment references so we can add M-Pesa Daraja callbacks without migration pain.
- **Data residency.** Primary Postgres lives on Contabo (EU) or a Kenyan provider. For call recordings, we default to a Kenyan or South African S3-compatible bucket to stay inside the Data Protection Act (Kenya, 2019) comfort zone.

---

## 5. Data Layer — Schema, State, and Events

### 5.1 The ER model

*Figure 3. Core entities. Each box is a PostgreSQL table. PK = primary key. FK = foreign key. Some v2 tables (notes, feedback, tags) are omitted for readability.*

### 5.2 What lives where

Deciding where each piece of data lives is the most important design work on this project. The rule: if you need it to make a routing decision in the next 200 milliseconds, it lives in Redis. If you need to answer it on a weekly report, it lives in Postgres. If it is larger than a few kilobytes and binary (a recording, a transcript file), it lives in object storage and Postgres holds only the pointer.

| Data | Where it lives | Why |
|---|---|---|
| Agent availability (online / busy / wrap-up) | Redis hash per agent + sets | Read on every routing decision; rewritten on every state change. |
| Active queues (who is waiting, how long) | Redis sorted set (ZSET) | ZADD gives O(log N) insertion and O(1) "who has waited longest". |
| Call in progress | Redis hash (call:{id}) | Frequent partial updates (mute on, hold off, transfer initiated). |
| Historical call records | Postgres table calls | Source of truth. Immutable once the call ends. |
| Call events (ring, answer, mute, end) | Postgres table call_events (JSONB) | Append-only log we can replay for audits. |
| Customer profiles and CRM data | Postgres customers + related tables | Transactional, queryable. |
| Recordings (audio file) | MinIO / S3 | Binary, large, not searchable as SQL. |
| Transcripts (text) | MinIO or Postgres (if small) | Text search can use Postgres FTS if kept in-DB. |
| Async jobs (send email, upload recording) | RabbitMQ | Durable queue with retries; isolates slow work from the hot path. |
| Pub/sub for UI updates | Redis Pub/Sub | Fan-out WebSocket events across multiple Go instances. |

### 5.3 Postgres schema excerpts

The full schema is checked in as migration files. The excerpts below show the shape of the two most load-bearing tables.

```sql
-- 001_users_and_agents.sql
CREATE TABLE users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email         TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  full_name     TEXT NOT NULL,
  role          TEXT NOT NULL CHECK (role IN ('admin','supervisor','agent')),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at    TIMESTAMPTZ
);

CREATE TABLE agents (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  extension       TEXT UNIQUE NOT NULL,
  skills          TEXT[]   NOT NULL DEFAULT '{}',
  status          TEXT     NOT NULL DEFAULT 'offline'
                           CHECK (status IN ('offline','available','busy','wrap_up','break')),
  max_concurrent  SMALLINT NOT NULL DEFAULT 1,
  team_id         UUID REFERENCES teams(id),
  last_seen_at    TIMESTAMPTZ
);
CREATE INDEX idx_agents_status        ON agents(status)
  WHERE status IN ('available','busy');
CREATE INDEX idx_agents_skills_gin    ON agents USING GIN (skills);
```

```sql
-- 003_calls.sql
CREATE TABLE calls (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  provider_sid   TEXT UNIQUE NOT NULL,   -- external call id from Africa's Talking/Twilio
  customer_id    UUID REFERENCES customers(id),
  agent_id       UUID REFERENCES agents(id),
  queue_id       UUID REFERENCES queues(id),
  direction      TEXT NOT NULL CHECK (direction IN ('inbound','outbound')),
  status         TEXT NOT NULL CHECK (status IN
                   ('queued','ringing','in_progress','completed','failed','no_answer','abandoned')),
  started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  answered_at    TIMESTAMPTZ,
  ended_at       TIMESTAMPTZ,
  wait_seconds   INT,
  talk_seconds   INT,
  recording_url  TEXT,
  cost_cents     INT,
  from_number    TEXT NOT NULL,
  to_number      TEXT NOT NULL
);
CREATE INDEX idx_calls_customer_started ON calls(customer_id, started_at DESC);
CREATE INDEX idx_calls_agent_started    ON calls(agent_id, started_at DESC);
CREATE INDEX idx_calls_started_at_brin  ON calls USING BRIN (started_at);
```

> **A note on BRIN indexes:** Call volume grows indefinitely. A BRIN index on `started_at` gives us fast date-range queries (for reports) while being dozens of times smaller than a B-tree index. Use BRIN for any append-only time-series column.

### 5.4 Redis key design

Keys are namespaced consistently so it is always obvious who owns each key. Every write has a TTL where it makes sense — stale agent state is worse than no state.

```bash
# Agent presence
HSET agent:{agent_id}  status "available"  last_seen "2026-04-22T09:15:00Z"
SADD agents:available  {agent_id}
EXPIRE agent:{agent_id}:heartbeat 30        # renewed every 10s by the agent's WS

# Queue state
ZADD queue:sales:waiting  {enqueue_ts}  {call_id}
HSET call:{call_id}  state "queued"  from "+254712345678"  customer_id "..."

# Pub/sub channel for live supervisor dashboards
PUBLISH ws:supervisor  '{"event":"queue_update","queue":"sales","waiting":7}'
```

### 5.5 Event sourcing — only where it earns its keep

We use append-only event logging for two specific areas: call events and agent state transitions. Every time an agent goes from available to busy, we write a row to `agent_events`. Every time a call changes state, we write a row to `call_events`. This lets us replay history for audits, reconstruct exactly what the supervisor saw at 14:32 on Tuesday, and compute KPIs without polling.

We deliberately do not event-source the rest of the system. The CRM and reporting tables are updated in place. Full-blown event sourcing across an entire domain is a six-month architectural commitment that this project does not need.

---

## 6. API Surface

The backend exposes three surfaces: a REST API for the UI and admin tools, WebSocket channels for real-time push, and webhook endpoints that our telephony provider calls. Each surface has a versioned URL prefix so we can evolve without breaking existing clients.

### 6.1 REST — key endpoints

| Method & Path | Purpose | Auth |
|---|---|---|
| POST /api/v1/auth/login | Exchange email+password for a JWT | Public |
| POST /api/v1/auth/refresh | Rotate refresh token | Refresh token |
| GET  /api/v1/me | Current user + agent profile | Bearer JWT |
| POST /api/v1/agents/:id/status | Set own status (available / break / offline) | Agent (self) |
| POST /api/v1/calls/outbound | Initiate a click-to-dial call | Agent |
| POST /api/v1/calls/:id/hold | Put a live call on hold | Agent |
| POST /api/v1/calls/:id/transfer | Warm or cold transfer to another agent or queue | Agent |
| POST /api/v1/calls/:id/disposition | Save wrap-up notes + category | Agent |
| GET  /api/v1/calls | Paginated historical list (filters + search) | Supervisor |
| GET  /api/v1/calls/:id/recording | Pre-signed URL to play recording | Supervisor/Agent |
| GET  /api/v1/queues/:id/live | Live queue snapshot (waiting, SLA breach count) | Supervisor |
| POST /api/v1/supervisor/:call_id/whisper | Start whispering to an agent mid-call | Supervisor |
| POST /api/v1/supervisor/:call_id/barge | Full-duplex join into a live call | Supervisor |
| POST /api/v1/campaigns | Create a new outbound campaign | Campaign Manager |
| POST /api/v1/campaigns/:id/contacts | Upload contacts CSV for a campaign | Campaign Manager |
| GET  /api/v1/reports/daily | Daily KPI report (JSON or CSV) | Supervisor |
| POST /webhooks/voice/inbound | Africa's Talking posts here when a call rings in | HMAC signature |
| POST /webhooks/voice/status | Provider posts call status updates | HMAC signature |
| POST /webhooks/recording/ready | Provider posts when a recording is available | HMAC signature |

### 6.2 WebSocket channels

We run two WebSocket endpoints. Both use the same Gorilla-based hub under the hood but carry different event shapes to different audiences.

- `/ws/agent` — one connection per logged-in agent. Events the server pushes: `incoming_call`, `call_ended`, `whisper_started`, `status_changed`, `broadcast_message`.
- `/ws/supervisor` — one connection per logged-in supervisor. Events: `queue_update`, `agent_state_change`, `sla_breach`, `campaign_progress`, `kpi_tick`.

All event payloads follow a consistent envelope so the frontend can route them with a single switch statement:

```json
{
  "id":        "evt_01H8X...",
  "type":      "incoming_call",
  "timestamp": "2026-04-22T09:15:03.214Z",
  "data": {
    "call_id":      "c_01H8X...",
    "from_number":  "+254712345678",
    "customer":     { "id": "...", "name": "Jane Mwangi", "tags": ["vip"] },
    "queue":        "sales",
    "wait_seconds": 4
  }
}
```

### 6.3 Webhook security

Webhook endpoints are the only parts of our API that the internet can reach without a session. They must be defended aggressively.

- **HMAC signature verification** on every webhook. Africa's Talking and Twilio both support this. We reject any request with a missing, malformed, or replayed signature.
- **Idempotency keys.** Providers retry on timeout. We dedupe by `provider_sid` so the same event is never processed twice.
- **Strict IP allowlist** at the NGINX layer; only known provider ranges can reach `/webhooks/*`.
- **Timeouts.** Webhooks have 5 seconds to return. Anything slower gets a 200 immediately and the work goes onto RabbitMQ.

---

## 7. Data Flow and Workflow Diagrams

Architecture diagrams show structure. Flow diagrams show behaviour. A team member who understands the flows below can contribute to any part of the codebase.

### 7.1 Flow A — inbound call from ring to wrap-up

*Figure 4. A customer dials the business number. The backend selects an agent, rings their browser, and the telephony provider bridges the audio directly to the agent.*

Key observations:

- Only steps 1, 9, 10, and 11 involve audio. Everything else is signaling — small JSON payloads. This is what keeps our servers small.
- Redis is the deciding authority for "who answers". Postgres is only consulted for context (customer history, queue configuration).
- The WebSocket push in step 8 is what makes the agent's phone ring. If WebSockets are down, the agent never hears the call — this is why WS health is a tier-1 alert.

### 7.2 Flow B — outbound click-to-dial

*Figure 5. An agent clicks a customer's phone number in the CRM view. The backend checks consent and compliance before placing the call.*

Two compliance checks matter on every outbound call and must be enforced server-side (the frontend is not trusted for this):

1. Is the customer on the do-not-call list? If yes, refuse. Audit log the attempt.
2. Is the current time within the permitted dialing window for the customer's jurisdiction? In Kenya this is typically 08:00–20:00 local time. Reject outside that range.

### 7.3 Flow C — recording pipeline

Recordings are the single largest data-volume concern. A one-minute call is roughly 500 KB in Opus or G.711. A thousand calls a day is about 500 MB of new data. This must be handled asynchronously.

1. Telephony provider finishes recording and posts a webhook with a temporary URL.
2. Backend writes a RabbitMQ message: `{ call_id, provider_recording_url }`.
3. A Go worker consumes the message, downloads the audio, uploads to MinIO under `recordings/YYYY/MM/DD/{call_id}.opus`, and updates `calls.recording_url` to the stable internal path.
4. If transcription is enabled, the worker enqueues a second message to a `transcribe` queue. A later worker runs speech-to-text and writes the result to `recordings.transcript_url`.
5. Lifecycle policy on the MinIO bucket moves recordings older than 30 days to a cheaper storage tier and deletes after 90 days (configurable per tenant later).

### 7.4 Flow D — supervisor whisper and barge

Whisper (supervisor speaks to the agent only) and barge (supervisor joins the conversation fully) are critical coaching tools. They use the same bridging mechanism as a three-way conference with mute flags per participant. Supervisors with the `supervisor` role in JWT can trigger these via `/api/v1/supervisor/:call_id/whisper` or `/barge`. The backend re-issues TwiML that adds the supervisor as a conferenced participant. Muting the supervisor's outgoing leg gives "listen only"; muting the supervisor to the customer gives "whisper"; muting no one gives "barge".

> **Tip:** The fastest way to verify the whole flow works end-to-end is to run a call twice a day during development. This catches subtle regressions (a dropped WebSocket event, a wrong TwiML shape) far earlier than any unit test.

---

## 8. Frontend Design and UX

The agent workspace is where the product is won or lost. An agent is on the phone for six hours a day. Every extra click, every half-second of lag, and every confusing label is an attack on their productivity. The frontend has to be fast, simple, and self-explanatory.

### 8.1 The three primary screens

| Screen | Who uses it | What it must do in one view |
|---|---|---|
| Agent workspace | Agents | Softphone · incoming call card · customer history · notes · disposition buttons · status toggle. |
| Supervisor cockpit | Supervisors | Live queues with SLA colors · agent grid with real-time states · listen/whisper/barge actions · shift KPIs. |
| Admin & reports | Admins / BMs | Users and teams · queue configuration · campaign manager · historical reports and CSV export. |

### 8.2 Agent workspace layout

The agent workspace is a single page with three columns. The left column is a compact activity feed (recent calls, unread notes). The centre is dominated by the current call card: big caller name, small caller number, customer tags, and a prominent timer. The right column shows the customer's history — recent calls, open tickets, last agent to speak with them. A slim footer holds the softphone dialpad, mute, hold, transfer, and end buttons. No modal ever covers the current call card. The agent can always see who is on the line.

### 8.3 Real-time updates without glitches

The single biggest frontend engineering problem is keeping the UI in sync with Redis state. Our rule: the WebSocket is the source of truth while the agent is active, and TanStack Query hydrates the initial page. On reconnect, the client requests a snapshot to catch up — it does not replay missed events, because missed events are too cheap to bother tracking in v1.

```typescript
// hooks/use-agent-socket.ts (sketch)
export function useAgentSocket() {
  const queryClient = useQueryClient();
  const socket = useMemo(() => new WebSocket(WS_URL + '/agent'), []);

  useEffect(() => {
    socket.onmessage = (m) => {
      const evt = JSON.parse(m.data);
      switch (evt.type) {
        case 'incoming_call':
          queryClient.setQueryData(['current_call'], evt.data);
          ringAudio.play();
          break;
        case 'call_ended':
          queryClient.setQueryData(['current_call'], null);
          queryClient.invalidateQueries(['recent_calls']);
          break;
        case 'status_changed':
          queryClient.setQueryData(['me'], (prev) => ({ ...prev, status: evt.data.status }));
          break;
      }
    };
    return () => socket.close();
  }, [socket]);
}
```

### 8.4 WebRTC softphone

The softphone embedded in the agent UI is a thin wrapper around the telephony provider's official JavaScript SDK — `africastalking-client` for Africa's Talking or `@twilio/voice-sdk` for Twilio. These SDKs handle the heavy lifting (signaling, ICE, audio codecs). We expose a minimal interface the UI can depend on.

```typescript
// softphone/client.ts
export interface Softphone {
  connect(accessToken: string): Promise<void>;
  accept(): void;
  reject(): void;
  hangup(): void;
  toggleMute(): boolean;    // returns new muted state
  toggleHold(): boolean;
  dtmf(digit: string): void;
}
```

### 8.5 Design system and component choices

- **shadcn/ui** as the base. We copy components into our repo instead of depending on a versioned package. This removes whole categories of dependency drift.
- **Tailwind** for layout. No hand-written CSS files except a 20-line global stylesheet. Design tokens live in `tailwind.config.ts`.
- **Radix primitives** for dropdowns, dialogs, and tooltips because they get accessibility right by default.
- **Lucide icons.** Consistent line weight, tree-shakeable, permissive license.
- **High contrast palette.** Agents often work in call-centre lighting. We bias toward readable, not fashionable.

---

## 9. Deployment Architecture

We deploy to two VPS nodes running Docker Compose. This is enough to carry the MVP into production with real customers. When volume justifies it, the same container images move unchanged to a managed Kubernetes service — no code changes needed.

*Figure 6. Two-node production deployment. Node A runs application containers, Node B runs stateful services. Backups and object storage are off-node.*

### 9.1 Why two nodes, not one, and not ten

One node is a single point of failure — the first power cut kills the platform. Ten nodes means a Kubernetes cluster, and a Kubernetes cluster means a part-time platform engineer we do not have. Two nodes give us hardware isolation between the app and the database, reasonable failure domains, and a bill we can defend.

### 9.2 Environments

| Environment | Host | Purpose |
|---|---|---|
| Local (dev) | Developer laptop | `docker-compose up` spins the full stack locally with seed data. |
| Staging | Contabo VPS (shared) | Mirror of production at smaller size. Every main-branch merge auto-deploys here. |
| Production | Contabo VPS nodes A+B | Serves real customers. Tagged releases only. Change-controlled. |

### 9.3 Container topology

The production Docker Compose file defines the following containers on the application node:

- `api` — the Go backend, 3 replicas behind NGINX. Stateless.
- `worker` — same binary as `api`, started with the `worker` subcommand. Consumes RabbitMQ.
- `web` — Next.js in production mode. 2 replicas. Stateless.
- `nginx` — reverse proxy, TLS, rate limiting, WebSocket upgrade.
- `n8n` — existing automations server, reused for non-critical workflows such as daily email digests.
- `prometheus`, `grafana`, `loki`, `jaeger` — observability stack, pinned versions.

On the data node:

- `postgres-primary` — PostgreSQL 16 with streaming replication to a hot standby.
- `postgres-standby` — read replica used for reports and a fast failover target.
- `redis` — Redis 7 with AOF persistence. Memory limit set below RAM ceiling.
- `rabbitmq` — persistent queues, management UI on a non-public port.
- `minio` — S3-compatible storage with lifecycle policies for recordings.

### 9.4 CI/CD pipeline

```yaml
# .github/workflows/deploy.yml (excerpt)
on:
  push:
    branches: [main]
    tags: ['v*']

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5   # Go 1.22
      - run: go test ./... -race -cover
      - run: cd web && npm ci && npm run lint && npm test

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: docker/build-push-action@v5
        with:
          tags: ghcr.io/${{ github.repository }}/api:${{ github.sha }}
          push: true

  deploy-staging:
    if: github.ref == 'refs/heads/main'
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: ssh deploy@staging 'cd /opt/cc && docker compose pull && docker compose up -d'

  deploy-production:
    if: startsWith(github.ref, 'refs/tags/v')
    needs: build
    environment: production      # GitHub approval gate
    runs-on: ubuntu-latest
    steps:
      - run: ssh deploy@prod 'cd /opt/cc && ./deploy.sh ${{ github.ref_name }}'
```

### 9.5 Zero-downtime deploys

NGINX is configured with `least_conn` upstream and the Go binary supports graceful shutdown (SIGTERM drains in-flight HTTP and WebSocket connections with a 30-second window). Deploying one `api` replica at a time gives us zero dropped requests. The Next.js rolling update uses the same mechanism.

### 9.6 Backup and restore

- Hourly `pg_basebackup` + WAL archiving to an off-site Backblaze B2 bucket, encrypted with SOPS-managed keys.
- Daily full Postgres dump kept for 7 days locally and 30 days off-site.
- MinIO bucket replication nightly to an off-site bucket.
- Restore drill executed monthly. The runbook includes the exact commands. A backup we have not practised restoring is a backup we do not have.

---

## 10. Security and Compliance

A call center handles personally identifiable information on nearly every call and often takes payment details verbally. Our obligations under the Kenya Data Protection Act (2019), GDPR (for any European callers), and PCI-DSS (if we ever handle card data) are real and personal — the data controller is legally accountable. Security is not a feature, it is the cost of being allowed to operate.

### 10.1 Threat model in one page

| Threat | Impact if realized | Primary mitigation |
|---|---|---|
| Stolen agent credentials | Attacker makes fraudulent outbound calls | Short-lived JWTs + refresh rotation + MFA for supervisors. |
| Webhook spoofing | Attacker injects fake inbound calls / steals routing | HMAC signature verification on every webhook; IP allowlist. |
| SQL injection | Data theft or deletion | pgx parameterized queries exclusively; no string interpolation. |
| Recording leak (bucket misconfig) | Major privacy breach, regulatory penalties | Private bucket, signed URLs with short TTL, quarterly access review. |
| Insider abuse by supervisor | Mass export of recordings | Audit log on every playback; anomaly alerts on download volume. |
| DDoS on /webhooks/* | Service unavailable | Cloudflare in front; provider-IP allowlist; per-IP rate limit in NGINX. |
| Lost recording backup | Cannot honour customer disclosure requests | Two geographic backups; monthly restore drill. |

### 10.2 Authentication and authorization

- **JWT (access) + opaque refresh token.** Access tokens live 15 minutes; refresh tokens 7 days and rotate on every use.
- **Argon2id** for password hashing. bcrypt is acceptable if Argon2 is not available in the chosen libraries.
- **MFA via TOTP** for supervisors and admins. Optional for agents in v1, mandatory in v2.
- **RBAC at the handler level.** Every protected handler declares the minimum role it accepts; a middleware enforces it before the handler runs.

### 10.3 Data protection at every layer

- **In transit.** TLS 1.2+ everywhere. Let's Encrypt certificates via cert-manager or manual certbot. HSTS and strong cipher suites in NGINX.
- **At rest.** Postgres on encrypted volumes (LUKS). MinIO server-side encryption with per-object keys. Backups encrypted with a separate key held in Doppler.
- **In memory.** PII is not logged, ever. A log scrubber middleware redacts known PII fields before the log line is written.
- **Retention.** Recordings default to 90 days. Customer request to delete is honoured within 30 days and logged in the audit trail.

### 10.4 The things that will actually bite you

> **Common mistakes to avoid:** Putting the recording bucket behind a long-lived signed URL; forgetting to rotate the JWT signing key after an engineer leaves; allowing webhook endpoints to accept any content-type; logging full HTTP request bodies in debug mode. Every one of these has been the root cause of a real incident at a real call center. Do not assume it will be different here.

---

## 11. Observability and Runbooks

If we cannot see the system, we cannot run it. Observability in this platform is not an afterthought — it is installed on day one, before the first call ever rings.

### 11.1 The three pillars

- **Metrics** — Prometheus scrapes `/metrics` on every container. Grafana dashboards visualize RED and USE metrics (rate, errors, duration; utilization, saturation, errors).
- **Logs** — structured JSON from Go (zerolog) and Next.js, shipped by Promtail into Loki. Queryable in Grafana alongside metrics.
- **Traces** — OpenTelemetry instruments every HTTP handler and outbound call. Jaeger surfaces slow requests and bottleneck components.

### 11.2 The twelve dashboards that must exist

1. Live call center overview — calls in progress, queue depth, average wait, agent availability.
2. Backend health — request rate, 4xx/5xx rates, p95 latency per endpoint.
3. WebSocket hub — open connections, message throughput, reconnect rate.
4. Postgres — connections, slow queries, WAL lag, table bloat.
5. Redis — memory usage, hit rate, evictions, slowlog.
6. RabbitMQ — queue depths, consumer throughput, unacked messages.
7. Telephony — calls per minute, provider error rate, per-minute cost.
8. Storage — MinIO bucket size, upload throughput, failed uploads.
9. Agent performance (per-agent) — calls handled, average handle time, adherence.
10. SLA compliance — percent of calls answered within target window.
11. Error budgets — actuals vs. SLOs for the month.
12. Cost — telephony minutes, infra spend, projected month-end total.

### 11.3 Alerts that actually wake people up

Alert fatigue is the enemy. We start with the minimum alert set and add only after we have felt the pain of not having one.

| Alert | Threshold | Severity |
|---|---|---|
| Error rate on /webhooks/* | > 2% for 3 minutes | Critical (page) |
| WebSocket hub disconnected | Any instance down > 60 seconds | Critical (page) |
| Postgres primary unreachable | > 60 seconds | Critical (page) |
| Redis memory above 85% | Sustained 5 minutes | High (Slack) |
| Queue wait time above SLA | Median above SLA for 5 minutes | High (Slack) |
| Disk above 80% on any node | Any point | Medium (Slack) |
| Recording upload backlog | > 100 messages for 10 minutes | Medium (Slack) |

### 11.4 SLOs and error budgets

- **Inbound call answer latency** (ring-in to ring-out to agent): < 3 seconds at p95, 30-day window.
- **Backend availability:** 99.5% month one, 99.9% after month three.
- **WebSocket delivery:** 99% of UI events delivered within 500 ms.
- **Recording availability:** 100% of completed calls have a recording within 15 minutes of end.

---

## 12. Project Management and Delivery Approach

We have 28 days to ship a production-ready MVP. That is not a lot of time. Every choice in this section is made with one goal: minimize time spent coordinating and maximize time spent building. The person who feels "too busy to plan" is the person who will ship six weeks late.

*Figure 7. The 28-day sprint plan. Four weekly milestones anchor the team and force early integration.*

### 12.1 Operating model: weekly sprints with daily standups

Four one-week sprints. Each sprint starts with a 45-minute planning meeting on Monday morning and ends with a 45-minute demo + retro on Friday afternoon. Daily standup is 15 minutes, same time every day, same three questions: what did I finish, what am I doing today, what is blocking me. That is the entire ceremony overhead.

### 12.2 Milestones

| Milestone | Target day | Acceptance criteria |
|---|---|---|
| M1 — Auth live | Day 7 | Admin can create a user, that user can log in, JWT flow works end-to-end, session persists across refresh. |
| M2 — First call connected | Day 14 | An inbound call rings in the telephony provider and lights up a test agent's browser with two-way audio. |
| M3 — Feature freeze | Day 21 | All v1 features merged and passing CI. Only bug fixes after this date. |
| M4 — Production go-live | Day 28 | Real business traffic on the platform, runbook validated, on-call rotation active. |

### 12.3 Week-by-week breakdown

#### Week 1 — Foundations (Days 1–7)

Goal: the team can build. Nothing customer-facing yet. We bootstrap the infrastructure, set up CI/CD, scaffold the Go backend and Next.js frontend, and ship working authentication.

- Provision both VPS nodes, install Docker, bring up Postgres / Redis / RabbitMQ containers with the production configuration.
- Register accounts with Africa's Talking and Twilio. Purchase two test DIDs, one Kenyan and one international.
- Create the two repos (backend, web) and seed them with skeleton code, linting, and formatters.
- Build the GitHub Actions pipeline. Enforce code review on main.
- Implement user registration, login, refresh, and RBAC middleware. Seed an admin user.
- Next.js scaffold with the design system, a login page, and a placeholder dashboard.
- Deploy to staging by Friday. Smoke-test login in a browser from two different networks.

#### Week 2 — Core routing (Days 8–14)

Goal: an inbound call can complete. This is the single highest-risk week. Everything we build afterwards depends on calls working.

- Webhook handler for Africa's Talking inbound POST, returning valid XML that dials a test agent.
- Agent state manager with Redis — online, available, busy, wrap-up transitions and heartbeats.
- WebSocket hub on the backend. Agent UI opens a socket on login and reflects state changes.
- Queue orchestrator with one default queue and "longest-idle agent" routing. Skills-based routing added if time allows.
- Agent softphone UI using the telephony provider's JS SDK. Accept, reject, mute, hang up.
- End-to-end smoke test: call the test DID from a mobile phone, see the agent browser ring, answer, talk, hang up.

#### Week 3 — Outbound, CRM, and supervisor tools (Days 15–21)

Goal: the platform is usable by a real agent and supervisor in a real shift.

- Outbound API + click-to-dial from the CRM view.
- Campaign module: upload CSV, schedule, start, pause, resume, export results.
- Customer/CRM service: search, view, tag, merge duplicates.
- Disposition capture workflow with mandatory category + free-text notes.
- Supervisor cockpit: live queue view, agent grid, listen/whisper/barge actions.
- Recording upload worker moving audio from provider to MinIO; updates `calls.recording_url`.

#### Week 4 — Harden and ship (Days 22–28)

Goal: hand the platform to real users without regret. Feature freeze Monday; no new features after Tuesday.

- Reporting and KPI service: daily, weekly, monthly reports with CSV export. Grafana dashboards for live KPIs.
- Load test with 500 concurrent simulated calls. Profile the Go backend, fix any p99 hotspots.
- Security review: dependency audit, secret rotation, OWASP top-10 sweep on the web app.
- UAT with 3–5 real agents for two half-days; run a bug bash afterwards.
- Write and rehearse the runbook. Record three short training videos (agent, supervisor, admin).
- Production cutover on Friday morning. Monitor closely through the weekend.

### 12.4 How we move this fast without breaking

- **Ruthless scoping.** The "explicitly out of scope" list is defended every week.
- **Pair programming on risky work.** Telephony webhooks, WebSocket reconnection, auth — always two people at the keyboard.
- **Vertical slices.** Every ticket delivers a visible user-facing change, even if small. No "the foundation is done, now I'll start the UI" week-long tickets.
- **End-to-end tests from day five.** The e2e suite catches regressions faster than the team notices them.
- **Ship on Fridays.** Nothing goes to production on Friday afternoon. Production deploys happen Tuesday through Thursday.
- **No meetings outside standup and sprint planning/demo.** Engineers need four-hour uninterrupted blocks. Protect them.

> **A hard truth:** A four-week timeline has zero tolerance for being clever. Use the boring solution for everything non-critical. Save your cleverness for the one or two places it genuinely matters — almost certainly the routing logic.

---

## 13. Team Structure, Roles, and Communication

### 13.1 The minimum viable team

This platform can be built by four to five focused people. More than six and coordination cost eats the speed gain; fewer than four and key skills are missing.

| Role | Count | Primary responsibilities |
|---|---|---|
| Tech Lead / Backend Engineer | 1 | Architecture decisions, Go backend core, code review, production readiness. |
| Backend Engineer | 1 | Go services (CRM, campaigns, reporting), webhook integrations, workers. |
| Frontend Engineer | 1 | Next.js app, design system, softphone integration, supervisor cockpit. |
| DevOps / SRE (fractional) | 0.5 | VPS provisioning, CI/CD, observability stack, incident response on call. |
| Product / Delivery Manager | 1 | Backlog, stakeholder communication, daily standup facilitation, UAT. |
| QA (from week 3) | 0.5 | Manual exploratory testing, UAT coordination, bug triage during the bug bash. |

### 13.2 Decision-making

We use a two-tier decision framework so small choices do not require meetings and large ones do not get made accidentally.

- **Reversible decisions (Type 2):** made by whoever is doing the work. Documented in the pull request description. No meeting required.
- **Irreversible or cross-cutting decisions (Type 1):** require a short (1-page) decision record in `docs/decisions/`, reviewed by the tech lead and one other engineer before merging.

### 13.3 Communication channels

| Channel | Purpose | Expected response time |
|---|---|---|
| #cc-standup (Slack) | Async standup notes when someone cannot attend live. | N/A |
| #cc-engineering (Slack) | Day-to-day engineering chatter, code questions. | Business hours, same day |
| #cc-alerts (Slack) | All automated alerts from Prometheus, Sentry, uptime. | Immediate for critical |
| #cc-incidents (Slack) | Live incident coordination only. | Immediate during incident |
| Linear or GitHub Projects | Backlog, sprints, tickets. | Daily grooming |
| Weekly demo (video call) | Friday afternoon, all stakeholders. | Scheduled |

### 13.4 The on-call rotation

Starting the day we go live, one engineer is on call 24/7 on a one-week rotation. The on-call engineer responds to pages within 15 minutes, acknowledges in #cc-incidents, and leads triage. Every incident produces a blameless postmortem within 48 hours, filed in the repository.

---

## 14. Testing Strategy

A call center cannot ship broken. Every minute of downtime is a customer who got a busy tone. Our testing strategy is pragmatic: we invest heavily where the cost of a bug is highest, and lightly where a bug is cheap to notice and fix.

### 14.1 The testing pyramid for this project

| Layer | Where | Coverage target | Tools |
|---|---|---|---|
| Unit tests | Pure functions, domain logic, helpers | Tight on routing, queue, dialer; thin elsewhere. | Go testing + testify |
| Integration tests | DB queries, Redis ops, handler → service | Every REST endpoint has one happy and one sad path. | testcontainers-go, dockertest |
| Contract tests | Webhook payloads from Africa's Talking/Twilio | One test per event type with a captured fixture. | Fixture-based |
| End-to-end tests | Playwright-driven UI + mocked telephony | Login, inbound call, outbound call, disposition. | Playwright |
| Load tests | API + WS under sustained concurrency | 500 concurrent WS, 100 calls/minute sustained. | k6, Vegeta |
| Manual UAT | Real agents on real phones | Pre-launch bug bash, 2 half-days. | Checklist |

### 14.2 The tests we never skip

- **Routing determinism test.** Given a specific Redis state, the same call must always pick the same agent. Flakiness here causes production chaos.
- **WebSocket reconnection test.** Simulate a dropped connection and verify the agent sees the current call within 2 seconds of reconnect.
- **Webhook signature negative tests.** Invalid signature, replayed event, malformed payload — all must be rejected with the right status code.
- **Database migration dry-run** against a snapshot of production data before every release.

### 14.3 Test data

We maintain a seed dataset that reflects realistic production shapes (not toy data): 50 agents across 3 teams, 5 queues, 10,000 customers, and 100,000 historical calls. This is loaded into staging and local dev. Performance tests that pass against 10 customers but fail against 10,000 are worse than no tests at all.

### 14.4 Load testing the routing engine

```javascript
// k6 script sketch — scripts/load/route.js
import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 50 },   // ramp to 50 virtual agents
    { duration: '3m', target: 500 },  // sustain 500 WS connections
    { duration: '1m', target: 0  },
  ],
  thresholds: {
    'ws_msgs_received': ['rate>0'],
    'ws_session_duration': ['p(95)<60000'],
  },
};
export default function () {
  const url = __ENV.WS_URL + '/agent';
  ws.connect(url, { headers: { Authorization: 'Bearer ' + __ENV.TOKEN } },
    (socket) => { socket.on('open', () => socket.setInterval(() => socket.ping(), 10000)); });
}
```

---

## 15. Go-Live Checklist

This checklist is executed top to bottom on the day of production cutover. Nothing is skipped. Every line has a named owner and a checkbox in the actual runbook document.

### 15.1 T-7 days before launch

- [ ] All production secrets rotated and stored in Doppler/SOPS. No dev values remain.
- [ ] DNS records for `api.example.com`, `app.example.com`, and `status.example.com` created with 60-second TTL during cutover week.
- [ ] Production TLS certificates issued and installed. Auto-renewal confirmed.
- [ ] Pager rotation published in PagerDuty (or equivalent). Primary and secondary identified.
- [ ] Runbook reviewed and signed off by the tech lead and at least one other engineer.
- [ ] Backup restore drill completed successfully in staging.

### 15.2 T-1 day before launch

- [ ] Final staging-to-production database sync or fresh-start decision confirmed.
- [ ] Production seed data loaded (teams, queues, initial users).
- [ ] Phone numbers provisioned with Africa's Talking and configured to point at production webhook URLs.
- [ ] Monitoring dashboards checked for false positives. Alert channels verified.
- [ ] All known bugs either fixed or triaged with a documented workaround.
- [ ] Go / no-go meeting held. Decision documented.

### 15.3 Launch day

1. Freeze main branch. No deploys except the launch deploy.
2. Deploy production containers. Verify all health checks pass.
3. Walk through the smoke-test checklist: admin login, create user, user login, inbound call, outbound call, disposition, recording playback, report export.
4. Have a real agent complete a real call end-to-end under observation.
5. Announce launch to stakeholders. Publish the status page.
6. Engineering team stays online for the first four hours. The tech lead stays on call for 24 hours.

### 15.4 Post-launch first 72 hours

- Watch error rate, latency, and WebSocket connection stability every 30 minutes.
- Daily 30-minute incident review with the whole team.
- Any user-reported issue becomes a ticket within an hour of being reported.
- First weekly retrospective on day 7 with a written summary distributed to stakeholders.

---

## 16. Risk Register

Every project has risks. The only safe thing to do with a risk is name it, estimate it, and decide in advance how to handle it. Unnamed risks are the ones that blow up schedules.

| Risk | Likelihood | Impact | Mitigation / response |
|---|---|---|---|
| Telephony provider API changes mid-sprint | Medium | High | Pin SDK versions. Capture sample webhook payloads as fixtures. Keep a secondary provider integrated as fallback. |
| Agent browser WebRTC quirks (Safari on iOS) | High | Medium | Declare Chrome and Firefox as tier-1. iOS Safari is tier-2 for v1. Document the restriction clearly. |
| Underestimated load test results | Medium | High | Run load tests weekly from week 2, not just in week 4. Fix issues as they appear, not at the finish line. |
| Key engineer sick or unavailable | Medium | High | Pair programming on critical paths. Every PR is reviewed. No knowledge silos. |
| Kenya Data Protection Act audit finding | Low | High | Register as a data controller before go-live. Publish the privacy policy. Keep DPA documentation in a shared drive. |
| Budget overrun on telephony minutes during load testing | High | Low | Use a low-cost sandbox number and simulate traffic with mocks in CI. Reserve real minute tests for one scheduled window. |
| Scope creep from stakeholders | Very high | High | Out-of-scope list in Section 1.4 is the contract. Every new request is a v2 ticket by default. |
| Postgres outage during launch week | Low | Severe | Hot standby configured from day one. Practise failover in staging before launch. |
| Recording bucket grows beyond budget | Medium | Medium | Lifecycle policy to cold storage at 30 days, delete at 90. Monthly cost review. |
| On-call burnout after launch | Medium | Medium | Rotate weekly. Any after-hours page is followed by compensating time off. |

---

## 17. Post-Launch Roadmap

The four-week MVP is the start line, not the finish line. These are the investments we expect to make in months 2–6, roughly in priority order. Each theme is a month of work for the same team size.

### 17.1 Month 2 — Stability and polish

- Bug fixes surfaced from real usage.
- Mobile-responsive supervisor cockpit.
- Bulk import/export for users and customers.
- First pass at multi-queue skills-based routing with priority weights.
- Automated email digests for supervisors (daily shift summary).

### 17.2 Month 3 — Integrations

- WhatsApp inbound channel integrated into the agent workspace (conversations, not only notifications).
- SMS outbound and inbound routed through the same queue engine.
- Zapier / n8n webhook endpoints so third-party systems can react to call events.
- HubSpot/Salesforce two-way contact sync (optional, depending on demand).

### 17.3 Month 4 — Intelligence

- Automatic call transcription (Deepgram or Whisper).
- Post-call summarization using an LLM; agent reviews and edits, not generates from scratch.
- Real-time sentiment indicator in the supervisor cockpit.
- IVR visual builder (drag-and-drop voice flow editor).

### 17.4 Month 5 — Scale

- Migrate from Docker Compose to a managed Kubernetes cluster.
- Read replicas for reporting queries; dedicated analytics database.
- Multi-tenancy with row-level security (if commercializing as SaaS).
- Evaluate moving to self-hosted FreeSWITCH to reduce telephony costs at scale.

### 17.5 Month 6 — Operational maturity

- SOC 2 readiness: formal access reviews, audit logs, change management records.
- Disaster recovery site in a second region.
- SLAs published to customers with monthly uptime reports.
- Published status page with historical incident record.

> **The discipline for the roadmap:** Nothing on this list is promised to anyone outside the team until it is in progress. "In the roadmap" means we have studied it and believe it is the right next step. "Shipping" means code is merged and deployed. Confusing those two creates debt with customers that is harder to repay than code debt.

---

## 18. Appendices

### 18.1 Glossary

| Term | Meaning |
|---|---|
| ACD | Automatic Call Distributor — the part of the platform that routes calls to agents. |
| AHT | Average Handle Time — mean duration agents spend on a call, including wrap-up. |
| Barge | A supervisor joining a live call as a full participant. |
| CSAT | Customer Satisfaction score — typically a post-call survey rating. |
| DID | Direct Inward Dial — a phone number that routes straight to the platform. |
| DNIS | Dialed Number Identification Service — the number the customer dialed. |
| DTMF | Touch-tone keypad input, used in IVRs. |
| IVR | Interactive Voice Response — the menu a caller navigates ("Press 1 for sales"). |
| PSTN | Public Switched Telephone Network — the traditional phone network. |
| RBAC | Role-Based Access Control — permissions granted by job role. |
| SDP | Session Description Protocol — the WebRTC handshake payload. |
| SIP | Session Initiation Protocol — the signaling protocol used between phone systems. |
| SLA | Service Level Agreement — a measurable commitment (e.g., 80% of calls answered within 20 s). |
| SLO | Service Level Objective — our internal target that drives SLA compliance. |
| TwiML | Twilio's XML instruction format (Africa's Talking uses an equivalent format). |
| WebRTC | Browser-native real-time audio/video technology used for the softphone. |
| Whisper | Supervisor speaking only to the agent during a live call. |
| Wrap-up | Post-call state where the agent finishes notes before taking the next call. |

### 18.2 Repository structure

```
callcenter/
├── backend/                  # Go module
│   ├── cmd/
│   │   ├── api/main.go       # REST + WS server entry point
│   │   └── worker/main.go    # RabbitMQ consumer entry point
│   ├── internal/
│   │   ├── auth/             # JWT, password hashing, RBAC
│   │   ├── routing/          # Call router, skills matcher
│   │   ├── agents/           # Agent state manager
│   │   ├── queues/           # Queue orchestrator
│   │   ├── calls/            # Call lifecycle, dispositions
│   │   ├── campaigns/        # Outbound campaigns + dialer
│   │   ├── crm/              # Customer service
│   │   ├── recordings/       # Upload worker, lifecycle
│   │   ├── reports/          # KPI computation
│   │   ├── telephony/        # Provider adapters (africas_talking, twilio)
│   │   └── transport/
│   │       ├── http/         # chi handlers
│   │       └── ws/           # Gorilla WebSocket hub
│   ├── migrations/           # SQL migrations
│   ├── sqlc/                 # Generated type-safe queries
│   └── go.mod
├── web/                      # Next.js 15 app
│   ├── app/
│   │   ├── (agent)/          # Agent workspace
│   │   ├── (supervisor)/     # Supervisor cockpit
│   │   └── (admin)/          # Admin + reports
│   ├── components/ui/        # shadcn components
│   ├── lib/                  # API client, WS hook, utils
│   └── package.json
├── infra/
│   ├── docker-compose.yml    # Production
│   ├── docker-compose.dev.yml
│   ├── nginx/
│   └── grafana/              # Dashboard JSON files
├── docs/
│   ├── runbook.md
│   ├── decisions/            # ADRs
│   └── this_document.docx
└── README.md
```

### 18.3 Further reading

- *Designing Data-Intensive Applications* — Martin Kleppmann. Core chapters on replication, consistency, and batch vs stream processing.
- *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi. Short, practical, highly applicable.
- Africa's Talking Voice API documentation — the official reference for our primary telephony provider.
- Twilio Voice docs — equally useful; TwiML concepts apply to Africa's Talking's XML responses.
- The PostgreSQL documentation, specifically the chapters on indexes, concurrency, and the query planner.
- *Google SRE Book* — especially the chapters on SLOs, alerting, and postmortems.

### 18.4 Living document policy

This document is owned by the tech lead. It is reviewed at the end of every sprint and updated whenever an architectural decision, a tool choice, or a workflow changes. The docx file is regenerated from markdown in the repository so the same content is viewable in a browser and in Word.

### 18.5 Revision history

| Version | Date | Author | Summary |
|---|---|---|---|
| 1.0 | April 2026 | Claude + project team | Initial blueprint. Covers architecture, stack, data flows, security, delivery plan, and operations for the 28-day MVP. |
