#!/bin/bash
# Deploy all services to GCP Kubernetes cluster

set -e

# Check if PROJECT_ID is set
if [ -z "$GCP_PROJECT_ID" ]; then
    echo "Error: GCP_PROJECT_ID environment variable is not set"
    echo "Usage: export GCP_PROJECT_ID=your-project-id && ./scripts/deploy-gcp.sh"
    exit 1
fi

echo "Using GCP Project: $GCP_PROJECT_ID"

# Build and push Docker images
echo "Building Docker images..."
./gradlew buildAll

echo "Tagging images for GCR..."
docker tag keycloak:1.0.0 gcr.io/$GCP_PROJECT_ID/keycloak:1.0.0
docker tag reverse-proxy:1.0.0 gcr.io/$GCP_PROJECT_ID/reverse-proxy:1.0.0

echo "Pushing images to GCR..."
docker push gcr.io/$GCP_PROJECT_ID/keycloak:1.0.0
docker push gcr.io/$GCP_PROJECT_ID/reverse-proxy:1.0.0

# Update deployment files with project ID
echo "Updating deployment files..."
sed -i.bak "s/PROJECT_ID/$GCP_PROJECT_ID/g" Keycloak/k8s/gcp/deployment.yaml
sed -i.bak "s/PROJECT_ID/$GCP_PROJECT_ID/g" ReverseProxy/k8s/gcp/deployment.yaml

# Deploy to GKE
echo "Deploying to GKE..."
kubectl apply -f Keycloak/k8s/gcp/
kubectl apply -f ReverseProxy/k8s/gcp/

# Restore original files
mv Keycloak/k8s/gcp/deployment.yaml.bak Keycloak/k8s/gcp/deployment.yaml
mv ReverseProxy/k8s/gcp/deployment.yaml.bak ReverseProxy/k8s/gcp/deployment.yaml

echo ""
echo "Deployment complete!"
echo ""
echo "Check status with: kubectl get pods,svc"
echo "Get external IP: kubectl get svc reverse-proxy"
