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
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/log"

	"github.com/urfave/cli/v2"
)

const DefaultPageSize = 15

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List logs according to the specified options",
	UsageText: `List logs according to the specified options.

Examples:
1. Query all logs:
$ swctl logs list

2. Query the logs related to trace id 3d56f33f-bcd3-4e40-9e4f-5dc547998ef5
$ swctl logs list --trace-id 3d56f33f-bcd3-4e40-9e4f-5dc547998ef5`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.InstanceFlags,
		flags.EndpointFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "trace-id",
				Usage:    "related trace id",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "tags",
				Usage:    "key=value,key=value",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseInstance(false),
		interceptor.ParseEndpoint(false),
	),
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
		serviceInstanceID := ctx.String("instance-id")
		endpointID := ctx.String("endpoint-id")
		traceID := ctx.String("trace-id")
		pageNum := 1

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
			PageNum:  &pageNum,
			PageSize: DefaultPageSize,
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
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: logs, Condition: condition})
	},
}
