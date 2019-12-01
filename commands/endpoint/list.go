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
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/graphql/client"
)

var ListCommand = cli.Command{
	Name:        "list",
	ShortName:   "ls",
	Usage:       "List endpoints",
	Description: "list all endpoints if no <endpoint id> is given, otherwise, only list the given endpoint",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "service-id",
			Usage:    "`<service id>` whose endpoints are to be searched",
			Required: true,
		},
		cli.IntFlag{
			Name:     "limit",
			Usage:    "returns at most `<limit>` endpoints",
			Required: false,
			Value:    100,
		},
		cli.StringFlag{
			Name:     "keyword",
			Usage:    "`<keyword>` of the endpoint name to search for, empty to search all",
			Required: false,
			Value:    "",
		},
	},
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		limit := ctx.Int("limit")
		keyword := ctx.String("keyword")

		endpoints := client.SearchEndpoints(ctx, serviceID, keyword, limit)

		return display.Display(ctx, endpoints)
	},
}
