#!/bin/sh
set -e

CONFIG_FILE=/usr/share/nginx/html/env-config.js

cat <<EOF > "$CONFIG_FILE"
window.__ENV__ = {
  VITE_API_BASE: "${VITE_API_BASE:-}",
  VITE_OIDC_AUTHORITY: "${VITE_OIDC_AUTHORITY:-}",
  VITE_OIDC_CLIENT_ID: "${VITE_OIDC_CLIENT_ID:-map-app}",
  VITE_OIDC_REDIRECT_URI: "${VITE_OIDC_REDIRECT_URI:-}",
};
EOF

exec "$@"
