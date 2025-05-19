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

package asyncprofiler

import (
	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var getTaskListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Query async-profiler task list",
	UsageText: `Query async-profiler task list

Examples:
1. Query all async-profiler tasks
$ swctl profiling async list --service-name=service-name`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.DurationFlags,
		[]cli.Flag{
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit defines the number of the tasks to be returned.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step")
		cold := ctx.Bool("cold")
		duration := query.Duration{
			Start:     start,
			End:       end,
			Step:      step.(*model.StepEnumValue).Selected,
			ColdStage: &cold,
		}
		var limit *int
		if limitArg := ctx.Int("limit"); limitArg != 0 {
			limit = &limitArg
		}

		request := &query.AsyncProfilerTaskListRequest{
			ServiceID:     serviceID,
			QueryDuration: &duration,
			Limit:         limit,
		}

		tasks, err := profiling.GetAsyncProfilerTaskList(ctx.Context, request)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: tasks, Condition: request})
	},
}
