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

// Package inspect wraps the OAP admin-server `inspect` feature module: browse the
// metric catalog and enumerate the entities currently holding values for a metric.
// The entity rows carry an MQE-ready payload that pastes into a follow-up
// execExpression query on the public GraphQL surface.
package inspect

import (
	"context"
	"net/url"
	"strconv"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

// MaxLimit is the server-side hard cap on /inspect/entities rows.
const MaxLimit = 300

// Metric is a single entry of the metric catalog.
type Metric struct {
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Catalog         string   `json:"catalog"`
	ScopeID         int      `json:"scopeId"`
	Scope           string   `json:"scope"`
	ValueColumnName string   `json:"valueColumnName"`
	Downsamplings   []string `json:"downsamplings"`
}

// Metrics is the response of GET /inspect/metrics.
type Metrics struct {
	Metrics []Metric `json:"metrics"`
}

// MetricsOptions holds the optional filters of GET /inspect/metrics.
type MetricsOptions struct {
	Regex        string
	Types        []string
	Catalogs     []string
	MQEQueryable bool
}

// Entity is one row of GET /inspect/entities: the decoded entity plus an MQE-ready
// payload to feed back into execExpression.
type Entity struct {
	EntityID  string `json:"entityId"`
	Decoded   any    `json:"decoded"`
	Layer     string `json:"layer"`
	MqeEntity any    `json:"mqeEntity"`
}

// Entities is the response of GET /inspect/entities.
type Entities struct {
	Metric string   `json:"metric"`
	Scope  string   `json:"scope"`
	Step   string   `json:"step"`
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Rows   []Entity `json:"rows"`
}

// EntitiesOptions holds the parameters of GET /inspect/entities.
type EntitiesOptions struct {
	Metric string
	Start  string
	End    string
	Step   string
	Limit  int
	// ValueColumn / ValueType are required only when the metric is NOT defined on the target OAP
	// (a metric persisted by another OAP). When set, the OAP resolves the metric from storage
	// using this caller-supplied metadata instead of its local registry.
	ValueColumn string
	ValueType   string
}

// ListMetrics lists the registered metric catalog (GET /inspect/metrics).
func ListMetrics(ctx context.Context, opts MetricsOptions) (*Metrics, error) {
	query := url.Values{}
	if opts.Regex != "" {
		query.Set("regex", opts.Regex)
	}
	for _, t := range opts.Types {
		query.Add("type", t)
	}
	for _, c := range opts.Catalogs {
		query.Add("catalog", c)
	}
	if opts.MQEQueryable {
		query.Set("mqeQueryable", "true")
	}

	var out Metrics
	err := client.GetJSON(ctx, "/inspect/metrics", query, &out)
	return &out, err
}

// ListEntities enumerates the entities holding values for a metric over a time range
// (GET /inspect/entities). Only REGULAR_VALUE / LABELED_VALUE metrics are accepted.
func ListEntities(ctx context.Context, opts *EntitiesOptions) (*Entities, error) {
	query := url.Values{}
	query.Set("metric", opts.Metric)
	query.Set("start", opts.Start)
	query.Set("end", opts.End)
	query.Set("step", opts.Step)
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.ValueColumn != "" {
		query.Set("valueColumn", opts.ValueColumn)
	}
	if opts.ValueType != "" {
		query.Set("valueType", opts.ValueType)
	}

	var out Entities
	err := client.GetJSON(ctx, "/inspect/entities", query, &out)
	return &out, err
}
