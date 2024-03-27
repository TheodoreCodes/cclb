package loadbalancer

import (
	"cclb/log"
	"context"
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("Received request from %s", r.Host), map[string]any{
				"Method":     r.Method,
				"Host":       r.Host,
				"User-Agent": r.UserAgent(),
				"Accept":     r.Header.Get("Accept"),
			})
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func timeoutMiddleware(logger log.Logger, timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithDeadline(r.Context(), time.Now().Add(timeout))
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
