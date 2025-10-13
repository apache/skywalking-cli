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

package pprof

import (
	"strings"

	"github.com/urfave/cli/v2"
	"skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"
)

var analysisCommand = &cli.Command{
	Name:    "analysis",
	Aliases: []string{"a"},
	Usage:   "Query pprof analysis",
	UsageText: `Query pprof analysis

Examples:
1. Query the flame graph produced by pprof
$ swctl profiling pprof analysis --service-name=service-name --task-id=task-id \
	--instance-name-list=instance-name1,instance-name2`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceListFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Usage:    "pprof task id",
				Required: true,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseInstanceList(true),
	),
	Action: func(ctx *cli.Context) error {
		taskID := ctx.String("task-id")
		instances := strings.Split(ctx.String("instance-id-list"), ",")

		request := &query.PprofAnalyzationRequest{
			TaskID:      taskID,
			InstanceIds: instances,
		}

		analyze, err := profiling.GetPprofAnalyze(ctx.Context, request)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: analyze, Condition: request})
	},
}
