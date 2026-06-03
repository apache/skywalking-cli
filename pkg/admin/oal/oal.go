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

// Package oal wraps the read-only OAL listing endpoints (`/runtime/oal/*`) hosted by
// the OAP admin-server `dsl-debugging` feature module: the rule-picker source for the
// OAL live debugger. OAL hot-update is upstream-deferred; these endpoints are read-only.
package oal

import (
	"context"
	"net/http"
	"net/url"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

// ListFiles lists the loaded .oal file names (GET /runtime/oal/files).
func ListFiles(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/runtime/oal/files", nil, &out)
	return out, err
}

// GetFile returns the raw .oal text of a single file (GET /runtime/oal/files/{name}).
func GetFile(ctx context.Context, name string) (string, error) {
	resp, err := client.Do(ctx, &client.Request{
		Method: http.MethodGet,
		Path:   "/runtime/oal/files/" + url.PathEscape(name),
		Accept: "application/json, text/plain",
	})
	if err != nil {
		return "", err
	}
	if !resp.IsSuccess() {
		return "", client.ParseError(resp)
	}
	return string(resp.Body), nil
}

// ListSources lists the per-dispatcher OAL source catalog (GET /runtime/oal/rules).
func ListSources(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/runtime/oal/rules", nil, &out)
	return out, err
}

// GetSource returns one source's per-metric holder status
// (GET /runtime/oal/rules/{source}).
func GetSource(ctx context.Context, source string) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/runtime/oal/rules/"+url.PathEscape(source), nil, &out)
	return out, err
}
