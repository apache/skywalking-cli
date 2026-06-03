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

package runtimerule

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/admin/runtimerule"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/util"
)

var Command = &cli.Command{
	Name:  "runtime-rule",
	Usage: "Hot-update MAL / LAL rules via the admin-server `receiver-runtime-rule` module",
	UsageText: `Add, override, inactivate and delete MAL / LAL rule files at runtime without
restarting OAP, and inspect the live and bundled rule state.

Catalogs: otel-rules, log-mal-rules, telegraf-rules, lal.`,
	Subcommands: []*cli.Command{
		listCommand,
		bundledCommand,
		getCommand,
		addCommand,
		inactivateCommand,
		deleteCommand,
		dumpCommand,
	},
}

func catalogFlag(required bool) cli.Flag {
	return &cli.StringFlag{
		Name:     "catalog",
		Usage:    "rule `catalog`: otel-rules / log-mal-rules / telegraf-rules / lal",
		Required: required,
	}
}

func nameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "name",
		Usage:    "rule `name`",
		Required: true,
	}
}

var listCommand = &cli.Command{
	Name:  "list",
	Usage: "List the live rule state per node (GET /runtime/rule/list)",
	Flags: []cli.Flag{catalogFlag(false)},
	Action: func(ctx *cli.Context) error {
		result, err := runtimerule.List(ctx.Context, ctx.String("catalog"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var bundledCommand = &cli.Command{
	Name:  "bundled",
	Usage: "List the static (bundled) rule twins for a catalog (GET /runtime/rule/bundled)",
	Flags: []cli.Flag{
		catalogFlag(true),
		&cli.BoolFlag{
			Name:  "with-content",
			Usage: "include the bundled rule content",
			Value: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		result, err := runtimerule.ListBundled(ctx.Context, ctx.String("catalog"), ctx.Bool("with-content"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var getCommand = &cli.Command{
	Name:  "get",
	Usage: "Fetch a single rule's content + metadata (GET /runtime/rule)",
	Flags: []cli.Flag{
		catalogFlag(true),
		nameFlag(),
		&cli.StringFlag{
			Name:  "source",
			Usage: "`source` to read: runtime (default, DAO first) or bundled (static twin only)",
		},
	},
	Action: func(ctx *cli.Context) error {
		rule, err := runtimerule.Get(ctx.Context, ctx.String("catalog"), ctx.String("name"), ctx.String("source"), "")
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: rule})
	},
}

var addCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"add-or-update"},
	Usage:   "Push a new or updated rule from a YAML file (POST /runtime/rule/addOrUpdate)",
	UsageText: `Push a new or updated MAL / LAL rule. The file holds the raw rule YAML.

Examples:
1. Apply a MAL rule:
$ swctl admin runtime-rule add --catalog otel-rules --name vm -f vm.yaml`,
	Flags: []cli.Flag{
		catalogFlag(true),
		nameFlag(),
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Usage:    "`path` to the rule YAML",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "allow-storage-change",
			Usage: "approve an edit that moves the rule's storage identity (structural change)",
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "re-run apply even when the content is byte-identical",
		},
	},
	Action: func(ctx *cli.Context) error {
		body, err := readFile(ctx.String("file"))
		if err != nil {
			return err
		}
		result, err := runtimerule.AddOrUpdate(ctx.Context, ctx.String("catalog"), ctx.String("name"),
			body, ctx.Bool("allow-storage-change"), ctx.Bool("force"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var inactivateCommand = &cli.Command{
	Name:  "inactivate",
	Usage: "Turn a rule off (POST /runtime/rule/inactivate)",
	Flags: []cli.Flag{catalogFlag(true), nameFlag()},
	Action: func(ctx *cli.Context) error {
		result, err := runtimerule.Inactivate(ctx.Context, ctx.String("catalog"), ctx.String("name"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var deleteCommand = &cli.Command{
	Name:  "delete",
	Usage: "Remove a rule (POST /runtime/rule/delete)",
	UsageText: `Remove a rule. An active rule must be inactivated first (the server returns a
409 requires_inactivate_first otherwise).

Examples:
1. Delete an inactivated rule:
$ swctl admin runtime-rule delete --catalog otel-rules --name vm

2. Revert a rule to its bundled twin:
$ swctl admin runtime-rule delete --catalog otel-rules --name vm --mode revertToBundled`,
	Flags: []cli.Flag{
		catalogFlag(true),
		nameFlag(),
		&cli.StringFlag{
			Name:  "mode",
			Usage: "deletion `mode`; pass revertToBundled to restore the static twin",
		},
	},
	Action: func(ctx *cli.Context) error {
		result, err := runtimerule.Delete(ctx.Context, ctx.String("catalog"), ctx.String("name"), ctx.String("mode"))
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var dumpCommand = &cli.Command{
	Name:  "dump",
	Usage: "Download a tar.gz snapshot of all rules (GET /runtime/rule/dump)",
	UsageText: `Download a tar.gz snapshot of all rules, or one catalog's.

Examples:
1. Dump every rule to a file:
$ swctl admin runtime-rule dump -o rules.tar.gz

2. Dump one catalog:
$ swctl admin runtime-rule dump --catalog otel-rules -o otel.tar.gz`,
	Flags: []cli.Flag{
		catalogFlag(false),
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "`path` to write the tar.gz to (default: stdout)",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		data, err := runtimerule.Dump(ctx.Context, ctx.String("catalog"))
		if err != nil {
			return explain(ctx, err)
		}
		out := util.ExpandFilePath(ctx.String("output"))
		if err := os.WriteFile(out, data, 0o600); err != nil {
			return fmt.Errorf("failed to write dump to %q: %w", out, err)
		}
		fmt.Printf("wrote %d bytes to %s\n", len(data), out)
		return nil
	},
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(util.ExpandFilePath(path))
	if err != nil {
		return "", fmt.Errorf("failed to read rule file %q: %w", path, err)
	}
	// Send the rule bytes verbatim. The runtime-rule API hashes the raw body for its
	// contentHash / no-change detection, so the CLI must not normalize whitespace —
	// trimming or re-adding a trailing newline would make a byte-identical rule look
	// like a change.
	return string(content), nil
}

func explain(ctx *cli.Context, err error) error {
	return preflight.Explain(ctx.Context, err, preflight.ModuleRuntimeRule, "SW_RECEIVER_RUNTIME_RULE")
}
