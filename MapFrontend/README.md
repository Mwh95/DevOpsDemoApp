# Map Frontend (React + TypeScript)

SPA for the Map Markers app: OpenStreetMap via Leaflet, edit mode to add markers by clicking, label and note per marker, Keycloak login.

## Setup

```bash
npm install
```

## Dev

```bash
npm run dev
```

Uses Vite proxy: `/api` and `/login` go to the Map API (default `http://localhost:8080`). Configure Keycloak (see Keycloak README) and set env if needed:

- `VITE_OIDC_AUTHORITY` – Keycloak realm URL (e.g. `http://localhost:50594/login/realms/users`)
- `VITE_OIDC_CLIENT_ID` – client id (e.g. `map-app`)
- `VITE_OIDC_REDIRECT_URI` – SPA origin (e.g. `http://localhost:5173/`)
- `VITE_API_BASE` – API base URL (default `''` for same origin/proxy)

## Build

```bash
npm run build
```

Output is written to `dist/`.

## Container build

Build the standalone frontend container:

```bash
docker build -t map-frontend:dev .
```

## Local Kubernetes

Apply the local frontend Deployment and Service:

```bash
kubectl apply -f ./k8s/local/deployment.yaml
kubectl rollout status deployment/map-frontend
```

## Tests

- **Run tests**: `npm run test` (Vitest, single run).
- **Watch mode**: `npm run test:watch`.
- **Coverage**: `npm run test:coverage` (reports in `coverage/`).
