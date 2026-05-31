package http

import (
	"net/http"
	"time"

	"github.com/alekparkhomenko/investor/platform/pkg/logger"
	"github.com/gorilla/mux"
)

// statusRecorder captures the HTTP status code written by a handler.
type statusRecorder struct {
	http.ResponseWriter
	status   int
	wroteHeader bool
}

// WriteHeader implements http.ResponseWriter.WriteHeader.
func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(code)
}

// Write implements http.ResponseWriter.Write.
func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// LoggingMiddleware logs each HTTP request with method, path, status, and duration.
func LoggingMiddleware(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			log.Info(r.Context(), "http request", logger.Fields{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rec.status,
				"duration_ms": time.Since(start).Milliseconds(),
			})
		})
	}
}

// CORSMiddleware adds CORS headers for development and handles preflight requests.
func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 error.
func RecoveryMiddleware(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error(r.Context(), "panic recovered", logger.Fields{
						"method": r.Method,
						"path":   r.URL.Path,
						"panic":  rec,
					})
					errorResponse(w, http.StatusInternalServerError, "internal_error", "Internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
