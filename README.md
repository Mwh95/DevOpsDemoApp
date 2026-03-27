# DemoApp

A demo application with Keycloak authentication, Kubernetes Ingress routing, PostgreSQL persistence, and a React map UI backed by a Go API.

You can label any point on a map and add additional comments. Labels are stored in the database and associated with a user id.

#### Screenshot
![Screenshot of the DemoApp](docs/img/DemoApp.png)

## Project Structure

```text
DemoApp/
├── ARCHITECTURE.md              # Architecture and deployment model
├── DEPLOYMENT.md                # Deployment guide and troubleshooting
├── Database/                    # Local PostgreSQL via Docker Compose
│   ├── docker-compose.yml
│   ├── k8s/local/
│   └── README.md
├── Keycloak/                    # Keycloak image and standalone manifests
│   ├── Dockerfile
│   ├── k8s/local/
│   ├── k8s/gcp/
│   └── README.md
├── Liquibase/                   # Database bootstrap and schema migrations
│   ├── Dockerfile
│   ├── modules/
│   ├── k8s/local/
│   └── README.md
├── MapFrontend/                 # React/Vite single-page app
│   ├── Dockerfile
│   ├── src/
│   └── k8s/local/
├── MapService/                  # Go REST API for marker data
│   ├── Dockerfile
│   ├── cmd/map-api/
│   ├── internal/
│   └── k8s/local/
├── deploy/helm/demoapp/         # Primary local Helm chart
├── docs/                        # Runbooks and screenshots
├── k8s/                         # Shared ingress manifests
└── scripts/                     # Local and GCP deployment helpers
```

## Documentation

- **[ARCHITECTURE.md](ARCHITECTURE.md)** — System design and components.
- **[DEPLOYMENT.md](DEPLOYMENT.md)** — Deployment guide and troubleshooting.
- **[Database/README.md](Database/README.md)** — Local PostgreSQL details.
- **[Keycloak/README.md](Keycloak/README.md)** — Keycloak notes.
- **[Liquibase/README.md](Liquibase/README.md)** — Migration details.

## Features

- **Keycloak authentication** with local ingress routing under `/login`
- **Separate frontend and API services** for the map application
- **Liquibase migrations** for bootstrap and schema management
- **Helm-based local deployment** for the Kubernetes workloads
- **PostgreSQL persistence** with dedicated schemas and users
- **Draft GCP path** for future cloud deployment work

## Prerequisites

### Local Development

- Docker
- Kubernetes (Docker Desktop, Minikube, or Kind)
- kubectl
- Helm

## Quick Start

### Local Deployment

1. Start a local Kubernetes cluster.

2. Start PostgreSQL from the repository root:

```bash
docker compose -f Database/docker-compose.yml up -d
```

3. Deploy the local stack:

```bash
./scripts/deploy-local.sh
```

If you need fresh images, rebuild them during deploy with:

```bash
./scripts/deploy-local.sh --build
```

The deploy script:

- verifies or rebuilds the local `keycloak:dev`, `map-api:dev`, `map-frontend:dev`, and `demoapp-liquibase:dev` images
- loads those images automatically when the current Kubernetes context is `minikube`
- installs `ingress-nginx` with Helm
- deploys the DemoApp Helm chart from `deploy/helm/demoapp`
- runs the Liquibase job before the application workloads start

4. In a separate shell, port-forward ingress on `50594`:

```bash
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 50594:80
```

5. Wait for pods to be ready, then access the services:

- Map app: `http://localhost:50594/`
- Map API: `http://localhost:50594/api`
- Keycloak: `http://localhost:50594/login`
- Admin Console: `http://127.0.0.1:50594/login/admin`

6. Check status:

```bash
kubectl get pods,svc
kubectl get ingress
kubectl get jobs
```

7. Clean up Kubernetes resources when finished:

```bash
./scripts/cleanup-local.sh
```

To stop the local database too:

```bash
docker compose -f Database/docker-compose.yml down
```

### DRAFT: GCP Deployment

The GCP path is still draft-quality. Today, `scripts/deploy-gcp.sh` only builds and deploys Keycloak plus the shared GKE ingress resources; it is not yet a full-stack Helm deployment like the local setup.

1. Set your GCP project and configure kubectl for GKE:

```bash
export GCP_PROJECT_ID=your-gcp-project-id
gcloud container clusters get-credentials your-cluster-name --region your-region
```

2. Update secrets in `Keycloak/k8s/gcp/deployment.yaml`:

- Replace placeholder passwords such as `CHANGE_ME_BEFORE_DEPLOYMENT`.
- Set the database connection and hostname values for your environment.

3. Deploy the current draft GCP resources:

```bash
./scripts/deploy-gcp.sh
```

4. Get the ingress address:

```bash
kubectl get ingress demoapp-ingress
```

Optional: run `./scripts/setup-gke-ingress.sh` for ingress setup with an optional static IP. For Console-based setup, see [docs/gke-ingress-console.md](docs/gke-ingress-console.md).

## Building Individual Components

### Keycloak

```bash
docker build -t keycloak:dev Keycloak/
```

### Map API

```bash
docker build -t map-api:dev MapService/
```

### Map Frontend

```bash
docker build -t map-frontend:dev MapFrontend/
```

### Liquibase

```bash
docker build -f Liquibase/Dockerfile -t demoapp-liquibase:dev .
```

### Database (Local)

```bash
docker compose -f Database/docker-compose.yml up -d
```

See [Database/README.md](Database/README.md) for details.

## Architecture

### Local Environment

- PostgreSQL runs locally via Docker Compose in `Database/docker-compose.yml`.
- Keycloak, Map Frontend, Map API, and Liquibase run in Kubernetes.
- Ingress is installed by `scripts/setup-local-ingress.sh`, while application routing is managed by the Helm chart in `deploy/helm/demoapp`.
- Local ingress routes `/` to `map-frontend`, `/api` to `map-api`, `/login` to Keycloak, and `/keycloak/health` plus `/keycloak/metrics` to Keycloak management endpoints.
- The default Helm values point workloads at `host.minikube.internal` for PostgreSQL access from the cluster.

### DRAFT: GCP Environment

- PostgreSQL on Cloud SQL is the intended target.
- Keycloak has draft GKE deployment manifests.
- Shared ingress resources live in `k8s/gcp/`.
- The current GCP path is not yet a full production-ready deployment for the whole stack.

## Configuration

### Keycloak

- Default admin: `tmpadmin/admin` (local)
- Database: PostgreSQL
- Ports: `8080` (HTTP), `9000` (management)
- Main route: `/login`
- Health checks: `/keycloak/health/ready` and `/keycloak/health/live`

### Map App

- Frontend route: `/`
- API route: `/api`
- OIDC authority and redirect URI are injected into the frontend container at runtime

### Database

- Engine: PostgreSQL
- Port: `5432`
- Default database: `MapMarkerDb`
- Bootstrap user: `postgres/postgres` (init only)
- Runtime users: `keycloak` for schema `keycloak`, `mapservice` for schema `mapservice`

## Customization

### Adding Routes

Edit `deploy/helm/demoapp/templates/ingress.yaml` for the local Helm deployment or `k8s/gcp/ingress.yaml` for the draft GCP ingress, then redeploy.

### Keycloak Configuration

Add realm export files to `Keycloak/` and update the Dockerfile:

```dockerfile
COPY realm-export.json /opt/keycloak/data/import/
```

### Database Initialization

Add bootstrap changelogs under `Liquibase/modules/database/bootstrap/` or service migrations under `Liquibase/modules/mapservice/changelog/`.

## Monitoring and Troubleshooting

### View Logs

```bash
# Keycloak
kubectl logs -l app=keycloak

# Map API
kubectl logs -l app=map-api

# Map Frontend
kubectl logs -l app=map-frontend

# Ingress controller (local)
kubectl logs -n ingress-nginx -l app.kubernetes.io/component=controller

# Database (local)
docker compose -f Database/docker-compose.yml logs
```

### Port Forwarding (for testing)

```bash
kubectl port-forward svc/keycloak 8080:8080
# Or use ingress-nginx:
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 50594:80
```

## Development

- Build Keycloak: `docker build -t keycloak:dev Keycloak/`
- Build Map API: `docker build -t map-api:dev MapService/`
- Build Map Frontend: `docker build -t map-frontend:dev MapFrontend/`
- Build Liquibase: `docker build -f Liquibase/Dockerfile -t demoapp-liquibase:dev .`
- Start local database: `docker compose -f Database/docker-compose.yml up -d`
- Deploy locally: `./scripts/deploy-local.sh`
- Rebuild and deploy locally: `./scripts/deploy-local.sh --build`
- Clean up Kubernetes resources: `./scripts/cleanup-local.sh`
- Helm chart: `deploy/helm/demoapp`

For more commands and troubleshooting, see [DEPLOYMENT.md](DEPLOYMENT.md).

## License

See [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request