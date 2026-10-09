package traefik_plugin_request_id

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func init() {
	uuid.EnableRandPool()
}

const defaultHeader = "X-Request-ID"
const defaultEnabled = true

type Config struct {
	HeaderName        string `json:"headerName,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	AddResponseHeader bool   `json:"addResponseHeader,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		HeaderName: defaultHeader,
		Enabled:    defaultEnabled,
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !config.Enabled || request.Header.Get(config.HeaderName) != "" {
			next.ServeHTTP(writer, request)
			return
		}

		reqUUID, err := uuid.NewRandom()
		if err != nil {
			http.Error(writer, fmt.Sprintf("HTTP server plugin %v: Failed to generate UUID: %v",
				name, err.Error()), http.StatusInternalServerError)
			return
		}

		reqUUIDStr := reqUUID.String()
		request.Header.Add(config.HeaderName, reqUUIDStr)

		if config.AddResponseHeader {
			writer.Header().Add(config.HeaderName, reqUUIDStr)
		}

		next.ServeHTTP(writer, request)
	}), nil
}
