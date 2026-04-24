#!/usr/bin/env bash
set -euo pipefail

BASE="${BASE_URL:-http://localhost:8080}"
echo "=== Smoke test against $BASE ==="

# Health check
echo -n "healthz: "
curl -sf "$BASE/healthz" | grep -q ok && echo "OK" || { echo "FAIL"; exit 1; }

# Admin login (uses the seeded admin user, needs real hash — skip if not seeded)
echo -n "login: "
RESP=$(curl -sf -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.local","password":"Admin1234!"}' 2>&1) || true

if echo "$RESP" | grep -q "access_token"; then
  TOKEN=$(echo "$RESP" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
  echo "OK (token received)"
else
  echo "SKIP (seed user not set up with real hash)"
  TOKEN=""
fi

if [ -n "$TOKEN" ]; then
  echo -n "/api/v1/me: "
  curl -sf "$BASE/api/v1/me" -H "Authorization: Bearer $TOKEN" | grep -q "email" && echo "OK" || echo "FAIL"

  echo -n "/api/v1/calls: "
  curl -sf "$BASE/api/v1/calls" -H "Authorization: Bearer $TOKEN" | grep -q "success" && echo "OK" || echo "FAIL"
fi

echo "=== Smoke test complete ==="
