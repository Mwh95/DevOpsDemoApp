#!/bin/bash
# Deploy all services to local Kubernetes cluster

set -e

echo "Building Docker images..."
./gradlew buildAll

echo "Starting local database..."
cd Database
docker-compose up -d
cd ..

echo "Deploying to Kubernetes..."
kubectl apply -f Database/k8s/local/ 2>/dev/null || echo "Database k8s config not found, using docker-compose"
kubectl apply -f Keycloak/k8s/local/

echo "Setting up Ingress (controller + routing)..."
./scripts/setup-local-ingress.sh

echo ""
echo "Deployment complete!"
echo ""
echo "Keycloak: http://localhost:<nodeport>/auth  (get nodeport: kubectl get svc -n ingress-nginx ingress-nginx-controller)"
echo "Check status: kubectl get pods,svc && kubectl get ingress"
