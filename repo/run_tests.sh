#!/usr/bin/env bash
# run_tests.sh — Execute all test suites and print a pass/fail summary.
#
# ALL tests run inside Docker for a reproducible, isolated environment.
# No local Go, Node, or gcc installation is required on the host.
#
# Usage:
#   ./run_tests.sh            # run all suites
#   ./run_tests.sh -v         # verbose (stream test output)
#   ./run_tests.sh -race      # enable Go race detector
#   ./run_tests.sh -v -race   # both
#   REBUILD=1 ./run_tests.sh  # force rebuild the test Docker image
#
# Performance notes:
#   • Pre-built test image (w2t86-test) — Node.js baked in, no apt-get per run.
#   • Named Docker volumes cache Go modules + npm node_modules across runs.

set -euo pipefail

# ---------------------------------------------------------------------------
# Docker wrapper — always run inside a known-good container
# ---------------------------------------------------------------------------
if [[ -z "${IN_DOCKER:-}" ]]; then
  if ! command -v docker &>/dev/null; then
    echo "ERROR: Docker is required to run the test suite."
    echo "       Install Docker Desktop (https://docs.docker.com/get-docker/) and retry."
    exit 1
  fi

  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  TEST_IMAGE="w2t86-test:latest"

  # Build the test image if it doesn't exist or REBUILD=1 is set.
  # The Dockerfile is embedded inline so no external file is required.
  if [[ "${REBUILD:-0}" == "1" ]] || ! docker image inspect "$TEST_IMAGE" &>/dev/null; then
    echo "Building test image (${TEST_IMAGE})…"
    docker build -q -t "$TEST_IMAGE" - <<'DOCKERFILE'
FROM golang:1.22-bookworm
RUN apt-get update -qq && \
    apt-get install -y -qq --no-install-recommends nodejs npm && \
    rm -rf /var/lib/apt/lists/*
ENV CGO_ENABLED=1
DOCKERFILE
    echo "Test image ready."
  fi

  echo "Launching test container (${TEST_IMAGE})…"
  exec docker run --rm \
    -e IN_DOCKER=1 \
    -e CGO_ENABLED=1 \
    -v "${SCRIPT_DIR}:/workspace" \
    -v "w2t86-gomod:/root/go/pkg/mod" \
    -v "w2t86-npm:/workspace/web/node_modules" \
    -w /workspace \
    "$TEST_IMAGE" \
    bash /workspace/run_tests.sh "$@"
fi

# ---------------------------------------------------------------------------
# Colour helpers (running inside Docker from here on)
# ---------------------------------------------------------------------------
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
RESET="\033[0m"

pass() { echo -e "${GREEN}PASS${RESET}  $1"; }
fail() { echo -e "${RED}FAIL${RESET}  $1"; }
info() { echo -e "${YELLOW}----${RESET}  $1"; }

# ---------------------------------------------------------------------------
# Parse flags
# ---------------------------------------------------------------------------
VERBOSE=""
RACE=""
for arg in "$@"; do
  case "$arg" in
    -v)     VERBOSE="-v" ;;
    -race)  RACE="-race" ;;
    *)      echo "Unknown flag: $arg"; exit 1 ;;
  esac
done

# ---------------------------------------------------------------------------
# Locate the project root.
# ---------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

export MIGRATION_PATH="$SCRIPT_DIR/migrations/001_schema.sql"

# ---------------------------------------------------------------------------
# Suite definitions.
#
# The seven integration sub-suites are merged into a single invocation —
# they all target the same package and running them separately just repeats
# compilation overhead.
# ---------------------------------------------------------------------------
declare -a SUITE_LABELS=(
  "Unit tests (state machine, inventory, auth, validation, rate-limiter)"
  "API functional tests (normal inputs, bad params, permission errors)"
  "Integration tests (auth, materials, orders, distribution, messaging, moderation, admin)"
  "Config & scheduler unit tests"
  "Middleware unit tests (auth, RBAC)"
  "Service unit tests (analytics, admin, courses, moderation)"
)
declare -a SUITE_PKGS=(
  "./unit_tests/..."
  "./API_tests/..."
  "./internal/integration/"
  "./internal/config/... ./internal/scheduler/..."
  "./internal/middleware/..."
  "./internal/services/..."
)
declare -a SUITE_RUN=(
  ""
  ""
  ""
  ""
  ""
  ""
)

# ---------------------------------------------------------------------------
# Run each Go suite sequentially.
# ---------------------------------------------------------------------------
TOTAL=0
PASSED=0
FAILED=0
FAILED_SUITES=()

echo ""
echo "========================================================"
echo " w2t86 — Test Runner (Docker)"
echo " $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================================"
echo ""

for i in "${!SUITE_LABELS[@]}"; do
  label="${SUITE_LABELS[$i]}"
  pkg="${SUITE_PKGS[$i]}"
  run_filter="${SUITE_RUN[$i]}"

  TOTAL=$((TOTAL + 1))
  info "Running: $label"

  cmd=(go test -tags sqlite_fts5 $RACE $VERBOSE -timeout 300s)
  [[ -n "$run_filter" ]] && cmd+=(-run "$run_filter")
  IFS=' ' read -ra pkg_list <<< "$pkg"
  cmd+=("${pkg_list[@]}")

  tmp=$(mktemp)
  if "${cmd[@]}" >"$tmp" 2>&1; then
    pass "$label"
    PASSED=$((PASSED + 1))
    [[ -n "$VERBOSE" ]] && cat "$tmp"
  else
    fail "$label"
    FAILED=$((FAILED + 1))
    FAILED_SUITES+=("$label")
    cat "$tmp"
  fi
  rm -f "$tmp"
  echo ""
done

# ---------------------------------------------------------------------------
# Frontend tests (Vitest — app.js, map.js)
# ---------------------------------------------------------------------------
FE_LABEL="Frontend unit tests (Vitest — app.js, map.js)"
TOTAL=$((TOTAL + 1))
info "Running: $FE_LABEL"

tmp=$(mktemp)
if (cd web && npm install --silent && npm test) >"$tmp" 2>&1; then
  pass "$FE_LABEL"
  PASSED=$((PASSED + 1))
  [[ -n "$VERBOSE" ]] && cat "$tmp"
else
  fail "$FE_LABEL"
  FAILED=$((FAILED + 1))
  FAILED_SUITES+=("$FE_LABEL")
  cat "$tmp"
fi
rm -f "$tmp"
echo ""

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo "========================================================"
echo " Results: ${PASSED}/${TOTAL} suites passed"
echo "========================================================"

if [[ $FAILED -gt 0 ]]; then
  echo ""
  echo -e "${RED}Failed suites:${RESET}"
  for s in "${FAILED_SUITES[@]}"; do
    echo -e "  ${RED}✘${RESET} $s"
  done
  echo ""
  exit 1
else
  echo ""
  echo -e "${GREEN}All suites passed.${RESET}"
  echo ""
  exit 0
fi
