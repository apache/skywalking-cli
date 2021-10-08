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

package endpoint

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"

	"github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/graphql/metadata"

	"github.com/apache/skywalking-cli/pkg/display"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `list all monitored endpoints of the given "--service-id" or "--service-name"`,
	UsageText: `This command lists all endpoints of the service, via service id or service name.

Examples:
1. List all endpoints of the service by service name "business-zone::projectC":
$ swctl endpoint ls --service-name business-zone::projectC

2. List all endpoints of the service by service id "YnVzaW5lc3Mtem9uZTo6cHJvamVjdEM=.1":
$ swctl endpoint ls --service-id YnVzaW5lc3Mtem9uZTo6cHJvamVjdEM=.1

3. Search endpoints like "projectC" of the service "business-zone::projectC":
$ swctl endpoint ls --service-name business-zone::projectC --keyword projectC`,
	Flags: flags.Flags(
		flags.ServiceFlags,

		[]cli.Flag{
			&cli.IntFlag{
				Name:     "limit",
				Usage:    "returns at most `<limit>` endpoints",
				Required: false,
				Value:    100,
			},
			&cli.StringFlag{
				Name:     "keyword",
				Usage:    "`<keyword>` of the endpoint name to search for, empty to search all",
				Required: false,
				Value:    "",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		limit := ctx.Int("limit")
		keyword := ctx.String("keyword")

		endpoints, err := metadata.SearchEndpoints(ctx, serviceID, keyword, limit)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: endpoints})
	},
}
