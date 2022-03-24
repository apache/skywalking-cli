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
	"strconv"
	"strings"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

var getProfiledAnalyzeCommand = &cli.Command{
	Name:      "analysis",
	Aliases:   []string{"pa"},
	Usage:     "Analyze profiled segment.",
	ArgsUsage: "[parameters...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "segment-id",
			Usage: "profiled segment id.",
		},
		&cli.StringFlag{
			Name:  "time-ranges",
			Usage: "need to analyze time ranges in the segment: start-end,start-end",
		},
	},
	Action: func(ctx *cli.Context) error {
		segmentID := ctx.String("segment-id")

		tagStr := ctx.String("time-ranges")
		var timeRanges []*api.ProfileAnalyzeTimeRange = nil
		if tagStr != "" {
			tagArr := strings.Split(tagStr, ",")
			for _, tag := range tagArr {
				kv := strings.Split(tag, "-")
				start, err := strconv.ParseInt(kv[0], 10, 64)
				if err != nil {
					return err
				}
				end, err := strconv.ParseInt(kv[1], 10, 64)
				if err != nil {
					return err
				}
				timeRanges = append(timeRanges, &api.ProfileAnalyzeTimeRange{Start: start, End: end})
			}
		}

		analysis, err := profiling.GetTraceProfilingAnalyze(ctx, segmentID, timeRanges)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: analysis, Condition: segmentID})
	},
}
