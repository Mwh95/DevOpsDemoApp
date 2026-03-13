ARG KC_VERSION=26.5.0
FROM quay.io/keycloak/keycloak:${KC_VERSION} AS builder

# Set the working directory
WORKDIR /opt/keycloak

# Set environment variables
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres
ENV KC_HTTP_RELATIVE_PATH=/login
ENV KC_HTTP_MANAGEMENT_RELATIVE_PATH=/keycloak
# Build optimized Keycloak
RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:${KC_VERSION}
COPY --from=builder /opt/keycloak/ /opt/keycloak/

# Set environment variables
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres
ENV KC_HTTP_RELATIVE_PATH=/login
ENV KC_HTTP_MANAGEMENT_RELATIVE_PATH=/keycloak
ENV KC_HTTP_ACCESS_LOG_ENABLED=true
ENV KC_HTTP_ACCESS_LOG_PATTERN=combined
ENV KC_CACHE=ispn

# Expose ports
EXPOSE 8080
EXPOSE 8443

# Set the entrypoint
ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]
