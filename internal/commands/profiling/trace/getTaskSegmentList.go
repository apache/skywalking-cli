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

var getTaskSegmentListCommand = &cli.Command{
	Name:    "segment-list",
	Aliases: []string{"sl"},
	Usage:   "Query profiling trace task segments list",
	UsageText: `Query profiling trace task segments list

Examples:
1. Query profiled segments list
$ swctl profiling trace segment-list --service-name=service-name --endpoint-name=endpoint`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "task-id",
			Usage: "`<task id>` whose profiled segment are to be searched",
		},
	},
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")
		segmentList, err := profiling.GetTraceProfilingTaskSegmentList(ctx.Context, taskID)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: segmentList, Condition: taskID})
	},
}
