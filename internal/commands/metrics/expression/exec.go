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

package expression

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

const expressionParameterName = "expression"

var ExecCommand = &cli.Command{
	Name:    "execute",
	Aliases: []string{"exec"},
	Usage:   `Execute a metrics expression.`,
	UsageText: `Execute a metrics expression.

Examples:
1. Execute the given expression of service "business-zone::projectB"
$ swctl metrics execute --expression="service_resp_time" --service-name business-zone::projectB`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
		flags.ProcessRelationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     expressionParameterName,
				Usage:    "metrics expression",
				Required: true,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseEndpointRelation(false),
		interceptor.ParseInstanceRelation(false),
		interceptor.ParseProcessRelation(false),
	),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		expression := ctx.String(expressionParameterName)
		entity, err := interceptor.ParseEntity(ctx)
		if err != nil {
			return err
		}

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		result, err := metrics.Execute(ctx, expression, entity, duration)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: result})
	},
}
