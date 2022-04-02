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

package endpoint

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var GetCommand = &cli.Command{
	Name:  "get",
	Usage: `get monitored endpoint of the given <endpoint-id>`,
	UsageText: `This command get single endpoint, via endpoint-id.

Examples:
1. get single endpoint by endpoint id "cHJvdmlkZXI=.1_L3VzZXJz":
$ swctl endpoint get cHJvdmlkZXI=.1_L3VzZXJz`,
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() == 0 {
			return fmt.Errorf("endpoint-id must be provide")
		}

		endpointInfo, err := metadata.GetEndpointInfo(ctx, ctx.Args().First())
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: endpointInfo})
	},
}