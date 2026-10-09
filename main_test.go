// Based on https://github.com/traefik/plugindemo/blob/master/demo_test.go

package traefik_plugin_request_id_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	plugin "github.com/digital-h-gmbh/traefik-plugin-request-id"
	"github.com/google/uuid"
)

func TestDefaultOptions(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("response was not OK")
	}

	hdr := req.Header.Get(cfg.HeaderName)
	t.Log("header value:", hdr)
	if err := uuid.Validate(hdr); err != nil {
		t.Fatalf("header %v contains invalid UUID: %v", cfg.HeaderName, err)
	}

	respHdr := resp.Header.Get(cfg.HeaderName)
	if respHdr != "" {
		t.Fatalf("response header %v should have been unset", cfg.HeaderName)
	}
}

func TestDisabledPlugin(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	cfg.Enabled = false

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("response was not OK")
	}

	hdr := req.Header.Get(cfg.HeaderName)
	if hdr != "" {
		t.Fatalf("header %v was set, but should not be", cfg.HeaderName)
	}
}

func TestCustomHeaderName(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	cfg.HeaderName = "X-My-Custom-ID"

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("response was not OK")
	}

	hdr := req.Header.Get(cfg.HeaderName)
	t.Log("header value:", hdr)
	if err := uuid.Validate(hdr); err != nil {
		t.Fatalf("header %v contains invalid UUID: %v", cfg.HeaderName, err)
	}
}

func TestResponseHeader(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	cfg.AddResponseHeader = true

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("response was not OK")
	}

	respHdr := resp.Header.Get(cfg.HeaderName)
	t.Log("response header value:", respHdr)
	if err := uuid.Validate(respHdr); err != nil {
		t.Fatalf("response header %v contains invalid UUID: %v", cfg.HeaderName, err)
	}
}

func TestExistingRequestHeaderIsPreserved(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

	handler, err := plugin.New(ctx, next, cfg, "plugin-x-request-id")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(cfg.HeaderName, "hello world")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("response was not OK")
	}

	hdr := req.Header.Get(cfg.HeaderName)
	t.Log("header value:", hdr)
	if hdr != "hello world" {
		t.Fatalf("request header %v was not preserved", cfg.HeaderName)
	}
}

func TestErrorWithoutFailSafe(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	cfg.GenerateUUID = func() (uuid.UUID, error) {
		return uuid.Nil, fmt.Errorf("no UUIDs for you today")
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("response body:", string(respBody))
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response should have been an error")
	}
}

func TestErrorWithFailSafe(t *testing.T) {
	ctx := context.Background()

	cfg := plugin.CreateConfig()
	cfg.FailSafe = true

	wasUUIDGenerated := false
	cfg.GenerateUUID = func() (uuid.UUID, error) {
		wasUUIDGenerated = true
		return uuid.Nil, fmt.Errorf("no UUIDs for you today")
	}

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	})

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

	resp := recorder.Result()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("response body:", string(respBody))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("response should have been successful")
	}

	hdr := req.Header.Get(cfg.HeaderName)
	if hdr != "" {
		t.Fatalf("header %v was set, but should not be", cfg.HeaderName)
	}

	if !wasUUIDGenerated {
		t.Fatal("UUID generation was never attempted!")
	}
}
