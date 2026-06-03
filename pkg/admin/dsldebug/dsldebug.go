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

// Package dsldebug wraps the OAP admin-server `dsl-debugging` feature module: a
// sample-based live debugger for MAL / LAL / OAL rules. The deeply-nested per-stage
// capture payloads are passed through as generic JSON for display; only the request
// parameters and client-side limits are modeled.
package dsldebug

import (
	"context"
	"net/http"
	"net/url"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

const (
	// MaxRecordCap mirrors OAP's SessionLimits.MAX_RECORD_CAP.
	MaxRecordCap = 100
	// MaxRetentionMillis mirrors OAP's SessionLimits.MAX_RETENTION_MILLIS (1 hour).
	MaxRetentionMillis = 60 * 60 * 1000
)

// Catalogs accepted by a debug session.
var Catalogs = []string{"otel-rules", "log-mal-rules", "telegraf-rules", "lal", "oal"}

// StartArgs holds the inputs of POST /dsl-debugging/session. Catalog, Name, RuleName
// and ClientID are mandatory query params; RecordCap / RetentionMillis are optional and
// sent as a JSON body only when set. Granularity (LAL only) is sent as a query param.
type StartArgs struct {
	ClientID        string
	Catalog         string
	Name            string
	RuleName        string
	Granularity     string
	RecordCap       int
	RetentionMillis int
}

// StartSession opens a debug capture session.
func StartSession(ctx context.Context, a *StartArgs) (any, error) {
	query := url.Values{
		"catalog":  []string{a.Catalog},
		"name":     []string{a.Name},
		"ruleName": []string{a.RuleName},
		"clientId": []string{a.ClientID},
	}
	if a.Granularity != "" {
		query.Set("granularity", a.Granularity)
	}

	body := map[string]any{}
	if a.RecordCap > 0 {
		body["recordCap"] = a.RecordCap
	}
	if a.RetentionMillis > 0 {
		body["retentionMillis"] = a.RetentionMillis
	}
	var in any
	if len(body) > 0 {
		in = body
	}

	var out any
	err := client.SendJSON(ctx, http.MethodPost, "/dsl-debugging/session", query, in, &out)
	return out, err
}

// GetSession polls a session's captured records (GET /dsl-debugging/session/{id}).
func GetSession(ctx context.Context, id string) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/dsl-debugging/session/"+url.PathEscape(id), nil, &out)
	return out, err
}

// StopSession stops a session (POST /dsl-debugging/session/{id}/stop). Idempotent.
func StopSession(ctx context.Context, id string) (any, error) {
	var out any
	err := client.SendJSON(ctx, http.MethodPost, "/dsl-debugging/session/"+url.PathEscape(id)+"/stop", nil, nil, &out)
	return out, err
}

// ListSessions lists the active sessions (GET /dsl-debugging/sessions).
func ListSessions(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/dsl-debugging/sessions", nil, &out)
	return out, err
}

// Status returns the module health snapshot (GET /dsl-debugging/status).
func Status(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/dsl-debugging/status", nil, &out)
	return out, err
}
