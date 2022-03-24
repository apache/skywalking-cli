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
	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

func CreateTraceTask(ctx *cli.Context, condition *api.ProfileTaskCreationRequest) (api.ProfileTaskCreationResult, error) {
	var response map[string]api.ProfileTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/CreateTask.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetTraceProfilingTaskList(ctx *cli.Context, serviceID, endpointName string) ([]*api.ProfileTask, error) {
	var response map[string][]*api.ProfileTask

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/GetTaskList.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("endpointName", endpointName)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetTraceProfilingTaskLogList(ctx *cli.Context, taskID string) ([]*api.ProfileTaskLog, error) {
	var response map[string][]*api.ProfileTaskLog

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/GetProfileTaskLogs.graphql"))
	request.Var("taskID", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetTraceProfilingTaskSegmentList(ctx *cli.Context, taskID string) ([]*api.BasicTrace, error) {
	var response map[string][]*api.BasicTrace

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/GetTaskSegmentList.graphql"))
	request.Var("taskId", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetTraceProfilingSegment(ctx *cli.Context, segmentID string) (api.ProfiledSegment, error) {
	var response map[string]api.ProfiledSegment

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/GetProfiledSegment.graphql"))
	request.Var("segmentId", segmentID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetTraceProfilingAnalyze(ctx *cli.Context, segmentID string, timeRanges []*api.ProfileAnalyzeTimeRange) (api.ProfileAnalyzation, error) {
	var response map[string]api.ProfileAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/profiling/trace/GetProfileAnalyze.graphql"))
	request.Var("segmentId", segmentID)
	request.Var("timeRanges", timeRanges)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
