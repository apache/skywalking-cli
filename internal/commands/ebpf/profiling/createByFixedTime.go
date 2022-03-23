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

package profiling

import (
	"time"

	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/ebpf"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	ebpf_graphql "github.com/apache/skywalking-cli/pkg/graphql/ebpf"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

var CreateCommand = &cli.Command{
	Name:    "createByFixedTime",
	Aliases: []string{"cft"},
	Usage:   "Create a new ebpf profiling fixed time task",
	UsageText: `Create a new ebpf profiling fixed time task

Examples:
1. Create ebpf profiling fixed time task
$ swctl ebpf-profiling createByFixedTime --process-finder=PROCESS_ID --process-id=abc --duration=1m target-type=ON_CPU`,
	Flags: flags.Flags(
		[]cli.Flag{
			&cli.GenericFlag{
				Name:  "process-finder",
				Usage: "the `process-finder` by the way to address the target process",
				Value: &ebpf.ProfilingProcessFinderTypeEnumValue{
					Enum:     api.AllEBPFProfilingProcessFinderType,
					Default:  api.EBPFProfilingProcessFinderTypeProcessID,
					Selected: api.EBPFProfilingProcessFinderTypeProcessID,
				},
			},
			&cli.StringFlag{
				Name:     "process-id",
				Usage:    "the `process-id` by which process ID need to be profiling",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "duration",
				Usage:    "profiling task continuous time.",
				Required: true,
			},
			&cli.Int64Flag{
				Name:  "start-time",
				Usage: "profiling task start time(millisecond).",
			},
			&cli.GenericFlag{
				Name:  "target-type",
				Usage: "the `target-type` by the way of profiling the process",
				Value: &ebpf.ProfilingTargetTypeEnumValue{
					Enum:     api.AllEBPFProfilingTargetType,
					Default:  api.EBPFProfilingTargetTypeOnCPU,
					Selected: api.EBPFProfilingTargetTypeOnCPU,
				},
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		processID := ctx.String("process-id")
		duration, err := time.ParseDuration(ctx.String("duration"))
		if err != nil {
			return err
		}
		request := &api.EBPFProfilingTaskFixedTimeCreationRequest{
			ProcessFinder: &api.EBPFProfilingProcessFinder{
				FinderType: ctx.Generic("process-finder").(*ebpf.ProfilingProcessFinderTypeEnumValue).Selected,
				ProcessID:  &processID,
			},
			StartTime:  ctx.Int64("start-time"),
			Duration:   int(duration.Seconds()),
			TargetType: ctx.Generic("target-type").(*ebpf.ProfilingTargetTypeEnumValue).Selected,
		}

		task, err := ebpf_graphql.CreateEBPFProfilingFixedTimeTask(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
