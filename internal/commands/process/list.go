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

package process

import (
	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `list all monitored processes under the instance`,
	UsageText: `This command lists all processes of the service-instance.

Examples:
3. List all processes by instance name "provider-01" and service name "provider":
$ swctl process ls --instance-name provider-01 --service-name provider

4. List all processes by instance id "cHJvdmlkZXI=.1_cHJvdmlkZXIx":
$ swctl process ls --instance-id cHJvdmlkZXI=.1_cHJvdmlkZXIx`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceFlags,
		flags.DurationFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseInstance(true),
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		instanceID := ctx.String("instance-id")
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		processes, err := metadata.Processes(ctx, instanceID, api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: processes})
	},
}
