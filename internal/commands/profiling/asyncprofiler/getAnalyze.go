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

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/asyncprofiler"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	query "skywalking.apache.org/repo/goapi/query"
)

var analysisCommand = &cli.Command{
	Name:    "analysis",
	Aliases: []string{"a"},
	Usage:   "Query async-profiler analysis",
	UsageText: `Query async-profiler analysis

Examples:
1. Query the flame graph produced by async-profiler
$ swctl profiling async analysis --service-name=service-name --task-id=task-id --service-instance-ids=instance-name1,instance-name2 --event=execution_sample`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceSliceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "async-profiler task id",
				Required: true,
			},
			&cli.GenericFlag{
				Name:     "event",
				Usage:    "which event types this task needs to collect.",
				Required: true,
				Value: &asyncprofiler.JFREventTypeEnumValue{
					Enum: query.AllJFREventType,
				},
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseInstanceSlice(true),
	),
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")
		instances := strings.Split(ctx.String("instance-id-slice"), ",")
		eventType := ctx.Generic("event").(*asyncprofiler.JFREventTypeEnumValue).Selected

		request := &query.AsyncProfilerAnalyzationRequest{
			TaskID:      taskID,
			InstanceIds: instances,
			EventType:   eventType,
		}

		analyze, err := profiling.GetAsyncProfilerAnalyze(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: analyze, Condition: request})
	},
}
