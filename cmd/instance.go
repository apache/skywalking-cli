package cmd

import (
	"github.com/urfave/cli"
)

var InstanceCmd = cli.Command{
	Name:    "instance",
	Aliases: []string{"inst"},
	Usage:   "instance of service",
	Action: func(c *cli.Context) error {
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "list,ls",
			Usage: "List all available instances.",
		},
		cli.StringFlag{
			Name:  "service,svc",
			Usage: "Set the service code in current command context.",
		},
		cli.StringFlag{
			Name:  "search",
			Usage: "Filter the instance from the existing service instance list by given --service parameter.",
		},
		cli.StringFlag{
			Name:  "metrics-value,mv",
			Usage: "Metrics value in the given duration and metrics name.",
		},
		cli.StringFlag{
			Name:  "metrics-linear,ml",
			Usage: "Show the metrics linear graph based on response values.",
		},
	},
}

func runInstance(c *cli.Context) error {
	return nil
}
