# Keycloak Authentication Service

This module contains the Keycloak authentication service setup with Docker and Kubernetes configurations.

## Build Docker Image

```bash
./gradlew :Keycloak:buildDockerImage
```

## Deploy to Kubernetes

### Local Deployment
```bash
kubectl apply -f k8s/local/deployment.yaml
```

### GCP Deployment
1. Update the PROJECT_ID in `k8s/gcp/deployment.yaml`
2. Update secrets in `k8s/gcp/deployment.yaml`
3. Deploy:
```bash
kubectl apply -f k8s/gcp/deployment.yaml
```

## Configuration

- **Replicas**: 2 nodes for high availability
- **Database**: PostgreSQL
- **Ports**: 8080 (HTTP), 8443 (HTTPS)
- **Health checks**: Ready and liveness probes configured

## Environment Variables

### Local
- Admin credentials set to `admin:admin`
- Database connection to local postgres service

### GCP
- Admin credentials via Kubernetes secrets
- Database connection via Cloud SQL proxy
- HTTPS enforced with proper hostname
