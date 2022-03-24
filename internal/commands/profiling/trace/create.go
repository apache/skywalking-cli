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

package trace

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new trace profiling task",
	UsageText: `Create a new trace profiling task

Examples:
1. Create trace profiling task
$ swctl profiling trace create --service-name=service-name --endpoint=endpoint --start-time=1627656127860 --duration=5 \
	--min-duration-threshold=0 --dump-period=10 --max-sampling-count=9`,
	Flags: flags.Flags(
		flags.EndpointFlags,

		[]cli.Flag{
			&cli.Int64Flag{
				Name:  "start-time",
				Usage: "profile task start time(millisecond).",
			},
			&cli.IntFlag{
				Name:  "duration",
				Usage: "profile task continuous time(minute).",
			},
			&cli.IntFlag{
				Name:  "min-duration-threshold",
				Usage: "profiled endpoint must greater duration(millisecond).",
			},
			&cli.IntFlag{
				Name:  "dump-period",
				Usage: "profiled endpoint dump period(millisecond).",
			},
			&cli.IntFlag{
				Name:  "max-sampling-count",
				Usage: "profile task max sampling count.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseEndpoint(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		startTime := ctx.Int64("start-time")
		request := &api.ProfileTaskCreationRequest{
			ServiceID:            serviceID,
			EndpointName:         ctx.String("endpoint-name"),
			StartTime:            &startTime,
			Duration:             ctx.Int("duration"),
			MinDurationThreshold: ctx.Int("min-duration-threshold"),
			DumpPeriod:           ctx.Int("dump-period"),
			MaxSamplingCount:     ctx.Int("max-sampling-count"),
		}

		task, err := profiling.CreateTraceTask(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
