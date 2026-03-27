---
name: demoapp-standards
description: Applies DemoApp repository standards for Go, PostgreSQL, React, Kubernetes, and GCP deployment. Use when adding or changing backend services, frontend features, database schemas or Liquibase migrations, deployment scripts, Kubernetes manifests, or architecture decisions in this repository.
---

# DemoApp Standards

## Quick Start

When working in this repository:

1. Read `README.md` and `ARCHITECTURE.md` before making architecture or deployment changes.
2. Follow clean architecture patterns for backend services.
3. Prefer the repository defaults unless the user asks otherwise:
   - Backend: Go
   - Frontend: TypeScript and React
   - Database: PostgreSQL
   - Auth: OAuth 2.0 or OpenID Connect when authentication is needed
   - Platform: Kubernetes and cloud-native tooling
4. Keep local workflows containerized and runnable from the repository scripts.

## Project Constraints

- Each service uses its own PostgreSQL schema and database user.
- Database changes are versioned with Liquibase.
- Local development stays containerized and startable from a single script.
- Local data must survive container rebuilds and restarts.
- Production targets Google Cloud and should stay easy to turn on and off without losing data.

## Software Quality

- Follow Google Go style for Go code.
- Use current stable dependencies and base images unless the repo already pins something specific.
- Add unit and integration tests for source code when practical.
- Frontend tests use Vitest.
- Configuration files and database bootstrap/migration files do not need direct tests unless the user asks for them.

## Security

- Apply security best practices by default.
- Do not add insecure fallback values when missing configuration could trust the wrong system or identity.
- Keep environment variable delivery consistent between local and production setups.

## Service Specific Checklist
### Backend: Go Server Checklist

Use this when changing `MapService` or another Go backend:

- Keep domain, application, and adapter concerns separated.
- Prefer explicit interfaces at boundaries instead of leaking infrastructure concerns inward.
- Use PostgreSQL for persistence.
- If the service needs auth, use OAuth 2.0 or OpenID Connect patterns already present in the repo.
- Avoid insecure default config values for issuer URLs, allowed origins, credentials, or trusted identities.
- Add or update unit tests with the code change.
- Add integration coverage when the change affects persistence, HTTP wiring, or auth behavior.

Before finishing a Go backend change, verify:

- Request and response behavior still matches the public API contract.
- New config is documented and wired consistently for local and production deployment.
- Schema or permission changes are reflected in Liquibase rather than ad hoc SQL drift.

### Liquibase And Database Checklist

Use this when changing schemas, tables, seed data, or service persistence:

- Keep PostgreSQL as the default target database.
- Give each service its own schema and database user.
- Version schema changes with Liquibase.
- Prefer additive, reviewable migrations over manual one-off steps.
- Keep bootstrap or shared database setup aligned with service ownership boundaries.

Before finishing a database change, verify:

- The migration belongs in the correct module or changelog chain.
- The affected schema and runtime user match the owning service.
- Local and production assumptions remain consistent.
- Application code and migration names reflect the same data model.

### Kubernetes And Script Checklist

Use this when changing manifests, deployment scripts, or local environment automation:

- Prefer Kubernetes-native and container-first workflows.
- Keep local setup runnable from the repository scripts.
- Preserve persistent local data across rebuilds and restarts.
- Keep environment variable injection consistent between local and production styles.
- Prefer containerized services in production unless a non-containerized database is clearly justified.
- Keep GCP deployment choices compatible with shutting environments down and restoring them later without data loss.

Before finishing a deployment or script change, verify:

- Local setup still works from the documented scripts.
- Manifest changes match the documented architecture and routing.
- Secrets and config are not hardcoded in insecure ways beyond intentionally local-only defaults.
- New operational steps are reflected in `README.md` or `DEPLOYMENT.md` when needed.

## Decision Defaults

When several approaches are possible, prefer the one that:

- Fits the existing repository architecture.
- Minimizes special-case local setup.
- Keeps production deployment cloud-native.
- Preserves security over convenience.
- Reduces hidden operational knowledge by documenting the workflow in repo files.

## Additional Resources

- For system structure, see `ARCHITECTURE.md`.
- For local and deployment workflows, see `README.md` and `DEPLOYMENT.md`.
