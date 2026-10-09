/*
 * X-Request-ID header plugin for Traefik.
 *
 * Copyright 2026 Digital H GmbH
 * Copyright 2023 M.D. Klapwijk
 * Copyright 2020 Vladimir Buyanov
 * Copyright 2020 Felipe Martínez
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package traefik_plugin_request_id

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func init() {
	uuid.EnableRandPool()
}

const defaultHeader = "X-Request-ID"
const defaultEnabled = true

type Config struct {
	Enabled           bool   `json:"enabled,omitempty"`
	HeaderName        string `json:"headerName,omitempty"`
	AddResponseHeader bool   `json:"addResponseHeader,omitempty"`

	// If true, the plugin will (try to) never prevent requests from succeeding, and instead only logs errors if
	// something goes wrong.
	FailSafe bool `json:"failSafe,omitempty"`

	GenerateUUID func() (uuid.UUID, error) `json:"-"`
}

func CreateConfig() *Config {
	return &Config{
		HeaderName:   defaultHeader,
		Enabled:      defaultEnabled,
		GenerateUUID: uuid.NewRandom,
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !config.Enabled || request.Header.Get(config.HeaderName) != "" {
			next.ServeHTTP(writer, request)
			return
		}

		reqUUID, err := config.GenerateUUID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate UUID: %v\n", err.Error())

			if config.FailSafe {
				next.ServeHTTP(writer, request)
			} else {
				http.Error(writer, fmt.Sprintf("Fatal error: HTTP server plugin %v: Failed to generate UUID: %v",
					name, err.Error()), http.StatusInternalServerError)
			}

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
