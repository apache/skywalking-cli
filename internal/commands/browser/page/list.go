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

package page

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `list all monitored browser page of the give "--service-id" or "--service-name"`,
	UsageText: `This command lists all page of the browser service, via service id or service name.

Examples:
1. List all page of the browser service by service name "test-ui":
$ swctl browser page ls --service-name test-ui

2. List all page of the browser service by service id "dGVzdC11aQ==.1":
$ swctl browser page ls --service-id dGVzdC11aQ==.1`,
	Flags: flags.Flags(
		flags.ServiceFlags,

		[]cli.Flag{
			&cli.IntFlag{
				Name:     "limit",
				Usage:    "returns at most `<limit>` endpoints",
				Required: false,
				Value:    100,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseBrowserService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		limit := ctx.Int("limit")

		endpoints, err := metadata.SearchEndpoints(ctx, serviceID, "", limit)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: endpoints})
	},
}
