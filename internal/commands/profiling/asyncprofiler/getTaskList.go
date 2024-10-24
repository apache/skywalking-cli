package asyncprofiler

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"
)

var getTaskListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Query async-profiler task list",
	UsageText: `Query async-profiler task list

Examples:
1. Query all async-profiler tasks
$ swctl profiling asyncprofiler list --service-name=TEST_AGENT`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.Int64Flag{
				Name:  "start-time",
				Usage: "The start time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
			},
			&cli.Int64Flag{
				Name:  "end-time",
				Usage: "The end time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC.",
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit defines the number of the tasks to be returned.",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		startTime := ctx.Int64("start-time")
		endTime := ctx.Int64("end-time")
		limit := ctx.Int("limit")
		request := &query.AsyncProfilerTaskListRequest{
			ServiceID: serviceID,
			StartTime: &startTime,
			EndTime:   &endTime,
			Limit:     &limit,
		}

		tasks, err := profiling.GetAsyncProfilerTaskList(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: tasks, Condition: request})
	},
}
