#!/bin/bash
# Install ingress-nginx for the local DemoApp cluster with Helm.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

if ! command -v helm >/dev/null 2>&1; then
  echo "Required command not found: helm"
  exit 1
fi

echo "Configuring the ingress-nginx Helm repository..."
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx >/dev/null 2>&1 || true
helm repo update ingress-nginx >/dev/null

echo "Installing ingress-nginx controller..."
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=NodePort \
  --wait \
  --timeout 5m

echo ""
echo "Local Ingress is ready."
echo "Port-forward ingress for a stable local URL:"
echo "  kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 50594:80"
