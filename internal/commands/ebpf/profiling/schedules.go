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
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	ebpf_graphql "github.com/apache/skywalking-cli/pkg/graphql/ebpf"
)

var ListScheduleCommand = &cli.Command{
	Name:    "schedules",
	Aliases: []string{"ss"},
	Usage:   `query ebpf profiling task schedules`,
	UsageText: `This command lists all schedule of the ebpf profiling task, via id of task.

Example：
1. Query profiling schedules of task id "abc"
$ swctl ebpf-profiling schedules --task-id=abc
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "the `task-id` by which task are scheduled",
				Required: true,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step")

		schedules, err := ebpf_graphql.QueryEBPFProfilingScheduleList(ctx, taskID, &api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: schedules})
	},
}
