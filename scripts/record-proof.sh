#!/usr/bin/env bash
set -euo pipefail

MCP_URL="${AEGISFLOW_MCP_URL:-http://localhost:8082/mcp}"
ADMIN_URL="${AEGISFLOW_ADMIN_URL:-http://localhost:8081}"
API_KEY="${AEGISFLOW_API_KEY:-pr-writer-key-001}"

GREEN='\033[38;5;78m'
YELLOW='\033[38;5;221m'
RED='\033[38;5;203m'
BLUE='\033[38;5;75m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

post_tool() {
  local id="$1"
  local name="$2"
  local arguments="$3"
  jq -nc \
    --argjson id "$id" \
    --arg name "$name" \
    --argjson arguments "$arguments" \
    '{jsonrpc:"2.0",id:$id,method:"tools/call",params:{name:$name,arguments:$arguments}}' |
    curl -fsS -X POST "$MCP_URL" -H "Content-Type: application/json" --data-binary @-
}

step() {
  printf '\n%b%s%b\n' "$BOLD" "$1" "$RESET"
}

printf '\033[2J\033[H'
printf '%bAegisFlow: governed PR workflow%b\n' "$BOLD" "$RESET"
printf '%bLive MCP calls, local mock upstream, signed session evidence%b\n' "$DIM" "$RESET"
sleep 1

step "1. Agent reads repository metadata"
allow=$(post_tool 1 "github.list_repos" '{"owner":"saivedant169"}')
repo=$(jq -r '.result.content[0].text | fromjson | .[0].full_name' <<< "$allow")
printf '   %-34s %bALLOW%b\n' "github.list_repos" "$GREEN" "$RESET"
printf '   result: %s\n' "$repo"
sleep 2

step "2. Agent requests destructive deletion"
block=$(post_tool 2 "github.delete_repo" '{"repo":"saivedant169/AegisFlow"}' || true)
printf '   %-34s %bBLOCK%b\n' "github.delete_repo" "$RED" "$RESET"
printf '   error %s: %s\n' "$(jq -r '.error.code' <<< "$block")" "$(jq -r '.error.message' <<< "$block")"
sleep 2

step "3. Agent asks to open a pull request"
args='{"repo":"saivedant169/AegisFlow","title":"docs: clarify policy boundary","head":"docs/policy-boundary","base":"main"}'
review=$(post_tool 3 "github.create_pull_request" "$args" || true)
approval_id=$(jq -r '.error.data.approval_id // empty' <<< "$review")
if [[ -z "$approval_id" || "$approval_id" == "null" ]]; then
  printf 'No pending approval found\n' >&2
  exit 1
fi
printf '   %-34s %bREVIEW%b\n' "github.create_pull_request" "$YELLOW" "$RESET"
printf '   approval: %s\n' "${approval_id:0:12}"
sleep 2

step "4. Reviewer approves exact action"
approval=$(curl -fsS -X POST "$ADMIN_URL/admin/v1/approvals/$approval_id/approve" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"reviewer":"release-reviewer","comment":"diff scope checked"}')
printf '   reviewer: %-24s %b%s%b\n' \
  "$(jq -r '.reviewer' <<< "$approval")" "$GREEN" "$(jq -r '.status | ascii_upcase' <<< "$approval")" "$RESET"
sleep 2

step "5. Agent retries identical request"
resumed=$(post_tool 4 "github.create_pull_request" "$args")
pr_url=$(jq -r '.result.content[0].text | fromjson | .html_url' <<< "$resumed")
printf '   %-34s %bALLOW%b\n' "github.create_pull_request" "$GREEN" "$RESET"
printf '   mock result: %s\n' "$pr_url"
sleep 2

step "6. Verify signed evidence"
verify=$(curl -fsS -X POST "$ADMIN_URL/admin/v1/evidence/sessions/default/verify" \
  -H "X-API-Key: $API_KEY")
printf '   chain valid: %-8s records: %s\n' \
  "$(jq -r '.valid' <<< "$verify")" "$(jq -r '.total_records' <<< "$verify")"
printf '   %b%s%b\n' "$BLUE" "$(jq -r '.message' <<< "$verify")" "$RESET"
sleep 5
