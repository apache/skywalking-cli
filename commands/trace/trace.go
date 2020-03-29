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

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/display"
	"github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/graphql/trace"
)

var Command = cli.Command{
	Name:      "trace",
	ShortName: "t",
	Usage:     "trace related sub-command",
	ArgsUsage: "trace id",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return fmt.Errorf("command trace without sub command requires 1 trace id as argument")
		}

		trace := trace.Trace(ctx, ctx.Args().First())

		return display.Display(ctx, &displayable.Displayable{Data: trace})
	},
	Subcommands: cli.Commands{
		ListCommand,
	},
}
