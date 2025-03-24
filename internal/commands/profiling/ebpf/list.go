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

package ebpf

import (
	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/ebpf"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var ListTaskCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `query ebpf profiling task list`,
	UsageText: `This command lists all ebpf profiling task, via id or name in service, instance or process.

Example：
1. Query profiling tasks of service "business-zone::projectC"
$ swctl profiling ebpf list --service-name=service-name
`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.GenericFlag{
				Name:  "trigger",
				Usage: "the trigger type for the profiling task",
				Value: &ebpf.ProfilingTriggerTypeEnumValue{
					Enum:     api.AllEBPFProfilingTriggerType,
					Default:  api.EBPFProfilingTriggerTypeFixedTime,
					Selected: api.EBPFProfilingTriggerTypeFixedTime,
				},
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		processes, err := profiling.QueryEBPFProfilingTaskList(
			ctx.Context, serviceID,
			ctx.Generic("trigger").(*ebpf.ProfilingTriggerTypeEnumValue).Selected)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: processes})
	},
}
