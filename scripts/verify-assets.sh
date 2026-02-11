#!/usr/bin/env bash
set -euo pipefail

# Verify that all expected release artifacts exist in dist/.
# Usage: ./scripts/verify-assets.sh [VERSION]

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

VERSION="${1:-$(cat "$ROOT_DIR/.beat-version" | tr -d '[:space:]')}"

DIST="$ROOT_DIR/dist"
ERRORS=0

expect() {
  if [[ ! -f "$DIST/$1" ]]; then
    echo "MISSING: $1"
    ERRORS=$((ERRORS + 1))
  else
    echo "     OK: $1"
  fi
}

echo "Verifying artifacts for safecastbeat v${VERSION}"
echo "---"

# tar.gz + sha512
for combo in linux-x86_64 linux-aarch64 darwin-x86_64 darwin-aarch64; do
  expect "safecastbeat-${VERSION}-${combo}.tar.gz"
  expect "safecastbeat-${VERSION}-${combo}.tar.gz.sha512"
done

# zip + sha512
expect "safecastbeat-${VERSION}-windows-x86_64.zip"
expect "safecastbeat-${VERSION}-windows-x86_64.zip.sha512"

# deb + sha512 (optional — only if nfpm was available)
for deb_arch in amd64 arm64; do
  if [[ -f "$DIST/safecastbeat_${VERSION}_${deb_arch}.deb" ]]; then
    expect "safecastbeat_${VERSION}_${deb_arch}.deb"
    expect "safecastbeat_${VERSION}_${deb_arch}.deb.sha512"
  fi
done

# rpm + sha512 (optional — only if nfpm was available)
for rpm_arch in x86_64 aarch64; do
  if [[ -f "$DIST/safecastbeat-${VERSION}-1.${rpm_arch}.rpm" ]]; then
    expect "safecastbeat-${VERSION}-1.${rpm_arch}.rpm"
    expect "safecastbeat-${VERSION}-1.${rpm_arch}.rpm.sha512"
  fi
done

echo "---"
if [[ "$ERRORS" -gt 0 ]]; then
  echo "FAIL: ${ERRORS} missing artifact(s)"
  exit 1
else
  echo "PASS: all artifacts present"
fi
