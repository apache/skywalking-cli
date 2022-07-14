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

func CreateEBPFProfilingFixedTimeTask(ctx *cli.Context,
	condition *api.EBPFProfilingTaskFixedTimeCreationRequest) (api.EBPFProfilingTaskCreationResult, error) {
	var response map[string]api.EBPFProfilingTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/CreateEBPFProfilingFixedTimeTask.graphql"))
	request.Var("request", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func CreateEBPFNetworkProfilingTask(ctx *cli.Context, condition *api.EBPFProfilingNetworkTaskRequest) (api.EBPFProfilingTaskCreationResult, error) {
	var response map[string]api.EBPFProfilingTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/CreateEBPFNetworkProfilingTask.graphql"))
	request.Var("request", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryPrepareCreateEBPFProfilingTaskData(ctx *cli.Context, serviceID string) (*api.EBPFProfilingTaskPrepare, error) {
	var response map[string]*api.EBPFProfilingTaskPrepare

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/QueryPrepareCreateEBPFProfilingTaskData.graphql"))
	request.Var("serviceId", serviceID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryEBPFProfilingTaskList(ctx *cli.Context, serviceID string) ([]*api.EBPFProfilingTask, error) {
	var response map[string][]*api.EBPFProfilingTask

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/QueryEBPFProfilingTaskList.graphql"))
	request.Var("serviceId", serviceID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryEBPFProfilingScheduleList(ctx *cli.Context, taskID string) ([]*api.EBPFProfilingSchedule, error) {
	var response map[string][]*api.EBPFProfilingSchedule

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/QueryEBPFProfilingScheduleList.graphql"))
	request.Var("taskID", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func AnalysisEBPFProfilingResult(ctx *cli.Context, scheduleIDList []string,
	timeRanges []*api.EBPFProfilingAnalyzeTimeRange, aggregateType api.EBPFProfilingAnalyzeAggregateType) (*api.EBPFProfilingAnalyzation, error) {
	var response map[string]*api.EBPFProfilingAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/AnalysisEBPFProfilingResult.graphql"))
	request.Var("scheduleIdList", scheduleIDList)
	request.Var("timeRanges", timeRanges)
	request.Var("aggregateType", aggregateType)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func KeepNetworkProfilingTask(ctx *cli.Context, taskID string) (*api.EBPFNetworkKeepProfilingResult, error) {
	var response map[string]*api.EBPFNetworkKeepProfilingResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/ebpf/KeepNetworkProfilingTask.graphql"))
	request.Var("taskId", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
