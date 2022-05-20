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

package ebpf

import (
	"strconv"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model/ebpf"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var AnalyzationCommand = &cli.Command{
	Name:    "analysis",
	Aliases: []string{"as"},
	Usage:   `analyze ebpf profiling task`,
	UsageText: `This command analysis profiling task, via id of task and time ranges.

Example:
1. Analysis profiling tasks of task id "abc" and time range in 1648020042869 to 1648020100764.
$ swctl profiling ebpf analysis --schedule-id=abc --time-ranges=1648020042869-1648020100764
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "schedule-id",
				Usage:    "the `schedule-id` list by which task are scheduled, multiple schedule split by ',': schedule-id,schedule-id",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "time-ranges",
				Usage: "need to analyze time ranges in the segment: start-end,start-end",
			},
			&cli.GenericFlag{
				Name:  "aggregate",
				Usage: "the aggregate type for the profiling data anlysis",
				Value: &ebpf.ProfilingAnalyzeAggregateTypeEnumValue{
					Enum:     api.AllEBPFProfilingAnalyzeAggregateType,
					Default:  api.EBPFProfilingAnalyzeAggregateTypeCount,
					Selected: api.EBPFProfilingAnalyzeAggregateTypeCount,
				},
			},
		},
	),
	Action: func(ctx *cli.Context) error {
		scheduleIDListStr := ctx.String("schedule-id")
		scheduleIDList := strings.Split(scheduleIDListStr, ",")

		timeRangeStr := ctx.String("time-ranges")
		var timeRanges []*api.EBPFProfilingAnalyzeTimeRange = nil
		if timeRangeStr != "" {
			tagArr := strings.Split(timeRangeStr, ",")
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

		aggregateType := ctx.Generic("aggregate").(*ebpf.ProfilingAnalyzeAggregateTypeEnumValue).Selected
		analyzation, err := profiling.AnalysisEBPFProfilingResult(ctx, scheduleIDList, timeRanges, aggregateType)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: analyzation})
	},
}
