#!/bin/sh
set -e

# Stop the service before removal
if command -v systemctl >/dev/null 2>&1; then
  systemctl stop safecastbeat 2>/dev/null || true
  systemctl disable safecastbeat 2>/dev/null || true
fi
