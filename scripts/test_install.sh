#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

mkdir -p "$TMP_DIR/bin" "$TMP_DIR/fixtures" "$TMP_DIR/install"

cat > "$TMP_DIR/fixtures/aegisflow-linux-amd64" <<'BINARY'
#!/bin/sh
echo "aegisflow test"
BINARY
chmod +x "$TMP_DIR/fixtures/aegisflow-linux-amd64"

checksum=$(sha256sum "$TMP_DIR/fixtures/aegisflow-linux-amd64" | awk '{ print $1 }')
printf '%s  aegisflow-linux-amd64\n' "$checksum" > "$TMP_DIR/fixtures/SHA256SUMS"

cat > "$TMP_DIR/bin/uname" <<'STUB'
#!/bin/sh
case "$1" in
  -s) echo Linux ;;
  -m) echo x86_64 ;;
  *) exit 1 ;;
esac
STUB

cat > "$TMP_DIR/bin/curl" <<'STUB'
#!/bin/sh
set -eu
destination=""
url=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    -o)
      destination="$2"
      shift 2
      ;;
    -*) shift ;;
    *)
      url="$1"
      shift
      ;;
  esac
done

case "$url" in
  */SHA256SUMS)
    if [ "${BAD_CHECKSUM:-0}" = "1" ]; then
      printf '%064d  aegisflow-linux-amd64\n' 0 > "$destination"
    else
      cp "$FIXTURE_DIR/SHA256SUMS" "$destination"
    fi
    ;;
  */aegisflow-linux-amd64)
    cp "$FIXTURE_DIR/aegisflow-linux-amd64" "$destination"
    ;;
  *) exit 1 ;;
esac
STUB

chmod +x "$TMP_DIR/bin/uname" "$TMP_DIR/bin/curl"

env \
  PATH="$TMP_DIR/bin:$PATH" \
  FIXTURE_DIR="$TMP_DIR/fixtures" \
  AEGISFLOW_BIN_DIR="$TMP_DIR/install" \
  sh "$ROOT_DIR/scripts/install.sh" > "$TMP_DIR/success.log"

test -x "$TMP_DIR/install/aegisflow"
grep -q "checksum verified" "$TMP_DIR/success.log"

rm -f "$TMP_DIR/install/aegisflow"
if env \
  PATH="$TMP_DIR/bin:$PATH" \
  FIXTURE_DIR="$TMP_DIR/fixtures" \
  BAD_CHECKSUM=1 \
  AEGISFLOW_BIN_DIR="$TMP_DIR/install" \
  sh "$ROOT_DIR/scripts/install.sh" > "$TMP_DIR/failure.log" 2>&1; then
  echo "installer accepted an invalid checksum" >&2
  exit 1
fi

grep -q "checksum verification failed" "$TMP_DIR/failure.log"
test ! -e "$TMP_DIR/install/aegisflow"

echo "installer verification tests passed"
