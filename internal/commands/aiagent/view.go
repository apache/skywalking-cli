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

package aiagent

import (
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/aiagent/view"
)

var viewCommand = &cli.Command{
	Name:  "view",
	Usage: "Read a whole conversation as one asz.view document",
	UsageText: `Read the whole conversation, once, as one asz.view 1.0 document from the OAP's route
GET /ai-agent/conversations/{conversation}/v1/view, on the "--base-url" host. The body
is streamed to stdout, or to "--output", as it arrives: JSON, or YAML with "--yaml".
The "--display" option does not apply; the document is printed as the OAP sends it.

Examples:
1. A conversation as JSON, into a file:
$ swctl ai-agent view --service-name "Claude Code" --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 --output conversation.json

2. As YAML, on the terminal:
$ swctl ai-agent view --service-name "Claude Code" --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 --yaml`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "conversation",
				Usage:    "`id` of the conversation",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "yaml",
				Usage: "ask for the document as YAML instead of JSON",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "write the document to this `file` instead of stdout",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
		interceptor.ParseInstance(false),
	),
	Action: func(ctx *cli.Context) error {
		var out io.Writer = os.Stdout
		if path := ctx.String("output"); path != "" {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}
		_, err := view.Fetch(ctx.Context, ctx.String("conversation"), ctx.String("service-name"), ctx.String("instance-name"), ctx.Bool("yaml"), out)
		return err
	},
}
