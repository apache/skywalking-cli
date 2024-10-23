package asyncprofiler

import (
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:      "asyncprofiler",
	Usage:     "async profiler related sub-command",
	UsageText: `todo`,
	Subcommands: []*cli.Command{
		createCommand,
		getTaskListCommand,
		getTaskProgressCommand,
		analysisCommand,
	},
}
