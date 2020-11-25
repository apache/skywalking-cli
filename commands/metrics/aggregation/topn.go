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

package aggregation

import (
	"fmt"
	"strconv"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/graphql/metrics"
	"github.com/apache/skywalking-cli/graphql/schema"

	"github.com/urfave/cli"
)

var TopN = cli.Command{
	Name:      "top",
	Usage:     "query top `n` entities",
	ArgsUsage: "<n>",
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		[]cli.Flag{
			cli.GenericFlag{
				Name:  "order",
				Usage: "the `order` by which the top entities are sorted",
				Value: &model.OrderEnumValue{
					Enum:     schema.AllOrder,
					Default:  schema.OrderDes,
					Selected: schema.OrderDes,
				},
			},
		},
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step").(*model.StepEnumValue).Selected

		metricsName := ctx.String("name")
		normal := !ctx.Bool("unnoraml")
		scope := interceptor.ParseScope(metricsName)
		order := ctx.Generic("order").(*model.OrderEnumValue).Selected
		topN := 5
		parentService := ctx.String("service")

		if ctx.NArg() > 0 {
			nn, err := strconv.Atoi(ctx.Args().First())
			if err != nil {
				return fmt.Errorf("the 1st argument must be a number")
			}
			topN = nn
		}

		duration := schema.Duration{
			Start: start,
			End:   end,
			Step:  step,
		}

		metricsValues := metrics.SortMetrics(ctx, schema.TopNCondition{
			Name:          metricsName,
			ParentService: &parentService,
			Normal:        &normal,
			Scope:         &scope,
			TopN:          topN,
			Order:         order,
		}, duration)

		return display.Display(ctx, &displayable.Displayable{Data: metricsValues})
	},
}
