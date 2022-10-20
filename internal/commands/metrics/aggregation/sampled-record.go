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

package aggregation

import (
	"fmt"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"
)

var SampledRecords = &cli.Command{
	Name:      "sampled-record",
	Usage:     "query the top <n> entities sorted by the specified records",
	ArgsUsage: "<n>",
	UsageText: `Query the top <n> entities sorted by the specified records.

Examples:
1. Query the top 5 database statements whose execute duration are largest:
$ swctl metrics sampled-record --name top_n_database_statement 5
`,
	Flags: flags.Flags(
		flags.DurationFlags,
		flags.MetricsFlags,
		flags.InstanceRelationFlags,
		flags.EndpointRelationFlags,
		flags.ProcessRelationFlags,
		[]cli.Flag{
			&cli.GenericFlag{
				Name:  "order",
				Usage: "the `order` by which the top entities are sorted",
				Value: &model.OrderEnumValue{
					Enum:     api.AllOrder,
					Default:  api.OrderDes,
					Selected: api.OrderDes,
				},
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
		interceptor.ParseEndpointRelation(false),
		interceptor.ParseInstanceRelation(false),
		interceptor.ParseProcessRelation(false),
	),
	Action: func(ctx *cli.Context) error {
		// read OAP version
		major, minor, err := metadata.BackendVersion(ctx)
		if err != nil {
			return fmt.Errorf("read backend version failure: %v", err)
		}

		// since 9.3.0, use new record query API
		if major >= 9 && minor >= 3 {
			condition, duration, err1 := buildReadRecordsCondition(ctx)
			if err1 != nil {
				return err1
			}
			logger.Log.Debugln(condition.Name, condition.TopN)

			records, err1 := metrics.ReadRecords(ctx, *condition, *duration)
			if err1 != nil {
				return err1
			}

			return display.Display(ctx, &displayable.Displayable{Data: records})
		}

		condition, duration, err := buildSortedCondition(ctx, false)
		if err != nil {
			return err
		}

		logger.Log.Debugln(condition.Name, condition.Scope, condition.TopN)
		sampledRecords, err := metrics.SampledRecords(ctx, *condition, *duration)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: sampledRecords})
	},
}
