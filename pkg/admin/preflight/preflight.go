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

// Package preflight performs admin-server feature detection by reading the effective
// configuration dump, mirroring how Horizon UI degrades gracefully when an admin
// feature module is disabled or the admin host is unreachable.
package preflight

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

// Known admin feature-module selector keys (the first dotted segment of a config-dump
// key) and the environment variables that enable them.
const (
	ModuleAdminServer = "admin-server"
	ModuleStatus      = "status"
	ModuleInspect     = "inspect"
	ModuleUIManage    = "ui-management"
	ModuleDSLDebug    = "dsl-debugging"
	ModuleRuntimeRule = "receiver-runtime-rule"
)

// Module describes one admin feature module and how to enable it.
type Module struct {
	Name    string `json:"name"`
	EnvVar  string `json:"envVar"`
	Enabled bool   `json:"enabled"`
}

// Result is the outcome of a preflight check against the admin host.
type Result struct {
	AdminURL       string   `json:"adminURL"`
	AdminReachable bool     `json:"adminReachable"`
	Modules        []Module `json:"modules"`

	enabled map[string]bool
}

// IsEnabled reports whether a feature module (by selector key) is enabled.
func (r *Result) IsEnabled(module string) bool { return r.enabled[module] }

var knownModules = []struct{ name, env string }{
	{ModuleAdminServer, "SW_ADMIN_SERVER"},
	{ModuleStatus, "SW_STATUS"},
	{ModuleInspect, "SW_INSPECT"},
	{ModuleUIManage, "SW_UI_MANAGEMENT"},
	{ModuleDSLDebug, "SW_DSL_DEBUGGING"},
	{ModuleRuntimeRule, "SW_RECEIVER_RUNTIME_RULE"},
}

// Run reads GET /debugging/config/dump from the admin host and reports which feature
// modules are enabled. A module is considered enabled when any dotted config key
// starts with `<module>.`. When the dump cannot be fetched, AdminReachable is false,
// every module reports disabled, and the transport error is returned alongside the
// (still useful) Result so callers can surface the admin URL.
func Run(ctx context.Context) (*Result, error) {
	r := &Result{AdminURL: client.BaseURL(ctx), enabled: map[string]bool{}}

	var dump map[string]any
	err := client.GetJSON(ctx, "/debugging/config/dump", nil, &dump)
	if err == nil {
		r.AdminReachable = true
		for k := range dump {
			prefix := k
			if i := strings.IndexByte(k, '.'); i >= 0 {
				prefix = k[:i]
			}
			r.enabled[prefix] = true
		}
	}
	for _, m := range knownModules {
		r.Modules = append(r.Modules, Module{Name: m.name, EnvVar: m.env, Enabled: r.enabled[m.name]})
	}
	return r, err
}

// Explain enriches an admin-call error with operator-actionable context. A transport
// failure is reported as an unreachable admin host. A 404 with no recognizable JSON
// error envelope is reported as a likely-disabled module (the route is not registered),
// whereas a 404 that carries an error body (e.g. {"error":"not_found"}) is a real
// resource miss from an enabled module and is returned unchanged — as are all other API
// errors (400/409/421/...), which are already specific.
func Explain(ctx context.Context, err error, module, envVar string) error {
	if err == nil {
		return nil
	}
	adminURL := client.BaseURL(ctx)
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusNotFound && apiErr.Message == "" && apiErr.Status == "" {
			return fmt.Errorf("the `%s` admin feature module appears disabled on OAP "+
				"(HTTP 404 with no error body at %s); enable it with %s=default. original error: %w",
				module, adminURL, envVar, err)
		}
		return err
	}
	return fmt.Errorf("could not reach the OAP admin-server at %s; "+
		"verify --admin-url and that admin-server is enabled (SW_ADMIN_SERVER=default). "+
		"original error: %w", adminURL, err)
}
