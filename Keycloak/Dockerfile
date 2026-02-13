FROM quay.io/keycloak/keycloak:23.0

# Set the working directory
WORKDIR /opt/keycloak

# Copy custom configuration if needed
# COPY realm-export.json /opt/keycloak/data/import/

# Set environment variables
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres

# Build optimized Keycloak
RUN /opt/keycloak/bin/kc.sh build

# Expose ports
EXPOSE 8080
EXPOSE 8443

# Set the entrypoint
ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]
