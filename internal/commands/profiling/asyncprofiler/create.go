package asyncprofiler

import (
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"
	"strings"
)

var createCommand = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new async profiler task",
	UsageText: `Create a new async profiler task

Examples:
1. Create async-profiler task
$ swctl profiling asyncprofiler create todo`,
	Flags: flags.Flags(
		flags.EndpointFlags,

		[]cli.Flag{
			&cli.StringSliceFlag{
				Name:  "service-instance-ids",
				Usage: "which instances to execute task.",
			},
			&cli.IntFlag{
				Name:  "duration",
				Usage: "task continuous time(second).",
			},
			&cli.StringSliceFlag{
				Name:  "events",
				Usage: "which event types this task needs to collect.",
			},
			&cli.StringFlag{
				Name:  "exec-args",
				Usage: "other async-profiler execution options, e.g. alloc=2k,lock=2s.",
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		events := ctx.StringSlice("events")
		eventTypes := make([]query.AsyncProfilerEventType, 0)
		for _, event := range events {
			upperCaseEvent := strings.ToUpper(event)
			eventTypes = append(eventTypes, query.AsyncProfilerEventType(upperCaseEvent))
		}

		execArgs := ctx.String("exec-args")

		request := &query.AsyncProfilerTaskCreationRequest{
			ServiceID:          serviceID,
			ServiceInstanceIds: ctx.StringSlice("service-instance-ids"),
			Duration:           ctx.Int("duration"),
			Events:             make([]query.AsyncProfilerEventType, 0),
			ExecArgs:           &execArgs,
		}
		task, err := profiling.CreateAsyncProfilerTask(ctx, request)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}
