package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Common API error codes (see SwapKit docs).
var (
	ErrNoRoutesFound            = errors.New("noRoutesFound")
	ErrBlackListAsset           = errors.New("blackListAsset")
	ErrAPIKeyInvalid            = errors.New("apiKeyInvalid")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrInvalidRequest           = errors.New("invalidRequest")
	ErrInsufficientBalance      = errors.New("insufficientBalance")
	ErrInsufficientAllowance    = errors.New("insufficientAllowance")
	ErrUnableToBuildTransaction = errors.New("unableToBuildTransaction")
)

// APIError represents a structured error from the SwapKit API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Data       map[string]interface{}
	body       []byte
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("swapkit api: %s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("swapkit api: %s", e.Message)
}

func (e *APIError) Unwrap() error {
	switch e.Code {
	case "noRoutesFound":
		return ErrNoRoutesFound
	case "blackListAsset":
		return ErrBlackListAsset
	case "apiKeyInvalid":
		return ErrAPIKeyInvalid
	case "unauthorized":
		return ErrUnauthorized
	case "invalidRequest":
		return ErrInvalidRequest
	case "insufficientBalance":
		return ErrInsufficientBalance
	case "insufficientAllowance":
		return ErrInsufficientAllowance
	case "unableToBuildTransaction":
		return ErrUnableToBuildTransaction
	default:
		return nil
	}
}

// apiErrorBody is the JSON shape of error responses.
type apiErrorBody struct {
	Message string                 `json:"message"`
	Error   string                 `json:"error"`
	Data    map[string]interface{} `json:"data"`
}

// parseAPIError builds an APIError from an HTTP response and body.
func parseAPIError(resp *http.Response, body []byte) *APIError {
	out := &APIError{StatusCode: resp.StatusCode, Message: resp.Status, body: body}
	if len(body) == 0 {
		return out
	}
	var parsed apiErrorBody
	if err := json.Unmarshal(body, &parsed); err != nil {
		out.Message = string(body)
		return out
	}
	if parsed.Message != "" {
		out.Message = parsed.Message
	}
	out.Code = parsed.Error
	out.Data = parsed.Data
	return out
}
