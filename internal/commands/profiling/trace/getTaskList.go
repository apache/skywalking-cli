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

package trace

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"
)

var getTaskListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Query trace profiling task list",
	UsageText: `Query trace profiling task list

Examples:
1. Query all trace profiling tasks
$ swctl profiling trace list --service-name=service-name --endpoint-name=endpoint

2. Query trace profiling tasks of service "business-zone::projectC", endpoint "/projectC/{value}"
$ swctl profiling trace list --service-name=business-zone::projectC --endpoint-name=/projectC/{value}`,
	Flags: flags.Flags(
		flags.EndpointFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseEndpoint(false),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		endpoint := ctx.String("endpoint-name")

		task, err := profiling.GetTraceProfilingTaskList(ctx, serviceID, endpoint)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: serviceID})
	},
}
