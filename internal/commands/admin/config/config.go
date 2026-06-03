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

package config

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/admin/status"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "config",
	Usage: "Inspect the OAP effective configuration and TTL from the admin-server",
	Subcommands: []*cli.Command{
		dumpCommand,
		ttlCommand,
	},
}

var dumpCommand = &cli.Command{
	Name:  "dump",
	Usage: "Dump the OAP node's effective, secrets-redacted configuration (GET /debugging/config/dump)",
	UsageText: `Dump the effective configuration of the OAP node as a flat map of
"<module>.<provider>.<property>" keys. Secrets are redacted by OAP.

Examples:
1. Dump the effective configuration:
$ swctl admin config dump`,
	Action: func(ctx *cli.Context) error {
		dump, err := status.ConfigDump(ctx.Context)
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleStatus, "SW_STATUS")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: dump})
	},
}

var ttlCommand = &cli.Command{
	Name:  "ttl",
	Usage: "Show the effective TTL configuration (GET /status/config/ttl)",
	UsageText: `Show the effective metric / record / trace / log TTL bounds.

Examples:
1. Show the effective TTL configuration:
$ swctl admin config ttl`,
	Action: func(ctx *cli.Context) error {
		ttl, err := status.ConfigTTL(ctx.Context)
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleStatus, "SW_STATUS")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: ttl})
	},
}
