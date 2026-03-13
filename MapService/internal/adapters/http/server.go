package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/demoapp/map-service/internal/adapters/repository"
	"github.com/demoapp/map-service/internal/usecases"
	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server configures and runs the map API HTTP server.
type Server struct {
	Router *chi.Mux
}

type healthResponse struct {
	Status string `json:"status"`
}

// NewServer builds the server with auth, DB, and routes.
func NewServer(auth *KeycloakJWKSVerifier, pool *pgxpool.Pool, _ string) (*Server, error) {
	r := chi.NewRouter()
	r.Use(chimid.RequestID, chimid.RealIP, chimid.Logger, chimid.Recoverer)

	repo := repository.NewPostgresMarkerRepository(pool)
	uc := usecases.NewMarkerUseCases(repo)
	handler := NewMarkerHandler(uc)

	r.Get("/public/health/ready", func(w http.ResponseWriter, r *http.Request) {
		writeReadinessJSON(w, r, pool)
	})
	r.Get("/public/health/live", func(w http.ResponseWriter, r *http.Request) {
		writeHealthJSON(w)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Get("/markers", handler.List)
		r.Post("/markers", handler.Create)
		r.Get("/markers/{id}", handler.Get)
		r.Put("/markers/{id}", handler.Update)
		r.Delete("/markers/{id}", handler.Delete)
	})

	return &Server{Router: r}, nil
}

// Run starts the HTTP server on the given addr. If addr is empty, uses PORT env or ":8090".
func (s *Server) Run(addr string) error {
	if addr == "" {
		addr = os.Getenv("PORT")
		if addr == "" {
			addr = ":8090"
		}
		if addr[0] != ':' && addr[0] != '0' {
			addr = ":" + addr
		}
	}
	return http.ListenAndServe(addr, s.Router)
}

func writeHealthJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "UP"})
}

func writeReadinessJSON(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	w.Header().Set("Content-Type", "application/json")
	if pool == nil || pool.Ping(r.Context()) != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "DOWN"})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(healthResponse{Status: "UP"})
}
