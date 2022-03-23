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
	"fmt"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	ebpf_graphql "github.com/apache/skywalking-cli/pkg/graphql/ebpf"
)

var ListTaskCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `query ebpf profiling task list`,
	UsageText: `This command lists all ebpf profiling task, via id or name in service, instance or process.

Example：
1. Query profiling tasks of service "business-zone::projectC"
$ swctl ebpf-profiling list --service-name=service-name

2. Query profiling tasks of instance name "provider-01" and service name "provider":
$ swctl ebpf-profiling list --instance-name provider-01 --service-name provider

2. Query profiling tasks of process id "abc"
$ swctl ebpf-profiling list --process-id=abc
`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "process-id",
				Usage:    "the `process-id` by which process ID need to be profiling",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(false),
		interceptor.ParseInstance(false),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		instanceID := ctx.String("instance-id")
		processID := ctx.String("process-id")
		if serviceID == "" && instanceID == "" && processID == "" {
			return fmt.Errorf("service, instance, or process must provide one")
		}

		processes, err := ebpf_graphql.QueryEBPFProfilingTaskList(ctx, &api.EBPFProfilingTaskCondition{
			FinderType: nil,
			ServiceID:  &serviceID,
			InstanceID: &instanceID,
			ProcessID:  &processID,
		})
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: processes})
	},
}
