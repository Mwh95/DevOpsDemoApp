package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "github.com/demoapp/map-service/internal/adapters/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envDatabaseURL       = "DATABASE_URL"
	envPGHost            = "PG_HOST"
	envPGPort            = "PG_PORT"
	envPGUser            = "PG_USER"
	envPGPassword        = "PG_PASSWORD"
	envPGDatabase        = "PG_DATABASE"
	envDatabaseURLSuffix = "DATABASE_URL_SUFFIX"
	envKeycloakIssuer    = "KEYCLOAK_ISSUER"
	envStaticDir         = "STATIC_DIR"
	envPort              = "PORT"

	defaultPGHost            = "localhost"
	defaultPGPort            = "5432"
	defaultPGUser            = "mapservice"
	defaultPGDatabase        = "MapMarkerDb"
	defaultHTTPClientTimeout = 10 * time.Second
	defaultServerPort        = "8090"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv(envDatabaseURL)
	if dbURL == "" {
		host := getEnv(envPGHost, defaultPGHost)
		port := getEnv(envPGPort, defaultPGPort)
		user := getEnv(envPGUser, defaultPGUser)
		pass := mustGetEnv(envPGPassword)
		dbName := getEnv(envPGDatabase, defaultPGDatabase)
		dbURLSuffix := mustGetEnv(envDatabaseURLSuffix)
		dbURL = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbName + dbURLSuffix
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	issuer := mustGetEnv(envKeycloakIssuer)
	auth, err := httpadapter.NewKeycloakJWKSVerifier(issuer, &http.Client{Timeout: defaultHTTPClientTimeout})
	if err != nil {
		log.Fatalf("oidc: %v", err)
	}

	staticDir := os.Getenv(envStaticDir)
	srv, err := httpadapter.NewServer(auth, pool, staticDir)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	addr := os.Getenv(envPort)
	if addr == "" {
		addr = defaultServerPort
	}
	if addr[0] != ':' {
		addr = ":" + addr
	}
	log.Printf("listening on %s", addr)
	if err := srv.Run(addr); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("missing required env var %s", key)
	return ""
}
