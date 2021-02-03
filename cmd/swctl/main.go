// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"io/ioutil"
	"os"

	"github.com/apache/skywalking-cli/internal/commands/dashboard"
	"github.com/apache/skywalking-cli/internal/commands/endpoint"
	"github.com/apache/skywalking-cli/internal/commands/event"
	"github.com/apache/skywalking-cli/internal/commands/healthcheck"
	"github.com/apache/skywalking-cli/internal/commands/install"
	"github.com/apache/skywalking-cli/internal/commands/instance"
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/commands/metrics"
	"github.com/apache/skywalking-cli/internal/commands/service"
	"github.com/apache/skywalking-cli/internal/commands/trace"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var log *logrus.Logger
var version string // Will be initialized when building

func init() {
	log = logger.Log
}

func main() {
	app := cli.NewApp()
	app.Usage = "The CLI (Command Line Interface) for Apache SkyWalking."
	app.Version = version

	flags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "config",
			Value: "~/.skywalking.yml",
			Usage: "load configuration `FILE`",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "base-url",
			Required: false,
			Usage:    "base `url` of the OAP backend graphql",
			Value:    "http://127.0.0.1:12800/graphql",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "grpcAddr",
			Usage:    "`host:port` to connect",
			Value:    "127.0.0.1:11800",
			Required: false,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "username",
			Required: false,
			Usage:    "username of basic authorization",
			Value:    "",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "password",
			Required: false,
			Usage:    "password of basic authorization",
			Value:    "",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "authorization",
			Required: false,
			Usage:    "authorization to the OAP backend",
			Value:    "",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "timezone",
			Required: false,
			Hidden:   true,
			Usage:    "the timezone of the server side",
		}),
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:     "debug",
			Required: false,
			Usage:    "enable debug mode, will print more detailed logs",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "display",
			Required: false,
			Usage:    "display `style` of the result, supported styles are: json, yaml, table, graph",
			Value:    "json",
		}),
	}

	app.Commands = []cli.Command{
		endpoint.Command,
		instance.Command,
		service.Command,
		metrics.Command,
		trace.Command,
		healthcheck.Command,
		dashboard.Command,
		install.Command,
		event.Command,
	}

	app.Before = interceptor.BeforeChain([]cli.BeforeFunc{
		setUpCommandLineContext,
		expandConfigFile,
		tryConfigFile(flags),
	})

	app.Flags = flags

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func expandConfigFile(c *cli.Context) error {
	return c.Set("config", util.ExpandFilePath(c.String("config")))
}

func tryConfigFile(flags []cli.Flag) cli.BeforeFunc {
	return func(c *cli.Context) error {
		configFile := c.String("config")
		if bytes, err := ioutil.ReadFile(configFile); err == nil {
			log.Debug("Using configurations:\n", string(bytes))

			err = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
			if err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			log.Debugf("open %s no such file, skip loading configuration file\n", c.GlobalString("config"))
		} else {
			return err
		}

		return nil
	}
}

func setUpCommandLineContext(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(logrus.DebugLevel)
		log.Debugln("Debug mode is enabled")
	}

	return nil
}
