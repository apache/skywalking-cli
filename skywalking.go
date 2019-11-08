package main

import (
	"github.com/apache/skywalking-cli/cmd"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/logger"
)

const (
	version     = "0.0.1"
	name        = "Apache Skywalking CLI"
	usage       = "Pleace refer to https://github.com/apache/skywalking-cli"
	description = "SkyWalking CLI is a command interaction tool for the SkyWalking user or OPS team, as an alternative besides using browser GUI.\n" +
		"It is based on SkyWalking GraphQL query protocol, same as GUI."
)

var log *logrus.Logger

func init() {
	log = logger.Log
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = name
	app.Usage = usage
	app.Description = description
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Usage: "(optional)Path of settings.yml config. Use the file in the same folder as default. ",
		},
		cli.StringFlag{
			Name:  "debug",
			Usage: "Show all interaction steps, including GraphQL request, response, the local filter execution result.",
		},
		cli.StringFlag{
			Name:  "start-time,st",
			Usage: "start time",
		},
		cli.StringFlag{
			Name:  "end-time,et",
			Usage: "end time",
		},
	}

	app.Commands = []cli.Command{
		cmd.ServiceCmd,
		cmd.InstanceCmd,
		cmd.EndpointCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
