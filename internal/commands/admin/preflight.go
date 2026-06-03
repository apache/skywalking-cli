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

package admin

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var preflightCommand = &cli.Command{
	Name:  "preflight",
	Usage: "Detect which admin feature modules are enabled on the OAP admin-server",
	UsageText: `Reads the effective configuration from the admin host and reports which feature
modules (status, inspect, ui-management, dsl-debugging, runtime-rule) are enabled.

Examples:
1. Check admin feature availability:
$ swctl admin preflight`,
	Action: func(ctx *cli.Context) error {
		// Run reports per-module enablement; a transport error means the admin host
		// is unreachable, in which case we still print the (all-disabled) result so
		// the user sees which admin URL was probed.
		result, err := preflight.Run(ctx.Context)
		if err != nil && !result.AdminReachable {
			return preflight.Explain(ctx.Context, err, preflight.ModuleAdminServer, "SW_ADMIN_SERVER")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}
