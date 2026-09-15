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

// Package aiagent wraps the GraphQL query of ai-agent-conversation.graphqls: the list
// page of the conversations the AI Sessionizer lands. The conversation document and its
// files are not GraphQL queries; see pkg/aiagent/view and pkg/aiagent/files.
package aiagent

import (
	"context"

	"github.com/machinebox/graphql"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
)

// ListConversations lists one row per conversation of a service active in the duration,
// newest first, from the newest round's attributes.
func ListConversations(ctx context.Context, condition *api.ConversationListCondition, duration api.Duration) (api.ConversationList, error) {
	var response map[string]api.ConversationList

	request := graphql.NewRequest(assets.Read("graphqls/aiagent/ListConversations.graphql"))
	request.Var("condition", condition)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)
	return response["result"], err
}
