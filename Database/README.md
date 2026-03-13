# PostgreSQL Database

This module contains the PostgreSQL database setup for local development.

## Local Development

Or directly with docker-compose:
```bash
cd Database
docker-compose up -d
```

Or directly:
```bash
cd Database
docker-compose down
```


## Configuration

- **Image**: PostgreSQL 18.2
- **Port**: 5432
- **Database**: MapMarkerDb
- **Data Volume**: postgres-data (persists data between restarts)

## Connection Details

### From Host Machine
```
Host: localhost
Port: 5432
Database: MapMarkerDb
Username: postgres
Password: ...
```

### From Kubernetes
The database runs as a container in local development. For GCP, use Cloud SQL.

## DRAFT: GCP Cloud SQL

For production deployment on GCP:

1. Create a Cloud SQL PostgreSQL instance
2. Configure the Cloud SQL Proxy sidecar in Kubernetes
3. Update Keycloak connection settings in `Keycloak/k8s/gcp/deployment.yaml`

See [Cloud SQL Proxy](https://cloud.google.com/sql/docs/postgres/connect-kubernetes-engine) documentation for details.
