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

package single

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"

	api "skywalking.apache.org/repo/goapi/query"
)

var Command = &cli.Command{
	Name:  "single",
	Usage: "query single-value metrics defined in backend OAL",
	UsageText: `Query single-value metrics defined in backend OAL.

Examples:
1. Query the traffic load (calls per minute) of service "business-zone:projectC"
$ swctl metrics single --name service_cpm --service-name business-zone::projectC

2. Query the traffic load (calls per minute) of service "business-zone:projectC" endpoint "/projectC/{value}"
$ swctl metrics single --name endpoint_cpm --service-name business-zone::projectC --endpoint-name /projectC/{value}`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
		flags.ProcessRelationFlags,
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

		metricsName := ctx.String("name")
		entity, err := interceptor.ParseEntity(ctx)
		if err != nil {
			return err
		}

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		metricsValue, err := metrics.IntValues(ctx, api.MetricsCondition{
			Name:   metricsName,
			Entity: entity,
		}, duration)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: metricsValue})
	},
}
