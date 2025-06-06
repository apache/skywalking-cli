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

package instance

import (
	"regexp"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var SearchCommand = &cli.Command{
	Name:  "search",
	Usage: "Filter the instance from the existing service instance list",
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.SearchRegexFlags,
		flags.ServiceFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")
		coldStage := ctx.Bool("cold")
		regex := ctx.String("regex")

		instances, err := metadata.Instances(ctx.Context, serviceID, api.Duration{
			Start:     start,
			End:       end,
			Step:      step.(*model.StepEnumValue).Selected,
			ColdStage: &coldStage,
		})
		if err != nil {
			return err
		}

		var result []api.ServiceInstance
		if len(instances) > 0 {
			for _, instance := range instances {
				if ok, _ := regexp.Match(regex, []byte(instance.Name)); ok {
					result = append(result, instance)
				}
			}
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}
