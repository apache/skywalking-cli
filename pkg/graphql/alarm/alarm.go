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

package alarm

import (
	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

type ListAlarmCondition struct {
	Duration *api.Duration
	Keyword  string
	Scope    api.Scope
	Tags     []*api.AlarmTag
	Paging   *api.Pagination
}

func Alarms(ctx *cli.Context, condition *ListAlarmCondition) (api.Alarms, error) {
	var response map[string]api.Alarms

	request := graphql.NewRequest(assets.Read("graphqls/alarm/alarms.graphql"))
	request.Var("paging", condition.Paging)
	request.Var("tags", condition.Tags)
	request.Var("duration", condition.Duration)
	request.Var("keyword", condition.Keyword)
	if condition.Scope != "" {
		request.Var("scope", condition.Scope)
	}

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
