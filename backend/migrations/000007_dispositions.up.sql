CREATE TABLE disposition_categories (
  id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name   TEXT NOT NULL UNIQUE,
  active BOOL NOT NULL DEFAULT true
);

INSERT INTO disposition_categories (name) VALUES
  ('resolved'),
  ('callback_requested'),
  ('escalated'),
  ('no_action'),
  ('wrong_number');

CREATE TABLE dispositions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  call_id    UUID NOT NULL UNIQUE REFERENCES calls (id) ON DELETE CASCADE,
  agent_id   UUID NOT NULL REFERENCES agents (id),
  category   TEXT NOT NULL,
  notes      TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dispositions_agent ON dispositions (agent_id);
