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

package global

import (
	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/dashboard"
)

var GlobalCommand = &cli.Command{
	Name:    "global",
	Aliases: []string{"g"},
	Usage:   "Display global dashboard",
	UsageText: `Display global dashboard

Examples:
1. Display the global dashboard
$ swctl dashboard global

2. Display the global dashboard with a customized template
$ swctl dashboard global --template my-global-template.yml
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "template",
				Usage:    "load dashboard UI template",
				Required: false,
				Value:    dashboard.DefaultTemplatePath,
			},
			&cli.IntFlag{
				Name:     "refresh",
				Usage:    "the auto refreshing interval (s)",
				Required: false,
				Value:    6,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		globalData, err := dashboard.Global(ctx, api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		})

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: globalData})
	},
}
