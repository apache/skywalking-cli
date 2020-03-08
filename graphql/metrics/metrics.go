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

	"github.com/apache/skywalking-cli/graphql/client"

	"github.com/apache/skywalking-cli/graphql/schema"
)

func IntValues(ctx *cli.Context, condition schema.BatchMetricConditions, duration schema.Duration) schema.IntValues {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(`
		query ($metric: BatchMetricConditions!, $duration: Duration!) {
			metrics: getValues(metric: $metric, duration: $duration) {
				values { id value }
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["metrics"]
}

func LinearIntValues(ctx *cli.Context, condition schema.MetricCondition, duration schema.Duration) schema.IntValues {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(`
		query ($metric: MetricCondition!, $duration: Duration!) {
			metrics: getLinearIntValues(metric: $metric, duration: $duration) {
				values { value }
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["metrics"]
}

func MultipleLinearIntValues(ctx *cli.Context, condition schema.MetricCondition, numOfLinear int, duration schema.Duration) []schema.IntValues {
	request := graphql.NewRequest(`
		query ($metric: MetricCondition!, $numOfLinear: Int!, $duration: Duration!) {
			metrics: getMultipleLinearIntValues(metric: $metric, numOfLinear: $numOfLinear, duration: $duration) {
				values { value }
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("numOfLinear", numOfLinear)
	request.Var("duration", duration)

	var response map[string][]schema.IntValues

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["metrics"]
}

func Thermodynamic(ctx *cli.Context, condition schema.MetricCondition, duration schema.Duration) schema.Thermodynamic {
	request := graphql.NewRequest(`
		query ($metric: MetricCondition!, $duration: Duration!) {
			metrics: getThermodynamic(metric: $metric, duration: $duration) {
				nodes axisYStep
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("duration", duration)

	var response map[string]schema.Thermodynamic

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["metrics"]
}
