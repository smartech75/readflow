package middleware

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/pkg/constant"
)

// Tracing is a middleware to trace HTTP request
func Tracing(nextRequestID func() string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), constant.ContextRequestID, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
