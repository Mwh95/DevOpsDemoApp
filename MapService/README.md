# Map API (Go)

REST API for the Map Markers app.
## TLDR
```bash
docker build -t map-api:dev .
minikube image load map-api:dev
kubectl apply -f ./k8s/local/deployment.yaml
kubectl rollout status deployment/map-api
# some times you need to restart if the previous deployment was not successful
# kubectl rollout restart deployment/map-api

```

## Build the container image

Build from the `MapService/` directory:

```bash
docker build -t map-api:dev .
```

## Deploy this service to local Kubernetes

These steps deploy only `MapService`. They assume your local cluster already exists and that the standalone `Liquibase` job has already provisioned the database schema.

1. Review `MapService/k8s/local/deployment.yaml` and set the runtime values you want in `map-api-config` and `map-api-db-secret`.

2. Load the image into your local cluster if needed:

For `minikube`:

```bash
minikube image load map-api:dev
```

For Docker Desktop Kubernetes, the local Docker image is usually available directly.

3. Apply this service manifest:

```bash
kubectl apply -f ./k8s/local/deployment.yaml
```

4. Wait for the rollout:

```bash
kubectl rollout status deployment/map-api
```

5. Access the service directly:

```bash
kubectl port-forward svc/map-api 8090:8090
```

Then use:

- `http://localhost:8090/api`
- `http://localhost:8090/public/health/ready`
- `http://localhost:8090/public/health/live`

## Local manifest contents

`MapService/k8s/local/deployment.yaml` creates:

- a `Deployment` for `map-api`
- a `Service` on port `8090`
- a `ConfigMap` for non-secret runtime values
- a `Secret` for database credentials

For local Kubernetes, keep `KEYCLOAK_ISSUER` set to the ingress URL that appears in tokens, and set `KEYCLOAK_JWKS_URL` to the in-cluster Keycloak service URL so the pod can fetch signing keys without calling back to `localhost`.

## Useful commands

```bash
kubectl get pods,svc
kubectl describe deployment map-api
kubectl logs deployment/map-api
kubectl delete -f MapService/k8s/local/deployment.yaml
```

## Tests

- Unit tests: `go test -short ./internal/...`
- Integration tests: `make test-integration`
- Coverage: `make test-coverage`