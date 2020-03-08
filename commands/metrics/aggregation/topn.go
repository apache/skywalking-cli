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
	"strings"

	"github.com/apache/skywalking-cli/display/displayable"

	"github.com/apache/skywalking-cli/commands/interceptor"

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/aggregation"
	"github.com/apache/skywalking-cli/graphql/schema"
)

var TopN = cli.Command{
	Name:      "top",
	Usage:     "query top `n` entities",
	ArgsUsage: "<n>",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    "`metrics name`, which should be defined in OAL script",
				Required: true,
			},
			cli.GenericFlag{
				Name:  "order",
				Usage: "the `order` by which the top entities are sorted",
				Value: &model.OrderEnumValue{
					Enum:     schema.AllOrder,
					Default:  schema.OrderDes,
					Selected: schema.OrderDes,
				},
			},
			cli.StringFlag{
				Name:     "service-id",
				Usage:    "the `service id` whose instances/endpoints are to be fetch, if applicable",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		name := ctx.String("name")
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step").(*model.StepEnumValue).Selected
		order := ctx.Generic("order").(*model.OrderEnumValue).Selected
		serviceID := ctx.String("service-id")

		topN := 5

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

		var metricsValues []schema.TopNEntity

		if strings.HasPrefix(name, "service_instance") {
			if serviceID == "" {
				metricsValues = aggregation.AllServiceInstanceTopN(ctx, name, topN, duration, order)
			} else {
				metricsValues = aggregation.ServiceInstanceTopN(ctx, serviceID, name, topN, duration, order)
			}
		} else if strings.HasPrefix(name, "endpoint_") {
			if serviceID == "" {
				metricsValues = aggregation.AllEndpointTopN(ctx, name, topN, duration, order)
			} else {
				metricsValues = aggregation.EndpointTopN(ctx, serviceID, name, topN, duration, order)
			}
		} else if strings.HasPrefix(name, "service_") {
			metricsValues = aggregation.ServiceTopN(ctx, name, topN, duration, order)
		}

		return display.Display(ctx, &displayable.Displayable{Data: metricsValues})
	},
}
