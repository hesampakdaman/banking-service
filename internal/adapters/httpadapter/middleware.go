package httpadapter

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs the details of incoming HTTP requests and responses.
func LoggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Attach request details to logger
		reqLogger := logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		reqLogger.Info("Handled request",
			"status", rw.statusCode,
			"duration", time.Since(start),
		)
	})
}

// responseWriter is a wrapper to capture status codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader override to capture statusCode when set
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
