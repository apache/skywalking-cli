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

package global

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/dashboard"

	api "skywalking.apache.org/repo/goapi/query"
)

var Metrics = &cli.Command{
	Name:  "global-metrics",
	Usage: "Query global metrics",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "template",
				Usage:    "load dashboard UI template",
				Required: false,
				Value:    dashboard.DefaultTemplatePath,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(cliCtx *cli.Context) error {
		end := cliCtx.String("end")
		start := cliCtx.String("start")
		step := cliCtx.Generic("step")

		ctx := cliCtx.Context
		ctx = context.WithValue(ctx, contextkey.DashboardTemplate{}, cliCtx.String("template"))
		cliCtx.Context = ctx

		globalMetrics, err := dashboard.Metrics(cliCtx.Context, api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})
		if err != nil {
			return err
		}

		return display.Display(cliCtx.Context, &displayable.Displayable{Data: globalMetrics})
	},
}
