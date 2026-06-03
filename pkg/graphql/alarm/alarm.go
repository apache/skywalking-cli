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
	"context"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"

	api "skywalking.apache.org/repo/goapi/query"
)

type ListAlarmCondition struct {
	Duration  *api.Duration
	Keyword   string
	Tags      []*api.AlarmTag
	Paging    *api.Pagination
	Layer     string
	RuleNames []string
	Entities  []*api.Entity
}

// alarmQueryCondition mirrors the GraphQL `AlarmQueryCondition` input of `queryAlarms`.
// It is JSON-encoded as the `condition` variable, so the field tags must match the
// schema field names exactly.
type alarmQueryCondition struct {
	Duration  *api.Duration   `json:"duration"`
	Paging    *api.Pagination `json:"paging"`
	Entities  []*api.Entity   `json:"entities,omitempty"`
	Layer     *string         `json:"layer,omitempty"`
	RuleNames []string        `json:"ruleNames,omitempty"`
	Keyword   *string         `json:"keyword,omitempty"`
	Tags      []*api.AlarmTag `json:"tags,omitempty"`
}

func Alarms(ctx context.Context, condition *ListAlarmCondition) (api.Alarms, error) {
	var response map[string]api.Alarms

	queryCondition := alarmQueryCondition{
		Duration:  condition.Duration,
		Paging:    condition.Paging,
		Entities:  condition.Entities,
		RuleNames: condition.RuleNames,
		Tags:      condition.Tags,
	}
	if condition.Keyword != "" {
		queryCondition.Keyword = &condition.Keyword
	}
	if condition.Layer != "" {
		queryCondition.Layer = &condition.Layer
	}

	request := graphql.NewRequest(assets.Read("graphqls/alarm/alarms.graphql"))
	request.Var("condition", queryCondition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func TagAutocompleteKeys(ctx context.Context, duration api.Duration) ([]string, error) {
	var response map[string][]string

	request := graphql.NewRequest(assets.Read("graphqls/alarm/tagAutocompleteKeys.graphql"))
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func TagAutocompleteValues(ctx context.Context, duration api.Duration, key string) ([]string, error) {
	var response map[string][]string

	request := graphql.NewRequest(assets.Read("graphqls/alarm/tagAutocompleteValues.graphql"))
	request.Var("duration", duration)
	request.Var("tagKey", key)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
