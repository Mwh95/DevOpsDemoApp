#!/bin/bash
# Cleanup local deployment

set -e

echo "Removing Kubernetes resources..."
kubectl delete -f k8s/local/ingress.yaml --ignore-not-found=true
kubectl delete -f MapFrontend/k8s/local/ --ignore-not-found=true
kubectl delete -f MapService/k8s/local/ --ignore-not-found=true
kubectl delete -f Liquibase/k8s/local/ --ignore-not-found=true
kubectl delete -f Keycloak/k8s/local/ --ignore-not-found=true
kubectl delete -f Database/k8s/local/ --ignore-not-found=true
kubectl delete namespace ingress-nginx --ignore-not-found=true --timeout=60s 2>/dev/null || true

echo "Cleanup complete!"
