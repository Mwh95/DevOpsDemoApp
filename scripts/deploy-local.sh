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
kubectl apply -f Database/k8s/local/ || echo "Database k8s config not found, using docker-compose"
kubectl apply -f Keycloak/k8s/local/
kubectl apply -f ReverseProxy/k8s/local/

echo ""
echo "Deployment complete!"
echo ""
echo "Services:"
echo "  - Reverse Proxy: http://localhost:80"
echo "  - Keycloak: http://localhost:80/auth"
echo ""
echo "Check status with: kubectl get pods,svc"
