package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer server.Close()

	conf := Config{server.URL, 10, 2, 10 * time.Second}

	wr := Run(ctx, conf)
	if len(wr) != conf.N {
		t.Errorf("%d должен быть равен %d", len(wr), conf.N)
	}
}
