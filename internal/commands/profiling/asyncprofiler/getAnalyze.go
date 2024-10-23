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
	Name:      "analysis",
	Aliases:   []string{"pa"},
	Usage:     "Analyze profiled stacktrace",
	UsageText: ``,
	Flags: flags.Flags(
		flags.EndpointFlags,

		[]cli.Flag{
			&cli.IntFlag{
				Name:  "task-id",
				Usage: "async-profiler task id",
			},
			&cli.StringSliceFlag{
				Name:  "service-instance-ids",
				Usage: "select which service instances' jfr files to analyze",
			},
			&cli.StringFlag{
				Name:  "event",
				Usage: "which event types this task needs to collect.",
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		event := ctx.String("events")
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
