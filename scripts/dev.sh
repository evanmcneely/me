#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BLOG_ADDR="${BLOG_ADDR:-:8080}"
BROWSER_SYNC_PORT="${BROWSER_SYNC_PORT:-3000}"

port_from_addr() {
  local addr="$1"
  echo "${addr##*:}"
}

require_free_port() {
  local port="$1"
  local label="$2"
  if lsof -iTCP:"$port" -sTCP:LISTEN -n -P >/dev/null 2>&1; then
    echo "$label port $port is already in use" >&2
    exit 1
  fi
}

cd "$ROOT_DIR"

blog_pid=""
blog_port="$(port_from_addr "$BLOG_ADDR")"

require_free_port "$blog_port" "Blog"
require_free_port "$BROWSER_SYNC_PORT" "BrowserSync"

cleanup() {
  if [[ -n "$blog_pid" ]] && kill -0 "$blog_pid" 2>/dev/null; then
    kill "$blog_pid" 2>/dev/null || true
    wait "$blog_pid" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

ADDR="$BLOG_ADDR" go run ./cmd/blog &
blog_pid=$!

for _ in {1..50}; do
  if curl -sf "http://localhost${BLOG_ADDR}" >/dev/null; then
    break
  fi
  sleep 0.2
done

if ! curl -sf "http://localhost${BLOG_ADDR}" >/dev/null; then
  echo "Blog server did not become ready on http://localhost${BLOG_ADDR}" >&2
  exit 1
fi

exec npx --yes browser-sync start \
  --proxy "http://localhost${BLOG_ADDR}" \
  --port "$BROWSER_SYNC_PORT" \
  --files "content/**/*.md" \
  --files "internal/web/templates/**/*.html" \
  --no-open
