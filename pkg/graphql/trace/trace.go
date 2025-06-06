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

package trace

import (
	"context"

	"github.com/machinebox/graphql"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
)

func Trace(ctx context.Context, traceID string) (api.Trace, error) {
	var response map[string]api.Trace

	request := graphql.NewRequest(assets.Read("graphqls/trace/Trace.graphql"))
	request.Var("traceId", traceID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ColdTrace(ctx context.Context, duration api.Duration, traceID string) (api.Trace, error) {
	var response map[string]api.Trace

	request := graphql.NewRequest(assets.Read("graphqls/trace/ColdTrace.graphql"))
	request.Var("traceId", traceID)
	request.Var("duration", duration)
	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func Traces(ctx context.Context, condition *api.TraceQueryCondition) (api.TraceBrief, error) {
	var response map[string]api.TraceBrief

	request := graphql.NewRequest(assets.Read("graphqls/trace/Traces.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
