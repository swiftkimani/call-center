CREATE TABLE calls (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  provider_sid  TEXT UNIQUE NOT NULL,
  customer_id   UUID REFERENCES customers (id),
  agent_id      UUID REFERENCES agents (id),
  queue_id      UUID REFERENCES queues (id),
  direction     TEXT NOT NULL CHECK (direction IN ('inbound', 'outbound')),
  status        TEXT NOT NULL CHECK (status IN (
                  'queued', 'ringing', 'in_progress', 'completed',
                  'failed', 'no_answer', 'abandoned')),
  started_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  answered_at   TIMESTAMPTZ,
  ended_at      TIMESTAMPTZ,
  wait_seconds  INT,
  talk_seconds  INT,
  recording_url TEXT,
  cost_cents    INT,
  from_number   TEXT NOT NULL,
  to_number     TEXT NOT NULL
);

CREATE INDEX idx_calls_customer_started ON calls (customer_id, started_at DESC);
CREATE INDEX idx_calls_agent_started    ON calls (agent_id, started_at DESC);
CREATE INDEX idx_calls_started_brin     ON calls USING BRIN (started_at);
CREATE INDEX idx_calls_status           ON calls (status);
