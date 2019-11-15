package instance

import (
	"github.com/apache/skywalking-cli/commands/flags"
	"github.com/apache/skywalking-cli/commands/interceptor"
	"github.com/apache/skywalking-cli/commands/model"
	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/urfave/cli"
	"regexp"
)

var SearchCommand = cli.Command{
	Name:  "search",
	Usage: "Filter the instance from the existing service instance list by given --service-id or --service-name parameter",
	Flags: append(flags.DurationFlags, append(flags.InstanceServiceIDFlags, flags.SearchRegexFlags...)...),
	Before: interceptor.BeforeChain([]cli.BeforeFunc{
		interceptor.DurationInterceptor,
	}),
	Action: func(ctx *cli.Context) error {
		serviceID := verifyAndSwitch(ctx)

		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		regex := ctx.String("regex")

		instances := client.Instances(ctx, serviceID, schema.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		var result []schema.ServiceInstance
		if len(instances) > 0 {
			for _, instance := range instances {
				if ok, _ := regexp.Match(regex, []byte(instance.Name)); ok {
					result = append(result, instance)
				}
			}
		}
		return display.Display(ctx, result)
	},
}
