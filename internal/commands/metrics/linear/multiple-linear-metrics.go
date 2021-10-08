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

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"

	"github.com/urfave/cli/v2"
)

var Multiple = &cli.Command{
	Name:  "multiple-linear",
	Usage: "Query multiple linear-type metrics defined in backend OAL",
	UsageText: `Query multiple linear-type metrics defined in backend OAL.

Examples:
1. Query the global percentiles:
$ swctl metrics multiple-linear --name all_percentile`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "labels",
				Usage:    "the labels you need to query",
				Required: false,
				Value:    "0,1,2,3,4",
			},
		},
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
		labels := ctx.String("labels")
		entity, err := interceptor.ParseEntity(ctx)
		if err != nil {
			return err
		}

		if *entity.ServiceName == "" && entity.Scope != api.ScopeAll {
			return fmt.Errorf("the name of service should be specified when metrics' scope is not `All`")
		}

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		metricsValuesArray, err := metrics.MultipleLinearIntValues(ctx, api.MetricsCondition{
			Name:   metricsName,
			Entity: entity,
		}, strings.Split(labels, ","), duration)

		if err != nil {
			return err
		}

		reshaped := utils.MetricsValuesArrayToMap(duration, metricsValuesArray)
		return display.Display(ctx, &displayable.Displayable{Data: reshaped})
	},
}
