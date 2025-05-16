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
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/alarm"

	api "skywalking.apache.org/repo/goapi/query"
)

var autocompleteKeysCommand = &cli.Command{
	Name:    "autocomplete-keys",
	Aliases: []string{"ks"},
	Usage:   "Query autocomplete Keys",
	UsageText: `Query autocomplete keys

Examples:
1. Query autocomplete keys:
$ swctl alarm autocomplete-keys
`,
	Flags: flags.Flags(
		flags.DurationFlags,
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
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

		autocompleteKeys, err := alarm.TagAutocompleteKeys(ctx.Context, duration)
		if err != nil {
			return err
		}

		ctx.Lineage()
		return display.Display(ctx.Context, &displayable.Displayable{Data: autocompleteKeys, Condition: duration})
	},
}
