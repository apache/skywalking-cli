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

package event

import (
	event "skywalking.apache.org/repo/goapi/collect/event/v3"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	eventQl "github.com/apache/skywalking-cli/pkg/graphql/event"

	"github.com/urfave/cli"
)

const DefaultPageSize = 15

var listCommand = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "List events",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			cli.StringFlag{
				Name:     "service",
				Usage:    "service name",
				Required: false,
			},
			cli.StringFlag{
				Name:     "instance",
				Usage:    "service instance name",
				Required: false,
			},
			cli.StringFlag{
				Name:     "endpoint",
				Usage:    "endpoint name",
				Required: false,
			},
			cli.StringFlag{
				Name:     "name",
				Usage:    "event name",
				Required: false,
			},
			cli.GenericFlag{
				Name:  "type",
				Usage: "the type of the event",
				Value: &model.EventTypeEnumValue{
					Enum:     []event.Type{event.Type_Normal, event.Type_Error},
					Default:  event.Type_Normal,
					Selected: event.Type_Normal,
				},
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
		serviceName := ctx.String("service")
		serviceInstanceName := ctx.String("instance")
		endpointName := ctx.String("endpoint")
		name := ctx.String("name")
		eventType := api.EventType(ctx.Generic("type").(*model.EventTypeEnumValue).String())
		pageNum := 1
		needTotal := true

		paging := api.Pagination{
			PageNum:   &pageNum,
			PageSize:  DefaultPageSize,
			NeedTotal: &needTotal,
		}
		condition := &api.EventQueryCondition{
			Source: &api.SourceInput{
				Service:         &serviceName,
				ServiceInstance: &serviceInstanceName,
				Endpoint:        &endpointName,
			},
			Name:   &name,
			Type:   &eventType,
			Time:   &duration,
			Order:  nil,
			Paging: &paging,
		}

		events, err := eventQl.Events(ctx, condition)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		return display.Display(ctx, &displayable.Displayable{Data: events, Condition: condition})
	},
}
