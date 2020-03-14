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

package thermodynamic

import (
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/display/displayable"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/metrics"
	"github.com/apache/skywalking-cli/graphql/schema"
)

var Command = cli.Command{
	Name:      "thermodynamic",
	ShortName: "td",
	Usage:     "Query thermodynamic metrics defined in backend OAL",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    "metrics `name`, which should be defined in OAL script",
				Required: true,
			},
			cli.StringFlag{
				Name:     "id",
				Usage:    "metrics `id` if the metrics require one",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")
		metricsName := ctx.String("name")

		var id *string = nil
		if ctx.String("id") != "" {
			idString := ctx.String("id")
			id = &idString
		}

		duration := schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		metricsValues := metrics.Thermodynamic(ctx, schema.MetricCondition{
			Name: metricsName,
			ID:   id,
		}, duration)

		return display.Display(ctx, &displayable.Displayable{
			Data:     metricsValues,
			Duration: duration,
			Title:    metricsName,
		})
	},
}
