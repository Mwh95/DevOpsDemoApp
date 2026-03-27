#!/bin/bash
# Optionally build local images and deploy the DemoApp stack with Helm.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

CURRENT_CONTEXT="$(kubectl config current-context 2>/dev/null || true)"
CHART_PATH="$REPO_ROOT/deploy/helm/demoapp"
ROLL_OUT_TOKEN="$(date +%s)"
BUILD_IMAGES=false

require_command() {
  local command_name="$1"

  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Required command not found: $command_name"
    exit 1
  fi
}

usage() {
  echo "Usage: ./scripts/deploy-local.sh [--build]"
  echo ""
  echo "  --build    Rebuild the local Docker images before deploying"
}

build_image() {
  local service_dir="$1"
  local image_name="$2"
  local label="$3"

  echo "Building $label image..."
  cd "$service_dir"
  docker build -t "$image_name" .
  cd "$REPO_ROOT"
}

ensure_local_image_exists() {
  local image_name="$1"
  local label="$2"

  if ! docker image inspect "$image_name" >/dev/null 2>&1; then
    echo "Local $label image not found: $image_name"
    echo "Re-run with --build to create the expected local images before deploying."
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

require_command kubectl
require_command helm

while [[ $# -gt 0 ]]; do
  case "$1" in
    --build)
      BUILD_IMAGES=true
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
  shift
done

if [[ "$BUILD_IMAGES" == true ]]; then
  require_command docker
  build_image "Keycloak" "keycloak:dev" "Keycloak"
  build_image "MapService" "map-api:dev" "Map API"
  build_image "MapFrontend" "map-frontend:dev" "Map Frontend"
  build_image "Liquibase" "demoapp-liquibase:dev" "Liquibase"
else
  require_command docker
  echo "Skipping image builds; pass --build to rebuild local images."
  ensure_local_image_exists "keycloak:dev" "Keycloak"
  ensure_local_image_exists "map-api:dev" "Map API"
  ensure_local_image_exists "map-frontend:dev" "Map Frontend"
  ensure_local_image_exists "demoapp-liquibase:dev" "Liquibase"
fi

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
  --timeout 2m

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
