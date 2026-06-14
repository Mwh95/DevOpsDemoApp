# SystemTest

End-to-end regression tests for DemoApp. The whole stack is started with
[Testcontainers](https://testcontainers.com/) (PostgreSQL, Liquibase, Keycloak, Map API, Map
Frontend, plus an nginx gateway), test scenarios are described in Cucumber/Gherkin, and the UI
flows are driven by [Playwright for Java](https://playwright.dev/java/). REST checks hit the Map
API directly.

## What is covered

- **Authentication** (`features/authentication.feature`, `@ui`): the sign-in prompt, signing in
  through Keycloak, and signing out.
- **Marker UI** (`features/markers_ui.feature`, `@ui`): map loads, create a marker via the
  "Add marker at center" button and via map clicks, edit a marker, delete a marker.
- **Marker API** (`features/markers_api.feature`, `@api`): health endpoints, unauthenticated
  rejection, full create/read/update/list/delete lifecycle, and per-user isolation.

## How it works

A reverse-proxy gateway (nginx) is bound to a **fixed host port** (default `58080`) and routes:

| Path        | Upstream            |
|-------------|---------------------|
| `/`         | `map-frontend:8080` |
| `/api`      | `map-api:8090`      |
| `/public`   | `map-api:8090`      |
| `/login`    | `keycloak:8080`     |
| `/keycloak` | `keycloak:9000`     |

Using one origin keeps the browser URL, the OIDC `redirect_uri`, and the JWT `iss` claim
consistent, mirroring the real ingress.

Keycloak is configured by a dedicated one-shot **Terraform** container
(`keycloak-terraform/`, image `demoapp-keycloak-terraform:systemtest`) that runs after Keycloak
is healthy. Using the [Keycloak Terraform provider](https://registry.terraform.io/providers/keycloak/keycloak),
it provisions realm `users`, the public client `map-app` (standard + PKCE flow, plus direct access
grants so the API tests can use the password grant), and two users: `testuser` and `otheruser`
(both `Test1234!`). Connection details and values are passed as `TF_VAR_*` environment variables,
so re-running the container reconciles Keycloak whenever the configuration changes.

The stack uses dedicated `:systemtest` image tags so the local `:dev` images used by
`scripts/deploy-local.sh` are never overwritten.

## Prerequisites

- Docker (running)
- **Java 21 as the active JDK** — the scripts run Gradle with whatever `java`/`JAVA_HOME` is on
  your shell, and the Gradle wrapper (8.10) requires JDK ≤ 23. The launcher scripts check this and
  fail fast if the active JDK is not Java 21; set `JAVA_HOME` to a Java 21 install if needed.
- Internet access on first run (Gradle dependencies, Playwright browser, base images)

## Running

From this directory:

```bash
./run-system-tests.sh
```

This builds the five `:systemtest` images (Keycloak, Map API, Map Frontend, Liquibase, and the
Keycloak Terraform configurator), installs the Playwright browser, and runs every scenario.
Useful variants:

```bash
./run-system-tests.sh --no-build        # reuse existing :systemtest images
./run-system-tests.sh --tags @api       # only API scenarios
./run-system-tests.sh --tags @ui        # only UI scenarios
```

You can also run via Gradle directly once images exist:

```bash
./gradlew installPlaywright
./gradlew test -Dcucumber.filter.tags="@api"
```

Reports are written to `build/reports/cucumber/cucumber.html` (and the standard Gradle test
report under `build/reports/tests/test/`).

## Configuration

| Variable                  | Default      | Purpose                                            |
|---------------------------|--------------|----------------------------------------------------|
| `SYSTEMTEST_IMAGE_TAG`    | `systemtest` | Image tag to build and run                         |
| `SYSTEMTEST_GATEWAY_PORT` | `58080`      | Fixed host port for the gateway                    |

> The redirect URIs are derived from the gateway port automatically and passed to the Terraform
> container as `TF_VAR_redirect_uris`, so changing `SYSTEMTEST_GATEWAY_PORT` needs no manual edits.

## Cleanup

Testcontainers cleans up automatically on JVM exit. After hard kills, or to remove the built
images and build output:

```bash
./cleanup.sh            # remove leftover testcontainers resources + gradle clean
./cleanup.sh --images   # also remove the :systemtest images
```
