package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errResp        ErrorResponse
		wantStatus     int
		wantDetailsKey bool
	}{
		{
			name:       "400 with details",
			statusCode: http.StatusBadRequest,
			errResp: ErrorResponse{
				ErrorCode:        "MISSING_CORRELATION_ID",
				ErrorDescription: "X-Correlation-ID header is required",
				Details:          "Provide a valid UUIDv4",
			},
			wantStatus:     http.StatusBadRequest,
			wantDetailsKey: true,
		},
		{
			name:       "404 without details",
			statusCode: http.StatusNotFound,
			errResp: ErrorResponse{
				ErrorCode:        "BOOK_NOT_FOUND",
				ErrorDescription: "Book does not exist",
			},
			wantStatus:     http.StatusNotFound,
			wantDetailsKey: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteError(w, tt.statusCode, tt.errResp)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			var got map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if got["error_code"] != tt.errResp.ErrorCode {
				t.Errorf("error_code = %q, want %q", got["error_code"], tt.errResp.ErrorCode)
			}
			if got["error_description"] != tt.errResp.ErrorDescription {
				t.Errorf("error_description = %q, want %q", got["error_description"], tt.errResp.ErrorDescription)
			}

			_, hasDetails := got["details"]
			if hasDetails != tt.wantDetailsKey {
				t.Errorf("details present = %v, want %v", hasDetails, tt.wantDetailsKey)
			}
			if tt.wantDetailsKey && got["details"] != tt.errResp.Details {
				t.Errorf("details = %q, want %q", got["details"], tt.errResp.Details)
			}
		})
	}
}
