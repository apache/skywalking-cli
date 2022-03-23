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

package profiling

import (
	"strconv"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	ebpf_graphql "github.com/apache/skywalking-cli/pkg/graphql/ebpf"
)

var AnalyzationCommand = &cli.Command{
	Name:    "analysis",
	Aliases: []string{"as"},
	Usage:   `analyze ebpf profiling task`,
	UsageText: `This command analysis profiling task, via id of task and time ranges.

Example:
1. Analysis profiling tasks of task id "abc" and time range in 1648020042869 to 1648020100764.
$ swctl ebpf-profiling analysis --task-id=abc --time-ranges=1648020042869-1648020100764
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "the `task-id` by which task are scheduled",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "time-ranges",
				Usage: "need to analyze time ranges in the segment: start-end,start-end",
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")

		tagStr := ctx.String("time-ranges")
		var timeRanges []*api.EBPFProfilingAnalyzeTimeRange = nil
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
				timeRanges = append(timeRanges, &api.EBPFProfilingAnalyzeTimeRange{Start: start, End: end})
			}
		}

		analyzation, err := ebpf_graphql.QueryEBPFProfilingAnalyzation(ctx, taskID, timeRanges)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: analyzation})
	},
}
