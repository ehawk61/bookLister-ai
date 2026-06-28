package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehawk61/bookLister-ai/internal/response"
)

func TestCorrelationID(t *testing.T) {
	const validUUID = "550e8400-e29b-41d4-a716-446655440000"
	const validUUIDUpper = "550E8400-E29B-41D4-A716-446655440000"
	const uuidV1 = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	tests := []struct {
		name           string
		headerValue    string
		setHeader      bool
		wantStatus     int
		wantErrorCode  string
		wantNextCalled bool
	}{
		{
			name:           "missing header",
			setHeader:      false,
			wantStatus:     http.StatusBadRequest,
			wantErrorCode:  "MISSING_CORRELATION_ID",
			wantNextCalled: false,
		},
		{
			name:           "empty header",
			headerValue:    "",
			setHeader:      true,
			wantStatus:     http.StatusBadRequest,
			wantErrorCode:  "MISSING_CORRELATION_ID",
			wantNextCalled: false,
		},
		{
			name:           "invalid format",
			headerValue:    "not-a-uuid",
			setHeader:      true,
			wantStatus:     http.StatusBadRequest,
			wantErrorCode:  "INVALID_CORRELATION_ID",
			wantNextCalled: false,
		},
		{
			name:           "UUID v1 rejected",
			headerValue:    uuidV1,
			setHeader:      true,
			wantStatus:     http.StatusBadRequest,
			wantErrorCode:  "INVALID_CORRELATION_ID",
			wantNextCalled: false,
		},
		{
			name:           "valid UUIDv4",
			headerValue:    validUUID,
			setHeader:      true,
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "valid UUIDv4 uppercase",
			headerValue:    validUUIDUpper,
			setHeader:      true,
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			var ctxCorrelationID string

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				ctxCorrelationID = GetCorrelationID(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			handler := CorrelationID(next)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.setHeader {
				req.Header.Set("X-Correlation-ID", tt.headerValue)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNextCalled)
			}

			if tt.wantErrorCode != "" {
				var errResp response.ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.ErrorCode != tt.wantErrorCode {
					t.Errorf("error_code = %q, want %q", errResp.ErrorCode, tt.wantErrorCode)
				}
			}

			if tt.wantNextCalled {
				if ctxCorrelationID != tt.headerValue {
					t.Errorf("context correlation ID = %q, want %q", ctxCorrelationID, tt.headerValue)
				}
				if got := w.Header().Get("X-Correlation-ID"); got != tt.headerValue {
					t.Errorf("response X-Correlation-ID = %q, want %q", got, tt.headerValue)
				}
			}
		})
	}
}
