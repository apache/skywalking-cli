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

package display

import (
	"context"

	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/display"
	d "github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/urfave/cli/v2"
)

// The variable style sets the output style for the command.
var style = map[string]string{
	"dashboard global":         "graph",
	"dashboard global-metrics": "graph",
	"metrics top":              "table",
	"metrics sorted":           "table",
	"metrics linear":           "graph",
	"metrics list":             "table",
	"service list":             "table",
	"t":                        "graph",
	"trace":                    "graph",
	"ebpf analysis":            "graph",
	"trace analysis":           "graph",
}

// Display the object in the style specified in flag --display
func Display(cliCtx *cli.Context, displayable *d.Displayable) error {
	displayStyle := cliCtx.String("display")
	if displayStyle == "" {
		commandFullName := cliCtx.Command.FullName()
		if commandFullName != "" {
			displayStyle = getDisplayStyle(commandFullName)
		} else {
			for _, c := range cliCtx.Lineage() {
				if s := getDisplayStyle(c.Args().First()); s != "" {
					displayStyle = s
					break
				}
			}
		}
	}
	if displayStyle == "" {
		displayStyle = "json"
	}
	ctx := cliCtx.Context
	ctx = context.WithValue(ctx, contextkey.Display{}, displayStyle)
	return display.Display(ctx, displayable)
}

// getDisplayStyle gets the default display settings.
func getDisplayStyle(fullName string) string {
	if command, ok := style[fullName]; ok {
		return command
	}
	return ""
}
