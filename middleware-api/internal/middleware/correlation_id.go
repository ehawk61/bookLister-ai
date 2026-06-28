package middleware

import (
	"context"
	"net/http"
	"regexp"

	"github.com/ehawk61/bookLister-ai/internal/response"
)

type contextKey string

const correlationIDKey contextKey = "correlationID"

var uuidV4Re = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Correlation-ID")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "MISSING_CORRELATION_ID",
				ErrorDescription: "X-Correlation-ID header is required",
				Details:          "Provide a valid UUIDv4 in the X-Correlation-ID header",
			})
			return
		}

		if !uuidV4Re.MatchString(id) {
			response.WriteError(w, http.StatusBadRequest, response.ErrorResponse{
				ErrorCode:        "INVALID_CORRELATION_ID",
				ErrorDescription: "X-Correlation-ID must be a valid UUIDv4",
				Details:          id,
			})
			return
		}

		w.Header().Set("X-Correlation-ID", id)
		ctx := context.WithValue(r.Context(), correlationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCorrelationID(ctx context.Context) string {
	v, _ := ctx.Value(correlationIDKey).(string)
	return v
}
