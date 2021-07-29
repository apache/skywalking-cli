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

package logs

import (
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/log"

	"github.com/urfave/cli"
)

const DefaultPageSize = 15

var ListCommand = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "List logs",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "service-id",
				Usage:    "service id",
				Required: false,
			},
			cli.StringFlag{
				Name:     "service-instance-id",
				Usage:    "service instance id",
				Required: false,
			},
			cli.StringFlag{
				Name:     "endpoint-id",
				Usage:    "endpoint id",
				Required: false,
			},
			cli.StringFlag{
				Name:     "trace-id",
				Usage:    "relate trace id",
				Required: false,
			},
			cli.StringFlag{
				Name:     "tags",
				Usage:    "key=value,key=value",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.TimezoneInterceptor,
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step")

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}
		serviceID := ctx.String("service-id")
		serviceInstanceID := ctx.String("service-instance-id")
		endpointID := ctx.String("endpoint-id")
		traceID := ctx.String("trace-id")
		pageNum := 1
		needTotal := true

		tagStr := ctx.String("tags")
		var tags []*api.LogTag = nil
		if tagStr != "" {
			tagArr := strings.Split(tagStr, ",")
			for _, tag := range tagArr {
				kv := strings.Split(tag, "=")
				tags = append(tags, &api.LogTag{Key: kv[0], Value: &kv[1]})
			}
		}

		paging := api.Pagination{
			PageNum:   &pageNum,
			PageSize:  DefaultPageSize,
			NeedTotal: &needTotal,
		}

		condition := &api.LogQueryCondition{
			ServiceID:         &serviceID,
			ServiceInstanceID: &serviceInstanceID,
			EndpointID:        &endpointID,
			RelatedTrace:      &api.TraceScopeCondition{TraceID: traceID},
			Tags:              tags,
			QueryDuration:     &duration,
			Paging:            &paging,
		}
		logs, err := log.Logs(ctx, condition)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		return display.Display(ctx, &displayable.Displayable{Data: logs, Condition: condition})
	},
}
