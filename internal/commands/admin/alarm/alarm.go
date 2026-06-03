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

// Package alarm exposes the admin-server alarm runtime status (loaded rule
// definitions and per-entity evaluation/window state). This is distinct from the
// top-level `swctl alarm` command, which reads fired alarm records via GraphQL.
package alarm

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/admin/status"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "alarm",
	Usage: "Inspect alarm runtime status from the admin-server `status` module",
	UsageText: `Inspect the alarm-running kernel: loaded rule definitions and per-entity
evaluation/window state. This differs from "swctl alarm list", which returns fired
alarm records from the GraphQL surface.`,
	Subcommands: []*cli.Command{
		rulesCommand,
		ruleCommand,
	},
}

var rulesCommand = &cli.Command{
	Name:  "rules",
	Usage: "List the loaded alarm rules per OAP node (GET /status/alarm/rules)",
	UsageText: `List the loaded alarm rules, fanned out across every OAP node.

Examples:
1. List alarm rules:
$ swctl admin alarm rules`,
	Action: func(ctx *cli.Context) error {
		rules, err := status.AlarmRules(ctx.Context)
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleStatus, "SW_STATUS")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: rules})
	},
}

var ruleCommand = &cli.Command{
	Name:      "rule",
	Usage:     "Show one alarm rule's definition and running state (GET /status/alarm/{ruleId}[/{entityName}])",
	ArgsUsage: "<ruleId> [<entityName>]",
	UsageText: `Show the definition and running state of a single alarm rule. When an
entity name is given, the per-entity evaluation/window state is returned instead.

Examples:
1. Show a rule's running state:
$ swctl admin alarm rule service_resp_time_rule

2. Show the per-entity state of a rule:
$ swctl admin alarm rule service_resp_time_rule mock_b_service`,
	Action: func(ctx *cli.Context) error {
		args := ctx.Args()
		ruleID := args.Get(0)
		if ruleID == "" {
			return fmt.Errorf("a <ruleId> argument is required")
		}
		entityName := args.Get(1)

		var (
			result *status.ClusterAlarmStatus
			err    error
		)
		if entityName == "" {
			result, err = status.AlarmRule(ctx.Context, ruleID)
		} else {
			result, err = status.AlarmRuleEntity(ctx.Context, ruleID, entityName)
		}
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleStatus, "SW_STATUS")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}
