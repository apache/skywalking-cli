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

package create

import (
	"strings"
	"time"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/ebpf"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

var FixedTimeCreateCommand = &cli.Command{
	Name:    "fixed",
	Aliases: []string{"ft"},
	Usage:   "Create a new ebpf profiling fixed time task",
	UsageText: `Create a new ebpf profiling fixed time task

Examples:
1. Create ebpf profiling fixed time task
$ swctl profiling ebpf fixed --service-id=abc --process-id=abc --duration=1m --target-type=ON_CPU`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:  "labels",
				Usage: "the `labels` by which labels of the process need to be profiling, multiple labels split by ',': l1,l2",
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
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		labelStr := ctx.String("labels")
		labels := strings.Split(labelStr, ",")

		duration, err := time.ParseDuration(ctx.String("duration"))
		if err != nil {
			return err
		}
		request := &api.EBPFProfilingTaskFixedTimeCreationRequest{
			ServiceID:     serviceID,
			ProcessLabels: labels,
			StartTime:     ctx.Int64("start-time"),
			Duration:      int(duration.Seconds()),
			TargetType:    ctx.Generic("target-type").(*ebpf.ProfilingTargetTypeEnumValue).Selected,
		}

		task, err := profiling.CreateEBPFProfilingFixedTimeTask(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
