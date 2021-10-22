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

package version

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   `list all monitored browser version of the given "--service-id" or "--service-name"`,
	UsageText: `This command lists all version of the browser service, via service id or service name.

Examples:
1. List all version of the browser service by service name "provider":
$ swctl browser version ls --service-name test-ui

2. List all version of the browser service by service id "dGVzdC11aQ==.1"":
$ swctl browser version ls --service-id dGVzdC11aQ==.1`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.ServiceFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseBrowserService(true),
	),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")
		serviceID := ctx.String("service-id")

		instances, err := metadata.Instances(ctx, serviceID, api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: instances})
	},
}
