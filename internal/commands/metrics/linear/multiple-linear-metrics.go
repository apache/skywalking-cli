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
$ swctl metrics multiple-linear --name all_percentile

2. Relabel the labels for better readability:
$ swctl metrics multiple-linear --name all_percentile --labels=0,1,2,3,4 --relabels=P50,P75,P90,P95,P99
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
		flags.ProcessRelationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "labels",
				Usage:    "the labels you need to query",
				Required: false,
			},
		},
		[]cli.Flag{
			&cli.StringFlag{
				Name: "relabels",
				Usage: `the new labels to map to the original "--labels", must be in same size and is order-sensitive. ` +
					`"labels[i]" will be mapped to "relabels[i]"`,
				Required: false,
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

		metricsName := ctx.String("name")
		labelsString := ctx.String("labels")
		relabelsString := ctx.String("relabels")

		labels := strings.Split(labelsString, ",")
		relabels := strings.Split(relabelsString, ",")

		labelMapping := make(map[string]string)
		switch {
		case labelsString == "" && relabelsString != "":
			return fmt.Errorf(`"--labels" cannot be empty when "--relabels" is given`)
		case labelsString != "" && relabelsString != "" && len(labels) != len(relabels):
			return fmt.Errorf(`"--labels" and "--relabels" must be in same size if both specified, but was %v != %v`, len(labels), len(relabels))
		case relabelsString != "":
			for i := 0; i < len(labels); i++ {
				labelMapping[labels[i]] = relabels[i]
			}
		}

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
		}, labels, duration)

		if err != nil {
			return err
		}

		reshaped := utils.MetricsValuesArrayToMap(duration, metricsValuesArray, labelMapping)
		return display.Display(ctx, &displayable.Displayable{Data: reshaped})
	},
}
