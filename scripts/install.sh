#!/bin/sh
# AegisFlow installer: downloads and verifies a prebuilt release binary.
#
#   curl -fsSL https://raw.githubusercontent.com/saivedant169/AegisFlow/main/scripts/install.sh | sh
#
# Options (env vars):
#   AEGISFLOW_VERSION   release tag to install (default: latest)
#   AEGISFLOW_BIN_DIR   install directory (default: /usr/local/bin, falls back to ~/.local/bin)
#
# No Go toolchain required.

set -eu

REPO="saivedant169/AegisFlow"
BIN="aegisflow"
VERSION="${AEGISFLOW_VERSION:-latest}"

err() { echo "error: $*" >&2; exit 1; }
info() { echo "==> $*"; }

# --- detect OS / arch ---
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux) os=linux ;;
  darwin) os=darwin ;;
  *) err "unsupported OS: $os (only linux and darwin have prebuilt binaries)" ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) err "unsupported architecture: $arch" ;;
esac

asset="${BIN}-${os}-${arch}"

# --- resolve the download URL ---
if [ "$VERSION" = "latest" ]; then
  url="https://github.com/${REPO}/releases/latest/download/${asset}"
  checksums_url="https://github.com/${REPO}/releases/latest/download/SHA256SUMS"
else
  url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"
  checksums_url="https://github.com/${REPO}/releases/download/${VERSION}/SHA256SUMS"
fi

# --- pick an install dir we can write to ---
bindir="${AEGISFLOW_BIN_DIR:-/usr/local/bin}"
if [ ! -d "$bindir" ] || [ ! -w "$bindir" ]; then
  bindir="$HOME/.local/bin"
  mkdir -p "$bindir"
fi

tmpdir=$(mktemp -d)
binary_file="$tmpdir/$asset"
checksums_file="$tmpdir/SHA256SUMS"
trap 'rm -rf "$tmpdir"' EXIT

download() {
  source_url="$1"
  destination="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fSL "$source_url" -o "$destination" || err "download failed: $source_url"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$destination" "$source_url" || err "download failed: $source_url"
  else
    err "need curl or wget"
  fi
}

info "downloading ${asset} (${VERSION})"
download "$url" "$binary_file"
download "$checksums_url" "$checksums_file"

expected=$(awk -v name="$asset" '$2 == name { print $1 }' "$checksums_file")
[ -n "$expected" ] || err "checksum missing for $asset"

if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$binary_file" | awk '{ print $1 }')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$binary_file" | awk '{ print $1 }')
else
  err "need sha256sum or shasum"
fi
[ "$actual" = "$expected" ] || err "checksum verification failed for $asset"
info "checksum verified"

chmod +x "$binary_file"
mv "$binary_file" "$bindir/$BIN"
trap - EXIT
rm -rf "$tmpdir"

info "installed $bindir/$BIN"
if ! command -v "$BIN" >/dev/null 2>&1; then
  echo "note: $bindir is not on your PATH. Add it, for example:"
  echo "  export PATH=\"$bindir:\$PATH\""
fi
echo
echo "Next: grab a config and run it. The quickest governed demo (no API keys):"
echo "  git clone https://github.com/${REPO}.git && cd AegisFlow/starter-kit && ./install-pr-writer.sh"
echo "Or point a Claude client at the gateway: export ANTHROPIC_BASE_URL=http://localhost:8080"
