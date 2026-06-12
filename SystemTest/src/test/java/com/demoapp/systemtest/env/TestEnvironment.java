package com.demoapp.systemtest.env;

import java.time.Duration;
import java.util.List;

import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.Network;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.containers.startupcheck.OneShotStartupCheckStrategy;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.utility.DockerImageName;
import org.testcontainers.utility.MountableFile;

/**
 * Brings up the full DemoApp stack once for the whole test run using Testcontainers.
 *
 * <p>The containers run on a shared Docker network and are fronted by an nginx gateway bound to a
 * fixed host port, so the browser, the OIDC {@code redirect_uri}, and the JWT {@code iss} claim all
 * share one origin (mirroring the real ingress). Testcontainers' Ryuk reaper tears everything down
 * when the JVM exits, so no explicit shutdown is required.
 */
public final class TestEnvironment {

    /** Fixed host port for the gateway; must match the redirect URIs in the realm export. */
    public static final int GATEWAY_PORT = Integer.parseInt(env("SYSTEMTEST_GATEWAY_PORT", "58080"));

    private static final String IMAGE_TAG = env("SYSTEMTEST_IMAGE_TAG", "systemtest");

    private static final String DB_NAME = "MapMarkerDb";
    private static final String REALM = "users";

    private static volatile TestEnvironment instance;

    private final Network network = Network.newNetwork();

    private PostgreSQLContainer<?> postgres;
    private GenericContainer<?> liquibase;
    private GenericContainer<?> keycloak;
    private GenericContainer<?> keycloakConfig;
    private GenericContainer<?> mapApi;
    private GenericContainer<?> mapFrontend;
    private GenericContainer<?> gateway;

    private TestEnvironment() {
    }

    /** Lazily starts the stack on first access and returns the shared instance. */
    public static TestEnvironment get() {
        if (instance == null) {
            synchronized (TestEnvironment.class) {
                if (instance == null) {
                    TestEnvironment env = new TestEnvironment();
                    env.start();
                    instance = env;
                }
            }
        }
        return instance;
    }

    public String baseUrl() {
        return "http://localhost:" + GATEWAY_PORT;
    }

    public String issuer() {
        return baseUrl() + "/login/realms/" + REALM;
    }

    public String tokenEndpoint() {
        return issuer() + "/protocol/openid-connect/token";
    }

    @SuppressWarnings("resource")
    private void start() {
        postgres = new PostgreSQLContainer<>(
                DockerImageName.parse("postgres:18.2-trixie").asCompatibleSubstituteFor("postgres"))
                .withDatabaseName(DB_NAME)
                .withUsername("postgres")
                .withPassword("postgres")
                .withNetwork(network)
                .withNetworkAliases("postgres")
                .waitingFor(Wait.forLogMessage(".*database system is ready to accept connections.*", 2)
                        .withStartupTimeout(Duration.ofMinutes(2)));
        postgres.start();

        liquibase = new GenericContainer<>(image("demoapp-liquibase"))
                .withNetwork(network)
                .withEnv("PG_HOST", "postgres")
                .withEnv("PG_PORT", "5432")
                .withEnv("PG_DATABASE", DB_NAME)
                .withEnv("PG_JDBC_PARAMS", "")
                .withEnv("MAPSERVICE_SCHEMA", "mapservice")
                .withEnv("PG_BOOTSTRAP_USER", "postgres")
                .withEnv("PG_BOOTSTRAP_PASSWORD", "postgres")
                .withStartupCheckStrategy(new OneShotStartupCheckStrategy().withTimeout(Duration.ofMinutes(3)));
        liquibase.start();

        keycloak = new GenericContainer<>(image("keycloak"))
                .withNetwork(network)
                .withNetworkAliases("keycloak")
                .withExposedPorts(8080, 9000)
                .withEnv("KC_DB", "postgres")
                .withEnv("KC_DB_URL", "jdbc:postgresql://postgres:5432/" + DB_NAME + "?currentSchema=keycloak")
                .withEnv("KC_DB_USERNAME", "keycloak")
                .withEnv("KC_DB_PASSWORD", "keycloak")
                .withEnv("KC_BOOTSTRAP_ADMIN_USERNAME", "tmpadmin")
                .withEnv("KC_BOOTSTRAP_ADMIN_PASSWORD", "admin")
                .withEnv("KC_HOSTNAME", baseUrl() + "/login")
                .withEnv("KC_HOSTNAME_STRICT", "false")
                .withEnv("KC_HTTP_ENABLED", "true")
                .withEnv("KC_PROXY_HEADERS", "xforwarded")
                .withEnv("KC_HEALTH_ENABLED", "true")
                .withCommand("start", "--optimized")
                .waitingFor(Wait.forHttp("/keycloak/health/ready").forPort(9000)
                        .withStartupTimeout(Duration.ofMinutes(3)));
        keycloak.start();

        // Configure the realm/client/users with the dedicated Terraform container (one-shot).
        keycloakConfig = new GenericContainer<>(image("demoapp-keycloak-terraform"))
                .withNetwork(network)
                .withEnv("TF_VAR_keycloak_url", "http://keycloak:8080")
                .withEnv("TF_VAR_keycloak_base_path", "/login")
                .withEnv("TF_VAR_keycloak_admin_user", "tmpadmin")
                .withEnv("TF_VAR_keycloak_admin_password", "admin")
                .withEnv("TF_VAR_realm", REALM)
                .withEnv("TF_VAR_client_id", "map-app")
                .withEnv("TF_VAR_redirect_uris",
                        "[\"" + baseUrl() + "/*\",\"http://localhost:5173/*\"]")
                .withStartupCheckStrategy(new OneShotStartupCheckStrategy().withTimeout(Duration.ofMinutes(3)));
        keycloakConfig.start();

        mapApi = new GenericContainer<>(image("map-api"))
                .withNetwork(network)
                .withNetworkAliases("map-api")
                .withExposedPorts(8090)
                .withEnv("PG_HOST", "postgres")
                .withEnv("PG_PORT", "5432")
                .withEnv("PG_USER", "mapservice")
                .withEnv("PG_PASSWORD", "mapservice")
                .withEnv("PG_DATABASE", DB_NAME)
                .withEnv("DATABASE_URL_SUFFIX", "?sslmode=disable&search_path=mapservice")
                .withEnv("PORT", "8090")
                .withEnv("KEYCLOAK_ISSUER", issuer())
                .withEnv("KEYCLOAK_JWKS_URL",
                        "http://keycloak:8080/login/realms/" + REALM + "/protocol/openid-connect/certs")
                .withEnv("CORS_ALLOWED_ORIGINS", baseUrl())
                .waitingFor(Wait.forHttp("/public/health/ready").forPort(8090)
                        .withStartupTimeout(Duration.ofMinutes(2)));
        mapApi.start();

        mapFrontend = new GenericContainer<>(image("map-frontend"))
                .withNetwork(network)
                .withNetworkAliases("map-frontend")
                .withExposedPorts(8080)
                .withEnv("VITE_API_BASE", "")
                .withEnv("VITE_OIDC_AUTHORITY", issuer())
                .withEnv("VITE_OIDC_CLIENT_ID", "map-app")
                .withEnv("VITE_OIDC_REDIRECT_URI", baseUrl() + "/")
                .waitingFor(Wait.forHttp("/").forPort(8080)
                        .withStartupTimeout(Duration.ofMinutes(1)));
        mapFrontend.start();

        gateway = new GenericContainer<>(DockerImageName.parse("nginx:1.27-alpine"))
                .withNetwork(network)
                .withNetworkAliases("gateway")
                .withCopyFileToContainer(
                        MountableFile.forClasspathResource("gateway/nginx.conf"),
                        "/etc/nginx/conf.d/default.conf")
                .withExposedPorts(80)
                .waitingFor(Wait.forHttp("/").forPort(80).forStatusCode(200)
                        .withStartupTimeout(Duration.ofMinutes(1)));
        gateway.setPortBindings(List.of(GATEWAY_PORT + ":80"));
        gateway.start();
    }

    private DockerImageName image(String repository) {
        return DockerImageName.parse(repository + ":" + IMAGE_TAG);
    }

    private static String env(String key, String fallback) {
        String value = System.getenv(key);
        return value != null && !value.isBlank() ? value : fallback;
    }
}
