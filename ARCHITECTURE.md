# Architecture Documentation

## Overview

This application is a multi-service platform consisting of:

1. **Keycloak** - Authentication and authorization service (OIDC)
2. **Map Service (Map Markers)** - Frontend (SPA) and backend (Go REST API) for the map markers app; frontend is embedded in the backend image and served by the same process
3. **Ingress** - Kubernetes Ingress for routing (ingress-nginx locally, GKE Ingress on GCP)
4. **Database** - PostgreSQL (MapMarkerDb) for Keycloak and Map Service data

All components are containerized and can be deployed on both local Kubernetes clusters and Google Cloud Platform (GCP).

## Target Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Internet                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ HTTP/HTTPS
                         ▼
          ┌──────────────────────────────┐
          │  Ingress (Local: NodePort /   │
          │  GCP: GKE Ingress + LB)      │
          └──────────┬──────────┬─────────┘
                     │          │
         /login      │          │  /  and  /api
         /keycloak/* │          │
                     ▼          ▼
    ┌─────────────────┐   ┌─────────────────────────────┐
    │ Keycloak Service │   │   Map API (Map Markers)     │
    │  2 Replicas     │   │   Backend (Go) + Frontend   │
    │  8080, 9000     │   │   (embedded SPA), port 8090 │
    └────────┬────────┘   └──────────────┬──────────────┘
             │                           │
             │ JDBC                      │ JDBC
             │ keycloak schema           │ mapservice schema
             ▼                           ▼
          ┌──────────────────────────────────────────────┐
          │         PostgreSQL (MapMarkerDb)             │
          │   (Docker Compose / Cloud SQL), Port: 5432   │
          └──────────────────────────────────────────────┘
```

## Component Details

### 1. Keycloak

**Purpose**: Identity and Access Management (IAM)

**Technology**: 
- Base Image: `quay.io/keycloak/keycloak:23.0`
- Runtime: Java
- Database: PostgreSQL

**Configuration**:
- **Development Mode** (`start-dev`): Used for local deployment
  - Simplified setup
  - Built-in H2 database fallback
  - Relaxed hostname checks
  
- **Production Mode** (`start`): Used for GCP deployment
  - Requires external database
  - Strict HTTPS enforcement
  - Hostname verification enabled

**Kubernetes Resources**:
- **Deployment**: 2 replicas for high availability
- **Service**: ClusterIP (internal) on ports 8080, 8443
- **ConfigMap**: Database connection strings, hostname configuration
- **Secret**: Admin credentials, database passwords

**Health Checks** (context path `/login`; management at `/keycloak`):
- Readiness probe: `/login/health/ready` (or `/keycloak/health/ready` in local)
- Liveness probe: `/login/health/live` (or `/keycloak/health/live` in local)

**Resource Allocation** (GCP):
- Requests: 512Mi memory, 500m CPU
- Limits: 1Gi memory, 1000m CPU

### 2. Ingress

**Purpose**: Path-based routing, single entry point, and (on GCP) SSL termination

**Local (ingress-nginx)**:
- Controller: [ingress-nginx](https://kubernetes.github.io/ingress-nginx/) (containerized, runs in-cluster)
- Installed via `scripts/setup-local-ingress.sh` (applies official bare-metal manifest)
- Service type: NodePort (ports 80/443 mapped to node ports)

**GCP (GKE Ingress)**:
- GKE Ingress resource creates a Google Cloud HTTP(S) Load Balancer
- Configured via `k8s/gcp/ingress.yaml` and optionally `scripts/setup-gke-ingress.sh`
- Console runbook: `docs/gke-ingress-console.md`

**Routing** (local ingress):
- `/login` → Keycloak Service (port 8080)
- `/keycloak/health` → Keycloak (port 9000); `/keycloak/metrics` → Keycloak (port 9000, optional IP allowlist)
- `/` → Map Frontend (port 80) — serves the Map Markers SPA
- `/api` → Map API (port 8090) — Map Markers REST API

Keycloak is configured with `KC_HTTP_RELATIVE_PATH=/login` (and management at `/keycloak`).

### 3. Map Service (Map Markers)

**Purpose**: Map Markers application — interactive map with markers; users sign in via Keycloak.

**Components**:
- **Frontend**: SPA (`MapFrontend/`, image `map-frontend:dev`) served separately at `/`.
- **Backend**: Go REST API (`MapService/`, image `map-api:dev`) — CRUD for markers, JWT validation via Keycloak JWKS, Liquibase for DB migrations; serves API at `/api`.

**Technology**:
- Frontend: Vite/React; Docker build: `docker build -t map-frontend:dev MapFrontend/`.
- Backend: Go; Docker build: `docker build -t map-api:dev MapService/`.
- Database: PostgreSQL (MapMarkerDb, schema `mapservice`).

**Kubernetes Resources**:
- Frontend: Deployment (`map-frontend`), Service (port 80).
- Backend: Deployment (`map-api`), Service (port 8090), ConfigMap (PG_*, KEYCLOAK_ISSUER), Secret (DB credentials).
- Migrations: standalone `Liquibase` Job that bootstraps users/schemas and applies the `MapService` changelog before the app workloads start.

**Authentication**: API accepts JWTs from Keycloak (realm/client used by the frontend, e.g. client `map-app`). No separate backend client; frontend uses OIDC with Keycloak.

### 4. Database (PostgreSQL)

**Purpose**: Persistent data storage for Keycloak and Map Service

**Technology**: 
- Base Image: `postgres:15-alpine` (or 18)
- Database: **MapMarkerDb** — single database with separate schemas/users for Keycloak and Map Service

**Local Deployment**:
- Runs as Docker container via docker-compose (`Database/docker-compose.yml`)
- Data persisted in named volume `postgres-data`
- Exposed on host port 5432
- Bootstrap user: `postgres`; runtime users: `keycloak` (schema keycloak), `mapservice` (schema mapservice) (⚠️ change for production)

**GCP Deployment**:
- Google Cloud SQL for PostgreSQL
- Managed service with automatic backups
- Accessed via Cloud SQL Proxy or private IP

**Configuration**:
- Database name: `MapMarkerDb`
- Port: 5432
- Initialization: `Liquibase/modules/database/bootstrap/changelog.xml` (creates users and schemas)

## Deployment Environments

### Local Environment

**Characteristics**:
- Development-focused
- Simplified configuration
- All services in local Kubernetes
- Database via docker-compose
- No TLS requirements
- Default credentials

**Network Architecture**:
```
localhost:<nodeport> → ingress-nginx → Ingress
                         ├─ /login, /keycloak/* → keycloak:8080/9000 → MapMarkerDb (keycloak schema)
                         └─ /, /api → map-api:8090 → MapMarkerDb (mapservice schema)
```

**Access** (get nodeport: `kubectl get svc -n ingress-nginx ingress-nginx-controller`):
- Map app (frontend): `http://localhost:<nodeport>/`
- Map API: `http://localhost:<nodeport>/api`
- Keycloak login: `http://localhost:<nodeport>/login`

### GCP Environment

**Characteristics**:
- Production-ready
- Security-hardened
- Secrets management
- TLS/HTTPS enforced
- Cloud-native services (Cloud SQL)
- Scalability and monitoring

**Network Architecture**:
```
Internet → GCP Load Balancer (from Ingress) → keycloak:8080 → Cloud SQL
```

**Access**:
- Map app: `http://<INGRESS-ADDRESS>/` or `https://your-domain.com/`
- Map API: `http://<INGRESS-ADDRESS>/api` or `https://your-domain.com/api`
- Keycloak: `http://<INGRESS-ADDRESS>/login` or `https://your-domain.com/login`

**Required GCP Services**:
- Google Kubernetes Engine (GKE)
- Google Container Registry (GCR)
- Cloud SQL for PostgreSQL
- Cloud Load Balancing

## Build and Deployment Pipeline

### Build Structure

```
DemoApp (root)
├── k8s/                    # Ingress manifests (shared)
│   ├── local/              # Local Ingress + controller (script applies controller from URL)
│   └── gcp/                # GKE Ingress
├── Keycloak/
│   ├── Dockerfile          # Keycloak image
│   └── k8s/                # Kubernetes manifests (local/, gcp/)
├── MapService/             # Map Markers backend (Go API)
│   ├── Dockerfile          # Build from repo root to embed MapFrontend
│   ├── db/changelog/       # Liquibase migrations
│   └── k8s/local/          # Deployment, Service, ConfigMap, Secret
├── MapFrontend/             # Map Markers SPA (embedded into Map API image)
├── Database/
│   ├── docker-compose.yml  # Local PostgreSQL (MapMarkerDb)
│   └── config/             # Init scripts (e.g. init.sql)
├── scripts/                # Deployment automation
│   ├── deploy-local.sh     # Deploy local images to K8s; add --build to rebuild first
│   ├── deploy-gcp.sh       # Build Keycloak, push to GCR, deploy to GKE
│   ├── setup-local-ingress.sh
│   ├── setup-gke-ingress.sh
│   └── cleanup-local.sh
└── docs/                   # Runbooks (e.g. GKE Ingress via Console)
```

**Build and run commands**:
- **Keycloak image**: `docker build -t keycloak:1.0.0 Keycloak/`
- **Map API image**: `docker build -t map-api:dev MapService/`
- **Map Frontend image**: `docker build -t map-frontend:dev MapFrontend/`
- **Liquibase image**: `docker build -f Liquibase/Dockerfile -t demoapp-liquibase:dev .`
- **Full local deploy**: `./scripts/deploy-local.sh` (deploys existing local images with Helm) or `./scripts/deploy-local.sh --build` (rebuilds Keycloak, map-api, map-frontend, and Liquibase first)
- **Full GCP deploy**: `./scripts/deploy-gcp.sh` (requires `GCP_PROJECT_ID` and updated secrets)

### Deployment Flow

#### Local Deployment
1. Optionally rebuild Docker images with `./scripts/deploy-local.sh --build`; otherwise `./scripts/deploy-local.sh` reuses existing local images.
2. Load the local images into Minikube when needed, install `ingress-nginx` with Helm, and run the DemoApp Helm chart.
3. Run the standalone Liquibase job to bootstrap users/schemas and apply the `MapService` schema changes before the rest of the application rolls out.
4. Services are accessible after port-forwarding ingress-nginx on `50594`: Map app at `/`, API at `/api`, Keycloak at `/login`.

#### GCP Deployment
1. Build Keycloak image; optionally build and push Map API image.
2. Tag and push images to Google Container Registry (GCR).
3. Apply Kubernetes manifests from `Keycloak/k8s/gcp/`, `MapService/k8s/` (if present), and `k8s/gcp/`.
4. Configure Cloud SQL connection (MapMarkerDb) and secrets.
5. Services accessible via external IP (GKE Ingress).

## Security Considerations

### Local Development
- ⚠️ Uses default credentials (acceptable for development)
- No TLS encryption
- Relaxed security policies
- Services exposed on localhost

### Production (GCP)
- ✅ Secrets stored in Kubernetes Secrets
- ✅ TLS/HTTPS enforced
- ✅ Database with strong passwords
- ✅ Network policies (should be added)
- ✅ Resource limits enforced
- ✅ Health checks configured
- ⚠️ Consider: Web Application Firewall (WAF)
- ⚠️ Consider: DDoS protection

### Recommended Security Enhancements
1. Enable Kubernetes Network Policies
2. Use cert-manager for automatic TLS certificates
3. Implement Pod Security Standards
4. Enable audit logging
5. Use Workload Identity for GCP service access
6. Implement rate limiting
7. Regular security scanning of container images

## Scaling Strategy

### Horizontal Scaling

**Keycloak**:
- Replicas: 2 (default) → 10 (max)
- Scaling trigger: CPU > 70%
- Session affinity: Not required (clustered mode)

**Ingress**:
- Local: ingress-nginx controller scales with cluster; Ingress resource is declarative.
- GCP: GKE Ingress and load balancer scale automatically.

**Database**:
- Local: Single container (not for production)
- GCP: Cloud SQL with read replicas

### Vertical Scaling

Adjust resource requests and limits in Kubernetes manifests:
```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

## Monitoring and Observability

### Health Endpoints

**Keycloak**:
- `/health/ready` - Readiness check
- `/health/live` - Liveness check
- `/metrics` - Prometheus metrics (if enabled)

**Ingress**:
- Local: ingress-nginx controller exposes `/healthz` internally.
- GCP: Ingress health is managed by the load balancer.

### Logging

All containers log to stdout/stderr, collected by:
- Local: `kubectl logs`
- GCP: Cloud Logging (formerly Stackdriver)

### Metrics

**GCP Integration**:
- GKE metrics in Cloud Monitoring
- Cloud SQL metrics
- Custom application metrics via Prometheus

## Disaster Recovery

### Backup Strategy

**Database**:
- Local: Manual `pg_dump` backups
- GCP: Automated Cloud SQL backups (configurable retention)

**Configuration**:
- Version controlled in Git
- Kubernetes manifests as code

### Recovery Procedures

1. **Application failure**: Kubernetes self-healing (restart pods)
2. **Database failure**: Restore from backup
3. **Complete disaster**: Redeploy from Git repository

## Performance Considerations

### Expected Capacity

**Keycloak** (per replica):
- ~1000 active users
- ~100 requests/second
- ~512MB RAM baseline + session data

**Ingress** (local ingress-nginx / GCP load balancer):
- Handles routing; capacity is determined by controller and backend (Keycloak).

### Optimization Tips

1. Enable Keycloak caching
2. Use CDN for static assets
3. Implement connection pooling
4. Tune database connections
5. Enable HTTP/2
6. Implement caching headers

## Maintenance and Updates

### Updating Components

1. **Update base images** in Dockerfiles
2. **Rebuild images** with Docker (`docker build -t keycloak:1.0.0 Keycloak/`, etc.)
3. **Test in local environment**
4. **Deploy to GCP** with rolling updates (e.g. `./scripts/deploy-gcp.sh` or push then `kubectl rollout restart`)
5. **Monitor** for issues
6. **Rollback** if needed

### Rolling Updates

Kubernetes performs rolling updates automatically:
```bash
kubectl rollout restart deployment keycloak
kubectl rollout status deployment keycloak
```

### Rollback

```bash
kubectl rollout undo deployment keycloak
kubectl rollout history deployment keycloak
```

## Future Enhancements

### Potential Improvements

1. **Service Mesh**: Implement Istio or Linkerd for advanced traffic management
2. **Observability**: Add distributed tracing with Jaeger or Zipkin
3. **GitOps**: Implement ArgoCD or Flux for automated deployments
4. **Multi-region**: Deploy across multiple GCP regions
5. **API Gateway**: Add Kong or Ambassador for advanced API management
6. **Message Queue**: Add RabbitMQ or Kafka for async processing
7. **Caching Layer**: Add Redis for session caching
8. **WAF**: Implement Cloud Armor for DDoS protection

## References

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [ingress-nginx Documentation](https://kubernetes.github.io/ingress-nginx/)
- [GKE Ingress](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-configuration)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/15/)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [GKE Best Practices](https://cloud.google.com/kubernetes-engine/docs/best-practices)
- [Docker Documentation](https://docs.docker.com/)
