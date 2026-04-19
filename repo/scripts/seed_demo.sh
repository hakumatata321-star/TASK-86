#!/usr/bin/env bash
# seed_demo.sh — Create demo accounts for local evaluation.
#
# Run this script AFTER docker-compose up and AFTER you have retrieved the
# admin password from the container logs:
#
#   docker-compose logs portal | grep SECURITY
#
# Usage:
#   ADMIN_PASSWORD=<pass> bash scripts/seed_demo.sh [BASE_URL]
#
# The script creates the following accounts (all with password Demo1234!):
#   demo_student    (student)
#   demo_instructor (instructor)
#   demo_clerk      (clerk)
#   demo_moderator  (moderator)
#   demo_manager    (manager)

set -euo pipefail

ADMIN_PASS="${ADMIN_PASSWORD:?Set ADMIN_PASSWORD to the admin password from docker-compose logs}"
BASE_URL="${1:-http://localhost:3000}"
COOKIE_JAR="$(mktemp)"
trap 'rm -f "$COOKIE_JAR"' EXIT

echo "Logging in as admin at ${BASE_URL}…"

# Step 1: GET /login to obtain a CSRF token cookie.
curl -s -c "$COOKIE_JAR" "${BASE_URL}/login" -o /dev/null

# Extract csrf_token cookie value.
CSRF=$(grep csrf_token "$COOKIE_JAR" | awk '{print $NF}' || true)

# Step 2: POST /login with credentials and CSRF token.
LOGIN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
  -c "$COOKIE_JAR" -b "$COOKIE_JAR" \
  -X POST "${BASE_URL}/login" \
  -H "X-Csrf-Token: ${CSRF}" \
  -d "username=admin&password=${ADMIN_PASS}" \
  --max-redirs 5 -L)

if [[ "$LOGIN_STATUS" != "200" ]]; then
  echo "ERROR: Admin login failed (HTTP ${LOGIN_STATUS})."
  echo "       Make sure the server is running and ADMIN_PASSWORD is correct."
  exit 1
fi

# Refresh CSRF token after login.
CSRF=$(grep csrf_token "$COOKIE_JAR" | awk '{print $NF}' || true)

create_user() {
  local username=$1 email=$2 password=$3 role=$4

  STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -c "$COOKIE_JAR" -b "$COOKIE_JAR" \
    -X POST "${BASE_URL}/admin/users" \
    -H "X-Csrf-Token: ${CSRF}" \
    -d "username=${username}&email=${email}&password=${password}&role=${role}" \
    --max-redirs 5 -L)

  if [[ "$STATUS" == "200" ]] || [[ "$STATUS" == "302" ]] || [[ "$STATUS" == "303" ]]; then
    echo "  OK   ${role}: ${username} (${email})  password=${password}"
  else
    echo "  SKIP ${role}: ${username} — HTTP ${STATUS} (account may already exist)"
  fi

  # Re-fetch CSRF token for next request.
  CSRF=$(grep csrf_token "$COOKIE_JAR" | awk '{print $NF}' || echo "$CSRF")
}

echo ""
echo "Creating demo accounts…"
create_user "demo_student"    "demo_student@portal.local"    "Demo1234!" "student"
create_user "demo_instructor" "demo_instructor@portal.local" "Demo1234!" "instructor"
create_user "demo_clerk"      "demo_clerk@portal.local"      "Demo1234!" "clerk"
create_user "demo_moderator"  "demo_moderator@portal.local"  "Demo1234!" "moderator"
create_user "demo_manager"    "demo_manager@portal.local"    "Demo1234!" "manager"

echo ""
echo "Done.  Credentials (all non-admin accounts):"
echo "  Password: Demo1234!"
echo ""
echo "  See README → Demo Credentials for the full table."
