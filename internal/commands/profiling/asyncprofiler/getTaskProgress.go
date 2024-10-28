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

package asyncprofiler

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var getTaskProgressCommand = &cli.Command{
	Name:    "progress",
	Aliases: []string{"p"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "task-id",
			Usage:    "async profiler task id.",
			Required: true,
		},
	},
	Usage: "Query async-profiler task progress",
	UsageText: `Query async-profiler task progress

Examples:
1. Query task progress, including task logs and successInstances and errorInstances
$ swctl profiling async progress --task-id=task-id`,
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")

		data, err := profiling.GetAsyncProfilerTaskProgress(ctx, taskID)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: data, Condition: taskID})
	},
}
