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

package dependency

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/graphql/dependency"

	api "skywalking.apache.org/repo/goapi/query"
)

var ServiceCommand = &cli.Command{
	Name:    "service",
	Aliases: []string{"svc"},
	Usage:   "Query the dependencies of given service",
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.ServiceFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseService(true),
	),

	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")
		coldStage := ctx.Bool("cold")

		duration := api.Duration{
			Start:     start,
			End:       end,
			Step:      step.(*model.StepEnumValue).Selected,
			ColdStage: &coldStage,
		}

		dependency, err := dependency.ServiceTopology(ctx.Context, serviceID, duration)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: dependency})
	},
}
