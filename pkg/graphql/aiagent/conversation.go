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

// Package aiagent wraps the GraphQL queries of ai-agent-conversation.graphqls: the
// list page and the raw-file export of the conversations the AI Sessionizer lands.
// The conversation document itself is not a GraphQL query; see pkg/aiagent/view.
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

// RawFiles lists every landed file and round of a conversation as stored, or only the
// named ones; with body, each file comes verbatim, which is the export path.
func RawFiles(ctx context.Context, condition *api.ConversationCondition, files []string, body bool) (api.ConversationRawFiles, error) {
	var response map[string]api.ConversationRawFiles

	request := graphql.NewRequest(assets.Read("graphqls/aiagent/ConversationRawFiles.graphql"))
	request.Var("condition", condition)
	request.Var("files", files)
	request.Var("body", body)

	err := client.ExecuteQuery(ctx, request, &response)
	return response["result"], err
}
