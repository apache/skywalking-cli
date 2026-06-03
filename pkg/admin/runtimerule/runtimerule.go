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

// Package runtimerule wraps the OAP admin-server `receiver-runtime-rule` feature
// module: hot-update of MAL / LAL rule files without restarting OAP, plus read access
// to the live and bundled rule state.
package runtimerule

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

// Catalogs accepted by the runtime-rule endpoints.
var Catalogs = []string{"otel-rules", "log-mal-rules", "telegraf-rules", "lal"}

// ApplyResult is the JSON envelope returned by addOrUpdate / inactivate / delete.
type ApplyResult struct {
	ApplyStatus string `json:"applyStatus"`
	Catalog     string `json:"catalog"`
	Name        string `json:"name"`
	Message     string `json:"message"`
}

// Rule is the result of a single-rule GET. Metadata comes from X-Sw-* response
// headers; Content is the raw YAML body. NotModified is set when a conditional GET
// (If-None-Match) returns 304.
type Rule struct {
	Status      string `json:"status"`
	Source      string `json:"source"`
	ContentHash string `json:"contentHash"`
	UpdateTime  int64  `json:"updateTime"`
	ETag        string `json:"etag"`
	Content     string `json:"content,omitempty"`
	NotModified bool   `json:"notModified,omitempty"`
}

// List returns the live rule state per node (GET /runtime/rule/list[?catalog=]).
func List(ctx context.Context, catalog string) (any, error) {
	var query url.Values
	if catalog != "" {
		query = url.Values{"catalog": []string{catalog}}
	}
	var out any
	err := client.GetJSON(ctx, "/runtime/rule/list", query, &out)
	return out, err
}

// ListBundled returns the static (bundled) rule twins for a catalog
// (GET /runtime/rule/bundled?catalog=&withContent=).
func ListBundled(ctx context.Context, catalog string, withContent bool) (any, error) {
	query := url.Values{
		"catalog":     []string{catalog},
		"withContent": []string{strconv.FormatBool(withContent)},
	}
	var out any
	err := client.GetJSON(ctx, "/runtime/rule/bundled", query, &out)
	return out, err
}

// Get fetches a single rule (GET /runtime/rule?catalog=&name=[&source=]). It negotiates
// raw YAML and reads metadata from the X-Sw-* headers. When ifNoneMatch is set and the
// server replies 304, the returned Rule has NotModified=true and an empty Content.
func Get(ctx context.Context, catalog, name, source, ifNoneMatch string) (*Rule, error) {
	query := url.Values{"catalog": []string{catalog}, "name": []string{name}}
	if source != "" {
		query.Set("source", source)
	}
	headers := map[string]string{}
	if ifNoneMatch != "" {
		headers["If-None-Match"] = ifNoneMatch
	}

	resp, err := client.Do(ctx, &client.Request{
		Method:  http.MethodGet,
		Path:    "/runtime/rule",
		Query:   query,
		Accept:  "application/x-yaml",
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotModified {
		return &Rule{
			NotModified: true,
			ETag:        resp.Header.Get("ETag"),
			ContentHash: resp.Header.Get("X-Sw-Content-Hash"),
			Status:      headerOr(resp, "X-Sw-Status", "n/a"),
		}, nil
	}
	if !resp.IsSuccess() {
		return nil, client.ParseError(resp)
	}
	updateTime, _ := strconv.ParseInt(resp.Header.Get("X-Sw-Update-Time"), 10, 64)
	return &Rule{
		Status:      headerOr(resp, "X-Sw-Status", "n/a"),
		Source:      headerOr(resp, "X-Sw-Source", "runtime"),
		ContentHash: resp.Header.Get("X-Sw-Content-Hash"),
		UpdateTime:  updateTime,
		ETag:        resp.Header.Get("ETag"),
		Content:     string(resp.Body),
	}, nil
}

// AddOrUpdate pushes a new or updated rule as raw YAML
// (POST /runtime/rule/addOrUpdate?catalog=&name=[&allowStorageChange=][&force=]).
func AddOrUpdate(ctx context.Context, catalog, name, body string, allowStorageChange, force bool) (*ApplyResult, error) {
	query := url.Values{"catalog": []string{catalog}, "name": []string{name}}
	if allowStorageChange {
		query.Set("allowStorageChange", "true")
	}
	if force {
		query.Set("force", "true")
	}
	resp, err := client.Do(ctx, &client.Request{
		Method:      http.MethodPost,
		Path:        "/runtime/rule/addOrUpdate",
		Query:       query,
		Body:        strings.NewReader(body),
		ContentType: "text/plain",
	})
	if err != nil {
		return nil, err
	}
	return applyResult(resp)
}

// Inactivate turns a rule off (POST /runtime/rule/inactivate?catalog=&name=).
func Inactivate(ctx context.Context, catalog, name string) (*ApplyResult, error) {
	query := url.Values{"catalog": []string{catalog}, "name": []string{name}}
	resp, err := client.Do(ctx, &client.Request{Method: http.MethodPost, Path: "/runtime/rule/inactivate", Query: query})
	if err != nil {
		return nil, err
	}
	return applyResult(resp)
}

// Delete removes a rule (POST /runtime/rule/delete?catalog=&name=[&mode=revertToBundled]).
// The server enforces a two-step gate: deleting an active rule returns a 409
// requires_inactivate_first.
func Delete(ctx context.Context, catalog, name, mode string) (*ApplyResult, error) {
	query := url.Values{"catalog": []string{catalog}, "name": []string{name}}
	if mode != "" {
		query.Set("mode", mode)
	}
	resp, err := client.Do(ctx, &client.Request{Method: http.MethodPost, Path: "/runtime/rule/delete", Query: query})
	if err != nil {
		return nil, err
	}
	return applyResult(resp)
}

// Dump returns a tar.gz snapshot of all rules, or one catalog's
// (GET /runtime/rule/dump[/{catalog}]).
func Dump(ctx context.Context, catalog string) ([]byte, error) {
	path := "/runtime/rule/dump"
	if catalog != "" {
		path += "/" + url.PathEscape(catalog)
	}
	resp, err := client.Do(ctx, &client.Request{Method: http.MethodGet, Path: path, Accept: "application/gzip"})
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, client.ParseError(resp)
	}
	return resp.Body, nil
}

func applyResult(resp *client.Response) (*ApplyResult, error) {
	if !resp.IsSuccess() {
		// Non-2xx still carries the ApplyResult envelope; ParseError lifts its
		// applyStatus/message so callers can switch on the semantic code.
		return nil, client.ParseError(resp)
	}
	var out ApplyResult
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			return nil, err
		}
	}
	return &out, nil
}

func headerOr(resp *client.Response, key, fallback string) string {
	if v := resp.Header.Get(key); v != "" {
		return v
	}
	return fallback
}
