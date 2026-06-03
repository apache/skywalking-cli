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

package uitemplate

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/admin/uitemplate"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/util"
)

var Command = &cli.Command{
	Name:  "ui-template",
	Usage: "Manage dashboard templates via the admin-server `ui-management` module",
	UsageText: `Manage OAP dashboard templates over REST. This replaces the GraphQL
UIConfigurationManagement template mutations retired in SkyWalking 11.0.0. There is no
delete; templates are soft-disabled.`,
	Subcommands: []*cli.Command{
		listCommand,
		getCommand,
		createCommand,
		updateCommand,
		disableCommand,
	},
}

var listCommand = &cli.Command{
	Name:  "list",
	Usage: "List all dashboard templates (GET /ui-management/templates)",
	UsageText: `List all dashboard templates.

Examples:
1. List enabled templates:
$ swctl admin ui-template list

2. Include soft-disabled templates:
$ swctl admin ui-template list --include-disabled`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "include-disabled",
			Usage: "also return soft-disabled templates",
		},
	},
	Action: func(ctx *cli.Context) error {
		templates, err := uitemplate.List(ctx.Context, ctx.Bool("include-disabled"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: templates})
	},
}

var getCommand = &cli.Command{
	Name:      "get",
	Usage:     "Get a single dashboard template by ID (GET /ui-management/templates/{id})",
	ArgsUsage: "<id>",
	Action: func(ctx *cli.Context) error {
		id := ctx.Args().Get(0)
		if id == "" {
			return fmt.Errorf("an <id> argument is required")
		}
		template, err := uitemplate.Get(ctx.Context, id)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: template})
	},
}

var createCommand = &cli.Command{
	Name:  "create",
	Usage: "Add a new dashboard template (POST /ui-management/templates)",
	UsageText: `Add a new dashboard template. The file holds the JSON-encoded template
configuration. An id is required by OAP; a random UUID is generated when --id is omitted.

Examples:
1. Create a template from a file:
$ swctl admin ui-template create -f my-dashboard.json`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Usage:    "`path` to the JSON-encoded template configuration",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "id",
			Usage: "template `id`; a random UUID is generated when omitted",
		},
	},
	Action: func(ctx *cli.Context) error {
		configuration, err := readFile(ctx.String("file"))
		if err != nil {
			return err
		}
		id := ctx.String("id")
		if id == "" {
			id = uuid.New().String()
		}
		status, err := uitemplate.Create(ctx.Context, id, configuration)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: status})
	},
}

var updateCommand = &cli.Command{
	Name:  "update",
	Usage: "Update an existing dashboard template (PUT /ui-management/templates)",
	UsageText: `Update an existing dashboard template by ID with a new configuration.

Examples:
1. Update a template:
$ swctl admin ui-template update --id <uuid> -f my-dashboard.json`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "`id` of the template to update",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Usage:    "`path` to the JSON-encoded template configuration",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		configuration, err := readFile(ctx.String("file"))
		if err != nil {
			return err
		}
		status, err := uitemplate.Update(ctx.Context, ctx.String("id"), configuration)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: status})
	},
}

var disableCommand = &cli.Command{
	Name:      "disable",
	Usage:     "Soft-disable a dashboard template (POST /ui-management/templates/{id}/disable)",
	ArgsUsage: "<id>",
	Action: func(ctx *cli.Context) error {
		id := ctx.Args().Get(0)
		if id == "" {
			return fmt.Errorf("an <id> argument is required")
		}
		status, err := uitemplate.Disable(ctx.Context, id)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: status})
	},
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(util.ExpandFilePath(path))
	if err != nil {
		return "", fmt.Errorf("failed to read template file %q: %w", path, err)
	}
	return strings.TrimSpace(string(content)), nil
}

func explain(ctx *cli.Context, err error) error {
	return preflight.Explain(ctx.Context, err, preflight.ModuleUIManage, "SW_UI_MANAGEMENT")
}
