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

package dependency

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/dependency"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"

	api "skywalking.apache.org/repo/goapi/query"
)

var GlobalCommand = &cli.Command{
	Name:    "global",
	Aliases: []string{"glb"},
	Usage:   "Query the global dependencies",
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "layer",
				Usage:    "Name of the layer to query the dependency of this layer",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		layer := ctx.String("layer")
		end := ctx.String("end")
		start := ctx.String("start")
		step := ctx.Generic("step")

		duration := api.Duration{
			Start: start,
			End:   end,
			Step:  step.(*model.StepEnumValue).Selected,
		}

		major, _, err := metadata.BackendVersion(ctx)
		if err != nil {
			return err
		}

		var topology api.Topology
		if major >= 10 {
			topology, err = dependency.GlobalTopology(ctx, layer, duration)
		} else {
			if layer != "" {
				return fmt.Errorf("the layer parameter only available when OAP version >= 10.0.0")
			}
			topology, err = dependency.GlobalTopologyWithoutLayer(ctx, duration)
		}

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: topology})
	},
}
