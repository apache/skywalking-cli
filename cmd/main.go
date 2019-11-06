package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/logger"
)

var log *logrus.Logger

func init() {
	log = logger.Log
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path of settings.yml config. Use the file in the same folder as default.",
		},
	}
	app.Commands = []cli.Command{serviceCmd}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
