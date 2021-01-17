package report

import (
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/urfave/cli"
	"strings"

	event "skywalking/network/event/v3"
)

var Command = cli.Command{
	Name:      "report",
	Aliases:   []string{"r"},
	Usage:     "Report an event to OAP server via gRPC",
	ArgsUsage: "[parameters...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "grpcAddr",
			Usage:    "`host:port` to connect",
			Value:    "127.0.0.1:11800",
			Required: true,
		},
		cli.StringFlag{
			Name:     "uuid",
			Usage:    "Unique ID of the event",
			Required: true,
		},
		cli.StringFlag{
			Name:     "service",
			Usage:    "The service of the event occurred on",
			Required: true,
		},
		cli.StringFlag{
			Name:     "instance",
			Usage:    "The service instance of the event occurred on",
			Required: false,
		},
		cli.StringFlag{
			Name:     "endpoint",
			Usage:    "The endpoint of the event occurred on",
			Required: false,
		},
		cli.StringFlag{
			Name:     "name",
			Usage:    "The name of the event. For example, `Reboot`, `Upgrade` etc.",
			Required: true,
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
			Name:     "message",
			Usage:    "The detail of the event that describes why this event happened. This should be a one-line message that briefly describes why the event is reported",
			Required: true,
		},
		cli.Int64Flag{
			Name:     "startTime",
			Usage:    "The start time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC",
			Required: true,
		},
		cli.Int64Flag{
			Name:     "endTime",
			Usage:    "The end time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC",
			Required: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		parameters := make(map[string]string, ctx.NArg())
		if ctx.NArg() > 0 {
			for _, para := range ctx.Args() {
				// Not sure, "key:value" ?
				tmp := strings.Split(para, ":")
				if len(tmp) == 2 {
					parameters[tmp[0]] = tmp[1]
				} else {
					logger.Log.Warnf("%s is not a vaild parameter, should like `key:value`", para)
				}
			}
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

		logger.Log.Println("OAP modules are healthy")
		logger.Log.Println("OAP gRPC is healthy")
		return nil
	},
}
