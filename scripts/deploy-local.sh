#!/bin/bash
# Deploy all services to local Kubernetes cluster

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

CURRENT_CONTEXT="$(kubectl config current-context 2>/dev/null || true)"

load_local_image() {
  local image="$1"

  case "$CURRENT_CONTEXT" in
    kind-*)
      kind load docker-image "$image"
      ;;
    minikube)
      minikube image load "$image"
      ;;
    docker-desktop|rancher-desktop|orbstack)
      echo "Using local Docker image directly in context: $CURRENT_CONTEXT"
      ;;
    *)
      echo "Current context '$CURRENT_CONTEXT' may require manual image loading for $image."
      ;;
  esac
}

echo "Building Docker images..."
docker build -t keycloak:dev Keycloak/

echo "Building Map API image..."
docker build -t map-api:dev MapService/

echo "Building Map Frontend image..."
docker build -t map-frontend:dev MapFrontend/

echo "Building Liquibase image..."
docker build -f Liquibase/Dockerfile -t demoapp-liquibase:dev .

echo "Loading local images into the Kubernetes cluster when needed..."
load_local_image keycloak:dev
load_local_image map-api:dev
load_local_image map-frontend:dev
load_local_image demoapp-liquibase:dev

echo "Deploying to Kubernetes..."
kubectl apply -f Database/k8s/local/
kubectl rollout status deployment/postgres --timeout=180s
kubectl delete -f Liquibase/k8s/local/ --ignore-not-found=true
kubectl apply -f Liquibase/k8s/local/
kubectl wait --for=condition=complete job/liquibase --timeout=180s
kubectl apply -f Keycloak/k8s/local/
kubectl apply -f MapService/k8s/local/
kubectl apply -f MapFrontend/k8s/local/
kubectl rollout status deployment/keycloak --timeout=180s
kubectl rollout status deployment/map-api --timeout=180s
kubectl rollout status deployment/map-frontend --timeout=180s

echo "Setting up Ingress (controller + routing)..."
./scripts/setup-local-ingress.sh

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
