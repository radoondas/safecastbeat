#!/bin/sh
set -e

# Create safecastbeat system user if it does not exist
if ! getent passwd safecastbeat >/dev/null 2>&1; then
  useradd --system --no-create-home --shell /usr/sbin/nologin safecastbeat
fi

# Ensure config directory permissions
chown -R root:safecastbeat /etc/safecastbeat
chmod 750 /etc/safecastbeat

# Reload systemd
if command -v systemctl >/dev/null 2>&1; then
  systemctl daemon-reload
fi
