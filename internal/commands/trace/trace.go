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

package trace

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/trace"
)

var Command = &cli.Command{
	Name:      "trace",
	Aliases:   []string{"t"},
	Usage:     "Trace related sub-command",
	ArgsUsage: "<trace id>",
	UsageText: `This command can be used to query the details of a trace,
and its sub-command "ls" can be used to list all or part of the traces
with specified options, like service name, endpoint name, etc.

Examples:
1. Query the trace details (spans) of id "321661b1-9a31-4e12-ad64-c8f6711f108d":
$ swctl trace "321661b1-9a31-4e12-ad64-c8f6711f108d"`,
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return fmt.Errorf("command trace without sub command requires 1 trace id as argument")
		}

		trace, err := trace.Trace(ctx, ctx.Args().First())

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: trace})
	},
	Subcommands: cli.Commands{
		ListCommand,
	},
}
