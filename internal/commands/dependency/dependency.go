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

var Command = cli.Command{
	Name:      "dependency",
	ShortName: "dp",
	Usage:     "Query the dependencies of given endpoint",
	ArgsUsage: "<endpointId>",
	Flags:      flags.Flags(
		flags.DurationFlags,
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return fmt.Errorf("command dependency requires endpointId as argument") 
		}

		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		dependency, err := dependency.Dependency(ctx, ctx.Args().First(), duration)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		return display.Display(ctx, &displayable.Displayable{Data: dependency})
	},
}
