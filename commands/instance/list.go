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
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
)

var ListCommand = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "List all available instance by given --service-id or --service-name parameter",
	Flags:     flags.InstanceServiceIDFlags,
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		serviceName := ctx.String("service-name")

		if serviceID == "" && serviceName == "" {
			logger.Log.Fatalf("flags \"service-id, service-name\" must set one")
		}

		if serviceID == "" && serviceName != "" {
			service, err := client.SearchService(ctx, serviceName)
			if err != nil {
				logger.Log.Fatalln(err)
			}
			serviceID = service.ID
		}

		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		instances := client.Instances(ctx, serviceID, schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		return display.Display(ctx, instances)
	},
}
