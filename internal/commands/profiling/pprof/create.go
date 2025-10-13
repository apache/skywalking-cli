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

package pprof

import (
	"strings"

	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/pprof"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new pprof task",
	UsageText: `Create a new pprof task

Examples:
1. Create pprof task
$ swctl profiling pprof create --service-name=service-name --duration=10 --events=mutex \ 
	--instance-name-list=instance-name1,instance-name2 --dump-period=1`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceListFlags,
		[]cli.Flag{
			&cli.IntFlag{
				Name:  "duration",
				Usage: "task continuous time(minute), required for cpu, block and mutex events. ",
			},
			&cli.GenericFlag{
				Name:     "events",
				Usage:    "which event types this task needs to collect.",
				Required: true,
				Value: &pprof.ProfilingEventTypeEnumValue{
					Enum:    query.AllPprofEventType,
					Default: query.PprofEventTypeCPU,
				},
			},
			&cli.StringFlag{
				Name:  "dump-period",
				Usage: "pprof dump period parameters, required for block and mutex events. ",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
		interceptor.ParseInstanceList(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		instanceIDs := strings.Split(ctx.String("instance-id-list"), ",")
		eventTypes := ctx.Generic("events").(*pprof.ProfilingEventTypeEnumValue).Selected
		var dumpPeriod, duration *int
		if durationtime := ctx.Int("duration"); durationtime != 0 {
			duration = &durationtime
		}
		if period := ctx.Int("dump-period"); period != 0 {
			dumpPeriod = &period
		}
		request := &query.PprofTaskCreationRequest{
			ServiceID:          serviceID,
			ServiceInstanceIds: instanceIDs,
			Duration:           duration,
			Events:             eventTypes,
			DumpPeriod:         dumpPeriod,
		}
		task, err := profiling.CreatePprofTask(ctx.Context, request)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: task, Condition: request})
	},
}
