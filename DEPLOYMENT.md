# Deployment Guide

## Prerequisites

Before deploying the application, ensure you have the following installed:

### Required
- Docker (20.10+)
- Docker Compose (2.0+)
- kubectl (1.24+)
- Gradle 8.5+ (or use `./gradlew`)

### For Local Deployment
- Local Kubernetes cluster:
  - Docker Desktop with Kubernetes enabled, OR
  - Minikube, OR
  - Kind (Kubernetes in Docker)

### For GCP Deployment
- Google Cloud SDK (`gcloud` CLI)
- Active GCP project
- GKE cluster created
- Appropriate IAM permissions

---

## Local Deployment

### Step 1: Verify Prerequisites

```bash
# Check Docker
docker --version
docker-compose --version

# Check kubectl
kubectl version --client

# Check Kubernetes cluster
kubectl cluster-info
```

### Step 2: Build Docker Images

```bash
# Build all Docker images
./gradlew buildAll

# Verify images are created
docker images | grep keycloak
```

### Step 3: Start Database

The database will start automatically when you run the deploy script, but you can also start it manually:

```bash
cd Database
docker-compose up -d
docker-compose ps
cd ..
```

### Step 4: Deploy to Kubernetes

```bash
# Deploy all services
./scripts/deploy-local.sh

# Wait for pods to be ready (may take 1-2 minutes)
kubectl get pods -w
```

Expected output: keycloak pods and (in namespace ingress-nginx) ingress-nginx-controller pod(s).

### Step 5: Access Services

Get service information:
```bash
kubectl get svc
```

Access the application:
- Get NodePort: `kubectl get svc -n ingress-nginx ingress-nginx-controller`
- **Keycloak**: http://localhost:\<nodeport\>/auth
- **Keycloak Admin Console**: http://localhost:\<nodeport\>/auth/admin (replace \<nodeport\> with the HTTP port from the command above)
  - Username: `admin`
  - Password: `admin`

### Step 6: Verify Deployment

```bash
# Check all resources
kubectl get all

# View logs
kubectl logs -l app=keycloak
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Test the Ingress (use the HTTP NodePort from ingress-nginx-controller service)
curl http://localhost:\<nodeport\>/auth
```

### Cleanup Local Deployment

```bash
./scripts/cleanup-local.sh
```

---

## GCP Deployment

### Step 1: Setup GCP Environment

```bash
# Set your project ID
export GCP_PROJECT_ID=your-gcp-project-id

# Authenticate with GCP
gcloud auth login
gcloud config set project $GCP_PROJECT_ID

# Enable required APIs
gcloud services enable container.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable sqladmin.googleapis.com
```

### Step 2: Create GKE Cluster

```bash
# Create a GKE cluster (if not already created)
gcloud container clusters create playground-cluster \
  --zone us-central1-a \
  --num-nodes 2 \
  --machine-type n1-standard-2 \
  --enable-autoscaling \
  --min-nodes 2 \
  --max-nodes 4

# Get credentials
gcloud container clusters get-credentials playground-cluster \
  --zone us-central1-a
```

### Step 3: Create Cloud SQL Instance

```bash
# Create PostgreSQL instance
gcloud sql instances create playground-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1

# Create database
gcloud sql databases create keycloak --instance=playground-db

# Create user
gcloud sql users create keycloak \
  --instance=playground-db \
  --password=YOUR_SECURE_PASSWORD

# Get connection name
gcloud sql instances describe playground-db \
  --format="value(connectionName)"
```

### Step 4: Configure Secrets

Edit the deployment files to update secrets and configuration:

**Keycloak secrets** (`Keycloak/k8s/gcp/deployment.yaml`):
```yaml
stringData:
  admin-username: "admin"
  admin-password: "YOUR_SECURE_ADMIN_PASSWORD"  # CHANGE THIS
  db-username: "keycloak"
  db-password: "YOUR_SECURE_DB_PASSWORD"        # CHANGE THIS
```

**Database connection** (`Keycloak/k8s/gcp/deployment.yaml`):
```yaml
data:
  db-url: "jdbc:postgresql://CLOUD_SQL_CONNECTION_NAME/keycloak"
  hostname: "your-domain.com"  # Update with your actual domain
```

### Step 5: Configure Docker for GCR

```bash
# Configure Docker to use gcloud as credential helper
gcloud auth configure-docker
```

### Step 6: Deploy to GCP

```bash
# Deploy (builds Keycloak image, pushes to GCR, deploys Keycloak and GKE Ingress)
./scripts/deploy-gcp.sh

# Monitor deployment
kubectl get pods -w
```

To configure or create the Ingress via Google Cloud Console, see [docs/gke-ingress-console.md](docs/gke-ingress-console.md). Optional: run `./scripts/setup-gke-ingress.sh` to enable APIs and apply Ingress with optional static IP.

### Step 7: Get Ingress Address

```bash
# Get the Ingress external address
kubectl get ingress playground-ingress

# Wait for ADDRESS to be assigned (may take a few minutes)
# Then access your application at http://<ADDRESS>/auth
```

### Step 8: Configure DNS (Optional)

Point your domain to the Ingress address:
```bash
# Get Ingress address
kubectl get ingress playground-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}'

# Create DNS A record pointing your-domain.com to that IP
# See docs/gke-ingress-console.md for HTTPS and managed certificates
```

### Step 9: Enable HTTPS (Production)

For production, you should enable HTTPS:

1. Create a Google-managed SSL certificate:
```bash
gcloud compute ssl-certificates create playground-cert \
  --domains=your-domain.com
```

2. Update the service to use HTTPS and the certificate.

---

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name>

# Check logs
kubectl logs <pod-name>

# Check events
kubectl get events --sort-by='.lastTimestamp'
```

### Database Connection Issues

```bash
# For local deployment
docker-compose -f Database/docker-compose.yml logs

# Test database connectivity
kubectl run -it --rm debug --image=postgres:15-alpine --restart=Never -- \
  psql -h postgres -U keycloak -d keycloak
```

### Image Pull Issues (GCP)

```bash
# Verify image exists in GCR
gcloud container images list --repository=gcr.io/$GCP_PROJECT_ID

# Check service account permissions
kubectl describe serviceaccount default
```

### Port Forwarding for Debugging

```bash
# Forward Keycloak port directly
kubectl port-forward svc/keycloak 8080:8080

# Or use the Ingress NodePort (local): kubectl get svc -n ingress-nginx ingress-nginx-controller
```

---

## Scaling

### Scale Deployments

```bash
# Scale Keycloak
kubectl scale deployment keycloak --replicas=3

# Verify
kubectl get pods
```

### Auto-scaling (GCP)

```bash
# Enable horizontal pod autoscaling for Keycloak
kubectl autoscale deployment keycloak \
  --cpu-percent=70 \
  --min=2 \
  --max=10
```

---

## Monitoring

### View Resource Usage

```bash
# CPU and memory usage
kubectl top nodes
kubectl top pods
```

### View Logs

```bash
# Keycloak logs
kubectl logs -l app=keycloak --tail=100 -f

# Ingress controller logs (local)
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller --tail=100 -f
```

---

## Backup and Recovery

### Backup Database (Local)

```bash
cd Database
docker-compose exec postgres pg_dump -U keycloak keycloak > backup.sql
```

### Backup Database (GCP Cloud SQL)

```bash
gcloud sql backups create --instance=playground-db
```

### Restore Database

```bash
# Local
docker-compose exec -T postgres psql -U keycloak keycloak < backup.sql

# GCP
gcloud sql backups restore <BACKUP_ID> --backup-instance=playground-db
```

---

## Maintenance

### Update Images

```bash
# Build new images
./gradlew buildAll

# For GCP, push to registry
./scripts/deploy-gcp.sh

# Rolling update
kubectl rollout restart deployment keycloak

# Check rollout status
kubectl rollout status deployment keycloak
```

### Clean Up Resources

```bash
# Local
./scripts/cleanup-local.sh

# GCP - Delete cluster
gcloud container clusters delete playground-cluster --zone us-central1-a

# GCP - Delete Cloud SQL
gcloud sql instances delete playground-db
```

---

## Security Considerations

1. **Change default passwords** in production
2. **Use Kubernetes Secrets** for sensitive data
3. **Enable HTTPS** for all external endpoints
4. **Configure network policies** to restrict traffic
5. **Regular security updates** for base images
6. **Enable Cloud SQL SSL** for database connections
7. **Use IAM roles** and service accounts properly
8. **Enable audit logging** in GCP

---

## Additional Resources

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [GKE Ingress](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-configuration)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [GKE Documentation](https://cloud.google.com/kubernetes-engine/docs)
