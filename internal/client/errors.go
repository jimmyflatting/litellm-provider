package client

import (
	"encoding/json"
	"fmt"
)

// Error types
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIError struct {
	StatusCode int
	Err        error
	Response   *ErrorResponse
}

func (e *APIError) Error() string {
	if e.Response != nil {
		return fmt.Sprintf("API error (status %d): %s - %s", e.StatusCode, e.Response.Code, e.Response.Message)
	}
	return fmt.Sprintf("API error (status %d): %v", e.StatusCode, e.Err)
}

func parseError(statusCode int, body []byte) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed to parse error response: %v", err),
		}
	}
	return &APIError{
		StatusCode: statusCode,
		Response:   &errResp,
	}
}
