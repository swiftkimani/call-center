-- Dev seed data — runs on fresh DB volume only
-- Password is "Admin1234!" (argon2id hash)

INSERT INTO teams (id, name) VALUES
  ('00000000-0000-0000-0000-000000000001', 'Support'),
  ('00000000-0000-0000-0000-000000000002', 'Sales')
ON CONFLICT DO NOTHING;

INSERT INTO users (id, email, password_hash, full_name, role) VALUES
  (
    '00000000-0000-0000-0000-000000000010',
    'admin@test.local',
    '$argon2id$v=19$m=65536,t=2,p=4$kTuV28t9cqSJi9lHcE0I3A$/IdrAZJpww4EScDVcyMPKF+WEiCt4nOEJsqGjmiuGVA',
    'Admin User',
    'admin'
  ),
  (
    '00000000-0000-0000-0000-000000000011',
    'supervisor@test.local',
    '$argon2id$v=19$m=65536,t=2,p=4$OCuStBCtxdOkKJvsML8VDQ$pZcoKjJtPGWaeNJCHxDPuE1XaPCNqRZa+Ak4byipa1U',
    'Supervisor User',
    'supervisor'
  ),
  (
    '00000000-0000-0000-0000-000000000012',
    'agent@test.local',
    '$argon2id$v=19$m=65536,t=2,p=4$5Hob6B9xespVxAKuXUxXAw$nYVnljIgnY35ePS2G16oy7g275ya9XeBldQMKOOsb9Q',
    'Agent User',
    'agent'
  )
ON CONFLICT DO NOTHING;

INSERT INTO agents (id, user_id, extension, skills, team_id) VALUES
  (
    '00000000-0000-0000-0000-000000000020',
    '00000000-0000-0000-0000-000000000012',
    '1001',
    ARRAY['general','sales'],
    '00000000-0000-0000-0000-000000000002'
  )
ON CONFLICT DO NOTHING;

INSERT INTO queues (id, name, description, sla_seconds) VALUES
  ('00000000-0000-0000-0000-000000000030', 'General', 'General support queue', 20),
  ('00000000-0000-0000-0000-000000000031', 'Sales', 'Sales queue', 15)
ON CONFLICT DO NOTHING;
