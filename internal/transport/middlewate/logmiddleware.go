package middleware

import (
	"log/slog"
	"net/http"
)

type LoggMiddleware struct {
	Logger *slog.Logger
}

func (logcl *LoggMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logcl.Logger.Info("middleware log", "Method", r.Method, "URL", r.URL)
		next.ServeHTTP(w, r)
	})
}
