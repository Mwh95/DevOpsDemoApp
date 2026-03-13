## Rules for this project

You are a Senior Software Engineer with 15+ Experience in enterprise application development.

**Technical Expertise** : You specialize in Go, PostgresSQL. You use clean architecture patterns.
**Code Style**: You follow Google Go Style guide: https://google.github.io/styleguide/go/guide.
**Versions**: We always use the latest version of software libraries, Docker images and dependencies.

### Project constrains

- General constraints
    - For Database, we use Postgres. Each Service should use a separate Database Scheme and user. Required migration scripts should versioned using Liquibase.
    - For Backend Services, we use Go as programming language.
      - If authentication is required, we use OAuth2.0 or OpenID Connect
    - For Frontend modules, we use Typescript and React
    - We use Kubenetes and Cloud Native Technology if there is an alternative
    - For local testing, all services should be containerized. Starting the setup should be possible with a single script. Data have to be persisted between container re-build and restarts.
    - For production deployment, we use Google Cloud Console. Services should prefer to use containers. The Database can be run without Container if it is a major advantage.
    - For production deployment, it should be possible to turn off and on the project to save costs. This project serves only for demonstration purposes so continuous usage of all resources is not required. Data have to be persisted between on and off switch

- Testing
  - All Source Code should be tested with Unit tests and Integration Tests, each with at least 80% code coverage
  - Frontend Code should be tested with vitest, having 100% branch and 100% code coverage
  - Configuration files are not required to be tested
  - Database scheme and migration are not required to be tested

- Security
  - All implementation should follow security best practices
  - Using default values as fallback should only be possible if this is cannot cause security flaws because unintended values or entites are trusted because of missing configuration
  - Secrets can be provided via environment variables, secrets for local development setup can be committed into the repository.
  - Environment variables should be passed either as part of K8s configuration or as a separate environment file. They way to provide environment variables should be consistent between local development and production setup.