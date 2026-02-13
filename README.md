# PlaygroundApp

A multi-service application with Keycloak authentication, Apache reverse proxy, and PostgreSQL database, deployable on both local Kubernetes and GCP.

## Project Structure

```
PlaygroundApp/
├── Keycloak/           # Keycloak authentication service
│   ├── Dockerfile
│   ├── build.gradle
│   └── k8s/            # Kubernetes manifests
│       ├── local/
│       └── gcp/
├── ReverseProxy/       # Apache HTTPD reverse proxy
│   ├── Dockerfile
│   ├── build.gradle
│   ├── config/         # Apache configuration
│   └── k8s/            # Kubernetes manifests
│       ├── local/
│       └── gcp/
├── Database/           # PostgreSQL database
│   ├── docker-compose.yml
│   ├── build.gradle
│   └── config/         # Database initialization scripts
└── scripts/            # Deployment automation scripts
```

## Features

- **Keycloak**: Authentication and authorization service (2 replicas in K8s)
- **Reverse Proxy**: Apache HTTPD acting as reverse proxy (2 replicas in K8s)
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
   - Reverse Proxy: http://localhost:80
   - Keycloak: http://localhost:80/auth
   - Admin Console: http://localhost:80/auth/admin (admin/admin)

4. **Check status:**
   ```bash
   kubectl get pods,svc
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

5. **Get external IP:**
   ```bash
   kubectl get svc reverse-proxy
   ```

## Building Individual Components

### Keycloak
```bash
./gradlew :Keycloak:buildDockerImage
```

### Reverse Proxy
```bash
./gradlew :ReverseProxy:buildDockerImage
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
- Reverse Proxy routes traffic to Keycloak
- All services deployed to local Kubernetes cluster

### GCP Environment
- PostgreSQL on Cloud SQL (recommended)
- Keycloak and Reverse Proxy on GKE
- 2 replicas for high availability
- Load balancer for external access
- Secrets managed via Kubernetes Secrets
- Images stored in Google Container Registry

## Configuration

### Keycloak
- Default admin: `admin/admin` (local)
- Database: PostgreSQL
- Ports: 8080 (HTTP), 8443 (HTTPS)
- Health checks: `/health/ready`, `/health/live`

### Reverse Proxy
- Routes: `/auth` → Keycloak
- Ports: 80 (HTTP), 443 (HTTPS)
- Status endpoint: `/server-status`

### Database
- Engine: PostgreSQL 15
- Port: 5432
- Default database: `keycloak`
- Default user: `keycloak/keycloak` (⚠️ change for production)

## Customization

### Adding Routes to Reverse Proxy
Edit `ReverseProxy/config/httpd-vhosts.conf` and rebuild:
```bash
./gradlew :ReverseProxy:buildDockerImage
```

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

# Reverse Proxy
kubectl logs -l app=reverse-proxy

# Database (local)
./gradlew :Database:logsDb
```

### Port Forwarding (for testing)
```bash
kubectl port-forward svc/keycloak 8080:8080
kubectl port-forward svc/reverse-proxy 8000:80
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