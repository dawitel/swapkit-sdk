package client

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Code: "noRoutesFound", Message: "No routes found"}
	if s := e.Error(); s == "" || len(s) < 10 {
		t.Errorf("Error() = %q", s)
	}
	e2 := &APIError{Message: "Bad request"}
	if e2.Error() == "" {
		t.Error("Error() should not be empty when Message is set")
	}
}

func TestAPIError_Unwrap_AllSentinels(t *testing.T) {
	tests := []struct {
		code string
		want error
	}{
		{"noRoutesFound", ErrNoRoutesFound},
		{"blackListAsset", ErrBlackListAsset},
		{"apiKeyInvalid", ErrAPIKeyInvalid},
		{"unauthorized", ErrUnauthorized},
		{"invalidRequest", ErrInvalidRequest},
		{"insufficientBalance", ErrInsufficientBalance},
		{"insufficientAllowance", ErrInsufficientAllowance},
		{"unableToBuildTransaction", ErrUnableToBuildTransaction},
		{"unknownCode", nil},
		{"", nil},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			e := &APIError{Code: tt.code}
			got := e.Unwrap()
			if tt.want != nil {
				if got != tt.want {
					t.Errorf("Unwrap() = %v, want %v", got, tt.want)
				}
				if !errors.Is(e, tt.want) {
					t.Errorf("errors.Is(apiErr, sentinel) = false")
				}
			} else if got != nil {
				t.Errorf("Unwrap() = %v, want nil", got)
			}
		})
	}
}

func TestParseAPIError_EmptyBody(t *testing.T) {
	resp := &http.Response{StatusCode: 500, Status: "500 Internal Server Error"}
	e := parseAPIError(resp, nil)
	if e.StatusCode != 500 {
		t.Errorf("StatusCode = %d", e.StatusCode)
	}
	// When body is nil, Message comes from resp.Status (or is left as Status from initial set)
	if e.Message == "" {
		t.Errorf("Message should be set")
	}
}

func TestParseAPIError_NonJSONBody(t *testing.T) {
	resp := &http.Response{StatusCode: 400}
	e := parseAPIError(resp, []byte("plain text error"))
	if e.Message != "plain text error" {
		t.Errorf("Message = %q", e.Message)
	}
}

func TestParseAPIError_JSONBody(t *testing.T) {
	resp := &http.Response{StatusCode: 404}
	body := []byte(`{"message":"Not found","error":"noRoutesFound","data":{"sellAsset":"X"}}`)
	e := parseAPIError(resp, body)
	if e.Code != "noRoutesFound" || e.Message != "Not found" {
		t.Errorf("got Code=%q Message=%q", e.Code, e.Message)
	}
	if e.Data == nil || e.Data["sellAsset"] != "X" {
		t.Errorf("Data = %v", e.Data)
	}
}
