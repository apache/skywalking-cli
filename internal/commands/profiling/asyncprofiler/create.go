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
	"strings"

	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/asyncprofiler"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new async profiler task",
	UsageText: `Create a new async profiler task

Examples:
1. Create async-profiler task
$ swctl profiling async create --service-name=service-name --duration=60 --events=cpu,alloc \ 
	--instance-name-list=instance-name1,instance-name2 --exec-args=interval=50ms`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceListFlags,
		[]cli.Flag{
			&cli.IntFlag{
				Name:     "duration",
				Usage:    "task continuous time(second).",
				Required: true,
			},
			&cli.GenericFlag{
				Name:     "events",
				Usage:    "which event types this task needs to collect.",
				Required: true,
				Value: &asyncprofiler.ProfilerEventTypeEnumValue{
					Enum: query.AllAsyncProfilerEventType,
				},
			},
			&cli.StringFlag{
				Name:  "exec-args",
				Usage: "other async-profiler execution options, e.g. alloc=2k,lock=2s.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseInstanceList(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		instanceIDs := strings.Split(ctx.String("instance-id-list"), ",")
		duration := ctx.Int("duration")
		eventTypes := ctx.Generic("events").(*asyncprofiler.ProfilerEventTypeEnumValue).Selected

		var execArgs *string
		if args := ctx.String("exec-args"); args != "" {
			execArgs = &args
		}

		request := &query.AsyncProfilerTaskCreationRequest{
			ServiceID:          serviceID,
			ServiceInstanceIds: instanceIDs,
			Duration:           duration,
			Events:             eventTypes,
			ExecArgs:           execArgs,
		}
		task, err := profiling.CreateAsyncProfilerTask(ctx.Context, request)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: task, Condition: request})
	},
}
