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

package logs

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/log"

	api "skywalking.apache.org/repo/goapi/query"
)

const DefaultPageSize = 15

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List browser error logs according to the specified options",
	UsageText: `List browser error logs according to the specified options.

Examples:
1. Query all logs:
$ swctl browser logs list`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.VersionFlags,
		flags.PageFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseVersion(false),
		interceptor.ParsePage(false),
	),
	Action: func(ctx *cli.Context) error {
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step")

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		serviceID := ctx.String("service-id")
		serviceVersionID := ctx.String("version-id")
		pageID := ctx.String("page-id")

		pageNum := 1

		paging := api.Pagination{
			PageNum:  &pageNum,
			PageSize: DefaultPageSize,
		}
		condition := &api.BrowserErrorLogQueryCondition{
			ServiceID:        &serviceID,
			ServiceVersionID: &serviceVersionID,
			PagePathID:       &pageID,
			QueryDuration:    &duration,
			Paging:           &paging,
		}

		logs, err := log.BrowserLogs(ctx, condition)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: logs, Condition: condition})
	},
}
