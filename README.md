# DemoApp

A multi-service application with Keycloak authentication, Kubernetes Ingress routing, and PostgreSQL database, deployable on both local Kubernetes and GCP.

## Project Structure

```
DemoApp/
├── Keycloak/           # Keycloak authentication service
│   ├── Dockerfile
│   └── k8s/            # Kubernetes manifests
│       ├── local/
├── k8s/                # Ingress manifests (shared)
│   ├── local/          # Local Ingress (ingress-nginx)
├── Database/           # PostgreSQL (Docker Compose locally; Cloud SQL on GCP)
│   ├── docker-compose.yml
│   └── config/         # Database initialization scripts (e.g. init.sql)
├── docs/               # Runbooks (e.g. GKE Ingress via Console)
└── scripts/            # Deployment automation scripts
```

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** — Get running in 5 minutes (local).
- **[DEPLOYMENT.md](DEPLOYMENT.md)** — Full deployment guide (local).
- **[ARCHITECTURE.md](ARCHITECTURE.md)** — System design and components.

## Features

- **Keycloak**: Authentication and authorization service (2 replicas in K8s)
- **Ingress**: Kubernetes Ingress for routing (ingress-nginx locally)
- **Database**: PostgreSQL for data persistence
- **Multi-environment**: Supports both local development and GCP deployment
- **Container-first**: All services are containerized with Docker
- **Kubernetes Ready**: Complete K8s manifests for orchestration

## Prerequisites

### Local Development
- Docker and Docker Compose
- Kubernetes (Docker Desktop, Minikube, or Kind)
- kubectl

## Quick Start

### Local Deployment

1. Start a local Kubernetes cluster.

2. Deploy the local stack from the repository root:

```bash
./scripts/deploy-local.sh
```

3. Run port-forwarding to ingress on port `50594` in a separate shell:
```shell
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 50594:80
```

4. Wait for pods to be ready, then access the services:
   - Map app: `http://localhost:50594/`
   - Map API: `http://localhost:50594/api`
   - Keycloak: `http://localhost:50594/login`
   - Admin Console: `http://127.0.0.1:50594/login/admin`

5. Check status:
   ```bash
   kubectl get pods,svc && kubectl get ingress
   ```

6. Cleanup:
   ```bash
   ./scripts/cleanup-local.sh
   ```

For step-by-step instructions and troubleshooting, see [QUICKSTART.md](QUICKSTART.md).

### DRAFT: GCP Deployment

1. **Set your GCP project and configure kubectl for GKE:**
   ```bash
   export GCP_PROJECT_ID=your-gcp-project-id
   gcloud container clusters get-credentials your-cluster-name --region your-region
   ```

2. **Update secrets** in `Keycloak/k8s/gcp/deployment.yaml`:
   - Replace placeholder passwords (e.g. `CHANGE_ME_BEFORE_DEPLOYMENT`) with secure values.
   - Set database connection (e.g. Cloud SQL) and hostname as needed.  
   The deploy script will refuse to run if placeholders are still present.

3. **Deploy to GCP** (builds image, pushes to GCR, deploys Keycloak and GKE Ingress):
   ```bash
   ./scripts/deploy-gcp.sh
   ```

4. **Get Ingress address:**
   ```bash
   kubectl get ingress demoapp-ingress
   ```
   Optional: run `./scripts/setup-gke-ingress.sh` for APIs and Ingress with optional static IP. For Console-based setup, see [docs/gke-ingress-console.md](docs/gke-ingress-console.md).

For full GCP steps (cluster, Cloud SQL, secrets, HTTPS), see [DEPLOYMENT.md](DEPLOYMENT.md#gcp-deployment).

## Building Individual Components

### Keycloak
```bash
docker build -t keycloak:1.0.0 Keycloak/
```

### Database (Local)
```bash
cd Database && docker compose up -d
```
See [Database/README.md](Database/README.md) for details.

## Architecture

### Local Environment
- PostgreSQL runs via Docker Compose (`Database/docker-compose.yml`); not deployed to Kubernetes.
- Keycloak runs in Kubernetes and connects to the local PostgreSQL (host from cluster).
- Ingress (ingress-nginx) is installed by `scripts/setup-local-ingress.sh` and routes `/auth` to Keycloak.
- Access: http://localhost:\<nodeport\>/auth (nodeport from `kubectl get svc -n ingress-nginx ingress-nginx-controller`).

### DRAFT: GCP Environment
- PostgreSQL on Cloud SQL (recommended)
- Keycloak on GKE; GKE Ingress provides load balancer and routing
- 2 replicas for high availability
- Secrets managed via Kubernetes Secrets
- Images stored in Google Container Registry

## Configuration

### Keycloak
- Default admin: `tmpadmin/admin` (local)
- Database: PostgreSQL
- Ports: 8080 (HTTP), 8443 (HTTPS)
- Health checks (with context path `/login`): `/login/health/ready`, `/login/health/live`

### Ingress
- Local: ingress-nginx controller; path `/login` → Keycloak (NodePort).
- GCP: GKE Ingress; path `/login` → Keycloak. See `docs/gke-ingress-console.md` for Console setup.

### Database
- Engine: PostgreSQL 18.2
- Port: 5432
- Default database: `MapMarkerDb` (Keycloak and MapService each use a dedicated user and schema)
- Bootstrap user: `postgres/postgres` (init only); runtime: `keycloak` (schema keycloak), `mapservice` (schema mapservice)

## Customization

### Adding Routes
Edit `k8s/local/ingress.yaml` or `k8s/gcp/ingress.yaml` to add path rules and backend services, then re-apply.

### Keycloak Configuration
Add realm export files to `Keycloak/` and update Dockerfile:
```dockerfile
COPY realm-export.json /opt/keycloak/data/import/
```

### Database Initialization
Add bootstrap changelogs to `Liquibase/modules/database/bootstrap/changelog.xml`

## Monitoring and Troubleshooting

### View logs
```bash
# Keycloak
kubectl logs -l app=keycloak

# Ingress controller (local)
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Database (local)
docker logs <container-id>
```

### Port Forwarding (for testing)
```bash
kubectl port-forward svc/keycloak 8080:8080
# Or use the Ingress NodePort: kubectl get svc -n ingress-nginx ingress-nginx-controller
```

## Development

- Build Keycloak: `docker build -t keycloak:1.0.0 Keycloak/`
- Build Map API: `docker build -t map-api:dev MapService/`
- Build Map Frontend: `docker build -t map-frontend:dev MapFrontend/`
- Build Liquibase: `docker build -f Liquibase/Dockerfile -t demoapp-liquibase:dev .`
- Deploy locally: `./scripts/deploy-local.sh`; cleanup: `./scripts/cleanup-local.sh`

For more commands and troubleshooting, see [QUICKSTART.md](QUICKSTART.md) and [DEPLOYMENT.md](DEPLOYMENT.md).

## License

See [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request