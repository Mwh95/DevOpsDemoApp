#!/usr/bin/env bash
# Build the dedicated :systemtest images and run the Playwright + Cucumber system tests.
#
# Usage:
#   ./run-system-tests.sh                 # build images, install browsers, run all tests
#   ./run-system-tests.sh --no-build      # skip image builds (reuse existing :systemtest images)
#   ./run-system-tests.sh --tags @api     # only run scenarios with the given Cucumber tag(s)
#
# Environment:
#   SYSTEMTEST_IMAGE_TAG     image tag to build/use (default: systemtest)
#   SYSTEMTEST_GATEWAY_PORT  fixed host port for the gateway (default: 58080)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

IMAGE_TAG="${SYSTEMTEST_IMAGE_TAG:-systemtest}"
BUILD_IMAGES=true
TAGS=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-build)
      BUILD_IMAGES=false
      ;;
    --tags)
      shift
      TAGS="${1:-}"
      ;;
    -h|--help)
      grep '^#' "$0" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
  shift
done

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

# The Gradle wrapper (8.10) requires JDK <= 23 and the project targets Java 21, so fail fast unless
# the active JDK (the one ./gradlew will use) is Java 21. Set JAVA_HOME to a Java 21 install if not.
require_java_21() {
  local java_bin="java"
  [[ -n "${JAVA_HOME:-}" ]] && java_bin="$JAVA_HOME/bin/java"
  if ! command -v "$java_bin" >/dev/null 2>&1; then
    echo "Java not found; this project requires Java 21. Set JAVA_HOME to a Java 21 install." >&2
    exit 1
  fi
  local major
  major="$("$java_bin" -version 2>&1 | awk -F'[".]' '/version/{print $2; exit}')"
  if [[ "$major" != 21 ]]; then
    echo "Java 21 is required, but the active JDK is ${major:-unknown}. Set JAVA_HOME to a Java 21 install." >&2
    exit 1
  fi
}

require_command docker
require_java_21

# Testcontainers' Java client cannot auto-detect non-default Docker daemons (e.g. Rancher Desktop
# or a custom docker context). Derive DOCKER_HOST from the active context if it is not already set,
# and bind-mount the in-VM socket for Ryuk.
ensure_docker_host() {
  if [[ -z "${DOCKER_HOST:-}" ]]; then
    local endpoint
    endpoint="$(docker context inspect --format '{{ .Endpoints.docker.Host }}' 2>/dev/null || true)"
    if [[ -n "$endpoint" ]]; then
      export DOCKER_HOST="$endpoint"
    fi
  fi
  if [[ -n "${DOCKER_HOST:-}" && "$DOCKER_HOST" == unix://* ]]; then
    export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="${TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE:-/var/run/docker.sock}"
  fi
  # The bundled docker-java client otherwise falls back to API 1.32, which modern daemons reject.
  if [[ -z "${DOCKER_API_VERSION:-}" ]]; then
    local api
    api="$(docker version --format '{{ .Server.APIVersion }}' 2>/dev/null || true)"
    if [[ -n "$api" ]]; then
      export DOCKER_API_VERSION="$api"
    fi
  fi
  echo "==> Using DOCKER_HOST=${DOCKER_HOST:-<default>} DOCKER_API_VERSION=${DOCKER_API_VERSION:-<negotiated>}"
}

ensure_docker_host

build_image() {
  local context="$1" dockerfile="$2" image="$3" label="$4"
  echo "==> Building $label image ($image)"
  docker build -f "$dockerfile" -t "$image" "$context"
}

if [[ "$BUILD_IMAGES" == true ]]; then
  build_image "$REPO_ROOT/Keycloak"               "$REPO_ROOT/Keycloak/Dockerfile"               "keycloak:${IMAGE_TAG}"                  "Keycloak"
  build_image "$REPO_ROOT/MapService"             "$REPO_ROOT/MapService/Dockerfile"             "map-api:${IMAGE_TAG}"                   "Map API"
  build_image "$REPO_ROOT/MapFrontend"            "$REPO_ROOT/MapFrontend/Dockerfile"            "map-frontend:${IMAGE_TAG}"              "Map Frontend"
  build_image "$REPO_ROOT/Liquibase"              "$REPO_ROOT/Liquibase/Dockerfile"              "demoapp-liquibase:${IMAGE_TAG}"         "Liquibase"
  build_image "$SCRIPT_DIR/keycloak-terraform"    "$SCRIPT_DIR/keycloak-terraform/Dockerfile"    "demoapp-keycloak-terraform:${IMAGE_TAG}" "Keycloak Terraform"
else
  echo "==> Skipping image builds (--no-build); expecting existing :${IMAGE_TAG} images"
fi

cd "$SCRIPT_DIR"

echo "==> Installing Playwright browsers"
./gradlew --console=plain --no-daemon installPlaywright

echo "==> Running system tests"
GRADLE_ARGS=(test)
if [[ -n "$TAGS" ]]; then
  GRADLE_ARGS+=("-Dcucumber.filter.tags=$TAGS")
fi

SYSTEMTEST_IMAGE_TAG="$IMAGE_TAG" ./gradlew --console=plain --no-daemon "${GRADLE_ARGS[@]}"

echo "==> Done. Reports: SystemTest/build/reports/cucumber/cucumber.html"
