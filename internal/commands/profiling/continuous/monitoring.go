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

package continuous

import (
	"fmt"
	"strings"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

var Monitoring = &cli.Command{
	Name:    "monitoring",
	Aliases: []string{"monitor"},
	Usage:   `query all continuous profiling monitoring instances through service and policy target`,
	UsageText: `Query all continuous profiling monitoring instances through service and policy target.

Example:
1. Query all continuous profiling monitoring instances through service "business-zone::projectC" and policy target "ON_CPU"
$ swctl profiling continuous monitoring --service-name=business-zone::projectC --target=ON_CPU
`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "target",
				Usage:    "policy target type",
				Required: true,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")
		targetString := ctx.String("target")
		targetHasSet := false
		var target api.ContinuousProfilingTargetType
		for _, enum := range api.AllContinuousProfilingTargetType {
			if strings.EqualFold(enum.String(), targetString) {
				target = enum
				targetHasSet = true
				break
			}
		}
		if !targetHasSet {
			return fmt.Errorf("unknown target type: %s", targetString)
		}

		instances, err := profiling.QueryContinuousProfilingMonitoringInstances(ctx.Context, serviceID, target)
		if err != nil {
			return err
		}

		return display.Display(ctx.Context, &displayable.Displayable{Data: instances})
	},
}
