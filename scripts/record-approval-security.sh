#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT="$ROOT/docs/assets/approval-security-e2e.gif"
CAST="$(mktemp "${TMPDIR:-/tmp}/aegisflow-approval-e2e.XXXXXX.cast")"
trap 'rm -f "$CAST"' EXIT

for command in asciinema agg go node curl jq; do
  if ! command -v "$command" >/dev/null 2>&1; then
    printf 'missing command: %s\n' "$command" >&2
    exit 1
  fi
done

cd "$ROOT"
asciinema rec \
  --headless \
  --return \
  --overwrite \
  --idle-time-limit 0.75 \
  --window-size 98x27 \
  --command "AEGISFLOW_DEMO=1 bash scripts/e2e_approval_security.sh" \
  "$CAST"

agg \
  --quiet \
  --theme github-dark \
  --font-size 18 \
  --idle-time-limit 0.75 \
  --last-frame-duration 2 \
  --cols 98 \
  --rows 27 \
  "$CAST" \
  "$OUTPUT"

printf 'wrote %s\n' "$OUTPUT"
