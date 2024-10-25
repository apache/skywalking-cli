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
	"strings"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new async profiler task",
	UsageText: `Create a new async profiler task

Examples:
1. Create async-profiler task
$ swctl profiling asyncprofiler create --service-name=someservicename --duration=60 --events=cpu --service-instance-ids=someinstance`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringSliceFlag{
				Name:     "service-instance-ids",
				Usage:    "which instances to execute task.",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "duration",
				Usage:    "task continuous time(second).",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "events",
				Usage:    "which event types this task needs to collect.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "exec-args",
				Usage: "other async-profiler execution options, e.g. alloc=2k,lock=2s.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(false),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		events := ctx.StringSlice("events")
		eventTypes := make([]query.AsyncProfilerEventType, 0)
		for _, event := range events {
			upperCaseEvent := strings.ToUpper(event)
			eventTypes = append(eventTypes, query.AsyncProfilerEventType(upperCaseEvent))
		}

		execArgs := ctx.String("exec-args")

		request := &query.AsyncProfilerTaskCreationRequest{
			ServiceID:          serviceID,
			ServiceInstanceIds: ctx.StringSlice("service-instance-ids"),
			Duration:           ctx.Int("duration"),
			Events:             eventTypes,
			ExecArgs:           &execArgs,
		}
		task, err := profiling.CreateAsyncProfilerTask(ctx, request)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
