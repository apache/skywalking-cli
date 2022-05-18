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

package alarm

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/alarm"

	api "skywalking.apache.org/repo/goapi/query"
)

const DefaultPageSize = 15

var listCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List alarms",
	UsageText: `List alarms

Examples:
1. List all alarms:
$ swctl alarm list
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "keyword",
				Usage:    "alarm keyword",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "tags",
				Usage:    "`tags` of the alarm, in form of `key=value,key=value`",
				Required: false,
			},
			&cli.GenericFlag{
				Name:  "scope",
				Usage: "the `scope` of the alarm entity",
				Value: &model.ScopeEnumValue{
					Enum:     api.AllScope,
					Default:  "",
					Selected: "",
				},
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		start := ctx.String("start")
		end := ctx.String("end")
		step := ctx.Generic("step")

		keyword := ctx.String("keyword")
		tagStr := ctx.String("tags")
		scope := ctx.Generic("scope").(*model.ScopeEnumValue).Selected

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		var tags []*api.AlarmTag
		if tagStr != "" {
			tagArr := strings.Split(tagStr, ",")
			for _, tag := range tagArr {
				kv := strings.SplitN(tag, "=", 2)
				if len(kv) != 2 {
					return fmt.Errorf("invalid tag, cannot be splitted into 2 parts. %s", tag)
				}
				tags = append(tags, &api.AlarmTag{Key: kv[0], Value: &kv[1]})
			}
		}

		pageNum := 1
		paging := api.Pagination{
			PageNum:  &pageNum,
			PageSize: DefaultPageSize,
		}

		condition := &alarm.ListAlarmCondition{
			Duration: &duration,
			Keyword:  keyword,
			Scope:    scope,
			Tags:     tags,
			Paging:   &paging,
		}
		alarms, err := alarm.Alarms(ctx, condition)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: alarms, Condition: condition})
	},
}
