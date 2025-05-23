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
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"
)

var PrepareCreateCommand = &cli.Command{
	Name:    "prepare",
	Aliases: []string{"pre"},
	Usage:   "Prepare to a new ebpf profiling task, typically used to query data before creating a task",
	UsageText: `Prepare to a new ebpf profiling task, typically used to query data before creating a task

Examples:
1. Prepare ebpf profiling fixed time task
$ swctl profiling ebpf create prepare --service-id=abc`,
	Flags: flags.Flags(
		flags.ServiceFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		prepare, err := profiling.QueryPrepareCreateEBPFProfilingTaskData(ctx.Context, serviceID)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: prepare, Condition: serviceID})
	},
}
