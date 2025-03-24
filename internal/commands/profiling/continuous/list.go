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

package continuous

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `query service continuous profiling policy`,
	UsageText: `This command lists all service continuous profiling policy, via id or name in service.

Example：
1. Query continuous profiling policy of service "business-zone::projectC"
$ swctl profiling continuous list --service-name=business-zone::projectC
`,
	Flags: flags.Flags(
		flags.ServiceFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		result, err := profiling.QueryContinuousProfilingServiceTargets(ctx.Context, serviceID)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}
