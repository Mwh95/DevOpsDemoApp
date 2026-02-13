# Architecture Documentation

## Overview

This application is a multi-service platform consisting of three main components:

1. **Keycloak** - Authentication and authorization service
2. **Reverse Proxy** - Apache HTTPD load balancer and reverse proxy
3. **Database** - PostgreSQL database for persistent storage

All components are containerized and can be deployed on both local Kubernetes clusters and Google Cloud Platform (GCP).

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Internet                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ HTTP/HTTPS
                         ▼
          ┌──────────────────────────────┐
          │  Load Balancer (K8s Service) │
          │      Type: LoadBalancer      │
          └──────────────┬───────────────┘
                         │
                         │
          ┌──────────────▼───────────────┐
          │     Reverse Proxy (httpd)    │
          │         2 Replicas           │
          │    Port: 80, 443             │
          └──────────────┬───────────────┘
                         │
                         │ /auth → :8080
                         ▼
          ┌──────────────────────────────┐
          │       Keycloak Service       │
          │         2 Replicas           │
          │    Port: 8080, 8443          │
          └──────────────┬───────────────┘
                         │
                         │ JDBC
                         ▼
          ┌──────────────────────────────┐
          │    PostgreSQL Database       │
          │  (Docker Compose / Cloud SQL) │
          │         Port: 5432           │
          └──────────────────────────────┘
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

**Health Checks**:
- Readiness probe: `/health/ready` endpoint
- Liveness probe: `/health/live` endpoint

**Resource Allocation** (GCP):
- Requests: 512Mi memory, 500m CPU
- Limits: 1Gi memory, 1000m CPU

### 2. Reverse Proxy

**Purpose**: Load balancing, SSL termination, and routing

**Technology**: 
- Base Image: `httpd:2.4`
- Server: Apache HTTP Server

**Features**:
- Reverse proxy to Keycloak
- Health status endpoint
- Extensible routing configuration

**Enabled Apache Modules**:
- `mod_proxy` - Proxy support
- `mod_proxy_http` - HTTP proxy
- `mod_proxy_balancer` - Load balancing
- `mod_lbmethod_byrequests` - Load balancing method
- `mod_slotmem_shm` - Shared memory for balancing
- `mod_rewrite` - URL rewriting
- `mod_headers` - HTTP header manipulation
- `mod_ssl` - SSL/TLS support

**Routing Rules**:
```
/auth → http://keycloak:8080/
```

**Kubernetes Resources**:
- **Deployment**: 2 replicas for high availability
- **Service**: LoadBalancer (external) on ports 80, 443

**Health Checks**:
- Readiness probe: `/server-status` endpoint
- Liveness probe: `/server-status` endpoint

**Resource Allocation** (GCP):
- Requests: 128Mi memory, 100m CPU
- Limits: 256Mi memory, 200m CPU

### 3. Database (PostgreSQL)

**Purpose**: Persistent data storage for Keycloak

**Technology**: 
- Base Image: `postgres:15-alpine`
- Database: PostgreSQL 15

**Local Deployment**:
- Runs as Docker container via docker-compose
- Data persisted in named volume `postgres-data`
- Exposed on host port 5432
- Default credentials: `keycloak/keycloak` (⚠️ not for production)

**GCP Deployment**:
- Google Cloud SQL for PostgreSQL
- Managed service with automatic backups
- High availability with failover
- Accessed via Cloud SQL Proxy or private IP

**Configuration**:
- Database name: `keycloak`
- Port: 5432
- Initialization: `config/init.sql`

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
localhost:80 → reverse-proxy:80 → keycloak:8080 → postgres:5432
```

**Access**:
- External: `http://localhost:80`
- Keycloak: `http://localhost:80/auth`

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
Internet → GCP Load Balancer → reverse-proxy:80 → keycloak:8080 → Cloud SQL
```

**Access**:
- External: `http://<EXTERNAL-IP>` or `https://your-domain.com`
- Keycloak: `https://your-domain.com/auth`

**Required GCP Services**:
- Google Kubernetes Engine (GKE)
- Google Container Registry (GCR)
- Cloud SQL for PostgreSQL
- Cloud Load Balancing

## Build and Deployment Pipeline

### Gradle Build Structure

```
PlaygroundApp (root)
├── settings.gradle         # Multi-project configuration
├── build.gradle            # Root build configuration
├── Keycloak/
│   └── build.gradle       # Keycloak-specific tasks
├── ReverseProxy/
│   └── build.gradle       # Reverse Proxy-specific tasks
└── Database/
    └── build.gradle       # Database-specific tasks
```

**Available Gradle Tasks**:
- `buildAll` - Build all Docker images
- `cleanAll` - Clean all project artifacts
- `:Keycloak:buildDockerImage` - Build Keycloak image
- `:ReverseProxy:buildDockerImage` - Build Reverse Proxy image
- `:Database:startLocalDb` - Start local database
- `:Database:stopLocalDb` - Stop local database

### Deployment Flow

#### Local Deployment
1. Build Docker images with Gradle
2. Start PostgreSQL with docker-compose
3. Deploy Kubernetes manifests from `k8s/local/`
4. Services accessible via localhost

#### GCP Deployment
1. Build Docker images with Gradle
2. Tag images for Google Container Registry
3. Push images to GCR
4. Apply Kubernetes manifests from `k8s/gcp/`
5. Configure Cloud SQL connection
6. Services accessible via external IP

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

**Reverse Proxy**:
- Replicas: 2 (default) → 5 (max)
- Scaling trigger: CPU > 70%
- Stateless, no session affinity needed

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

**Reverse Proxy**:
- `/server-status` - Apache status page

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

**Reverse Proxy** (per replica):
- ~5000 concurrent connections
- ~1000 requests/second
- Low memory footprint (~128MB)

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
2. **Rebuild images** with Gradle
3. **Test in local environment**
4. **Deploy to GCP** with rolling updates
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
- [Apache HTTPD Documentation](https://httpd.apache.org/docs/2.4/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/15/)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [GKE Best Practices](https://cloud.google.com/kubernetes-engine/docs/best-practices)
- [Gradle User Manual](https://docs.gradle.org/current/userguide/userguide.html)
