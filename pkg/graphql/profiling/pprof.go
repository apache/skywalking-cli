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

func CreatePprofTask(ctx context.Context, condition *api.PprofTaskCreationRequest) (api.PprofTaskCreationResult, error) {
	var response map[string]api.PprofTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/pprof/CreateTask.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetPprofTaskProgress(ctx context.Context, taskID string) (api.PprofTaskProgress, error) {
	var response map[string]api.PprofTaskProgress

	request := graphql.NewRequest(assets.Read("graphqls/profiling/pprof/GetTaskProgress.graphql"))
	request.Var("taskId", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetPprofTaskList(ctx context.Context, condition *api.PprofTaskListRequest) (api.PprofTaskListResult, error) {
	var response map[string]api.PprofTaskListResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/pprof/GetTaskList.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetPprofAnalyze(ctx context.Context, condition *api.PprofAnalyzationRequest) (api.PprofAnalyzation, error) {
	var response map[string]api.PprofAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/profiling/pprof/GetAnalysis.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
