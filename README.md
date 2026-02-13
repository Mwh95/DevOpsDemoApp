# PlaygroundApp

A multi-service application with Keycloak authentication, Kubernetes Ingress routing, and PostgreSQL database, deployable on both local Kubernetes and GCP.

## Project Structure

```
PlaygroundApp/
├── Keycloak/           # Keycloak authentication service
│   ├── Dockerfile
│   ├── build.gradle
│   └── k8s/            # Kubernetes manifests
│       ├── local/
│       └── gcp/
├── k8s/                # Ingress manifests
│   ├── local/          # Local Ingress (ingress-nginx)
│   └── gcp/            # GKE Ingress
├── Database/           # PostgreSQL database
│   ├── docker-compose.yml
│   ├── build.gradle
│   └── config/         # Database initialization scripts
├── docs/               # Runbooks (e.g. GKE Ingress via Console)
└── scripts/            # Deployment automation scripts
```

## Features

- **Keycloak**: Authentication and authorization service (2 replicas in K8s)
- **Ingress**: Kubernetes Ingress for routing (ingress-nginx locally, GKE Ingress on GCP)
- **Database**: PostgreSQL for data persistence
- **Multi-environment**: Supports both local development and GCP deployment
- **Gradle Integration**: All projects use Gradle for build automation
- **Container-first**: All services are containerized with Docker
- **Kubernetes Ready**: Complete K8s manifests for orchestration

## Prerequisites

### Local Development
- Docker and Docker Compose
- Kubernetes (Docker Desktop, Minikube, or Kind)
- kubectl
- Gradle 8.5+ (or use included wrapper: `./gradlew`)

### GCP Deployment
- Google Cloud SDK (`gcloud`)
- GKE cluster
- Docker registry access (GCR)
- kubectl configured for GKE

## Quick Start

### Local Deployment

1. **Build all projects:**
   ```bash
   ./gradlew buildAll
   ```

2. **Deploy to local Kubernetes:**
   ```bash
   ./scripts/deploy-local.sh
   ```

3. **Access the services:**
   - Get the NodePort: `kubectl get svc -n ingress-nginx ingress-nginx-controller`
   - Keycloak: http://localhost:NODEPORT/auth (use the HTTP NodePort from the command above)
   - Admin Console: http://localhost:NODEPORT/auth/admin (admin/admin)

4. **Check status:**
   ```bash
   kubectl get pods,svc && kubectl get ingress
   ```

5. **Cleanup:**
   ```bash
   ./scripts/cleanup-local.sh
   ```

### GCP Deployment

1. **Set your GCP project:**
   ```bash
   export GCP_PROJECT_ID=your-gcp-project-id
   ```

2. **Configure kubectl for GKE:**
   ```bash
   gcloud container clusters get-credentials your-cluster-name --region your-region
   ```

3. **Update secrets in deployment files:**
   - Edit `Keycloak/k8s/gcp/deployment.yaml`
   - Update passwords and database connection strings

4. **Deploy to GCP:**
   ```bash
   ./scripts/deploy-gcp.sh
   ```

5. **Get Ingress address:**
   ```bash
   kubectl get ingress playground-ingress
   ```

## Building Individual Components

### Keycloak
```bash
./gradlew :Keycloak:buildDockerImage
```

### Database (Local)
```bash
./gradlew :Database:startLocalDb
./gradlew :Database:stopLocalDb
```

## Architecture

### Local Environment
- PostgreSQL runs as Docker container (docker-compose)
- Keycloak connects to local PostgreSQL
- Ingress (ingress-nginx) routes traffic to Keycloak at /auth
- All services deployed to local Kubernetes cluster

### GCP Environment
- PostgreSQL on Cloud SQL (recommended)
- Keycloak on GKE; GKE Ingress provides load balancer and routing
- 2 replicas for high availability
- Secrets managed via Kubernetes Secrets
- Images stored in Google Container Registry

## Configuration

### Keycloak
- Default admin: `admin/admin` (local)
- Database: PostgreSQL
- Ports: 8080 (HTTP), 8443 (HTTPS)
- Health checks: `/health/ready`, `/health/live`

### Ingress
- Local: ingress-nginx controller; path `/auth` → Keycloak (NodePort).
- GCP: GKE Ingress; path `/auth` → Keycloak. See `docs/gke-ingress-console.md` for Console setup.

### Database
- Engine: PostgreSQL 15
- Port: 5432
- Default database: `keycloak`
- Default user: `keycloak/keycloak` (⚠️ change for production)

## Customization

### Adding Routes
Edit `k8s/local/ingress.yaml` or `k8s/gcp/ingress.yaml` to add path rules and backend services, then re-apply.

### Keycloak Configuration
Add realm export files to `Keycloak/` and update Dockerfile:
```dockerfile
COPY realm-export.json /opt/keycloak/data/import/
```

### Database Initialization
Add SQL scripts to `Database/config/init.sql`

## Monitoring and Troubleshooting

### View logs
```bash
# Keycloak
kubectl logs -l app=keycloak

# Ingress controller (local)
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Database (local)
./gradlew :Database:logsDb
```

### Port Forwarding (for testing)
```bash
kubectl port-forward svc/keycloak 8080:8080
# Or use the Ingress NodePort: kubectl get svc -n ingress-nginx ingress-nginx-controller
```

## Development

### Project Tasks
```bash
# List all tasks
./gradlew tasks

# Build all projects
./gradlew buildAll

# Clean all projects
./gradlew cleanAll
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request