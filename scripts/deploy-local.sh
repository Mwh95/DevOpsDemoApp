#!/bin/bash
# Build local images and deploy the DemoApp stack with Helm.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

CURRENT_CONTEXT="$(kubectl config current-context 2>/dev/null || true)"
CHART_PATH="$REPO_ROOT/deploy/helm/demoapp"
ROLL_OUT_TOKEN="$(date +%s)"

require_command() {
  local command_name="$1"

  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Required command not found: $command_name"
    exit 1
  fi
}

load_local_image() {
  local image="$1"

  case "$CURRENT_CONTEXT" in
    minikube)
      minikube image load "$image"
      ;;
    *)
      echo "Current context '$CURRENT_CONTEXT' may require manual image loading for $image."
      ;;
  esac
}

require_command docker
require_command kubectl
require_command helm

echo "Building Keycloak images..."
cd Keycloak/
docker build -t keycloak:dev .
cd ..

echo "Building Map API image..."
cd MapService
docker build -t map-api:dev .
cd ..

echo "Building Map Frontend image..."
cd MapFrontend
docker build -t map-frontend:dev .
cd ..

echo "Building Liquibase image..."
cd Liquibase
docker build -t demoapp-liquibase:dev .
cd ..

echo "Loading local images into the Kubernetes cluster when needed..."
load_local_image keycloak:dev
load_local_image map-api:dev
load_local_image map-frontend:dev
load_local_image demoapp-liquibase:dev

echo "Installing ingress-nginx with Helm..."
./scripts/setup-local-ingress.sh

echo "Deploying DemoApp with Helm..."
helm upgrade --install demoapp "$CHART_PATH" \
  --set-string rolloutToken="$ROLL_OUT_TOKEN" \
  --wait \
  --wait-for-jobs \
  --timeout 5m

echo ""
echo "Deployment complete!"
echo ""
echo "Port-forward ingress for a stable local URL:"
echo "  kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 50594:80"
echo ""
echo "Then open:"
echo "  Map App:   http://localhost:50594/"
echo "  Keycloak:  http://localhost:50594/login"
echo ""
echo "Check status with:"
echo "  kubectl get pods,svc && kubectl get ingress"
