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

	"github.com/apache/skywalking-cli/api"
	"github.com/apache/skywalking-cli/internal/logger"

	"github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/graphql/metadata"

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
)

var SearchCommand = cli.Command{
	Name:  "search",
	Usage: "Filter the instance from the existing service instance list",
	Flags: append(flags.DurationFlags, append(flags.SearchRegexFlags, flags.InstanceServiceIDFlags...)...),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		serviceID := verifyAndSwitchServiceParameter(ctx)

		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		regex := ctx.String("regex")

		instances, err := metadata.Instances(ctx, serviceID, api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		if err != nil {
			logger.Log.Fatalln(err)
		}

		var result []api.ServiceInstance
		if len(instances) > 0 {
			for _, instance := range instances {
				if ok, _ := regexp.Match(regex, []byte(instance.Name)); ok {
					result = append(result, instance)
				}
			}
		}
		return display.Display(ctx, &displayable.Displayable{Data: result})
	},
}
