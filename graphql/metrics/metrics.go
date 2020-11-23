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

package metrics

import (
	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"

	"github.com/machinebox/graphql"
	"github.com/urfave/cli"
)

func IntValues(ctx *cli.Context, condition schema.MetricsCondition, duration schema.Duration) int {
	var response map[string]int

	request := graphql.NewRequest(assets.Read("graphqls/metrics/MetricsValue.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func LinearIntValues(ctx *cli.Context, condition schema.MetricsCondition, duration schema.Duration) schema.MetricsValues {
	var response map[string]schema.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/MetricsValues.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func MultipleLinearIntValues(ctx *cli.Context, condition schema.MetricsCondition, labels []string, duration schema.Duration) []schema.MetricsValues {
	var response map[string][]schema.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/LabeledMetricsValues.graphql"))

	request.Var("duration", duration)
	request.Var("condition", condition)
	request.Var("labels", labels)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func Thermodynamic(ctx *cli.Context, condition schema.MetricsCondition, duration schema.Duration) schema.HeatMap {
	var response map[string]schema.HeatMap

	request := graphql.NewRequest(assets.Read("graphqls/metrics/HeatMap.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func SortMetrics(ctx *cli.Context, condition schema.TopNCondition, duration schema.Duration) []*schema.SelectedRecord {
	var response map[string][]*schema.SelectedRecord

	request := graphql.NewRequest(assets.Read("graphqls/metrics/SortMetrics.graphql"))
	request.Var("condition", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func ListMetrics(ctx *cli.Context, regex string) []*schema.MetricDefinition {
	var response map[string][]*schema.MetricDefinition
	request := graphql.NewRequest(assets.Read("graphqls/metrics/ListMetrics.graphql"))
	request.Var("regex", regex)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}
