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

package estimate

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

var ScaleCommand = &cli.Command{
	Name:  "scale",
	Usage: `estimate monitored process scale of the given id or name in service and labels`,
	UsageText: `This command estimate monitored process scale, via id or name in service and labels.

Examples:
1. Estimate process scale by service name "abc" with labels "t1,t2":
$ swctl process estimate scale --service-name abc --labels t1,t2`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:  "labels",
				Usage: "the `labels` by which labels of the process, multiple labels split by ',': l1,l2",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		labelString := ctx.String("labels")
		labels := make([]string, 0)
		if labelString != "" {
			labels = strings.Split(labelString, ",")
		}

		scale, err := metadata.EstimateProcessScale(ctx, serviceID, labels)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: scale})
	},
}
