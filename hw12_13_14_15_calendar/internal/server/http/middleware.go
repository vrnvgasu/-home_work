package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		msg := fmt.Sprintf(
			"%s [%s] %s %s %s %d %s '%s'",
			ReadUserIP(r),
			time.Now().UTC(),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			time.Since(start),
			r.UserAgent(),
		)
		s.logger.Info(msg)
		s.logger.File(msg)
	})
}
