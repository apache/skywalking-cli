package cmd

import (
	"github.com/apache/skywalking-cli/swctl/service"
	"github.com/urfave/cli"
	"log"
)

var ServiceCmd = cli.Command{
	Name:    "service",
	Aliases: []string{"svc"},
	Usage:   "service",
	Action: func(c *cli.Context) {
		ctl := service.NewService(c)

		err := ctl.Exec()
		if err != nil {
			log.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "list,ls",
			Usage: "List all available services.",
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

func runService(c *cli.Context) error {
	return nil
}
