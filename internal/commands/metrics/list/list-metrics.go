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

package list

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"
)

var Command = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List metrics that could be queried",
	UsageText: `With this command, you can list all available metrics
as well as their types and catalogs that can be used in other metrics command
to get the details.

Examples:
1. Query the metrics names start with "service_"
$ swctl metrics list --regex "service_.*"`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "regex",
			Usage: "filter metrics with regex",
		},
	},
	Action: func(ctx *cli.Context) error {
		regex := ctx.String("regex")

		metricsValue, err := metrics.ListMetrics(ctx, regex)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: metricsValue})
	},
}
