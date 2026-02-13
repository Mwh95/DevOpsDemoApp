# Quick Start Guide

Get the PlaygroundApp up and running in 5 minutes!

## Prerequisites Check

```bash
# Check Docker
docker --version
# Expected: Docker version 20.10+

# Check Docker Compose
docker-compose --version
# Expected: Docker Compose version 2.0+

# Check kubectl
kubectl version --client
# Expected: Client Version v1.24+

# Check Kubernetes cluster
kubectl cluster-info
# Should show cluster information
```

## 5-Minute Local Setup

### Option 1: Using Deployment Script (Recommended)

```bash
# Clone the repository (or navigate to your existing clone)
git clone https://github.com/Mwh95/PlaygroundApp.git
cd PlaygroundApp

# Deploy everything
./scripts/deploy-local.sh

# Wait for pods to be ready (1-2 minutes)
kubectl get pods -w
```

**Access the application:**
- Get NodePort: `kubectl get svc -n ingress-nginx ingress-nginx-controller`
- Keycloak Admin: http://localhost:NODEPORT/auth/admin (admin/admin; replace NODEPORT with the HTTP port)
- API endpoint: http://localhost:NODEPORT/auth

### Option 2: Step-by-Step

```bash
# 1. Build Docker images
./gradlew buildAll

# 2. Start database
./gradlew :Database:startLocalDb

# 3. Deploy to Kubernetes
kubectl apply -f Keycloak/k8s/local/deployment.yaml
./scripts/setup-local-ingress.sh

# 4. Wait for ready status
kubectl wait --for=condition=ready pod -l app=keycloak --timeout=120s

# 5. Access the application (use NodePort from: kubectl get svc -n ingress-nginx ingress-nginx-controller)
open http://localhost:NODEPORT/auth/admin
```

## Verify Installation

```bash
# Check all services are running
kubectl get pods,svc

# Expected output: keycloak pods and (in ingress-nginx namespace) ingress-nginx-controller

# Test the endpoints (replace NODEPORT with the HTTP port from ingress-nginx-controller service)
curl http://localhost:NODEPORT/auth

# View logs if needed
kubectl logs -l app=keycloak --tail=20
```

## Common Issues

### Pods not starting?
```bash
# Check pod status
kubectl describe pod -l app=keycloak

# Check logs
kubectl logs -l app=keycloak
```

### Database connection failed?
```bash
# Check database is running
docker ps | grep postgres

# Restart database
./gradlew :Database:stopLocalDb
./gradlew :Database:startLocalDb
```

### Port already in use?
```bash
# Ingress controller uses NodePort by default; get the port:
kubectl get svc -n ingress-nginx ingress-nginx-controller
```

## Clean Up

```bash
# Remove all resources
./scripts/cleanup-local.sh

# Or manually:
kubectl delete -f k8s/local/ingress.yaml
kubectl delete -f Keycloak/k8s/local/deployment.yaml
kubectl delete namespace ingress-nginx
./gradlew :Database:stopLocalDb
```

## Next Steps

- 📖 Read [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions
- 🏗️ Read [ARCHITECTURE.md](ARCHITECTURE.md) to understand the system architecture
- ☁️ Deploy to GCP following the [GCP Deployment Guide](DEPLOYMENT.md#gcp-deployment)
- 🔧 Customize Keycloak configuration and add realms
- 🔒 Configure SSL/TLS for production use

## Useful Commands

```bash
# Gradle tasks
./gradlew tasks                          # List all available tasks
./gradlew buildAll                       # Build all Docker images
./gradlew :Keycloak:buildDockerImage    # Build Keycloak image only

# Kubernetes operations
kubectl get all                          # View all resources
kubectl logs -l app=keycloak -f          # Follow Keycloak logs
kubectl exec -it <pod-name> -- bash      # Shell into a pod
kubectl port-forward svc/keycloak 8080:8080  # Port forward for testing

# Database operations
./gradlew :Database:startLocalDb         # Start database
./gradlew :Database:stopLocalDb          # Stop database
./gradlew :Database:logsDb              # View database logs
```

## Getting Help

- Check the [README.md](README.md) for overview
- Review [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment steps
- Read [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Open an issue on GitHub for bugs or questions

---

**🎉 Congratulations!** You now have a fully functional authentication system running on Kubernetes!
