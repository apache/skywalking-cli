package instance

import (
	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "List all available instance by given --service parameter",
	Flags:     flags.InstanceServiceIdFlags,
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		serviceId := ctx.String("service")
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")
		services := client.Instances(ctx, serviceId, schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		return display.Display(ctx, services)
	},
}
