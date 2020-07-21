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

package dashboard

import (
	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
)

type GlobalMetrics struct {
	ServiceLoad       []*schema.SelectedRecord `json:"serviceLoad"`
	SlowServices      []*schema.SelectedRecord `json:"slowServices"`
	UnhealthyServices []*schema.SelectedRecord `json:"unhealthyServices"`
	SlowEndpoints     []*schema.SelectedRecord `json:"slowEndpoints"`
}

type GlobalData struct {
	Metrics         *GlobalMetrics
	ResponseLatency []*schema.MetricsValues `json:"responseLatency"`
	HeatMap         schema.HeatMap          `json:"heatMap"`
}

func Metrics(ctx *cli.Context, duration schema.Duration) *GlobalMetrics {
	var response map[string][]*schema.SelectedRecord

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/GlobalMetrics.graphql"))
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return &GlobalMetrics{
		ServiceLoad:       response["serviceLoad"],
		SlowServices:      response["slowServices"],
		UnhealthyServices: response["unhealthyServices"],
		SlowEndpoints:     response["slowEndpoints"],
	}
}

func responseLatency(ctx *cli.Context, duration schema.Duration) []*schema.MetricsValues {
	var response map[string][]*schema.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/ResponseLatency.graphql"))
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func heatMap(ctx *cli.Context, duration schema.Duration) schema.HeatMap {
	var response map[string]schema.HeatMap

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/HeatMap.graphql"))
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func Global(ctx *cli.Context, duration schema.Duration) GlobalData {
	var globalData GlobalData

	globalData.Metrics = Metrics(ctx, duration)
	globalData.ResponseLatency = responseLatency(ctx, duration)
	globalData.HeatMap = heatMap(ctx, duration)

	return globalData
}
