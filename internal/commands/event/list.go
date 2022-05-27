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
	"strings"

	event "skywalking.apache.org/repo/goapi/collect/event/v3"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	eventQl "github.com/apache/skywalking-cli/pkg/graphql/event"

	"github.com/urfave/cli/v2"
)

const DefaultPageSize = 15
const EventTypeAll event.Type = -1

func init() {
	event.Type_name[-1] = "All"
	event.Type_value["All"] = -1
}

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List events",
	UsageText: `List events

Examples:
1. List all events:
$ swctl event list
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.InstanceFlags,
		flags.EndpointFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "event name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "layer",
				Usage:    "Name of the layer to which the event belongs (case-insensitive), which can be queried via 'swctl layer list'",
				Required: false,
			},
			&cli.GenericFlag{
				Name:  "type",
				Usage: "the type of the event",
				Value: &model.EventTypeEnumValue{
					Enum:     []event.Type{event.Type_Normal, event.Type_Error},
					Default:  EventTypeAll,
					Selected: EventTypeAll,
				},
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseService(false),
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
		serviceName := ctx.String("service-name")
		serviceInstanceName := ctx.String("instance-name")
		endpointName := ctx.String("endpoint-name")
		name := ctx.String("name")
		eventType := ctx.Generic("type").(*model.EventTypeEnumValue).Selected
		layer := strings.ToUpper(ctx.String("layer"))
		pageNum := 1

		paging := api.Pagination{
			PageNum:  &pageNum,
			PageSize: DefaultPageSize,
		}
		condition := &api.EventQueryCondition{
			Source: &api.SourceInput{
				Service:         &serviceName,
				ServiceInstance: &serviceInstanceName,
				Endpoint:        &endpointName,
			},
			Name:   &name,
			Time:   &duration,
			Layer:  &layer,
			Order:  nil,
			Paging: &paging,
		}
		if eventType != EventTypeAll {
			t := api.EventType(eventType.String())
			condition.Type = &t
		}

		events, err := eventQl.Events(ctx, condition)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: events, Condition: condition})
	},
}
