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

package process

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var GetCommand = &cli.Command{
	Name:  "get",
	Usage: `get monitored process of the given <process-id>`,
	UsageText: `This command get single process, via process-id.

Examples:
1. get single process by process id "2b9e46c13c91803695a4364257415e523af7cbf17bf4058e025c16b944a6a85b":
$ swctl process get 2b9e46c13c91803695a4364257415e523af7cbf17bf4058e025c16b944a6a85b`,
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() == 0 {
			return fmt.Errorf("process-id must be provided")
		}

		instance, err := metadata.GetProcess(ctx, ctx.Args().First())
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: instance})
	},
}
