#!/bin/bash
# Cleanup local deployment

set -e

echo "Removing Kubernetes resources..."
kubectl delete -f ReverseProxy/k8s/local/ --ignore-not-found=true
kubectl delete -f Keycloak/k8s/local/ --ignore-not-found=true

echo "Stopping local database..."
cd Database
docker-compose down -v
cd ..

echo "Cleanup complete!"
