# Apache HTTPD Reverse Proxy

This module contains the Apache HTTPD reverse proxy setup with Docker and Kubernetes configurations.

## Build Docker Image

```bash
./gradlew :ReverseProxy:buildDockerImage
```

## Deploy to Kubernetes

### Local Deployment
```bash
kubectl apply -f k8s/local/deployment.yaml
```

### GCP Deployment
1. Update the PROJECT_ID in `k8s/gcp/deployment.yaml`
2. Deploy:
```bash
kubectl apply -f k8s/gcp/deployment.yaml
```

## Configuration

- **Replicas**: 2 nodes for high availability
- **Ports**: 80 (HTTP), 443 (HTTPS)
- **Health checks**: Ready and liveness probes on /server-status endpoint
- **Proxy routes**: 
  - `/auth` → Keycloak service

## Customization

Edit `config/httpd-vhosts.conf` to add more proxy routes or configure SSL/TLS.
