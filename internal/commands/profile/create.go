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

package profile

import (
	"fmt"

	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
	"github.com/apache/skywalking-cli/pkg/graphql/profile"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name:      "create",
	Aliases:   []string{"c"},
	Usage:     "create a new profile task",
	ArgsUsage: "[parameters...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service-id",
			Usage: "<service-id> whose endpoints are to be profile.",
		},
		cli.StringFlag{
			Name:  "service-name",
			Usage: "<service-name> whose endpoints are to be profile.",
		},
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "which endpoint should profile.",
		},
		cli.Int64Flag{
			Name:  "start-time",
			Usage: "profile task start time(millisecond).",
		},
		cli.IntFlag{
			Name:  "duration",
			Usage: "profile task continuous time(minute).",
		},
		cli.IntFlag{
			Name:  "min-duration-threshold",
			Usage: "profiled endpoint must greater duration(millisecond).",
		},
		cli.IntFlag{
			Name:  "dump-period",
			Usage: "profiled endpoint dump period(millisecond).",
		},
		cli.IntFlag{
			Name:  "max-sampling-count",
			Usage: "profile task max sampling count.",
		},
	},
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		if serviceID == "" {
			serviceName := ctx.String("service-name")
			if serviceName == "" {
				return fmt.Errorf(`either flags "service-id" or "service-name" must be set`)
			}
			service, err := metadata.SearchService(ctx, serviceName)
			if err != nil {
				return err
			}
			serviceID = service.ID
		}

		startTime := ctx.Int64("start-time")
		request := &api.ProfileTaskCreationRequest{
			ServiceID:            serviceID,
			EndpointName:         ctx.String("endpoint"),
			StartTime:            &startTime,
			Duration:             ctx.Int("duration"),
			MinDurationThreshold: ctx.Int("min-duration-threshold"),
			DumpPeriod:           ctx.Int("dump-period"),
			MaxSamplingCount:     ctx.Int("max-sampling-count"),
		}

		task, err := profile.CreateTask(ctx, request)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
