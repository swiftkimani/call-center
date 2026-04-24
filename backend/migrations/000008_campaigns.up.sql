CREATE TABLE campaigns (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name         TEXT NOT NULL,
  status       TEXT NOT NULL DEFAULT 'draft'
               CHECK (status IN ('draft', 'active', 'paused', 'completed')),
  scheduled_at TIMESTAMPTZ,
  created_by   UUID NOT NULL REFERENCES users (id),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE campaign_contacts (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  campaign_id  UUID NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE,
  customer_id  UUID NOT NULL REFERENCES customers (id),
  status       TEXT NOT NULL DEFAULT 'pending'
               CHECK (status IN ('pending', 'dialing', 'completed', 'failed', 'dnc')),
  attempted_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ
);

CREATE INDEX idx_campaign_contacts_campaign_status ON campaign_contacts (campaign_id, status);
