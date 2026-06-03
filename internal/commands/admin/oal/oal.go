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

package oal

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/oal"
	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "oal",
	Usage: "Browse the read-only OAL catalog (the OAL debugger's rule picker)",
	UsageText: `Read-only listing of loaded OAL files and the per-dispatcher source catalog used
by the OAL live debugger. Hosted by the admin-server dsl-debugging module.`,
	Subcommands: []*cli.Command{
		filesCommand,
		fileCommand,
		rulesCommand,
		ruleCommand,
	},
}

var filesCommand = &cli.Command{
	Name:  "files",
	Usage: "List the loaded .oal file names (GET /runtime/oal/files)",
	Action: func(ctx *cli.Context) error {
		result, err := oal.ListFiles(ctx.Context)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var fileCommand = &cli.Command{
	Name:      "file",
	Usage:     "Print the raw .oal text of one file (GET /runtime/oal/files/{name})",
	ArgsUsage: "<name>",
	Action: func(ctx *cli.Context) error {
		name := ctx.Args().Get(0)
		if name == "" {
			return fmt.Errorf("a <name> argument is required")
		}
		content, err := oal.GetFile(ctx.Context, name)
		if err != nil {
			return explain(ctx, err)
		}
		fmt.Println(content)
		return nil
	},
}

var rulesCommand = &cli.Command{
	Name:  "rules",
	Usage: "List the per-dispatcher OAL source catalog (GET /runtime/oal/rules)",
	Action: func(ctx *cli.Context) error {
		result, err := oal.ListSources(ctx.Context)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var ruleCommand = &cli.Command{
	Name:      "rule",
	Usage:     "Show one OAL source's per-metric holder status (GET /runtime/oal/rules/{source})",
	ArgsUsage: "<source>",
	Action: func(ctx *cli.Context) error {
		source := ctx.Args().Get(0)
		if source == "" {
			return fmt.Errorf("a <source> argument is required")
		}
		result, err := oal.GetSource(ctx.Context, source)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

func explain(ctx *cli.Context, err error) error {
	return preflight.Explain(ctx.Context, err, preflight.ModuleDSLDebug, "SW_DSL_DEBUGGING")
}
