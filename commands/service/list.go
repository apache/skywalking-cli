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
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/display/displayable"

	"github.com/apache/skywalking-cli/graphql/metadata"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/schema"
)

var ListCommand = cli.Command{
	Name:        "list",
	ShortName:   "ls",
	Usage:       "List services",
	ArgsUsage:   "<service name>",
	Description: "list all services if no <service name> is given, otherwise, only list the given service",
	Flags:       flags.DurationFlags,
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		var services []schema.Service

		if args := ctx.Args(); len(args) == 0 {
			services = metadata.AllServices(ctx, schema.Duration{
				Start: start,
				End:   end,
				Step:  step.(*model.StepEnumValue).Selected,
			})
		} else {
			service, _ := metadata.SearchService(ctx, args.First())
			services = []schema.Service{service}
		}

		return display.Display(ctx, &displayable.Displayable{Data: services})
	},
}
