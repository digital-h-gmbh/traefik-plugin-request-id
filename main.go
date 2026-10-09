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

func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if config.Enabled && request.Header.Get(config.HeaderName) == "" {
			value := uuid.Must(uuid.NewRandom()).String()
			request.Header.Add(config.HeaderName, value)
			if config.AddResponseHeader {
				writer.Header().Add(config.HeaderName, value)
			}
		}

		next.ServeHTTP(writer, request)
	}), nil
}
