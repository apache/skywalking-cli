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
	"fmt"

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/graphql/dependency"

	api "skywalking.apache.org/repo/goapi/query"
)

var InstanceCommand = cli.Command{
	Name:      "instance",
	ShortName: "instc",
	Usage:     "Query the instance topology, based on the given clientServiceId and serverServiceId",
	ArgsUsage: "<clientServiceId> <serverServiceId>",
	Flags: flags.Flags(
		flags.DurationFlags,
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),

	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 2 {
			return fmt.Errorf("command instance requires both clientServiceId and serverServiceId as arguments")
		}

		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		dependency, err := dependency.InstanceTopology(ctx, ctx.Args().First(), ctx.Args().Get(1), duration)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		return display.Display(ctx, &displayable.Displayable{Data: dependency})
	},
}
