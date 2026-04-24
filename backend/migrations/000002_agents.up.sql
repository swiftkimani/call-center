CREATE TABLE agents (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
  extension      TEXT UNIQUE NOT NULL,
  skills         TEXT[]   NOT NULL DEFAULT '{}',
  status         TEXT     NOT NULL DEFAULT 'offline'
                          CHECK (status IN ('offline', 'available', 'busy', 'wrap_up', 'break')),
  max_concurrent SMALLINT NOT NULL DEFAULT 1,
  team_id        UUID REFERENCES teams (id),
  last_seen_at   TIMESTAMPTZ
);

CREATE INDEX idx_agents_status_partial ON agents (status)
  WHERE status IN ('available', 'busy');

CREATE INDEX idx_agents_skills_gin ON agents USING GIN (skills);
