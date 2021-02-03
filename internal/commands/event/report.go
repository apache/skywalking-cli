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
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/grpc"

	"github.com/urfave/cli"

	event "skywalking/network/event/v3"
)

var reportCommand = cli.Command{
	Name:      "report",
	Aliases:   []string{"r"},
	Usage:     "Report an event to OAP server via gRPC",
	ArgsUsage: "[parameters...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "uuid",
			Usage: "Unique `ID` of the event.",
		},
		cli.StringFlag{
			Name:  "service",
			Usage: "The service of the event occurred on.",
		},
		cli.StringFlag{
			Name:  "instance",
			Usage: "The service instance of the event occurred on.",
		},
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "The endpoint of the event occurred on",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "The name of the event. For example, 'Reboot' and 'Upgrade' etc.",
		},
		cli.GenericFlag{
			Name:  "type",
			Usage: "The type of the event.",
			Value: &model.EventTypeEnumValue{
				Enum:     []event.Type{event.Type_Normal, event.Type_Error},
				Default:  event.Type_Normal,
				Selected: event.Type_Normal,
			},
		},
		cli.StringFlag{
			Name:  "message",
			Usage: "The detail of the event. This should be a one-line message that briefly describes why the event is reported.",
		},
		cli.Int64Flag{
			Name:  "startTime",
			Usage: "The start time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
		},
		cli.Int64Flag{
			Name:  "endTime",
			Usage: "The end time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
		},
	},
	Action: func(ctx *cli.Context) error {
		parameters, err := interceptor.ParseParameters(ctx.Args())
		if err != nil {
			return err
		}

		event := event.Event{
			Uuid: ctx.String("uuid"),
			Source: &event.Source{
				Service:         ctx.String("service"),
				ServiceInstance: ctx.String("instance"),
				Endpoint:        ctx.String("endpoint"),
			},
			Name:       ctx.String("name"),
			Type:       ctx.Generic("type").(*model.EventTypeEnumValue).Selected,
			Message:    ctx.String("message"),
			Parameters: parameters,
			StartTime:  ctx.Int64("startTime"),
			EndTime:    ctx.Int64("endTime"),
		}

		reply, err := grpc.ReportEvent(ctx.GlobalString("grpcAddr"), &event)
		if err != nil {
			logger.Log.Fatalln(err)
			return err
		}

		logger.Log.Println("Report the event successfully, whose uuid is ", ctx.String("uuid"))
		return display.Display(ctx, &displayable.Displayable{Data: reply})
	},
}
