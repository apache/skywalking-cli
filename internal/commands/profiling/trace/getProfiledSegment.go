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
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"
)

var getProfiledSegmentCommand = &cli.Command{
	Name:    "profiled-segment",
	Aliases: []string{"ps"},
	Usage:   "Analyze profiled segment with time ranges",
	UsageText: `Analyze profiled segment with time ranges.

Examples:
1. Analyze profiled segment with time ranges.
$ swctl profiling trace profiled-segment --segment-id=<segment-id>`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "segment-id",
			Usage: "profiled segment id.",
		},
	},
	Action: func(ctx *cli.Context) error {
		segmentID := ctx.String("segment-id")
		segment, err := profiling.GetTraceProfilingSegment(ctx, segmentID)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: segment, Condition: segmentID})
	},
}
