#!/bin/bash
# Configure GKE Ingress for the PlaygroundApp (Keycloak at /auth).
# Run from repository root. Requires gcloud and kubectl configured for your GKE cluster.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# Optional: set to reserve and use a global static IP
STATIC_IP_NAME="${GKE_INGRESS_STATIC_IP_NAME:-}"

echo "Enabling required APIs (if not already enabled)..."
if command -v gcloud &>/dev/null; then
  PROJECT_ID="${GCP_PROJECT_ID:-$(gcloud config get-value project 2>/dev/null)}"
  if [ -n "$PROJECT_ID" ]; then
    gcloud services enable container.googleapis.com compute.googleapis.com --project="$PROJECT_ID" 2>/dev/null || true
  fi
fi

if [ -n "$STATIC_IP_NAME" ]; then
  echo "Reserving global static IP: $STATIC_IP_NAME"
  gcloud compute addresses create "$STATIC_IP_NAME" --global 2>/dev/null || echo "Address may already exist."
  echo "Update k8s/gcp/ingress.yaml annotation: kubernetes.io/ingress.global-static-ip-name: $STATIC_IP_NAME"
fi

echo "Applying GKE Ingress manifest..."
kubectl apply -f k8s/gcp/ingress.yaml

echo ""
echo "GKE Ingress applied. Waiting for external address (this can take a few minutes)..."
echo "Check status: kubectl get ingress playground-ingress"
echo ""
echo "Once an ADDRESS is assigned:"
echo "  - Keycloak: http://<ADDRESS>/auth"
echo "  - For HTTPS: create a ManagedCertificate and add annotation networking.gke.io/managed-certificates to the Ingress (see docs/gke-ingress-console.md)"
echo "  - Point your DNS A record to the Ingress ADDRESS if using a custom domain."
