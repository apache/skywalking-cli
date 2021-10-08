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
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"
)

var Single = &cli.Command{
	Name:  "linear",
	Usage: "Query linear-type metrics defined in backend OAL",
	UsageText: `Query linear-type metrics defined in backend OAL

Examples:
1. Query the response time of service "business-zone::projectB"
$ swctl metrics linear --name=service_resp_time --service-name business-zone::projectB

2. Query the response time of service instance
$ swctl metrics linear --name=service_instance_resp_time --service-name business-zone::projectB \
	--instance-name d708c6bfea9f4d50902d1743302a6f50@10.170.0.12
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseInstanceRelation(false),
		interceptor.ParseEndpointRelation(false),
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

		metricsValues, err := metrics.LinearIntValues(ctx, api.MetricsCondition{
			Name:   metricsName,
			Entity: entity,
		}, duration)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: utils.MetricsValuesToMap(duration, metricsValues)})
	},
}
