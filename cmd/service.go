package main

import (
	"github.com/apache/skywalking-cli/swctl/service"
	"github.com/urfave/cli"
)

var serviceCmd = cli.Command{
	Name: "service",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "list",
			Usage: "list all available services.",
		},
	},
	Action: func(c *cli.Context) {
		ctl := service.NewService(c)

		err := ctl.Exec()
		if err != nil {
			log.Fatal(err)
		}
	},
}
