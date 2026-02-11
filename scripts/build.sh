#!/usr/bin/env bash
set -euo pipefail

# Cross-compile safecastbeat for all supported platforms.
# Usage: ./scripts/build.sh [VERSION]
#   BEAT_VERSION env var takes priority, then positional arg, then .beat-version file.
# Output: dist/safecastbeat-{VERSION}-{os}-{arch}/safecastbeat[.exe]

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Resolve version
if [[ -n "${BEAT_VERSION:-}" ]]; then
  VERSION="$BEAT_VERSION"
elif [[ -n "${1:-}" ]]; then
  VERSION="$1"
else
  VERSION="$(cat "$ROOT_DIR/.beat-version" | tr -d '[:space:]')"
fi

echo "Building safecastbeat v${VERSION}"

MODULE="github.com/elastic/beats/v7/libbeat/version"
COMMIT="$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || echo unknown)"
NOW="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

LDFLAGS="-s -w -X ${MODULE}.buildTime=${NOW} -X ${MODULE}.commit=${COMMIT}"

PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%/*}"
  GOARCH="${platform#*/}"

  EXT=""
  if [[ "$GOOS" == "windows" ]]; then
    EXT=".exe"
  fi

  OUTDIR="$ROOT_DIR/dist/safecastbeat-${VERSION}-${GOOS}-${GOARCH}"
  mkdir -p "$OUTDIR"

  echo "  -> ${GOOS}/${GOARCH}"
  CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -ldflags "$LDFLAGS" -o "${OUTDIR}/safecastbeat${EXT}" "$ROOT_DIR"
done

echo "Done. Binaries in dist/"
