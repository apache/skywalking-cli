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

package ebpf

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

	request := graphql.NewRequest(assets.Read("graphqls/ebpf/CreateEBPFProfilingFixedTimeTask.graphql"))
	request.Var("request", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryEBPFProfilingTaskList(ctx *cli.Context,
	condition *api.EBPFProfilingTaskCondition) ([]*api.EBPFProfilingTask, error) {
	var response map[string][]*api.EBPFProfilingTask

	request := graphql.NewRequest(assets.Read("graphqls/ebpf/QueryEBPFProfilingTaskList.graphql"))
	request.Var("query", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryEBPFProfilingScheduleList(ctx *cli.Context, taskID string,
	duration *api.Duration) ([]*api.EBPFProfilingSchedule, error) {
	var response map[string][]*api.EBPFProfilingSchedule

	request := graphql.NewRequest(assets.Read("graphqls/ebpf/QueryEBPFProfilingScheduleList.graphql"))
	request.Var("taskID", taskID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func QueryEBPFProfilingAnalyzation(ctx *cli.Context, taskID string,
	timeRanges []*api.EBPFProfilingAnalyzeTimeRange) (*api.EBPFProfilingAnalyzation, error) {
	var response map[string]*api.EBPFProfilingAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/ebpf/QueryEBPFProflingAnalyzation.graphql"))
	request.Var("taskID", taskID)
	request.Var("timeRanges", timeRanges)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
