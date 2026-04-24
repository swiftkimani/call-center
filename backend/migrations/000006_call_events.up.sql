CREATE TABLE call_events (
  id          BIGSERIAL PRIMARY KEY,
  call_id     UUID NOT NULL REFERENCES calls (id) ON DELETE CASCADE,
  event_type  TEXT NOT NULL,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  payload     JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_call_events_call_occurred ON call_events (call_id, occurred_at);
