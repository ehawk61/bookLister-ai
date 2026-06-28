package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIVersion(t *testing.T) {
	tests := []struct {
		name            string
		headerValue     string
		setHeader       bool
		wantCtxVersion  string
		wantRespVersion string
	}{
		{
			name:            "header provided",
			headerValue:     "2.0",
			setHeader:       true,
			wantCtxVersion:  "2.0",
			wantRespVersion: "2.0",
		},
		{
			name:            "header missing",
			setHeader:       false,
			wantCtxVersion:  "1.0",
			wantRespVersion: "1.0",
		},
		{
			name:            "empty header defaults",
			headerValue:     "",
			setHeader:       true,
			wantCtxVersion:  "1.0",
			wantRespVersion: "1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctxVersion string

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctxVersion = GetAPIVersion(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			handler := APIVersion(next)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.setHeader {
				req.Header.Set("X-API-Version", tt.headerValue)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if ctxVersion != tt.wantCtxVersion {
				t.Errorf("context API version = %q, want %q", ctxVersion, tt.wantCtxVersion)
			}

			if got := w.Header().Get("X-API-Version"); got != tt.wantRespVersion {
				t.Errorf("response X-API-Version = %q, want %q", got, tt.wantRespVersion)
			}
		})
	}
}
