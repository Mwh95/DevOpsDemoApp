#!/bin/bash
# Install ingress-nginx controller (containerized) and apply DemoApp Ingress for local K8s.
# Run from repository root. Keycloak must already be deployed.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

INGRESS_NGINX_URL="https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/baremetal/deploy.yaml"

echo "Installing ingress-nginx controller (bare metal / NodePort)..."
kubectl apply -f "$INGRESS_NGINX_URL"

echo "Waiting for ingress-nginx controller to be ready..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s 2>/dev/null || {
  echo "Waiting for controller pods..."
  sleep 10
  kubectl get pods -n ingress-nginx
}

echo "Applying DemoApp Ingress (path /auth -> Keycloak)..."
kubectl apply -f k8s/local/ingress.yaml

echo ""
echo "Local Ingress is ready."
echo "Get the NodePort for HTTP: kubectl get svc -n ingress-nginx ingress-nginx-controller"
echo "On Docker Desktop / Minikube you may be able to use port 80 if the controller got a LoadBalancer."
