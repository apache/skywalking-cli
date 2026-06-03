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

// Package client is the REST client for the OAP admin-server (default port 17128),
// the HTTP host that bundles the status, inspect, ui-management, dsl-debugging and
// runtime-rule feature modules. It is the REST counterpart of pkg/graphql/client
// and shares its TLS / basic-auth handling via pkg/transport.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/transport"
)

const (
	// DefaultAdminURL is the fallback admin-server REST base URL.
	DefaultAdminURL = "http://127.0.0.1:17128"
	// DefaultAdminPort is the admin-server REST port (admin-server `port`).
	DefaultAdminPort = "17128"

	defaultTimeout = 30 * time.Second

	contentTypeJSON = "application/json"
)

// DeriveAdminURL resolves the admin REST base URL. When adminURL is non-empty it is
// used verbatim (whitespace and trailing slash trimmed). Otherwise the admin URL is
// derived from the GraphQL base URL by reusing its scheme and host and swapping in the
// admin port, dropping any path such as `/graphql`. A non-parseable base URL falls
// back to DefaultAdminURL.
func DeriveAdminURL(baseURL, adminURL string) string {
	if s := strings.TrimRight(strings.TrimSpace(adminURL), "/"); s != "" {
		return s
	}
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || u.Hostname() == "" {
		return DefaultAdminURL
	}
	scheme := u.Scheme
	if scheme == "" {
		scheme = "http"
	}
	// net.JoinHostPort brackets IPv6 hosts, e.g. http://[::1]:17128.
	return scheme + "://" + net.JoinHostPort(u.Hostname(), DefaultAdminPort)
}

// BaseURL returns the admin REST base URL stored in the context (set by the
// interceptor), trailing slash trimmed, falling back to DefaultAdminURL.
func BaseURL(ctx context.Context) string {
	return strings.TrimRight(transport.GetValue(ctx, contextkey.AdminURL{}, DefaultAdminURL), "/")
}

// Request describes a single admin REST call.
type Request struct {
	Method      string
	Path        string // appended to the admin base URL, e.g. "/status/cluster/nodes"
	Query       url.Values
	Body        io.Reader
	ContentType string            // request Content-Type, when a body is sent
	Accept      string            // Accept header; defaults to application/json
	Headers     map[string]string // extra request headers (e.g. If-None-Match)
}

// Response is the raw result of an admin REST call, exposed so callers can read
// status-dependent headers (e.g. runtime-rule X-Sw-* / ETag) and bodies (tar.gz).
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	URL        string
}

// IsSuccess reports whether the response carries a 2xx status code.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// APIError is returned for non-2xx admin responses. It decodes the common admin
// error envelopes so callers can switch on the semantic Status code (e.g.
// requires_inactivate_first, session_not_found, cluster_view_split).
type APIError struct {
	StatusCode int
	URL        string
	Status     string // applyStatus / code / status from the JSON envelope
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.Body
	}
	if e.Status != "" {
		return fmt.Sprintf("admin API %s: HTTP %d (%s): %s", e.URL, e.StatusCode, e.Status, msg)
	}
	return fmt.Sprintf("admin API %s: HTTP %d: %s", e.URL, e.StatusCode, msg)
}

// ParseError builds an APIError from a non-2xx response, decoding the admin error
// envelopes: {applyStatus,message} (runtime-rule), {status,code,message}
// (dsl-debugging / oal) and {error} (inspect). A non-JSON body is kept as raw text.
func ParseError(resp *Response) *APIError {
	e := &APIError{StatusCode: resp.StatusCode, URL: resp.URL, Body: strings.TrimSpace(string(resp.Body))}
	var env struct {
		ApplyStatus string `json:"applyStatus"`
		Status      string `json:"status"`
		Code        string `json:"code"`
		Message     string `json:"message"`
		Error       string `json:"error"`
	}
	if json.Unmarshal(resp.Body, &env) == nil {
		switch {
		case env.ApplyStatus != "":
			e.Status = env.ApplyStatus
		case env.Code != "":
			e.Status = env.Code
		case env.Status != "":
			e.Status = env.Status
		}
		switch {
		case env.Message != "":
			e.Message = env.Message
		case env.Error != "":
			e.Message = env.Error
		}
	}
	return e
}

// Do performs req against the admin base URL and returns the raw Response for any
// HTTP status. Only transport-level failures are returned as an error; callers
// decide how to treat non-2xx (see Response.IsSuccess / ParseError). It reuses the
// shared TLS and basic-auth handling from pkg/transport.
func Do(ctx context.Context, req *Request) (*Response, error) {
	full := BaseURL(ctx) + req.Path
	if len(req.Query) > 0 {
		full += "?" + req.Query.Encode()
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, full, req.Body)
	if err != nil {
		return nil, err
	}
	accept := req.Accept
	if accept == "" {
		accept = contentTypeJSON
	}
	httpReq.Header.Set("Accept", accept)
	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	}
	if auth := transport.AuthHeader(ctx); auth != "" {
		httpReq.Header.Set("Authorization", auth)
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	httpClient := transport.HTTPClient(ctx)
	httpClient.Timeout = defaultTimeout
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: body, URL: full}, nil
}

// GetJSON issues a GET to path and decodes a JSON 2xx response into out (which may be
// nil to discard the body). A non-2xx response is returned as an *APIError.
func GetJSON(ctx context.Context, path string, query url.Values, out any) error {
	resp, err := Do(ctx, &Request{Method: http.MethodGet, Path: path, Query: query})
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return ParseError(resp)
	}
	if out == nil || len(resp.Body) == 0 {
		return nil
	}
	return json.Unmarshal(resp.Body, out)
}

// SendJSON issues method to path with an optional JSON body and decodes a JSON 2xx
// response into out (which may be nil). A non-2xx response is returned as an *APIError.
func SendJSON(ctx context.Context, method, path string, query url.Values, in, out any) error {
	req := &Request{Method: method, Path: path, Query: query}
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return err
		}
		req.Body = strings.NewReader(string(data))
		req.ContentType = contentTypeJSON
	}
	resp, err := Do(ctx, req)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return ParseError(resp)
	}
	if out == nil || len(resp.Body) == 0 {
		return nil
	}
	return json.Unmarshal(resp.Body, out)
}
