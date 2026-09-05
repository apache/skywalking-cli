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

package aiagent

import (
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/aiagent"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List the conversations of an AI agent service",
	UsageText: `List the conversations of an AI agent service active in the duration, newest first,
one row per conversation from its newest round.

Examples:
1. The conversations of service "Claude Code" in the last 30 minutes:
$ swctl ai-agent list --service-name "Claude Code"

2. Only those pushed by one Sessionizer, in a day:
$ swctl ai-agent list --service-name "Claude Code" --instance-name laptop --start 2026-09-01 --end 2026-09-02`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.ServiceFlags,
		flags.InstanceFlags,
		[]cli.Flag{
			&cli.IntFlag{
				Name:  "limit",
				Usage: "at most this many rounds are read, newest first, before folding to one row per conversation; 0 for the OAP's default",
				Value: 0,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseService(true),
		interceptor.ParseInstance(false),
	),
	Action: func(ctx *cli.Context) error {
		duration := api.Duration{
			Start: ctx.String("start"),
			End:   ctx.String("end"),
			Step:  ctx.Generic("step").(*model.StepEnumValue).Selected,
		}
		condition := &api.ConversationListCondition{
			Service:  &api.ServiceCondition{ServiceName: ctx.String("service-name")},
			Instance: instanceCondition(ctx),
		}
		if limit := ctx.Int("limit"); limit > 0 {
			condition.Limit = &limit
		}

		list, err := aiagent.ListConversations(ctx.Context, condition, duration)
		if err != nil {
			return err
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: list, Condition: condition, Duration: duration})
	},
}

// instanceCondition names the sender when "--instance-name" (or "--instance-id",
// resolved to the name by the interceptor) was given.
func instanceCondition(ctx *cli.Context) *api.InstanceCondition {
	name := ctx.String("instance-name")
	if name == "" {
		return nil
	}
	return &api.InstanceCondition{ServiceName: ctx.String("service-name"), InstanceName: name}
}
