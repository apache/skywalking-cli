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

	"github.com/apache/skywalking-cli/display/displayable"

	"github.com/apache/skywalking-cli/graphql/metrics"
	"github.com/apache/skywalking-cli/graphql/utils"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/schema"
)

var Single = cli.Command{
	Name:  "linear",
	Usage: "Query linear metrics defined in backend OAL",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    "metrics `NAME`, such as `all_p99`",
				Required: true,
			},
			cli.GenericFlag{
				Name:  "scope",
				Usage: "the scope of the query, which follows the metrics `name`",
				Value: &model.ScopeEnumValue{
					Enum:     schema.AllScope,
					Default:  schema.ScopeAll,
					Selected: schema.ScopeAll,
				},
				Required: true,
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
		scope := ctx.Generic("scope").(*model.ScopeEnumValue).Selected

		duration := schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		metricsValues := metrics.LinearIntValues(ctx, schema.MetricsCondition{
			Name: metricsName,
			Entity: &schema.Entity{
				Scope: scope,
			},
		}, duration)

		return display.Display(ctx, &displayable.Displayable{Data: utils.MetricsValuesToMap(duration, metricsValues)})
	},
}
