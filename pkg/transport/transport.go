// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package transport holds the connection bits shared by the GraphQL client and
// the admin REST client: TLS handling, basic-auth resolution, and context helpers.
package transport

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"net/http"

	"github.com/apache/skywalking-cli/pkg/contextkey"
)

// Insecure reports whether TLS certificate verification should be skipped,
// according to the `--insecure` global option stored in the context.
func Insecure(ctx context.Context) bool {
	return GetValue(ctx, contextkey.Insecure{}, false)
}

// HTTPClient builds an *http.Client honoring the `--insecure` option from context.
// When insecure is not set it returns a plain client, equivalent to http.DefaultClient.
func HTTPClient(ctx context.Context) *http.Client {
	if Insecure(ctx) {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402
		return &http.Client{Transport: customTransport}
	}
	return &http.Client{}
}

// AuthHeader resolves the value of the `Authorization` header from the global
// `--authorization`, `--username` and `--password` options stored in the context,
// returning an empty string when no credentials are configured. A raw
// `--authorization` takes precedence over `--username`/`--password`.
func AuthHeader(ctx context.Context) string {
	username := GetValue(ctx, contextkey.Username{}, "")
	password := GetValue(ctx, contextkey.Password{}, "")
	authorization := GetValue(ctx, contextkey.Authorization{}, "")

	if authorization == "" && username != "" && password != "" {
		authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	}
	return authorization
}

// GetValue safely extracts a typed value from the context, falling back to
// defaultValue when the key is absent or holds a different type.
func GetValue[T any](ctx context.Context, key any, defaultValue T) T {
	val := ctx.Value(key)
	if v, ok := val.(T); ok {
		return v
	}
	return defaultValue
}
