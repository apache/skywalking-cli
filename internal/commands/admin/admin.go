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

	"github.com/apache/skywalking-cli/internal/commands/admin/alarm"
	"github.com/apache/skywalking-cli/internal/commands/admin/cluster"
	"github.com/apache/skywalking-cli/internal/commands/admin/config"
	"github.com/apache/skywalking-cli/internal/commands/admin/dsldebug"
	"github.com/apache/skywalking-cli/internal/commands/admin/inspect"
	"github.com/apache/skywalking-cli/internal/commands/admin/oal"
	"github.com/apache/skywalking-cli/internal/commands/admin/runtimerule"
	"github.com/apache/skywalking-cli/internal/commands/admin/uitemplate"
)

// Command is the parent of every sub-command that talks to the OAP admin-server
// REST host (default port 17128), as opposed to the public GraphQL surface on
// `--base-url` (default port 12800). The admin host bundles the status, inspect,
// ui-management, dsl-debugging and runtime-rule feature modules. Its address comes
// from `--admin-url` (or is derived from `--base-url` with port 17128).
var Command = &cli.Command{
	Name:  "admin",
	Usage: "Admin (REST) sub-commands that talk to the OAP admin-server (default port 17128)",
	UsageText: `Admin sub-commands call the OAP admin-server REST host, a separate surface from the
public GraphQL endpoint used by the other commands.

The admin host address defaults to the "--base-url" host with port 17128; override it
with the global "--admin-url" flag (or the SW_ADMIN_URL env var / "admin-url" config key).
The admin host has no built-in authentication and is expected to sit behind a gateway;
"--username"/"--password"/"--authorization" and "--insecure" apply to it the same way.`,
	Subcommands: []*cli.Command{
		preflightCommand,
		cluster.Command,
		config.Command,
		alarm.Command,
		inspect.Command,
		uitemplate.Command,
		runtimerule.Command,
		dsldebug.Command,
		oal.Command,
	},
}
