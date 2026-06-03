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

package dsldebug

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/dsldebug"
	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "dsl-debug",
	Usage: "Live MAL / LAL / OAL debugger via the admin-server `dsl-debugging` module",
	UsageText: `Run sample-based debug sessions that capture how MAL / LAL / OAL rules transform
live ingest, with per-stage captured records.`,
	Subcommands: []*cli.Command{
		statusCommand,
		sessionsCommand,
		sessionCommand,
	},
}

var statusCommand = &cli.Command{
	Name:  "status",
	Usage: "Show the dsl-debugging module health snapshot (GET /dsl-debugging/status)",
	Action: func(ctx *cli.Context) error {
		result, err := dsldebug.Status(ctx.Context)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var sessionsCommand = &cli.Command{
	Name:  "sessions",
	Usage: "List the active debug sessions (GET /dsl-debugging/sessions)",
	Action: func(ctx *cli.Context) error {
		result, err := dsldebug.ListSessions(ctx.Context)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var sessionCommand = &cli.Command{
	Name:  "session",
	Usage: "Start / poll / stop a single debug session",
	Subcommands: []*cli.Command{
		sessionStartCommand,
		sessionGetCommand,
		sessionStopCommand,
	},
}

var sessionStartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start a debug capture session (POST /dsl-debugging/session)",
	UsageText: `Start a sample-based debug capture session.

Examples:
1. Debug a MAL metric:
$ swctl admin dsl-debug session start --catalog otel-rules --name vm --rule-name vm_memory_used`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "catalog",
			Usage:    "session `catalog`: otel-rules / log-mal-rules / telegraf-rules / lal / oal",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "rule file name (MAL/LAL) or OAL source class `name`",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "rule-name",
			Usage:    "the metric / rule `name` within the file (OAL: same as --name)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "client-id",
			Usage: "stable per-debug-context `id`; a random UUID is generated when omitted",
		},
		&cli.StringFlag{
			Name:  "granularity",
			Usage: "LAL capture `granularity`: block (default) or statement",
		},
		&cli.IntFlag{
			Name:  "record-cap",
			Usage: fmt.Sprintf("max records to capture before the session is full (1-%d)", dsldebug.MaxRecordCap),
		},
		&cli.IntFlag{
			Name:  "retention-millis",
			Usage: fmt.Sprintf("session wall-clock retention in ms (max %d = 1h)", dsldebug.MaxRetentionMillis),
		},
	},
	Action: func(ctx *cli.Context) error {
		recordCap := ctx.Int("record-cap")
		if recordCap < 0 || recordCap > dsldebug.MaxRecordCap {
			return fmt.Errorf("--record-cap must be between 1 and %d", dsldebug.MaxRecordCap)
		}
		retention := ctx.Int("retention-millis")
		if retention < 0 || retention > dsldebug.MaxRetentionMillis {
			return fmt.Errorf("--retention-millis must be between 1 and %d", dsldebug.MaxRetentionMillis)
		}
		clientID := ctx.String("client-id")
		if clientID == "" {
			clientID = uuid.New().String()
		}

		result, err := dsldebug.StartSession(ctx.Context, &dsldebug.StartArgs{
			ClientID:        clientID,
			Catalog:         ctx.String("catalog"),
			Name:            ctx.String("name"),
			RuleName:        ctx.String("rule-name"),
			Granularity:     ctx.String("granularity"),
			RecordCap:       recordCap,
			RetentionMillis: retention,
		})
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var sessionGetCommand = &cli.Command{
	Name:      "get",
	Usage:     "Poll a session's captured records (GET /dsl-debugging/session/{id})",
	ArgsUsage: "<sessionId>",
	Action: func(ctx *cli.Context) error {
		id := ctx.Args().Get(0)
		if id == "" {
			return fmt.Errorf("a <sessionId> argument is required")
		}
		result, err := dsldebug.GetSession(ctx.Context, id)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

var sessionStopCommand = &cli.Command{
	Name:      "stop",
	Usage:     "Stop a session (POST /dsl-debugging/session/{id}/stop)",
	ArgsUsage: "<sessionId>",
	Action: func(ctx *cli.Context) error {
		id := ctx.Args().Get(0)
		if id == "" {
			return fmt.Errorf("a <sessionId> argument is required")
		}
		result, err := dsldebug.StopSession(ctx.Context, id)
		if err != nil {
			return explain(ctx, err)
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: result})
	},
}

func explain(ctx *cli.Context, err error) error {
	return preflight.Explain(ctx.Context, err, preflight.ModuleDSLDebug, "SW_DSL_DEBUGGING")
}
