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

package hierarchy

import (
	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/hierarchy"

	"github.com/urfave/cli/v2"
)

var serviceCommand = &cli.Command{
	Name:    "service",
	Aliases: []string{"svc"},
	Usage:   "Query the hierarchy of given service",
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "layer",
				Usage:    "Name of the layer to query the hierarchy of this layer",
				Required: true,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		layer := ctx.String("layer")

		hierarchy, err := hierarchy.ServiceHierarchy(ctx, serviceID, layer)
		if err != nil {
			return err
		}
		return display.Display(ctx, &displayable.Displayable{Data: hierarchy, Condition: serviceID})
	},
}
