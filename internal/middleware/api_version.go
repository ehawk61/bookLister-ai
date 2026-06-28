package middleware

import (
	"context"
	"net/http"
)

const (
	apiVersionKey     contextKey = "apiVersion"
	defaultAPIVersion            = "1.0"
)

func APIVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get("X-API-Version")
		if version == "" {
			version = defaultAPIVersion
		}

		w.Header().Set("X-API-Version", version)
		ctx := context.WithValue(r.Context(), apiVersionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAPIVersion(ctx context.Context) string {
	v, _ := ctx.Value(apiVersionKey).(string)
	return v
}
