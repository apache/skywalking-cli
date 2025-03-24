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

	"github.com/machinebox/graphql"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
)

func CreateAsyncProfilerTask(ctx context.Context, condition *api.AsyncProfilerTaskCreationRequest) (api.AsyncProfilerTaskCreationResult, error) {
	var response map[string]api.AsyncProfilerTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/CreateTask.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerTaskList(ctx context.Context, condition *api.AsyncProfilerTaskListRequest) (api.AsyncProfilerTaskListResult, error) {
	var response map[string]api.AsyncProfilerTaskListResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetTaskList.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerTaskProgress(ctx context.Context, taskID string) (api.AsyncProfilerTaskProgress, error) {
	var response map[string]api.AsyncProfilerTaskProgress

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetTaskProgress.graphql"))
	request.Var("taskId", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerAnalyze(ctx context.Context, condition *api.AsyncProfilerAnalyzationRequest) (api.AsyncProfilerAnalyzation, error) {
	var response map[string]api.AsyncProfilerAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetAnalysis.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
