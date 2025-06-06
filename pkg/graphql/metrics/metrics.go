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
	"context"

	"github.com/machinebox/graphql"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	api "skywalking.apache.org/repo/goapi/query"
)

func IntValues(ctx context.Context, condition api.MetricsCondition, duration api.Duration) (int, error) {
	var response map[string]int

	request := graphql.NewRequest(assets.Read("graphqls/metrics/MetricsValue.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func NullableIntValue(ctx context.Context, condition api.MetricsCondition, duration api.Duration) (api.NullableValue, error) {
	var response map[string]api.NullableValue

	request := graphql.NewRequest(assets.Read("graphqls/metrics/NullableMetricsValue.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func LinearIntValues(ctx context.Context, condition api.MetricsCondition, duration api.Duration) (api.MetricsValues, error) {
	var response map[string]api.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/MetricsValues.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func MultipleLinearIntValues(
	ctx context.Context,
	condition api.MetricsCondition,
	labels []string,
	duration api.Duration,
) ([]api.MetricsValues, error) {
	var response map[string][]api.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/metrics/LabeledMetricsValues.graphql"))

	request.Var("duration", duration)
	request.Var("condition", condition)
	request.Var("labels", labels)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func Thermodynamic(ctx context.Context, condition api.MetricsCondition, duration api.Duration) (api.HeatMap, error) {
	var response map[string]api.HeatMap

	request := graphql.NewRequest(assets.Read("graphqls/metrics/HeatMap.graphql"))

	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func SortMetrics(ctx context.Context, condition api.TopNCondition, duration api.Duration) ([]*api.SelectedRecord, error) {
	var response map[string][]*api.SelectedRecord

	request := graphql.NewRequest(assets.Read("graphqls/metrics/SortMetrics.graphql"))
	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func SampledRecords(ctx context.Context, condition api.TopNCondition, duration api.Duration) ([]*api.SelectedRecord, error) {
	var response map[string][]*api.SelectedRecord

	request := graphql.NewRequest(assets.Read("graphqls/metrics/SampledRecords.graphql"))
	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ReadRecords(ctx context.Context, condition api.RecordCondition, duration api.Duration) ([]*api.Record, error) {
	var response map[string][]*api.Record

	request := graphql.NewRequest(assets.Read("graphqls/metrics/ReadRecords.graphql"))
	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ListMetrics(ctx context.Context, regex string) ([]*api.MetricDefinition, error) {
	var response map[string][]*api.MetricDefinition
	request := graphql.NewRequest(assets.Read("graphqls/metrics/ListMetrics.graphql"))
	request.Var("regex", regex)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func Execute(ctx context.Context, expression string, entity *api.Entity, duration api.Duration) (api.ExpressionResult, error) {
	var response map[string]api.ExpressionResult

	request := graphql.NewRequest(assets.Read("graphqls/metrics/ExecuteExpression.graphql"))
	request.Var("expression", expression)
	request.Var("entity", *entity)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
