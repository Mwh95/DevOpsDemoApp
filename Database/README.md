# PostgreSQL Database

This module contains the PostgreSQL database setup for local development.

## Local Development

### Start Database
```bash
./gradlew :Database:startLocalDb
```

Or directly with docker-compose:
```bash
cd Database
docker-compose up -d
```

### Stop Database
```bash
./gradlew :Database:stopLocalDb
```

Or directly:
```bash
cd Database
docker-compose down
```

### View Logs
```bash
./gradlew :Database:logsDb
```

## Configuration

- **Image**: PostgreSQL 15 Alpine
- **Port**: 5432
- **Database**: keycloak
- **Username**: keycloak
- **Password**: keycloak (⚠️ Change for production!)
- **Data Volume**: postgres-data (persists data between restarts)

## Connection Details

### From Host Machine
```
Host: localhost
Port: 5432
Database: keycloak
Username: keycloak
Password: keycloak
```

### From Docker Containers (same network)
```
Host: postgres
Port: 5432
Database: keycloak
Username: keycloak
Password: keycloak
```

### From Kubernetes
The database runs as a container in local development. For GCP, use Cloud SQL.

## GCP Cloud SQL

For production deployment on GCP:

1. Create a Cloud SQL PostgreSQL instance
2. Configure the Cloud SQL Proxy sidecar in Kubernetes
3. Update Keycloak connection settings in `Keycloak/k8s/gcp/deployment.yaml`

See [Cloud SQL Proxy](https://cloud.google.com/sql/docs/postgres/connect-kubernetes-engine) documentation for details.
