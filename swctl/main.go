/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"encoding/json"
	"github.com/apache/skywalking-cli/commands/service"
	"github.com/apache/skywalking-cli/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/apache/skywalking-cli/logger"
)

var log *logrus.Logger

func init() {
	log = logger.Log
}

func main() {
	app := cli.NewApp()
	app.Usage = "The CLI (Command Line Interface) for Apache SkyWalking."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "./settings.yml",
			Usage: "load configuration `FILE`, default to ./settings.yml",
		},
		cli.BoolFlag{
			Name:     "debug",
			Required: false,
			Usage:    "enable debug mode, will print more detailed logs",
		},
	}
	app.Commands = []cli.Command{
		service.Command,
	}

	app.Before = BeforeCommand

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func BeforeCommand(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(logrus.DebugLevel)
		log.Debugln("Debug mode is enabled")
	}

	configFile := c.String("config")
	log.Debugln("Using configuration file:", configFile)

	if bytes, err := ioutil.ReadFile(configFile); err != nil {
		return err
	} else if err := yaml.Unmarshal(bytes, &config.Config); err != nil {
		return err
	}

	if bytes, err := json.Marshal(config.Config); err == nil {
		log.Debugln("Configurations: ", string(bytes))
	}

	return nil
}
