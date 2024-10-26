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
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"
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
		[]cli.Flag{
			&cli.Int64Flag{
				Name:  "start-time",
				Usage: "The start time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
			},
			&cli.Int64Flag{
				Name:  "end-time",
				Usage: "The end time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit defines the number of the tasks to be returned.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		var startTime *int64
		if startTimeArg := ctx.Int64("start-time"); startTimeArg != 0 {
			startTime = &startTimeArg
		}
		var endTime *int64
		if endTimeArg := ctx.Int64("end-time"); endTimeArg != 0 {
			endTime = &endTimeArg
		}
		var limit *int
		if limitArg := ctx.Int("limit"); limitArg != 0 {
			limit = &limitArg
		}

		request := &query.AsyncProfilerTaskListRequest{
			ServiceID: serviceID,
			StartTime: startTime,
			EndTime:   endTime,
			Limit:     limit,
		}

		tasks, err := profiling.GetAsyncProfilerTaskList(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: tasks, Condition: request})
	},
}
