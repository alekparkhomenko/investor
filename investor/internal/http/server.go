package http

import (
	"context"
	_ "embed"
	"net/http"
	"time"

	"github.com/alekparkhomenko/investor/platform/pkg/logger"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed docs/swagger.json
var swaggerSpec []byte

// swaggerJSON serves the swagger.json file from embedded bytes.
func swaggerJSON(w http.ResponseWriter, r *http.Request) {
	if len(swaggerSpec) == 0 {
		http.Error(w, "swagger spec not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(swaggerSpec)
}

// HTTPServer represents HTTP server.
type HTTPServer struct {
	addr   string
	router *mux.Router
	server *http.Server
	log    *logger.Logger
}

// Config holds server configuration.
type Config struct {
	Address string
	Handler *Handler
	Logger  *logger.Logger
}

// NewHTTPServer creates a new HTTPServer.
func NewHTTPServer(cfg Config) *HTTPServer {
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/tickers", cfg.Handler.ListTickers).Methods(http.MethodGet)
	api.HandleFunc("/portfolio", cfg.Handler.GetPortfolio).Methods(http.MethodGet)
	api.HandleFunc("/portfolio", cfg.Handler.AddToPortfolio).Methods(http.MethodPost)
	api.HandleFunc("/portfolio/{ticker}", cfg.Handler.RemoveFromPortfolio).Methods(http.MethodDelete)

	// Swagger JSON
	router.HandleFunc("/swagger/doc.json", swaggerJSON)

	// Swagger UI
	router.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &HTTPServer{
		addr:   cfg.Address,
		router: router,
		server: server,
		log:    cfg.Logger,
	}
}

// Start starts the HTTP server.
func (s *HTTPServer) Start() error {
	s.log.Info(context.Background(), "starting HTTP server", logger.Fields{
		"component": "http-server",
		"addr":       s.addr,
	})

	return s.server.ListenAndServe()
}

// Stop stops the HTTP server gracefully.
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.log.Info(ctx, "stopping HTTP server", logger.Fields{
		"component": "http-server",
	})

	return s.server.Shutdown(ctx)
}