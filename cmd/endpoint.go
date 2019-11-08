package cmd

import (
	"github.com/urfave/cli"
)

var EndpointCmd = cli.Command{
	Name:    "endpoint",
	Aliases: []string{"ep"},
	Usage:   "endpoint",
	Action: func(c *cli.Context) error {
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "service,svc",
			Usage: "(Required)Set the service code in current command context.",
		},
		cli.StringFlag{
			Name:  "search",
			Usage: "Use the backend endpoint search capability to get the endpoint list by given --service parameter.",
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

func runEndpoint(c *cli.Context) error {
	return nil
}
