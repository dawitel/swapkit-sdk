package api_test

import (
	"net/http/httptest"
	"testing"

	"github.com/dawitel/swapkit-sdk/client"
)

func newTestClient(t *testing.T, srv *httptest.Server) *client.Client {
	t.Helper()
	c, err := client.New(client.Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	return c
}
