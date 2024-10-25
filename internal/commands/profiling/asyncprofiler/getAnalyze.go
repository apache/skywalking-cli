package asyncprofiler

import (
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
	"github.com/urfave/cli/v2"
	query "skywalking.apache.org/repo/goapi/query"
	"strings"
)

var analysisCommand = &cli.Command{
	Name:    "analysis",
	Aliases: []string{"a"},
	Usage:   "Query async-profiler analysis",
	UsageText: `Query async-profiler analysis

Examples:
1. Query the flame graph produced by async-profiler
$ swctl profiling asyncprofiler analysis  --task-id=task-id --service-instance-ids=instanceIds --event=execution_sample`,
	Flags: flags.Flags(
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "async-profiler task id",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "service-instance-ids",
				Usage:    "select which service instances' jfr files to analyze",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "event",
				Usage:    "which event types this task needs to collect.",
				Required: true,
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		event := ctx.String("event")
		eventType := query.JFREventType(strings.ToUpper(event))

		request := &query.AsyncProfilerAnalyzationRequest{
			TaskID:      ctx.String("task-id"),
			InstanceIds: ctx.StringSlice("service-instance-ids"),
			EventType:   eventType,
		}

		analyze, err := profiling.GetAsyncProfilerAnalyze(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: analyze, Condition: request})
	},
}
