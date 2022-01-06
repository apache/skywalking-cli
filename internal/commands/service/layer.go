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

package service

import (
	"fmt"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var LayerCommand = &cli.Command{
	Name:      "layer",
	Aliases:   []string{"ly"},
	Usage:     `list the service list according to layer`,
	ArgsUsage: "<layer name>",
	UsageText: `This command lists the services matching the given "<layer name>".

Examples:
2. List services in "GENERAL" layer:
$ swctl svc ly GENERAL`,
	Action: func(ctx *cli.Context) error {
		var services []api.Service

		if args := ctx.Args(); args.Len() == 0 {
			return fmt.Errorf("layer must be provide")
		}
		services, err := metadata.ListLayerService(ctx, ctx.Args().First())
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: services})
	},
}
