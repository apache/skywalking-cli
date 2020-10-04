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
	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/assets"

	"github.com/apache/skywalking-cli/graphql/client"

	"github.com/apache/skywalking-cli/graphql/schema"
)

func IntValues(ctx *cli.Context, condition schema.BatchMetricConditions, duration schema.Duration) schema.IntValues {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/IntValues.graphql"))

	request.Var("metric", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func LinearIntValues(ctx *cli.Context, condition schema.MetricCondition, duration schema.Duration) schema.IntValues {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/LinearIntValues.graphql"))

	request.Var("metric", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func MultipleLinearIntValues(ctx *cli.Context, condition schema.MetricCondition, numOfLinear int, duration schema.Duration) []schema.IntValues {
	var response map[string][]schema.IntValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/MultipleLinearIntValues.graphql"))

	request.Var("metric", condition)
	request.Var("numOfLinear", numOfLinear)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func Thermodynamic(ctx *cli.Context, condition schema.MetricsCondition, duration schema.Duration) schema.HeatMap {
	var response map[string]schema.HeatMap

	request := graphql.NewRequest(assets.Read("graphqls/metrics/Thermodynamic.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}
