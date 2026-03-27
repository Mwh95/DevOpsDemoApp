#!/bin/bash
# Cleanup local Helm releases.

set -euo pipefail

echo "Removing Kubernetes resources..."
helm uninstall demoapp --ignore-not-found=true || true
helm uninstall ingress-nginx -n ingress-nginx --ignore-not-found=true || true

# non critical error, namespace may already be gone
kubectl delete namespace ingress-nginx --ignore-not-found=true --timeout=10s 2>/dev/null || true

echo "Cleanup complete!"
