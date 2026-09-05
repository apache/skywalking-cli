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

// Package aiagent holds the commands for the conversations of long-lived AI agents
// that the AI Sessionizer (apache/skywalking-ai-sessionizer) lands in the OAP under
// the AI_AGENT layer: the list page, the raw-file export, and the conversation itself
// as one asz.view document.
package aiagent

import (
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:  "ai-agent",
	Usage: "AI agent conversations landed by the AI Sessionizer",
	UsageText: `The AI Sessionizer collects an agent runtime's transcripts and pushes them to the OAP
under the AI_AGENT layer. "list" and "files" are GraphQL queries on the "--base-url"
endpoint; "view" reads the whole conversation as one asz.view document from the OAP's
streamed route on the same host, GET /ai-agent/conversations/{conversation}/v1/view.`,
	Subcommands: []*cli.Command{
		listCommand,
		filesCommand,
		viewCommand,
	},
}
