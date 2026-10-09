// Based on https://github.com/traefik/plugindemo/blob/master/demo_test.go

package traefik_plugin_request_id_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/digital-h-gmbh/traefik-plugin-request-id"
	"github.com/google/uuid"
)

func TestDefaultOptions(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := plugin.New(ctx, next, cfg, "plugin-x-request-id")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	hdr := req.Header.Get(cfg.HeaderName)
	t.Log("header value:", hdr)
	if err := uuid.Validate(hdr); err != nil {
		t.Fatalf("header %v contains invalid UUID: %v", cfg.HeaderName, err)
	}
}
