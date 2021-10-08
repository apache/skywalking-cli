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

package service

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
	Name:      "list",
	Aliases:   []string{"ls"},
	Usage:     `list the monitored services`,
	ArgsUsage: "<service name>",
	UsageText: `This command lists all services if no "<service name>" is given, 
otherwise, it only lists the services matching the given "<service name>".

Examples:
1. List all the services:
$ swctl svc ls

2. List a specific service named "projectC":
$ swctl svc ls projectC`,
	Flags: flags.DurationFlags,
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		var services []api.Service
		var err error

		if args := ctx.Args(); args.Len() == 0 {
			services, err = metadata.AllServices(ctx, api.Duration{
				Start: start,
				End:   end,
				Step:  step.(*model.StepEnumValue).Selected,
			})
			if err != nil {
				return err
			}
		} else {
			service, err := metadata.SearchService(ctx, args.First())
			if err != nil {
				return err
			}
			services = []api.Service{service}
		}

		return display.Display(ctx, &displayable.Displayable{Data: services})
	},
}
