#!/usr/bin/env bash
# Clean up after a system-test run.
#
# Testcontainers' Ryuk reaper already removes the started containers, the shared network and
# volumes when the test JVM exits (even on failure), so this script is only needed after hard
# kills / ctrl-c, or to remove the built test images and Gradle build artifacts.
#
# Usage:
#   ./cleanup.sh             # remove leftover testcontainers resources + ./gradlew clean
#   ./cleanup.sh --images    # also remove the five :systemtest images
#
# Environment:
#   SYSTEMTEST_IMAGE_TAG     image tag to remove with --images (default: systemtest)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_TAG="${SYSTEMTEST_IMAGE_TAG:-systemtest}"
REMOVE_IMAGES=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --images) REMOVE_IMAGES=true ;;
    -h|--help) grep '^#' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "Unknown argument: $1" >&2; exit 1 ;;
  esac
  shift
done

if command -v docker >/dev/null 2>&1; then
  echo "==> Removing leftover Testcontainers containers"
  leftover_containers="$(docker ps -aq --filter label=org.testcontainers=true)"
  if [[ -n "$leftover_containers" ]]; then
    docker rm -f $leftover_containers
  fi

  echo "==> Removing leftover Testcontainers networks"
  leftover_networks="$(docker network ls -q --filter label=org.testcontainers=true)"
  if [[ -n "$leftover_networks" ]]; then
    docker network rm $leftover_networks >/dev/null 2>&1 || true
  fi

  if [[ "$REMOVE_IMAGES" == true ]]; then
    echo "==> Removing :${IMAGE_TAG} test images"
    docker image rm -f \
      "keycloak:${IMAGE_TAG}" \
      "map-api:${IMAGE_TAG}" \
      "map-frontend:${IMAGE_TAG}" \
      "demoapp-liquibase:${IMAGE_TAG}" \
      "demoapp-keycloak-terraform:${IMAGE_TAG}" 2>/dev/null || true
  fi
else
  echo "docker not found; skipping container/image cleanup" >&2
fi

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

require_java_21

echo "==> Cleaning Gradle build output"
cd "$SCRIPT_DIR"
./gradlew --console=plain clean

echo "==> Cleanup complete"
