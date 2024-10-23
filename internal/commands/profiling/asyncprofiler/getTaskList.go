package asyncprofiler

import (
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"
)

var getTaskListCommand = &cli.Command{
	Name:      "list",
	Aliases:   []string{"l"},
	Usage:     "Query async-profiler task list",
	UsageText: ``,
	Flags:     flags.Flags(),
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
