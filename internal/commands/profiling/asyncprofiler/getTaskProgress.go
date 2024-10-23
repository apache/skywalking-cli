package asyncprofiler

import (
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
)

var getTaskProgressCommand = &cli.Command{
	Name:    "progress",
	Aliases: []string{"logs", "p"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "task-id",
			Usage: "async profiler task id.",
		},
	},
	Usage: "Query trace profiling task log list",
	UsageText: `Query trace profiling task log list

Examples:
1. Query all trace profiling logs of task id "task-id"
$ swctl profiling trace logs --task-id=task-id`,
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")

		data, err := profiling.GetAsyncProfilerTaskProgress(ctx, taskID)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: data, Condition: taskID})
	},
}
