#!/bin/sh
set -eu

: "${PG_HOST:?missing PG_HOST}"
: "${PG_PORT:=5432}"
: "${PG_DATABASE:?missing PG_DATABASE}"
: "${PG_JDBC_PARAMS:=}"
: "${MAPSERVICE_SCHEMA:=mapservice}"
: "${PG_BOOTSTRAP_USER:?missing PG_BOOTSTRAP_USER}"
: "${PG_BOOTSTRAP_PASSWORD:?missing PG_BOOTSTRAP_PASSWORD}"

JDBC_URL="jdbc:postgresql://${PG_HOST}:${PG_PORT}/${PG_DATABASE}${PG_JDBC_PARAMS}"
echo "Starting Liquibase database migration ..."

run_liquibase() {
  changelog_file="$1"
  default_schema="${2:-}"

  if [ -n "$default_schema" ]; then
    liquibase \
      --changelog-file="$changelog_file" \
      --url="$JDBC_URL" \
      --username="$PG_BOOTSTRAP_USER" \
      --password="$PG_BOOTSTRAP_PASSWORD" \
      --default-schema-name="$default_schema" \
      update
    return
  fi

  liquibase \
    --changelog-file="$changelog_file" \
    --url="$JDBC_URL" \
    --username="$PG_BOOTSTRAP_USER" \
    --password="$PG_BOOTSTRAP_PASSWORD" \
    update
}

run_liquibase modules/database/bootstrap/changelog.xml
run_liquibase modules/mapservice/changelog/changelog.xml "$MAPSERVICE_SCHEMA"

echo "Finished Liquibase database migration."
