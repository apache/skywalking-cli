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
	"fmt"
	"strings"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/graphql/metrics"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/graphql/utils"

	"github.com/urfave/cli"
)

var Multiple = cli.Command{
	Name:  "multiple-linear",
	Usage: "Query multiple linear metrics defined in backend OAL",
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "labels",
				Usage:    "the labels you need to query",
				Required: false,
				Value:    "0,1,2,3,4",
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
		serviceName := ctx.String("service")
		normal := ctx.BoolT("isNormal")
		instanceName := ctx.String("instance")
		endpointName := ctx.String("endpoint")
		labels := ctx.String("labels")
		scope := interceptor.ParseScope(metricsName)

		if serviceName == "" && scope != schema.ScopeAll {
			return fmt.Errorf("the name of service should be specified when metrics' scope is not `All`")
		}

		duration := schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		metricsValuesArray := metrics.MultipleLinearIntValues(ctx, schema.MetricsCondition{
			Name: metricsName,
			Entity: &schema.Entity{
				Scope:               scope,
				ServiceName:         &serviceName,
				Normal:              &normal,
				ServiceInstanceName: &instanceName,
				EndpointName:        &endpointName,
			},
		}, strings.Split(labels, ","), duration)

		reshaped := utils.MetricsValuesArrayToMap(duration, metricsValuesArray)
		return display.Display(ctx, &displayable.Displayable{Data: reshaped})
	},
}
