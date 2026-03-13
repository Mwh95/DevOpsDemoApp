# Keycloak Authentication Service

This module contains the Keycloak authentication service setup with Docker and Kubernetes configurations.

## Build Docker Image

```shell
eval $(minikube docker-env)
docker build . -t keycloak:dev
```

## Deploy to Kubernetes
###
Start local cluster
```shell
minikube start
```
### Local Deployment
```shell
kubectl apply -f k8s/local/deployment.yaml
```

rollout changes
```shell
kubectl rollout restart deployment keycloak -n default
```

Pod status
```shell
kubectl get pods
# or
kubectl describe pods
```

Rollout status
```shell
kubectl rollout status deployment/keycloak   
```

Check deployments
```shell
kubectl get deployments
```

Check services
```shell
kubectl get services
kubectl describe services/keycloak
```

Proxy to Cluster API via second terminal if no service is defined
```shell
kubectl proxy
```
Inspect the Keycloak container:
```shell
kubectl exec -it <pod-id> -- /bin/bash
cd /opt/keycloak
```


if needed during initial setup: direct portmapping into minikube
```shell
 kubectl port-forward svc/keycloak 8080:8080 9000:9000
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
- Admin credentials set to `tmpadmin:admin`
- Database connection to local postgres service

### GCP
- Admin credentials via Kubernetes secrets
- Database connection via Cloud SQL proxy
- HTTPS enforced with proper hostname

## Map Markers App (OIDC client)

The Map Markers frontend uses Keycloak for login. Configure a realm and client as follows.

### 1. Create a realm (optional)

- Admin Console: **Keycloak Admin** → **Create realm** (e.g. `users`).
- Note the realm name; the issuer URL will be `{KEYCLOAK_BASE}/realms/{REALM}` (e.g. `http://localhost:50594/login/realms/users`).

### 2. Create a client for the Map SPA

- In the realm: **Clients** → **Create client**.
- **Client type**: `OpenID Connect`.
- **Client ID**: `map-app` (or set `VITE_OIDC_CLIENT_ID` in the frontend to match).
- **Client authentication**: OFF (public client).
- **Valid redirect URIs**: Add the SPA origin, e.g. `http://localhost:NODEPORT/*`.
- **Web origins**: Add the same origin for CORS.
- Save.

### 3. Map API configuration

Set `KEYCLOAK_ISSUER` to the realm issuer URL (e.g. `http://localhost:50594/login/realms/users`). The frontend needs `VITE_OIDC_AUTHORITY` (same value), `VITE_OIDC_CLIENT_ID` (`map-app`), and `VITE_OIDC_REDIRECT_URI` (SPA origin).
