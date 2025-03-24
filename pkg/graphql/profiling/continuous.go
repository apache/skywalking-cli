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

package profiling

import (
	"context"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"

	api "skywalking.apache.org/repo/goapi/query"
)

func SetContinuousProfilingPolicy(ctx context.Context, creation *api.ContinuousProfilingPolicyCreation) (api.ContinuousProfilingSetResult, error) {
	var response map[string]api.ContinuousProfilingSetResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/continuous/SetContinuousProfilingPolicy.graphql"))
	request.Var("request", creation)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryContinuousProfilingServiceTargets(ctx context.Context, serviceID string) ([]*api.ContinuousProfilingPolicyTarget, error) {
	var response map[string][]*api.ContinuousProfilingPolicyTarget

	request := graphql.NewRequest(assets.Read("graphqls/profiling/continuous/QueryContinuousProfilingServiceTargets.graphql"))
	request.Var("serviceId", serviceID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryContinuousProfilingMonitoringInstances(ctx context.Context, serviceID string,
	target api.ContinuousProfilingTargetType,
) ([]api.ContinuousProfilingMonitoringInstance, error) {
	var response map[string][]api.ContinuousProfilingMonitoringInstance

	request := graphql.NewRequest(assets.Read("graphqls/profiling/continuous/QueryContinuousProfilingMonitoringInstances.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("target", target)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
