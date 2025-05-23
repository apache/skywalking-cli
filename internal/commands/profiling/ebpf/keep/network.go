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

package keep

import (
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"
)

var NetworkKeepCommand = &cli.Command{
	Name:    "network",
	Aliases: []string{"net"},
	Usage:   "Keep alive the exist ebpf network profiling task",
	UsageText: `Keep alive the exist ebpf network profiling task

Examples:
1. Keep alive the ebpf network profiling task
$ swctl profiling ebpf keep network --task-id=abc`,
	Flags: flags.Flags(
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "the `task-id` of the network profiling task",
				Required: true,
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")

		keepResult, err := profiling.KeepNetworkProfilingTask(ctx.Context, taskID)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: keepResult})
	},
}
