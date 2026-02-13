#!/bin/bash
# Deploy all services to GCP Kubernetes cluster

set -e

# Check if PROJECT_ID is set
if [ -z "$GCP_PROJECT_ID" ]; then
    echo "Error: GCP_PROJECT_ID environment variable is not set"
    echo "Usage: export GCP_PROJECT_ID=your-project-id && ./scripts/deploy-gcp.sh"
    exit 1
fi

# Validate that secrets have been changed from placeholders
echo "Validating deployment configuration..."
if grep -q "CHANGE_ME_BEFORE_DEPLOYMENT" Keycloak/k8s/gcp/deployment.yaml; then
    echo "❌ Error: Placeholder passwords found in Keycloak/k8s/gcp/deployment.yaml"
    echo "Please update the admin-password and db-password before deployment."
    echo "Edit Keycloak/k8s/gcp/deployment.yaml and replace CHANGE_ME_BEFORE_DEPLOYMENT with secure passwords."
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

# Deploy to GKE using kubectl with environment variable substitution
echo "Deploying to GKE..."

# Use temporary files to avoid modifying source files
echo "Creating temporary deployment files..."
export PROJECT_ID=$GCP_PROJECT_ID
envsubst < Keycloak/k8s/gcp/deployment.yaml > /tmp/keycloak-deployment.yaml || {
    # Fallback if envsubst is not available
    sed "s/PROJECT_ID/$GCP_PROJECT_ID/g" Keycloak/k8s/gcp/deployment.yaml > /tmp/keycloak-deployment.yaml
}

envsubst < ReverseProxy/k8s/gcp/deployment.yaml > /tmp/reverse-proxy-deployment.yaml || {
    # Fallback if envsubst is not available
    sed "s/PROJECT_ID/$GCP_PROJECT_ID/g" ReverseProxy/k8s/gcp/deployment.yaml > /tmp/reverse-proxy-deployment.yaml
}

# Apply the temporary files
kubectl apply -f /tmp/keycloak-deployment.yaml
kubectl apply -f /tmp/reverse-proxy-deployment.yaml

# Clean up temporary files
rm -f /tmp/keycloak-deployment.yaml /tmp/reverse-proxy-deployment.yaml

echo ""
echo "Deployment complete!"
echo ""
echo "Check status with: kubectl get pods,svc"
echo "Get external IP: kubectl get svc reverse-proxy"
