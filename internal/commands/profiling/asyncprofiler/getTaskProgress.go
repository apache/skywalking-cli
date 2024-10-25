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
	Usage: "Query async-profiler task progress",
	UsageText: `Query async-profiler task progress

Examples:
1. Query task progress, including task logs and successInstances and errorInstances
$ swctl profiling asyncprofiler progress --task-id=task-id`,
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")

		data, err := profiling.GetAsyncProfilerTaskProgress(ctx, taskID)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: data, Condition: taskID})
	},
}
