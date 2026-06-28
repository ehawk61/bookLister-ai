package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	Details          string `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, statusCode int, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResp)
}
