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

package linear

import (
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/graphql/metrics"
	"github.com/apache/skywalking-cli/graphql/utils"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/schema"
)

var Multiple = cli.Command{
	Name:  "multiple-linear",
	Usage: "Query multiple linear metrics defined in backend OAL",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    "metrics `NAME`, such as `all_percentile`",
				Required: true,
			},
			cli.StringFlag{
				Name:     "id",
				Usage:    "`ID`, the related id if the metrics requires one",
				Required: false,
			},
			cli.IntFlag{
				Name:     "num",
				Usage:    "`num`, the number of linear metrics to query, (default: 5)",
				Required: false,
				Value:    5,
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
		numOfLinear := ctx.Int("num")

		var id *string = nil

		if idString := ctx.String("id"); idString != "" {
			id = &idString
		}

		duration := schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		values := metrics.MultipleLinearIntValues(ctx, schema.MetricCondition{
			Name: metricsName,
			ID:   id,
		}, numOfLinear, duration)

		reshaped := make([]map[string]float64, len(values))

		for index, value := range values {
			reshaped[index] = utils.MetricsToMap(duration, value)
		}

		return display.Display(ctx, reshaped)
	},
}
