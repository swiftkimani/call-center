CREATE TABLE queues (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name              TEXT NOT NULL,
  description       TEXT NOT NULL DEFAULT '',
  skills_required   TEXT[]  NOT NULL DEFAULT '{}',
  max_wait_seconds  INT     NOT NULL DEFAULT 300,
  sla_seconds       INT     NOT NULL DEFAULT 20,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE queue_agents (
  queue_id UUID NOT NULL REFERENCES queues (id) ON DELETE CASCADE,
  agent_id UUID NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
  PRIMARY KEY (queue_id, agent_id)
);
