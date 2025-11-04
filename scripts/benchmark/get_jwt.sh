#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://127.0.0.1:8080"
EMAIL="loadtest@example.com"
PASSWORD="Password123"
FIRST_NAME="Load"
LAST_NAME="Tester"
PHONE="9999999999"

# ────────────────────────────────
# 1️ Signup (if not exists)
# ────────────────────────────────
echo "[+] Attempting signup for $EMAIL ..."
signup_resp=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -d "{\"firstName\":\"$FIRST_NAME\",\"lastName\":\"$LAST_NAME\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"phone\":\"$PHONE\"}")

signup_body=$(echo "$signup_resp" | sed '$d')
signup_code=$(echo "$signup_resp" | tail -n1)

if [[ "$signup_code" -ge 200 && "$signup_code" -lt 300 ]]; then
  echo "[+] Signup succeeded for $EMAIL"
else
  echo "[i] Signup returned HTTP $signup_code — user may already exist. Proceeding..."
fi

# ────────────────────────────────
# 2️ Login to get JWT
# ────────────────────────────────
echo "[+] Logging in to fetch JWT..."
login_resp=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

JWT=$(echo "$login_resp" | jq -r '.access_token // .token // .jwt // .data.access_token' 2>/dev/null)

if [[ -z "$JWT" || "$JWT" == "null" ]]; then
  echo " ERROR: Could not extract JWT from login response:"
  echo "$login_resp"
  exit 1
fi

export JWT_TOKEN="$JWT"
echo "[+] JWT_TOKEN captured successfully (length: ${#JWT_TOKEN})"

# ────────────────────────────────
# 3️ Patch Lua benchmark with token dynamically
# ────────────────────────────────
echo "[+] Injecting JWT token into benchmark script..."
tmpfile=$(mktemp)
sed "s|local JWT_TOKEN = nil|local JWT_TOKEN = '${JWT_TOKEN}'|" appointment_create_benchmark.lua > "$tmpfile"

# ────────────────────────────────
# 4️ Execute benchmark
# ────────────────────────────────
echo "[+] Running load test (20s, 8 threads, 50 connections)..."
wrk -t12 -c200 -d20s -s "$tmpfile" --latency "$BASE_URL"

# ────────────────────────────────
# 5️ Cleanup
# ────────────────────────────────
rm -f "$tmpfile"
