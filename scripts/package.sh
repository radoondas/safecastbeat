#!/usr/bin/env bash
set -euo pipefail

# Package safecastbeat binaries into distributable archives.
# Expects binaries already built in dist/ by build.sh.
# Usage: ./scripts/package.sh [VERSION]
#   BEAT_VERSION env var takes priority, then positional arg, then .beat-version file.
# Output: dist/safecastbeat-{VERSION}-{os}-{arch_name}.tar.gz|.zip + .sha512

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

echo "Packaging safecastbeat v${VERSION}"

# Common files to include in every archive
COMMON_FILES=(
  "safecastbeat.yml"
  "safecastbeat.reference.yml"
  "fields.yml"
  "CA.pem"
  "LICENSE.txt"
  "NOTICE.txt"
  "README.md"
)

# Map GOARCH to archive arch name
arch_name() {
  case "$1" in
    amd64) echo "x86_64" ;;
    arm64) echo "aarch64" ;;
    *)     echo "$1" ;;
  esac
}

# Create SHA-512 checksum sidecar
sha512() {
  local file="$1"
  if command -v sha512sum &>/dev/null; then
    (cd "$(dirname "$file")" && sha512sum "$(basename "$file")" > "$(basename "$file").sha512")
  else
    (cd "$(dirname "$file")" && shasum -a 512 "$(basename "$file")" > "$(basename "$file").sha512")
  fi
}

# Package tar.gz (linux, darwin)
package_tar() {
  local goos="$1" goarch="$2"
  local aname
  aname="$(arch_name "$goarch")"
  local builddir="$ROOT_DIR/dist/safecastbeat-${VERSION}-${goos}-${goarch}"
  local archive="$ROOT_DIR/dist/safecastbeat-${VERSION}-${goos}-${aname}.tar.gz"
  local staging="$ROOT_DIR/dist/_staging/safecastbeat-${VERSION}-${goos}-${aname}"

  if [[ ! -f "$builddir/safecastbeat" ]]; then
    echo "  SKIP ${goos}/${goarch} (binary not found)"
    return
  fi

  rm -rf "$staging"
  mkdir -p "$staging"

  cp "$builddir/safecastbeat" "$staging/"
  for f in "${COMMON_FILES[@]}"; do
    if [[ -f "$ROOT_DIR/$f" ]]; then
      cp "$ROOT_DIR/$f" "$staging/"
    fi
  done

  tar -czf "$archive" -C "$ROOT_DIR/dist/_staging" "safecastbeat-${VERSION}-${goos}-${aname}"
  sha512 "$archive"
  echo "  -> $(basename "$archive")"
}

# Package zip (windows)
package_zip() {
  local goarch="$1"
  local aname
  aname="$(arch_name "$goarch")"
  local builddir="$ROOT_DIR/dist/safecastbeat-${VERSION}-windows-${goarch}"
  local archive="$ROOT_DIR/dist/safecastbeat-${VERSION}-windows-${aname}.zip"
  local staging="$ROOT_DIR/dist/_staging/safecastbeat-${VERSION}-windows-${aname}"

  if [[ ! -f "$builddir/safecastbeat.exe" ]]; then
    echo "  SKIP windows/${goarch} (binary not found)"
    return
  fi

  rm -rf "$staging"
  mkdir -p "$staging"

  cp "$builddir/safecastbeat.exe" "$staging/"
  for f in "${COMMON_FILES[@]}"; do
    if [[ -f "$ROOT_DIR/$f" ]]; then
      cp "$ROOT_DIR/$f" "$staging/"
    fi
  done

  (cd "$ROOT_DIR/dist/_staging" && zip -qr "$archive" "safecastbeat-${VERSION}-windows-${aname}")
  sha512 "$archive"
  echo "  -> $(basename "$archive")"
}

# Clean staging
rm -rf "$ROOT_DIR/dist/_staging"

# tar.gz archives
for platform in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
  package_tar "${platform%/*}" "${platform#*/}"
done

# zip archive
package_zip amd64

# deb/rpm via nfpm (linux only)
if command -v nfpm &>/dev/null; then
  echo "Building deb/rpm packages..."
  for goarch in amd64 arm64; do
    builddir="$ROOT_DIR/dist/safecastbeat-${VERSION}-linux-${goarch}"
    if [[ ! -f "$builddir/safecastbeat" ]]; then
      echo "  SKIP nfpm linux/${goarch} (binary not found)"
      continue
    fi

    case "$goarch" in
      amd64) deb_arch="amd64"; rpm_arch="x86_64" ;;
      arm64) deb_arch="arm64"; rpm_arch="aarch64" ;;
    esac

    ARCH="$deb_arch" VERSION="$VERSION" GOARCH="$goarch" \
      nfpm package -p deb -f "$ROOT_DIR/packaging/nfpm.yaml" -t "$ROOT_DIR/dist/"
    sha512 "$ROOT_DIR/dist/safecastbeat_${VERSION}_${deb_arch}.deb"
    echo "  -> safecastbeat_${VERSION}_${deb_arch}.deb"

    ARCH="$rpm_arch" VERSION="$VERSION" GOARCH="$goarch" \
      nfpm package -p rpm -f "$ROOT_DIR/packaging/nfpm.yaml" -t "$ROOT_DIR/dist/"
    sha512 "$ROOT_DIR/dist/safecastbeat-${VERSION}-1.${rpm_arch}.rpm"
    echo "  -> safecastbeat-${VERSION}-1.${rpm_arch}.rpm"
  done
else
  echo "nfpm not found, skipping deb/rpm packages (install: go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.41.1)"
fi

# Clean staging
rm -rf "$ROOT_DIR/dist/_staging"

echo "Done. Archives in dist/"
