package goreq

import (
	"context"
	"log/slog"
	"net/http"
)

type requestKey struct{}

// RequestParse adds the request parameters defined by T into the request context.
// The request parameters can be retrieved using GetRequest()
func RequestParse[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := parseParameters[T](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("Error parsing request parameters", slog.String("error", err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), requestKey{}, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequest retrieves the request parameters from the context, if any.
func GetRequest[T any](ctx context.Context) (*T, bool) {
	val, ok := ctx.Value(requestKey{}).(*T)
	return val, ok
}
