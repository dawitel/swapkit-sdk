package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	swapkitapi "github.com/dawitel/swapkit-sdk/internal/api"
)

func TestGetSwapTo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/swapTo" || r.URL.Query().Get("sellAsset") != "ETH.ETH" {
			t.Errorf("path=%s query=%s", r.URL.Path, r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["BTC.BTC","SOL.SOL"]`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	list, err := swapkitapi.GetSwapTo(context.Background(), c, "ETH.ETH")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0] != "BTC.BTC" || list[1] != "SOL.SOL" {
		t.Errorf("list: %v", list)
	}
}
