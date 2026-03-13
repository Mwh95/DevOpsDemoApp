# Liquibase

Local database bootstrap and schema migration module for the DemoApp.

This module runs as a standalone Kubernetes `Job`. It:

- applies `Liquibase/modules/database/bootstrap/changelog.xml` as the PostgreSQL bootstrap changelog
- runs the `Liquibase/modules/mapservice/changelog/` Liquibase changelog after bootstrap completes

## Layout

- `modules/database/bootstrap/` contains shared database bootstrap Liquibase changelogs
- `modules/mapservice/changelog/` contains `MapService` schema migrations

## Build

Build from the repository root:

```bash
docker build -t demoapp-liquibase:dev .
```

## Deploy to local Kubernetes

1. Ensure the local PostgreSQL service is already running in the cluster.
2. Load the image into your local cluster if needed.

For `minikube`:

```bash
minikube image load demoapp-liquibase:dev
```

3. Apply the job:

```bash
kubectl delete job/liquibase --ignore-not-found=true
kubectl apply -f ./k8s/local/job.yaml
```

The manifest includes:

- `ConfigMap/liquibase-config` for database connection settings
- `Secret/liquibase-bootstrap-secret` for the bootstrap username and password

You can edit the `ConfigMap` values before applying, for example:

- `PG_HOST`
- `PG_PORT`
- `PG_DATABASE`
- `PG_JDBC_PARAMS` for extra JDBC options such as `?sslmode=disable`
- `MAPSERVICE_SCHEMA`

4. Watch completion and inspect logs:

```bash
kubectl wait --for=condition=complete job/liquibase --timeout=60s
kubectl logs job/liquibase
```

5. Re-run the job after changes:

```bash
kubectl delete job/liquibase --ignore-not-found=true
kubectl apply -f Liquibase/k8s/local/job.yaml
```
