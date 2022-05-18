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

package trace

import (
	"fmt"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/trace"
)

const DefaultPageSize = 15

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "Query the monitored traces",
	UsageText: `Query the monitored traces.

Examples:
1. Query all monitored traces:
$ swctl trace ls

2. Query all monitored traces of service "business-zone::projectB":
$ swctl trace ls --service-name "business-zone::projectB"

3. Query all monitored traces of endpoint "/projectB/{value}" of service "business-zone::projectB":
$ swctl trace ls --service-name "business-zone::projectB" --endpoint-name "/projectB/{value}"

3. Query the monitored trace of id "321661b1-9a31-4e12-ad64-c8f6711f108d":
$ swctl trace ls --trace-id "321661b1-9a31-4e12-ad64-c8f6711f108d"
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.InstanceFlags,
		flags.EndpointFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "trace-id",
				Usage:    "`id` of the trace",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "tags",
				Usage:    "`tags` of the trace, in form of `key=value,key=value`",
				Required: false,
			},
			&cli.StringFlag{
				Name:  "order",
				Usage: "`order` of the returned traces, can be `duration` or `startTime`",
				Value: "duration",
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
		endpointID := ctx.String("endpoint-id")
		serviceInstanceID := ctx.String("instance-id")
		traceID := ctx.String("trace-id")
		tagStr := ctx.String("tags")
		var tags []*api.SpanTag = nil
		if tagStr != "" {
			tagArr := strings.Split(tagStr, ",")
			for _, tag := range tagArr {
				kv := strings.Split(tag, "=")
				tags = append(tags, &api.SpanTag{Key: kv[0], Value: &kv[1]})
			}
		}
		pageNum := 1

		paging := api.Pagination{
			PageNum:  &pageNum,
			PageSize: DefaultPageSize,
		}

		var order api.QueryOrder
		switch orderStr := ctx.String("order"); orderStr {
		case "duration":
			order = api.QueryOrderByDuration
		case "startTime":
			order = api.QueryOrderByStartTime
		default:
			return fmt.Errorf(`invalid order %v, must be one of "duration" or "startTime"`, orderStr)
		}

		condition := &api.TraceQueryCondition{
			ServiceID:         &serviceID,
			ServiceInstanceID: &serviceInstanceID,
			TraceID:           &traceID,
			EndpointID:        &endpointID,
			QueryDuration:     &duration,
			MinTraceDuration:  nil,
			MaxTraceDuration:  nil,
			TraceState:        api.TraceStateAll,
			QueryOrder:        order,
			Tags:              tags,
			Paging:            &paging,
		}
		traces, err := trace.Traces(ctx, condition)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: traces, Condition: condition})
	},
}
