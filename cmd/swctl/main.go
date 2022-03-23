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
	"os"
	"runtime"

	"github.com/apache/skywalking-cli/internal/commands/alarm"
	"github.com/apache/skywalking-cli/internal/commands/browser"
	"github.com/apache/skywalking-cli/internal/commands/completion"
	"github.com/apache/skywalking-cli/internal/commands/dashboard"
	"github.com/apache/skywalking-cli/internal/commands/dependency"
	"github.com/apache/skywalking-cli/internal/commands/ebpf/profiling"
	"github.com/apache/skywalking-cli/internal/commands/endpoint"
	"github.com/apache/skywalking-cli/internal/commands/event"
	"github.com/apache/skywalking-cli/internal/commands/healthcheck"
	"github.com/apache/skywalking-cli/internal/commands/install"
	"github.com/apache/skywalking-cli/internal/commands/instance"
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/commands/layer"
	"github.com/apache/skywalking-cli/internal/commands/logs"
	"github.com/apache/skywalking-cli/internal/commands/metrics"
	"github.com/apache/skywalking-cli/internal/commands/process"
	"github.com/apache/skywalking-cli/internal/commands/profile"
	"github.com/apache/skywalking-cli/internal/commands/service"
	"github.com/apache/skywalking-cli/internal/commands/trace"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var log *logrus.Logger
var version string // Will be initialized when building

func init() {
	log = logger.Log

	if runtime.GOOS != "windows" {
		cli.AppHelpTemplate = util.AppHelpTemplate
		cli.CommandHelpTemplate = util.CommandHelpTemplate
		cli.SubcommandHelpTemplate = util.SubcommandHelpTemplate
	}
}

func main() {
	app := cli.NewApp()
	app.Usage = "The CLI (Command Line Interface) for Apache SkyWalking."
	app.UsageText = `Commands in SkyWalking CLI are organized into two levels,
in the form of "swctl --option <level1> --option <level2> --option",
there are options in each level, which should follow right after
the corresponding command, take the following command as example:

	$ swctl --debug service list --start="2019-11-11" --end="2019-11-12"

where "--debug" is is an option of "swctl", and since the "swctl" is
a top-level command, "--debug" is also called global option, and "--start"
is an option of the third level command "list", there is no option for the
second level command "service".

Generally, the second level commands are entity related, there are entities
like "service", "service instance", "metrics" in SkyWalking, and we have
corresponding sub-command like "service"; the third level commands are
operations on the entities, such as "list" command will list all the
services, service instances, etc.`
	app.Version = version

	flags := flags()

	app.Commands = []*cli.Command{
		browser.Command,
		endpoint.Command,
		instance.Command,
		service.Command,
		metrics.Command,
		trace.Command,
		healthcheck.Command,
		dashboard.Command,
		install.Command,
		event.Command,
		logs.Command,
		profile.Command,
		completion.Command,
		dependency.Command,
		alarm.Command,
		layer.Command,
		process.Command,
		profiling.Command,
	}

	app.Before = interceptor.BeforeChain(
		setUpCommandLineContext,
		expandConfigFile,
		tryConfigFile(flags),
	)

	app.Flags = flags
	app.CommandNotFound = util.CommandNotFound

	// Enable auto-completion.
	app.EnableBashCompletion = true
	cli.BashCompletionFlag = &cli.BoolFlag{
		Name:   "auto_complete",
		Hidden: true,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "config",
			Value: "~/.skywalking.yml",
			Usage: "file `path` of the default configurations.",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "base-url",
			Required: false,
			Usage:    "base `url` of the OAP backend graphql service",
			Value:    "http://127.0.0.1:12800/graphql",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "grpc-addr",
			Usage:    "backend gRPC service address `<host:port>`",
			Value:    "127.0.0.1:11800",
			Required: false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "username",
			Required: false,
			Usage:    "`username` of basic authorization",
			Value:    "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "password",
			Required: false,
			Usage:    "`password` of basic authorization",
			Value:    "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "authorization",
			Required: false,
			Usage: "`authorization` header, can be something like `Basic base64(username:password)` or `Bearer jwt-token`, " +
				"if `authorization` is set, `--username` and `--password` are ignored",
			Value: "",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "timezone",
			Required: false,
			Usage: "specifies the `timezone` where `--start` and `--end` are based, in the form of `+0800`. " +
				"If `--timezone` is given in the command line option, then it's used directly. " +
				"If the backend support the timezone API (since 6.5.0), CLI will try to get the timezone from backend, " +
				"and use it. " +
				"Otherwise, the CLI will use the current timezone in the current machine.",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "debug",
			Required: false,
			Usage:    "enable debug mode, will print more detailed information at runtime",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "display",
			Required: false,
			Usage: "display `style` of the result, supported styles are: `json`, `yaml`, `table`, `graph`. " +
				"(Not all display styles are supported in all commands due to data formats incompatibilities " +
				"and the limits of Ascii Graph, like coloring in terminal, " +
				"so in that cases please use `json`  or `yaml` instead.)",
			Value: "",
		}),
	}
}

func expandConfigFile(c *cli.Context) error {
	return c.Set("config", util.ExpandFilePath(c.String("config")))
}

func tryConfigFile(flags []cli.Flag) cli.BeforeFunc {
	return func(c *cli.Context) error {
		configFile := c.String("config")
		if bytes, err := os.ReadFile(configFile); err == nil {
			log.Debug("Using configurations:\n", string(bytes))

			err = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
			if err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			log.Debugf("open %s no such file, skip loading configuration file\n", c.String("config"))
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
