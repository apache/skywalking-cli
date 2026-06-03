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

// Package uitemplate wraps the OAP admin-server `ui-management` feature module: REST
// CRUD over dashboard templates. It is the replacement for the GraphQL
// UIConfigurationManagement template mutations retired in SkyWalking 11.0.0. There is
// no DELETE; templates are soft-disabled.
package uitemplate

import (
	"context"
	"net/http"
	"net/url"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

const basePath = "/ui-management/templates"

// Template is a dashboard template. Configuration is an opaque JSON-encoded string.
type Template struct {
	ID            string `json:"id"`
	Configuration string `json:"configuration"`
	Disabled      bool   `json:"disabled"`
}

// ChangeStatus is the acknowledgement of a create / update / disable write.
type ChangeStatus struct {
	ID      string `json:"id"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// List returns all templates (GET /ui-management/templates). When includingDisabled
// is true, soft-disabled templates are included as well.
func List(ctx context.Context, includingDisabled bool) ([]Template, error) {
	var query url.Values
	if includingDisabled {
		query = url.Values{"includingDisabled": []string{"true"}}
	}
	var out []Template
	err := client.GetJSON(ctx, basePath, query, &out)
	return out, err
}

// Get returns a single template by ID (GET /ui-management/templates/{id}).
func Get(ctx context.Context, id string) (*Template, error) {
	var out Template
	err := client.GetJSON(ctx, basePath+"/"+url.PathEscape(id), nil, &out)
	return &out, err
}

// Create adds a new template (POST /ui-management/templates). configuration is the
// JSON-encoded template body. Since OAP 11.0.0 (skywalking#13884) the id is required
// in the request body, so callers pass a client-generated id.
func Create(ctx context.Context, id, configuration string) (*ChangeStatus, error) {
	var out ChangeStatus
	body := map[string]string{"id": id, "configuration": configuration}
	err := client.SendJSON(ctx, http.MethodPost, basePath, nil, body, &out)
	return &out, err
}

// Update replaces an existing template (PUT /ui-management/templates).
func Update(ctx context.Context, id, configuration string) (*ChangeStatus, error) {
	var out ChangeStatus
	body := map[string]string{"id": id, "configuration": configuration}
	err := client.SendJSON(ctx, http.MethodPut, basePath, nil, body, &out)
	return &out, err
}

// Disable soft-disables a template (POST /ui-management/templates/{id}/disable). It is
// idempotent. The template row is preserved; only its disabled flag flips.
func Disable(ctx context.Context, id string) (*ChangeStatus, error) {
	out := ChangeStatus{ID: id, Status: true}
	path := basePath + "/" + url.PathEscape(id) + "/disable"
	err := client.SendJSON(ctx, http.MethodPost, path, nil, nil, &out)
	return &out, err
}
