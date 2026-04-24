CREATE TABLE customers (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  phone_number TEXT UNIQUE NOT NULL,
  full_name    TEXT NOT NULL DEFAULT '',
  email        TEXT,
  tags         TEXT[] NOT NULL DEFAULT '{}',
  dnc_listed   BOOL   NOT NULL DEFAULT false,
  timezone     TEXT   NOT NULL DEFAULT 'Africa/Nairobi',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_phone ON customers (phone_number);
CREATE INDEX idx_customers_tags_gin ON customers USING GIN (tags);
